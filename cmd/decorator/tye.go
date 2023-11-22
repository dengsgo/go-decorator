package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/types"
	"strconv"
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

type mapV[K comparable, V any] struct {
	items map[K]V
}

func newMapV[K comparable, V any]() *mapV[K, V] {
	return &mapV[K, V]{
		items: make(map[K]V),
	}
}

func (m *mapV[K, V]) put(key K, value V) bool {
	if _, ok := m.items[key]; ok {
		return false
	}
	m.items[key] = value
	return true
}

func (m *mapV[K, V]) get(key K) (v V) {
	if v, ok := m.items[key]; ok {
		return v
	}
	return
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

func (d *decorArg) passRequiredLint(value string) error {
	if d.required == nil {
		return nil
	}
	if !d.required.inEnum(value) {
		return errors.New(
			fmt.Sprintf("lint: key '%s' value '%s' can't pass lint enum", d.name, value))
	}
	if d.required.compare == nil {
		return nil
	}

	val := 0.0
	if d.typeKind() == types.IsString {
		val = float64(len(value) - 2)
	} else {
		val, _ = strconv.ParseFloat(value, 64)
	}
	compare := func(c lintComparableKey, v float64) bool {
		switch c {
		case lintCpGt:
			return val > v
		case lintCpGte:
			return val >= v
		case lintCpLt:
			return val < v
		case lintCpLte:
			return val <= v
		}
		return true
	}
	for c, v := range d.required.compare {
		if !compare(c, v) {
			return errors.New(
				fmt.Sprintf("lint: key '%s' value '%s' can't pass lint %s:%v", d.name, value, c, v))
		}
	}
	return nil
}

func (d *decorArg) passNonzeroLint(value string) error {
	isZero := func() bool {
		switch d.typeKind() {
		case types.IsInteger,
			types.IsFloat:
			value, _ := strconv.ParseFloat(value, 64)
			return value == 0
		case types.IsString:
			return value == `""`
		case types.IsBoolean:
			return value == "false"
		}
		return false
	}
	if d.nonzero && isZero() {
		return errors.New(
			fmt.Sprintf("lint: key '%s' value '%s' can't pass nonzero lint", d.name, value))
	}
	return nil
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

func (r *requiredLinter) inEnum(value string) bool {
	if r.enum == nil {
		return true
	}
	for _, v := range r.enum {
		if v == value {
			return true
		}
	}
	return false
}
