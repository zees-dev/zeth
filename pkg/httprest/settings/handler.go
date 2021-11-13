package settings

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/zees-dev/zeth/pkg/app"
	"github.com/zees-dev/zeth/pkg/httprest/rest"
	"github.com/zees-dev/zeth/pkg/node"
	"github.com/zees-dev/zeth/pkg/settings"
)

type settingsHandler struct {
	settings settings.Settings
	nodes    node.NodeService
}

func RegisterRoutes(app *app.App, baseRouter *mux.Router) {
	h := settingsHandler{
		settings: app.Services.Settings,
		nodes:    app.Services.Nodes,
	}

	baseRouter.HandleFunc("/settings", h.get).Methods(http.MethodGet)
	baseRouter.HandleFunc("/settings/node", h.updateDefaultNode).Methods(http.MethodPut)
}

/* curl request:
curl \
	-H "Content-Type: application/json" \
	http://localhost:7000/api/v1/settings
*/
func (h *settingsHandler) get(w http.ResponseWriter, r *http.Request) {
	settings, err := h.settings.Get(r.Context())
	if err != nil {
		http.Error(w, rest.HTTPInternalServerError, http.StatusInternalServerError)
		return
	}

	rest.JSON(w, settings)
}

type settingsNodeUpdateRequestBody struct {
	Uuid string `json:"uuid"`
}

func (s *settingsNodeUpdateRequestBody) Validate() url.Values {
	errs := url.Values{}

	if _, err := uuid.FromString(s.Uuid); err != nil {
		errs.Add("uuid", "invalid uuid")
	}

	return errs
}

/* curl request:
curl -X PUT \
	-H "Content-Type: application/json" \
	-d '{"uuid": "00000000-0000-0000-0000-000000000000"}' \
	http://localhost:7000/api/v1/settings/node
*/
func (h *settingsHandler) updateDefaultNode(w http.ResponseWriter, r *http.Request) {
	payload := settingsNodeUpdateRequestBody{}
	if ok := rest.DecodeAndValidateJSONPayload(w, r.Body, &payload); !ok {
		log.Debug().Msg("validation failed")
		return
	}

	// get id from request parameters
	uid, _ := uuid.FromString(payload.Uuid)

	if _, err := h.nodes.Get(r.Context(), uid); err != nil {
		http.Error(w, rest.HTTPNotFound, http.StatusNotFound)
		return
	}

	s, err := h.settings.Get(r.Context())
	if err != nil {
		http.Error(w, rest.HTTPInternalServerError, http.StatusInternalServerError)
		return
	}

	// update default node ID
	s.NodeSettings.DefaultNodeID = uid

	if err := h.settings.Update(r.Context(), s); err != nil {
		http.Error(w, rest.HTTPInternalServerError, http.StatusInternalServerError)
		return
	}

	rest.JSON(w, s)
}
