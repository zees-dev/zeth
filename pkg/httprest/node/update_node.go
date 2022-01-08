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
	Name        string   `json:"name"`
	Enabled     bool     `json:"enabled"`
	ExplorerURL string   `json:"explorerUrl"`
	RPC         node.RPC `json:"rpc"`
}

func (a *updateNodeRequestPayload) Validate() url.Values {
	errs := url.Values{}
	if len(strings.TrimSpace(a.Name)) == 0 {
		errs.Add("name", "The name field is required!")
	}
	return errs
}

/* curl request:
curl -X PUT \
	-H "Content-Type: application/json" \
	-d '{"name": "dark"}' \
	http://localhost:7000/api/v1/nodes/00000000-0000-0000-0000-000000000000
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

	if err := h.nodes.Update(r.Context(), node.ID, *node); err != nil {
		http.Error(w, rest.HTTPInternalServerError, http.StatusInternalServerError)
		return
	}

	rest.JSON(w, node)
}
