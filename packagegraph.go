package goref

import (
	"go/ast"
	"go/types"
	"path"
	"strings"

	log "github.com/sirupsen/logrus"
	"golang.org/x/tools/go/loader"
)

// PackageGraph represents a collection of Go packages and their
// mutual dependencies. All dependencies of a Package in the
// PackageGraph are also part of the PackageGraph.
type PackageGraph struct {
	// Map of package load-path to Package objects.
	Packages map[string]*Package

	// Map of file path to File objects.

	// Slice of corpora that files may be loaded from
	corpora []Corpus

	// versionF is a function that returns the version of the
	// provided Go package as an int64.
	//
	// A simple way to implement it, for a graph that will not
	// outlive its packages' versions, is to provide a function
	// that always returns 0. A more correct way to implement it
	// is to look at the package's files and check their mtime.
	// Both are provided in packagegraph_utils.go.
	//
	// If versionF returns an error for a given package, that
	// package is not loaded.
	versionF func(loader.Program, loader.PackageInfo) (int64, error)

	// filterF is a function that determines whether a package
	// version should be loaded into the graph.
	filterF func(loadpath string, version int64) bool
}

// CleanImportSpec takes an ast.ImportSpec and cleans the Path
// component by trimming the quotes (") that surround it.
func CleanImportSpec(spec *ast.ImportSpec) string {
	// FIXME we should make sure we understand what can cause Path
	// to be empty.
	if spec.Path != nil {
		s := spec.Path.Value
		s = strings.Trim(s, "\"")
		return s
	}
	return "<unknown>"
}

// CandidatePaths returns a slice enumerating all the possible import
// paths for a package. This means inserting the possible "vendor"
// directory location from the load path of the importing package.
//
// If package a/b imports c/d, the following paths are candidates:
// a/b/vendor/c/d
// a/vendor/c/d
// vendor/c/d
// c/d
//
// Order matters, as the most-specific vendored package is selected.
// Note that multi-level vendoring works, as PackageGraph will
// consider the full import path, including path/to/vendor/, as the
// package path when building the graph. In this sense we follow the
// go tool's convention to not try to detect when two packages loaded
// through different paths are the same package.
func CandidatePaths(loadpath, parent string) []string {
	const vendor = "vendor"
	paths := []string{}
	for parent != "." && parent != "" {
		paths = append(paths, path.Join(parent, vendor, loadpath))
		parent = path.Dir(parent)
	}
	// Some dependencies may be vendored under
	// $GOROOT/src/vendor. This is the case e.g. for
	// `golang_org/x/net/lex/httplex` which is imported by
	// `net/http`. This is the correct path: it's just vendored
	// that way in the standard library. See
	// https://github.com/golang/go/issues/16333 for background on
	// that.
	paths = append(paths, path.Join(vendor, loadpath))
	paths = append(paths, loadpath)
	return paths
}

// specialPackage returns hardcoded packages for packages wihtout a Go
// implementation. Currently this is only "unsafe".
func specialPackage(loadpath string) *Package {
	if loadpath == "unsafe" {
		return &Package{
			Name:   "unsafe",
			Path:   "unsafe",
			InRefs: make([]*Ref, 0),

			// There exists a src/unsafe/unsafe.go but it
			// is not loaded by go/loader, and so
			// referencing it would cause more issues than
			// not (as no Fileset would contain its
			// contents). So the Fset is empty.
			Fset: nil,

			// We assume Unsafe has no interfaces, types or outrefs.
			OutRefs:    make([]*Ref, 0),
			Interfaces: make([]*types.Named, 0),
			Impls:      make([]*types.Named, 0),

			// We assume Unsafe only has one version and
			// that version is immutable. If Unsafe's API
			// changes that's OK, since anything can
			// reference it at this version already.
			Version: 1,
		}
	}
	return nil
}

