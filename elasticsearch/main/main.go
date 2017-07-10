package main

import (
	"flag"

	"github.com/korfuri/goref"
	"github.com/korfuri/goref/elasticsearch"
	log "github.com/sirupsen/logrus"
	elastic "gopkg.in/olivere/elastic.v5"
)

const (
	Usage = `elastic_goref -version 42 -include_tests <true|false> \\
  -elastic_url http://localhost:9200/ -elastic_user elastic -elastic_password changeme \\
  github.com/korfuri/goref github.com/korfuri/goref/elastic/main`
)

var (
	version = flag.Int64("version", -1,
		"Version of the code being examined. Should increase monotonically when the code is updated.")
	includeTests = flag.Bool("include_tests", true,
		"Whether XTest packages should be included in the index.")
	elasticUrl = flag.String("elastic_url", "http://localhost:9200",
		"URL of the ElasticSearch cluster.")
	elasticUsername = flag.String("elastic_user", "elastic",
		"Username to authenticate with ElasticSearch.")
	elasticPassword = flag.String("elastic_password", "changeme",
		"Password to authenticate with ElasticSearch.")
)

func usage() {
	log.Fatal(Usage)
}

func main() {
	flag.Parse()
	args := flag.Args()

	if *version == -1 || len(args) == 0 {
		usage()
	}

	// Create a client
	client, err := elastic.NewClient(
		elastic.SetURL(*elasticUrl),
		elastic.SetBasicAuth(*elasticUsername, *elasticPassword))
	if err != nil {
		log.Fatal(err)
	}

	// Filter out packages that already exist at this version in
	// the index.
	packages := make([]string, 0)
	for _, a := range args {
		if !elasticsearch.PackageExists(a, *version, client) {
			packages = append(packages, a)
		}
	}

	// Index the requested packages
	log.Infof("Indexing packages: %v", packages)
	if *includeTests {
		log.Info("This index will include XTests.")
	}
	pg := goref.NewPackageGraph(0)
	pg.LoadPrograms(packages, *includeTests)
	log.Info("Computing the interface-implementation matrix.")
	pg.ComputeInterfaceImplementationMatrix()

	log.Infof("%d packages in the graph.", len(pg.Packages))
	log.Infof("%d files in the graph.", len(pg.Files))

	// Load the indexed references into ElasticSearch
	log.Info("Inserting references into ElasticSearch.")
	if missed, err := elasticsearch.LoadGraphToElastic(*pg, client); err != nil {
		log.Fatalf("Couldn't load %d references. Error: %s", len(missed), err)
	}
	log.Info("Done, bye.")
}
