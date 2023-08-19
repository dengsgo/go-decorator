package main

import "testing"

func TestDecorX(t *testing.T) {
	sucCases := []string{
		"log.Println",
		"fmt.Printf",
		"x.s",
		"log.",
		"decor.Context",
	}
	failCases := []string{
		"log",
		".Printf",
		"",
		"x.a.",
		"aaaa##c",
	}
	for i, s := range sucCases {
		if decorX(s) == "" {
			t.Fatalf("decorX('%s') should pass, case sucCases i: %d\n", s, i)
		}
	}
	for i, s := range failCases {
		if decorX(s) != "" {
			t.Fatalf("decorX('%s') should fail, case failCases i: %d\n", s, i)
		}
	}
}
