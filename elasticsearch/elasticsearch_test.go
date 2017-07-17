package elasticsearch_test

import (
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/korfuri/goref"
	"github.com/korfuri/goref/elasticsearch"
	"github.com/korfuri/goref/elasticsearch/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	elastic "gopkg.in/olivere/elastic.v5"
)

func TestPackageExists(t *testing.T) {
	client := &mocks.Client{}
	client.On("GetPackage", mock.Anything, "v1@1@fmt").Return(
		&elastic.GetResult{}, nil)
	assert.True(t, elasticsearch.PackageExists("fmt", 1, client))
	client.On("GetPackage", mock.Anything, "v1@2@log").Return(
		nil, errors.New("not found"))
	assert.False(t, elasticsearch.PackageExists("log", 2, client))
}

func TestLoadGraphToElastic_emptyGraph(t *testing.T) {
	pg := goref.PackageGraph{}
	client := &mocks.Client{}
	assert.NoError(t, elasticsearch.LoadGraphToElastic(pg, client))
}

func TestLoadGraphToElastic_allPkgsExist(t *testing.T) {
	const pkgpath = "github.com/korfuri/goref/testprograms/simple"

	client := &mocks.Client{}
	// All packages exist.
	client.On("GetPackage", mock.Anything, mock.Anything).Return(&elastic.GetResult{}, nil)

	pg := goref.NewPackageGraph(goref.ConstantVersion(0))
	pg.LoadPackages([]string{pkgpath}, false)

	assert.NoError(t, elasticsearch.LoadGraphToElastic(*pg, client))
}

func TestLoadGraphToElastic_somePkgsDontExist(t *testing.T) {
	const (
		pkgpath = "github.com/korfuri/goref/testprograms/simple"
		simple  = "v1@0@github.com/korfuri/goref/testprograms/simple"
		fmt     = "v1@0@fmt"
	)

	client := &mocks.Client{}

	// There are many other packages loaded transitively by
	// fmt. Let's say they already exist.
	client.On("GetPackage", mock.Anything, mock.MatchedBy(func(x string) bool { return (x != simple && x != fmt) })).Return(&elastic.GetResult{}, nil)
	// fmt and simple don't exist for this test.
	client.On("GetPackage", mock.Anything, simple).Return(nil, errors.New("not found"))
	client.On("GetPackage", mock.Anything, fmt).Return(nil, errors.New("not found"))

	// Creating packages, files and refs always works
	client.On("CreatePackage", mock.Anything, mock.Anything).Times(2).Return(nil)
	client.On("CreateFile", mock.Anything, mock.Anything).Return(&elastic.IndexResponse{}, nil)
	client.On("CreateRef", mock.Anything, mock.Anything).Return(&elastic.IndexResponse{}, nil)

	pg := goref.NewPackageGraph(goref.ConstantVersion(0))
	pg.LoadPackages([]string{pkgpath}, false)

	assert.NoError(t, elasticsearch.LoadGraphToElastic(*pg, client))
}

func TestLoadGraphToElastic_pkgFailsToInsert(t *testing.T) {
	const (
		pkgpath = "github.com/korfuri/goref/testprograms/interfaces"
		main    = "v1@0@github.com/korfuri/goref/testprograms/interfaces"
		lib     = "v1@0@github.com/korfuri/goref/testprograms/interfaces/lib"
	)

	client := &mocks.Client{}

	// fmt and simple don't exist for this test.
	client.On("GetPackage", mock.Anything, main).Return(nil, errors.New("not found"))
	client.On("GetPackage", mock.Anything, lib).Return(nil, errors.New("not found"))

	// Creating packages, files and refs always works, except to
	// create lib.
	matchMain := func(p *goref.Package) bool {
		t.Logf("Doc id: %s", p.DocumentID())
		return p.DocumentID() == main
	}
	matchLib := func(p *goref.Package) bool {
		return p.DocumentID() == lib
	}
	client.On("CreatePackage", mock.Anything, mock.MatchedBy(matchLib)).Return(errors.New("failed to insert lib"))
	client.On("CreatePackage", mock.Anything, mock.MatchedBy(matchMain)).Return(nil)

	// Files and refs are created without issues
	client.On("CreateFile", mock.Anything, mock.Anything).Return(&elastic.IndexResponse{}, nil)
	client.On("CreateRef", mock.Anything, mock.Anything).Return(&elastic.IndexResponse{}, nil)

	pg := goref.NewPackageGraph(goref.ConstantVersion(0))
	pg.LoadPackages([]string{pkgpath}, false)

	assert.Error(t, elasticsearch.LoadGraphToElastic(*pg, client))
}

