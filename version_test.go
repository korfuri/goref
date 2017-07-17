package goref_test

import (
	"go/ast"
	"go/token"
	"go/types"
	"testing"
	"time"

	"github.com/korfuri/goref"
	"github.com/stretchr/testify/assert"
	"golang.org/x/tools/go/loader"
)

// Now, as a number of nanoseconds since the Unix epoch.
var nowNanoS = int64(time.Now().UTC().Sub(time.Unix(0, 0)))

func getExampleProgram(t *testing.T) *loader.Program {
	conf := loader.Config{}
	_, err := conf.FromArgs([]string{"github.com/korfuri/goref/testprograms/simple"}, false)
	assert.NoError(t, err)
	program, err := conf.Load()
	assert.NoError(t, err)
	assert.NotNil(t, program)
	return program
}

func TestFileMTimeVersion(t *testing.T) {
	program := getExampleProgram(t)
	pi := program.Package("github.com/korfuri/goref/testprograms/simple")
	assert.NotNil(t, pi)
	v, err := goref.FileMTimeVersion(*program, *pi)
	assert.NoError(t, err)
	assert.True(t, v > 0)
	assert.True(t, v < nowNanoS)
}

func TestFileMTimeVersion_badPackage(t *testing.T) {
	program := getExampleProgram(t)
	pi := &loader.PackageInfo{
		Files: make([]*ast.File, 0),
		Pkg:   &types.Package{},
	}
	assert.NotNil(t, pi)
	_, err := goref.FileMTimeVersion(*program, *pi)
	assert.Error(t, err)
}

func TestFileMTimeVersion_badFile(t *testing.T) {
	program := getExampleProgram(t)
	pi := program.Package("github.com/korfuri/goref/testprograms/simple")
	pi.Files = append(pi.Files, &ast.File{
		Package: token.Pos(0),
	})
	assert.NotNil(t, pi)
	_, err := goref.FileMTimeVersion(*program, *pi)
	assert.Error(t, err)
}
