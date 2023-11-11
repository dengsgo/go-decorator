package main

import (
	_ "github.com/dengsgo/go-decorator/decor"
	_ "github.com/dengsgo/go-decorator/example/usages/externala"
)

//go:decor externala.OnlyPrintSelf
func useExternalaDecor() {
	// nothing to do
}
