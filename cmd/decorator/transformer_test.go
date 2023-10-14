package main

import (
	"go/parser"
	"go/token"
	"testing"
)

const importWays = `
package main
 import (
 	_ "github.com/dengsgo/go-decorator/decor"
 	"github.com/dengsgo/cmd/logs"
 	"gopkg.in/yaml.v3"
 	_ "github.com/dengsgo/runner/v2"
	. "log"
	o "os"
 )
`

func TestGetGoModPath(t *testing.T) {
	s := getGoModPath()
	if s != "github.com/dengsgo/go-decorator" {
		t.Fatalf("getGoModPath != 'github.com/dengsgo/go-decorator', now = %s\n", s)
	}
}

func TestImporter(t *testing.T) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, "main.go", []byte(importWays), parser.ParseComments)
	if err != nil {
		t.Fatal("parse importWays content error", err)
	}
	imp := newImporter(f)
	cLen := 6
	if len(imp.nameMap) != cLen || len(imp.pathMap) != cLen || len(imp.pathObjMap) != cLen {
		t.Fatal("newImporter() error,Clen=6 but got", len(imp.nameMap), len(imp.pathMap), len(imp.pathObjMap))
	}
	cases := []struct {
		name,
		pkg string
	}{
		{"decor", "github.com/dengsgo/go-decorator/decor"},
		{"logs", "github.com/dengsgo/cmd/logs"},
		{"yaml", "gopkg.in/yaml.v3"},
		{"runner", "github.com/dengsgo/runner/v2"},
		{"log", "log"},
		{"o", "os"},
	}
	for _, v := range cases {
		pkg, ok := imp.importedName(v.name)
		if !ok {
			t.Fatal("importedName() error, name not found", v.name)
		}
		if pkg != v.pkg {
			t.Fatalf("importedName() error, got %s, want %s\n", pkg, v.pkg)
		}
	}
	cass := []struct {
		pkg,
		name string
	}{
		{"github.com/dengsgo/go-decorator/decor", "_"},
		{"github.com/dengsgo/cmd/logs", "logs"},
		{"gopkg.in/yaml.v3", "yaml"},
		{"github.com/dengsgo/runner/v2", "_"},
		{"log", "."},
		{"os", "o"},
	}
	for _, v := range cass {
		name, ok := imp.importedPath(v.pkg)
		if !ok {
			t.Fatal("importedPath() error, pkg not found", v.name)
		}
		if name != v.name {
			t.Fatalf("importedPath() error, got %s, want %s\n", name, v.name)
		}
	}
}
