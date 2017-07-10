package goref

import (
	"encoding/json"
	"fmt"
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
	return fmt.Sprintf("%s of to:`%s.%s` at %s by from:`%s.%s` at %s",
		r.RefType,
		r.ToPackage, r.ToIdent, r.ToPosition,
		r.FromPackage, r.FromIdent, r.FromPosition)
}

// MarshalJSON implements encoding/json.Marshaler interface
func (r Ref) MarshalJSON() ([]byte, error) {
	type location struct {
		Position Position `json:"position"`
		Pkg      string   `json:"package"`
		Ident    string   `json:"ident"`
	}
	type refForJSON struct {
		From    location `json:"from"`
		To      location `json:"to"`
		Typ     RefType  `json:"type"`
		Version int64    `json:"version"`
	}
	return json.Marshal(refForJSON{
		Version: r.FromPackage.Version,
		From: location{
			Position: r.FromPosition,
			Pkg:      r.FromPackage.Path,
			Ident:    r.FromIdent,
		},
		To: location{
			Position: r.ToPosition,
			Pkg:      r.ToPackage.Path,
			Ident:    r.ToIdent,
		},
		Typ: r.RefType,
	})
}
