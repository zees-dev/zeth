package geth

// GethConfig is the geth configuration which to replicates official go-ethereum configuration

import (
	"math/big"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/ethereum/go-ethereum/eth/ethconfig"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
	"github.com/ethereum/go-ethereum/node"
	"github.com/ethereum/go-ethereum/p2p"
	"github.com/ethereum/go-ethereum/p2p/enode"
	"github.com/ethereum/go-ethereum/params"
)

// GethConfig is official un-exported struct/config from official go-ethereum codebase
type GethConfig struct {
	Eth      ethconfig.Config `json:"eth"`
	Node     node.Config      `json:"node"`
	Ethstats struct {
		URL string `toml:",omitempty"`
	} `json:",omitempty"`
	Metrics metrics.Config `json:"metrics"`
}

// DefaultGethConfig returns the default geth configuration - based on official unexported go-ethereum implementation
// source: https://github.com/ethereum/go-ethereum/blob/v1.10.11/cmd/geth/config.go#L110
var DefaultGethConfig = func() GethConfig {
	cfg := node.DefaultConfig
	// cfg.Name = clientIdentifier
	// cfg.Version = params.VersionWithCommit(gitCommit, gitDate)
	cfg.Version = params.Version
	cfg.HTTPModules = append(cfg.HTTPModules, "eth")
	cfg.WSModules = append(cfg.WSModules, "eth")
	cfg.IPCPath = "geth.ipc"

	return GethConfig{
		Eth:     ethconfig.Defaults,
		Node:    cfg,
		Metrics: metrics.DefaultConfig,
	}
}()

// SetupMainnetGethConfig creates a geth configuration
// note: only mainnet based setup is supported
func MainnetGethConfig() GethConfig {
	cfg := DefaultGethConfig

	setBootstrapNodes(&cfg.Node.P2P)
	setBootstrapNodesV5(&cfg.Node.P2P)

	utils.SetDNSDiscoveryDefaults(&cfg.Eth, params.MainnetGenesisHash)

	return cfg
}

// DevModeConfig returns a config with all settings at their defaults for the dev mode
// note: period is ignored as it is only used for the genesis block (which is not set at the moment)
func DevModeConfig(period int) GethConfig {
	cfg := DefaultGethConfig

	// ETH config
	cfg.Eth.NetworkId = 1337
	cfg.Eth.SyncMode = downloader.FullSync

	// var ks *keystore.KeyStore
	// if keystores := stack.AccountManager().Backends(keystore.KeyStoreType); len(keystores) > 0 {
	// 	ks = keystores[0].(*keystore.KeyStore)
	// }

	// developer, err = ks.NewAccount("")
	// if err != nil {
	// 	return GethConfig{}, err
	// }

	// if err := ks.Unlock(developer, passphrase); err != nil {
	// 	return GethConfig{}, err
	// }
	// log.Info("Using developer account", "address", developer.Address)

	// cfg.Eth.Genesis = core.DeveloperGenesisBlock(period, developer.Address)

	cfg.Eth.Miner.GasPrice = big.NewInt(1)

	// Node config
	cfg.Node.DataDir = ""
	cfg.Node.UseLightweightKDF = true
	cfg.Node.P2P.MaxPeers = 0
	cfg.Node.P2P.NoDiscovery = true
	cfg.Node.P2P.ListenAddr = ""
	cfg.Node.P2P.NoDial = true

	// bootstrap nodes
	setBootstrapNodes(&cfg.Node.P2P)
	setBootstrapNodesV5(&cfg.Node.P2P)

	return cfg
}

// EnableHTTP sets the HTTPHost for the config.
// This is also required to enable websocket and graphql based communication with node.
func (cfg *GethConfig) EnableHTTP() {
	cfg.Node.HTTPHost = "127.0.0.1"
}

// setBootstrapNodes creates a list of bootstrap nodes from the command line
// flags, reverting to pre-configured ones if none have been specified.
// source: https://github.com/ethereum/go-ethereum/blob/v1.10.11/cmd/utils/flags.go#L842
func setBootstrapNodes(cfg *p2p.Config) {
	urls := params.MainnetBootnodes
	// switch {
	// case ctx.GlobalIsSet(BootnodesFlag.Name):
	// 	urls = SplitAndTrim(ctx.GlobalString(BootnodesFlag.Name))
	// case ctx.GlobalBool(RopstenFlag.Name):
	// 	urls = params.RopstenBootnodes
	// case ctx.GlobalBool(RinkebyFlag.Name):
	// 	urls = params.RinkebyBootnodes
	// case ctx.GlobalBool(GoerliFlag.Name):
	// 	urls = params.GoerliBootnodes
	// case cfg.BootstrapNodes != nil:
	// 	return // already set, don't apply defaults.
	// }

	cfg.BootstrapNodes = make([]*enode.Node, 0, len(urls))
	for _, url := range urls {
		if url != "" {
			node, err := enode.Parse(enode.ValidSchemes, url)
			if err != nil {
				log.Crit("Bootstrap URL invalid", "enode", url, "err", err)
				continue
			}
			cfg.BootstrapNodes = append(cfg.BootstrapNodes, node)
		}
	}
}

// setBootstrapNodesV5 creates a list of bootstrap nodes from the command line
// flags, reverting to pre-configured ones if none have been specified.
// source: https://github.com/ethereum/go-ethereum/blob/v1.10.11/cmd/utils/flags.go#L872
func setBootstrapNodesV5(cfg *p2p.Config) {
	urls := params.V5Bootnodes
	// switch {
	// case ctx.GlobalIsSet(BootnodesFlag.Name):
	// 	urls = SplitAndTrim(ctx.GlobalString(BootnodesFlag.Name))
	// case cfg.BootstrapNodesV5 != nil:
	// 	return // already set, don't apply defaults.
	// }

	cfg.BootstrapNodesV5 = make([]*enode.Node, 0, len(urls))
	for _, url := range urls {
		if url != "" {
			node, err := enode.Parse(enode.ValidSchemes, url)
			if err != nil {
				log.Error("Bootstrap URL invalid", "enode", url, "err", err)
				continue
			}
			cfg.BootstrapNodesV5 = append(cfg.BootstrapNodesV5, node)
		}
	}
}
