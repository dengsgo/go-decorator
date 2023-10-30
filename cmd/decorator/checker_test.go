package main

import (
	"errors"
	"log"
	"testing"
)

func TestCheckDecorAndGetParam(t *testing.T) {
	param, err := checkDecorAndGetParam("github.com/dengsgo/go-decorator/decor", "find", nil)
	log.Println(param, err)
	param, err = checkDecorAndGetParam("github.com/dengsgo/go-decorator/cmd/decorator", "logging", nil)
	log.Println(param, err)
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
