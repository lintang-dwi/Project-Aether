package parser_test

import (
	"context"
	"testing"

	"aether/internal/parser"
)

func TestGoParser(t *testing.T) {
	p := parser.NewGoParser()
	if !p.CanParse("main.go") {
		t.Fatal("expected GoParser to handle main.go")
	}

	code := []byte(`
package main

import "fmt"

type Server struct {
	Port int
}

func (s *Server) Start() error {
	fmt.Println("Starting")
	return nil
}

func main() {
	s := &Server{Port: 8080}
	s.Start()
}
`)

	res, err := p.Parse(context.Background(), "cmd/server/main.go", code)
	if err != nil {
		t.Fatalf("failed to parse Go code: %v", err)
	}

	if res.Language != "go" {
		t.Errorf("expected language 'go', got '%s'", res.Language)
	}
	if len(res.Imports) != 1 || res.Imports[0] != "fmt" {
		t.Errorf("expected import 'fmt', got %v", res.Imports)
	}

	var hasPkgMain, hasStructServer, hasMethodStart, hasFuncMain bool
	for _, sym := range res.Symbols {
		if sym.Name == "main" && sym.Kind == parser.KindPackage {
			hasPkgMain = true
		}
		if sym.Name == "main" && sym.Kind == parser.KindFunction {
			hasFuncMain = true
		}
		if sym.Name == "Server" && sym.Kind == parser.KindStruct {
			hasStructServer = true
		}
		if sym.Name == "Start" && sym.Kind == parser.KindMethod {
			hasMethodStart = true
		}
	}

	if !hasPkgMain {
		t.Error("expected package symbol 'main'")
	}
	if !hasFuncMain {
		t.Error("expected function symbol 'main'")
	}
	if !hasStructServer {
		t.Error("expected struct symbol 'Server'")
	}
	if !hasMethodStart {
		t.Error("expected method symbol 'Start'")
	}
}

func TestTSParser(t *testing.T) {
	p := parser.NewTSParser()
	if !p.CanParse("app.ts") {
		t.Fatal("expected TSParser to handle app.ts")
	}

	code := []byte(`
import { useState } from 'react';

export class App {
	render() {}
}

export function helper() {
	return 42;
}
`)

	res, err := p.Parse(context.Background(), "src/app.ts", code)
	if err != nil {
		t.Fatalf("failed to parse TS code: %v", err)
	}

	if len(res.Imports) != 1 || res.Imports[0] != "react" {
		t.Errorf("expected import 'react', got %v", res.Imports)
	}
}

func TestPyParser(t *testing.T) {
	p := parser.NewPyParser()
	if !p.CanParse("script.py") {
		t.Fatal("expected PyParser to handle script.py")
	}

	code := []byte(`
import os

class Processor:
    def process(self):
        pass

def run():
    p = Processor()
    p.process()
`)

	res, err := p.Parse(context.Background(), "scripts/run.py", code)
	if err != nil {
		t.Fatalf("failed to parse Python code: %v", err)
	}

	if len(res.Entities) < 2 {
		t.Errorf("expected at least 2 entities, got %d", len(res.Entities))
	}
}
