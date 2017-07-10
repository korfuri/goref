package goref_test

import (
	"testing"

	"github.com/korfuri/goref"
	"github.com/korfuri/goref/testutils"
	"github.com/stretchr/testify/assert"
)

func TestRelativeImports(t *testing.T) {
	const (
		pkgpath = "github.com/korfuri/goref/testprograms/relativeimports"
	)

	pg := goref.NewPackageGraph(0)
	pg.LoadPrograms([]string{pkgpath+"/lib"}, true)
	assert.Contains(t, pg.Packages, pkgpath+"/lib")
	assert.Contains(t, pg.Packages, pkgpath)
	pkg := pg.Packages[pkgpath]
	lib := pg.Packages[pkgpath+"/lib"]

	pred := testutils.EqualRefPred(&goref.Ref{
		FromPackage: lib,
		FromIdent:   "r",
		ToPackage:   pkg,
		ToIdent:     "lib",
		RefType:     goref.Import,
	})
	assert.True(t, testutils.ContainsRefP(&lib.OutRefs, pred))
}
