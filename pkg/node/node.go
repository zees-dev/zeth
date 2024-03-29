package node

import (
	"context"
	"time"

	uuid "github.com/satori/go.uuid"
)

type NodeService interface {
	Create(context.Context, ZethNode) (ZethNode, error)
	Get(context.Context, uuid.UUID) (*ZethNode, error)
	GetAll(context.Context) ([]ZethNode, error)
	Update(context.Context, uuid.UUID, ZethNode) error
	Delete(context.Context, uuid.UUID) error
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
	ID          uuid.UUID `json:"id"`
	Name        string    `json:"name"`
	IsDev       bool      `json:"isDev"`
	Enabled     bool      `json:"enabled"`
	DateAdded   time.Time `json:"dateAdded"`
	ExplorerURL string    `json:"explorerUrl"`
	RPC         RPC       `json:"rpc"`
}
