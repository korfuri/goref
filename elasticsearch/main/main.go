package main

import (
	"context"
	"log"

	"github.com/korfuri/goref"
	elastic "gopkg.in/olivere/elastic.v5"
)

func main() {
	ctx := context.Background()

	// Create a client
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetBasicAuth("elastic", "changeme"))
	if err != nil {
		log.Fatal(err)
	}

	pg := goref.NewPackageGraph(0)
	pg.LoadProgram("github.com/korfuri/goref/main", []string{"main.go"})
	pg.ComputeInterfaceImplementationMatrix()

	for _, p := range pg.Packages {
		for _, r := range p.OutRefs {
			_, err := client.Index().
				Index("goref").
				Type("ref").
				BodyJson(r).
				Refresh("true").
				Do(ctx)
			if err != nil {
				log.Fatalf("FAIL: %s", err)
			}
		}
	}
}
