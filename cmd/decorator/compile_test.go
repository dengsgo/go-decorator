package main

import (
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestDecorX(t *testing.T) {
	sucCases := []string{
		"log.Println",
		"fmt.Printf",
		"x.s",
		"log.",
		"decor.Context",
	}
	failCases := []string{
		"log",
		".Printf",
		"",
		"x.a.",
		"aaaa##c",
	}
	for i, s := range sucCases {
		if decorX(s) == "" {
			t.Fatalf("decorX('%s') should pass, case sucCases i: %d\n", s, i)
		}
	}
	for i, s := range failCases {
		if decorX(s) != "" {
			t.Fatalf("decorX('%s') should fail, case failCases i: %d\n", s, i)
		}
	}
}

func TestReverseSlice(t *testing.T) {
	t.Run("IntSlice", func(t *testing.T) {
		in := []int{1, 2, 3, 4, 5}
		out := []int{5, 4, 3, 2, 1}
		for i, v := range reverseSlice(in) {
			if v != out[i] {
				t.Fatalf("TestReverseSlice IntSlice fail, i %+v\n", i)
			}
		}
	})
	t.Run("StringSlice", func(t *testing.T) {
		in := []string{"string", "int", "bool"}
		out := []string{"bool", "int", "string"}
		for i, v := range reverseSlice(in) {
			if v != out[i] {
				t.Fatalf("TestReverseSlice StringSlice fail, i %+v\n", i)
			}
		}
	})
}

func TestTypeDeclVisitor(t *testing.T) {
	src := `package main
type (
	S      string
	I      int
	Struct struct{}
)
type SI struct{}
var a string = ""
const name string = ""
func E(){}
func (m maps[K, V])W(){}
func (m *maps[K, V])M(){}
`
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "", src, parser.ParseComments)
	if err != nil {
		t.Fatal(err)
	}
	//ast.Print(fset, f)
	sym := []string{"S", "I", "Struct", "SI"}
	count := 0
	typeDeclVisitor(f.Decls, func(spec *ast.TypeSpec, group *ast.CommentGroup) {
		if group != nil {
			t.Fatal("TestTypeDeclVisitor group should be nil, but got ", group)
		}
		if !inSlice(sym, spec.Name.Name) {
			t.Fatal("TestTypeDeclVisitor inSlice should be true, but got false", spec.Name.Name)
		}
		count++
	})
	if count != len(sym) {
		t.Fatal("TestTypeDeclVisitor should found ", len(sym), ", but got", count)
	}
}

func inSlice[T comparable](in []T, target T) bool {
	for _, v := range in {
		if v == target {
			return true
		}
	}
	return false
}
