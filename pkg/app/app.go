package app

import (
	"context"
	"net/http"

	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/zees-dev/zeth/pkg/datastore"
	"github.com/zees-dev/zeth/pkg/defi"
	"github.com/zees-dev/zeth/pkg/defi/amm"
	"github.com/zees-dev/zeth/pkg/node"
	"github.com/zees-dev/zeth/pkg/settings"
)

const (
	DefaultPort   = 7000
	DefaultAppDir = "Zeth"
)

type (
	Services struct {
		Settings             settings.Settings
		Nodes                node.NodeService
		AutomatedMarketMaker defi.AutomatedMarketMaker
	}
	ServeSettings struct {
		Enabled    bool
		FileServer http.Handler
	}
)

type App struct {
	Port     int
	DataDir  string
	IsDev    bool
	Services Services
}

func NewApp(store datastore.Store, isDev bool) *App {
	return &App{
		Port:    DefaultPort,
		DataDir: DefaultAppDir,
		IsDev:   isDev,
		Services: Services{
			Settings:             settings.NewService(store),
			Nodes:                node.NewService(store),
			AutomatedMarketMaker: amm.NewService(store),
		},
	}
}

func (app *App) Init(isNew bool) error {
	// TODO: app configuration
	// configure the application on initialization
	// err := app.Configure()
	// if err != nil {
	// 	return errors.Wrap(err, "failed to configure app")
	// }

	// seed database if it is new
	if isNew {
		err := app.Seed()
		if err != nil {
			return errors.Wrap(err, "failed to seed datastore")
		}
	}

	return nil
}

// Configure will configure the application.
// - download geth binary if they do not exist
// - start up in-process geth nodes (if they were running before)
// func (app *App) Configure() error {
// 	// download geth binary if it does not exist
// 	gethBinaryDir := filepath.Join(app.DataDir, downloader.DefaultGethBinaryDir)
// 	gethBinaryPath := filepath.Join(gethBinaryDir, downloader.GethBinaryFilename(runtime.GOOS, params.Version))
// 	if _, err := os.Stat(gethBinaryPath); os.IsNotExist(err) {
// 		log.Info().Msgf("Downloading Geth binary for v%s...", params.Version)
// 		err := downloader.DownloadGethBinary(gethBinaryDir)
// 		if err != nil {
// 			return err
// 		}
// 	}

// 	return nil
// }

// Seed will seed the datastore with some initial data.
func (app *App) Seed() error {
	log.Info().Msg("Seeding database...")

	// seed default settings
	err := app.seedDefaultSettings()
	if err != nil {
		return err
	}

	return nil
}

func (app *App) seedDefaultSettings() error {
	defaultNode := node.NewNode(node.DefaultNodeHTTPRPC, node.DefaultNodeWSRPC)
	defaultNode.Name = node.DefaultNodeName

	_, err := app.Services.Nodes.Create(context.TODO(), *defaultNode)
	if err != nil {
		return errors.Wrap(err, "failed to create default node")
	}

	// TODO: seed
	defaultSetting := settings.Setting{
		NodeSettings: settings.NodeSettings{
			DefaultNodeID: defaultNode.ID,
			SupportedNodes: []settings.NodeTypeSetting{
				{},
			},
		},
	}

	return app.Services.Settings.Update(context.TODO(), defaultSetting)
}
