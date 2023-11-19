# go-decorator

[![Actions](https://github.com/dengsgo/go-decorator/workflows/Go/badge.svg)](https://github.com/dengsgo/go-decorator/actions)  [![Go Report Card](https://goreportcard.com/badge/github.com/dengsgo/go-decorator)](https://goreportcard.com/report/github.com/dengsgo/go-decorator)  [![godoc.org](https://godoc.org/github.com/dengsgo/go-decorator/decor?status.svg)](https://godoc.org/github.com/dengsgo/go-decorator/decor)  [![Go doc](https://img.shields.io/badge/go.dev-reference-brightgreen?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/dengsgo/go-decorator/decor)  [![Goproxy.cn](https://goproxy.cn/stats/github.com/dengsgo/go-decorator/badges/download-count.svg)](https://goproxy.cn)


中文 | [English](README.md)  

> 公开测试版本，谨慎用于生产环境. 欢迎 ⭐star 关注项目进展

`go-decorator`, Go 便捷使用装饰器的工具，装饰器能够切面(AOP)、代理(Proxy)任意的函数和方法，提供观察和控制函数的能力。


## Feature

- 添加注释 `//go:decor F` 即可使用装饰器（`F` 为装饰器函数），快速完成“样板代码注入、非侵入式改变函数行为、控制逻辑流程”等逻辑；  
- 可以自由定义函数作为装饰器，应用于任意一级函数和方法上（top-level function or method）;
- 支持使用多个（行） `//go:decor` 装饰器装饰目标函数;
- 支持注释 `type T types` 类型声明，decorator 会自动装饰代理 `T` 类型下的所有方法（即将上线）；  
- 装饰器支持可选参数，给开发带来更多可能；
- 支持编译时 `lint` 验证，保证 Go 编译代码的健壮性；
- 提供友好的错误提示，可在编译时发现问题并给出错误原因和错误行号（例如未定义的装饰器或未引用的包等）;   
- 仅在编译时增强目标函数，不会降低编译后程序的性能，亦没有反射操作;   

装饰器的使用场景，可以类比其他语言，比如 Python、TypeScript。（非常适合在缓存、鉴权、日志等场景使用，作为辅助手段解放重复编码的困扰）。

> `go-decorator` 是一种编译时代码注入技术。使用它不会影响您项目的源文件，并且不会在项目中生成额外的 `.go` 文件和其他冗余文件。这种注入方法与 `go:generate` 生成方式截然不同。


## Guide

查看： [中文文档](GUIDE.zh_cn.md#使用引导)  |  [English Guide](GUIDE.md#guide) 

## Install

通过 `go install` 安装:
```shell
$ go install github.com/dengsgo/go-decorator/cmd/decorator@latest
```

运行 `decorator`，显示 `decorator` 版本信息即为安装成功。
```shell
$ decorator
decorator 0.15.0 beta , https://github.com/dengsgo/go-decorator
```

提示：经常运行上述安装命令来安装最新版本，以获得 BUG 修复、增强体验和更多的新特性。

## Usage

`decorator` 依赖原生 `go` 命令来调用它，使用只需在 `go` 的子命令中加入 `-toolexec decorator` 参数即可。
例如：  

|原生命令| 使用 `decorator` |
|--------|--------|
| `go build` | `go build -toolexec decorator` |
| `go run main.go` | `go run -toolexec decorator main.go` |
| `go test -v` | `go test -toolexec decorator -v` |
| `go install` | `go install -toolexec decorator` |
| `go ... -flags...` | `go ... -toolexec decorator -flags...` |


## Code

在你的项目引入装饰器依赖（必须是 go.mod 项目）:

```shell
$ go get -u github.com/dengsgo/go-decorator
```

编写类似代码：

```go
package main

import (
	"github.com/dengsgo/go-decorator/decor"
	"log"
)

func main() {
	// 正常调用你的函数。
	// 由于这是一个声明使用装饰器logging的函数, 
	// decorator 编译链会在编译代码时注入装饰器方法logging的调用。
	// 所以使用上面的方式编译后运行，你会得到如下输出：
	// 
	// 2023/08/13 20:26:30 decorator function logging in []
	// 2023/08/13 20:26:30 this is a function: myFunc
	// 2023/08/13 20:26:30 decorator function logging out []
	// 
	// 而不是只有 myFunc 本身的一句输出。
	// 也就是说通过装饰器改变了这个方法的行为！
	myFunc() 
}

// 通过使用 go:decor 注释声明该函数将使用装饰器logging来装饰。
//
//go:decor logging
func myFunc() {
	log.Println("this is a function: myFunc")
}

// 这是一个普通的函数
// 但是它实现了 func(*decor.Context [, ...any]) 类型，因此它还是一个装饰器方法，
// 可以在其他函数上使用这个装饰器。
// 在函数中，ctx 是装饰器上下文，可以通过 ctx 获取到目标函数的函数名、出入参
// 以及执行目标函数。
// 如果函数中没有执行 ctx.TargetDo(), 那么意味着目标函数不会执行，
// 即使你代码里调用了被装饰的目标函数！这时候，目标函数返回的都是零值。
// 在 ctx.TargetDo() 之前，可以修改 ctx.TargetIn 来改变入参值。
// 在 ctx.TargetDo() 之后，可以修改 ctx.TargetOut 来改变返回值。
// 只能改变出入参的值。不要试图改变他们的类型和数量，这将会引发运行时 panic !!!
func logging(ctx *decor.Context) {
	log.Println("decorator function logging in", ctx.TargetIn)
	ctx.TargetDo()
	log.Println("decorator function logging out", ctx.TargetOut)
}

```

带有可选参数的装饰器用法：

```go
package main

import (
	"github.com/dengsgo/go-decorator/decor"
)

func main()  {
	optionalParametersFuncDemo()
}

// 它也是装饰器，不同于普通装饰器，它允许目标函数额外提供参数 `level` 。
// decorator 提供了 lint 语法给开发者，在编译代码时强制进行校验。比如:
// `required` 要求目标函数必须对该字段传值；
// `nonzero` 要求目标函数传值不能时空值。
// 如果编译时校验不通过，会编译失败。
// 使用方式如下:(更多用法查看 Guide.md)：
//
//go:decor-lint required: {level}
//go:decor-lint nonzero: {level}
func levelLogging(ctx *decor.Context, level string)  {
	if level == "debug" {
		// to do something
	}
	ctx.TargetDo()
}

// 这个方法使用了装饰器 levelLogging，并且额外传递了 `level` 参数值 "debug" 给装饰器。
// 
//go:decor levelLogging#{level: "debug"}
func optionalParametersFuncDemo()  {
	// function code
}
```

## Example

[example/usages](example/usages) Example 项目示范了如何正确编写代码使用 go-decorator 工具。 

```go
func main() {
	section("inner.go")
	// 示例：使用包内装饰器的函数
	useScopeInnerDecor("hello, world", 100)

	section("external.go")
	// 示例：使用其他包装饰器的函数
	useExternalaDecor()
	g.PrintfLn("plus(2, 3) = %+v", plus(2, 3))

	section("datetime.go")
	// 示例：文档 Guide.md 中演示使用装饰器的代码
	{
		t := 1692450000
		s := datetime(t)
		g.Printf("datetime(%d)=%s\n", t, s)
	}

	section("genericfunc.go")
	// 示例：泛型函数使用装饰器
	g.PrintfLn("Sum(1, 2, 3, 4, 5, 6, 7, 8, 9) = %+v", Sum(1, 2, 3, 4, 5, 6, 7, 8, 9))

	section("method.go")
	// 方法使用装饰器
	{
		m := &methodTestPointerStruct{}
		m.doSomething("main called")
	}
	{
		m := methodTestRawStruct{}
		m.doSomething("main called")
	}
	
	section("withdecorparams.go")
	// 示例：使用带有参数的装饰器，如何传值
	g.PrintfLn("useArgsDecor() = %+v", useArgsDecor())
	// 示例：装饰器如何使用 Lint 在编译时约束验证目标函数的参数
	g.Printf("useHitUseRequiredLint() = %+v", useHitUseRequiredLint())
	g.Printf("useHitUseNonzeroLint() = %+v", useHitUseNonzeroLint())
	g.Printf("useHitBothUseLint() = %+v", useHitBothUseLint())
	g.Printf("useHitUseMultilineLintDecor() = %+v", useHitUseMultilineLintDecor())
}

```

详细文档查看 [Guide](#guide) .

## Requirement

使用该工具必须满足：

- go 1.18 及其以上  
- go.mod 项目

## Issue

发现任何问题，都可以在这里反馈. [github issues](https://github.com/dengsgo/go-decorator/issues)  

## Contribute

项目仍在开发中，由于变动频繁，暂时不接受外部贡献。欢迎稳定后再提交 Pull Request .

## TODO

- More documents.
- IDE friendly tool support.  
- better performance.
- More testing coverage.  
- More clear error reminders.
- More...