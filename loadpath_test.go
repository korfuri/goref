package goref

import (
	"github.com/stretchr/testify/assert"
	"go/ast"
	"testing"
)

func TestCleanImportSpec(t *testing.T) {
	assert.Equal(t, "foo/bar/baz", cleanImportSpec(&ast.ImportSpec{Path: &ast.BasicLit{Value: "foo/bar/baz"}}))
	assert.Equal(t, "foo/bar/baz", cleanImportSpec(&ast.ImportSpec{Path: &ast.BasicLit{Value: "\"foo/bar/baz\""}}))
}

func TestCandidatePaths(t *testing.T) {
	r := []string{
		"a/b/vendor/c/d",
		"a/vendor/c/d",
		"vendor/c/d",
		"c/d",
	}
	assert.Equal(t, r, candidatePaths("c/d", "a/b"))
}
