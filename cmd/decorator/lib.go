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
        Kind:      decor.KFunc,
        TargetIn:  []any{${stringer .InArgNames}},
        TargetOut: []any{${stringer .OutArgNames}},
    }
    ${.DecorVarName}.Func = func() {
        ${if .HaveReturn}${stringer .DecorListOut} = ${end}${.FuncMain} (${stringer .DecorCallIn})
    }
    ${.DecorCallName}(${.DecorVarName})
    ${if .HaveReturn}return ${stringer .DecorCallOut}${end}`

type ReplaceArgs struct {
	HaveReturn    bool
	DecorVarName, // decor var
	DecorCallName, // decor function name . logging
	FuncMain string // (a, b, c) {raw func}
	InArgNames, // a, b, c
	OutArgNames, // c, d
	InArgTypes, // int, int, int
	OutArgTypes, // int, int
	DecorListOut, // decor.TargetOut[0], decor.TargetOut[1]
	DecorCallIn, // decor.TargetIn[0].(int), decor.TargetIn[1].(int), decor.TargetIn[2].(int)
	DecorCallOut []string // decor.TargetOut[0].(int), decor.TargetOut[1].(int)
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

func builderReplaceArgs(f *ast.FuncDecl, decorName string, gi *genIdentId) *ReplaceArgs {
	ra := &ReplaceArgs{false, gi.nextStr(), decorName, "", []string{}, []string{}, []string{}, []string{}, []string{}, []string{}, []string{}}
	//funcMain
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
				ra.OutArgNames = append(ra.OutArgNames, p.Name)
				ra.OutArgTypes = append(ra.OutArgTypes, typeString(r.Type))
				ra.DecorListOut = append(ra.DecorListOut,
					fmt.Sprintf("%s.TargetOut[%d]", ra.DecorVarName, count))
				ra.DecorCallOut = append(ra.DecorCallOut,
					fmt.Sprintf("%s.TargetOut[%d].(%s)", ra.DecorVarName, count, typeString(r.Type)))
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
				ra.InArgNames = append(ra.InArgNames, p.Name)
				ra.InArgTypes = append(ra.InArgTypes, typeString(r.Type))
				ra.DecorCallIn = append(ra.DecorCallIn,
					fmt.Sprintf("%s.TargetIn[%d].(%s)%s", ra.DecorVarName, count, typeString(r.Type), elString(r.Type)))
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
func funIsDecorator(fd *ast.FuncDecl) bool {
	if fd == nil ||
		fd.Recv != nil ||
		fd.Type == nil ||
		fd.Type.Params == nil ||
		fd.Type.Params.NumFields() != 1 {
		return false
	}
	if fd.Type.TypeParams == nil {
		return false
	}
	if fd.Type.TypeParams.NumFields() != 1 {
		return false
	}
	expr := fd.Type.TypeParams.List[0].Type
	buffer := bytes.NewBuffer([]byte{})
	err := printer.Fprint(buffer, emptyFset, expr)
	if err != nil {
		logs.Debug("funIsDecorator printer.Fprint fail", err)
		return false
	}
	// TODO

	return true
}

func getStmtList(s string) (r []ast.Stmt, i int, err error) {
	s = "func(){\n" + s + "\n}()"
	logs.Debug("getStmtList", s)
	expr, err := parser.ParseExpr(s)
	if err != nil {
		return
	}
	r = expr.(*ast.CallExpr).Fun.(*ast.FuncLit).Body.List
	i = 0
	return
}
