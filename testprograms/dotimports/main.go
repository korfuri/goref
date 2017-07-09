package main

import (
	. "github.com/korfuri/goref/testprograms/dotimports/lib"
)

func use(interface{}) {}

func main() {
	use(Typ{})
	Fun()
}
