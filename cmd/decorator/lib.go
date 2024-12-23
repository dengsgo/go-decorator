package main

import (
	"bytes"
	"fmt"
	"github.com/dengsgo/go-decorator/cmd/logs"
	"go/ast"
	"go/parser"
	"go/printer"
	"go/token"
	"math/rand"
	"strconv"
	"strings"
	"text/template"
)

const randSeeds = "abcdefghijklmnopqrstuvwxyz"

var emptyFset = token.NewFileSet()

const replaceTpl = `    ${.DecorVarName} := &decor.Context{
        Kind:       decor.${.TKind},
        TargetName: ${.TargetName},
        Receiver:   ${.ReceiverVarName},
        TargetIn:   []any{${stringer .InArgNames}},
        TargetOut:  []any{${stringer .OutArgNames}},
    }
    ${.DecorVarName}.Func = func() {
        ${if .HaveReturn}${stringer .DecorListOut} = ${end}${.FuncMain} (${stringer .DecorCallIn})
    }
    ${.DecorCallName}(${.DecorVarName}${if .HaveDecorParam}, ${stringer .DecorCallParams}${end})
    ${if .HaveReturn}return ${stringer .DecorCallOut}${end}`

type ReplaceArgs struct {
	HaveDecorParam,
	HaveReturn bool
	TKind, // target kind
	TargetName,
	ReceiverVarName, // Receiver var
	DecorVarName, // decor var
	DecorCallName, // decor function name . logging
	FuncMain string // (a, b, c) {raw func}
	DecorCallParams, // decor function parameters. like "", 0, true, options, default empty
	InArgNames, // a, b, c
	OutArgNames, // c, d
	InArgTypes, // int, int, int
	OutArgTypes, // int, int
	DecorListOut, // decor.TargetOut[0], decor.TargetOut[1]
	DecorCallIn, // decor.TargetIn[0].(int), decor.TargetIn[1].(int), decor.TargetIn[2].(int)
	DecorCallOut []string // decor.TargetOut[0].(int), decor.TargetOut[1].(int)
}

func newReplaceArgs(gi *genIdentId, targetName, decorName string) *ReplaceArgs {
	return &ReplaceArgs{
		false,
		false,
		"KFunc", // decor.TKind,
		`"` + targetName + `"`,
		"nil",
		gi.nextStr(),
		decorName,
		"",
		[]string{},
		[]string{},
		[]string{},
		[]string{},
		[]string{},
		[]string{},
		[]string{},
		[]string{},
	}
}

func replace(args *ReplaceArgs) (string, error) {
	tpl, err := template.
		New("decorReplace").
		Delims("${", "}").
		Funcs(map[string]any{"stringer": stringer}).
		Parse(replaceTpl)
	if err != nil {
		return "", err
	}
	bf := bytes.NewBuffer([]byte{})
	err = tpl.Execute(bf, args)
	if err != nil {
		return "", err
	}
	return bf.String(), nil
}

