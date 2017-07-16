package testutils

import (
	"fmt"
	"strings"
	"testing"

	"github.com/korfuri/goref"
	"github.com/stretchr/testify/assert"
)

// ContainsRefP asserts that a collection of Ref contains at least one
// element that satisfies a given predicate.
func ContainsRefP(slc *[]*goref.Ref, pred func(*goref.Ref) bool) bool {
	return FindRefP(slc, pred) != nil
}

// FindRefP returns the first Ref in a colletion that satisfies the given predicate
func FindRefP(slc *[]*goref.Ref, pred func(*goref.Ref) bool) *goref.Ref {
	for _, r := range *slc {
		if pred(r) {
			return r
		}
	}
	return nil
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

// ToPositionFilenamePred returns a predicated that accepts a *Ref as
// argument and returns whether the Ref's ToPosition's Filename ends
// with `suffix`.
func ToPositionFilenamePred(suffix string) func(*goref.Ref) bool {
	return func(r *goref.Ref) bool {
		f := r.ToPosition.File
		return strings.HasSuffix(f, suffix)
	}
}

// AssertPresenceOfRef asserts that there exists a ref (should=true)
// or doesn't exist any ref (should=false) of type RefType from
// package fromPkg, from identifier fromID, to package toPkg, to
// identifier toID. This verifies that the Ref exists in both
// fromPkg.OutRefs and toPkg.InRefs.
func AssertPresenceOfRef(t *testing.T, toPkg *goref.Package, toID string, fromPkg *goref.Package, fromID string, reftype goref.RefType, should bool) {
	refref := &goref.Ref{
		FromPackage: fromPkg,
		ToPackage:   toPkg,
		FromIdent:   fromID,
		ToIdent:     toID,
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

// GetRef finds a ref in the InRefs of the provided package that
// matches the provided parameters. It fails the test if no such ref
// can be found or if the ref doesn't exist in the corresponding
// OutRefs.
func GetRef(t *testing.T, toPkg *goref.Package, toID string, fromPkg *goref.Package, fromID string, reftype goref.RefType) *goref.Ref {
	pred := EqualRefPred(&goref.Ref{
		FromPackage: fromPkg,
		FromIdent:   fromID,
		ToPackage:   toPkg,
		ToIdent:     toID,
		RefType:     reftype,
	})
	r := FindRefP(&toPkg.InRefs, pred)
	assert.NotNil(t, r)
	assert.Contains(t, fromPkg.OutRefs, r)
	return r
}
