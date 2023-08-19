# go-decorator

中文 | [English](README.md) (Translated from Chinese document)

> 请勿应用于生产环境！！！  
> 项目仍在迭代开发，仅公开测试阶段

让 Go 语言实现装饰器特性的编译链工具。

只需 `//go:decor your-decorator-function` 来注释你的函数即可。让你快速完成样板代码注入、改变函数行为、控制逻辑流程等。

装饰器的使用场景，可以类比其他语言，比如 Python、TypeScript。

`go-decorator` 是在编译时进行的装饰器注入，因此它不会破坏你项目的源文件，也不会额外在项目中生成新的`.go`文件和其他多余文件。 和 `go:generate` 生成式完全不同。

## Install

项目仍在积极开发，当前最好的安装方式是源码编译（后续稳定后会提供二进制分发）。

```shell
$ git clone https://github.com/dengsgo/go-decorator.git
$ cd go-decorator/cmd/decorator
$ go build
```
编译成功会得到 `decorator` 二进制文件。将此文件路径加入你的环境变量 `Path` 中以便后续调用。

也可以直接 `go install`:
```shell
$ go install github.com/dengsgo/go-decorator/cmd/decorator@latest
```

运行 `decorator --help`，出现`decorator`的帮助信息即为安装成功。

注意：请经常更新以安装最新版本。获得最佳体验。

## Usage

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

[**single**](example/single): 这个一个单文件示例，装饰器定义和被装饰的函数都位于一个包内。这种情况无需考虑导入依赖包的问题，按示例代码使用即可。

[**packages**](example/packages)：该项目示例为装饰器定义和被装饰的函数不在同一个包内，需要使用匿名包导入。

更多示例和手册陆续添加中...

## Requirement

使用该工具必须满足：

- go 1.18 及其以上  
- go.mod 项目

## Issue

发现任何问题，都可以在这里反馈. [github issues](https://github.com/dengsgo/go-decorator/issues)  

## Contribute

项目仍在开发中，由于变动频繁，暂时不接受外部贡献。欢迎稳定后再提交 Pull Request .