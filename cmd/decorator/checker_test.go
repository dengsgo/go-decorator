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

func TestIsLetters(t *testing.T) {
	cas := []struct {
		s string
		r bool
	}{
		{"thisisastring", true},
		{"this isastring", false},
		{"thisisastring ", false},
		{" thisisastring", false},
		{"这是string", true},
		{"这 是string", false},
		{"这是 string", false},
		{"这是string\t", false},
		{"这是\vstring", false},
		{"\n这是string", false},
		{"thisisa字符串", true},
		{"", false},
		{"\r", false},
		{"😀/(ㄒoㄒ)/~~😊😄😄😄😄😄😄", false},
		{" 😀/(ㄒoㄒ)/~ ~   😊😄	😄\v😄  😄😄😄  ", false},
	}
	for i, v := range cas {
		if isLetters(v.s) != v.r {
			t.Fatal("isLetters(v.s)!=r, pos", i, ": ", v.s, isLetters(v.s), "!=", v.r)
		}
	}
}
