package parser

import (
	"testing"
)

func TestGoParser(t *testing.T) {
	src := []byte(`package main

import "fmt"
import "github.com/user/repo/helper"

func main() {
	fmt.Println(helper.Greet())
}

func Greet() string { return "hi" }

type User struct{}
`)
	p := NewGoParser()
	meta, err := p.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	if meta.Package != "main" {
		t.Errorf("package = %q, want main", meta.Package)
	}
	if len(meta.Imports) != 2 {
		t.Errorf("imports = %v, want 2", meta.Imports)
	}
	if len(meta.Exports) != 2 {
		t.Errorf("exports = %v, want 2", meta.Exports)
	}
}

func TestJavaScriptParser(t *testing.T) {
	src := []byte(`import { helper } from './helper';

export function greet() {
	return helper();
}

export class User {}
`)
	p := NewJavaScriptParser()
	meta, err := p.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(meta.Imports) != 1 {
		t.Errorf("imports = %v, want 1", meta.Imports)
	}
	if len(meta.Exports) != 2 {
		t.Errorf("exports = %v, want 2", meta.Exports)
	}
}

func TestPythonParser(t *testing.T) {
	src := []byte(`import os
from . import helper

def greet():
	return helper.message()

class User:
	pass
`)
	p := NewPythonParser()
	meta, err := p.Parse(src)
	if err != nil {
		t.Fatal(err)
	}
	if len(meta.Imports) != 2 {
		t.Errorf("imports = %v, want 2", meta.Imports)
	}
	if len(meta.Exports) != 2 {
		t.Errorf("exports = %v, want 2", meta.Exports)
	}
}
