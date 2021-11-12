package node

import (
	"net/http"

	"github.com/zees-dev/zeth/pkg/httprest/rest"
	"github.com/zees-dev/zeth/pkg/node"
)

type allNodesResponse struct {
	Nodes []node.SupportedNode `json:"nodes"`
}

/* curl request:
curl \
	-H "Content-Type: application/json" \
	http://localhost:7000/api/v1/nodes
*/
func (h *nodesHandler) getNodes(w http.ResponseWriter, r *http.Request) {
	nodes, err := h.nodes.GetAll(r.Context())
	if err != nil {
		http.Error(w, rest.HTTPInternalServerError, http.StatusInternalServerError)
		return
	}

	response := allNodesResponse{
		Nodes: nodes,
	}

	rest.JSON(w, response)
}
