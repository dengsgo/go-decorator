package main

import (
	"log"
	"time"
)

func main() {
	// 这是一个使用包内装饰器的函数
	useScopeInnerDecor("hello, world", 100)

	// 这是一个使用其他包装饰器的函数
	useExternalaDecor()
}

func init() {
	log.SetFlags(0)
	time.Local = time.FixedZone("CST", 8*3600)
}
