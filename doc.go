package goref

// Goref is a package that analyzes a set of Go packages, starting
// from one or more `main` packages, and computes the reverse of the
// cross-package identifier usage graph. In other words, it indexes
// your Go code and tells you where an exported identifier is used. It
// can answer questions such as:
//
// * Where is this type instantiated?
// * Where is this function called?
// * Where are all the references to this identifier?
// * What types implement this interface?
// * What interfaces are implemented by this type?
