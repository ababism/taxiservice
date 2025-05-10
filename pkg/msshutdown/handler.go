package msshutdown

import (
	"os"
	"sync"
)

type handler struct {
	stop      chan os.Signal
	mu        sync.Mutex
	callbacks []*Callback
}

func newHandler() *handler {
	return &handler{
		stop: make(chan os.Signal, 1),
	}
}

func (h *handler) add(cb *Callback) {
	h.mu.Lock()
	defer h.mu.Unlock()

	h.callbacks = append(h.callbacks, cb)
}

func (h *handler) wait() {
	<-h.stop
}
