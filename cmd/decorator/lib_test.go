package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"testing"
)

func TestStringer(t *testing.T) {
	sucCases := []struct {
		cas []string
		r   string
	}{
		{[]string{}, ""},
		{[]string{"var1"}, "var1"},
		{[]string{"var1", "var2"}, "var1, var2"},
		{[]string{"var1", "var2", "var3"}, "var1, var2, var3"},
	}
	for i, s := range sucCases {
		r := stringer(s.cas)
		if r != s.r {
			t.Fatalf("stringer('%s') != %s, now = %s case fail i: %d\n", s.cas, s.r, r, i)
		}
	}
}

func TestRandStr(t *testing.T) {
	for i := 0; i < 100; i++ {
		if len(randStr(i)) != i {
			t.Fatalf("randStr(%d) != %d, case fail\n", i, i)
		}
	}
}

func TestGenIdentId(t *testing.T) {
	gi := newGenIdentId()
	id := gi.nextStr()
	if id != gi.ident+"1" {
		t.Fatalf("first call gi.nextStr() != %s, now %s case fail\n", gi.ident+"1", id)
	}
	maps := map[string]bool{}
	for i := 0; i < 100; i++ {
		id := gi.nextStr()
		if _, ok := maps[id]; ok {
			t.Fatalf("result gi.nextStr() has already %s, Repeated. case fail\n", id)
		}
		maps[id] = true
	}
}

func TestGetStmtList(t *testing.T) {
	cases := []struct {
		text string
		leng int
	}{
		{
			``,
			0,
		},
		{
			`a:=0`,
			1,
		},
		{
			`a:=0
b:=2`,
			2,
		},
	}
	for _, cas := range cases {
		r, _, err := getStmtList(cas.text)
		if err != nil {
			t.Fatalf("getStmtList('%s') has error %s, case fail\n", cas.text, err.Error())
		}
		if len(r) != cas.leng {
			t.Fatalf("getStmtList('%s') result length != %d, now = %d, case fail\n", cas.text, cas.leng, len(r))
		}
	}

	failCases := []string{
		"ssssssssssssssssss+",
		"###edddd",
		"{{{{{ssss",
	}
	for _, cas := range failCases {
		if _, _, err := getStmtList(cas); err == nil {
			t.Fatalf("getStmtList('%s') should err, now = nil, case fail\n", cas)
		}
	}
}

func TestFunIsDecorator(t *testing.T) {
	check := func(name, pkgName string) {
		code := testGetCode(name, pkgName)
		f, err := parser.ParseFile(token.NewFileSet(), "main.go", code, parser.ParseComments)
		if err != nil || f == nil || len(f.Decls) == 0 {
			t.Fatal("TestFunIsDecorator testGetCode parse error", err)
		}
		i := 0
		for _, v := range f.Decls {
			fd, ok := v.(*ast.FuncDecl)
			if !ok {
				continue
			}
			i++
			if funIsDecorator(fd, pkgName) && fd.Name.Name != "isDecorator" {
				t.Fatal(fd.Name.Name, "should not be a decorator function")
			}
		}
		if i == 0 {
			t.Fatal("f.Decls have type *ast.FuncDecl functions. but got 0")
		}
	}
	check("", "decor")
	check("dec", "dec")
	check("a", "a")
}

func testGetCode(name, pkgName string) string {
	return fmt.Sprintf(`
package main
import %s "github.com/dengsgo/go-decorator/decor"
func isDecorator(ctx %s.Context) {}
func notDecorator1(ctx %s.Context, a int) {}
func notDecorator2(ctx %s.Contex) {}
func notDecorator3(a int) {}
`, name, pkgName, pkgName, pkgName)
}
