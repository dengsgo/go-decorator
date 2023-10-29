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
	"unicode"
	"unicode/utf8"
)

var (
	errUsedDecorSyntaxErrorLossFunc  = errors.New("syntax error using decorator: miss decorator name")
	errUsedDecorSyntaxErrorLossValue = errors.New("syntax error using decorator: miss parameters value")
	errUsedDecorSyntaxError          = errors.New("syntax error using decorator")
	errCalledDecorNotDecorator       = errors.New("used decor is not a decorator function")
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

func parseDecorAndParameters(s string) (string, map[string]string, error) {
	// s like:
	//   function
	//   function#{}
	//   function#{key=""}
	//   function#{key="", name=""}
	//   function#{key="", name="", age=100}
	//   function#{key="", name="", age=100, b = false}
	if s == "" {
		return "", nil, errUsedDecorSyntaxErrorLossFunc
	}

	callName, pStr, _ := strings.Cut(s, "#")
	callName = cleanSpaceChar(callName)
	if callName == "" {
		// non
		return callName, nil, errUsedDecorSyntaxErrorLossFunc
	}
	p := map[string]string{}
	if pStr == "" {
		if strings.HasSuffix(s, "#") {
			// #
			return callName, p, errUsedDecorSyntaxError
		}
		return callName, p, nil
	}
	if pStr[0] != '{' || pStr[len(pStr)-1] != '}' {
		return callName, nil, errUsedDecorSyntaxError
	}
	if len(pStr) == 2 {
		// {}
		return callName, p, nil
	}
	if len(pStr) < 5 {
		return callName, p, errUsedDecorSyntaxError
	}
	consumerString := func(s string) (string, string) { // TODO
		offset := 0
		for offset < len(s) {
			r, size := utf8.DecodeRuneInString(s[offset:])
			if r == utf8.RuneError {
				return s[0:offset], s[offset:]
			}
			offset += size
			if r == '"' && s[offset-1:offset] != "\\" {
				return s[0:offset], s[offset:]
			}
		}
		return s[0:offset], s[offset:]
	}
	for {
		key, value, _ := strings.Cut(pStr, "=")
		key = strings.TrimSpace(key)
		if !isLetters(key) {
			return callName, p, errUsedDecorSyntaxError
		}
		value = strings.TrimLeftFunc(value, unicode.IsSpace)
		if len(value) < 2 {
			return callName, p, errUsedDecorSyntaxErrorLossValue
		}
		// TODO
		if value[0] == '"' {
			// consumer string
			consumerString(value[1:])
		} else if r, _ := utf8.DecodeRuneInString(value); r != utf8.RuneError && unicode.IsLetter(r) {
			// consumer rune
			p[key] = value

		}
	}

	//str := pStr[1 : len(pStr)-1]
	//for {
	//	if isLet() {
	//
	//	}
	//}

	return callName, p, nil
}

func isLetters(s string) (b bool) {
	for offset := 0; offset < len(s); {
		r, size := utf8.DecodeRuneInString(s[offset:])
		if r == utf8.RuneError {
			return b
		}
		offset += size
		if !unicode.IsLetter(r) {
			return false
		}
		b = true
	}
	return b
}

func cleanSpaceChar(s string) string {
	bf := bytes.NewBuffer([]byte{})
	offset := 0
	for offset < len(s) {
		r, size := utf8.DecodeRuneInString(s[offset:])
		offset += size
		if unicode.IsSpace(r) {
			continue
		}
		bf.WriteRune(r)
	}
	return bf.String()
}

func isLet() {

}

func checkDecorAndGetParam(pkgPath, funName string, annotationMap map[string]string) (args string, error error) {
	decl, file, err := pkgILoader.findFunc(pkgPath, funName)
	if err != nil {
		return "", err
	}
	imp := newImporter(file)
	pkgName, ok := imp.importedPath(decoratorPackagePath)
	if !ok {
		return "", errors.New(msgDecorPkgNotFound)
	}
	m := collDeclFuncParamsAnfTypes(decl)
	if len(m) < 1 {
		return "", errCalledDecorNotDecorator
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
	pi, err := getPackageInfo(pkgPath)
	if err != nil {
		return nil, err
	}
	set.fset = token.NewFileSet()
	set.pkgs, err = parser.ParseDir(set.fset, pi.Dir, nil, 0)
	d.pkg[pkgPath] = set
	return
}
