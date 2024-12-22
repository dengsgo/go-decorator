# Q&A  

- [Q\&A](#qa)
  - [我在什么场景下需要用到 go-decorator ?](#我在什么场景下需要用到-go-decorator-)
  - [只有这个场景才能用到 go-decorator 吗？](#只有这个场景才能用到-go-decorator-吗)
  - [go-decorator 是 Go 官方库吗？](#go-decorator-是-go-官方库吗)
  - [什么时候发布正式版本？](#什么时候发布正式版本)
  - [go-decorator 如何保证稳定性？](#go-decorator-如何保证稳定性)
  - [迭代节奏是怎样的？](#迭代节奏是怎样的)
  - [还会给 go-decorator 加入更多特性吗？](#还会给-go-decorator-加入更多特性吗)
  - [go-decorator 对用户来说似乎是个黑魔法，可以详细说下实现原理吗？](#go-decorator-对用户来说似乎是个黑魔法可以详细说下实现原理吗)
  - [可以邀请在公司或组织内部分享布道吗？](#可以邀请在公司或组织内部分享布道吗)

## 我在什么场景下需要用到 go-decorator ?

如果目前 go 标准工具链已经能够满足你的需求，那么就无需考虑引入 `go-decorator`。不要因为手里有锤子，就到处找钉子，不要创造需求。

通常情况，只有你的项目到达一定规模后，在此基础上实现一些比较通用的逻辑时，发现通过现有工具链、三方库无法解决，或者要动大量现有逻辑，存在大量重复编码的情况，这时可以考虑使用 `go-decorator`。它最大的优点就是允许开发者以非侵入的方式控制函数，记录、改变函数的行为。

比如有如下场景：生产环境有一 Web 程序出现了一些性能问题，但是出现的概率较低，呈现随机偶发性。经过编码CR，怀疑是某些请求的参数导致一些函数响应延迟高导致的。现在要排查这个问题，找出是哪些函数、哪些参数导致的，要如何在不影响线上稳定性的情况下，去 debug 呢？

大致方案：
- 日志方案。把用户请求链路全部或者采样的落下来，采集离线来分析。  
- 改代码，埋点监控。在认为有问题的代码里手动插入一些埋点，然后离线来分析。  
- 代码里已经有各种 log, 直接线上某个节点开 debug 级别日志，然后离线分析。  
- 其他方案...

这些方案或多或少存在一些现实问题：
- **风险**。为了解决 bug 要做代码修改，尤其是大量修改现有代码逻辑，到处埋点，这不仅严重破坏了函数单一职能原则，加剧代码理解的成本，还对其他比如单测、高敏感逻辑等造成不可预测的风险。  
- **成本**。代码修改成本，日志储存成本，投入的人力成本等等。  
- **复用性**。事实上，这种问题一般都是 case by case, 很难说抽象出通用的模型来处理。一般的做法也就是继续完善日志系统、加更多埋点上报等，尽量采集更多的信息应对以后可能出现的问题。  

但是现在 `go-decorator` 为你提供了新的思路，使用它来处理此类问题就很方便。处理流程为：

**一、编写函数 LogFuncInfo，记录出入参和耗时** 

伪代码：

```go
func LogFuncInfo(elapsed int, ctx *decor.Context) {
    in = getInArgsMsg(ctx.TargetIn)
    out = getOutArgsMsg(ctx.TargetOut)
    msg = format("elapsed:%s, funcName:%s, In: %s, Out: %s", elapsed, ctx.TargetName, in, out)
    Log.send(msg)
}
```

**二、编写装饰器函数 RecordSlowFunc, 记录慢函数**

伪代码： 

```go 
func RecordSlowFunc(ctx *decor.Context) {
    startTime = now()
    ctx.TargetDo() // 执行原函数(目标函数)
    elapsed = now() - startTime // 耗时
    const slowTime = getSlowTimeConfig() // ms
    if (elapsed >= slowTime) {
        LogFuncInfo(elapsed, ctx) // 记录慢日志
    }
}
```

**三、给任意函数加上注释，使用装饰器 RecordSlowFunc**  

伪代码：

```go
//go:decor RecordSlowFunc
func maybeSlow(g *Context, state *State, req *Request, other *Other) {
    // code ...
}
```

代码修改结束！  

我们不需要改动任何现有代码的逻辑，要做的只是给它们加上一行注释 `//go:decor RecordSlowFunc`！

通过 `go-decorator` 编译后，所有加上此注释的函数都会被自动注入逻辑：当执行时间超过配置的慢函数时长就会被捕获，把此次函数的出入参、函数名、具体耗时上报给日志系统。

因此可以看到，使用 `go-decorator` 之后。不仅能够精确的采集到想要的数据，还几乎让我们没有心智负担的来修改代码。更重要的是，我们非常自然的抽象出了可复用模型。后面再发生此类问题进行排查，此经验直接复用，对于个人、团队来说，都是宝贵的财富。  

风险、成本和可复用性三者有了较优解。  

## 只有这个场景才能用到 go-decorator 吗？

当然不是。上面的例子只是给大家一个初步的认知，可以使用非侵入的方式来处理问题，而不用大动干戈到处改代码逻辑。  

理论上，这些场景都非常适合考虑使用 `go-decorator`：
- 存在大量重复代码的编写  
- 存在大量模板化的代码  
- 解决逻辑杂糅问题  
- 保持函数职能单一原则，保持简洁不被污染  

所以，一些独立的模块，像 日志、缓存、ORM、配置、埋点等，都可以将他们包装成装饰器函数从而给需要的函数使用，这样既天然代码解耦合，还能隐藏实现细节，专注于逻辑本身。  

总之，要权衡利弊，在适合自己的场景考虑。  

## go-decorator 是 Go 官方库吗？

不是。目前为止它仍然只有一个人（就是我）参与维护的产品。

但好在 `go-decorator` 是完全开源的，并且基于友好的 `MIT License`，任何组织和个人都可获取代码并修改完善。在此也呼吁感兴趣的同学参与其中，贡献力量。

## 什么时候发布正式版本？

很遗憾，无法提供具体日期。尽管在设计和实现上，`go-decorator` 一直致力于保持稳定的生产可用性，但由于目前涉及的测试用户样本和反馈不足，我对此仍持谨慎态度，无法在短时间内发布正式版本。一旦条件成熟，将及时调整节奏和发布。

## go-decorator 如何保证稳定性？

- 更多的测试用例，尽可能覆盖更多的边界场景；  
- 基于 Issue 的用户反馈；  
- 外部贡献；  

在未收到更多测试样本和反馈之前，`go-decorator` 将保持 v0 测试版本的迭代。 

## 迭代节奏是怎样的？

当前所有我认为需要的核心特性均已实现。在发布正式版本之前，不再考虑加新特性（除非 Go 上游有重大变更），因此迭代均以修复 BUG 为主，比较依赖 Issue 的反馈。

## 还会给 go-decorator 加入更多特性吗？

Less is more. 

保持克制，尽可能完成最核心的功能和稳定的体验，是当前的重点工作。如果你有新特性的想法或者构思可以通过 Issue 提交反馈。

## go-decorator 对用户来说似乎是个黑魔法，可以详细说下实现原理吗？

非常愿意分享。但受限于两点：
- 需要用户在此之前对 Go 的编译、内部机制有了解才能比较容易理解 `go-decorator` 的实现原理  
- 受众用户太少  

因此暂无博文发布，请关注后续动态。可以先通过代码阅读了解实现原理。

## 可以邀请在公司或组织内部分享布道吗？

可以。通过邮件方式联系我。

E: WkdWdVozTm5ieU5uYldGcGJDTmpiMjA9  
Title: go-decorator 内部分享邀请
