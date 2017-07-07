package goref_test

import (
	"github.com/korfuri/goref"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestDotImports(t *testing.T) {
	const (
		pkgpath  = "github.com/korfuri/goref/testprograms/dotimports"
		filepath = "testprograms/dotimports/main.go"
	)

	pg := goref.NewPackageGraph()
	pg.LoadProgram(pkgpath, []string{filepath})
	assert.Contains(t, pg.Packages, pkgpath)
	assert.Contains(t, pg.Packages, pkgpath + "/lib")
	pkg := pg.Packages[pkgpath]
	assert.Empty(t, pkg.InRefs)
	assert.Len(t, pkg.OutRefs, 1)
	r := pkg.OutRefs[0]
	assert.True(t, goref.Call == r.RefType, r.String())
	assert.Equal(t, "Foo", r.ToIdent)
	assert.Equal(t, "Foo", r.FromIdent)
	assert.Equal(t, pg.Packages[pkgpath + "/lib"], r.ToPackage)
	assert.Equal(t, pkg, r.FromPackage)
}
