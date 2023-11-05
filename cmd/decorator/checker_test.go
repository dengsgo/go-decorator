package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"log"
	"testing"
)

func TestCheckDecorAndGetParam(t *testing.T) {
	cas := []struct {
		in map[string]string
		r  []string
	}{
		{
			map[string]string{"s": `"value"`},
			[]string{`"value"`, "0", "false"},
		},
		{
			map[string]string{"a": "11111"},
			[]string{`""`, "11111", "false"},
		},
		{
			map[string]string{"b": "true"},
			[]string{`""`, "0", "true"},
		},
		{
			map[string]string{"s": `"value"`, "a": "0", "b": "true"},
			[]string{`"value"`, "0", "true"},
		},
		{
			map[string]string{"a": "0", "s": `"value"`, "b": "true"},
			[]string{`"value"`, "0", "true"},
		},
		{
			map[string]string{"b": "true", "a": "0", "s": `"value"`},
			[]string{`"value"`, "0", "true"},
		},
		{
			map[string]string{"b": "false", "a": "0", "s": `"kkkk"`},
			[]string{`"kkkk"`, "0", "false"},
		},
	}

	targetPkg := "github.com/dengsgo/go-decorator/cmd/decorator"
	for index, c := range cas {
		param, err := checkDecorAndGetParam(targetPkg,
			"logging", c.in)
		if err != nil {
			t.Fatal("checkDecorAndGetParam should err == nil but got error", err)
		}
		for i, v := range c.r {
			if param[i] != v {
				t.Fatalf("checkDecorAndGetParam should param == r but got: %s != %s, case index: %+v, i: %+v",
					param[i], v, index, i)
			}
		}
	}

	_, err := checkDecorAndGetParam("github.com/dengsgo/go-decorator/decor", "find", nil)
	if err == nil {
		t.Fatal("checkDecorAndGetParam should return err but got nil")
	}

	// TODO
	//failed := []map[string]string{
	//	{"s": `value`, "a": "0", "b": "true"},
	//	{"s": `if`, "a": "0", "b": "true"},
	//	{"s": `"value"`, "a": "0.0", "b": "true"},
	//	{"s": `"value"`, "a": "0", "b": "true1"},
	//}
	//for i, v := range failed {
	//	_, err := checkDecorAndGetParam(targetPkg, "logging", v)
	//	if err == nil {
	//		t.Fatal("checkDecorAndGetParam should return err but got nil, index: ", i)
	//	}
	//}
}

func TestCleanSpaceChar(t *testing.T) {
	cas := []struct {
		s,
		r string
	}{
		{"helloworld", "helloworld"},
		{"hello world", "helloworld"},
		{"hello 世界", "hello世界"},
		{" he l l owo      rld    ", "helloworld"},
		{"hello 世 界  这是测		试\t用     例 	  ", "hello世界这是测试用例"},
		{" 😀/(ㄒoㄒ)/~ ~   😊😄	😄\v😄  😄😄😄  ", "😀/(ㄒoㄒ)/~~😊😄😄😄😄😄😄"},
		{"if a > 1 {\necho ''\n}", "ifa>1{echo''}"},
	}
	for i, v := range cas {
		if cleanSpaceChar(v.s) != v.r {
			t.Fatal("cleanSpaceChar(v.s)!=r, pos", i, ": ", cleanSpaceChar(v.s), "!=", v.r)
		}
	}
}

func TestIsLetters(t *testing.T) {
	cas := []struct {
		s string
		r bool
	}{
		{"thisisastring", true},
		{"this isastring", false},
		{"thisisastring ", false},
		{" thisisastring", false},
		{"这是string", true},
		{"这 是string", false},
		{"这是 string", false},
		{"这是string\t", false},
		{"这是\vstring", false},
		{"\n这是string", false},
		{"thisisa字符串", true},
		{"", false},
		{"\r", false},
		{"😀/(ㄒoㄒ)/~~😊😄😄😄😄😄😄", false},
		{" 😀/(ㄒoㄒ)/~ ~   😊😄	😄\v😄  😄😄😄  ", false},
	}
	for i, v := range cas {
		if isLetters(v.s) != v.r {
			t.Fatal("isLetters(v.s)!=r, pos", i, ": ", v.s, isLetters(v.s), "!=", v.r)
		}
	}
}

