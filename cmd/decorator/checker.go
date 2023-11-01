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
	"go/types"
	"strings"
	"unicode"
	"unicode/utf8"
)

var (
	errUsedDecorSyntaxErrorLossFunc  = errors.New("syntax error using decorator: miss decorator name")
	errUsedDecorSyntaxErrorLossValue = errors.New("syntax error using decorator: miss parameters value")
	errUsedDecorSyntaxErrorInvalidP  = errors.New("syntax error using decorator: invalid parameter format")
	errUsedDecorSyntaxError          = errors.New("syntax error using decorator")
	errCalledDecorNotDecorator       = errors.New("used decor is not a decorator function")
)

var (
	decorOptionParamTypeMap = map[string]types.BasicInfo{
		"bool": types.IsBoolean,

		"int":    types.IsInteger,
		"int8":   types.IsInteger,
		"in16":   types.IsInteger,
		"int32":  types.IsInteger,
		"int64":  types.IsInteger,
		"unit":   types.IsInteger,
		"unit8":  types.IsInteger,
		"unit16": types.IsInteger,
		"unit32": types.IsInteger,
		"unit64": types.IsInteger,

		"float32": types.IsFloat,
		"float64": types.IsFloat,

		"string": types.IsString,
	}
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
	//   function#{key:""}
	//   function#{key:"", name:""}
	//   function#{key:"", name:"", age:100}
	//   function#{key:"", name:"", age:100, b : false}
	if s == "" {
		return "", nil, errUsedDecorSyntaxErrorLossFunc
	}

	_callName, pStr, _ := strings.Cut(s, "#")
	cAst, err := parser.ParseExpr(_callName)
	if err != nil {
		return "", nil, errUsedDecorSyntaxError
	}
	callName := ""
	switch a := cAst.(type) {
	case *ast.SelectorExpr, *ast.Ident:
		callName = typeString(a)
	default:
		return "", nil, errUsedDecorSyntaxError
	}
	if callName == "" {
		// non
		return callName, nil, errUsedDecorSyntaxErrorLossFunc
	}
	p := map[string]string{}
	pStr = strings.TrimSpace(pStr)
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

	exprList, err := parseDecorParameterStringToExprList(pStr)
	if err != nil {
		return callName, p, err
	}
	px := mapx(p)
	if err := decorStmtListToMap(exprList, px); err != nil {
		return callName, p, err
	}
	return callName, px, nil
}

func decorStmtListToMap(exprList []ast.Expr, p mapx) error {
	ident := func(v ast.Expr) string {
		if v == nil {
			return ""
		}
		id, ok := v.(*ast.Ident)
		if !ok {
			return ""
		}
		return id.Name
	}
	consumerKeyValue := func(expr *ast.KeyValueExpr) error {
		key := ident(expr.Key)
		if key == "" {
			return errors.New("invalid parameter name") // error
		}
		switch value := expr.Value.(type) {
		case *ast.BasicLit:
			switch value.Kind {
			// a:"b"
			// a: 0
			// a: 0.0
			case token.STRING, token.INT, token.FLOAT:
				if !p.put(key, value.Value) {
					return errors.New("duplicate parameters key '" + key + "'")
				}
			default:
				return errors.New("invalid parameter type") // error
			}
		case *ast.Ident:
			val := ident(value)
			if val != "true" && val != "false" {
				return errors.New("invalid parameter value, should be bool")
			}
			if !p.put(key, val) {
				return errors.New("duplicate parameters key '" + key + "'")
			}
		default:
			return errors.New("invalid parameter value")
		}
		return nil
	}
	for _, v := range exprList {
		switch expr := v.(type) {
		case *ast.KeyValueExpr: // a:b
			if err := consumerKeyValue(expr); err != nil {
				return err
			}
		default:
			return errUsedDecorSyntaxErrorInvalidP // error
		}
	}
	return nil // error
}

// s = {xxxxx}
func parseDecorParameterStringToExprList(s string) ([]ast.Expr, error) {
	s = "map[any]any" + s
	stmts, _, err := getStmtList(s)
	if err != nil {
		return nil, errUsedDecorSyntaxErrorInvalidP
	}
	if len(stmts) != 1 {
		return nil, errUsedDecorSyntaxErrorInvalidP
	}
	stmt, ok := stmts[0].(*ast.ExprStmt)
	if !ok || stmt == nil {
		return nil, errUsedDecorSyntaxErrorInvalidP
	}
	clit, ok := stmt.X.(*ast.CompositeLit)
	if !ok || clit == nil {
		return nil, errUsedDecorSyntaxErrorInvalidP
	}
	if clit.Elts == nil {
		return nil, errUsedDecorSyntaxErrorInvalidP
	}
	return clit.Elts, nil
}

func checkDecorAndGetParam(pkgPath, funName string, annotationMap map[string]string) ([]string, error) {
	decl, file, err := pkgILoader.findFunc(pkgPath, funName)
	if err != nil {
		return nil, err
	}
	imp := newImporter(file)
	pkgName, ok := imp.importedPath(decoratorPackagePath)
	if !ok {
		return nil, errors.New(msgDecorPkgNotFound)
	}
	m := collDeclFuncParamsAnfTypes(decl)
	if len(m) < 1 {
		return nil, errCalledDecorNotDecorator
	}
	for _, v := range m {
		if v.index == 0 && v.typ != fmt.Sprintf("*%s.Context", pkgName) {
			return nil, errors.New("used decor is not a decorator function")
		}
	}
	if len(m) == 1 {
		return []string{}, nil
	}
	params := make([]string, len(m))
	for _, v := range m {
		if v.index == 0 {
			continue
		}
		if value, ok := annotationMap[v.name]; ok {
			params[v.index] = value
		} else {
			typ, ok := decorOptionParamTypeMap[v.typ]
			if !ok {
				return nil, errors.New("unsupported types '" + v.typ + "'")
			}
			switch typ {
			case types.IsInteger:
				params[v.index] = "0"
			case types.IsFloat:
				params[v.index] = "0.0"
			case types.IsString:
				params[v.index] = `""`
			case types.IsBoolean:
				params[v.index] = "false"
			}
		}
	}

	//go:decor logging#(key : "")   func(key, name, instance string)
	return params[1:], nil
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
