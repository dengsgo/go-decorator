package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const out = `useArgsDecor()=hit received: msg=message from decor, count=10, repeat=true, f=1.000000, opt=

`

func TestUseArgsDecor(t *testing.T) {
	args := []string{
		"go", "run", "-toolexec", "decorator", "./main.go",
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = os.Environ()
	//wd, _ := os.Getwd()
	//cmd.Dir = wd
	log.Println("args", args)
	bf := bytes.NewBuffer([]byte{})
	cmd.Stdout = bf
	cmd.Stderr = bf
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Run useArgsDecor fail %s", err)
	}
	if strings.ReplaceAll(bf.String(), "\r\n", "\n") !=
		strings.ReplaceAll(out, "\r\n", "\n") {
		t.Fatalf("useArgsDecor out fail, out %s", bf.String())
	}
}
