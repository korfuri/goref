package goref_test

import (
	"github.com/korfuri/goref"
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestSimplePackage(t *testing.T) {
	const (
		pkgpath = "github.com/korfuri/goref/testprograms/simple/main"
		filepath = "testprograms/simple/main.go"
	)

	pg := goref.NewPackageGraph()
	pg.LoadProgram(pkgpath, []string{filepath})
	assert.Contains(t, pg.Packages, pkgpath)
	assert.Contains(t, pg.Packages, "fmt")
	pkg := pg.Packages[pkgpath]
	assert.Empty(t, pkg.InRefs)
	assert.Len(t, pkg.OutRefs, 1)
	for p, r := range pkg.OutRefs {
		assert.Equal(t, 6, p.PosL)
		assert.Equal(t, 6, p.PosC)
		assert.Equal(t, 6, p.EndL)
		assert.Equal(t, 13, p.EndC)
		assert.Contains(t, p.File, filepath)

		assert.True(t, goref.Call == r.RefType)
		assert.Equal(t, r.FromPosition.PosL, 6)
		assert.Equal(t, r.FromPosition.PosC, 6)
		assert.Contains(t, r.FromPosition.File, filepath)
		assert.Equal(t, "Println", r.ToIdent)
		assert.Equal(t, pg.Packages["fmt"], r.ToPackage)
		assert.Equal(t, pkg, r.FromPackage)
	}
}
