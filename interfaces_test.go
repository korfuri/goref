package goref_test

import (
	"fmt"
	"testing"

	"github.com/korfuri/goref"
	"github.com/stretchr/testify/assert"
)

// Asserts that a collection of Ref contains at least one element that
// satisfies a given predicate.
func ContainsRefP(slc *[]*goref.Ref, pred func(*goref.Ref) bool) bool {
	for _, r := range *slc {
		if pred(r) {
			return true
		}
	}
	return false
}

func stringifyRefslice(slc []*goref.Ref) string {
	s := ""
	for _, r := range slc {
		s = s + fmt.Sprintf("[%s]\n", r)
	}
	return s
}

func assertImplementsOrExtends(t *testing.T, ifacePkg *goref.Package, iface string, implPkg *goref.Package, impl string, reftype goref.RefType, should bool) {
	refref := &goref.Ref{
		FromPackage: implPkg,
		ToPackage:   ifacePkg,
		FromIdent:   impl,
		ToIdent:     iface,
		RefType:     reftype,
	}
	pred := func(r *goref.Ref) bool {
		return (r.FromPackage == refref.FromPackage &&
			r.FromIdent == refref.FromIdent &&
			r.ToIdent == refref.ToIdent &&
			r.ToPackage == refref.ToPackage &&
			r.RefType == refref.RefType)
	}
	if should {
		assert.True(t, ContainsRefP(&ifacePkg.InRefs, pred), "ifacePkg(%s).InRefs does not contain a Ref matching the expected Ref: [%s].\nInRefs was: \n%s", ifacePkg, refref.String(), stringifyRefslice(implPkg.OutRefs))
		assert.True(t, ContainsRefP(&implPkg.OutRefs, pred), "implPkg(%s).OutRefs does not contain a Ref matching the expected Ref: [%s].\nOutRefs was: \n%s", implPkg, refref.String(), stringifyRefslice(implPkg.OutRefs))
	} else {
		assert.False(t, ContainsRefP(&ifacePkg.InRefs, pred), "ifacePkg(%s).InRefs contains a Ref matching the forbidden Ref: [%s].\nInRefs was: \n%s", ifacePkg, refref.String(), stringifyRefslice(implPkg.OutRefs))
		assert.False(t, ContainsRefP(&implPkg.OutRefs, pred), "implPkg(%s).OutRefs contains a Ref matching the forbidden Ref: [%s].\nOutRefs was: \n%s", implPkg, refref.String(), stringifyRefslice(implPkg.OutRefs))
	}
}

func TestInterfaceImplMatrix(t *testing.T) {
	const (
		pkgpath  = "github.com/korfuri/goref/testprograms/interfaces/main"
		filepath = "testprograms/interfaces/main.go"
	)

	pg := goref.NewPackageGraph()
	pg.LoadProgram(pkgpath, []string{filepath})
	assert.Contains(t, pg.Packages, pkgpath)
	pg.ComputeInterfaceImplementationMatrix()

	pkg := pg.Packages[pkgpath]
	lib := pg.Packages["github.com/korfuri/goref/testprograms/interfaces/lib"]
	assert.Len(t, pkg.InRefs, 17)

	// A/LibA implement IfaceA/IfaceLibA
	assertImplementsOrExtends(t, pkg, "IfaceA", pkg, "A", goref.Implementation, true)
	assertImplementsOrExtends(t, pkg, "IfaceA", lib, "LibA", goref.Implementation, true)
	assertImplementsOrExtends(t, lib, "IfaceLibA", pkg, "A", goref.Implementation, true)
	assertImplementsOrExtends(t, lib, "IfaceLibA", lib, "LibA", goref.Implementation, true)

	assertImplementsOrExtends(t, pkg, "IfaceB", lib, "libB", goref.Implementation, true)

	// assertContainsRefP(t, &pkg.InRefs, func(r *goref.Ref) bool {
	// 	return (r.FromPackage == lib &&
	// 		r.ToIdent == "IfaceA" &&
	// 		r.RefType == goref.Implementation &&
	// 		r.ToPackage == pkg)
	// })
	// assertContainsRefP(t, &pkg.OutRefs, func(r *goref.Ref) bool {
	// 	return (r.FromPackage == pkg &&
	// 		r.ToIdent == "IfaceLibA" &&
	// 		r.RefType == goref.Implementation &&
	// 		r.ToPackage == lib)
	// })
	// assertContainsRefP(t, &lib.InRefs, func(r *goref.Ref) bool {
	// 	return (r.FromPackage == pkg &&
	// 		r.ToIdent == "IfaceLibA" &&
	// 		r.RefType == goref.Implementation &&
	// 		r.ToPackage == lib)
	// })
	// assertContainsRefP(t, &pkg.InRefs, func(r *goref.Ref) bool {
	// 	return (r.FromPackage == lib &&
	// 		r.ToIdent == "AB" &&
	// 		r.RefType == goref.Implementation &&
	// 		r.ToPackage == pkg)
	// })
	// assertContainsRefP(t, &pkg.InRefs, func(r *goref.Ref) bool {
	// 	return (r.FromPackage == lib &&
	// 		r.ToIdent == "IfaceA" &&
	// 		r.RefType == goref.Extension &&
	// 		r.ToPackage == pkg)
	// })
}
