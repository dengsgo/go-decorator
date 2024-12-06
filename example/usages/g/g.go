package g

// 该文件提供公共标识和方法

import (
	"bytes"
	"fmt"
	"io"
	"os"
)

var (
	InTest      = true
	TestBuffers = bytes.NewBuffer([]byte{})
)

func Printf(format string, a ...any) {
	_, _ = fmt.Fprintf(getIOWriter(), format, a...)
}

func PrintfLn(format string, a ...any) {
	Printf(format+"\n", a...)
}

func getIOWriter() io.Writer {
	if InTest {
		return TestBuffers
	}
	return os.Stdout
}

func ResetTestBuffers() {
	if InTest {
		TestBuffers.Reset()
	}
}
