package elasticsearch

import (
	"context"
	"errors"
	"fmt"

	"github.com/korfuri/goref"
	log "github.com/sirupsen/logrus"
)

const (
	// Max number of errors reported in one call to
	// LoadGraphToElastic
	maxErrorsReported = 20
)

// PackageExists returns whether the provided loadpath + version tuple
// exists in this index.
func PackageExists(loadpath string, version int64, client Client) bool {
	ctx := context.Background()
	docID := fmt.Sprintf("v1@%d@%s", version, loadpath)
	pkgDoc, _ := client.GetPackage(ctx, docID)
	// TODO: handle errors better. Right now we assume that any
	// error is a 404 and can be ignored safely.
	return pkgDoc != nil
}

// LoadGraphToElastic loads all Packages and Refs from a PackageGraph
// to the provided ES index.
func LoadGraphToElastic(pg goref.PackageGraph, client Client) error {
	ctx := context.Background()
	missedRefs := make([]*goref.Ref, 0)
	missedFiles := make([]string, 0)
	errs := make([]error, 0)

	for _, p := range pg.Packages {
		if PackageExists(p.Path, p.Version, client) {
			log.Infof("Package %s already exists in this index.", p)
			continue
		}

		log.Debugf("Creating Package %s in the index", p)
		if err := client.CreatePackage(ctx, p); err != nil {
			return err
		}

		for _, f := range p.Files {
			entry := File{
				Filename: f,
				Package:  p.Name,
			}
			refDoc, err := client.CreateFile(ctx, entry)
			if err != nil {
				missedFiles = append(missedFiles, f)
				errs = append(errs, err)
				log.Debugf("Create file document failed with err:[%s] for file:[%s]", err, f)
			} else {
				log.Debugf("Created file document with docID:[%s] for file:[%s]", refDoc.Id, f)
			}
		}

		for _, r := range p.OutRefs {
			refDoc, err := client.CreateRef(ctx, r)
			if err != nil {
				missedRefs = append(missedRefs, r)
				errs = append(errs, err)
				log.Debugf("Create Ref document failed with err:[%s] for Ref:[%s]", err, r)
			} else {
				log.Debugf("Created Ref document with docID:[%s] for Ref:[%s]", refDoc.Id, r)
			}
		}
	}
	if len(errs) > 0 {
		errStr := fmt.Sprintf("%d entries couldn't be imported. Errors were:\n", len(missedRefs))
		c := 0
		for _, e := range errs {
			errStr = errStr + e.Error() + "\n"
			c = c + 1
			if c >= maxErrorsReported {
				break
			}
		}
		return errors.New(errStr)
	}
	return nil
}
