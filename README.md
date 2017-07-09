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
