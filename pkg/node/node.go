package node

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
)

type NodeService interface {
	Create(ctx context.Context, s SupportedNode) (SupportedNode, error)
	Get(ctx context.Context, id uuid.UUID) (SupportedNode, error)
	GetAll(ctx context.Context) ([]SupportedNode, error)
	Update(ctx context.Context, id uuid.UUID, Node SupportedNode) error
	Delete(ctx context.Context, id uuid.UUID) error
	ReverseProxyCache() ReverseProxyCache
}

// Following sites can be used to point to free ethereum nodes
// - https://docs.linkpool.io/docs/public_rpc
// - https://ethereumnodes.com/
const (
	DefaultNodeName    = "Ethereum"
	DefaultNodeHTTPRPC = "https://main-light.eth.linkpool.io"
	DefaultNodeWSRPC   = "wss://main-light.eth.linkpool.io/ws"
)

// NodeType represents the types of nodes that Zeth supports for setup
type NodeType int

const (
	_ NodeType = iota
	// GethNodeInProcess represents a geth node that runs in-process
	TypeGethNodeInProcess
	// GethNode is a ethereum node that is run externally from Zeth process
	TypeGethNode
	// GethRemoteNode is a node that is run externally from Zeth process
	TypeRemoteNode
	// // ErigonNodeInProcess (turbo geth) represents an erigon node that runs in-process
	// ErigonNodeInProcess
	// // ErigonNode is an erigon node that is run externally from Zeth process
	// ErigonNode
)

// Ethereum mainnet and sidechain Network IDs
type NetworkID uint64

const (
	// Mainnet is the default eth network ID
	Mainnet           NetworkID = 1
	BinanceSmartChain NetworkID = 56
	Polygon           NetworkID = 137
)

func (id NetworkID) IsSupported() bool {
	switch id {
	case Mainnet, BinanceSmartChain, Polygon:
		return true
	default:
		return false
	}
}

type ZethNode struct {
	ID       uuid.UUID `json:"id"`
	Name     string    `json:"name"`
	NodeType NodeType  `json:"nodeType"`
	IsDev    bool      `json:"isDev"`
	Running  bool      `json:"running"`

	DateAdded time.Time `json:"dateAdded"`
}

// SupportedNode is the interface that must be satisfied by all supported node types
type SupportedNode interface {
	Properties() ZethNode
	SetProperties(ZethNode)
}
