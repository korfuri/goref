package goref

import (
	"fmt"
	"go/build"
	"path/filepath"
	"strings"
)

// A Corpus represents a prefix from which Go packages may be loaded.
// Default corpora are $GOROOT/src and each of $GOPATH/src
type Corpus string

// NewCorpus creates a new Corpus.
func NewCorpus(basepath string) (Corpus, error) {
	if !filepath.IsAbs(basepath) {
		return Corpus(""), fmt.Errorf("Corpus %s has a relative basepath", basepath)
	}
	return Corpus(basepath), nil
}

// Contains returns whether the provided filepath exists under this
// Corpus.
func (c Corpus) Contains(fpath string) bool {
	rel, err := filepath.Rel(string(c), fpath)
	if err != nil {
		return false
	}
	return !strings.HasPrefix(rel, "../") && rel != ".."
}

// Rel returns the relative path of a file within a Corpus.
// If the string does not belong to the corpus, it returns fpath.
func (c Corpus) Rel(fpath string) string {
	if !c.Contains(fpath) {
		return fpath
	}
	rel, err := filepath.Rel(string(c), fpath)
	if err != nil {
		return fpath
	}
	return rel
}

// Returns the set of default corpora based on GOROOT and GOPATH.
func DefaultCorpora() []Corpus {
	srcdirs := build.Default.SrcDirs()
	corpora := make([]Corpus, len(srcdirs))
	for n, s := range srcdirs {
		corpora[n] = Corpus(s)
	}
	return corpora
}
