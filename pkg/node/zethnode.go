package node

import (
	"context"

	"github.com/ethereum/go-ethereum/ethclient"
	uuid "github.com/satori/go.uuid"
)

type RPC struct {
	HTTP string `json:"http"`
	WS   string `json:"ws"`
}

func NewNode(httpRPCURL, wsRPCURL string) *ZethNode {
	return &ZethNode{
		ID:       uuid.NewV4(),
		NodeType: TypeRemoteNode,
		Enabled:  true,
		RPC: RPC{
			HTTP: httpRPCURL,
			WS:   wsRPCURL,
		},
	}
}

// TestConnection returns true if the node can be connected to
func (n *ZethNode) TestConnection(ctx context.Context) error {
	client, err := ethclient.Dial(n.RPC.HTTP)
	if err != nil {
		return err
	}

	_, err = client.BlockNumber(ctx)
	if err != nil {
		return err
	}

	return nil
}
