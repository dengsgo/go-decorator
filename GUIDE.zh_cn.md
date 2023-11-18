# 使用引导

- [使用引导](#使用引导)
  - [要求](#要求)
  - [安装](#安装)
  - [使用](#使用)
    - [添加参数](#添加参数)
    - [添加依赖](#添加依赖)
    - [目标函数和装饰器](#目标函数和装饰器)
    - [自定义装饰器](#自定义装饰器)
    - [使用装饰器](#使用装饰器)
    - [使用多个装饰器](#使用多个装饰器)
  - [Context](#context)
    - [ctx.Kind](#ctxkind)
    - [ctx.TargetIn](#ctxtargetin)
    - [ctx.TargetOut](#ctxtargetout)
    - [ctx.TargetDo()](#ctxtargetdo)
    - [ctx.DoRef()](#ctxdoref)
  - [包引用](#包引用)
  - [条件和限制](#条件和限制)
  - [开发与调试](#开发与调试)
  - [性能](#性能)
  - [更多](#更多)


`go-decorator`, 让 Go 便捷使用装饰器的工具，装饰器能够切面(AOP)、代理(Proxy)任意的函数和方法，提供观察和控制函数的能力。

## 要求

- go1.18+
- go.mod project

## 安装

```shell
$ go install github.com/dengsgo/go-decorator/cmd/decorator@latest
```

运行 `decorator` 显示当前版本。

```shell
$ decorator
decorator 0.12.0 beta , https://github.com/dengsgo/go-decorator
```

提示：经常运行上述安装命令来安装最新版本，以获得 BUG 修复、增强体验和更多的新特性。

## 使用 

`decorator` 是 `go` 的编译链工具，依靠 `go` 命令来调用它，进行代码的编译。

### 添加参数

在 `go` 子命令中添加参数 `-toolexec decorator` 即可使用该工具。

例如：

|原生命令| 使用 `decorator` |
|--------|--------|
| `go build` | `go build -toolexec decorator` |
| `go run main.go` | `go run -toolexec decorator main.go` |
| `go test -v` | `go test -toolexec decorator -v` |
| `go install` | `go install -toolexec decorator` |
| `go ... -flags...` | `go ... -toolexec decorator -flags...` |

它适用于大多数 `go` 的子命令。

### 添加依赖

在你的项目根目录，添加 `go-decorator` 依赖。

```shell
$ go get -u github.com/dengsgo/go-decorator
```

### 了解目标函数和装饰器

> 目标函数：即使用了装饰器的函数或方法，也称为被装饰的函数、目标对象。  
> 比如 A 函数使用了装饰器 B 来装饰自己，A 即为目标函数。

装饰器也是函数，它通常用来对目标函数进行包含装饰。当代码运行到目标函数的时候，并不会真的执行它，而是运行它所使用的装饰器。实际的目标函数逻辑被包装到了装饰器中，并允许装饰器来控制它。

### 自定义装饰器

装饰器是普通的 `go` 一级函数(Top-level Function)，它的类型是 `func(*decor.Context [, ...any])`. 只要函数满足这个类型，它即是合法的装饰器，可以在项目代码里来装饰其他函数。

例如, 这里定义一个 logging 装饰器，它可以打印被调用函数的参数：

```go
package main

import "github.com/dengsgo/go-decorator/decor"

func logging(ctx *decor.Context) {
	log.Println("logging print target in", ctx.TargetIn)
	ctx.TargetDo()
	log.Println("logging print target out", ctx.TargetOut)
}
```

这个函数 `logging` 就是一个合法的装饰器，它可以用在任意的一级函数上。

关于 `ctx *decor.Context` 参数，跳转这里 [Context](#context)。

### 使用装饰器

在任意一级函数上都可以通过注释 `//go:decor ` 来使用装饰器。

例如，我们有一个函数 `datetime`, 它可以把时间戳转换成字符串日期格式。通过 `//go:decor logging` 来使用 `logging` 装饰器打印函数出入参：

```go
// Omitted code ...

// Convert timestamp to string date format.
//
//go:decor logging
func datetime(timestamp int) string {
    return time.Unix(int64(timestamp), 0).String()
}

// Omitted code ...
```

`datetime` 在编译时会被识别并注入`logging`的调用。当其他地方调用 `datetime` 函数时，就会自动执行 `logging`。

比如，我们在 `main` 入口函数中调用 `datetime`:

```go
func main() {
	t := 1692450000
	s := datetime(t)
	log.Printf("datetime(%d)=%s\n", t, s)
}
```

使用下面的命令编译、运行:

```shell
$ go build -toolexec decorator
$ ./datetime
```
将会看到如下输出：

```shell
2023/08/19 21:12:21 logging print target in [1692450000]
2023/08/19 21:12:21 logging print target out [2023-08-19 21:00:00 +0800 CST]
2023/08/19 21:12:21 datetime(1692450000)=2023-08-19 21:00:00 +0800 CST
```
我们代码中只调用了 `datetime` 函数，但是可以看到使用的 logging 装饰器也被执行了！

完整代码可以查看 [example/usages](example/usages). 

### 使用多个装饰器

`decorator` 允许同时使用多个装饰器来装饰目标函数。 

通过多个 `//go:decor ` 来使用多个装饰器。

例如，下面的 `datetime` 使用了3个装饰器,分别是 `logging`、`appendFile`、`timeFollowing`：

```go
// Omitted code ...

// Convert timestamp to string date format.
//
//go:decor logging
//go:decor appendFile
//go:decor timeFollowing
func datetime(timestamp int64) string {
    return time.Unix(timestamp, 0).String()
}

// Omitted code ...
```

如果使用了多个装饰器，装饰器执行的优先级为从上往下，也就是说先定义的先被执行。上面的装饰器中，执行顺序为 `logging` -> `appendFile` -> `timeFollowing`.

多个装饰器的使用，可能会导致代码的可读性变差，加大逻辑流程理解成本，尤其是装饰器本身的代码又特别复杂的情况。因此并不推荐这样使用。

### 带有额外参数的装饰器

顾名思义，装饰器允许定义除了第一个参数 `*decor.Context` 外的额外参数, 如：

```go
package main
import "github.com/dengsgo/go-decorator/decor"

func hit(ctx *decor.Context, msg string, count int64, repeat bool, f float64, opt string) {
	// code...
} 
```

`hit` 函数即为一个合法带有可选参数的装饰器，它允许目标函数调用时传入相应值， `hit` 函数就能获取到目标函数的参数值。

以下的参数类型是被允许的：

| 类型  | 关键字|     
|-----|-----|
| 整数  | int,int8,int16.int32,int64,unit,unit8,unit16,unit32,unit64 |
| 浮点数 | float32,float64 |
| 字符串 | string |
| 布尔值 | bool |

如果超出以上类型，无法通过编译。

### 使用带有参数的装饰器

使用 `//go:decor function#{}` 的方式给装饰器传参。相对于无参调用，多了 `#{}` 这一部分，这部分我们称为参数域。

参数域以 `#` 标识开始，后跟 `{key:value, key1: value1}` 这样的键值对。其中键为装饰器的形参名，值为要传递的字符串、布尔值、或者数值。

例如，我们要调用上面定义的 `hit` 装饰器：

```go
package main
import "github.com/dengsgo/go-decorator/decor"

func hit(ctx *decor.Context, msg string, count int64, repeat bool, f float64, opt string) {
	// code...
}

//go:decor hit#{msg: "message from decor", repeat: true, count: 10, f:1}
func useArgsDecor()  {}
```

`decorator` 在编译时就会把 `{msg: "message from decor", repeat: true, count: 10, f:1}` 参数按形参名自动对应传递到装饰器中。

参数域中的参数顺序和装饰器的形参顺序无关，你可以按自己的习惯组织代码。

当参数域中没有对应的形参值时，比如上面的 `opt` ，`decorator` 会默认传递对应类型的零值。

### 装饰器约束和验证

`decorator` 允许在装饰器上使用注释 `//go:decor-lint linter: {}` 来添加装饰器约束。这个约束可以在编译时用来验证目标函数的调用是否合法。

目前内置了两种装饰器约束：

#### required

验证参数是必须要传的。例如：

```go
//go:decor-lint required: {msg, count, repeat, f}
func hit(ctx *decor.Context, msg string, count int64, repeat bool, f float64, opt string) {
	// code...
}
```

`msg, count, repeat, f` 这四个参数就要求目标函数在调用时是必须要传的，否则编译无法通过。

不仅如此，`required` 还支持验证枚举和范围。比如：

**枚举值限制**：  
`//go:decor-lint required: {msg: {"hello", "world", "yeah"}, count, repeat, f}`：要求 `msg` 的参数必须是 `"hello", "world", "yeah"` 这三个值中的一个。  

**范围限制**：  
`//go:decor-lint required: {msg: {gte: 8, lte: 24}, count, repeat, f}`：要求 `msg` 的字符串长度范围在 `[8,24]` 之间。
目前支持的范围指令有四个：

| 范围指令  | 说明 |
|-------|----|
| `gte` | `>=` |
| `gt`  | `>`  |
| `lte` | `<=` |
| `lt`  | `<`  |


#### nonzero

验证参数值不能为零值。例如：

```go
//go:decor-lint nonzero: {msg, count, f}
func hit(ctx *decor.Context, msg string, count int64, repeat bool, f float64, opt string) {
	// code...
}
```

`msg, count, f` 三个参数要求目标函数在调用时传值不能为零值。

> 可以在装饰器上多次添加 `//go:decor-lint` 规则约束，这意味着目标函数在调用装饰器时，必须全部满足这些约束才能正常编译。


## Context

`ctx *decor.Context` 是装饰器函数的入参，它是目标函数（即使用了这个装饰器的函数，也称为被装饰的函数）上下文。

在装饰器中可以使用这个上下文来修改目标函数的出入参，调整执行逻辑等操作。

### ctx.Kind

目标函数类型。
`decor.KFunc`: 函数, 目标函数是函数。  
`decor.KMethod`: 方法, 目标函数是个方法，此时的 `ctx.Receiver` 值为方法的接收者。

### ctx.TargetName

目标函数的函数名或方法名。

### ctx.Receiver

目标函数的接收者。如果 `ctx.Kind == decor.KFunc` （即函数类型），值为 nil。

### ctx.TargetIn

目标函数的入参列表。它是一个[]any slice, 其中每个元素的类型和目标函数的入参类型一致。 如果目标函数没有入参，列表为空。

`ctx.TargetDo()` 会使用这个 slice 来当作真实调用的入参，因此改变它的元素值可以修改目标函数的入参。只在 `ctx.TargetDo()` 调用前修改有效。

### ctx.TargetOut

目标函数的出参列表。它是一个[]any slice, 其中每个元素的类型和目标函数的出参类型一致。 如果目标函数没有出参，列表为空。

`ctx.TargetDo()` 会使用这个 slice 来接收真实调用的结果，因此改变它的元素值可以修改目标函数的出参。只在 `ctx.TargetDo()` 调用后修改有效。

### ctx.TargetDo()

执行目标函数。它是对目标函数的无参化包装，调用它才会真正的执行目标函数逻辑。

它从 `ctx.TargetIn` 获取到目标函数入参值，然后执行目标函数代码，再把得到的结果赋值到 `ctx.TargetOut`。

如果装饰器中没有执行 `ctx.TargetDo()` ，意味着目标函数真实的逻辑不会被执行，调用目标函数得到的结果是零值（在没有修改 ctx.TargetOut 的情况下）。  

### ctx.DoRef()  

获取匿名包装类被执行的次数。通常它代表着装饰器函数中执行 TargetDo()的次数。

> 在编写装饰器代码时要注意，一定要对 ctx.TargetIn、ctx.TargetOut 的元素值断言类型，任何类型错误的赋值都会产生 runtime panic。  
> 不要改变 ctx.TargetIn、ctx.TargetOut 值（赋值/追加/删除等），这会导致 ctx.TargetDo()  调用时产生严重错误 panic。

## 包引用

上面的 `datetime` [example/usages](example/usages) 例子中，我们的装饰器和目标函数都是在一个包中的，我们无需考虑包的问题。

当我们有很多包时，需要考虑包引用。

go 规范中，没有被当前文件里的代码使用到的包无法导入，这就导致了 `//go:decor` 这样的注释无法真正的导入包，因此需要我们使用匿名导入包的方式来导入对应的包。像这样 `import _ "path/to/your/package"`.

有下面两种情况需要使用匿名导入包：

一、包的函数使用了装饰器注释，但没有导入 `github.com/dengsgo/go-decorator/decor` 这个包，需要我们添加匿名导入包：

```go
import _ "github.com/dengsgo/go-decorator/decor"
```

二、如果包(A)引用了另外一个包(B)的装饰器, 而 B 没有被 A `import`,我们需要使用匿名导入包的方式导入它。

例如：

```go
package main

// other imports
import _ "github.com/dengsgo/go-decorator/example/packages/fun1"

//go:decor fun1.DecorHandlerFunc
func test() {
	//...
}
```

当然，如果包已经被文件里其他代码用到了，已经导入，那么就不需要再匿名导入了。

完整的例子可以查看 [example/usages](example/usages) .

## 条件和限制

以下几种情况需要注意：

- 使用装饰器的目标函数范围**仅限当前项目内**。依赖的其他库即使使用 `//go:decor`也**无法**被装饰。

例如，你的项目module名称是 `a/b/c` ，那么 `//go:decor` 只在 `a/b/c` 及其子包中生效（`a/b/c/d` 有效，`a/m/`无效）。

但是`//go:decor`可以使用任意包的装饰器，没有范围限制。

- **不能**在同一个目标函数上同时使用相同的装饰器重复装饰；  
- **不能**对装饰器函数应用装饰器；  
- 升级 `decorator` 后或者调整编译参数可能需要在 go 命令中追加 `-a` 参数**强制编译**一次，以覆盖旧的编译缓存。

## 开发与调试

`decorator` 作为 go 编译链中的一环，编译时被 go 编译器加载使用。它与 go 的编译链保持兼容，不会产生副作用。

开发流程中要改变的只是给用到的 go 命令增加 `-toolexec decorator` 参数，其他完全一致，感觉不到有变化。  

你也可以随时取消这个参数。放弃项目对 go 装饰器的使用。即使代码中保留了 `//go:decor ` 注释也不会有任何副作用（因为它对于标准工具链来说只是无意义的注释而已）。

调试同理。

例如，在 vscode 中，编辑 `launch.json`:

```json
{
    "version": "0.2.0",
    "configurations": [
        {
            "name": "Launch file",
            "type": "go",
            "request": "launch",
            "mode": "debug",
            "program": "${file}",
            "buildFlags": "-toolexec decorator"
        }
    ]
}
```

添加 `"buildFlags": "-toolexec decorator"` 这一行以启用 `decorator` 的装饰器编译。

然后正常断点调试即可。

> 调试体验会不断完善，如果发现问题请让我知道 [Issues](https://github.com/dengsgo/go-decorator/issues)。

## 性能

尽管 `decorator` 在编译时会对目标函数做额外的处理，但它仅仅只构建必要的上下文参数，没有额外开销，更没有反射。相对于原始go代码直接调用装饰器函数来讲，性能几乎是一致的。

// TODO 提供性能指标对比

## 更多

// TODO