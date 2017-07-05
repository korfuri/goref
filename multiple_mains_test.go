package goref_test

import (
	"github.com/korfuri/goref"
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestMultipleMains(t *testing.T) {
	const (
		pkgbase = "github.com/korfuri/goref/testdata/multiple_mains/"
		filebase = "testdata/multiple_mains/"
	)

	pg := goref.NewPackageGraph()
	pg.LoadProgram(pkgbase + "1/main", filebase + "1/main.go")
	assert.Contains(t, pg.Packages, pkgbase + "1/main")
	assert.Contains(t, pg.Packages, pkgbase + "common")
	common := pg.Packages[pkgbase + "common"]
	assert.Len(t, pg.Packages, 2)
	pg.LoadProgram(pkgbase + "2/main", filebase + "2/main.go")
	assert.Contains(t, pg.Packages, pkgbase + "1/main")
	assert.Contains(t, pg.Packages, pkgbase + "2/main")
	assert.Contains(t, pg.Packages, pkgbase + "common")
	assert.Len(t, pg.Packages, 3)
	assert.Equal(t, common, pg.Packages[pkgbase + "common"])
}
