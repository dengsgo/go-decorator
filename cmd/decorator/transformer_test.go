package main

import "testing"

func TestGetGoModPath(t *testing.T) {
	s := getGoModPath()
	if s != "github.com/dengsgo/go-decorator" {
		t.Fatalf("getGoModPath != 'github.com/dengsgo/go-decorator', now = %s\n", s)
	}
}

func TestImporter(t *testing.T) {
	// TODO
}
