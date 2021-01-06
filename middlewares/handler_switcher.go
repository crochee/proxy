// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/6

package middlewares

import (
	"net/http"
	"sync"
)

// HTTPHandlerSwitcher allows hot switching of http.ServeMux.
type HTTPHandlerSwitcher struct {
	handler http.Handler
	lock    sync.RWMutex
}

// NewHandlerSwitcher builds a new instance of HTTPHandlerSwitcher.
func NewHandlerSwitcher(newHandler http.Handler) (hs *HTTPHandlerSwitcher) {
	return &HTTPHandlerSwitcher{
		handler: newHandler,
	}
}

func (h *HTTPHandlerSwitcher) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	handlerBackup := h.GetHandler()
	handlerBackup.ServeHTTP(rw, req)
}

// GetHandler returns the current http.ServeMux.
func (h *HTTPHandlerSwitcher) GetHandler() http.Handler {
	h.lock.RLock()
	defer h.lock.RUnlock()
	return h.handler
}

// UpdateHandler safely updates the current http.ServeMux with a new one.
func (h *HTTPHandlerSwitcher) UpdateHandler(newHandler http.Handler) {
	h.lock.Lock()
	defer h.lock.Unlock()
	h.handler = newHandler
}
