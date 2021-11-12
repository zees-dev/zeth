package node

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/cmd/utils"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/eth/downloader"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/zees-dev/zeth/pkg/geth"
	"github.com/zees-dev/zeth/pkg/httprest/rest"
	"github.com/zees-dev/zeth/pkg/node"
)

type createNodeRequestPayload struct {
	NodeTypeSetting node.NodeType `json:"nodeType"`

	Name      string              `json:"name"`
	NetworkID node.NetworkID      `json:"networkID"`
	DataDir   string              `json:"dataDir"`  // utils.DirectoryString
	SyncMode  downloader.SyncMode `json:"syncMode"` // downloader.SyncMode
	// Connections
	IPC     node.IPCOptions       `json:"ipc"`
	HTTP    node.HTTPOptions      `json:"http"`
	WS      node.WebSocketOptions `json:"ws"`
	GraphQL node.GraphQLOptions   `json:"graphql"`
	// Dev mode
	DevMode node.DevModeOptions `json:"devMode"`
	// Main port
	Port int              `json:"port"`
	Mine node.MineOptions `json:"mine"`
	// GasPriceOracle node.GasPriceOracleOptions `json:"gasPriceOracle"`
	Metrics node.MetricsOptions `json:"metrics"`
}

func (payload *createNodeRequestPayload) Validate() url.Values {
	errs := url.Values{}

	// TODO: check if NodeTypeSetting exists as supported node

	// TODO: check if name doesnt have alphanumeric characters
	// TODO: name cannot contain `/`
	if payload.Name == "" {
		errs.Add("name", "name is required")
	}

	if !payload.NetworkID.IsSupported() {
		errs.Add("networkID", "provided networkID is not supported")
	}

	if !payload.SyncMode.IsValid() {
		errs.Add("syncMode", "provided syncMode is not valid")
	}

	if payload.IPC.Enabled {
		if payload.IPC.Path == "" {
			errs.Add("ipc.path", "ipc.path is required")
		}
	}

	if payload.HTTP.Enabled {
		if payload.HTTP.Address == "" {
			errs.Add("http.address", "http.address is required")
		}
	}

	if payload.WS.Enabled {
		if payload.WS.Address == "" {
			errs.Add("ws.address", "ws.address is required")
		}
	}

	if payload.GraphQL.Enabled {
		if !payload.HTTP.Enabled {
			errs.Add("http", "http must be enabled for graphql")
		}
	}

	if payload.Mine.Etherbase != "" {
		if !common.IsHexAddress(payload.Mine.Etherbase) {
			errs.Add("mine.etherbase", "mine.etherbase is not a valid address")
		}
	}

	return errs
}

/* curl request:
curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"title": "dark"}' \
	http://localhost:7000/api/v1/nodes
*/
func (h *nodesHandler) createNode(w http.ResponseWriter, r *http.Request) {
	payload := createNodeRequestPayload{}
	if ok := rest.DecodeAndValidateJSONPayload(w, r.Body, &payload); !ok {
		log.Debug().Msg("validation failed")
		return
	}

	cfg := geth.MainnetGethConfig()

	// create new zeth node with sane defaults
	gipNode := node.GethInProcessNode{
		ZethNode: node.ZethNode{
			ID:       uuid.NewV4(),
			NodeType: node.TypeGethNodeInProcess,
		},
		GethConfig: cfg,
	}

	// override default node config with payload
	overrideConfig(&gipNode, payload)

	node, err := h.nodes.Create(r.Context(), &gipNode)
	if err != nil {
		http.Error(w, rest.HTTPInternalServerError, http.StatusInternalServerError)
		return
	}

	rest.JSON(w, node)
}

func overrideConfig(node *node.GethInProcessNode, payload createNodeRequestPayload) {
	if payload.DevMode.Enabled {
		cfg := geth.DevModeConfig(payload.DevMode.Period)
		node.IsDev = true
		node.GethConfig = cfg
	}

	if payload.Name != "" {
		node.Name = payload.Name
		node.GethConfig.Node.Name = payload.Name
		node.GethConfig.Node.UserIdent = payload.Name
	}

	node.GethConfig.Eth.NetworkId = uint64(payload.NetworkID)

	if payload.DataDir != "" {
		node.GethConfig.Node.DataDir = payload.DataDir
	}

	node.GethConfig.Eth.SyncMode = payload.SyncMode

	if payload.IPC.Enabled {
		node.GethConfig.Node.IPCPath = payload.IPC.Path
	} else {
		node.GethConfig.Node.IPCPath = ""
	}

	if payload.HTTP.Enabled {
		node.GethConfig.Node.HTTPHost = payload.HTTP.Address
		if payload.HTTP.Port != 0 {
			node.GethConfig.Node.HTTPPort = payload.HTTP.Port
		}
		if payload.HTTP.API != "" {
			node.GethConfig.Node.HTTPModules = utils.SplitAndTrim(payload.HTTP.API)
		}
		if payload.HTTP.CORSDomain != "" {
			node.GethConfig.Node.HTTPCors = utils.SplitAndTrim(payload.HTTP.CORSDomain)
		}
		if payload.HTTP.RPCPrefix != "" {
			node.GethConfig.Node.HTTPPathPrefix = payload.HTTP.RPCPrefix
		}
	}

	if payload.WS.Enabled {
		node.GethConfig.Node.WSHost = payload.WS.Address
		if payload.WS.Port != 0 {
			node.GethConfig.Node.WSPort = payload.WS.Port
		}
		if payload.WS.API != "" {
			node.GethConfig.Node.WSModules = utils.SplitAndTrim(payload.WS.API)
		}
		if payload.WS.Origins != "" {
			node.GethConfig.Node.WSOrigins = utils.SplitAndTrim(payload.WS.Origins)
		}
		if payload.WS.RPCPrefix != "" {
			node.GethConfig.Node.WSPathPrefix = payload.WS.RPCPrefix
		}
	}

	if payload.GraphQL.Enabled {
		if payload.GraphQL.CORSDomain != "" {
			node.GethConfig.Node.GraphQLCors = utils.SplitAndTrim(payload.GraphQL.CORSDomain)
		}
	}

	if payload.Port != 0 {
		node.GethConfig.Node.P2P.ListenAddr = fmt.Sprintf(":%d", payload.Port)
	}

	if payload.Mine.Enabled {
		if payload.Mine.GasLimit != 0 {
			node.GethConfig.Eth.Miner.GasCeil = payload.Mine.GasLimit
		}
		if payload.Mine.GasPrice != nil {
			node.GethConfig.Eth.Miner.GasPrice = payload.Mine.GasPrice
		}
		if payload.Mine.Etherbase != "" {
			node.GethConfig.Eth.Miner.Etherbase = common.HexToAddress(payload.Mine.Etherbase)
		}
	}

	if payload.Metrics.Enabled {
		node.GethConfig.Metrics.Enabled = payload.Metrics.Enabled
		node.GethConfig.Metrics.EnabledExpensive = payload.Metrics.Expensive
		if payload.Metrics.Port != 0 {
			node.GethConfig.Metrics.Port = payload.Metrics.Port
		}
	}
}
