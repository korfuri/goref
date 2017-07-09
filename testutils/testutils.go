package testutils

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

// stringifyRefslice returns a string representation of a []*Ref.
func stringifyRefslice(slc []*goref.Ref) string {
	s := ""
	for _, r := range slc {
		s = s + fmt.Sprintf("[%s]\n", r)
	}
	return s
}

// EqualRefPred returns a predicate that accepts a *Ref as argument
// and returns whether it is equal to compRef.
func EqualRefPred(compRef *goref.Ref) func(*goref.Ref) bool {
	return func(r *goref.Ref) bool {
		return (r.FromPackage == compRef.FromPackage &&
			r.FromIdent == compRef.FromIdent &&
			r.ToIdent == compRef.ToIdent &&
			r.ToPackage == compRef.ToPackage &&
			r.RefType == compRef.RefType)
	}
}

// Asserts that there exists a ref (should=true) or doesn't exist any
// ref (should=false) of type RefType from package fromPkg, from
// identifier fromId, to package toPkg, to identifier toId. This
// verifies that the Ref exists in both fromPkg.OutRefs and
// toPkg.InRefs.
func AssertPresenceOfRef(t *testing.T, toPkg *goref.Package, toId string, fromPkg *goref.Package, fromId string, reftype goref.RefType, should bool) {
	refref := &goref.Ref{
		FromPackage: fromPkg,
		ToPackage:   toPkg,
		FromIdent:   fromId,
		ToIdent:     toId,
		RefType:     reftype,
	}
	pred := EqualRefPred(refref)
	if should {
		assert.True(t, ContainsRefP(&toPkg.InRefs, pred), "ifacePkg(%s).InRefs does not contain a Ref matching the expected Ref: [%s].\nInRefs was: \n%s", toPkg, refref.String(), stringifyRefslice(toPkg.InRefs))
		assert.True(t, ContainsRefP(&fromPkg.OutRefs, pred), "implPkg(%s).OutRefs does not contain a Ref matching the expected Ref: [%s].\nOutRefs was: \n%s", fromPkg, refref.String(), stringifyRefslice(fromPkg.OutRefs))
	} else {
		assert.False(t, ContainsRefP(&toPkg.InRefs, pred), "ifacePkg(%s).InRefs contains a Ref matching the forbidden Ref: [%s].\nInRefs was: \n%s", toPkg, refref.String(), stringifyRefslice(toPkg.InRefs))
		assert.False(t, ContainsRefP(&fromPkg.OutRefs, pred), "implPkg(%s).OutRefs contains a Ref matching the forbidden Ref: [%s].\nOutRefs was: \n%s", fromPkg, refref.String(), stringifyRefslice(fromPkg.OutRefs))
	}
}
