package goref

import (
	"fmt"
	"os"
	"time"

	"golang.org/x/tools/go/loader"
)

var (
	epoch = time.Unix(0, 0)
)

// ConstantVersion returns a versionF function that always replies
// with a constant version. Useful for experimenting, or for graphs
// who load from an immutable snapshot of the Go universe.
func ConstantVersion(v int64) (func(loader.Program, loader.PackageInfo) (int64, error)) {
	return func(prog loader.Program, pi loader.PackageInfo) (int64, error) {
		return v, nil
	}
}

// FileMTimeVersion is a versionF function that processes all files in
// the provided PackageInfo and returns the newest mtime's second as a
// time.Time-compatible int64.
func FileMTimeVersion(prog loader.Program, pi loader.PackageInfo) (int64, error) {
	newestMTime := time.Time{}
	for _, f := range pi.Files {
		filepath := prog.Fset.File(f.Package).Name()
		fi, err := os.Stat(filepath)
		if err != nil {
			return -1, err
		}
		fmt.Printf("%s mtime: %v\n", filepath, fi.ModTime())
		if fi.ModTime().After(newestMTime) {
			fmt.Printf("win!\n")
			newestMTime = fi.ModTime()
		}
	}
	if newestMTime == (time.Time{}) {
		return -1, fmt.Errorf("Unable to determine the version of package %s", pi.Pkg.Path())
	}
	// newestMTime - epoch gives us a duration which is an int64
	// of nanoseconds since the Unix epoch
	return int64(newestMTime.UTC().Sub(epoch)), nil
}

// FilterPass is a filterF function that always says yes.
func FilterPass(loadpath string, version int64) bool {
	return true
}
