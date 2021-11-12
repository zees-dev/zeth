package settings

import (
	"net/http"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	"github.com/zees-dev/zeth/pkg/app"
	"github.com/zees-dev/zeth/pkg/httprest/rest"
	"github.com/zees-dev/zeth/pkg/settings"
)

type settingsHandler struct {
	settings settings.Settings
}

func RegisterRoutes(app *app.App, baseRouter *mux.Router) {
	h := settingsHandler{
		settings: app.Services.Settings,
	}

	baseRouter.HandleFunc("/settings", h.get).Methods(http.MethodGet)
	baseRouter.HandleFunc("/settings", h.put).Methods(http.MethodPut)
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

type settingsRequestBody struct {
	Title string `json:"title"`
}

func (a *settingsRequestBody) Validate() url.Values {
	errs := url.Values{}
	// check if the title empty
	if a.Title == "" {
		errs.Add("title", "The title field is required!")
	}
	if a.Title == "" {
		errs.Add("titlea", "The title field is requireds!")
	}
	return errs
}

/* curl request:
curl -X PUT \
	-H "Content-Type: application/json" \
	-d '{"theme": "dark"}' \
	http://localhost:7000/api/v1/settings
*/
func (h *settingsHandler) put(w http.ResponseWriter, r *http.Request) {
	if ok := rest.DecodeAndValidateJSONPayload(w, r.Body, &settingsRequestBody{}); !ok {
		log.Debug().Msg("validation failed")
		return
	}

	setting, err := h.settings.Get(r.Context())
	if err != nil {
		http.Error(w, rest.HTTPInternalServerError, http.StatusInternalServerError)
		return
	}

	// TODO update the settings with request payload

	// body, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	rest.InternalServerError(w, err)
	// 	return
	// }

	// var settings settings.Setting
	// if err := json.Unmarshal(body, &settings); err != nil {
	// 	rest.BadRequest(w, err)
	// 	return
	// }

	if err := h.settings.Update(r.Context(), setting); err != nil {
		http.Error(w, rest.HTTPInternalServerError, http.StatusInternalServerError)
		return
	}

	rest.JSON(w, setting)
}
