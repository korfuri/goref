package goref_test

import (
	"github.com/korfuri/goref"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCanImportEmptyPackage(t *testing.T) {
	const (
		emptypkgpath = "github.com/korfuri/goref/testprograms/empty/main"
	)

	pg := goref.NewPackageGraph()
	pg.LoadProgram(emptypkgpath, []string{"testprograms/empty/main.go"})
	assert.Len(t, pg.Packages, 1)
	assert.Len(t, pg.Files, 1)
	assert.Empty(t, pg.Packages[emptypkgpath].InRefs)
	assert.Empty(t, pg.Packages[emptypkgpath].OutRefs)
}