// loadPackage recursively loads a Go package into the Package
// Graph. If the package was already loaded, it returns early. It
// always returns the Package object for the loaded package.
func (pg *PackageGraph) loadPackage(prog *loader.Program, loadpath string, pi *loader.PackageInfo) *Package {
	var pkg *Package = nil
	// First find whether we already know about this package,
	// either because it's already in the graph (then return it)
	// or because it has a hardcoded definition.
	if existingPkg, in := pg.Packages[loadpath]; in {
		return existingPkg
	}
	if specialPkg := specialPackage(loadpath); specialPkg != nil {
		pkg = specialPkg
	}

	// If we didn't find a hardcoded package, rely on versionF to
	// tell us what version this package should have. Otherwise
	// use the hardcoded Package's Version.
	var version int64
	if pkg == nil {
		var err error
		version, err = pg.versionF(*prog, *pi)
		if err != nil {
			return nil
		}
	} else {
		version = pkg.Version
	}

	// Find what corpus this package was loaded from.
	var corpus Corpus
	if len(pi.Files) > 0 {
		fpath := prog.Fset.File(pi.Files[0].Package).Name()
		for _, c := range pg.corpora {
			if c.Contains(fpath) {
				corpus = c
				break
			}
		}
	}

	// If pkg was a hardcoded package, return it now. Otherwise,
	// create it.
	if pkg != nil {
		return pkg
	}
	pkg = newPackage(pi, prog.Fset, version, corpus)
	pg.Packages[loadpath] = pkg

	// Apply filterF to stop loading any package that doesn't pass
	// the filter. Note that if a package was already present in
	// the graph, filterF is not called.
	if !pg.filterF(loadpath, version) {
		return pkg
	}

	// Iterate over all files in that package.
	for _, f := range pi.Files {
		// Add that file to the package's file list
		pkg.Files = append(pkg.Files, corpus.Rel(prog.Fset.File(f.Package).Name()))

		// Iterate over all imports in that file
		for _, imported := range f.Imports {
			// Find the import's load-path and load that
			// package into the graph.
			ipath := CleanImportSpec(imported)
			candidatePaths := CandidatePaths(ipath, loadpath)
			var i *loader.PackageInfo
			for _, c := range candidatePaths {
				i = prog.Package(c)
				if i != nil {
					ipath = c
					break
				}
			}
			if i == nil {
				log.Warnf("Tried to load package `%s` imported by package `%s` but it wasn't found anywhere in the load path. The candidate load paths were: %s\n", ipath, loadpath, candidatePaths)
				continue
			}
			importedPkg := pg.loadPackage(prog, ipath, i)
			if importedPkg == nil {
				// This happens if versionF fails to
				// determine the pacakge's version.
				// Often this would mean that the
				// package is a C package.
				continue
			}

			// Set up the edges on the package dependency graph
			var importAs string
			// If the import is unqualified
			if imported.Name == nil {
				importAs = i.Pkg.Name()
			} else {
				importAs = imported.Name.String()
			}

			// Create a Ref to each file in the imported
			// package. This is useful in two ways:
			// reverse lookups can happen from any file in
			// a package, and users are free to decide
			// what file(s) they want to look up after
			// finding `Import` OutRefs in a package.
			for _, f := range i.Files {
				if f.Name != nil {
					r := &Ref{
						RefType:      Import,
						FromPosition: NewPosition(corpus, prog.Fset, imported.Pos(), imported.End()),
						ToPosition:   NewPosition(corpus, prog.Fset, f.Name.Pos(), f.Name.End()),
						FromIdent:    importAs,
						ToIdent:      i.Pkg.Name(),
						FromPackage:  pkg,
						ToPackage:    importedPkg,
					}
					pkg.OutRefs = append(pkg.OutRefs, r)
					importedPkg.InRefs = append(importedPkg.InRefs, r)
				}
			}
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
						RefType:      refTypeForIdent(prog, id),
						ToIdent:      obj.Name(),
						ToPackage:    foreignPkg,
						ToPosition:   NewPosition(corpus, prog.Fset, obj.Pos(), NoPos),
						FromIdent:    id.Name,
						FromPackage:  pkg,
						FromPosition: NewPosition(corpus, prog.Fset, id.Pos(), id.End()),
					}

					foreignPkg.InRefs = append(foreignPkg.InRefs, ref)
					pkg.OutRefs = append(pkg.OutRefs, ref)
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

// LoadPrograms loads the specified packages and their transitive
// dependencies, as well as XTests (as defined by go/loader) if
// includeTests is true.  It may be called multiple times to load
// multiple package sets in the PackageGraph.
func (pg *PackageGraph) LoadPrograms(packages []string, includeTests bool) error {
	conf := loader.Config{}
	if _, err := conf.FromArgs(packages, includeTests); err != nil {
		return err
	}

	prog, err := conf.Load()
	if err != nil {
		return err
	}

	for k, v := range prog.AllPackages {
		pg.loadPackage(prog, k.Path(), v)
	}

	return nil
}

// ComputeInterfaceImplementationMatrix processes all loaded types and
// adds cross-package and intra-package Refs for Implementation and
// Extension edges of the graph.
func (pg *PackageGraph) ComputeInterfaceImplementationMatrix() {
	for _, pa := range pg.Packages {
		for _, iface := range pa.Interfaces {
			for _, pb := range pg.Packages {
				for _, typ := range pb.Impls {
					if typ == iface {
						continue
					}
					if types.AssignableTo(typ, iface) {
						r := &Ref{
							RefType:      Implementation,
							ToIdent:      iface.Obj().Name(),
							ToPackage:    pa,
							ToPosition:   NewPosition(pa.Corpus, pa.Fset, iface.Obj().Pos(), NoPos),
							FromIdent:    typ.Obj().Name(),
							FromPackage:  pb,
							FromPosition: NewPosition(pb.Corpus, pb.Fset, typ.Obj().Pos(), NoPos),
						}
						pa.InRefs = append(pa.InRefs, r)
						pb.OutRefs = append(pb.OutRefs, r)
					}
				}
				for _, ifaceb := range pb.Interfaces {
					if ifaceb == iface {
						continue
					}
					if types.AssignableTo(ifaceb, iface) {
						r := &Ref{
							RefType:    Extension,
							ToIdent:    iface.Obj().Name(),
							ToPackage:  pa,
							ToPosition: NewPosition(pa.Corpus, pa.Fset, ifaceb.Obj().Pos(), NoPos),

							FromIdent:    ifaceb.Obj().Name(),
							FromPackage:  pb,
							FromPosition: NewPosition(pb.Corpus, pb.Fset, ifaceb.Obj().Pos(), NoPos),
						}
						pa.InRefs = append(pa.InRefs, r)
						pb.OutRefs = append(pb.OutRefs, r)
					}
				}
			}
		}
	}
}

// NewPackageGraph returns a new, empty PackageGraph.
func NewPackageGraph(versionF func(loader.Program, loader.PackageInfo) (int64, error)) *PackageGraph {
	p := &PackageGraph{
		Packages: make(map[string]*Package),
		versionF: versionF,
		filterF:  FilterPass,
		corpora:  DefaultCorpora(),
	}
	return p
}

// SetFilterF sets the filterF for this PackageGraph
func (pg *PackageGraph) SetFilterF(f func(string, int64) bool) {
	pg.filterF = f
}
