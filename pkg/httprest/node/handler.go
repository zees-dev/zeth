package node

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/zees-dev/zeth/pkg/app"
	"github.com/zees-dev/zeth/pkg/node"
)

type nodesHandler struct {
	nodes          node.NodeService
	nodeRPCMonitor *NodeRPCMonitor
}

func RegisterRoutes(app *app.App, baseRouter *mux.Router) {
	nodeRPCMonitor := NewNodeRPCMonitor()

	h := nodesHandler{
		nodes:          app.Services.Nodes,
		nodeRPCMonitor: nodeRPCMonitor,
	}

	baseRouter.HandleFunc("/nodes", h.getNodes).Methods(http.MethodGet)
	baseRouter.HandleFunc("/nodes", h.createNode).Methods(http.MethodPost)
	baseRouter.HandleFunc("/nodes/{uuid}", h.getNode).Methods(http.MethodGet)
	baseRouter.HandleFunc("/nodes/{uuid}", h.updateNode).Methods(http.MethodPut)
	baseRouter.HandleFunc("/nodes/{uuid}", h.removeNode).Methods(http.MethodDelete)

	baseRouter.HandleFunc("/nodes/rpc/{uuid}", h.rpcNode)
	baseRouter.HandleFunc("/nodes/rpc/{uuid}/sse", h.nodeRPCMonitor.handleSSE).Methods(http.MethodGet)
}
