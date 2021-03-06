package goref_test

import (
	"testing"

	"github.com/korfuri/goref"
	"github.com/stretchr/testify/assert"
)

func TestVendoredPackage(t *testing.T) {
	const (
		pkgpath  = "github.com/korfuri/goref/testprograms/vendored"
		filepath = "testprograms/vendored/main.go"
	)

	pg := goref.NewPackageGraph(goref.ConstantVersion(0))
	pg.LoadPackages([]string{pkgpath}, false)
	assert.Contains(t, pg.Packages, pkgpath)
	assert.Contains(t, pg.Packages, "github.com/korfuri/goref/testprograms/vendored/vendor/github.com/korfuri/somedep")
}
