package main

import (
	"strings"
	"testing"
)

func TestUseArgsDecor(t *testing.T) {
	s := `hit received: msg=message from decor, count=10, repeat=true, f=1.000000, opt=`
	r := useArgsDecor()
	if strings.TrimSpace(r) != s {
		t.Fatalf("TestUseArgsDecor fail")
	}
}

func TestUseHitUseRequiredLint(t *testing.T) {
	s := `hit received: msg=你好, count=10, repeat=false, f=1.000000, opt=`
	r := useHitUseRequiredLint()
	if strings.TrimSpace(r) != s {
		t.Fatalf("TestUseArgsDecor fail")
	}
}

func TestUseHitUseNonzeroLint(t *testing.T) {
	s := `hit received: msg=你好, count=150, repeat=false, f=1.000000, opt=`
	r := useHitUseNonzeroLint()
	if strings.TrimSpace(r) != s {
		t.Fatalf("TestUseArgsDecor fail")
	}
}

func TestUseHitBothUseLint(t *testing.T) {
	s := `hit received: msg=message from decor, useHitBothUseLint, count=10, repeat=true, f=1.000000, opt=`
	r := useHitBothUseLint()
	if strings.TrimSpace(r) != s {
		t.Fatalf("TestUseArgsDecor fail")
	}
}

func TestUseHitUseMultilineLintDecor(t *testing.T) {
	s := `hit received: msg=hello, count=150, repeat=true, f=1.000000, opt=`
	r := useHitUseMultilineLintDecor()
	if strings.TrimSpace(r) != s {
		t.Fatalf("TestUseArgsDecor fail")
	}
}
