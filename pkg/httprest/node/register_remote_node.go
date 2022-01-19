package node

import (
	"context"
	"net/http"
	"net/url"
	"time"

	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/zees-dev/zeth/pkg/httprest/rest"
	"github.com/zees-dev/zeth/pkg/node"
)

type registerNodeRequestPayload struct {
	Name string `json:"name"`
	// Enabled     bool     `json:"enabled"`
	ExplorerURL    string   `json:"explorerUrl"`
	RPC            node.RPC `json:"rpc"`
	TestConnection bool     `json:"test"`
}

func (payload *registerNodeRequestPayload) Validate() url.Values {
	errs := url.Values{}

	if payload.Name == "" {
		errs.Add("name", "name is required")
	}

	if payload.RPC.HTTP == "" {
		errs.Add("rpc.http", "rpc http url is required")
	} else {
		if _, err := url.Parse(payload.RPC.HTTP); err != nil {
			errs.Add("rpc.http", "rpc http url is invalid")
		}
	}

	if payload.RPC.WS != "" {
		if _, err := url.Parse(payload.RPC.WS); err != nil {
			errs.Add("rpc.ws", "rpc ws url is invalid")
		}
	}

	if payload.RPC.Default != node.DefaultHTTPRPC && payload.RPC.Default != node.DefaultWSRPC {
		errs.Add("rpc.default", "rpc default is invalid; must be 0 (http) or 1 (ws)")
	}

	return errs
}

// registerNode registers an externally running node with zeth
/* curl request:
curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "test", "rpc": { "http": "http://localhost:8545", "default": 0 } }' \
	http://localhost:7000/api/v1/nodes/remote
*/
func (h *nodesHandler) createNode(w http.ResponseWriter, r *http.Request) {
	payload := registerNodeRequestPayload{}
	if ok := rest.DecodeAndValidateJSONPayload(w, r.Body, &payload); !ok {
		log.Debug().Msg("validation failed")
		return
	}

	remoteNode := node.ZethNode{
		ID:        uuid.NewV4(),
		Name:      payload.Name,
		NodeType:  node.TypeRemoteNode,
		Enabled:   true,
		DateAdded: time.Now().UTC(),
		RPC:       payload.RPC,
	}

	exists, err := h.remoteNodeAlreadyExists(r.Context(), payload)
	if exists || err != nil {
		http.Error(w, "a node with rpcURL or name already exists", http.StatusBadRequest)
		return
	}

	if payload.TestConnection {
		if err := remoteNode.TestConnection(r.Context()); err != nil {
			log.Debug().Err(err).Msg("failed to connect to node")
			http.Error(w, "failed to connect to node", http.StatusBadRequest)
			return
		}
	}

	node, err := h.nodes.Create(r.Context(), remoteNode)
	if err != nil {
		http.Error(w, rest.HTTPInternalServerError, http.StatusInternalServerError)
		return
	}

	rest.JSON(w, node)
}

// remoteNodeAlreadyExists checks if node with the same name or http rpc URL is already registered.
func (h *nodesHandler) remoteNodeAlreadyExists(ctx context.Context, payload registerNodeRequestPayload) (bool, error) {
	nodes, err := h.nodes.GetAll(ctx)
	if err != nil {
		return false, err
	}
	for _, n := range nodes {
		if n.Name == payload.Name {
			return true, nil
		}
		if n.RPC.HTTP == payload.RPC.HTTP {
			return true, nil
		}
	}
	return false, nil
}
