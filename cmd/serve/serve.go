package main

import (
	"fmt"
	"golang.org/x/net/context"
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"path/filepath"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"github.com/korfuri/goref"
	pb "github.com/korfuri/goref/cmd/serve/proto"
	"google.golang.org/grpc/reflection"
)

const (
	gatewayListenAddr = ":6061"
	grpcListenAddr = "127.0.0.1:6062"
)


// server implements pb.GorefServer
type server struct{
	corpora []goref.Corpus
}

func (s server) GetAnnotations(ctx context.Context, req *pb.GetAnnotationsRequest) (*pb.GetAnnotationsResponse, error) {
	return nil, nil
}
func (s server) GetFile(ctx context.Context, req *pb.GetFileRequest) (*pb.GetFileResponse, error) {
	fmt.Printf("path: %s\n", req.Path)
	fpath := req.Path
	if filepath.Ext(fpath) != ".go" {
		return nil, fmt.Errorf("Not found: invalid extension")
	}
	var corpus goref.Corpus
	for _, c := range s.corpora {
		if c.ContainsRel(fpath) {
			corpus = c
			break
		}
	}
	if corpus == "" {
		return nil, fmt.Errorf("Not found under any corpus")
	}
	if f, err := ioutil.ReadFile(corpus.Abs(fpath)); err == nil {
		return &pb.GetFileResponse{
			Path: fpath,
			Contents: string(f),
		}, nil
	} else {
		return nil, fmt.Errorf("Internal server error")
	}


}

func runGRPC(grpcReady chan struct{}) error {
	lis, err := net.Listen("tcp", grpcListenAddr)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	pb.RegisterGorefServer(s, &server{
		corpora: goref.DefaultCorpora(),
	})
	// Register reflection service on gRPC server.
	reflection.Register(s)
	grpcReady <- struct{}{}
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
	return nil
}

func runGateway(grpcReady <-chan struct{}) error {
	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	<- grpcReady

	mux := runtime.NewServeMux()
	opts := []grpc.DialOption{grpc.WithInsecure()}
	err := pb.RegisterGorefHandlerFromEndpoint(ctx, mux, grpcListenAddr, opts)
	if err != nil {
		return err
	}

	return http.ListenAndServe(gatewayListenAddr, mux)
}

func Do() {
	grpcReady := make(chan struct{})
	go runGateway(grpcReady)
	runGRPC(grpcReady)
}

func main() {
	Do()
}
