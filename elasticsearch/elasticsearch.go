package elasticsearch

import (
	"context"
	"errors"
	"fmt"

	"github.com/korfuri/goref"
	log "github.com/sirupsen/logrus"
	elastic "gopkg.in/olivere/elastic.v5"
)

const (
	// Max number of errors reported in one call to
	// LoadGraphToElastic
	maxErrorsReported = 20

	// Types in the Elastic search index
	packageType = "package"
	refType = "ref"
)

// PackageExists returns whether the provided loadpath + version tuple
// exists in this index.
func PackageExists(loadpath string, version int64, client *elastic.Client, index string) bool {
	ctx := context.Background()
	docID := fmt.Sprintf("v1@%d@%s", version, loadpath)
	pkgDoc, _ := client.Get().
		Index(index).
		Type(packageType).
		Id(docID).
		Do(ctx)
	// TODO: handle errors better. Right now we assume that any
	// error is a 404 and can be ignored safely.
	return pkgDoc != nil
}

// LoadGraphToElastic loads all Packages and Refs from a PackageGraph
// to the provided ES index.
func LoadGraphToElastic(pg goref.PackageGraph, client *elastic.Client, index string) ([]*goref.Ref, error) {
	ctx := context.Background()
	missedRefs := make([]*goref.Ref, 0)
	errs := make([]error, 0)

	for _, p := range pg.Packages {
		log.Infof("Processing package %s", p.Path)

		if PackageExists(p.Path, p.Version, client, index) {
			log.Infof("Package %s already exists in this index.", p)
			continue
		}

		log.Infof("Creating Package %s in the index", p)
		if _, err := client.Index().
			Index(index).
			Type(packageType).
			Id(p.DocumentID()).
			BodyJson(p).
			Do(ctx); err != nil {
			return nil, err
		}

		for _, r := range p.OutRefs {
			log.Infof("Creating Ref document [%s] in the index", r)
			refDoc, err := client.Index().
				Index(index).
				Type(refType).
				BodyJson(r).
				Do(ctx)
			if err != nil {
				missedRefs = append(missedRefs, r)
				errs = append(errs, err)
				log.Infof("Create Ref document failed with err:[%s] for Ref:[%s]", err, r)
			} else {
				log.Infof("Created Ref document with docID:[%s] for Ref:[%s]", refDoc.Id, r)
			}
		}
	}
	if len(missedRefs) > 0 {
		errStr := fmt.Sprintf("%d refs couldn't be imported. Errors were:\n", len(missedRefs))
		c := 0
		for _, e := range errs {
			errStr = errStr + e.Error() + "\n"
			c = c + 1
			if c >= maxErrorsReported {
				break
			}
		}
		return missedRefs, errors.New(errStr)
	}
	return nil, nil
}
