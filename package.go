package goref

import (
	"encoding/json"
	"fmt"
	"go/token"
	"go/types"

	"golang.org/x/tools/go/loader"
)

// Package represents a Go Package, including its dependencies.
type Package struct {
	// Name of the package
	Name string

	// OutRefs and InRefs are slices of references. For OutRefs
	// the Ref is to an identifier in another package. For InRefs
	// the Ref is to an identifier within this package.  Most
	// RefTypes are not indexed if the ToPackage and the
	// FromPackage are the same, but some do such as
	// Implementation. This means that a ref can exist in both
	// OutRefs and InRefs of the same package.
	OutRefs []*Ref
	InRefs  []*Ref

	// Interfaces is the list of interface types in this package.
	//
	// This is used to compute the interface-implementation matrix.
	//
	// Only named interfaces matter, because an unnamed interface
	// can't be exported.
	//
	// Interfaces equivalent to interface{} are excluded.
	Interfaces []*types.Named

	// Impls is the list of non-interface types in this package.
	//
	// This is used to compute the interface-implementation matrix.
	//
	// Only named types matter, because an unnamed type can't have
	// methods.
	Impls []*types.Named

	// Fset is a reference to the token.FileSet that loaded this
	// package.
	Fset *token.FileSet

	// Version is the version of the package that was loaded.
	Version int64

	// Path is the package's load path
	Path string
}

// String implements the Stringer interface
func (p *Package) String() string {
	return p.Name
}

// DocumentID returns a consistent id for this package at this
// version. This can be used to index the package e.g. in
// ElasticSearch. The ID contains the document version and path.
func (p Package) DocumentID() string {
	// "v1" is a prefix to recognize this DocumentID format, in
	// case the format changes in the future.
	return fmt.Sprintf("v1@%d@%s", p.Version, p.Path)
}

// MarshalJSON implements encoding/json.Marshaler interface
func (p Package) MarshalJSON() ([]byte, error) {
	return json.Marshal(struct {
		Path    string `json:"loadpath"`
		Version int64  `json:"version"`
	}{
		Path:    p.Path,
		Version: p.Version,
	})
}

func newPackage(pi *loader.PackageInfo, fset *token.FileSet, version int64) *Package {
	return &Package{
		//PackageInfo:  pi,
		Name:       pi.Pkg.Name(),
		OutRefs:    make([]*Ref, 0),
		InRefs:     make([]*Ref, 0),
		Interfaces: make([]*types.Named, 0),
		Impls:      make([]*types.Named, 0),
		Fset:       fset,
		Version:    version,
		Path:       pi.Pkg.Path(),
	}
}
