package validation_test

import (
	"context"
	"testing"

	"aether/internal/parser"
	"aether/internal/validation"
)

func TestValidationEngine_ValidGoCode(t *testing.T) {
	pipe := parser.NewPipeline()
	pipe.RegisterParser(parser.NewGoParser())

	valEngine := validation.NewEngine(pipe)

	validGo := []byte(`
package main

import "fmt"

func main() {
	fmt.Println("Hello")
}
`)

	res := valEngine.ValidateFile(context.Background(), "main.go", validGo)
	if !res.IsValid {
		t.Errorf("expected valid Go code, got diagnostics: %v", res.Diagnostics)
	}
}
