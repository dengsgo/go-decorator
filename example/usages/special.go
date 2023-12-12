package main

import (
	_ "github.com/dengsgo/go-decorator/decor"
	"time"
)

// 这个文件仅用来测试 decorator 工具。
// 函数输入和输出参数中有 nil 值，decorator 会在编译时做安全的类型转换。
// 有变长参数的函数，decorator 能够正常识别编译，生成相应语句。

//go:decor logging
func nilInAndOut(a any, b *time.Time, s string) (bool, error) {
	return false, nil
}

//go:decor logging
func ellipsisIn(i int, s ...string) []int {
	return []int{2024, 1, 1}
}
