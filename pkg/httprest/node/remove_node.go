package node

import (
	"net/http"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
	"github.com/zees-dev/zeth/pkg/httprest/rest"
)

/* curl request:
curl -X DELETE \
	http://localhost:7000/api/v1/nodes/00000000-0000-0000-0000-000000000000
*/
func (h *nodesHandler) removeNode(w http.ResponseWriter, r *http.Request) {
	// get id from request parameters
	id := mux.Vars(r)["uuid"]

	uid, err := uuid.FromString(id)
	if err != nil {
		http.Error(w, rest.HTTPBadRequest, http.StatusBadRequest)
		return
	}

	_, err = h.nodes.Get(r.Context(), uid)
	if err != nil {
		http.Error(w, rest.HTTPNotFound, http.StatusNotFound)
		return
	}

	err = h.nodes.Delete(r.Context(), uid)
	if err != nil {
		http.Error(w, rest.HTTPInternalServerError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
