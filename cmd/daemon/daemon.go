package main

import (
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
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

const (
	gatewayListenAddr = ":6061"
	grpcListenAddr    = "127.0.0.1:6062"
)

var (
	includeTests = flag.Bool("include_tests", true,
		"Whether XTest packages should be included in the index.")
)

// server implements pb.GorefServer
type server struct {
	graph goref.PackageGraph
}

func (s server) GetAnnotations(ctx context.Context, req *pb.GetAnnotationsRequest) (*pb.GetAnnotationsResponse, error) {
	fpath := req.Path
	corpus, err := s.findCorpus(fpath)
	if err != nil {
		return nil, err
	}
	fmt.Printf("%s\n", corpus.Pkg(fpath))
	pkg, in := s.graph.Packages[corpus.Pkg(fpath)]
	if !in {
		return nil, fmt.Errorf("Internal server error")
	}

	res := &pb.GetAnnotationsResponse{
		Path: fpath,
	}
	for _, r := range pkg.InRefs {
		if r.ToPosition.File == fpath {
			res.Annotation = append(res.Annotation, r.ToProto())
		}
	}
	for _, r := range pkg.OutRefs {
		if r.FromPosition.File == fpath {
			res.Annotation = append(res.Annotation, r.ToProto())
		}
	}

	return res, nil
}

func (s server) GetFiles(ctx context.Context, req *pb.GetFilesRequest) (*pb.GetFilesResponse, error) {
	pkg, in := s.graph.Packages[req.Package]
	if !in {
		return nil, fmt.Errorf("Unknown package")
	}

	res := &pb.GetFilesResponse{
		Package: req.Package,
	}
	for _, f := range pkg.Files {
		res.Filename = append(res.Filename, f)
	}

	return res, nil
}

func (s server) GetPackages(ctx context.Context, req *pb.GetPackagesRequest) (*pb.GetPackagesResponse, error) {
	res := &pb.GetPackagesResponse{}
	for _, pkg := range s.graph.Packages {
		res.Package = append(res.Package, pkg.Path)
	}
	return res, nil
}

func (s server) findCorpus(fpath string) (goref.Corpus, error) {
	if filepath.Ext(fpath) != ".go" {
		return goref.Corpus(""), fmt.Errorf("Not found: invalid extension")
	}
	var corpus goref.Corpus
	for _, c := range s.graph.Corpora {
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
	args := flag.Args()

	// Index the requested packages
	pg := goref.NewPackageGraph(goref.FileMTimeVersion)
	pg.LoadPrograms(args, *includeTests)

	grpcReady := make(chan struct{})
	go runGateway(grpcReady)
	s := &server{
		graph: *pg,
	}
	runGRPC(s, grpcReady)
}
