package goref

import (
	"go/ast"
	"go/token"
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
	// Refers to the PackageInfo object provided by `loader`
	// during program loading.
	//*loader.PackageInfo

	// Dependencies is the map of load-paths imported within this
	// package to the corresponding Package objects.
	Dependencies map[string]*Package

	// Dependents is the map of packages' load-paths that load
	// this package through any load-path, mapped to their
	// corresponding Package objects.
	Dependents map[string]*Package

	// Files is a map of paths to File objects that make up this package.
	Files map[string]*File

	// OutRefs is the map of token.Pos (position in the
	// token.FileSet) to Ref objects.
	OutRefs map[token.Pos]*Ref

	// InRefs is the map of token.Position (absolute position, as
	// in, file:line:col) into this file.
	InRefs map[token.Position]*Ref
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
	Reference
)

func (rt RefType) String() string {
	switch rt {
	case Instantiation:
		return "Instantiation"
	case Call:
		return "Call"
	case Reference:
		return "Reference"
	default:
		panic("Unknown RefType used")
	}
}

// A Ref is a reference to an identifier whose definition lives in
// another package.
type Ref struct {
	// Type of reference
	RefType

	// Where this reference points to
	// FIXME: inrefs should point to a Pos, or to an Identifier,
	// not to a Position (so we can do inverse matching of "what
	// points to this identifier?")
	token.Position

	// What identifier this reference points to
	Ident string
}

func (r *Ref) String() string {
	return fmt.Sprintf("%s of `%s` at %s", r.RefType, r.Ident, r.Position)
}

func cleanImportSpec(spec *ast.ImportSpec) string {
	s := spec.Path.Value
	s = strings.Trim(s, "\"")
	return s
}

func newPackage(pi *loader.PackageInfo) *Package {
	return &Package{
		//PackageInfo:  pi,
		Dependencies: make(map[string]*Package),
		Dependents:   make(map[string]*Package),
		Files:        make(map[string]*File),
		OutRefs:      make(map[token.Pos]*Ref),
		InRefs:       make(map[token.Position]*Ref),
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
	if pkg, in := pg.Packages[loadpath]; in {
		return pkg
	}
	pkg := newPackage(pi)
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
				ref := &Ref{
					RefType:  refTypeForId(prog, id),
					Position: prog.Fset.Position(id.Pos()),
					Ident:    obj.Name(),
				}

				// // Walk the file's AST to find OutRefs and index those.
				// ast.Inspect(f, func(n ast.Node) bool {
				// 	switch n := n.(type) {
				// 	case *ast.CallExpr:
				// 		f := n.Fun
				// 		fmt.Printf("Call: %s, Fun has type: %s ", n, reflect.TypeOf(f))
				// 		switch f := f.(type) {
				// 		case *ast.SelectorExpr:
				// 			fmt.Printf("and value  %s [.] %s\n", f.Sel, f.X)
				// 		}
				// 	}
				// 	return true
				// })

				// FIXME: this pulls the foreign package's
				// object's Position from the current
				// package's Fileset. This will not work if
				// both packages were loaded as part of 2
				// different programs (which may happen as
				// packages only get loaded once).
				foreignPkg.InRefs[prog.Fset.Position(obj.Pos())] = ref
				pkg.OutRefs[id.Pos()] = ref
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
			switch n.Fun.(type) {
			case *ast.SelectorExpr:
				return Call
			}
		case *ast.CompositeLit:
			switch n.Type.(type) {
			case *ast.SelectorExpr:
				return Instantiation
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

// NewPackageGraph returns a new, empty PackageGraph.
func NewPackageGraph() *PackageGraph {
	return &PackageGraph{
		Packages: make(map[string]*Package),
		Files:    make(map[string]*File),
	}
}
