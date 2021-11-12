package node

import (
	"context"
	"net/http"
	"net/url"

	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/zees-dev/zeth/pkg/httprest/rest"
	"github.com/zees-dev/zeth/pkg/node"
)

type registerNodeRequestPayload struct {
	Name            string `json:"name"`
	HTTPRPCURL      string `json:"httpRPCURL"`
	WebsocketRPCURL string `json:"websocketRPCURL"`
	TestConnection  bool   `json:"testConnection"`
}

func (payload *registerNodeRequestPayload) Validate() url.Values {
	errs := url.Values{}

	if payload.Name == "" {
		errs.Add("name", "name is required")
	}

	if payload.HTTPRPCURL == "" {
		errs.Add("httpRPCURL", "httpRPCURL is required")
	} else {
		_, err := url.Parse(payload.HTTPRPCURL)
		if err != nil {
			errs.Add("httpRPCURL", "httpRPCURL is invalid")
		}
	}

	return errs
}

// registerNode registers an externally running node with zeth
/* curl request:
curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name": "test", "rpcURL": "http://localhost:8545"}' \
	http://localhost:7000/api/v1/nodes/remote
*/
func (h *nodesHandler) registerRemoteNode(w http.ResponseWriter, r *http.Request) {
	payload := registerNodeRequestPayload{}
	if ok := rest.DecodeAndValidateJSONPayload(w, r.Body, &payload); !ok {
		log.Debug().Msg("validation failed")
		return
	}

	remoteNode := node.RemoteNode{
		ZethNode: node.ZethNode{
			ID:       uuid.NewV4(),
			Name:     payload.Name,
			NodeType: node.TypeRemoteNode,
		},
		RPC: node.RPC{
			HTTP: payload.HTTPRPCURL,
			WS:   payload.WebsocketRPCURL,
		},
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

	node, err := h.nodes.Create(r.Context(), &remoteNode)
	if err != nil {
		http.Error(w, rest.HTTPInternalServerError, http.StatusInternalServerError)
		return
	}

	rest.JSON(w, node)
}

// remoteNodeAlreadyExists checks if node with the same name or rpcURL is already registered.
func (h *nodesHandler) remoteNodeAlreadyExists(ctx context.Context, payload registerNodeRequestPayload) (bool, error) {
	nodes, err := h.nodes.GetAll(ctx)
	if err != nil {
		return false, err
	}
	for _, n := range nodes {
		if n.Properties().NodeType == node.TypeRemoteNode {
			rn := n.(*node.RemoteNode)
			if rn.Name == payload.Name {
				return true, nil
			}
			if rn.RPC.HTTP == payload.HTTPRPCURL || rn.RPC.WS == payload.WebsocketRPCURL {
				return true, nil
			}
		}
	}
	return false, nil
}
