package main

import (
	"bytes"
	"errors"
	"go/ast"
	"go/printer"
	"go/token"
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

	decorWrappedCodeFilePath := ""
	if dpp, err := getPackageInfo(decoratorPackagePath); err == nil {
		decorWrappedCodeFilePath = dpp.Dir + "/wrapped_code.go"
		files = append(files, decorWrappedCodeFilePath)
	}

	logs.Debug("packageName", packageName, files, args)

	var originPath string

	fset := token.NewFileSet()
	pkg, err := parserGOFiles(fset, files...)
	if err != nil {
		logs.Error(err)
	}

	errPos, err := typeDecorRebuild(pkg)
	if err != nil {
		logs.Error(err, biSymbol,
			friendlyIDEPosition(fset, errPos))
	}

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
			mapDecors := newMapV[string, *ast.Comment]()
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
				if !mapDecors.put(decorName, doc) {
					logs.Error("cannot use the same decorator for repeated decoration", biSymbol,
						"Decor:", friendlyIDEPosition(fset, doc.Pos()), biSymbol,
						"Repeated:", friendlyIDEPosition(fset, mapDecors.get(decorName).Pos()))
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
						"Target:", friendlyIDEPosition(fset, fd.Pos()), biSymbol,
						"Decor:", friendlyIDEPosition(fset, da.doc.Pos()))
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
					logs.Error(err, biSymbol, "Decor:", friendlyIDEPosition(fset, da.doc.Pos()))
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
				if wcf, ok := pkg.Files[decorWrappedCodeFilePath]; ok {
					assignWrappedCodePos(genStmts, wcf.Decls[0].(*ast.FuncDecl).Body.List)
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

func assignWrappedCodePos(from, reset []ast.Stmt) {
	{
		partFrom := from[0].(*ast.AssignStmt)
		partReset := reset[0].(*ast.AssignStmt)
		partFrom.TokPos = partReset.Pos()
		partFrom.Tok = partReset.Tok
		assignStmtPos(partFrom.Lhs[0], partReset.Lhs[0], true)
		assignStmtPos(partFrom.Rhs[0], partReset.Rhs[0], false)
		{
			l := partFrom.Rhs[0].(*ast.UnaryExpr).X.(*ast.CompositeLit)
			r := partReset.Rhs[0].(*ast.CompositeLit)
			l.Lbrace = r.Lbrace
			l.Rbrace = r.Rbrace
			assignStmtPos(l.Type, r.Type, true)
			//l.Type.(*ast.SelectorExpr).X.(*ast.Ident).NamePos = r.Type.(*ast.Ident).NamePos
			for i, kv := range l.Elts {
				rv := r.Elts[i].(*ast.KeyValueExpr)
				v := kv.(*ast.KeyValueExpr)
				assignStmtPos(v, rv, true)
			}
		}
	}
	{
		partFrom := from[1].(*ast.AssignStmt)
		partReset := reset[1].(*ast.AssignStmt)
		assignStmtPos(partFrom.Lhs[0], partReset.Lhs[0], true)
		//partFrom.Lhs[0].(*ast.SelectorExpr).X.(*ast.Ident).NamePos = partReset.Lhs[0].(*ast.SelectorExpr).X.(*ast.Ident).NamePos
		//partFrom.Lhs[0].(*ast.SelectorExpr).Sel.NamePos = partReset.Lhs[0].(*ast.SelectorExpr).Sel.NamePos
		partFrom.Tok = partReset.Tok
		//partFrom.Rhs[0].(*ast.FuncLit)
		assignStmtPos(partFrom.Rhs[0], partReset.Rhs[0], true)
		var flit *ast.CallExpr
		r := partReset.Rhs[0].(*ast.FuncLit).Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr)
		if astmt, ok := partFrom.Rhs[0].(*ast.FuncLit).Body.List[0].(*ast.AssignStmt); ok {
			assignStmtPos(astmt.Lhs[0], r, true)
			flit = partFrom.Rhs[0].(*ast.FuncLit).Body.List[0].(*ast.AssignStmt).Rhs[0].(*ast.CallExpr)
		} else {
			flit = partFrom.Rhs[0].(*ast.FuncLit).Body.List[0].(*ast.ExprStmt).X.(*ast.CallExpr)
		}
		flit.Lparen = r.Lparen
		flit.Rparen = r.Rparen
		//TODO
		//if flit.Args != nil {
		//	for _, arg := range flit.Args {
		//		assignStmtPos(arg, r.Args[0], false)
		//	}
		//}
	}
	// has-return
	if len(from) >= 4 {
		l := from[3].(*ast.ReturnStmt)
		r := reset[2].(*ast.ReturnStmt)
		l.Return = r.Return
		if l.Results != nil {
			//for _, v := range l.Results {
			//	assignStmtPos(v.)
			//}
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

func typeDecorRebuild(pkg *ast.Package) (pos token.Pos, err error) {
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
	type errSet struct {
		pos token.Pos
		err error
	}
	errs := []*errSet{}
	for _, f := range pkg.Files {
		typeDeclVisitor(f.Decls, func(spec *ast.TypeSpec, typeDoc *ast.CommentGroup) {
			if (spec.Doc == nil || spec.Doc.List == nil) &&
				(typeDoc == nil || typeDoc.List == nil) {
				return
			}
			comments := findAndCollDecorComments(spec.Doc)
			//log.Printf("findAndCollDecorComments(spec.Doc): %+v \n", comments)
			comments = append(comments, findAndCollDecorComments(typeDoc)...)
			//log.Printf("append(comments, findAndCollDecorComments(typeDoc)...): %+v \n", comments)
			if len(comments) == 0 {
				return
			}
			if _, ok := typeNameMapDecorComments[spec.Name.Name]; ok {
				errs = append(errs, &errSet{
					pos: spec.Name.NamePos,
					err: errors.New("duplicate type definition: " + spec.Name.Name),
				})
				return
			}
			typeNameMapDecorComments[spec.Name.Name] = comments
		})
		if len(errs) > 0 {
			return errs[0].pos, errs[0].err
		}
	}
	//log.Printf("typeNameMapDecorComments: %+v \n", typeNameMapDecorComments)
	//log.Printf("errs: %+v \n", errs)
	if len(typeNameMapDecorComments) == 0 {
		return
	}
	identName := func(expr ast.Expr) string {
		switch expr := expr.(type) {
		case *ast.Ident: // normal: var
			return expr.Name
		case *ast.IndexListExpr: // var[T]
			if v, ok := expr.X.(*ast.Ident); ok {
				return v.Name
			}
			return ""
		case *ast.IndexExpr: //  var[K,V]
			if v, ok := expr.X.(*ast.Ident); ok {
				return v.Name
			}
			return ""
		case *ast.StarExpr: // pointer
			switch x := expr.X.(type) {
			case *ast.Ident: // *var
				return x.Name
			case *ast.IndexExpr: // *var[K]
				if v, ok := x.X.(*ast.Ident); ok {
					return v.Name
				}
				return ""
			case *ast.IndexListExpr: // *var[K,V]
				if v, ok := x.X.(*ast.Ident); ok {
					return v.Name
				}
				return ""
			default:
				return ""
			}
		}
		return ""
	}
	for _, f := range pkg.Files {
		visitAstDecl(f, func(decl *ast.FuncDecl) (r bool) {
			if decl.Recv == nil ||
				decl.Recv.List == nil ||
				len(decl.Recv.List) != 1 ||
				decl.Recv.List[0].Type == nil {
				return
			}
			typeIdName := identName(decl.Recv.List[0].Type)
			if typeIdName == "" {
				return
			}
			comments, ok := typeNameMapDecorComments[typeIdName]
			if !ok || len(comments) == 0 {
				return
			}
			//log.Printf("decl: %+v, comments: %+v\n", decl, comments)
			if decl.Doc == nil {
				decl.Doc = &ast.CommentGroup{List: comments}
			} else {
				decl.Doc.List = append(decl.Doc.List, comments...)
			}
			return
		})
	}
	return
}

func friendlyIDEPosition(fset *token.FileSet, p token.Pos) string {
	if runtime.GOOS == "windows" {
		return fset.Position(p).String()
	}
	return filepath.Join("./", fset.Position(p).String())
}
