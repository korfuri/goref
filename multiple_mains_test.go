package goref_test

import (
	"testing"

	"github.com/korfuri/goref"
	"github.com/stretchr/testify/assert"
)

func TestMultipleMainsFromMultipleCalls(t *testing.T) {
	const (
		pkgbase  = "github.com/korfuri/goref/testprograms/multiple_mains/"
		filebase = "testprograms/multiple_mains/"
	)

	pg := goref.NewPackageGraph(0)
	pg.LoadPrograms([]string{pkgbase+"1"}, false)
	assert.Contains(t, pg.Packages, pkgbase+"1")
	assert.Contains(t, pg.Packages, pkgbase+"common")
	common := pg.Packages[pkgbase+"common"]
	assert.Len(t, pg.Packages, 2)
	pg.LoadPrograms([]string{pkgbase+"2"}, false)
	assert.Contains(t, pg.Packages, pkgbase+"1")
	assert.Contains(t, pg.Packages, pkgbase+"2")
	assert.Contains(t, pg.Packages, pkgbase+"common")
	assert.Len(t, pg.Packages, 3)
	assert.Equal(t, common, pg.Packages[pkgbase+"common"])
}


func TestMultipleMainsFromSameCalls(t *testing.T) {
	const (
		pkgbase  = "github.com/korfuri/goref/testprograms/multiple_mains/"
		filebase = "testprograms/multiple_mains/"
	)

	pg := goref.NewPackageGraph(0)
	pg.LoadPrograms([]string{pkgbase+"1", pkgbase+"2"}, false)
	assert.Contains(t, pg.Packages, pkgbase+"1")
	assert.Contains(t, pg.Packages, pkgbase+"common")
	assert.Contains(t, pg.Packages, pkgbase+"1")
	assert.Contains(t, pg.Packages, pkgbase+"2")
	assert.Contains(t, pg.Packages, pkgbase+"common")
	assert.Len(t, pg.Packages, 3)
	common := pg.Packages[pkgbase+"common"]
	assert.Equal(t, common, pg.Packages[pkgbase+"common"])
}
