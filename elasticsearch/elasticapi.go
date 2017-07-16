package elasticsearch

import (
	"context"

	"github.com/korfuri/goref"
	elastic "gopkg.in/olivere/elastic.v5"
)

// Types in the Elastic search index
const (
	PackageType = "package"
	RefType     = "ref"
	FileType    = "file"
)

// Client is an abstraction over the underlying elastic.Client that
// provides higher-order methods to manipulate data stored in
// elasticsearch.
type Client interface {
	// GetPackage retrieves a package by docID.
	GetPackage(ctx context.Context, docID string) (*elastic.GetResult, error)

	// CreatePackage creates a goref.Package entry in the index.
	CreatePackage(ctx context.Context, p *goref.Package) error

	// CreateFile creates a File entry in the index.
	CreateFile(ctx context.Context, f File) (*elastic.IndexResponse, error)

	// CreateRef creates a goref.Ref entry in the index.
	CreateRef(ctx context.Context, r *goref.Ref) (*elastic.IndexResponse, error)
}

// File represents a mapping of a file in a package
type File struct {
	Filename string `json:"filename"`
	Package  string `json:"package"`
}

// clientImpl implements Client
type clientImpl struct {
	client *elastic.Client
	index  string
}

// NewClient initializes a clientImpl
func NewClient(client *elastic.Client, index string) Client {
	return clientImpl{
		client: client,
		index:  index,
	}
}

// GetPackage implements Client for clientImpl
func (c clientImpl) GetPackage(ctx context.Context, docID string) (*elastic.GetResult, error) {
	return c.client.Get().
		Index(c.index).
		Type(PackageType).
		Id(docID).
		Do(ctx)
}

// CreatePackage implements Client for clientImpl
func (c clientImpl) CreatePackage(ctx context.Context, p *goref.Package) error {
	_, err := c.client.Index().
		Index(c.index).
		Type(PackageType).
		Id(p.DocumentID()).
		BodyJson(p).
		Do(ctx)
	return err
}

// CreateFile implements Client for clientImpl
func (c clientImpl) CreateFile(ctx context.Context, entry File) (*elastic.IndexResponse, error) {
	return c.client.Index().
		Index(c.index).
		Type(FileType).
		BodyJson(entry).
		Do(ctx)
}

// CreateRef implements Client for clientImpl
func (c clientImpl) CreateRef(ctx context.Context, r *goref.Ref) (*elastic.IndexResponse, error) {
	return c.client.Index().
		Index(c.index).
		Type(RefType).
		BodyJson(r).
		Do(ctx)
}
