package main

import (
	"go/ast"
	"go/types"
	"strings"
)

type lintComparableKey = string

const (
	lintCpGt  lintComparableKey = "gt"
	lintCpGte lintComparableKey = "gte"
	lintCpLt  lintComparableKey = "lt"
	lintCpLte lintComparableKey = "lte"
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
	lintRequiredRangeAllowKeyMap = map[lintComparableKey]bool{
		lintCpGt:  true,
		lintCpGte: true,
		lintCpLt:  true,
		lintCpLte: true,
	}
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

type decorArg struct {
	index int
	name,
	typ string
	// decor lint rule
	required *requiredLinter
	nonzero  bool
}

func (d *decorArg) typeKind() types.BasicInfo {
	if t, ok := decorOptionParamTypeMap[d.typ]; ok {
		return t
	}
	return types.IsUntyped
}

type decorArgsMap map[string]*decorArg

type paramLint interface {
	valid(in string) bool
}

type lintAllowType interface {
	string | float64 | bool
}

type requiredLinter struct {
	//	gt,
	//	gte,
	//	le,
	//	lte
	compare map[lintComparableKey]float64
	enum    []string
}

func (r *requiredLinter) initCompare() {
	if r.compare == nil {
		r.compare = map[lintComparableKey]float64{}
	}
}
