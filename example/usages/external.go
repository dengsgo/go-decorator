package main

// 这个包演示了如何使用一个外部包（非当前包）的装饰器。
// go 规范中，没有被当前文件里的代码使用到的包无法导入，这就导致了 `//go:decor` 这样的注释无法真正的导入包，
// 因此需要我们使用匿名导入包的方式来导入对应的包。像这样 `import _ "path/to/your/package"`.
// 如下面用到的 externala.OnlyPrintSelf 装饰器，需要这样导入： _ "github.com/dengsgo/go-decorator/example/usages/externala"
// 另外，因为当前文件使用了 //go:decor 注释语法，还需要导入： _ "github.com/dengsgo/go-decorator/decor"
//
// 如果包已经被用到，正常导入了，就无需再次匿名导入。
// 如当前文件已经导入 "github.com/dengsgo/go-decorator/example/usages/externalb"，使用这个包的装饰器直接用就行了。

import (
	_ "github.com/dengsgo/go-decorator/decor"
	_ "github.com/dengsgo/go-decorator/example/usages/externala"
	"github.com/dengsgo/go-decorator/example/usages/externalb"
)

//go:decor externala.OnlyPrintSelf
func useExternalaDecor() {
	// nothing to do
}

//go:decor externalb.DoubleIntegerValue
func plus(a, b int) int {
	return externalb.MathIntegerPlus(a, b)
}
