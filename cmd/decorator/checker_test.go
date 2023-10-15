package main

import (
	"log"
	"testing"
)

func TestCheckDecorAndGetParam(t *testing.T) {
	param, err := checkDecorAndGetParam("github.com/dengsgo/go-decorator/decor", "find", nil)
	log.Println(param, err)
}
