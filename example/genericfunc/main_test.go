package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const out = `logging print target in [1 10]
logging print target out [11]
plus(1, 10) =  11
`

func TestPlus(t *testing.T) {
	args := []string{
		"go", "run", "-toolexec", "decorator", "./main.go",
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = os.Environ()
	log.Println("args", args)
	bf := bytes.NewBuffer([]byte{})
	cmd.Stdout = bf
	cmd.Stderr = bf
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Run plus fail %s", err)
	}
	if strings.ReplaceAll(bf.String(), "\r\n", "\n") !=
		strings.ReplaceAll(out, "\r\n", "\n") {
		t.Fatalf("plus out fail, out %s", bf.String())
	}
}
