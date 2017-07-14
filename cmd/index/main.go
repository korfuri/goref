package main

import (
	"flag"

	"github.com/korfuri/goref"
	"github.com/korfuri/goref/elasticsearch"
	log "github.com/sirupsen/logrus"
	elastic "gopkg.in/olivere/elastic.v5"
)

const (
	Usage = `elastic_goref -include_tests <true|false> \\
  -elastic_url http://localhost:9200/ -elastic_user elastic -elastic_password changeme \\
  github.com/korfuri/goref github.com/korfuri/goref/elastic/main`
)

var (
	includeTests = flag.Bool("include_tests", true,
		"Whether XTest packages should be included in the index.")
	elasticUrl = flag.String("elastic_url", "http://localhost:9200",
		"URL of the ElasticSearch cluster.")
	elasticUsername = flag.String("elastic_user", "elastic",
		"Username to authenticate with ElasticSearch.")
	elasticPassword = flag.String("elastic_password", "changeme",
		"Password to authenticate with ElasticSearch.")
	elasticIndex = flag.String("elastic_index", "goref",
		"Name of the index to use in ElasticSearch.")
)

func usage() {
	log.Fatal(Usage)
}

func main() {
	flag.Parse()
	args := flag.Args()

	if len(args) == 0 {
		usage()
	}

	// Create a client
	client, err := elastic.NewClient(
		elastic.SetURL(*elasticUrl),
		elastic.SetBasicAuth(*elasticUsername, *elasticPassword))
	if err != nil {
		log.Fatal(err)
	}

	packages := args

	// Index the requested packages
	log.Infof("Indexing packages: %v", packages)
	if *includeTests {
		log.Info("This index will include XTests.")
	}
	pg := goref.NewPackageGraph(goref.FileMTimeVersion)
	// Set FilterF to skip any packages that exist in our index
	pg.SetFilterF(elasticsearch.FilterF(client, *elasticIndex))
	pg.LoadPrograms(packages, *includeTests)
	log.Info("Computing the interface-implementation matrix.")
	pg.ComputeInterfaceImplementationMatrix()

	log.Infof("%d packages in the graph.", len(pg.Packages))

	// Load the indexed references into ElasticSearch
	log.Info("Inserting references into ElasticSearch.")
	if err := elasticsearch.LoadGraphToElastic(*pg, client, *elasticIndex); err != nil {
		log.Fatalf("Couldn't load some references. Error: %s", err)
	}
	log.Info("Done, bye.")
}
