package goref

import (
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestSimplePackage(t *testing.T) {
	const (
		pkgpath = "github.com/korfuri/goref/testdata/simple/main"
		filepath = "testdata/simple/main.go"
	)

	pg := NewPackageGraph()
	pg.LoadProgram(pkgpath, filepath)
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

		assert.True(t, Call == r.RefType)
		assert.Equal(t, r.Position.Line, 6)
		assert.Equal(t, r.Position.Column, 6)
		assert.Contains(t, r.Position.Filename, filepath)
		assert.Equal(t, "Println", r.Ident)
		assert.Equal(t, pg.Packages["fmt"], r.ToPackage)
		assert.Equal(t, pkg, r.FromPackage)
	}
}
