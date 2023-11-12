package main

// 这个文件演示了泛型函数使用装饰器的用法。
// 它和普通函数的用法没有任何区别。

import _ "github.com/dengsgo/go-decorator/decor"

//go:decor logging
func Sum[T int8 | int16 | int | int32 | int64 | float32 | float64](a ...T) T {
	var r T
	for _, v := range a {
		r += v
	}
	return r
}
