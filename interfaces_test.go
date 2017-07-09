package goref_test

import (
	"testing"

	"github.com/korfuri/goref"
	"github.com/korfuri/goref/testutils"
	"github.com/stretchr/testify/assert"
)

func TestInterfaceImplMatrix(t *testing.T) {
	const (
		pkgpath  = "github.com/korfuri/goref/testprograms/interfaces/main"
		filepath = "testprograms/interfaces/main.go"
	)

	pg := goref.NewPackageGraph()
	pg.LoadProgram(pkgpath, []string{filepath})
	assert.Contains(t, pg.Packages, pkgpath)
	pg.ComputeInterfaceImplementationMatrix()

	pkg := pg.Packages[pkgpath]
	lib := pg.Packages["github.com/korfuri/goref/testprograms/interfaces/lib"]
	assert.Len(t, pkg.InRefs, 17)

	// A/LibA implement IfaceA/IfaceLibA
	testutils.AssertPresenceOfRef(t, pkg, "IfaceA", pkg, "A", goref.Implementation, true)
	testutils.AssertPresenceOfRef(t, pkg, "IfaceA", lib, "LibA", goref.Implementation, true)
	testutils.AssertPresenceOfRef(t, lib, "IfaceLibA", pkg, "A", goref.Implementation, true)
	testutils.AssertPresenceOfRef(t, lib, "IfaceLibA", lib, "LibA", goref.Implementation, true)

	// B/libB implement IfaceB/IfaceLibB
	// libB is unexported but refs still exist
	testutils.AssertPresenceOfRef(t, pkg, "IfaceB", pkg, "B", goref.Implementation, true)
	testutils.AssertPresenceOfRef(t, pkg, "IfaceB", lib, "libB", goref.Implementation, true)
	testutils.AssertPresenceOfRef(t, lib, "IfaceLibB", pkg, "B", goref.Implementation, true)
	testutils.AssertPresenceOfRef(t, lib, "IfaceLibB", lib, "libB", goref.Implementation, true)

	// LibC does not implement ifaceC because ifaceC is not exported
	testutils.AssertPresenceOfRef(t, pkg, "ifaceC", lib, "LibC", goref.Implementation, false)

	// AB and LibAB implement all of: IfaceA, IfaceB, IfaceAB, IfaceLibA, IfaceLibB, IfaceLibAB
	testutils.AssertPresenceOfRef(t, pkg, "IfaceA", pkg, "AB", goref.Implementation, true)
	testutils.AssertPresenceOfRef(t, lib, "IfaceLibA", pkg, "AB", goref.Implementation, true)
	testutils.AssertPresenceOfRef(t, pkg, "IfaceB", pkg, "AB", goref.Implementation, true)
	testutils.AssertPresenceOfRef(t, lib, "IfaceLibB", pkg, "AB", goref.Implementation, true)
	testutils.AssertPresenceOfRef(t, pkg, "IfaceAB", pkg, "AB", goref.Implementation, true)
	testutils.AssertPresenceOfRef(t, lib, "IfaceLibAB", pkg, "AB", goref.Implementation, true)
	testutils.AssertPresenceOfRef(t, pkg, "IfaceA", lib, "LibAB", goref.Implementation, true)
	testutils.AssertPresenceOfRef(t, lib, "IfaceLibA", lib, "LibAB", goref.Implementation, true)
	testutils.AssertPresenceOfRef(t, pkg, "IfaceB", lib, "LibAB", goref.Implementation, true)
	testutils.AssertPresenceOfRef(t, lib, "IfaceLibB", lib, "LibAB", goref.Implementation, true)
	testutils.AssertPresenceOfRef(t, pkg, "IfaceAB", lib, "LibAB", goref.Implementation, true)
	testutils.AssertPresenceOfRef(t, lib, "IfaceLibAB", lib, "LibAB", goref.Implementation, true)

	// IfaceAB extends IfaceA, IfaceB, IfaceLibA and IfaceLibB.
	// It also extends IfaceLibAB by being equivalent to it.
	// It does not extend ifaceC as ifaceC is unexported.
	testutils.AssertPresenceOfRef(t, pkg, "IfaceA", pkg, "IfaceAB", goref.Extension, true)
	testutils.AssertPresenceOfRef(t, pkg, "IfaceB", pkg, "IfaceAB", goref.Extension, true)
	testutils.AssertPresenceOfRef(t, lib, "IfaceLibA", pkg, "IfaceAB", goref.Extension, true)
	testutils.AssertPresenceOfRef(t, lib, "IfaceLibB", pkg, "IfaceAB", goref.Extension, true)
	testutils.AssertPresenceOfRef(t, lib, "IfaceLibAB", pkg, "IfaceAB", goref.Extension, true)
	testutils.AssertPresenceOfRef(t, pkg, "ifaceC", pkg, "IfaceAB", goref.Extension, false)

	// Nothing implements not extends Empty, because we explicitly
	// filter out interface{}.
	testutils.AssertPresenceOfRef(t, pkg, "Empty", pkg, "A", goref.Implementation, false)
	testutils.AssertPresenceOfRef(t, pkg, "Empty", pkg, "B", goref.Implementation, false)
	testutils.AssertPresenceOfRef(t, pkg, "Empty", pkg, "AB", goref.Implementation, false)
	testutils.AssertPresenceOfRef(t, pkg, "Empty", lib, "LibA", goref.Implementation, false)
	testutils.AssertPresenceOfRef(t, pkg, "Empty", lib, "libB", goref.Implementation, false)
	testutils.AssertPresenceOfRef(t, pkg, "Empty", lib, "LibAB", goref.Implementation, false)
	testutils.AssertPresenceOfRef(t, pkg, "Empty", pkg, "IfaceA", goref.Extension, false)
	testutils.AssertPresenceOfRef(t, pkg, "Empty", pkg, "IfaceB", goref.Extension, false)
	testutils.AssertPresenceOfRef(t, pkg, "Empty", pkg, "IfaceAB", goref.Extension, false)
	testutils.AssertPresenceOfRef(t, pkg, "Empty", lib, "IfaceLibA", goref.Implementation, false)
	testutils.AssertPresenceOfRef(t, pkg, "Empty", lib, "IfaceLibB", goref.Implementation, false)
	testutils.AssertPresenceOfRef(t, pkg, "Empty", lib, "IfaceLibAB", goref.Implementation, false)
}
