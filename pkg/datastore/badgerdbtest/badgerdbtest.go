package badgerdbtest

import (
	"github.com/zees-dev/zeth/pkg/datastore"
)

func MustNewTestBadgerDB() (datastore.Store, func() error) {
	store, err := datastore.NewBadgerDB("")
	if err != nil {
		panic(err)
	}
	return store, store.Close
}
