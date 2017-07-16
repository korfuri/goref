package elasticsearch_test

import (
	"errors"
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
	// create fmt.
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
