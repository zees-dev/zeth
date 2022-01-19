package node

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	uuid "github.com/satori/go.uuid"
)

type DefaultRPC uint

const (
	DefaultHTTPRPC DefaultRPC = iota
	DefaultWSRPC
	DefaultRPCend
)

func (defaultRPC DefaultRPC) IsValid() bool {
	return defaultRPC < DefaultRPCend
}

type RPC struct {
	HTTP    string     `json:"http"`
	WS      string     `json:"ws"`
	Default DefaultRPC `json:"default"`
}

func NewNode(httpRPCURL, wsRPCURL string) *ZethNode {
	return &ZethNode{
		ID:       uuid.NewV4(),
		NodeType: TypeRemoteNode,
		Enabled:  true,
		RPC: RPC{
			HTTP:    httpRPCURL,
			WS:      wsRPCURL,
			Default: DefaultHTTPRPC,
		},
	}
}

// TestConnection returns true if the node can be connected to via RPC HTTP endpoint.
func (n *ZethNode) TestConnection(ctx context.Context) error {
	client, err := ethclient.Dial(n.RPC.HTTP)
	if err != nil {
		return err
	}

	// _, err = client.BlockNumber(ctx)
	_, err = client.HeaderByNumber(ctx, nil)
	if err != nil {
		return err
	}

	return nil
}
