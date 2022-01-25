package settings

import (
	"context"

	uuid "github.com/satori/go.uuid"
)

type (
	// TODO: modify/remove
	NodeTypeSetting struct {
		Version string `json:"version"`
	}
	NodeSettings struct {
		SupportedNodes []NodeTypeSetting `json:"supportedNodes"`
		DefaultNodeID  uuid.UUID         `json:"defaultNodeID"`
	}
	Setting struct {
		NodeSettings NodeSettings `json:"nodeSettings"`
	}
	Settings interface {
		Get(ctx context.Context) (Setting, error)
		Update(ctx context.Context, setting Setting) error
	}
)
