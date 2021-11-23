package node

import (
	"fmt"
	"net/http"
	"sync"

	"github.com/gorilla/mux"
	uuid "github.com/satori/go.uuid"
)

// source: https://gist.github.com/rikonor/e53a33c27ed64861c91a095a59f0aa44

type UnsubscribeFunc func() error

type Subscriber interface {
	Subscribe(c chan []byte) (UnsubscribeFunc, error)
}

type Publisher interface {
	Publish(b []byte) error
}

type NotificationCenter struct {
	subscribers   map[chan []byte]struct{}
	subscribersMu *sync.Mutex
}

func NewNotificationCenter() *NotificationCenter {
	return &NotificationCenter{
		subscribers:   map[chan []byte]struct{}{},
		subscribersMu: &sync.Mutex{},
	}
}

func (nc *NotificationCenter) Subscribe(c chan []byte) (UnsubscribeFunc, error) {
	nc.subscribersMu.Lock()
	nc.subscribers[c] = struct{}{}
	nc.subscribersMu.Unlock()

	unsubscribeFn := func() error {
		nc.subscribersMu.Lock()
		delete(nc.subscribers, c)
		nc.subscribersMu.Unlock()

		return nil
	}

	return unsubscribeFn, nil
}

func (nc *NotificationCenter) Publish(b []byte) error {
	nc.subscribersMu.Lock()
	defer nc.subscribersMu.Unlock()

	for c := range nc.subscribers {
		select {
		case c <- b:
		default:
		}
	}

	return nil
}

type PubSub interface {
	Publisher
	Subscriber
}

type NodeRPCMonitor struct {
	rpcs map[uuid.UUID]PubSub
}

func NewNodeRPCMonitor() *NodeRPCMonitor {
	return &NodeRPCMonitor{rpcs: make(map[uuid.UUID]PubSub)}
}

/*
curl request:
curl -v http://localhost:7000/api/v1/nodes/rpc/3475ce0e-0124-4e8b-a661-b4b6e22cdf34/sse
*/
func (n *NodeRPCMonitor) handleSSE(w http.ResponseWriter, r *http.Request) {
	id := mux.Vars(r)["uuid"]
	nodeUUID, err := uuid.FromString(id)
	if err != nil {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	subscriber, ok := n.rpcs[nodeUUID]
	if !ok {
		http.Error(w, fmt.Sprintf("node subscriber not found: %s", nodeUUID), http.StatusNotFound)
		return
	}

	// Subscribe
	c := make(chan []byte)
	unsubscribeFn, err := subscriber.Subscribe(c)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	// Signal SSE Support
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Access-Control-Allow-Origin", "*")

	for {
		select {
		case <-r.Context().Done():
			if err := unsubscribeFn(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
			return
		default:
			b := <-c
			fmt.Fprintf(w, "data: %s\n\n", b)
			// fmt.Fprintf(w, string(b))
			w.(http.Flusher).Flush()
		}
	}
}
