package defi

import (
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/zees-dev/zeth/pkg/defi/amm"
	"github.com/zees-dev/zeth/pkg/httprest/rest"
)

/* curl request:
curl \
	-H "Content-Type: application/json" \
	http://localhost:7000/api/v1/defi/amm
*/
func (h *handler) getAllAMMs(w http.ResponseWriter, r *http.Request) {
	amms, err := h.ammSvc.GetAll(r.Context())
	if err != nil {
		http.Error(w, rest.HTTPInternalServerError, http.StatusInternalServerError)
		return
	}

	rest.JSON(w, amms)
}

/* curl request:
curl \
	-H "Content-Type: application/json" \
	http://localhost:7000/api/v1/defi/amm/1
*/
func (h *handler) getAMMsByChainID(w http.ResponseWriter, r *http.Request) {
	// get id from request parameters
	id := mux.Vars(r)["chainID"]

	chainID, err := strconv.Atoi(id)
	if err != nil {
		http.Error(w, rest.HTTPBadRequest, http.StatusBadRequest)
		return
	}

	amms, err := h.ammSvc.GetByChainID(r.Context(), chainID)
	if err != nil {
		http.Error(w, rest.HTTPNotFound, http.StatusNotFound)
		return
	}

	rest.JSON(w, amms)
}

type registerAMMPayload struct {
	Name           string `json:"name"`
	URL            string `json:"url"`
	ChainID        uint   `json:"chainId"`
	RouterAddress  string `json:"routerAddress"`
	FactoryAddress string `json:"factoryAddress"`
}

func (payload *registerAMMPayload) Validate() url.Values {
	errs := url.Values{}

	if payload.Name == "" {
		errs.Add("name", "name is required")
	}

	if payload.ChainID == 0 {
		errs.Add("chainId", "chainId is required")
	}

	if payload.RouterAddress == "" {
		errs.Add("routerAddress", "routerAddress is required")
	} else if !strings.HasPrefix(payload.RouterAddress, "0x") {
		errs.Add("routerAddress", "routerAddress must start with 0x")
	}

	if payload.FactoryAddress == "" {
		errs.Add("factoryAddress", "factoryAddress is required")
	} else if !strings.HasPrefix(payload.FactoryAddress, "0x") {
		errs.Add("factoryAddress", "factoryAddress must start with 0x")
	}

	return errs
}

// registerNode registers an externally running node with zeth
/* curl request:
curl -X POST \
	-H "Content-Type: application/json" \
	-d '{"name":"test","url":"http://localhost:7000","chainId":1,"routerAddress":"0x1234567890123456789012345678901234567890","factoryAddress":"0x1234567890123456789012345678901234567890"}' \
	http://localhost:7000/api/v1/defi/amm
*/
func (h *handler) createAMM(w http.ResponseWriter, r *http.Request) {
	payload := registerAMMPayload{}
	if ok := rest.DecodeAndValidateJSONPayload(w, r.Body, &payload); !ok {
		log.Debug().Msg("validation failed")
		return
	}

	amm := amm.AMM{
		ID:             uuid.NewV4(),
		Name:           payload.Name,
		URL:            payload.URL,
		ChainID:        payload.ChainID,
		RouterAddress:  payload.RouterAddress,
		FactoryAddress: payload.FactoryAddress,
	}

	amm, err := h.ammSvc.Create(r.Context(), amm)
	if err != nil {
		http.Error(w, rest.HTTPInternalServerError, http.StatusInternalServerError)
		return
	}

	rest.JSON(w, amm)
}

/* curl request:
curl \
	-H "Content-Type: application/json" \
	http://localhost:7000/api/v1/defi/amm/21fa9b25-840a-40e9-acda-fad525615e58
*/
func (h *handler) getAMM(w http.ResponseWriter, r *http.Request) {
	// get id from request parameters
	id := mux.Vars(r)["uuid"]

	uid, err := uuid.FromString(id)
	if err != nil {
		http.Error(w, rest.HTTPBadRequest, http.StatusBadRequest)
		return
	}

	amm, err := h.ammSvc.Get(r.Context(), uid)
	if err != nil {
		http.Error(w, rest.HTTPNotFound, http.StatusNotFound)
		return
	}

	rest.JSON(w, amm)
}

/* curl request:
curl -X DELETE \
	http://localhost:7000/api/v1/defi/amm/21fa9b25-840a-40e9-acda-fad525615e58
*/
func (h *handler) removeAMM(w http.ResponseWriter, r *http.Request) {
	// get id from request parameters
	id := mux.Vars(r)["uuid"]

	uid, err := uuid.FromString(id)
	if err != nil {
		http.Error(w, rest.HTTPBadRequest, http.StatusBadRequest)
		return
	}

	_, err = h.ammSvc.Get(r.Context(), uid)
	if err != nil {
		http.Error(w, rest.HTTPNotFound, http.StatusNotFound)
		return
	}

	err = h.ammSvc.Delete(r.Context(), uid)
	if err != nil {
		http.Error(w, rest.HTTPInternalServerError, http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
