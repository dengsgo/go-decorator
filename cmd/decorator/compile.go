package main

import (
	"bytes"
	"errors"
	"go/ast"
	"go/printer"
	"go/token"
	"log"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/dengsgo/go-decorator/cmd/logs"
)

const msgDecorPkgNotImported = "decorator used but package not imported (need add `import _ \"" + decoratorPackagePath + "\"`)"
const msgDecorPkgNotFound = "decor package is not found"
const msgCantUsedOnDecoratorFunc = `decorators cannot be used on decorators`

var packageInfo *_packageInfo

var printerCfg = &printer.Config{Tabwidth: 8, Mode: printer.SourcePos}

func compile(args []string) error {
	{
		var err error
		packageInfo, err = getPackageInfo("")
		if err != nil || packageInfo.Module.Path == "" {
			logs.Error("doesn't seem to be a Go project:", err)
		}
	}
	files := make([]string, 0, len(args))
	projectName := packageInfo.Module.Path
	logs.Debug("projectName", projectName)
	//log.Printf("TOOLEXEC_IMPORTPATH %+v\n", os.Getenv("TOOLEXEC_IMPORTPATH"))
	packageName := ""
	for i, arg := range args {
		if arg == "-p" && i+1 < len(args) {
			packageName = args[i+1]
		}
		if strings.HasPrefix(arg, "-") {
			continue
		}
		if strings.HasPrefix(arg, projectDir+string(filepath.Separator)) && strings.HasSuffix(arg, ".go") {
			files = args[i:]
			break
		}
	}

	if (packageName != "main" && !strings.HasPrefix(packageName, projectName)) || len(files) == 0 {
		return nil
	}

	logs.Debug("packageName", packageName, files, args)

	var originPath string

	fset := token.NewFileSet()
	pkg, err := parserGOFiles(fset, files...)
	if err != nil {
		logs.Error(err)
	}

	typeDecorRebuild(pkg)

	for file, f := range pkg.Files {
		logs.Debug("file Parse", file)
		//f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		//if err != nil {
		//	continue
		//}
		logs.Debug(f.Decls)
		imp := newImporter(f)

		updated := false

		visitAstDecl(f, func(fd *ast.FuncDecl) (r bool) {
			if fd.Doc == nil || fd.Doc.List == nil || len(fd.Doc.List) == 0 {
				return
			}
			originPath = file
			//log.Printf("%+v\n", fd)
			var collDecors []*decorAnnotation
			mapDecors := mapx{}
			for i := len(fd.Doc.List) - 1; i >= 0; i-- {
				doc := fd.Doc.List[i]
				if !strings.HasPrefix(doc.Text, decoratorScanFlag) {
					break
				}
				logs.Debug("HIT:", doc.Text)
				decorName, decorArgs, err := parseDecorAndParameters(doc.Text[len(decoratorScanFlag):])
				logs.Debug(decorName, decorArgs, err)
				if err != nil {
					logs.Error(err, biSymbol, friendlyIDEPosition(fset, doc.Pos()))
				}
				if !mapDecors.put(decorName, "") {
					logs.Error("cannot use the same decorator for repeated decoration\n\t",
						friendlyIDEPosition(fset, doc.Pos()))
				}
				collDecors = append(collDecors, newDecorAnnotation(doc, decorName, decorArgs))
			}
			if len(collDecors) == 0 {
				return
			}

			logs.Info("find the entry for using the decorator", friendlyIDEPosition(fset, fd.Pos()))
			logs.Debug("collDecors", collDecors)
			gi := newGenIdentId()
			for _, da := range collDecors {
				logs.Debug("handler:", da.doc.Text)
				// 检查 decorName 是不是装饰器
				//if fd.Recv != nil {
				//	logs.Error("decorators cannot decorate struct method", biSymbol,
				//		friendlyIDEPosition(fset, fd.Recv.Pos()))
				//	continue
				//}
				decorName, decorParams := da.name, da.parameters
				logs.Debug(decorName, decorParams)
				// check self is not decorator function
				pkgDecorName, ok := imp.importedPath(decoratorPackagePath)
				if !ok {
					logs.Error(msgDecorPkgNotImported, biSymbol,
						friendlyIDEPosition(fset, da.doc.Pos()))
				} else if pkgDecorName == "_" {
					imp.pathObjMap[decoratorPackagePath].Name = nil // rewrite this package import way
					imp.pathMap[decoratorPackagePath] = "decor"     // mark finished
					pkgDecorName = "decor"
				}

				if funIsDecorator(fd, pkgDecorName) {
					logs.Error(msgCantUsedOnDecoratorFunc, biSymbol,
						friendlyIDEPosition(fset, fd.Pos()))
				}
				// got package path
				decorPkgPath := ""
				if x := decorX(decorName); x != "" {
					if xPath, ok := imp.importedName(x); ok {
						name, _ := imp.importedPath(xPath)
						if name == "_" {
							imp.pathObjMap[xPath].Name = nil
							imp.pathMap[xPath] = x
						}
						decorPkgPath = xPath
					} else {
						logs.Error(x, "package not found", biSymbol,
							friendlyIDEPosition(fset, da.doc.Pos()))
					}
				}
				params, err := checkDecorAndGetParam(decorPkgPath, decorName, decorParams)
				if err != nil {
					logs.Error(err, biSymbol, friendlyIDEPosition(fset, da.doc.Pos()))
				}
				ra := builderReplaceArgs(fd, decorName, params, gi)
				rs, err := replace(ra)
				if err != nil {
					logs.Error(err)
				}
				genStmts, _, err := getStmtList(rs)
				if err != nil {
					logs.Error("getStmtList err", err)
				}

				if len(ra.OutArgNames) == 0 {
					// non-return
					genStmts[1].(*ast.AssignStmt).Rhs[0].(*ast.FuncLit).Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr).Fun.(*ast.FuncLit).Body.List = fd.Body.List
				} else {
					// has return
					genStmts[1].(*ast.AssignStmt).Rhs[0].(*ast.FuncLit).Body.List[0].(*ast.AssignStmt).Rhs[0].(*ast.CallExpr).Fun.(*ast.FuncLit).Body.List = fd.Body.List
				}
				ce := genStmts[2].(*ast.ExprStmt).X.(*ast.CallExpr)
				assignCorrectPos(da.doc, ce)

				fd.Body.List = genStmts
				//x.Body.Rbrace = x.Body.Lbrace + token.Pos(ofs)
				//log.Printf("fd.Body.Pos() %+v\n", fd.Body.Pos())
				updated = true
			}
			return
		},
		)

		if !updated {
			continue
		}

		var output []byte
		buffer := bytes.NewBuffer(output)
		err = printerCfg.Fprint(buffer, fset, f)
		if err != nil {
			return errors.New("fprint original code")
		}
		tgDir := path.Join(tempDir, os.Getenv("TOOLEXEC_IMPORTPATH"))
		_ = os.MkdirAll(tgDir, 0777)
		tmpEntryFile := path.Join(tgDir, filepath.Base(originPath))
		logs.Debug("originPath", originPath, filepath.Base(originPath))
		err = os.WriteFile(tmpEntryFile, buffer.Bytes(), 0777)
		if err != nil {
			logs.Error("fail write into temporary file", err.Error())
		}
		// update go build args
		for i := range args {
			if args[i] == originPath {
				args[i] = tmpEntryFile
			}
		}
		logs.Debug("args updated", args)
		logs.Debug("rewrite file", originPath, "=>", tmpEntryFile)
	}

	return nil
}

