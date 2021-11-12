package datastore

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/dgraph-io/badger/v3"
	"github.com/rs/zerolog/log"
)

const (
	// source: https://gist.github.com/alexanderbez/d99fd0383ad57e991e9af9adcbb70b9d
	// Default BadgerDB discardRatio. It represents the discard ratio for the BadgerDB GC.
	// Ref: https://godoc.org/github.com/dgraph-io/badger#DB.RunValueLogGC
	badgerDiscardRatio = 0.5

	// Default BadgerDB GC interval
	badgerGCInterval = 10 * time.Minute

	// DefaultStoreDir is the default badgerdb store directory
	DefaultStoreDir = "zeth.db"
)

type (
	// badgerStore is a wrapper around a badgerDB backend database that implements the Database interface.
	badgerStore struct {
		db         *badger.DB
		ctx        context.Context
		cancelFunc context.CancelFunc
		isNew      bool
	}
)

// NewBadgerDB returns a new initialized BadgerDB database implementing the DB interface.
// If the database cannot be initialized, an error will be returned.
// `dbPath` is the path to the BadgerDB database. If empty string is provided, an in-memory DB will be used
func NewBadgerDB(dbPath string) (*badgerStore, error) {
	// TODO clean up this mess
	isNew := false

	var badgerOpts badger.Options
	if dbPath == "" {
		// https://dgraph.io/docs/badger/get-started/#in-memory-mode-diskless-mode
		badgerOpts = badger.DefaultOptions("").WithInMemory(true)
		isNew = true
	} else {
		if _, err := os.Stat(dbPath); os.IsNotExist(err) {
			isNew = true
		}
		badgerOpts = badger.DefaultOptions(dbPath)
	}
	badgerOpts.Logger = nil // disable badger logging

	badgerDB, err := badger.Open(badgerOpts)
	if err != nil {
		return nil, err
	}

	bdb := &badgerStore{
		db:    badgerDB,
		isNew: isNew,
	}
	bdb.ctx, bdb.cancelFunc = context.WithCancel(context.Background())

	go bdb.runGC()
	return bdb, nil
}

// GetGlobal implements the Store interface. It attempts to get a value for a given namespace
// If the namespace does not exist, an error is returned.
func (bdb *badgerStore) GetGlobal(namespace []byte) ([]byte, error) {
	var value []byte

	err := bdb.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get(namespace)
		if err != nil {
			return err
		}

		err = item.Value(func(v []byte) error {
			value = append([]byte{}, v...)
			return nil
		})
		return err
	})
	if err != nil {
		return nil, err
	}

	return value, nil
}

// Get implements the Store interface. It attempts to get a value for a given key
// and namespace. If the key does not exist in the provided namespace, an error
// is returned, otherwise the retrieved value.
func (bdb *badgerStore) Get(namespace, key []byte) ([]byte, error) {
	return bdb.GetGlobal(badgerNamespaceKey(namespace, key))
}

// GetAll implements the Store interface. It attempts to get all values in a given namespace.
func (bdb *badgerStore) GetAll(namespace []byte) (map[string][]byte, error) {
	values := make(map[string][]byte)

	err := bdb.db.View(func(txn *badger.Txn) error {
		it := txn.NewIterator(badger.DefaultIteratorOptions)
		defer it.Close()
		for it.Seek(namespace); it.ValidForPrefix(namespace); it.Next() {
			item := it.Item()
			k := item.Key()
			err := item.Value(func(v []byte) error {
				// fmt.Printf("key=%s, value=%s\n", k, v)

				// remove namespace from key
				id := strings.Split(string(k), "/")[1]
				values[id] = append([]byte{}, v...)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}

	return values, nil
}

// Set implements the Store interface. It attempts to store a value for a given key
// and namespace. If the key/value pair cannot be saved, an error is returned.
func (bdb *badgerStore) Set(namespace, key, value []byte) error {
	return bdb.SetGlobal(badgerNamespaceKey(namespace, key), value)
}

// SetGlobal implements the Store interface. It attempts to store a value for a given globally for a namespace.
// If the namespace/value pair cannot be saved, an error is returned.
func (bdb *badgerStore) SetGlobal(namespace, value []byte) error {
	err := bdb.db.Update(func(txn *badger.Txn) error {
		return txn.Set(namespace, value)
	})
	if err != nil {
		log.Err(err).Msgf("failed to set global value for namespace %s", namespace)
		return err
	}

	return nil
}

// Remove implements the Store interface. It attempts to remove a key from a given namespace.
func (bdb *badgerStore) RemovePrefix(namespace, key []byte) error {
	return bdb.db.DropPrefix(badgerNamespaceKey(namespace, key))
}

// Has implements the Store interface. It returns a boolean reflecting if the
// datbase has a given key for a namespace or not. An error is only returned if
// an error to Get would be returned that is not of type badger.ErrKeyNotFound.
func (bs *badgerStore) Has(namespace, key []byte) (bool, error) {
	_, err := bs.Get(namespace, key)
	switch err {
	case badger.ErrKeyNotFound:
		return false, nil
	case nil:
		return true, nil
	}
	return false, nil
}

// Close implements the Store interface. It closes the connection to the underlying
// BadgerDB database as well as invoking the context's cancel function.
func (bs *badgerStore) Close() error {
	bs.cancelFunc()
	return bs.db.Close()
}

// IsNew returns a boolean reflecting if the database was created or not.
func (bs *badgerStore) IsNew() bool {
	return bs.isNew
}

// DropAll drops all keys in the database, effectively wiping the entire database.
func (bs *badgerStore) DropAll() error {
	return bs.db.DropAll()
}

// runGC triggers the garbage collection for the BadgerDB backend database.
// It should be run in a goroutine.
// https://dgraph.io/docs/badger/get-started/#garbage-collection
func (bs *badgerStore) runGC() {
	ticker := time.NewTicker(badgerGCInterval)
	for {
		select {
		case <-ticker.C:
			err := bs.db.RunValueLogGC(badgerDiscardRatio)
			if err != nil {
				// don't report error when GC didn't result in any cleanup
				if err == badger.ErrNoRewrite {
					log.Debug().Msgf("no BadgerDB GC occurred: %v", err)
				} else {
					log.Debug().Err(err).Msg("failed to GC BadgerDB")
				}
			}
		case <-bs.ctx.Done():
			return
		}
	}
}

// badgerNamespaceKey returns a composite key used for lookup and storage for a given namespace and key.
// Namespaced keys will be stored in the form: `namespace`/`key`
func badgerNamespaceKey(namespace, key []byte) []byte {
	prefix := []byte(fmt.Sprintf("%s/", namespace))
	return append(prefix, key...)
}

// Dir returns db file location
func (bdb *badgerStore) Dir() string {
	return bdb.db.Opts().Dir
}
