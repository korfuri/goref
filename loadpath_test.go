package goref_test

import (
	"fmt"
	"go/ast"
	"testing"

	"github.com/korfuri/goref"
	"github.com/stretchr/testify/assert"
)

func TestCleanImportSpec(t *testing.T) {
	assert.Equal(t, "foo/bar/baz", goref.CleanImportSpec(&ast.ImportSpec{Path: &ast.BasicLit{Value: "foo/bar/baz"}}))
	assert.Equal(t, "foo/bar/baz", goref.CleanImportSpec(&ast.ImportSpec{Path: &ast.BasicLit{Value: "\"foo/bar/baz\""}}))
}

func ExampleCleanImportSpec() {
	fmt.Println(goref.CleanImportSpec(&ast.ImportSpec{Path: &ast.BasicLit{Value: "\"foo/bar/baz\""}}))
	// Output: foo/bar/baz
}

func TestCandidatePaths(t *testing.T) {
	r := []string{
		"a/b/vendor/c/d",
		"a/vendor/c/d",
		"vendor/c/d",
		"c/d",
	}
	assert.Equal(t, r, goref.CandidatePaths("c/d", "a/b"))
}

func ExampleCandidatePaths() {
	for _, p := range goref.CandidatePaths("lib/util", "program/bin") {
		fmt.Println(p)
	}
	// Output:
	// program/bin/vendor/lib/util
	// program/vendor/lib/util
	// vendor/lib/util
	// lib/util
}
