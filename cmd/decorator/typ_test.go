package main

import "testing"

func TestNewMapV(t *testing.T) {
	r := newMapV[string, string]()
	if r.items == nil {
		t.Fatal("r.items == nil should be false")
	}
	if r.put("a", "b") == false {
		t.Fatal("r.put(\"a\", \"b\") == false should be true")
	}
	if r.get("a") != "b" {
		t.Fatal("r.get(\"a\") != \"b\" should be true")
	}
	if r.put("a", "b") == true {
		t.Fatal("r.put(\"a\", \"b\") == false should be false")
	}
}
