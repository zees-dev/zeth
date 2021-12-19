package node

import (
	"bytes"
	"compress/gzip"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

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
	nc := NewNotificationCenter()
	h.nodeRPCMonitor.rpcs[uid] = nc
	proxy, err := createNodeReverseProxy(rpcURLs, h.nodeRPCMonitor.rpcs[uid], r)
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

// createNodeReverseProxy gets the relevant http or websocket reverse-proxy for the calling request.
func createNodeReverseProxy(rpc node.RPC, publisher Publisher, r *http.Request) (http.Handler, error) {
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
		p.Transport = rpcRoundTripper{rpcURL: rpc.HTTP, publisher: publisher}
		proxy = p
	}

	// strip the RPC request's prefix from the request to forward it to the target node
	proxy = http.StripPrefix(r.URL.Path, proxy)

	return proxy, nil
}

type callType int

const (
	_ callType = iota
	request
	response
)

type RPCEvent struct {
	ID      string `json:"id"`
	URI     string `json:"uri"`
	RPCURL  string `json:"rpcURL"`
	Request struct {
		Headers string `json:"headers"`
		Body    string `json:"body"`
	} `json:"request"`
	Response struct {
		Headers string `json:"headers"`
		Body    string `json:"body"`
		// Body       map[string]interface{} `json:"body"`
		StatusCode int `json:"statusCode"`
	} `json:"response"`
	Duration int64 `json:"duration,omitempty"` // duration in milliseconds
}

func NewRPCEvent(rpcURL string) *RPCEvent {
	return &RPCEvent{ID: uuid.NewV4().String(), RPCURL: rpcURL}
}

func (ev *RPCEvent) ParseRequest(r *http.Request) {
	reqHeadersBytes, _ := json.Marshal(r.Header)

	var buf bytes.Buffer
	io.Copy(&buf, r.Body)
	// r.Body.Close() //  must close

	// set event request properties
	ev.URI = r.RequestURI
	ev.Request.Headers = string(reqHeadersBytes)
	ev.Request.Body = buf.String()

	r.Body = ioutil.NopCloser(bytes.NewBuffer(buf.Bytes()))
}

func (ev *RPCEvent) ParseResponse(res *http.Response) {
	resHeadersBytes, _ := json.Marshal(res.Header)

	var buf bytes.Buffer
	io.Copy(&buf, res.Body)

	// set event response properties
	ev.Response.Headers = string(resHeadersBytes)
	ev.Response.StatusCode = res.StatusCode

	var responseBody string
	// check res header for content-encoding; if gzip, decompress the response body
	if res.Header.Get("Content-Encoding") == "gzip" {
		gzipReader, err := gzip.NewReader(bytes.NewReader(buf.Bytes()))
		if err != nil {
			log.Debug().Err(err).Msg("failed to decompress response body")
		}
		// read decompressed response body
		var decompressedBody bytes.Buffer
		io.Copy(&decompressedBody, gzipReader)
		responseBody = decompressedBody.String()
	} else {
		responseBody = buf.String()
	}
	ev.Response.Body = responseBody

	res.Body = ioutil.NopCloser(bytes.NewBuffer(buf.Bytes()))
}

func (ev RPCEvent) Bytes() []byte {
	b := new(bytes.Buffer)
	json.NewEncoder(b).Encode(ev)
	return b.Bytes()
}

// rpcRoundTripper satisfies the http.RoundTripper interface
type rpcRoundTripper struct {
	rpcURL    string
	publisher Publisher
}

// RoundTrip satisfies the http.RoundTripper interface
// During roundtrip, we modify request header, publish RPC'ed request and response to subscribers
func (rt rpcRoundTripper) RoundTrip(r *http.Request) (*http.Response, error) {
	reqStartTime := time.Now()

	// modify request headers
	r.Header.Set("Host", r.Host)
	// r.Header.Set("X-Forwarded-Host", r.Header.Get("Host"))

	event := NewRPCEvent(rt.rpcURL)

	// parse request, log request, publish to subscriber
	event.ParseRequest(r)
	log.Info().Msg(fmt.Sprintf(
		"proxying rpc request:\n\treq: %s\n\trpc: %s\n\theaders: %s\n\tbody: %s",
		event.URI,
		event.RPCURL,
		event.Request.Headers,
		event.Request.Body,
	))
	rt.publisher.Publish(event.Bytes())

	// perform roundtrip against actual underlying rpc endpoint
	res, err := http.DefaultTransport.RoundTrip(r)
	if err != nil {
		return nil, err
	}

	// parse response, calc duration, log response, publish to subscriber
	event.ParseResponse(res)
	event.Duration = time.Since(reqStartTime).Milliseconds()
	log.Info().Msg(fmt.Sprintf(
		"proxied rpc response:\n\trpc: %s\n\theaders: %s\n\tstatus: %d\n\tbody: %s\n\tduration: %d",
		event.RPCURL,
		event.Response.Headers,
		event.Response.StatusCode,
		event.Response.Body,
		event.Duration,
	))
	rt.publisher.Publish(event.Bytes())

	return res, nil
}
