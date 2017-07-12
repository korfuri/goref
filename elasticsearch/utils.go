package elasticsearch

import (
	elastic "gopkg.in/olivere/elastic.v5"
)

// FilterF returns a function suitable for goref.PackageGraph.FilterF
// that returns false if a package exists in this ElasticSearch index.
func FilterF(client *elastic.Client, index string) func(string, int64) bool {
	return func(loadpath string, version int64) bool {
		return !PackageExists(loadpath, version, client, index)
	}
}
