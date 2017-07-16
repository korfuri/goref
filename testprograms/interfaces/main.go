// Package main is a test program that demonstrates
// interface-implementation relations in goref.
package main

import (
	"github.com/korfuri/goref/testprograms/interfaces/lib"
)

// use is a function to avoid "unused variable" errors
func use(interface{}) {}

func main() {
	use(lib.LibA(0))
	acceptAB(AB(0))
}

// IfaceA is a simple interface
type IfaceA interface {
	A()
}

// IfaceB is a simple interface
type IfaceB interface {
	B()
}

// IfaceAB is a simple interface that extends IfaceA and IfaceB.
type IfaceAB interface {
	A()
	B()
}

// ifaceC is an unexported interface.
type ifaceC interface {
	C()
}

// Empty is an empty interface.
type Empty interface{}

// A A implements IfaceA.
// ^ golint doesn't differentiate between the leading article and the
// type name.
type A int

// B implements IfaceB.
type B int

// AB implements IfaceAB, IfaceA and IfaceB.
type AB int

// A implements IfaceA.
func (a A) A() {}

// B implements IfaceB.
func (b B) B() {}

// A implements IfaceA.
func (ab AB) A() {}

// B implements IfaceB.
func (ab AB) B() {}

// acceptAB tests for acceptance of IfaceAB.
func acceptAB(IfaceAB) {}
