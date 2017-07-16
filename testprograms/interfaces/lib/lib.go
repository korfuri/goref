// Package lib contains interfaces and types that implement them. It
// is very similar to package main, whose interfaces its types also
// implements.
package lib

// IfaceLibA is a simple test interface.
type IfaceLibA interface {
	A()
}

// IfaceLibB is a simple test interface.
type IfaceLibB interface {
	B()
}

// IfaceLibAB is a simple test interface that extends IfaceLibA and
// IfaceLibB.
type IfaceLibAB interface {
	A()
	B()
}

// LibA implements IfaceLibA
type LibA int

// libB implements IfaceLibB, and is unexported
type libB int

// LibC doesn't implement anything from this package but implements
// the unexported main.ifaceC.
type LibC int

// LibAB implements IfaceLibAB
type LibAB int

// A implements IfaceLibA
func (a LibA) A() {}

// B implements IfaceLibB
func (b libB) B() {}

// C is just a member functions of LibC
func (c LibC) C() {}

// A implements IfaceLibA
func (ab LibAB) A() {}

// B implements IfaceLibB
func (ab LibAB) B() {}
