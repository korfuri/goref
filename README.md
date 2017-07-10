# Goref

[![GoDoc](https://godoc.org/github.com/korfuri/goref?status.svg)](http://godoc.org/github.com/korfuri/goref)
[![Build Status](https://travis-ci.org/korfuri/goref.svg?branch=master)](https://travis-ci.org/korfuri/goref)
[![Coverage Status](https://coveralls.io/repos/github/korfuri/goref/badge.svg?branch=master)](https://coveralls.io/github/korfuri/goref?branch=master)

Goref is a Go package that analyzes a set of Go programs, starting
from one or more `main` packages, and computes the inverse of the
cross-package identifier usage graph. In other words, it indexes your
Go code and tells you where an exported identifier is used. It can
answer questions such as:

* Where is this type instantiated?
* Where is this function called?
* Where are all the references to this identifier?
* What types implement this interface?
* What interfaces are implemented by this type?

Goref can be used as a library, see
its [godoc](http://godoc.org/github.com/korfuri/goref) for usage
information.

Goref can also be used to index code into ElasticSearch. This is
currently a Work-In-Progress. The binary for this is at
elasticsearch/main/main for lack of a good name until now. Usage is:

    ./main
	  --version 42    # version of the code being indexed
	  --include_tests # or --noinclude_tests to avoid indexing [XTests](https://godoc.org/golang.org/x/tools/go/loader#hdr-CONCEPTS_AND_TERMINOLOGY)
	  --elasticsearch_url http://localhost:9200
	  --elasticsearch_user user
	  --elasticsearch_password hunter2
	  github.com/korfuri/goref
	  your/awesome/pacakge

This always imports dependencies recursively.

## Code versioning

When code is indexed, the concept of "version" is critical. Since code
is a living thing, indexes of code must be versioned. Since the code
we index lives in many repositories, we can't use the repositories'
history as a versioning tool.

So versions have to be provided to the index externally. It's
recommended to keep a global, monotonically increasing counter for
this. Every time you `go get` one or more packages, that counter
should be incremented. The counter is an int64, so an elegant way to
do this is to use the `time.Time()` at which you last sync'd your entire
Go tree.

TODO(korfuri): versions should be per-package, not per-graph. There
should be a way to avoid duplicating all packages if only one package
was updated. Probably a callback passed by the user code that returns
the version for a given package, so that callback could look up the
latest mtime of all files in that package. Need to find a convenient
API.

If you'll be doing operations to a completely immutable tree of
packages (typically, your PackageGraph remains in memory and is never
serialized to disk, and you don't `go get` or `git pull` while loading
packages), you can just set the version to a fixed number and ignore
that.

Goref is (obviously) not safe to use if you concurrently update the
code while it's analyzing it.

## Vendoring and goref

Vendored packages are treated as separate packages in goref. There is
no support to deduplicate `github.com/foo/bar` and
`github.com/baz/qux/vendor/foo/bar`. This follows the `go`
tool's
[philosophy on that question](https://docs.google.com/document/d/1Bz5-UB7g2uPBdOx-rw5t9MxJwkfpx90cqG9AFL0JAYo/edit).

## Types of references

Currently, goref provides the following kinds of references, defined
in `reftype.go`:

* `Import` represents an import of a package by a file. Since a
  package doesn't exist in a single position, there will be a Ref
  from the importing file to each file in the imported package.
* `Call` represents a call of a function by another function.
* `Instantiation` are generated for composite literals.
* `Implementation` represent a reference from a type implementing an
  interface to that interface.
* `Extension` represent a reference from interface A to interface B if
  interface A is a superset of interface B.
* `Reference` is the default enum value, used if goref can't figure
  out what kind of reference is used but detects that a package
  depends on an identifier in another package.
