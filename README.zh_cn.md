# go-decorator

[![Actions](https://github.com/dengsgo/go-decorator/workflows/Go/badge.svg)](https://github.com/dengsgo/go-decorator/actions)  [![Go Report Card](https://goreportcard.com/badge/github.com/dengsgo/go-decorator)](https://goreportcard.com/report/github.com/dengsgo/go-decorator)  [![godoc.org](https://godoc.org/github.com/dengsgo/go-decorator/decor?status.svg)](https://godoc.org/github.com/dengsgo/go-decorator/decor)  [![Go doc](https://img.shields.io/badge/go.dev-reference-brightgreen?logo=go&logoColor=white&style=flat)](https://pkg.go.dev/github.com/dengsgo/go-decorator/decor)  [![Goproxy.cn](https://goproxy.cn/stats/github.com/dengsgo/go-decorator/badges/download-count.svg)](https://goproxy.cn)


中文 | [English](README.md)  

> 请勿应用于生产环境！！！  
> 项目仍在迭代开发，仅公开测试阶段

`go-decorator` 是用于 Go 语言编译器的中间件工具，能够通过注释方式实现非侵入式的装饰器使用。


## Feature

- 使用 `//go:decor decoratorfunctionName` 注释函数即可使用装饰器 `decoratorfunctionName`，快速完成样板代码注入、非侵入式改变函数行为、控制逻辑流程等；  
- 定义 `func(*decor.Context)` 类型的函数，即可作为装饰器，应用于任意一级函数（top-level function）;  
- 支持使用多个（行）`//go:decor`装饰器来装饰函数;  
- 提供友好的错误提示，可在编译时发现问题并给出错误原因和错误行号（例如未定义的装饰器或未引用的包等）;   
- 仅在编译时增强目标函数，不会降低编译后程序的性能，也没有反射操作;  
- 提供基础使用指南。  

装饰器的使用场景，可以类比其他语言，比如 Python、TypeScript。（非常适合在缓存、鉴权、日志等场景使用，作为辅助手段解放重复编码的困扰）。

`go-decorator` 是一种编译时装饰器注入技术。使用它不会影响您项目的源文件，并且不会在项目中生成额外的 `.go` 文件或其他冗余文件。这种注入方法与 `go:generate` 生成方式截然不同。


## Guide

查看： [中文文档](GUIDE.zh_cn.md#使用引导)  |  [English Guide](GUIDE.md#guide)  |  More

## Install

通过 `go install` 安装:
```shell
$ go install github.com/dengsgo/go-decorator/cmd/decorator@latest
```

运行 `decorator`，显示 `decorator` 版本信息即为安装成功。

注意：请经常更新以安装最新版本。获得最佳体验。

## Usage

`decorator` 是 `go` 的编译链工具，依靠 `go` 命令来调用它运行，进行代码的编译。

在 `go build` 命令中加入 `-toolexec 'decorator'` 参数即可。

假如你平时就是使用 `go build`,那么现在只需要加上工具链参数变成 `go build -toolexec 'decorator'`。其他一切和以往一样，无需做任何更改！

go 的其他子命令也是同样的使用方法。

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
// 但是它实现了 func(*decor.Context) 类型，因此它还是一个装饰器方法，
// 可以在其他函数上使用这个装饰器。
// 在函数中，ctx 是装饰器上下文，可以通过 ctx 获取到目标函数的出入参
// 和目标方法的执行。
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

## Example

[example](example)这个目录示范了如何正确编写代码来使用 go-decorator 工具。

| Project                           | Notes                                                     |
|-----------------------------------|-----------------------------------------------------------|
| [**single**](example/single)      | 这个一个单文件示例，装饰器定义和被装饰的函数都位于一个包内。这种情况无需考虑导入依赖包的问题，按示例代码使用即可。 | 
| [**packages**](example/packages)  | 该项目示例为装饰器定义和被装饰的函数不在同一个包内，需要使用匿名包导入。                      |
| [**datetime**](example/datetime)  | Guide 里演示示例所用到的完整代码                                       |
| [**emptyfunc**](example/emptyfunc) | 演示装饰器中调用和不调用`TargetDo()` 的区别                              |


更多内容查看 [Guide](#guide) .

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