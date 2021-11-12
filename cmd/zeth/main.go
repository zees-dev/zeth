package main

import (
	"os"
	"path/filepath"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/zees-dev/zeth/pkg/app"
	"github.com/zees-dev/zeth/pkg/datastore"
	"github.com/zees-dev/zeth/pkg/httprest"
)

// initializeGlobalLogger initializes the global zerolog logger.
func initializeGlobalLogger() {
	// initialize global zerolog logger
	output := zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.RFC3339,
		NoColor:    false,
		PartsOrder: []string{
			"level",
			"time",
			"caller",
			"message",
		},
	}

	// mutate global zerolog logger
	log.Logger = zerolog.New(output).With().Timestamp().Logger()
	log.Logger = log.With().Caller().Logger()
}

func main() {
	// TODO: CLI flags

	initializeGlobalLogger()

	store, err := datastore.NewBadgerDB(filepath.Join(app.DefaultAppDir, datastore.DefaultStoreDir))
	if err != nil {
		log.Err(err).Msg("failed to create datastore")
		os.Exit(1)
	}
	defer store.Close()

	app := app.NewApp(store, true)

	err = app.Init(store.IsNew())
	if err != nil {
		log.Err(err).Msg("failed to init app")
		os.Exit(1)
	}

	r := httprest.Routing(app)
	if err := httprest.Start(app, r); err != nil {
		log.Err(err).Msg("failed to start http server")
		os.Exit(1)
	}
}
