package node

import (
	"context"
	"log"
	"os"
	"testing"

	"github.com/ethereum/go-ethereum/params"
	"github.com/stretchr/testify/assert"
	"github.com/zees-dev/zeth/pkg/geth"
)

func Test_EthNode(t *testing.T) {
	is := assert.New(t)

	datadir, _ := os.MkdirTemp("", "eth-node-test")
	defer os.RemoveAll(datadir)

	log.Println(params.Version)

	ecfg := NewGethInProcessNode()
	ecfg.GethConfig = geth.DevModeConfig(0)
	ecfg.GethConfig.EnableHTTP()

	err := ecfg.Start(context.Background())
	is.NoError(err)

	// err := dumpConfig(ecfg.GethConfig, os.Stdout)
	// is.NoError(err)

	// t.Error("TODO: test node")

	t.Run("updating node can be retrieved via UUID", func(t *testing.T) {

	})
}

func Test_versionCheck(t *testing.T) {
	is := assert.New(t)

	cfg := geth.DefaultGethConfig
	wantVersion := "1.10.11"

	is.Equal(wantVersion, cfg.Node.Version)
}
