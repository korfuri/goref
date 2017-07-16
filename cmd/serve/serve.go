package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path/filepath"

	"golang.org/x/net/context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/korfuri/goref"
	pb "github.com/korfuri/goref/cmd/serve/proto"
	gorefelastic "github.com/korfuri/goref/elasticsearch"
	gorefpb "github.com/korfuri/goref/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	elastic "gopkg.in/olivere/elastic.v5"
)

const (
	gatewayListenAddr = ":6061"
	grpcListenAddr    = "127.0.0.1:6062"
)

var (
	elasticURL = flag.String("elastic_url", "http://localhost:9200",
		"URL of the ElasticSearch cluster.")
	elasticUsername = flag.String("elastic_user", "elastic",
		"Username to authenticate with ElasticSearch.")
	elasticPassword = flag.String("elastic_password", "changeme",
		"Password to authenticate with ElasticSearch.")
	elasticIndex = flag.String("elastic_index", "goref",
		"Name of the index to use in ElasticSearch.")
)

// server implements pb.GorefServer
type server struct {
	corpora []goref.Corpus
	client  *elastic.Client
}

func (s server) GetAnnotations(ctx context.Context, req *pb.GetAnnotationsRequest) (*pb.GetAnnotationsResponse, error) {
	fpath := req.Path
	_, err := s.findCorpus(fpath)
	if err != nil {
		return nil, err
	}

	termQuery := elastic.NewTermQuery("to.position.filename.keyword", fpath)
	action := s.client.Search().
		Index(*elasticIndex).
		Query(termQuery).
		Type(gorefelastic.RefType).
		From(0).Size(1000).
		Pretty(false)
	searchResult, err := action.Do(ctx)
	if err != nil {
		return nil, err
	}

	res := &pb.GetAnnotationsResponse{
		Path: fpath,
	}
	for _, hit := range searchResult.Hits.Hits {
		var r gorefpb.Ref
		json.Unmarshal(*hit.Source, &r)
		res.Annotation = append(res.Annotation, &r)
	}

	return res, nil
}

func (s server) GetFiles(ctx context.Context, req *pb.GetFilesRequest) (*pb.GetFilesResponse, error) {
	termQuery := elastic.NewTermQuery("package.keyword", req.Package)
	action := s.client.Search().
		Index(*elasticIndex).
		Query(termQuery).
		Type(gorefelastic.FileType).
		From(0).Size(1000).
		Pretty(false)
	searchResult, err := action.Do(ctx)
	if err != nil {
		return nil, err
	}

	res := &pb.GetFilesResponse{
		Package: req.Package,
	}
	for _, hit := range searchResult.Hits.Hits {
		var f gorefelastic.File
		json.Unmarshal(*hit.Source, &f)
		res.Filename = append(res.Filename, f.Filename)
	}

	return res, nil
}

func (s server) GetPackages(ctx context.Context, req *pb.GetPackagesRequest) (*pb.GetPackagesResponse, error) {
	action := s.client.Search().
		Index(*elasticIndex).
		Type(gorefelastic.PackageType).
		From(0).Size(1000).
		Pretty(false)
	if req.Prefix != "" {
		query := elastic.NewPrefixQuery("loadpath.keyword", req.Prefix)
		action = action.Query(query)
	}
	searchResult, err := action.Do(ctx)
	if err != nil {
		return nil, err
	}

	res := &pb.GetPackagesResponse{}
	for _, hit := range searchResult.Hits.Hits {
		var p goref.Package
		json.Unmarshal(*hit.Source, &p)
		res.Package = append(res.Package, p.Path)
	}

	return res, nil
}

func (s server) findCorpus(fpath string) (goref.Corpus, error) {
	if filepath.Ext(fpath) != ".go" {
		return goref.Corpus(""), fmt.Errorf("Not found: invalid extension")
	}
	var corpus goref.Corpus
	for _, c := range s.corpora {
		if c.ContainsRel(fpath) {
			corpus = c
			break
		}
	}
	if corpus == "" {
		return goref.Corpus(""), fmt.Errorf("Not found under any corpus")
	}
	return corpus, nil
}

func (s server) GetFile(ctx context.Context, req *pb.GetFileRequest) (*pb.GetFileResponse, error) {
	fpath := req.Path
	corpus, err := s.findCorpus(fpath)
	if err != nil {
		return nil, err
	}
	if f, err := ioutil.ReadFile(corpus.Abs(fpath)); err == nil {
		return &pb.GetFileResponse{
			Path:     fpath,
			Contents: string(f),
		}, nil
	}
	return nil, fmt.Errorf("Internal server error")
}

func runGRPC(s *server, grpcReady chan struct{}) error {
	lis, err := net.Listen("tcp", grpcListenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	srv := grpc.NewServer()
	pb.RegisterGorefServer(srv, s)
	// Register reflection service on gRPC server.
	reflection.Register(srv)
	grpcReady <- struct{}{}
	if err := srv.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return nil
}

func runGateway(grpcReady <-chan struct{}) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	<-grpcReady

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterGorefHandlerFromEndpoint(ctx, mux, grpcListenAddr, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(gatewayListenAddr, mux)
}

func main() {
	flag.Parse()
	grpcReady := make(chan struct{})
	ec, err := elastic.NewClient(
		elastic.SetURL(*elasticURL),
		elastic.SetBasicAuth(*elasticUsername, *elasticPassword))
	if err != nil {
		panic(err)
	}
	go runGateway(grpcReady)
	s := &server{
		corpora: goref.DefaultCorpora(),
		client:  ec,
	}
	runGRPC(s, grpcReady)
}
