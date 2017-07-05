package goref_test

import (
	"github.com/korfuri/goref"
	"github.com/stretchr/testify/assert"
	"testing"
)


func TestVendoredPackage(t *testing.T) {
	const (
		pkgpath = "github.com/korfuri/goref/testprograms/vendored/main"
		filepath = "testprograms/vendored/main.go"
	)

	pg := goref.NewPackageGraph()
	pg.LoadProgram(pkgpath, []string{filepath})
	assert.Contains(t, pg.Packages, pkgpath)
	assert.Contains(t, pg.Packages, "github.com/korfuri/goref/testprograms/vendored/vendor/github.com/korfuri/somedep")
}
