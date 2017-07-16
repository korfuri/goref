package goref_test

import (
	"testing"

	"github.com/korfuri/goref"
	"github.com/korfuri/goref/testutils"
	"github.com/stretchr/testify/assert"
)

func TestSimplePackage(t *testing.T) {
	const (
		pkgpath  = "github.com/korfuri/goref/testprograms/simple"
		filepath = "github.com/korfuri/goref/testprograms/simple/main.go"
	)

	pg := goref.NewPackageGraph(goref.ConstantVersion(0))
	pg.LoadPackages([]string{pkgpath}, false)
	assert.Contains(t, pg.Packages, pkgpath)
	assert.Contains(t, pg.Packages, "fmt")
	pkg := pg.Packages[pkgpath]
	assert.Empty(t, pkg.InRefs)
	assert.Contains(t, pkg.Files, filepath)

	testutils.AssertPresenceOfRef(t, pg.Packages["fmt"], "fmt", pkg, "fmt", goref.Import, true)

	r := testutils.GetRef(t, pg.Packages["fmt"], "Println", pkg, "Println", goref.Call)
	assert.NotNil(t, r)
	p := r.FromPosition
	assert.Equal(t, 6, p.PosL)
	assert.Equal(t, 6, p.PosC)
	assert.Equal(t, 6, p.EndL)
	assert.Equal(t, 13, p.EndC)
	assert.Equal(t, filepath, p.File)
}
