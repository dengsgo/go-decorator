package main

import _ "github.com/dengsgo/go-decorator/decor"

// this file is used to test underscores in parameters. decorator will rewrite the target code

//go:decor dumpDecorContext
func underscoresParamIn1Out1(_ string) (_ int) {
	return 1
}

//go:decor dumpDecorContext
func underscoresParamInHyperOut1(_ string, num int) (_ int) {
	return 1
}

//go:decor dumpDecorContext
func underscoresParamIn2Out2(_ string, num int) (_ int, f float32) {
	return 1, 1.0
}

//go:decor dumpDecorContext
func underscoresParamIn0Out2() (_ int, f float32) {
	return 1, 1.0
}

//go:decor dumpDecorContext
func underscoresParamIn2Out0(_ int, f float32) {
	//nothing
}
