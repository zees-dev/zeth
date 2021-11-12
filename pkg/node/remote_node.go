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

type RemoteNode struct {
	ZethNode
	RPC RPC `json:"rpc"`
}

func NewRemoteNode(httpRPCURL, wsRPCURL string) *RemoteNode {
	return &RemoteNode{
		ZethNode: ZethNode{
			ID:       uuid.NewV4(),
			NodeType: TypeRemoteNode,
		},
		RPC: RPC{
			HTTP: httpRPCURL,
			WS:   wsRPCURL,
		},
	}
}

func (n *RemoteNode) Properties() ZethNode {
	return n.ZethNode
}

func (n *RemoteNode) SetProperties(properties ZethNode) {
	n.ZethNode = properties
}

// TestConnection returns true if the node can be connected to
func (n *RemoteNode) TestConnection(ctx context.Context) error {
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
