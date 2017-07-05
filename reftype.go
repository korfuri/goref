package goref

import (
	"go/ast"
	"golang.org/x/tools/go/loader"
)

// RefType is an enum of the various ways a package can reference an
// identifier in another package (e.g. as a call, an instantiation,
// etc.)
type RefType int

const (
	Instantiation = iota
	Call
	Implementation
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
	case Reference:
		return "Reference"
	default:
		panic("Unknown RefType used")
	}
}

func refTypeForId(prog *loader.Program, id *ast.Ident) RefType {
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
			}
		}
	}
	return Reference
}
