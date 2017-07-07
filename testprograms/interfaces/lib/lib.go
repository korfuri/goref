package lib

type IfaceLibA interface {
	A()
}

type IfaceLibB interface {
	B()
}

type IfaceLibAB interface {
	A()
	B()
}

type LibA int
type libB int
type LibC int
type LibAB int

func (a LibA) A()   {}
func (b libB) B()   {}
func (c LibC) C()   {}
func (ab LibAB) A() {}
func (ab LibAB) B() {}
