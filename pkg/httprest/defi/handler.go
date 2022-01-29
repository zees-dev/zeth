package defi

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zees-dev/zeth/pkg/app"
	"github.com/zees-dev/zeth/pkg/defi"
)

type handler struct {
	ammSvc defi.AutomatedMarketMaker
}

func RegisterRoutes(app *app.App, baseRouter *mux.Router) {
	h := handler{
		ammSvc: app.Services.AutomatedMarketMaker,
	}

	baseRouter.HandleFunc("/defi/amm", h.getAllAMMs).Methods(http.MethodGet)
	baseRouter.HandleFunc("/defi/amm", h.createAMM).Methods(http.MethodPost)
	baseRouter.HandleFunc("/defi/amm/{uuid}", h.getAMM).Methods(http.MethodGet)
	baseRouter.HandleFunc("/defi/amm/{uuid}", h.removeAMM).Methods(http.MethodDelete)
	baseRouter.HandleFunc("/defi/{chainID}/amm", h.getAMMsByChainID).Methods(http.MethodGet)
}
