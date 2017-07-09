package json

import (
	"encoding/json"

	"github.com/korfuri/goref"
)

func GraphAsJSON(pg goref.PackageGraph, outch chan<- []byte, errch chan<- error, done chan<- struct{}) {
	for _, p := range pg.Packages {
		for _, r := range p.InRefs {
			if j, err := json.Marshal(r); err == nil {
				outch <- j
			} else {
				errch <- err
			}
		}
	}
	done <- struct{}{}
}
