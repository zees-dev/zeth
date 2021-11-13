package node

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"

	badger "github.com/dgraph-io/badger/v3"
	uuid "github.com/satori/go.uuid"
	"github.com/zees-dev/zeth/pkg/datastore"
)

var nodeKey = []byte("node")

var ErrUnknownNodeType = errors.New("unknown node type")
var ErrNodeTypeNotFound = errors.New("node type not found")

type nodeService struct {
	store datastore.Store
	cache ReverseProxyCache
}

func NewService(store datastore.Store) *nodeService {
	rpCache := NewRPCProxyCache(defaultProxyCacheSize)
	return &nodeService{
		store: store,
		cache: rpCache,
	}
}

// Create assigns a UUID to the node and saves it to the database.
// UUID assignment mutates the node struct.
func (ns *nodeService) Create(ctx context.Context, n SupportedNode) (SupportedNode, error) {
	id := uuid.NewV4()

	// set ID on existing props
	props := n.Properties()
	props.ID = id
	n.SetProperties(props)

	bodyBytes := new(bytes.Buffer)

	json.NewEncoder(bodyBytes).Encode(n)

	return n, ns.store.Set(nodeKey, id.Bytes(), bodyBytes.Bytes())
}

func (ns *nodeService) Get(ctx context.Context, id uuid.UUID) (SupportedNode, error) {
	dbNode, err := ns.store.Get(nodeKey, id.Bytes())
	if err != nil {
		return nil, err
	}

	node, err := unmarshal(dbNode)
	if err != nil {
		return nil, err
	}

	return node, nil
}

func (ns *nodeService) Update(ctx context.Context, id uuid.UUID, node SupportedNode) error {
	ok, err := ns.store.Has(nodeKey, id.Bytes())
	if err != nil {
		return err
	}
	if !ok {
		return badger.ErrKeyNotFound
	}

	bodyBytes := new(bytes.Buffer)
	json.NewEncoder(bodyBytes).Encode(node)
	return ns.store.Set(nodeKey, id.Bytes(), bodyBytes.Bytes())
}

func (ns *nodeService) Delete(ctx context.Context, id uuid.UUID) error {
	return ns.store.RemovePrefix(nodeKey, id.Bytes())
}

func (ns *nodeService) GetAll(ctx context.Context) ([]SupportedNode, error) {
	results := []SupportedNode{}

	nodesMap, err := ns.store.GetAll(nodeKey)
	if err != nil {
		return nil, err
	}

	for _, s := range nodesMap {
		node, err := unmarshal(s)
		if err != nil {
			return nil, err
		}
		results = append(results, node)
	}

	return results, nil
}

func (ns *nodeService) ReverseProxyCache() ReverseProxyCache {
	return ns.cache
}

func unmarshal(b []byte) (SupportedNode, error) {
	var anyStruct map[string]interface{}
	if err := json.Unmarshal(b, &anyStruct); err != nil {
		return nil, err
	}

	nodeTypeVal, ok := anyStruct["nodeType"]
	if !ok {
		return nil, ErrNodeTypeNotFound
	}

	switch NodeType(nodeTypeVal.(float64)) {
	case TypeGethNodeInProcess:
		var node GethInProcessNode
		if err := json.Unmarshal(b, &node); err != nil {
			return nil, err
		}
		return &node, nil
	case TypeRemoteNode:
		var node RemoteNode
		if err := json.Unmarshal(b, &node); err != nil {
			return nil, err
		}
		return &node, nil
	default:
		return nil, ErrUnknownNodeType
	}
}
