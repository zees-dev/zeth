package node

import (
	"net/http"
	"net/url"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/zees-dev/zeth/pkg/httprest/rest"
	"github.com/zees-dev/zeth/pkg/node"
)

type updateNodeRequestPayload struct {
	Name           string   `json:"name"`
	Enabled        bool     `json:"enabled"`
	ExplorerURL    string   `json:"explorerUrl"`
	RPC            node.RPC `json:"rpc"`
	TestConnection bool     `json:"test"`
}

func (payload *updateNodeRequestPayload) Validate() url.Values {
	errs := url.Values{}
	if len(strings.TrimSpace(payload.Name)) == 0 {
		errs.Add("name", "The name field is required!")
	}

	if len(strings.TrimSpace(payload.RPC.HTTP)) == 0 {
		errs.Add("rpc.http", "The rpc http url is required!")
	} else {
		if _, err := url.Parse(payload.RPC.HTTP); err != nil {
			errs.Add("rpc.http", "The rpc http url is invalid!")
		}
	}

	if len(strings.TrimSpace(payload.RPC.WS)) != 0 {
		if _, err := url.Parse(payload.RPC.WS); err != nil {
			errs.Add("rpc.ws", "The rpc ws url is invalid")
		}
		if payload.RPC.Default == node.DefaultRPCWS {
			errs.Add("rpc.ws", "rpc default is 1 (ws) but ws url is empty")
		}
	}

	if payload.RPC.Default != node.DefaultRPCHTTP && payload.RPC.Default != node.DefaultRPCWS {
		errs.Add("rpc.default", "rpc default is invalid; must be 0 (http) or 1 (ws)")
	}

	return errs
}

/* curl request:
curl -X PUT \
	-H "Content-Type: application/json" \
	-d '{ "name": "avalanche mainnet", "enabled": true, "explorerUrl": "", "rpc": { "http":"https://api.avax.network/ext/bc/C/rpc", "ws": "", "default": 0 }, "test": true }' \
	http://localhost:7000/api/v1/nodes/f0470395-76a6-4d30-a79a-57387d93fb1b
*/
func (h *nodesHandler) updateNode(w http.ResponseWriter, r *http.Request) {
	payload := updateNodeRequestPayload{}
	if ok := rest.DecodeAndValidateJSONPayload(w, r.Body, &payload); !ok {
		log.Debug().Msg("validation failed")
		return
	}

	// get id from request parameters
	id := mux.Vars(r)["uuid"]

	uid, err := uuid.FromString(id)
	if err != nil {
		http.Error(w, rest.HTTPBadRequest, http.StatusBadRequest)
		return
	}

	node, err := h.nodes.Get(r.Context(), uid)
	if err != nil {
		http.Error(w, rest.HTTPNotFound, http.StatusNotFound)
		return
	}

	node.Name = payload.Name
	node.Enabled = payload.Enabled
	node.ExplorerURL = payload.ExplorerURL
	node.RPC = payload.RPC

	if payload.TestConnection {
		if err := node.TestConnection(r.Context()); err != nil {
			log.Debug().Err(err).Msg("failed to connect to node")
			http.Error(w, "failed to connect to node", http.StatusBadRequest)
			return
		}
	}

	if err := h.nodes.Update(r.Context(), node.ID, *node); err != nil {
		http.Error(w, rest.HTTPInternalServerError, http.StatusInternalServerError)
		return
	}

	rest.JSON(w, node)
}