func TestLoadGraphToElastic_fileFailsToInsert(t *testing.T) {
	const (
		pkgpath = "github.com/korfuri/goref/testprograms/interfaces"
		main    = "v1@0@github.com/korfuri/goref/testprograms/interfaces"
		lib     = "v1@0@github.com/korfuri/goref/testprograms/interfaces/lib"
	)

	client := &mocks.Client{}

	// fmt and simple don't exist for this test.
	client.On("GetPackage", mock.Anything, main).Return(nil, errors.New("not found"))
	client.On("GetPackage", mock.Anything, lib).Return(nil, errors.New("not found"))

	// Creating packages, files and refs always works, except to
	// create lib's file.
	client.On("CreatePackage", mock.Anything, mock.Anything).Return(nil)
	matchLibGo := func(f elasticsearch.File) bool {
		return f.Filename == "github.com/korfuri/goref/testprograms/interfaces/lib/lib.go"
	}
	matchNotLibGo := func(f elasticsearch.File) bool {
		return !matchLibGo(f)
	}

	client.On("CreateFile", mock.Anything, mock.MatchedBy(matchNotLibGo)).Return(&elastic.IndexResponse{}, nil)
	client.On("CreateFile", mock.Anything, mock.MatchedBy(matchLibGo)).Return(nil, errors.New("cannot create file lib.go"))
	client.On("CreateRef", mock.Anything, mock.Anything).Return(&elastic.IndexResponse{}, nil)

	pg := goref.NewPackageGraph(goref.ConstantVersion(0))
	pg.LoadPackages([]string{pkgpath}, false)

	err := elasticsearch.LoadGraphToElastic(*pg, client)
	assert.Error(t, err)
	assert.Equal(t, "1 entries couldn't be imported. Errors were:\ncannot create file lib.go\n", err.Error())
}

func TestLoadGraphToElastic_refFailsToInsert(t *testing.T) {
	const (
		pkgpath = "github.com/korfuri/goref/testprograms/interfaces"
		main    = "v1@0@github.com/korfuri/goref/testprograms/interfaces"
		lib     = "v1@0@github.com/korfuri/goref/testprograms/interfaces/lib"
	)

	client := &mocks.Client{}

	// fmt and simple don't exist for this test.
	client.On("GetPackage", mock.Anything, main).Return(nil, errors.New("not found"))
	client.On("GetPackage", mock.Anything, lib).Return(nil, errors.New("not found"))

	// Creating packages, files and refs always works, except to
	// create main's outrefs to lib.
	client.On("CreatePackage", mock.Anything, mock.Anything).Return(nil)
	client.On("CreateFile", mock.Anything, mock.Anything).Return(&elastic.IndexResponse{}, nil)
	matchLibToMain := func(r *goref.Ref) bool {
		return (r.FromPackage.DocumentID() == main &&
			r.ToPackage.DocumentID() == lib)
	}
	matchNotLibToMain := func(r *goref.Ref) bool {
		return !matchLibToMain(r)
	}
	client.On("CreateRef", mock.Anything, mock.MatchedBy(matchNotLibToMain)).Return(&elastic.IndexResponse{}, nil)
	client.On("CreateRef", mock.Anything, mock.MatchedBy(matchLibToMain)).Return(nil, errors.New("cannot create ref"))

	pg := goref.NewPackageGraph(goref.ConstantVersion(0))
	pg.LoadPackages([]string{pkgpath}, false)

	err := elasticsearch.LoadGraphToElastic(*pg, client)
	assert.Error(t, err)
	assert.Equal(t, "2 entries couldn't be imported. Errors were:\ncannot create ref\ncannot create ref\n", err.Error())
}

func TestLoadGraphToElastic_allRefsFailToInsert(t *testing.T) {
	const (
		pkgpath = "github.com/korfuri/goref/testprograms/simple"
	)

	client := &mocks.Client{}

	client.On("GetPackage", mock.Anything, mock.Anything).Return(nil, errors.New("not found"))
	client.On("CreatePackage", mock.Anything, mock.Anything).Return(nil)
	client.On("CreateFile", mock.Anything, mock.Anything).Return(&elastic.IndexResponse{}, nil)
	client.On("CreateRef", mock.Anything, mock.Anything).Return(nil, errors.New("cannot create ref"))

	pg := goref.NewPackageGraph(goref.ConstantVersion(0))
	pg.LoadPackages([]string{pkgpath}, false)

	err := elasticsearch.LoadGraphToElastic(*pg, client)
	assert.Error(t, err)
	// Errors are capped at 20 reported in the error message
	assert.Contains(t, err.Error(), "entries couldn't be imported. Errors were:")
	var n int
	fmt.Sscanf(err.Error(), "%d entries couldn't be imported.", &n)
	assert.True(t, n > 20)
	r := regexp.MustCompile("\n")
	// There's an extra \n due to the leading message.
	assert.Equal(t, 21, len(r.FindAllStringIndex(err.Error(), -1)))
}
