package goref_test

import (
	"testing"

	"github.com/korfuri/goref"
	"github.com/korfuri/goref/testutils"
	"github.com/stretchr/testify/assert"
)

func TestDotImports(t *testing.T) {
	const (
		pkgpath = "github.com/korfuri/goref/testprograms/dotimports"
	)

	pg := goref.NewPackageGraph(goref.ConstantVersion(0))
	pg.LoadPrograms([]string{pkgpath}, true)
	assert.Contains(t, pg.Packages, pkgpath)
	assert.Contains(t, pg.Packages, pkgpath+"/lib")
	pkg := pg.Packages[pkgpath]
	lib := pg.Packages[pkgpath+"/lib"]
	assert.Empty(t, pkg.InRefs)
	assert.Len(t, pkg.OutRefs, 4)
	testutils.AssertPresenceOfRef(t, lib, "Typ", pkg, "Typ", goref.Instantiation, true)
	testutils.AssertPresenceOfRef(t, lib, "Fun", pkg, "Fun", goref.Call, true)

	basepred := testutils.EqualRefPred(&goref.Ref{
		FromPackage: pkg,
		FromIdent:   ".",
		ToPackage:   lib,
		ToIdent:     "lib",
		RefType:     goref.Import,
	})
	{
		pospred := testutils.ToPositionFilenamePred("lib.go")
		r := testutils.FindRefP(&lib.InRefs, func(r *goref.Ref) bool {
			return basepred(r) && pospred(r)
		})
		assert.NotNil(t, r)
	}
	{
		pospred := testutils.ToPositionFilenamePred("lib2.go")
		r := testutils.FindRefP(&lib.InRefs, func(r *goref.Ref) bool {
			return basepred(r) && pospred(r)
		})
		assert.NotNil(t, r)
	}
}
