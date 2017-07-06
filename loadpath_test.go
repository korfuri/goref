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
