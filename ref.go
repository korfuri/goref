package goref

import (
	"fmt"
	"go/token"
)

// A Ref is a reference to an identifier whose definition lives in
// another package.
type Ref struct {
	// Type of reference
	RefType

	// Where this reference points from, i.e. where the identifier
	// was used in another package.
	token.Position

	// What identifier this reference points to, i.e. what
	// identifier is referred to by another package. For most
	// references the name in the other package is identical, but
	// for Implementation references this is the name of the
	// interface.
	Ident string

	// What package the identifier is in
	ToPackage *Package

	// What package the ref is from, i.e. what foreign package was
	// this identifier used in.
	FromPackage *Package
}

func (r *Ref) String() string {
	return fmt.Sprintf("%s of `%s.%s` in %s at %s", r.RefType, r.ToPackage, r.Ident, r.FromPackage, r.Position)
}
