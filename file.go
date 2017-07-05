package goref

import (
	"go/token"
)

// File represents a file that is part of a package loaded in a
// PackageGraph.
type File struct {
	// ImportPkgs is the map of load-paths imported within this
	// file to the corresponding Package objects.
	ImportPkgs map[string]*Package

	// Fset is a reference to the token.FileSet that loaded this
	// file.
	Fset *token.FileSet
}

func newFile(fset *token.FileSet) *File {
	return &File{
		Fset:       fset,
		ImportPkgs: make(map[string]*Package),
	}
}
