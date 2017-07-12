package goref_test

import (
	"testing"

	"github.com/korfuri/goref"
	"github.com/stretchr/testify/assert"
)

func TestCanImportEmptyPackage(t *testing.T) {
	const (
		emptypkgpath = "github.com/korfuri/goref/testprograms/empty"
	)

	pg := goref.NewPackageGraph(goref.ConstantVersion(0))
	pg.LoadPrograms([]string{emptypkgpath}, false)
	assert.Len(t, pg.Packages, 1)
	assert.Empty(t, pg.Packages[emptypkgpath].InRefs)
	assert.Empty(t, pg.Packages[emptypkgpath].OutRefs)
}
