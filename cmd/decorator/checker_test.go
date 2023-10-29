package main

import (
	"log"
	"testing"
)

func TestCheckDecorAndGetParam(t *testing.T) {
	param, err := checkDecorAndGetParam("github.com/dengsgo/go-decorator/decor", "find", nil)
	log.Println(param, err)
	param, err = checkDecorAndGetParam("github.com/dengsgo/go-decorator/cmd/decorator", "logging", nil)
	log.Println(param, err)
}

func TestCleanSpaceChar(t *testing.T) {
	cas := []struct {
		s,
		r string
	}{
		{"helloworld", "helloworld"},
		{"hello world", "helloworld"},
		{"hello 世界", "hello世界"},
		{" he l l owo      rld    ", "helloworld"},
		{"hello 世 界  这是测		试\t用     例 	  ", "hello世界这是测试用例"},
		{" 😀/(ㄒoㄒ)/~ ~   😊😄	😄\v😄  😄😄😄  ", "😀/(ㄒoㄒ)/~~😊😄😄😄😄😄😄"},
		{"if a > 1 {\necho ''\n}", "ifa>1{echo''}"},
	}
	for i, v := range cas {
		if cleanSpaceChar(v.s) != v.r {
			t.Fatal("cleanSpaceChar(v.s)!=r, pos", i, ": ", cleanSpaceChar(v.s), "!=", v.r)
		}
	}
}
