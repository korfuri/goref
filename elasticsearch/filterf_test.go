package elasticsearch_test

import (
	"errors"
	"testing"

	"github.com/korfuri/goref/elasticsearch"
	"github.com/korfuri/goref/elasticsearch/mocks"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	elastic "gopkg.in/olivere/elastic.v5"
)

func TestFilterF(t *testing.T) {
	client := &mocks.Client{}
	client.On("GetPackage", mock.Anything, "v1@0@a").Return(&elastic.GetResult{}, nil)
	client.On("GetPackage", mock.Anything, "v1@0@b").Return(nil, errors.New("404"))

	f := elasticsearch.FilterF(client)
	assert.False(t, f("a", 0))
	assert.True(t, f("b", 0))
}
