package goref

import (
	"go/ast"
	"go/token"
	"go/types"
	"golang.org/x/tools/go/loader"

	"fmt"
	"log"
	"strings"
)

// PackageGraph represents a collection of Go packages and their
// mutual dependencies. All dependencies of a Package in the
// PackageGraph are also part of the PackageGraph.
type PackageGraph struct {
	// Map of package load-path to Package objects.
	Packages map[string]*Package

	// Map of file path to File objects.
	Files map[string]*File
}

// Package represents a Go Package, including its dependencies.
type Package struct {
	// Name of the package
	Name string

	// Dependencies is the map of load-paths imported within this
	// package to the corresponding Package objects.
	Dependencies map[string]*Package

	// Dependents is the map of packages' load-paths that load
	// this package through any load-path, mapped to their
	// corresponding Package objects.
	Dependents map[string]*Package

	// Files is a map of paths to File objects that make up this package.
	Files map[string]*File

	// OutRefs and InRefs are maps of references from a Position
	// (file:line:column). For OutRefs the Position is local to
	// the package and the Ref is to an identifier in another
	// package. For InRefs the Position is external to the package
	// and the Ref is to an identifier within this package.
	OutRefs map[Position]*Ref
	InRefs  map[Position]*Ref

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
}

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

// Valid RefTypes
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

// A Position is similar to token.Position in that it gives an
// absolute position within a file, but it may also denote the Pos +
// End concept of token.Pos.
//
// The End is optional. If NoPos is used as the End, Position only
// contains file:line:column.
//
// The Pos is not optional and must resolve to file:line:column.
type Position struct {
	File       string
	PosL, PosC int
	EndL, EndC int
}

func (p Position) String() string {
	if p.EndL >= 0 {
		return fmt.Sprintf("%s:[%d:%d]-[%d:%d]", p.File, p.PosL, p.PosC, p.EndL, p.EndC)
	} else {
		return fmt.Sprintf("%s:%d:%d", p.File, p.PosL, p.PosC)
	}
}

func NewPosition(fset *token.FileSet, pos, end token.Pos) Position {
	ppos := fset.Position(pos)
	if end == token.NoPos {
		return Position{
			File: ppos.Filename,
			PosL: ppos.Line,
			PosC: ppos.Column,
			EndL: -1,
			EndC: -1,
		}
	}
	pend := fset.Position(end)
	if ppos.Filename != pend.Filename {
		panic("Invalid pair of {pos,end} for NewPosition: pos and end come from different files!")
	}
	return Position{
		File: ppos.Filename,
		PosL: ppos.Line,
		PosC: ppos.Column,
		EndL: pend.Line,
		EndC: pend.Column,
	}
}

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

func (p *Package) String() string {
	return p.Name
}

func (r *Ref) String() string {
	return fmt.Sprintf("%s of `%s.%s` in %s at %s", r.RefType, r.ToPackage, r.Ident, r.FromPackage, r.Position)
}

func cleanImportSpec(spec *ast.ImportSpec) string {
	// FIXME we should make sure we understand what can cause Path
	// to be empty.
	if spec.Path != nil {
		s := spec.Path.Value
		s = strings.Trim(s, "\"")
		return s
	}
	return "<unknown>"
}

func newPackage(pi *loader.PackageInfo, fset *token.FileSet) *Package {
	return &Package{
		//PackageInfo:  pi,
		Name:         pi.Pkg.Name(),
		Dependencies: make(map[string]*Package),
		Dependents:   make(map[string]*Package),
		Files:        make(map[string]*File),
		OutRefs:      make(map[Position]*Ref),
		InRefs:       make(map[Position]*Ref),
		Interfaces:   make([]*types.Named, 0),
		Impls:        make([]*types.Named, 0),
		Fset:         fset,
	}
}

func newFile(fset *token.FileSet) *File {
	return &File{
		Fset:       fset,
		ImportPkgs: make(map[string]*Package),
	}
}

