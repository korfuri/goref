package main

import (
	"github.com/korfuri/goref/testprograms/interfaces/lib"
)

func use(interface{}) {}

func main() {
	use(lib.LibA(0))
	acceptAB(AB(0))
}

type IfaceA interface {
	A()
}

type IfaceB interface {
	B()
}

type IfaceAB interface {
	A()
	B()
}

type ifaceC interface {
	C()
}

type Empty interface{}

type A int
type B int
type AB int

func (a A) A()   {}
func (b B) B()   {}
func (ab AB) A() {}
func (ab AB) B() {}

func acceptAB(IfaceAB) {}
