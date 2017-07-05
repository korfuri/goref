package goref

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCanImportEmptyPackage(t *testing.T) {
	const (
		emptypkgpath = "github.com/korfuri/goref/testdata/empty/main"
	)

	pg := NewPackageGraph()
	pg.LoadProgram(emptypkgpath, "testdata/empty/main.go")
	assert.Len(t, pg.Packages, 1)
	assert.Len(t, pg.Files, 1)
	assert.Empty(t, pg.Packages[emptypkgpath].InRefs)
	assert.Empty(t, pg.Packages[emptypkgpath].OutRefs)
}
