package main

import (
	"github.com/dengsgo/go-decorator/example/usages/g"
	"log"
	"strings"
	"testing"
)

func TestStructType(t *testing.T) {
	s := &structType{"main say hello"}
	g.PrintfLn("s.Name() = %+v", s.Name())
	s.StrName("StrName() set")
	g.PrintfLn("s.Name() = %+v", s.Name())
	s.empty()
	out := strings.TrimSpace(g.TestBuffers.String())
	r := `dumpTargetType say: Receiver: &{name:main say hello}, TargetName: Name
structType: main say hello
s.Name() = main say hello
dumpTargetType say: Receiver: &{name:main say hello}, TargetName: StrName
dumpTargetType say: Receiver: &{name:StrName() set}, TargetName: Name
structType: StrName() set
s.Name() = StrName() set
dumpTargetType say: Receiver: &{name:StrName() set}, TargetName: empty`
	if out != r {
		t.Fatalf("TestStructType fail, out : %s, \nshould : %s", out, r)
	}
	g.ResetTestBuffers()
}

func TestVarIntType(t *testing.T) {
	v := varIntType(100)
	g.PrintfLn("v.value() = %+v", v.value())
	v.zeroSelf()
	g.PrintfLn("v.value() = %+v", v.value())
	out := strings.TrimSpace(g.TestBuffers.String())
	r := `dumpTargetType say: Receiver: 100, TargetName: value
v.value() = 100
dumpTargetType say: Receiver: 100, TargetName: zeroSelf
dumpTargetType say: Receiver: 100, TargetName: value
v.value() = 100`
	if out != r {
		t.Fatalf("TestVarIntType fail, out : %s, \nshould : %s", out, r)
	}
	g.ResetTestBuffers()
}

func TestVarStringType(t *testing.T) {
	v := VarStringType("hello")
	s := v.value()
	g.PrintfLn("v.value() = %v", s)
	out := strings.TrimSpace(g.TestBuffers.String())
	r := `dumpTargetType say: Receiver: hello, TargetName: value
v.value() = hello`
	if out != r {
		t.Fatalf("TestVarStringType fail, out : %s, \nshould : %s", out, r)
	}
	g.ResetTestBuffers()
}

func TestNonMethodType(t *testing.T) {
	v := nonMethodType{}
	g.PrintfLn("nonMethodType = %+v", v)
	out := strings.TrimSpace(g.TestBuffers.String())
	r := `nonMethodType = {}`
	if out != r {
		t.Fatalf("TestNonMethodType fail, out : %s, \nshould : %s", out, r)
	}
	g.ResetTestBuffers()
}

func TestOtherFileDefMethodType(t *testing.T) {
	o := &otherFileDefMethodType{}
	g.PrintfLn("o.string() = %+v", o.string())
	out := strings.TrimSpace(g.TestBuffers.String())
	r := `dumpTargetType say: Receiver: &{}, TargetName: string
o.string() = otherFileDefMethodType string()`
	if out != r {
		t.Fatalf("TestOtherFileDefMethodType fail, out : %s, \nshould : %s", out, r)
	}
	g.ResetTestBuffers()
}

func TestGenericType(t *testing.T) {
	t.Run("Part1", func(t *testing.T) {
		genInt := &genericType[int]{}
		g.PrintfLn("genInt.value() = %+v", genInt.value())
		genStr := &genericType[string]{}
		g.PrintfLn("genStr.value() = %+v", genStr.value())
		//genBool := &genericType[bool]{}
		//g.PrintfLn("genBool.value() = %+v", genBool.value())
		//genStruct := &genericType[struct{}]{}
		//g.PrintfLn("genStruct.value() = %+v", genStruct.value())
		out := strings.TrimSpace(g.TestBuffers.String())
		r := `dumpTargetType say: Receiver: &{t:0}, TargetName: value
genInt.value() = 0
dumpTargetType say: Receiver: &{t:}, TargetName: value
genStr.value() =`
		if out != r {
			t.Fatalf("TestGenericType fail, out : %s, \nshould : %s", out, r)
		}
		g.ResetTestBuffers()
	})

	t.Run("Part2", func(t *testing.T) {
		genBool := &genericType[bool]{}
		g.PrintfLn("genBool.value() = %+v", genBool.value())
		genStruct := &genericType[struct{}]{}
		g.PrintfLn("genStruct.value() = %+v", genStruct.value())
		out := strings.TrimSpace(g.TestBuffers.String())
		log.Println(out)
		r := `dumpTargetType say: Receiver: &{t:false}, TargetName: value
genBool.value() = false
dumpTargetType say: Receiver: &{t:{}}, TargetName: value
genStruct.value() = {}`
		if out != r {
			t.Fatalf("TestGenericType1 fail, out : %s, \nshould : %s", out, r)
		}
		g.ResetTestBuffers()
	})

}
