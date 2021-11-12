package datastore

import "context"

// Store is the interface that wraps common datastore operations.
type Store interface {
	Get(namespace, key []byte) ([]byte, error)
	GetGlobal(namespace []byte) ([]byte, error)
	GetAll(namespace []byte) (map[string][]byte, error)
	Set(namespace, key, value []byte) error
	SetGlobal(namespace, value []byte) error
	Has(namespace, key []byte) (bool, error)
	RemovePrefix(namespace, key []byte) error
	Close() error
	DropAll() error
	Dir() string
}

// TODO: Seeder must be implemented by services that need to seed DB with data before server startup.
type Seeder interface {
	Seed(ctx context.Context) error
}
