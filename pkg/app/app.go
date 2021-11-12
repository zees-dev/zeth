package app

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"runtime"

	"github.com/ethereum/go-ethereum/params"
	"github.com/pkg/errors"
	"github.com/rs/zerolog/log"
	"github.com/zees-dev/zeth/pkg/datastore"
	"github.com/zees-dev/zeth/pkg/geth/downloader"
	"github.com/zees-dev/zeth/pkg/node"
	"github.com/zees-dev/zeth/pkg/settings"
)

const (
	DefaultPort   = 7000
	DefaultAppDir = "Zeth"
)

type (
	Services struct {
		Settings settings.Settings
		Nodes    node.NodeService
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
			Settings: settings.NewService(store),
			Nodes:    node.NewService(store),
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
func (app *App) Configure() error {
	// download geth binary if it does not exist
	gethBinaryDir := filepath.Join(app.DataDir, downloader.DefaultGethBinaryDir)
	gethBinaryPath := filepath.Join(gethBinaryDir, downloader.GethBinaryFilename(runtime.GOOS, params.Version))
	if _, err := os.Stat(gethBinaryPath); os.IsNotExist(err) {
		log.Info().Msgf("Downloading Geth binary for v%s...", params.Version)
		err := downloader.DownloadGethBinary(gethBinaryDir)
		if err != nil {
			return err
		}
	}

	// start any in-process nodes if they were previously running
	nodes, err := app.Services.Nodes.GetAll(context.TODO())
	if err != nil {
		return errors.Wrap(err, "failed to get all nodes")
	}
	for _, n := range nodes {
		props := n.Properties()
		if props.NodeType == node.TypeGethNodeInProcess && props.Running {
			gipNode := n.(interface{}).(*node.GethInProcessNode)
			log.Info().Msgf("Starting node: %s", props.ID)
			err = gipNode.Start(context.TODO())
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to start node: %s", props.ID))
			}
		}
	}

	return nil
}

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
	defaultNode := node.NewRemoteNode(node.DefaultNodeHTTPRPC, node.DefaultNodeWSRPC)
	defaultNode.Name = node.DefaultNodeName

	_, err := app.Services.Nodes.Create(context.TODO(), defaultNode)
	if err != nil {
		return errors.Wrap(err, "failed to create default node")
	}

	defaultSetting := settings.Setting{
		NodeSettings: settings.NodeSettings{
			DefaultNodeID: defaultNode.ID,
			SupportedNodes: []settings.NodeTypeSetting{
				{
					NodeType: node.TypeGethNodeInProcess,
					Version:  params.Version,
				},
				{
					NodeType: node.TypeRemoteNode,
				},
			},
		},
	}

	return app.Services.Settings.Update(context.TODO(), defaultSetting)
}