func TestParseDecorAndParameters(t *testing.T) {
	cas := []struct {
		s        string
		callName string
		params   map[string]string
	}{
		{"function", "function", map[string]string{}},
		{"fun.DO", "fun.DO", map[string]string{}},
		{"fun.DO#{}", "fun.DO", map[string]string{}},
		{"a.b.c.d.DO#{}", "a.b.c.d.DO", map[string]string{}},
		{"function#{}", "function", map[string]string{}},
		{`function#{key:""}`, "function", map[string]string{"key": `""`}},
		{`function#{age:100}`, "function", map[string]string{"age": "100"}},
		{`function#{f:0.110}`, "function", map[string]string{"f": "0.110"}},
		{`function#{b:true}`, "function", map[string]string{"b": "true"}},
		{`function#{b:true, key:"", f:0.110, age:100}`, "function", map[string]string{"b": "true", "key": `""`, "age": "100", "f": "0.110"}},
		{`function#{b:true, key:"", f:0.110, age:100,   }`, "function", map[string]string{"b": "true", "key": `""`, "age": "100", "f": "0.110"}},
		{`function#{   b:true, key:"", f:0.110, age:100}`, "function", map[string]string{"b": "true", "key": `""`, "age": "100", "f": "0.110"}},
		{`function#{   b:true, key:"", f:0.110, age:100   }`, "function", map[string]string{"b": "true", "key": `""`, "age": "100", "f": "0.110"}},
		{`function#{   b:true, key:"", f:0.110, age:100   }   `, "function", map[string]string{"b": "true", "key": `""`, "age": "100", "f": "0.110"}},
		{`function #{   b:true, key:"", f:0.110, age:100   }   `, "function", map[string]string{"b": "true", "key": `""`, "age": "100", "f": "0.110"}},
		{`function # {   b:true, key:"", f:0.110, age:100   }   `, "function", map[string]string{"b": "true", "key": `""`, "age": "100", "f": "0.110"}},
	}
	for _, v := range cas {
		name, p, err := parseDecorAndParameters(v.s)
		if err != nil {
			log.Fatalf("parseDecorAndParameters(v.s) parse error, err: %+v, case: %s, callName: %+v, params: %+v,\n",
				err, v.s, v.callName, v.params)
		}
		if name != v.callName {
			log.Fatalf("parseDecorAndParameters(v.s) parse ok but callName failed, case: %s, callName: %+v, params: %+v,\n",
				v.s, v.callName, v.params)
		}
		if v.params == nil {
			log.Fatalf("parseDecorAndParameters(v.s) parse ok but v.params == nil, case: %s, callName: %+v, params: %+v,\n",
				v.s, v.callName, v.params)
		}
		if len(v.params) != len(p) {
			log.Fatalf("parseDecorAndParameters(v.s) parse ok but len(v.params) != len(p), case: %s, callName: %+v, params: %+v,\n",
				v.s, v.callName, v.params)
		}
		for k, value := range v.params {
			if _v, ok := p[k]; ok && _v == value {
				continue
			}
			log.Fatalf("parseDecorAndParameters(v.s) parse ok but v.params keyOrValue not exist, key:%s, value:%s, case: %s, callName: %+v, params: %+v,\n",
				k, value, v.s, v.callName, v.params)
		}
	}

	failed := []struct {
		s   string
		err error
	}{
		{"", errUsedDecorSyntaxErrorLossFunc},
		{"      ", errUsedDecorSyntaxError},
		{"     f f ", errUsedDecorSyntaxError},
		{"{k:v}", errUsedDecorSyntaxError},
		{"{k:}", errUsedDecorSyntaxError},
		{"{k}", errUsedDecorSyntaxError},
		{"{}", errUsedDecorSyntaxError},
		{"{", errUsedDecorSyntaxError},
		{"#", errUsedDecorSyntaxError},
		{"#####", errUsedDecorSyntaxError},
		{"function#", errUsedDecorSyntaxError},
		{"function##", errUsedDecorSyntaxError},
		{"function #", errUsedDecorSyntaxError},
		{"function ##", errUsedDecorSyntaxError},
		{"function #  {", errUsedDecorSyntaxError},
		{"function #  }", errUsedDecorSyntaxError},
		{"function #  }{", errUsedDecorSyntaxError},
		{"function{}", errUsedDecorSyntaxError},
		{"function#{{}}", errUsedDecorSyntaxError},
		{"function{}#", errUsedDecorSyntaxError},
		{"function#{#}", errUsedDecorSyntaxError},
		{"function#{key}", errUsedDecorSyntaxErrorInvalidP},
		{"function#{key:}", errUsedDecorSyntaxErrorInvalidP},
		{"function#{k ey:}", errUsedDecorSyntaxErrorInvalidP},
		{"function#{key：}", errUsedDecorSyntaxErrorInvalidP},
		{"function#{:}", errUsedDecorSyntaxError},
		{"function#{:value}", errUsedDecorSyntaxErrorInvalidP},
		{"function#{:val ue}", errUsedDecorSyntaxErrorInvalidP},
		{"function#{:val#ue}", errUsedDecorSyntaxErrorInvalidP},
		{"function#{:va\"l#ue}", errUsedDecorSyntaxErrorInvalidP},
		{"function#{key:vv v}", errUsedDecorSyntaxErrorInvalidP},
		{"function#{key:vv v, ,}", errUsedDecorSyntaxErrorInvalidP},
		{"function#{key:vv v, ssd,}", errUsedDecorSyntaxErrorInvalidP},
		{"function#{key:true1,s:false,}", errors.New("invalid parameter value, should be bool")},
		{"function#{key:vv,key:vv,}", errors.New("invalid parameter value, should be bool")},
		{`function#{name:"vv",name:"vvccc"}`, errors.New("duplicate parameters key 'name'")},
		{"function#{key:vv,keys:vv,,,}", errUsedDecorSyntaxErrorInvalidP},
		{"function#{,,,key:vv,keys:vv,,,}", errUsedDecorSyntaxErrorInvalidP},
		{"function#{Name:$}", errUsedDecorSyntaxErrorInvalidP},
		{"function#{Name:<>}", errUsedDecorSyntaxErrorInvalidP},
		{"function#{Name:<>},", errUsedDecorSyntaxError},
		{`function#""`, errUsedDecorSyntaxError},
		{`function#{""}`, errUsedDecorSyntaxError},
		{`function#{"}`, errUsedDecorSyntaxError},
		{`function#{"Name"}`, errUsedDecorSyntaxErrorInvalidP},
		{`function#{"Name":""}`, errors.New("invalid parameter name")},
		{`function#{"Name"=""}`, errUsedDecorSyntaxErrorInvalidP},
		{`function#{key=""}`, errUsedDecorSyntaxErrorInvalidP},
		{`function#{key:=""}`, errUsedDecorSyntaxErrorInvalidP},
		{`function#{key:if}`, errUsedDecorSyntaxErrorInvalidP},
		{`function#{for:if}`, errUsedDecorSyntaxErrorInvalidP},
		{`function#{for:true}`, errUsedDecorSyntaxErrorInvalidP},
		{".DO#{}", errUsedDecorSyntaxError},
		{"a.b.c.#{}", errUsedDecorSyntaxError},
		{"a,b.c.#{}", errUsedDecorSyntaxError},
	}
	for i, v := range failed {
		_, _, err := parseDecorAndParameters(v.s)
		if err == nil {
			log.Fatalf("parseDecorAndParameters(v.s) should be fail but pass, case: %s\n",
				v.s)
		}
		if err.Error() != v.err.Error() {
			log.Fatalf("parseDecorAndParameters(v.s) err not match case, i:%+v, err: %+v, should: %+v, case: %s\n",
				i, err, v.err, v.s)
		}
	}
}