// loadPackage recursively loads a Go package into the Package
// Graph. If the package was already loaded, it returns early. It
// always returns the Package object for the loaded package.
func (pg *PackageGraph) loadPackage(prog *loader.Program, loadpath string, pi *loader.PackageInfo) *Package {
	if pi == nil {
		fmt.Printf("??? No PackageInfo for loadpath=%s\n", loadpath)
		return nil
	}
	if pkg, in := pg.Packages[loadpath]; in {
		return pkg
	}
	pkg := newPackage(pi, prog.Fset)
	pg.Packages[loadpath] = pkg

	// Iterate over all files in that package.
	for _, f := range pi.Files {
		// Get the File object to represent that file
		filepath := prog.Fset.File(f.Package).Name()
		ff := newFile(prog.Fset)

		// Add this File to the maps for the Package and for
		// the PackageGraph.
		pkg.Files[filepath] = ff
		pg.Files[filepath] = ff

		// Iterate over all imports in that file
		for _, imported := range f.Imports {
			// Find the import's load-path and load that
			// package into the graph.
			ipath := cleanImportSpec(imported)
			i := prog.Package(ipath)
			importedPkg := pg.loadPackage(prog, ipath, i)

			// Set up the edges on the package dependency graph
			importedPkg.Dependents[loadpath] = pkg
			pkg.Dependencies[ipath] = importedPkg
		}
	}

	// Iterate over all object uses in that package and filter for
	// non-local references only.
	for id, obj := range pi.Uses {
		// the object's Pkg will be nil for builtins
		if obj.Pkg() != nil {
			pkgLoadPath := obj.Pkg().Path()
			if pkgLoadPath != loadpath {
				foreignPkg := pg.Packages[pkgLoadPath]
				if foreignPkg != nil {
					ref := &Ref{
						RefType:     refTypeForId(prog, id),
						Position:    prog.Fset.Position(id.Pos()),
						Ident:       obj.Name(),
						ToPackage:   foreignPkg,
						FromPackage: pkg,
					}

					refpos := NewPosition(prog.Fset, id.Pos(), id.End())
					foreignPkg.InRefs[refpos] = ref
					pkg.OutRefs[refpos] = ref
				}
			}
		}
	}

	// Iterate over all types in that package and insert them as
	// needed into Structs and Interfaces.
	for _, name := range pi.Pkg.Scope().Names() {
		if obj, ok := pi.Pkg.Scope().Lookup(name).(*types.TypeName); ok {
			if named, ok := obj.Type().(*types.Named); ok {
				if types.IsInterface(named) {
					i := named.Obj().Type().Underlying().(*types.Interface)
					// We only care about interfaces that are exported, and that are not interface{}.
					if named.Obj().Exported() && i.NumMethods() > 0 {
						pkg.Interfaces = append(pkg.Interfaces, named)
					}
				} else {
					pkg.Impls = append(pkg.Impls, named)
				}
			}
		}
	}

	return pkg
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

// LoadProgram loads recursively packages used from a `main` package.
// It may be called multiple times to load multiple programs'
// package sets in the PackageGraph.
func (p *PackageGraph) LoadProgram(loadpath string, filename string) {
	conf := loader.Config{}
	conf.CreateFromFilenames(loadpath, filename)

	prog, err := conf.Load()
	if err != nil {
		log.Fatal(err)
	}

	for k, v := range prog.AllPackages {
		p.loadPackage(prog, k.Path(), v)
	}
}

func (p *PackageGraph) ComputeInterfaceImplementationMatrix() {
	for _, pa := range p.Packages {
		for _, iface := range pa.Interfaces {
			for _, pb := range p.Packages {
				for _, typ := range pb.Impls {
					if typ == iface {
						continue
					}
					if types.AssignableTo(typ, iface) {
						fset := pb.Fset
						pos := NewPosition(fset, typ.Obj().Pos(), token.NoPos)
						r := &Ref{
							RefType:     Implementation,
							Position:    fset.Position(typ.Obj().Pos()),
							Ident:       iface.Obj().Name(),
							ToPackage:   pa,
							FromPackage: pb,
						}
						pa.InRefs[pos] = r
						pb.OutRefs[pos] = r
					}
				}
			}
		}
	}
}

// NewPackageGraph returns a new, empty PackageGraph.
func NewPackageGraph() *PackageGraph {
	return &PackageGraph{
		Packages: make(map[string]*Package),
		Files:    make(map[string]*File),
	}
}
