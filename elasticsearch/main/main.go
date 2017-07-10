package main

import (
	"log"

	"github.com/korfuri/goref"
	"github.com/korfuri/goref/elasticsearch"
	elastic "gopkg.in/olivere/elastic.v5"
)

func main() {
	// Create a client
	client, err := elastic.NewClient(
		elastic.SetURL("http://localhost:9200"),
		elastic.SetBasicAuth("elastic", "changeme"))
	if err != nil {
		log.Fatal(err)
	}

	pg := goref.NewPackageGraph(0)
	pg.LoadPrograms([]string{"github.com/korfuri/goref/main/main"}, true)
	pg.ComputeInterfaceImplementationMatrix()

	if missed, err := elasticsearch.LoadGraphToElastic(*pg, client); err != nil {
		log.Fatalf("Couldn't load %d references. Error: %s", len(missed), err)
	}
}