func builderReplaceArgs(f *ast.FuncDecl, decorName string, decorParams []string, gi *genIdentId) *ReplaceArgs {
	ra := newReplaceArgs(gi, f.Name.Name, decorName)
	// decor params
	if decorParams != nil && len(decorParams) > 0 {
		ra.HaveDecorParam = true
		ra.DecorCallParams = decorParams
	}
	// target TKind
	if f.Recv != nil && f.Recv.List != nil && len(f.Recv.List) > 0 {
		ra.TKind = "KMethod"
		ra.ReceiverVarName = f.Recv.List[0].Names[0].Name
	}
	//funcMain
	var tp *ast.FieldList
	if f.Type != nil && f.Type.TypeParams != nil {
		tp = f.Type.TypeParams
		f.Type.TypeParams = nil
	}
	closure := &ast.FuncLit{
		Type: f.Type,
		Body: f.Body,
	}
	var output []byte
	buffer := bytes.NewBuffer(output)
	err := printer.Fprint(buffer, token.NewFileSet(), closure)
	if err != nil {
		logs.Error("builderReplaceArgs printer.Fprint fail", decorName, err)
	}
	f.Type.TypeParams = tp
	ra.FuncMain = buffer.String()

	// in result
	if f.Type.Results != nil && f.Type.Results.List != nil {
		for _, r := range f.Type.Results.List {
			if r.Names != nil {
				continue
			}
			r.Names = []*ast.Ident{
				{
					NamePos: 0,
					Name:    gi.nextStr(),
					Obj:     nil,
				},
			}
		}
		count := 0
		for _, r := range f.Type.Results.List {
			if len(r.Names) == 0 {
				continue
			}
			for _, p := range r.Names {
				if p.Name == "_" {
					// fix issue #10. If the parameter name is “_”, we need to create a new name to replace it since the context will use this variable
					p.Name = gi.nextStr()
				}
				ra.OutArgNames = append(ra.OutArgNames, p.Name)
				ra.OutArgTypes = append(ra.OutArgTypes, typeString(r.Type))
				ra.DecorListOut = append(ra.DecorListOut,
					fmt.Sprintf("%s.TargetOut[%d]", ra.DecorVarName, count))
				ra.DecorCallOut = append(ra.DecorCallOut,
					//fmt.Sprintf("%s.TargetOut[%d].(%s)", ra.DecorVarName, count, typeString(r.Type)))
					fmt.Sprintf(
						"func() %s {o,_ := %s.TargetOut[%d].(%s); return o}()",
						typeString(r.Type),
						ra.DecorVarName,
						count,
						typeString(r.Type),
					),
				)
				count++
			}
		}
	}

	// in args
	if f.Type.Params.List != nil && len(f.Type.Params.List) > 0 {
		count := 0
		for _, r := range f.Type.Params.List {
			if len(r.Names) == 0 {
				continue
			}
			for _, p := range r.Names {
				if p.Name == "_" {
					// fix issue #10. If the parameter name is “_”, we need to create a new name to replace it since the context will use this variable
					p.Name = gi.nextStr()
				}
				ra.InArgNames = append(ra.InArgNames, p.Name)
				ra.InArgTypes = append(ra.InArgTypes, typeString(r.Type))
				ra.DecorCallIn = append(ra.DecorCallIn,
					//fmt.Sprintf("%s.TargetIn[%d].(%s)%s", ra.DecorVarName, count, typeString(r.Type), elString(r.Type)))
					fmt.Sprintf(
						"func() %s {o,_ := %s.TargetIn[%d].(%s); return o}()%s",
						typeString(r.Type),
						ra.DecorVarName,
						count,
						typeString(r.Type),
						elString(r.Type),
					),
				)
				count++
			}
		}
	}

	ra.HaveReturn = len(ra.OutArgNames) != 0
	return ra
}

func typeString(expr ast.Expr) string {
	var output []byte
	buffer := bytes.NewBuffer(output)
	err := printer.Fprint(buffer, emptyFset, expr)
	if err != nil {
		logs.Error("typeString printer.Fprint fail", err)
	}
	s := buffer.String()
	if strings.HasPrefix(s, "...") {
		return "[]" + s[3:]
	}
	return s
}

func elString(expr ast.Expr) string {
	if strings.HasPrefix(typeString(expr), "[]") {
		return "..."
	}
	return ""
}

func stringer(elems []string) string {
	if elems == nil {
		return ""
	}
	return strings.Join(elems, ", ")
}

func randStr(le int) string {
	s := ""
	for i := 0; i < le; i++ {
		index := rand.Intn(len(randSeeds))
		s += string(randSeeds[index])
	}
	return s
}

type genIdentId struct {
	id    int
	ident string
}

func newGenIdentId() *genIdentId {
	suf := randStr(6)
	return &genIdentId{
		id:    0,
		ident: "_decorGenIdent" + suf,
	}
}

func (g *genIdentId) next() int {
	g.id++
	return g.id
}

func (g *genIdentId) nextStr() string {
	g.next()
	return g.ident + strconv.Itoa(g.id)
}

