package node

import (
	"context"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/eth"
	"github.com/ethereum/go-ethereum/node"
	"github.com/zees-dev/zeth/pkg/geth"
)

// TODO: Ethereum testnet IDs
type TestNetworkID int

// TODO: remove this??
type (
	IPCOptions struct {
		Enabled bool   `json:"enabled"` // an empty path disables IPC
		Path    string `json:"path"`
	}

	HTTPOptions struct {
		Enabled    bool   `json:"enabled"`
		Address    string `json:"address"`
		Port       int    `json:"port"`
		API        string `json:"api"`
		CORSDomain string `json:"corsDomain"`
		RPCPrefix  string `json:"rpcPrefix"`
	}

	WebSocketOptions struct {
		Enabled   bool   `json:"enabled"`
		Address   string `json:"address"`
		Port      int    `json:"port"`
		API       string `json:"api"`
		Origins   string `json:"origins"`
		RPCPrefix string `json:"rpcPrefix"`
	}

	GraphQLOptions struct {
		Enabled    bool   `json:"enabled"`
		CORSDomain string `json:"corsDomain"`
	}

	DevModeOptions struct {
		Enabled bool `json:"enabled"`
		Period  int  `json:"period"`
	}

	MineOptions struct {
		Enabled   bool     `json:"enabled"`
		GasLimit  uint64   `json:"gasLimit"`
		GasPrice  *big.Int `json:"gasPrice"`
		Etherbase string   `json:"etherbase"`
	}

	GasPriceOracleOptions struct {
		Blocks int `json:"blocks"`
	}

	MetricsOptions struct {
		Enabled   bool `json:"enabled"`
		Expensive bool `json:"expensive"`
		Port      int  `json:"port"`
	}
)

type GethInProcessNode struct {
	ZethNode
	GethConfig geth.GethConfig `json:"gethConfig"`
	stack      *node.Node      // non-nil value only for in-process nodes
	mux        *sync.RWMutex
}

func NewGethInProcessNode() *GethInProcessNode {
	return &GethInProcessNode{
		ZethNode: ZethNode{
			NodeType: TypeGethNodeInProcess,
		},
	}
}

func (n *GethInProcessNode) Properties() ZethNode {
	return n.ZethNode
}

func (n *GethInProcessNode) SetProperties(properties ZethNode) {
	n.ZethNode = properties
}

func (n *GethInProcessNode) Stack() *node.Node {
	return n.stack
}

// Start starts the node
func (n *GethInProcessNode) Start(ctx context.Context) error {
	n.mux.Lock()
	defer n.mux.Unlock()

	n.Enabled = true

	stack, err := node.New(&n.GethConfig.Node)
	if err != nil {
		return err
	}
	defer stack.Close() // you may want to keep this running

	_, err = eth.New(stack, &n.GethConfig.Eth)
	if err != nil {
		return err
	}

	if err = stack.Start(); err != nil {
		return err
	}

	n.stack = stack

	// log.Debug().Msgf("My datadir: %s", stack.DataDir())
	// log.Debug().Msgf("My address: %s", stack.Server().NodeInfo().ListenAddr)

	// rpcClient, err := stack.Attach()
	// if err != nil {
	// 	return err
	// }

	// ethClient := ethclient.NewClient(rpcClient)

	// progress, err := ethClient.SyncProgress(ctx)
	// if err != nil {
	// 	return err
	// }

	// log.Debug().Msgf("Sync progress: %d", progress)

	// block, err := ethClient.BlockByNumber(ctx, nil)
	// if err != nil {
	// 	return err
	// }

	// log.Debug().Msgf("Latest block: %d", block.Number())

	go stack.Wait()

	return nil
}

func (n *GethInProcessNode) Stop(ctx context.Context) error {
	n.mux.Lock()
	defer n.mux.Unlock()

	n.Enabled = false
	return n.stack.Close()
}
