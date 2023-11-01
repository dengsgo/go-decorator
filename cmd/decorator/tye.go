package main

import (
	"go/ast"
	"strings"
)

type mapx map[string]string

func (p mapx) put(key, value string) bool {
	if _, ok := p[key]; ok {
		return false
	}
	p[key] = value
	return true
}

type decorAnnotation struct {
	doc        *ast.Comment      // ast node for doc
	name       string            // decorator name
	parameters map[string]string // options parameters
}

func newDecorAnnotation(doc *ast.Comment, name string, parameters map[string]string) *decorAnnotation {
	return &decorAnnotation{
		doc:        doc,
		name:       name,
		parameters: parameters,
	}
}

func (d *decorAnnotation) splitName() []string {
	return strings.Split(d.name, ".")
}
