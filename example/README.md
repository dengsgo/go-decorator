## example

[usages](usages) 目录示范了如何正确编写代码来使用 go-decorator 工具。  

```go
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
	// 给 `type T types` 类型声明添加注释//go:decor F，decorator 会自动使用装饰器 F 装饰代理以 T 或者 *T 为接收者的所有方法：
	{
		// 结构体
		s := &structType{"main say hello"}
		g.PrintfLn("s.Name() = %+v", s.Name())
		s.StrName("StrName() set")
		g.PrintfLn("s.Name() = %+v", s.Name())
		s.empty()
	}
	{
		// 基础类型
		v := varIntType(100)
		g.PrintfLn("v.value() = %+v", v.value())
		v.zeroSelf()
		g.PrintfLn("v.value() = %+v", v.value())
	}
	{
		// 泛型结构体 T = int
		genInt := &genericType[int]{}
		g.PrintfLn("genInt.value() = %+v", genInt.value())
		// 泛型结构体 T = struct
		genStruct := &genericType[struct{}]{}
		g.PrintfLn("genStruct.value() = %+v", genStruct.value())
	}

	section("types_multiple.go")
	// `type T types` 和它的方法同时使用装饰器。
	// 方法上的装饰器优先执行
	{
		m := multipleStructStandType{}
		m.sayHello()
	}
	{
		m := multipleStructWrapType{}
		m.sayNiHao()
	}

}
```