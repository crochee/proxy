// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/30

package router

import (
	"net/http"
	"net/http/httputil"
	"sync"

	"github.com/crochee/proxy/config"
)

type Redirect struct {
}

func (r Redirect) ServeHTTP(w http.ResponseWriter, request *http.Request) {
	reverseProxy.ServeHTTP(w, request)
}

var reverseProxy = &httputil.ReverseProxy{
	Director: func(r *http.Request) {

	},
	Transport: http.DefaultTransport,
}

type HostMap struct {
	sync.RWMutex
	HostList map[string]Balancer
}

func (h *HostMap) Server(host string) string {
	h.RLock()
	list := h.HostList[host]
	h.RUnlock()
	return list.Next()
}

var (
	singleton *HostMap
	once      sync.Once
)

func NewHostMap(list []*config.ProxyHost) *HostMap {
	once.Do(func() {
		singleton = &HostMap{
			HostList: make(map[string]Balancer, len(list)),
		}
		for _, proxyHost := range list {
			for _, origin := range proxyHost.Origin {
				hostList, ok := singleton.HostList[origin]
				if !ok {
					singleton.HostList[origin] = &Robin{serverList: make([]*server, 0, len(list))}
				}
				for _, target := range proxyHost.Target {
					_ = hostList.UpsertServer(target)
				}
			}
		}
	})
	return singleton
}
