package goref_test

import (
	"testing"

	"github.com/korfuri/goref"
	"github.com/stretchr/testify/assert"
)

func TestContains(t *testing.T) {
	assert.True(t, goref.Corpus("/a/b/c").Contains("/a/b/c/d"))
	assert.True(t, goref.Corpus("/a/b/c/d").Contains("/a/b/c/d"))
	assert.True(t, goref.Corpus("/").Contains("/a/b/c/d"))
	assert.False(t, goref.Corpus("/a/b/c").Contains("/a/b/d"))
	assert.False(t, goref.Corpus("/a/b/c").Contains("/"))
	assert.False(t, goref.Corpus("/a/b/c").Contains("/a/b/"))
	assert.False(t, goref.Corpus("").Contains("/a/b"))
	assert.False(t, goref.Corpus("").Contains("/a/b/"))
}

func TestRel(t *testing.T) {
	assert.Equal(t, "d",
		goref.Corpus("/a/b/c").Rel("/a/b/c/d"))
	assert.Equal(t, ".",
		goref.Corpus("/a/b/c/d").Rel("/a/b/c/d"))
	assert.Equal(t, "a/b/c/d",
		goref.Corpus("/").Rel("/a/b/c/d"))
	assert.Equal(t, "/a/b/d",
		goref.Corpus("/a/b/c").Rel("/a/b/d"))
	assert.Equal(t, "/",
		goref.Corpus("/a/b/c").Rel("/"))
	assert.Equal(t, "/a/b/",
		goref.Corpus("/a/b/c").Rel("/a/b/"))
}

func TestNewCorpus_positive(t *testing.T) {
	c, err := goref.NewCorpus("/a/b/c")
	assert.NoError(t, err)
	assert.Equal(t, "/a/b/c", string(c))
}

func TestNewCorpus_relativePath(t *testing.T) {
	_, err := goref.NewCorpus("a/b/c")
	assert.Error(t, err)
}

func TestNewCorpus_emptyPath(t *testing.T) {
	_, err := goref.NewCorpus("")
	assert.Error(t, err)
}

func TestAbs(t *testing.T) {
	assert.Equal(t, "/a/b/c/d", goref.Corpus("/a/b").Abs("c/d"))
	assert.Equal(t, "/a/b/c/d", goref.Corpus("/a/b").Abs("/c/d"))
	assert.Equal(t, "/a/b", goref.Corpus("/a/b").Abs(""))
	assert.Equal(t, "", goref.Corpus("").Abs("c/d"))
}

func TestPkg(t *testing.T) {
	assert.Equal(t, "c/d", goref.Corpus("/a/b").Pkg("c/d/e"))
	assert.Equal(t, "c/d", goref.Corpus("/a/b").Pkg("c/d/"))
	assert.Equal(t, "", goref.Corpus("").Pkg("c/d/"))
	assert.Equal(t, "", goref.Corpus("/a/b").Pkg(""))
}

func TestDefaultCorpora(t *testing.T) {
	assert.NotEmpty(t, goref.DefaultCorpora())
}

func TestContainsRel(t *testing.T) {
	corpora := goref.DefaultCorpora()
	for _, c := range corpora {
		if c.ContainsRel("github.com/korfuri/goref/corpus_test.go") {
			// If ContainsRel returns true for any corpus,
			// this is considered a success.
			return
		}
	}
	t.Fail()
}
