package goref

import (
	"encoding/json"
	"go/ast"

	"golang.org/x/tools/go/loader"
)

// RefType is an enum of the various ways a package can reference an
// identifier in another package (e.g. as a call, an instantiation,
// etc.)
type RefType int

// These are the possible types of edges in a graph.
const (
	// Instantiation of a type in another package.
	Instantiation = iota

	// Call of a function in another package.
	Call

	// Implementation of an interface by a type.
	Implementation

	// Extension of an interface by another interface.
	Extension

	// Import is the import of a package by another. `fromIdent`
	// may differ from the name of the target package in the case
	// of named imports. For dot-imports, `fromIdent` is ".".
	Import

	// Reference is the default, used when we can't determine the
	// type of reference.
	Reference
)

func (rt RefType) String() string {
	switch rt {
	case Instantiation:
		return "Instantiation"
	case Call:
		return "Call"
	case Implementation:
		return "Implementation"
	case Extension:
		return "Extension"
	case Import:
		return "Import"
	case Reference:
		return "Reference"
	}
	panic("Unknown RefType used")
}

// refTypeForIdent walks the AST from a given Ident and deducts what
// type of Reference it is performing.
func refTypeForIdent(prog *loader.Program, id *ast.Ident) RefType {
	_, path, _ := prog.PathEnclosingInterval(id.Pos(), id.End())
	// Walk the file's AST to find OutRefs and index those.
	for _, n := range path {
		//fmt.Printf("Id: %s, Node: %v\n", id.Name, n)
		switch n := n.(type) {
		case *ast.CallExpr:
			// If this identifier appears in a CallExpr,
			// we make sure it appears as the function as
			// part of a SelectorExpr (because it will be
			// of the form package.Function(args).
			switch f := n.Fun.(type) {
			case *ast.SelectorExpr:
				if f.Sel == id {
					return Call
				}
			case *ast.Ident:
				// It's possible that Foo() refers to
				// an imported identifier, in the case
				// of dot imports.
				if f == id {
					return Call
				}
			}
		case *ast.CompositeLit:
			// A CompositeLit is an expression of the form
			// Type{...}. We check that the Type is a
			// SelectorExpr because we are looking for
			// package.Type{...}.
			switch t := n.Type.(type) {
			case *ast.SelectorExpr:
				if t.Sel == id {
					return Instantiation
				}
			case *ast.Ident:
				// It's possible that Foo{} refers to
				// an imported identifier, in the case
				// of dot imports.
				if t == id {
					return Instantiation
				}
			}
		}
	}
	return Reference
}

// MarshalJSON implements encoding/json.Marshaler interface
func (rt RefType) MarshalJSON() ([]byte, error) {
	return json.Marshal(rt.String())
}
