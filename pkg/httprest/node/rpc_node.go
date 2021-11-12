package node

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"time"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/yhat/wsutil"
	"github.com/zees-dev/zeth/pkg/httprest/rest"
	"github.com/zees-dev/zeth/pkg/node"
)

var (
	// TODO: improve this cache (use LRU cache if required - https://github.com/hashicorp/golang-lru)
	// nodeReverseProxyCache is a concurrency-safe cache of RPC reverse proxies for nodes.
	nodeReverseProxyCache sync.Map // type is sync.Map[uuid.UUID]http.Handler

	// the set key(s) in the cache are automatically evicted after specified cacheEvictionDuration.
	cacheEvictionDuration = time.Minute * 5
)

/* curl request:
curl -v localhost:7000/api/v1/nodes/rpc/b38bad92-619f-41e4-b01f-36a3de1b3a52 \
	-X POST \
	-H "Content-Type: application/json" \
	-d '{"jsonrpc":"2.0","method":"eth_accounts","params":[],"id":1}'
*/
func (h *nodesHandler) rpcNode(w http.ResponseWriter, r *http.Request) {
	// get id from request parameters
	id := mux.Vars(r)["uuid"]

	uid, err := uuid.FromString(id)
	if err != nil {
		http.Error(w, rest.HTTPBadRequest, http.StatusBadRequest)
		return
	}

	// set x-forwarded-host header to flag this request as proxied
	r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

	// get rpcReverseProxy from cache if possible
	if nodeRPCReverseProxy, ok := nodeReverseProxyCache.Load(uid); ok {
		log.Debug().Msgf("nodeRPCReverseProxy found in cache for node: %s", uid)
		// TODO: log the request/response body/data
		nodeRPCReverseProxy.(http.Handler).ServeHTTP(w, r)
		return
	}

	rpcURLs, err := h.getRPC(r, uid)
	if err != nil {
		http.Error(w, rest.HTTPInternalServerError, http.StatusInternalServerError)
		return
	}

	// TODO: make this into a middleware
	// note: this is only applicable to HTTP RPC requests
	bodyBytes, _ := ioutil.ReadAll(r.Body)
	r.Body.Close() //  must close
	r.Body = ioutil.NopCloser(bytes.NewBuffer(bodyBytes))
	log.Info().Msgf(
		"proxying rpc request:\n\tnode: %s\n\tbody: %s\n\tRPCs: %v",
		uid,
		string(bodyBytes),
		rpcURLs,
	)

	// serve the reverse proxy
	proxy, err := getNodeReverseProxy(rpcURLs, r)
	if err != nil {
		http.Error(w, rest.HTTPInternalServerError, http.StatusInternalServerError)
	}

	// store proxy in cache for quicker subsequent lookups; automatically evict cache after x minutes
	nodeReverseProxyCache.Store(uid, proxy)
	go func() {
		<-time.After(cacheEvictionDuration)
		log.Debug().Msgf("evicting nodeReverseProxyCache for node: %s", uid)
		nodeReverseProxyCache.Delete(uid)
	}()

	// note: ServeHttp is a blocking call
	proxy.ServeHTTP(w, r)
}

// getRPCURL returns the RPC URL for the given node.
func (h *nodesHandler) getRPC(r *http.Request, uuid uuid.UUID) (node.RPC, error) {
	var httpRPCURL, websocketRPCURL string

	if n, err := h.nodes.Get(r.Context(), uuid); err == nil {
		switch n.Properties().NodeType {
		// case node.TypeGethNodeInProcess:
		// 	// check request for `Upgrade` header to identify whether its http request or websocket
		// 	gipNode := n.(*node.GethInProcessNode)
		// 	var scheme, host string
		// 	var port int
		// 	if wsutil.IsWebSocketRequest(r) {
		// 		scheme, host, port = "http", gipNode.GethConfig.Node.HTTPHost, gipNode.GethConfig.Node.HTTPPort
		// 	} else {
		// 		scheme, host, port = "ws", gipNode.GethConfig.Node.WSHost, gipNode.GethConfig.Node.WSPort
		// 	}
		// 	rpcURL = fmt.Sprintf("%s://%s:%d", scheme, host, port)
		case node.TypeRemoteNode:
			remoteNode := n.(*node.RemoteNode)
			if _, err := url.Parse(remoteNode.RPC.HTTP); err != nil {
				return node.RPC{}, err
			}
			httpRPCURL = remoteNode.RPC.HTTP
			if _, err := url.Parse(remoteNode.RPC.WS); err == nil {
				websocketRPCURL = remoteNode.RPC.WS
			}
		default:
			return node.RPC{}, fmt.Errorf("unsupported node type: %v", n.Properties().NodeType)
		}
	}

	return node.RPC{HTTP: httpRPCURL, WS: websocketRPCURL}, nil
}

// getNodeReverseProxy gets the relevant http or websocket reverse-proxy for the calling request.
func getNodeReverseProxy(rpc node.RPC, r *http.Request) (http.Handler, error) {
	var proxy http.Handler

	// Override http/websocket reverse proxies to support https/wss - https://stackoverflow.com/a/53007606/10813908
	if wsutil.IsWebSocketRequest(r) {
		url, err := url.Parse(rpc.WS)
		if err != nil {
			return nil, err
		}
		p := wsutil.NewSingleHostReverseProxy(url) // ws(s) reverse proxy
		d := p.Director
		p.Director = func(r *http.Request) {
			d(r)              // call default director
			r.Host = url.Host // set Host header as expected by target
			r.URL.Scheme = url.Scheme
			r.URL.Path = url.Path
		}
		proxy = p
	} else {
		url, err := url.Parse(rpc.HTTP)
		if err != nil {
			return nil, err
		}

		p := httputil.NewSingleHostReverseProxy(url) // http(s) proxy reverse proxy
		d := p.Director
		p.Director = func(r *http.Request) {
			d(r)              // call default director
			r.Host = url.Host // set Host header as expected by target
			r.URL.Scheme = url.Scheme
			r.URL.Path = url.Path
		}
		proxy = p
	}

	// strip the RPC request's prefix from the request to forward it to the target node
	proxy = http.StripPrefix(r.URL.Path, proxy)

	return proxy, nil
}
