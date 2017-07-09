package goref_test

import (
	"testing"

	"github.com/korfuri/goref"
	"github.com/stretchr/testify/assert"
)

func TestMultipleMains(t *testing.T) {
	const (
		pkgbase  = "github.com/korfuri/goref/testprograms/multiple_mains/"
		filebase = "testprograms/multiple_mains/"
	)

	pg := goref.NewPackageGraph()
	pg.LoadProgram(pkgbase+"1/main", []string{filebase + "1/main.go"})
	assert.Contains(t, pg.Packages, pkgbase+"1/main")
	assert.Contains(t, pg.Packages, pkgbase+"common")
	common := pg.Packages[pkgbase+"common"]
	assert.Len(t, pg.Packages, 2)
	pg.LoadProgram(pkgbase+"2/main", []string{filebase + "2/main.go"})
	assert.Contains(t, pg.Packages, pkgbase+"1/main")
	assert.Contains(t, pg.Packages, pkgbase+"2/main")
	assert.Contains(t, pg.Packages, pkgbase+"common")
	assert.Len(t, pg.Packages, 3)
	assert.Equal(t, common, pg.Packages[pkgbase+"common"])
}
