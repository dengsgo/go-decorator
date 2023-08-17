package main

import (
	"bytes"
	"log"
	"os"
	"os/exec"
	"strings"
	"testing"
)

const singleOut = `MAIN call myFunc1()
call decorHandlerFunc in 0 []
call myFunc1
call decorHandlerFunc out []
MAIN call myFunc2UseMultipleDecor()
call decorHandlerFunc in 0 []
call yetDecorHandlerFunc in 0 []
call myFunc1
call yetDecorHandlerFunc out []
call decorHandlerFunc out []
MAIN call myFunc2UseMultipleDecor()
call decorHandlerFunc in 0 [100 this is test case []]
call decorHandlerFunc out`

func TestExampleSingle(t *testing.T) {
	args := []string{
		"go", "run", "-toolexec", "decorator", "./main.go",
	}
	cmd := exec.Command(args[0], args[1:]...)
	cmd.Env = os.Environ()
	log.Println(args)
	//cmd.Dir = "./"
	bf := bytes.NewBuffer([]byte{})
	cmd.Stdout = bf
	cmd.Stderr = bf
	err := cmd.Run()
	if err != nil {
		t.Fatalf("Run single fail %s", err)
	}
	if !strings.HasPrefix(
		strings.ReplaceAll(bf.String(), "\r\n", "\n"),
		strings.ReplaceAll(singleOut, "\r\n", "\n")) {
		t.Fatalf("single out fail, out %s", bf.String())
	}
}
