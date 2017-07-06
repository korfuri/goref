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
	FromPosition Position

	// Where this reference points to, i.e. where the definition
	// is
	ToPosition Position
	
	// What identifier points to this Ref
	FromIdent string

	// What identifier this Ref points to
	ToIdent string

	// What package contains what the identifier refers to
	ToPackage *Package

	// What package the ref is from, i.e. what foreign package was
	// this identifier used in.
	FromPackage *Package
}

func (r *Ref) String() string {
	return fmt.Sprintf("%s of to:`%s.%s` at %s by from:`%s.%s` at %s", r.RefType, r.ToPackage, r.ToIdent, r.ToPosition, r.FromPackage, r.FromIdent, r.FromPosition)
}
