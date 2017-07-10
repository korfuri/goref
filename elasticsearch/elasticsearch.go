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
)

func LoadGraphToElastic(pg goref.PackageGraph, client *elastic.Client) ([]*goref.Ref, error) {
	ctx := context.Background()
	missedRefs := make([]*goref.Ref, 0)
	errs := make([]error, 0)

	for _, p := range pg.Packages {
		docID := p.DocumentID()
		pkgDoc, err := client.Get().
			Index("goref").
			Type("package").
			Id(docID).
			Do(ctx)

		if err == nil {
			return nil, err
		}

		if pkgDoc != nil {
			log.Infof("Package %s already exists in this index.", docID)
			continue
		}

		log.Infof("Creating Package document %s in the index", docID)
		if _, err := client.Index().
			Index("goref").
			Type("package").
			Id(docID).
			BodyJson(p).
			Do(ctx); err != nil {
			return nil, err
		}

		for _, r := range p.OutRefs {
			log.Infof("Creating Ref document [%s] in the index", r)
			refDoc, err := client.Index().
				Index("goref").
				Type("ref").
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
