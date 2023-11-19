package main

import (
	"github.com/dengsgo/go-decorator/example/usages/g"
	"log"
	"time"
)

func main() {
	section("inner.go")
	// 这是一个使用包内装饰器的函数
	useScopeInnerDecor("hello, world", 100)

	section("external.go")
	// 这是一个使用其他包装饰器的函数
	useExternalaDecor()
	g.PrintfLn("plus(2, 3) = %+v", plus(2, 3))

	section("datetime.go")
	// 文档 Guide.md 中演示使用装饰器的代码
	{
		t := 1692450000
		s := datetime(t)
		g.Printf("datetime(%d)=%s\n", t, s)
	}

	section("genericfunc.go")
	// 泛型函数使用装饰器
	g.PrintfLn("Sum(1, 2, 3, 4, 5, 6, 7, 8, 9) = %+v", Sum(1, 2, 3, 4, 5, 6, 7, 8, 9))

	section("method.go")
	// 结构体方法使用装饰器
	{
		m := &methodTestPointerStruct{}
		m.doSomething("main called")
	}
	{
		m := methodTestRawStruct{}
		m.doSomething("main called")
	}

	section("withdecorparams.go")
	// 使用带有参数的装饰器，如何传值
	g.PrintfLn("useArgsDecor() = %+v", useArgsDecor())
	// 装饰器如何使用 Lint 在编译时约束验证目标函数的参数
	g.Printf("useHitUseRequiredLint() = %+v", useHitUseRequiredLint())
	g.Printf("useHitUseNonzeroLint() = %+v", useHitUseNonzeroLint())
	g.Printf("useHitBothUseLint() = %+v", useHitBothUseLint())
	g.Printf("useHitUseMultilineLintDecor() = %+v", useHitUseMultilineLintDecor())

	section("types.go")
	// 注释 `type T types` 类型声明，decorator 会自动装饰代理 `T` 类型下的所有方法：
	{
		// 注释结构体的用法
		s := &structType{"main say hello"}
		g.PrintfLn("s.Name() = %+v", s.Name())
		s.StrName("StrName() set")
		g.PrintfLn("s.Name() = %+v", s.Name())
		s.empty()
	}
	{
		// 注释基础类型的用法
		v := varIntType(100)
		g.PrintfLn("v.value() = %+v", v.value())
		v.zeroSelf()
		g.PrintfLn("v.value() = %+v", v.value())
	}
	{
		// 注释基础类型的用法
		v := VarStringType("hello")
		g.PrintfLn("v.value() = %+v", v.value())
	}
	{
		// 注释结构体的用法，该类型没有任何方法
		v := nonMethodType{}
		g.PrintfLn("nonMethodType = %+v", v)
	}
	{
		// 注释结构体的用法，该类型的方法定义在包内的其他文件里
		o := &otherFileDefMethodType{}
		g.PrintfLn("o.string() = %+v", o.string())
	}
}

func section(s string) {
	g.PrintfLn("\n++++++++++ " + s + " ++++++++++")
}

func init() {
	log.SetFlags(0)
	time.Local = time.FixedZone("CST", 8*3600)
}
