package goref

import (
	"encoding/json"
	"fmt"

	pb "github.com/korfuri/goref/proto"
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
	return json.Marshal(r.ToProto())
}

// ToProto marshals a Ref as a pb.Ref
func (r Ref) ToProto() *pb.Ref {
	return &pb.Ref{
		Version: r.FromPackage.Version,
		From: &pb.Location{
			Position: r.FromPosition.ToProto(),
			Package:  r.FromPackage.Path,
			Ident:    r.FromIdent,
		},
		To: &pb.Location{
			Position: r.ToPosition.ToProto(),
			Package:  r.ToPackage.Path,
			Ident:    r.ToIdent,
		},
		Type: pb.Type(r.RefType),
	}
}
