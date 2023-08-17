package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const packagesOut = `MAIN call myFunc1()
call fun.decorHandlerFunc in 0 []
call myFunc1
call fun1.decorHandlerFunc in 0 []
call Ts yeah
call fun1.decorHandlerFunc out []
call fun.decorHandlerFunc out []
call fun1.decorHandlerFunc in 0 []
call Ts yeah
call fun1.decorHandlerFunc out []
`

func TestExamplePackages(t *testing.T) {
	args := []string{
		"go", "run", "-toolexec", "decorator", "main.go",
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
		t.Fatalf("Run packages fail %s", err)
	}
	if strings.ReplaceAll(bf.String(), "\r\n", "\n") !=
		strings.ReplaceAll(packagesOut, "\r\n", "\n") {
		t.Fatalf("packages out fail, out %s", bf.String())
	}
}
