package main

import (
	"bytes"
	"errors"
	"fmt"
	"github.com/dengsgo/go-decorator/cmd/logs"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"strings"
)

func isDecoratorFunc(fd *ast.FuncDecl, pkgName string) bool {
	if pkgName == "" ||
		fd == nil ||
		fd.Recv != nil ||
		fd.Type == nil ||
		fd.Type.Params == nil ||
		fd.Type.Params.NumFields() != 1 ||
		fd.Type.Params.List[0] == nil ||
		fd.Type.Params.List[0].Type == nil {
		return false
	}
	expr := fd.Type.Params.List[0].Type
	buffer := bytes.NewBuffer([]byte{})
	err := printer.Fprint(buffer, emptyFset, expr)
	if err != nil {
		logs.Debug("funIsDecorator printer.Fprint fail", err)
		return false
	}
	return strings.TrimSpace(buffer.String()) == fmt.Sprintf("*%s.Context", pkgName)
}

func checkDecorAndGetParam(pkgPath, funName string, annotationMap map[string]string) (args string, error error) {
	decl, file, err := pkgILoader.findFunc(pkgPath, funName)
	if err != nil {
		return "", err
	}
	imp := newImporter(file)
	pkgName, ok := imp.importedPath(decoratorPackagePath)
	if !ok {
		return "", errors.New(msgDecorPkgNotImported)
	}
	m := collDeclFuncParamsAnfTypes(decl)
	if len(m) < 1 {
		return "", errors.New("used decor is not a decorator function")
	}
	for _, v := range m {
		if v.index == 0 && v.typ != fmt.Sprintf("*%s.Context", pkgName) {
			return "", errors.New("used decor is not a decorator function")
		}
	}
	if len(m) == 1 {
		return "", nil
	}
	params := make([]string, len(m))
	for _, v := range m {
		if v.index == 0 {
			continue
		}
		if value, ok := annotationMap[v.name]; ok {
			params[v.index] = value
		} else {
			switch {
			case strings.HasPrefix(v.typ, "int"):
				params[v.index] = "0"
			case strings.HasPrefix(v.typ, "float"):
				params[v.index] = "0.0"
			case v.typ == "string":
				params[v.index] = `""`
			case v.typ == "bool":
				params[v.index] = "false"
			default:
				return
			}

		}
	}

	//go:decor logging#(key = "")   func(key, name, instance string)
	return strings.Join(params[1:], ", "), nil
}

type decorParamType struct {
	index int
	name,
	typ string
}

func collDeclFuncParamsAnfTypes(fd *ast.FuncDecl) (m map[string]*decorParamType) {
	m = map[string]*decorParamType{}
	if fd == nil ||
		fd.Type == nil ||
		fd.Type.Params == nil ||
		fd.Type.Params.NumFields() == 0 ||
		fd.Type.Params.List[0] == nil {
		return m
	}
	index := 0
	for _, field := range fd.Type.Params.List {
		typ := typeString(field.Type)
		for _, id := range field.Names {
			m[id.Name] = &decorParamType{index, id.Name, typ}
			index++
		}
	}
	return m
}

type pkgSet struct {
	fset *token.FileSet
	pkgs map[string]*ast.Package
}

var pkgILoader = newPkgLoader()

type pkgLoader struct {
	pkg   map[string]*pkgSet
	funcs map[string]*ast.FuncDecl
}

func newPkgLoader() *pkgLoader {
	return &pkgLoader{
		pkg:   map[string]*pkgSet{},
		funcs: map[string]*ast.FuncDecl{},
	}
}

func (d *pkgLoader) findFunc(pkgPath, funName string) (target *ast.FuncDecl, file *ast.File, err error) {
	return d.findTarget(pkgPath, funName)
}

func (d *pkgLoader) findTarget(pkgPath string, funName string) (target *ast.FuncDecl, afile *ast.File, err error) {
	set, err := d.loadPkg(pkgPath)
	if err != nil {
		return nil, nil, err
	}
	err = errors.New("target not found")
	for _, v := range set.pkgs {
		if v == nil || v.Files == nil {
			continue
		}
		for _, file := range v.Files {
			visitAstDecl(file, func(decl *ast.FuncDecl) bool {
				if decl == nil ||
					decl.Name == nil ||
					decl.Name.Name != funName ||
					decl.Recv != nil {
					return false
				}

				afile = file
				target = decl
				err = nil
				return true
			})
		}
	}
	return
}

func (d *pkgLoader) loadPkg(pkgPath string) (set *pkgSet, err error) {
	if _set, ok := d.pkg[pkgPath]; ok {
		set = _set
		return
	}
	set = &pkgSet{}
	pi := getPkgCompiledInfo(pkgPath)
	if pi.dir == "" {
		s := "getPkgCompiledInfo fail"
		err = errors.New(s)
		logs.Debug(s)
		return
	}
	set.fset = token.NewFileSet()
	set.pkgs, err = parser.ParseDir(set.fset, pi.dir, nil, 0)
	d.pkg[pkgPath] = set
	return
}