func TestParseLinterFromAnnotation(t *testing.T) {
	args := decorArgsMap{
		"name":     &decorArg{1, "name", "string", nil, false},
		"intVal":   &decorArg{2, "intVal", "int", nil, false},
		"floatVal": &decorArg{3, "floatVal", "float64", nil, false},
		"boolVal":  &decorArg{4, "boolVal", "bool", nil, false},
		"rangeVal": &decorArg{4, "rangeVal", "int64", nil, false},
		"emptyVal": &decorArg{5, "emptyVal", "string", nil, false},
	}
	cas := []string{
		`required: {intVal}`,
		`required: {name, floatVal}`,
		`required: {name: {}, floatVal: {}}`,
		`required: {name: {"a"}, intVal: {gt:1}, floatVal: {lte:0}}`,
		`required: {name: {"a", "b", lte:10}, intVal: {1,2,3,4,5,10}, floatVal: {lte:0, 999.0}}`,
		`required: {boolVal: {false}}`,
		`required: {rangeVal: {100,200, gt:1, lt:999, gte: 0, lte: 1000}}`,
		`nonzero: {name}`,
		`nonzero: {name, floatVal}`,
		`nonzero: {name, intVal}`,
		`nonzero: {name, intVal, rangeVal}`,
	}
	for i, v := range cas {
		err := parseLinterFromAnnotation(v, args)
		if err != nil {
			log.Fatalf("parseLinterFromAnnotation(s) should pass, i:%+v, err: %+v, case: %s\n",
				i, err, v)
		}
	}
	result := map[string]string{
		"name":     `&{compare:map[lte:10] enum:["a" "a" "b"]}`,
		"intVal":   `&{compare:map[gt:1] enum:[1 2 3 4 5 10]}`,
		"floatVal": `&{compare:map[lte:0] enum:[999.0]}`,
		"boolVal":  `&{compare:map[] enum:[false]}`,
		"rangeVal": `&{compare:map[gt:1 gte:0 lt:999 lte:1000] enum:[100 200]}`,
		"emptyVal": `<nil>`,
	}
	for k, v := range args {
		fmt.Printf("k:%s, v:%+v\n", k, v.required)
		if fmt.Sprintf("%+v", v.required) != result[k] {
			log.Fatalf("parseLinterFromAnnotation(s) required result should pass, k:%+v, result[k]:%s, case: %+v\n",
				k, result[k], v.required)
		}
	}

	result = map[string]string{
		"name":     `true`,
		"intVal":   `true`,
		"floatVal": `true`,
		"boolVal":  `false`,
		"rangeVal": `true`,
		"emptyVal": `false`,
	}
	for k, v := range args {
		fmt.Printf("k:%s, v:%+v\n", k, v.nonzero)
		if fmt.Sprintf("%+v", v.nonzero) != result[k] {
			log.Fatalf("parseLinterFromAnnotation(s) nonzero result should pass, k:%+v, result[k]:%s, case: %+v\n",
				k, result[k], v.nonzero)
		}
	}
}

func TestA(t *testing.T) {
	s := `map[any]any{a, b:{"str", 1, 1.0, true, gte: 1}, c}`
	a, err := parser.ParseExpr(s)
	log.Println(err)
	ast.Print(token.NewFileSet(), a)
}
