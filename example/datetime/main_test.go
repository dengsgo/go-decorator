package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const out = `logging print target in [1692450000]
logging print target out [2023-08-19 21:00:00 +0800 CST]
datetime(1692450000)=2023-08-19 21:00:00 +0800 CST
mike
cook
`

func TestExampleDatetime(t *testing.T) {
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
		t.Fatalf("Run datetime fail %s", err)
	}
	if strings.ReplaceAll(bf.String(), "\r\n", "\n") !=
		strings.ReplaceAll(out, "\r\n", "\n") {
		t.Fatalf("datetime out fail, out %s", bf.String())
	}
}
