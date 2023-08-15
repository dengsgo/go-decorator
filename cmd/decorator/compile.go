package main

import (
	"bytes"
	"errors"
	"github.com/dengsgo/go-decorator/cmd/logs"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"golang.org/x/tools/go/ast/astutil"
	"os"
	"path"
	"path/filepath"
	"strings"
)

func compile(args []string) error {
	files := make([]string, 0, len(args))
	var cfg string
	projectName := getGoModPath()
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
		if strings.Contains(arg, filepath.Join("b001", "importcfg")) {
			cfg = arg
		} else if strings.HasPrefix(arg, projectDir+string(filepath.Separator)) && strings.HasSuffix(arg, ".go") {
			files = args[i:]
			break
		}
	}

	if (packageName != "main" && !strings.HasPrefix(packageName, projectName)) || cfg == "" || len(files) == 0 {
		return nil
	}

	logs.Debug("packageName", packageName, files, cfg, args)

	var originPath string
	imp := newImporter(cfg)

	for _, file := range files {
		fset := token.NewFileSet()
		logs.Debug("file Parse", file)
		f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			continue
		}
		logs.Debug(f.Decls)

		// decorators imports
		decorImports := []*ast.ImportSpec{}
		updated := false

		visitAstDecl(
			f,
			func(fd *ast.FuncDecl) {
				if fd.Doc == nil || fd.Doc.List == nil || len(fd.Doc.List) == 0 {
					return
				}
				originPath = file
				//log.Printf("%+v\n", fd)
				collDecors := []*ast.Comment{}
				for i := len(fd.Doc.List) - 1; i >= 0; i-- {
					doc := fd.Doc.List[i]
					if !strings.HasPrefix(doc.Text, decoratorScanFlag) {
						break
					}
					logs.Debug("HIT:", doc.Text)
					decorName, decorArgs, ok := parseGoDecComment(doc.Text)
					logs.Debug(decorName, decorArgs, ok)
					collDecors = append(collDecors, doc)
				}
				if len(collDecors) == 0 {
					return
				}
				logs.Debug("find comment FuncDecl entry", fset.Position(fd.Pos()))
				logs.Debug("collDecors", collDecors)
				gi := newGenIdentId()
				for _, doc := range collDecors {
					logs.Debug("handler:", doc.Text)
					decorName, decorArgs, ok := parseGoDecComment(doc.Text)
					logs.Debug(decorName, decorArgs, ok)
					// TODO 检查 decorName 是不是装饰器
					if fd.Recv != nil {
						logs.Warn("decor function is have Recv, ignore")
						continue
					}
					ra := builderReplaceArgs(fd, decorName, gi)
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
					assignCorrectPos(doc, ce)

					fd.Body.List = genStmts
					//x.Body.Rbrace = x.Body.Lbrace + token.Pos(ofs)
					//log.Printf("fd.Body.Pos() %+v\n", fd.Body.Pos())
					updated = true
				}
			},
			func(gd *ast.GenDecl) {
				if gd.Doc == nil || gd.Doc.List == nil || len(gd.Doc.List) == 0 {
					return
				}
				for i := len(gd.Doc.List) - 1; i >= 0; i-- {
					if !strings.HasPrefix(gd.Doc.List[i].Text, decoratorScanFlag) {
						break
					}
					arr := strings.Split(gd.Doc.List[i].Text, " ")
					if len(arr) == 2 {
						decorImports = append(decorImports, &ast.ImportSpec{
							Path: &ast.BasicLit{Value: arr[1]},
						})
					} else if len(arr) == 3 {
						decorImports = append(decorImports, &ast.ImportSpec{
							Name: &ast.Ident{Name: arr[1]},
							Path: &ast.BasicLit{Value: arr[2]},
						})
					} else {
						logs.Warn("import format fail", gd.Doc.List[i].Text)
					}

				}
			},
		)

		if !updated {
			continue
		}

		decorImports = append(decorImports, &ast.ImportSpec{
			Path: &ast.BasicLit{Value: decoratorPackagePath},
		})
		for _, v := range decorImports {
			if v.Name == nil {
				astutil.AddImport(fset, f, v.Path.Value)
			} else {
				astutil.AddNamedImport(fset, f, v.Name.Name, v.Path.Value)
			}
			err := imp.addImport(v.Path.Value)
			if err != nil {
				logs.Error("imp.addImport(v.Path.Value) error", v.Path.Value, err)
			}
		}

		if err := imp.sync(); err != nil {
			logs.Error("imp.sync()", err)
		}

		var output []byte
		buffer := bytes.NewBuffer(output)
		err = (&printer.Config{Tabwidth: 8, Mode: printer.SourcePos}).Fprint(buffer, fset, f)
		if err != nil {
			return errors.New("fprint original code")
		}
		tgDir := path.Join(tempDir, os.Getenv("TOOLEXEC_IMPORTPATH"))
		_ = os.MkdirAll(tgDir, 0777)
		tmpEntryFile := path.Join(tgDir, filepath.Base(originPath))
		logs.Debug("originPath", originPath, filepath.Base(originPath))
		err = os.WriteFile(tmpEntryFile, buffer.Bytes(), 0777)
		if err != nil {
			return errors.New("fail write into temporary file" + err.Error())
		}
		// update go build args
		for i := range args {
			if args[i] == originPath {
				args[i] = tmpEntryFile
			}
		}
		logs.Debug("args updated", args)
		logs.Info("rewrite file", originPath, "=>", tmpEntryFile)
	}

	return nil
}

func visitAstDecl(f *ast.File, funVisitor func(*ast.FuncDecl), genVisitor func(*ast.GenDecl)) {
	if f.Decls == nil {
		return
	}
	for _, t := range f.Decls {
		if t == nil {
			continue
		}
		switch decl := t.(type) {
		case *ast.FuncDecl:
			funVisitor(decl)
		case *ast.GenDecl:
			genVisitor(decl)
		}
	}
}

func assignCorrectPos(doc *ast.Comment, ce *ast.CallExpr) {
	ce.Lparen = doc.Pos()
	offset := token.Pos(0)
	if t, ok := ce.Fun.(*ast.Ident); ok {
		t.NamePos = doc.Pos()
		offset = token.Pos(len(t.Name))
	} else if t, ok := ce.Fun.(*ast.SelectorExpr); ok {
		if id, ok := t.X.(*ast.Ident); ok {
			id.NamePos = doc.Pos()
			offset = token.Pos(len(id.Name))
		}
		t.Sel.NamePos = doc.Pos() + offset + 1
		offset += token.Pos(len(t.Sel.Name)) + 1
	}
	for _, arg := range ce.Args {
		if id, ok := arg.(*ast.Ident); ok {
			id.NamePos = doc.Pos() + offset
		}
	}
}

func parseGoDecComment(s string) (decName string, args map[string]string, ok bool) {
	bodys := strings.Split(s[len(decoratorScanFlag):], "#")
	if len(bodys) > 2 || len(bodys) == 0 {
		return
	}
	decName = bodys[0]
	args = map[string]string{}
	if len(bodys) == 2 {
		a := strings.Split(bodys[1], ",")
		if len(a) > 0 {
			for _, v := range a {
				if len(v) > 0 {
					kv := strings.Split(v, "=")
					if len(kv) == 2 {
						args[kv[0]] = kv[1]
					}
				}
			}
		}
	}
	return decName, args, true
}
