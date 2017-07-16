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
