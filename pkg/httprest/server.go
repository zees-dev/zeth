package httprest

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/zees-dev/zeth/pkg/app"
	"github.com/zees-dev/zeth/pkg/httprest/node"
	"github.com/zees-dev/zeth/pkg/httprest/settings"
)

// Routing will return the router for the server. It will load all routes and sub-routes.
func Routing(app *app.App) *mux.Router {
	log.Info().Msgf("loading routes...")
	r := mux.NewRouter()

	// register server health check route
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
	})

	// sub-route all API calls under versioned API path
	apiRouter := r.PathPrefix("/api/v1").Subrouter()

	// use some middleware
	// h.baseRouter.Use()

	settings.RegisterRoutes(app, apiRouter)
	node.RegisterRoutes(app, apiRouter)

	// serve UI on root URL if UIDir provided
	if app.ServeSettings.Enabled {
		r.PathPrefix("/").Handler(app.ServeSettings.FileServer)
	}

	return r
}

// Start the server in a blocking manner.
func Start(app *app.App, r http.Handler) error {
	server := http.Server{
		Addr:         fmt.Sprintf(":%d", app.Port),
		Handler:      r,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}

	serverErrors := make(chan error, 1)

	// Start the service listening for requests.
	go func() {
		log.Info().Msgf("starting zeth on port %d...", app.Port)
		serverErrors <- server.ListenAndServe()
	}()

	// Shutdown
	shutdown := make(chan os.Signal, 1)
	signal.Notify(shutdown, os.Interrupt, syscall.SIGTERM)

	select {
	// Make a channel to listen for errors coming from the listener. Use a
	// buffered channel so the goroutine can exit if we don'usertransport collect this error.
	case err := <-serverErrors:
		return errors.Wrap(err, "server error")

	// Blocking main and waiting for shutdown.
	case sig := <-shutdown:
		log.Info().Msgf("start shutdown: %v", sig)

		// Give outstanding requests a deadline for completion.
		ctx, cancel := context.WithTimeout(context.Background(), 10)
		defer cancel()

		// Asking listener to shutdown and load shed.
		err := server.Shutdown(ctx)
		if err != nil {
			log.Info().Msgf("graceful shutdown did not complete in %v : %v", 10, err)
			err = server.Close()
			return err
		}

		switch {
		case sig == syscall.SIGSTOP:
			return errors.New("integrity issue caused shutdown")
		case err != nil:
			return errors.Wrap(err, "could not stop server gracefully")
		}
	}

	return nil
}