// TODO
func funIsDecorator(fd *ast.FuncDecl, pkgName string) bool {
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

func getStmtList(s string) (r []ast.Stmt, i int, err error) {
	s = "func(){\n" + s + "\n}()"
	//logs.Debug("getStmtList", s)
	expr, err := parser.ParseExpr(s)
	if err != nil {
		return
	}
	r = expr.(*ast.CallExpr).Fun.(*ast.FuncLit).Body.List
	i = 0
	return
}

// same like /usr/local/go/src/go/parser/interface.go:139#ParseDir
func parserGOFiles(fset *token.FileSet, files ...string) (*ast.Package, error) {
	var pkg *ast.Package
	for _, file := range files {
		f, err := parser.ParseFile(fset, file, nil, parser.ParseComments)
		if err != nil {
			return pkg, err
		}
		if pkg == nil {
			pkg = &ast.Package{
				Name:  f.Name.Name,
				Files: make(map[string]*ast.File),
			}
		}
		pkg.Files[file] = f
	}
	return pkg, nil
}

func assignStmtPos(f, t ast.Node, depth bool) {
	if f == nil || t == nil {
		return
	}
	switch v := f.(type) {
	case nil:
		return
	case *ast.Ident:
		v.NamePos = t.Pos()
	case *ast.BasicLit:
		if v != nil {
			v.ValuePos = t.Pos()
		}
	case *ast.UnaryExpr:
		v.OpPos = t.Pos()
		if depth {
			assignStmtPos(v.X, t, depth)
		}
	case *ast.IndexExpr:
		v.Lbrack = t.Pos()
		v.Rbrack = t.Pos()
		assignStmtPos(v.X, t, depth)
		assignStmtPos(v.Index, t, depth)
	case *ast.AssignStmt:
		v.TokPos = t.Pos()
		if depth {
			for _, lhs := range v.Lhs {
				assignStmtPos(lhs, t, depth)
			}
			for _, rhs := range v.Rhs {
				assignStmtPos(rhs, t, depth)
			}
		}
	case *ast.CompositeLit:
		v.Lbrace = t.Pos()
		v.Rbrace = t.End()
		if depth {
			assignStmtPos(v.Type, t, depth)
			if v.Elts != nil {
				for _, els := range v.Elts {
					assignStmtPos(els, t, depth)
				}
			}
		}
	case *ast.KeyValueExpr:
		v.Colon = t.Pos()
		if depth {
			assignStmtPos(v.Key, t, depth)
			assignStmtPos(v.Value, t, depth)
		}
	case *ast.ArrayType:
		v.Lbrack = t.Pos()
		if depth {
			assignStmtPos(v.Len, t, depth)
			assignStmtPos(v.Elt, t, depth)
		}
	case *ast.SelectorExpr:
		assignStmtPos(v.Sel, t, depth)
		assignStmtPos(v.X, t, depth)
	case *ast.FuncLit:
		assignStmtPos(v.Type, t, depth)
		if depth {
			assignStmtPos(v.Body, t, depth)
		}
	case *ast.FuncType:
		v.Func = t.Pos()
		assignStmtPos(v.Params, t, depth)
		assignStmtPos(v.Results, t, depth)
	case *ast.BlockStmt:
		v.Lbrace = t.Pos()
		v.Rbrace = t.End()
		if depth && v.List != nil {
			for _, st := range v.List {
				assignStmtPos(st, t, depth)
			}
		}
	case *ast.TypeAssertExpr:
		v.Lparen = t.Pos()
		v.Rparen = t.End()
		assignStmtPos(v.Type, t, depth)
		assignStmtPos(v.X, t, depth)
	case *ast.FieldList:
		if v == nil {
			return
		}
		v.Opening = t.Pos()
		v.Closing = t.Pos()
		if depth && v.List != nil {
			for _, field := range v.List {
				assignStmtPos(field, t, depth)
			}
		}
	case *ast.Field:
		if v == nil {
			return
		}
		assignStmtPos(v.Type, t, depth)
		assignStmtPos(v.Tag, t, depth)
		if v.Names != nil {
			for _, name := range v.Names {
				assignStmtPos(name, t, depth)
			}
		}
	case *ast.CallExpr:
		v.Lparen = t.Pos()
		v.Rparen = t.Pos()
		if v.Args != nil {
			for _, arg := range v.Args {
				assignStmtPos(arg, t, depth)
			}
		}
		if depth {
			assignStmtPos(v.Fun, t, depth)
		}
	default:
		logs.Info("can`t support type from assignStmtPos")
	}
}
