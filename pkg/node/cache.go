package node

import (
	"net/http"
	"strings"

	lru "github.com/hashicorp/golang-lru"
	uuid "github.com/satori/go.uuid"
)

type ReverseProxyCache interface {
	Get(r *http.Request, uuid uuid.UUID) (http.Handler, bool)
	Set(r *http.Request, uuid uuid.UUID, reverseProxy http.Handler) bool
	Delete(uuid uuid.UUID) bool
}

const defaultProxyCacheSize = 1024

// entry is a tuple containing the http and ws reverse proxies associated to a node
type entry struct {
	http http.Handler
	ws   http.Handler
}

// rpcProxyCache is a concurrency-safe, in-memory cache to reduce proxy-struct re-creations
type rpcProxyCache struct {
	// cache type [uuid]entry cache (key: uuid, value: http/ws reverse proxies entry)
	// note: []byte keys are not supported by golang-lru Cache
	cache *lru.Cache
}

// NewRPCProxyCache creates a new cache for rpc proxies
func NewRPCProxyCache(cacheSize int) *rpcProxyCache {
	cache, _ := lru.New(cacheSize)
	return &rpcProxyCache{cache: cache}
}

// Get returns the reverse-proxy associated to a node's UUID and incoming request type
func (c *rpcProxyCache) Get(r *http.Request, uuid uuid.UUID) (http.Handler, bool) {
	val, ok := c.cache.Get(uuid)
	if !ok {
		return nil, false
	}
	tuple := val.(entry)

	if isWebSocketRequest(r) {
		return tuple.ws, tuple.ws != nil
	}
	return tuple.http, tuple.http != nil
}

// Set persists reverse proxies associated to a node's UUID.
// If a previous entry is found, we need to update the relevant entry reverse proxy based on request (http or ws)
func (c *rpcProxyCache) Set(r *http.Request, uuid uuid.UUID, reverseProxy http.Handler) bool {
	val, ok := c.cache.Get(uuid)
	if !ok {
		if isWebSocketRequest(r) {
			return c.cache.Add(uuid, entry{ws: reverseProxy})
		}
		return c.cache.Add(uuid, entry{http: reverseProxy})
	}
	tuple := val.(entry)
	if isWebSocketRequest(r) {
		tuple.ws = reverseProxy
	} else {
		tuple.http = reverseProxy
	}
	return c.cache.Add(uuid, tuple)
}

// Delete evicts reverse proxies associated to a node's UUID
func (c *rpcProxyCache) Delete(uuid uuid.UUID) bool {
	return c.cache.Remove(uuid)
}

// source: TODO
func isWebSocketRequest(r *http.Request) bool {
	contains := func(key, val string) bool {
		vv := strings.Split(r.Header.Get(key), ",")
		for _, v := range vv {
			if val == strings.ToLower(strings.TrimSpace(v)) {
				return true
			}
		}
		return false
	}
	if !contains("Connection", "upgrade") {
		return false
	}
	if !contains("Upgrade", "websocket") {
		return false
	}
	return true
}
