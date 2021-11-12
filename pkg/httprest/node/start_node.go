package node

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/zees-dev/zeth/pkg/httprest/rest"
	pnode "github.com/zees-dev/zeth/pkg/node"
)

type startNodeRequestPayload struct {
	Name string `json:"name"`
}

func (a *startNodeRequestPayload) Validate() url.Values {
	errs := url.Values{}
	// check if the title empty
	if a.Name == "" {
		errs.Add("title", "The title field is required!")
	}
	return errs
}

/* curl request:
curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"title": "dark"}' \
	http://localhost:7000/api/v1/nodes/{uuid}/start
*/
func (h *nodesHandler) startNode(w http.ResponseWriter, r *http.Request) {
	payload := startNodeRequestPayload{}
	if ok := rest.DecodeAndValidateJSONPayload(w, r.Body, &payload); !ok {
		log.Debug().Msg("validation failed")
		return
	}

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

	if err := node.(interface{}).(*pnode.GethInProcessNode).Start(r.Context()); err != nil {
		http.Error(w, rest.HTTPInternalServerError, http.StatusInternalServerError)
		return
	}

	rest.JSON(w, node)
}
