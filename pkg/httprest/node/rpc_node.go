package node

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"
	uuid "github.com/satori/go.uuid"
	"github.com/yhat/wsutil"
	"github.com/zees-dev/zeth/pkg/httprest/rest"
	"github.com/zees-dev/zeth/pkg/node"
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

	// get rpcReverseProxy from cache if possible
	if nodeRPCReverseProxy, ok := h.nodes.ReverseProxyCache().Get(r, uid); ok {
		log.Debug().Msgf("nodeRPCReverseProxy found in cache for node: %s", uid)
		// TODO: log the request/response body/data
		nodeRPCReverseProxy.ServeHTTP(w, r)
		return
	}

	rpcURLs, err := h.getRPC(r, uid)
	if err != nil {
		http.Error(w, rest.HTTPInternalServerError, http.StatusInternalServerError)
		return
	}

	// serve the reverse proxy
	proxy, err := getNodeReverseProxy(rpcURLs, r)
	if err != nil {
		http.Error(w, rest.HTTPInternalServerError, http.StatusInternalServerError)
	}

	// store proxy in cache for quicker subsequent lookups; automatically evict cache after x minutes
	h.nodes.ReverseProxyCache().Set(r, uid, proxy)

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
			modifyRequest(r, rpc.WS)
		}
		// p.Transport = rpcRoundTripper{rpc.WS}
		proxy = p
	} else {
		url, err := url.Parse(rpc.HTTP)
		if err != nil {
			return nil, err
		}

		p := httputil.NewSingleHostReverseProxy(url) // http(s) proxy reverse proxy
		d := p.Director
		p.Director = func(r *http.Request) {
			d(r) // call default director
			// r.URL.Host = url.Host
			// r.URL.Scheme = url.Scheme
			r.URL.Path = url.Path
			r.Host = url.Host // set Host header as expected by target
		}
		p.Transport = rpcRoundTripper{rpc.HTTP}
		proxy = p
	}

	// strip the RPC request's prefix from the request to forward it to the target node
	proxy = http.StripPrefix(r.URL.Path, proxy)

	return proxy, nil
}

// rpcRoundTripper satisfies the http.RoundTripper interface
type rpcRoundTripper struct {
	rpcURL string
}

// RoundTrip satisfies the http.RoundTripper interface
// This allows us to simply log and/or modify the request end-to-end
func (rt rpcRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	log.Printf("request received. url=%s", r.URL)
	defer log.Printf("request complete. url=%s", r.URL)

	modifyRequest(r, rt.rpcURL)

	res, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	modifyResponse(res, rt.rpcURL)

	return res, nil
}

// modifyRequest set x-forwarded-host header to flag this request as proxied
func modifyRequest(r *http.Request, rpcURL string) {
	r.Header.Set("Host", r.Host)
	// r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

	// log request
	reqHeadersBytes, _ := json.Marshal(r.Header)

	var buf bytes.Buffer
	io.Copy(&buf, r.Body)
	// r.Body.Close() //  must close
	log.Info().Msgf(
		"proxying rpc request:\n\treq: %s\n\trpc: %s\n\theaders: %s\n\tbody: %s",
		r.RequestURI,
		rpcURL,
		string(reqHeadersBytes),
		buf.String(),
	)
	r.Body = ioutil.NopCloser(bytes.NewBuffer(buf.Bytes()))
}

// modifyResponse logs the response
// source: https://daryl-ng.medium.com/why-you-should-avoid-ioutil-readall-in-go-e6be4de180f8
func modifyResponse(res *http.Response, rpcURL string) error {
	// log response
	resHeadersBytes, _ := json.Marshal(res.Header)

	var buf bytes.Buffer
	_, err := io.Copy(&buf, res.Body)
	if err != nil {
		return err
	}

	log.Info().Msgf(
		"proxied rpc response:\n\trpc: %s\n\theaders: %s\n\tstatus: %d\n\tbody: %s",
		rpcURL,
		string(resHeadersBytes),
		res.StatusCode,
		buf.String(),
	)

	res.Body = ioutil.NopCloser(bytes.NewBuffer(buf.Bytes()))

	return nil
}
