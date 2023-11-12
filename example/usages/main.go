package main

import (
	"github.com/dengsgo/go-decorator/example/usages/g"
	"log"
	"time"
)

func main() {
	// 这是一个使用包内装饰器的函数
	useScopeInnerDecor("hello, world", 100)

	// 这是一个使用其他包装饰器的函数
	useExternalaDecor()
	g.PrintfLn("plus(2, 3) = %+v", plus(2, 3))

	// 文档 Guide.md 中演示使用装饰器的代码
	{
		t := 1692450000
		s := datetime(t)
		g.Printf("datetime(%d)=%s\n", t, s)
	}
}

func init() {
	log.SetFlags(0)
	time.Local = time.FixedZone("CST", 8*3600)
}