func decorX(decorName string) string {
	arr := strings.Split(decorName, ".")
	if len(arr) != 2 {
		return ""
	}
	return arr[0]
}

func visitAstDecl(f *ast.File, funVisitor func(*ast.FuncDecl) bool) {
	if f.Decls == nil || funVisitor == nil {
		return
	}
LOOP:
	for _, t := range f.Decls {
		if t == nil {
			continue
		}
		switch decl := t.(type) {
		case *ast.FuncDecl:
			if funVisitor(decl) {
				break LOOP
			}
		}
	}
}

// Reset the line of the behavior annotation where the decorator call is located
func assignCorrectPos(doc *ast.Comment, ce *ast.CallExpr) {
	ce.Lparen = doc.Pos()
	ce.Rparen = doc.Pos()
	offset := token.Pos(0)
	if t, ok := ce.Fun.(*ast.Ident); ok {
		t.NamePos = doc.Pos()
		offset = token.Pos(len(t.Name))
	} else if t, ok := ce.Fun.(*ast.SelectorExpr); ok {
		if id, ok := t.X.(*ast.Ident); ok {
			id.NamePos = doc.Pos()
			offset = token.Pos(len(id.Name))
		}
		//t.Sel.NamePos = doc.Pos() + offset + 1
		t.Sel.NamePos = doc.Pos()
		offset += token.Pos(len(t.Sel.Name)) + 1
	}
	for _, arg := range ce.Args {
		//ast.Print(token.NewFileSet(), arg)
		//if id, ok := arg.(*ast.Ident); ok {
		//	//id.NamePos = doc.Pos() + offset
		//	id.NamePos = doc.Pos()
		//}
		switch arg := arg.(type) {
		case *ast.Ident:
			arg.NamePos = doc.Pos()
		case *ast.BasicLit:
			arg.ValuePos = doc.Pos()
		case *ast.UnaryExpr:
			arg.OpPos = doc.Pos()
			if a, ok := arg.X.(*ast.Ident); ok {
				a.NamePos = doc.Pos()
			}
		}
	}
}

