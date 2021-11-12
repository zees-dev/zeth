package settings

import (
	"context"

	uuid "github.com/satori/go.uuid"
	"github.com/zees-dev/zeth/pkg/node"
)

type (
	// NodeTypeSetting represents the settings for a node type
	NodeTypeSetting struct {
		NodeType node.NodeType `json:"nodeType"`
		Version  string        `json:"version"`
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
