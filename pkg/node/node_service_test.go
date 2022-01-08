package node

import (
	"context"
	"testing"

	badger "github.com/dgraph-io/badger/v3"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"github.com/zees-dev/zeth/pkg/datastore/badgerdbtest"
)

const emptyUUIDV4 = "00000000-0000-0000-0000-000000000000"

func Test_SatisfiesNodeInterface(t *testing.T) {
	is := assert.New(t)
	is.Implements((*NodeService)(nil), NewService(nil))
}

func Test_NodeCrud(t *testing.T) {
	is := assert.New(t)

	store, cleanup := badgerdbtest.MustNewTestBadgerDB()
	defer cleanup()

	nodesService := NewService(store)

	t.Run("creating node assigns UUID", func(t *testing.T) {
		n := ZethNode{}
		is.Equal(n.ID.String(), emptyUUIDV4)
		is.Equal(n.ID.String(), emptyUUIDV4)

		n2, err := nodesService.Create(context.Background(), n)
		is.NoError(err)

		is.NotEqual(n.ID.String(), emptyUUIDV4)
		is.NotEqual(n2.ID.String(), emptyUUIDV4)

		store.DropAll()
	})

	t.Run("empty node without nodetype throws error on retrieval", func(t *testing.T) {
		n := GethInProcessNode{}
		_, err := nodesService.Create(context.Background(), &n)
		is.NoError(err)

		_, err = nodesService.Get(context.Background(), n.ID)
		is.Error(err)

		store.DropAll()
	})

	t.Run("created node can be retrieved via UUID", func(t *testing.T) {
		n := NewGethInProcessNode()
		_, err := nodesService.Create(context.Background(), n)
		is.NoError(err)

		_, err = nodesService.Get(context.Background(), n.ID)
		is.NoError(err)

		store.DropAll()
	})

	t.Run("updating node can be retrieved via UUID", func(t *testing.T) {
		n := NewGethInProcessNode()
		_, err := nodesService.Create(context.Background(), n)
		is.NoError(err)

		// updating existing key with new node which has empty UUID (dont do this in code)
		err = nodesService.Update(context.Background(), n.ID, NewGethInProcessNode())
		is.NoError(err)

		s2, err := nodesService.Get(context.Background(), n.ID)
		is.NoError(err)

		is.NotEqual(n.ID, s2.ID)

		store.DropAll()
	})

	t.Run("updating non-existing node returns error", func(t *testing.T) {
		err := nodesService.Update(context.Background(), uuid.NewV4(), &GethInProcessNode{})
		is.Error(err)
		is.ErrorIs(err, badger.ErrKeyNotFound)

		store.DropAll()
	})

	t.Run("create, remove and check if node successfully removed", func(t *testing.T) {
		n := GethInProcessNode{}
		_, err := nodesService.Create(context.Background(), &n)
		is.NoError(err)

		// updating existing key with new node which has empty UUID (dont do this in code)
		err = nodesService.Delete(context.Background(), n.ID)
		is.NoError(err)

		_, err = nodesService.Get(context.Background(), n.ID)
		is.Error(err)
		is.ErrorIs(err, badger.ErrKeyNotFound)

		store.DropAll()
	})

	t.Run("get all returns all creating nodes", func(t *testing.T) {
		_, err := nodesService.Create(context.Background(), NewGethInProcessNode())
		is.NoError(err)

		_, err = nodesService.Create(context.Background(), NewGethInProcessNode())
		is.NoError(err)

		_, err = nodesService.Create(context.Background(), NewGethInProcessNode())
		is.NoError(err)

		nodeList, err := nodesService.GetAll(context.Background())
		is.NoError(err)

		is.Equal(3, len(nodeList))

		store.DropAll()
	})

}