func reverseSlice[T any](ele []T) []T {
	r := make([]T, len(ele))
	for i, v := range ele {
		r[len(ele)-1-i] = v
	}
	return r
}

func typeDeclVisitor(decls []ast.Decl, fn func(*ast.TypeSpec, *ast.CommentGroup)) {
	if decls == nil || len(decls) == 0 {
		return
	}
	for _, decl := range decls {
		gd, ok := decl.(*ast.GenDecl)
		if !ok || gd.Specs == nil || len(gd.Specs) == 0 {
			continue
		}
		for _, spec := range gd.Specs {
			ts, ok := spec.(*ast.TypeSpec)
			if !ok {
				continue
			}
			fn(ts, gd.Doc)
		}
	}
}

func typeDecorRebuild(pkg *ast.Package) {
	findAndCollDecorComments := func(cg *ast.CommentGroup) []*ast.Comment {
		comments := make([]*ast.Comment, 0)
		if cg == nil || cg.List == nil {
			return comments
		}
		for i := len(cg.List) - 1; i >= 0; i-- {
			if !strings.HasPrefix(cg.List[i].Text, decoratorScanFlag) {
				break
			}
			comments = append(comments, cg.List[i])
		}
		return reverseSlice(comments)
	}
	typeNameMapDecorComments := map[string][]*ast.Comment{}
	errs := []error{}
	for _, f := range pkg.Files {
		typeDeclVisitor(f.Decls, func(spec *ast.TypeSpec, typeDoc *ast.CommentGroup) {
			if (spec.Doc == nil || spec.Doc.List == nil) &&
				(typeDoc == nil || typeDoc.List == nil) {
				return
			}
			comments := findAndCollDecorComments(spec.Doc)
			log.Printf("findAndCollDecorComments(spec.Doc): %+v \n", comments)
			comments = append(comments, findAndCollDecorComments(typeDoc)...)
			log.Printf("append(comments, findAndCollDecorComments(typeDoc)...): %+v \n", comments)
			if len(comments) == 0 {
				return
			}
			if _, ok := typeNameMapDecorComments[spec.Name.Name]; ok {
				errs = append(errs, errors.New("duplicate type definition"))
				return
			}
			typeNameMapDecorComments[spec.Name.Name] = comments
		})
	}
	log.Printf("typeNameMapDecorComments: %+v \n", typeNameMapDecorComments)
	log.Printf("errs: %+v \n", errs)
	if len(typeNameMapDecorComments) == 0 {
		return
	}
	for _, f := range pkg.Files {
		visitAstDecl(f, func(decl *ast.FuncDecl) (r bool) {
			if decl.Recv == nil ||
				decl.Recv.List == nil ||
				len(decl.Recv.List) != 1 ||
				decl.Recv.List[0].Type == nil {
				return
			}
			var typeId *ast.Ident
			switch typ := decl.Recv.List[0].Type.(type) {
			case *ast.Ident:
				typeId = typ
			case *ast.StarExpr:
				id, ok := typ.X.(*ast.Ident)
				if !ok {
					return
				}
				typeId = id
			}
			if typeId.Name == "" {
				return
			}
			comments, ok := typeNameMapDecorComments[typeId.Name]
			if !ok || len(comments) == 0 {
				return
			}
			log.Printf("decl: %+v, comments: %+v\n", decl, comments)
			if decl.Doc == nil {
				decl.Doc = &ast.CommentGroup{List: comments}
			} else {
				decl.Doc.List = append(decl.Doc.List, comments...)
			}
			return
		})
	}

}

func friendlyIDEPosition(fset *token.FileSet, p token.Pos) string {
	if runtime.GOOS == "windows" {
		return fset.Position(p).String()
	}
	return filepath.Join("./", fset.Position(p).String())
}
