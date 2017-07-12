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
