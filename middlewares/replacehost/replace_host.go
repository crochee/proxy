// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/17

package replacehost

import (
	"context"
	"fmt"
	"net/http"
	"net/url"

	"github.com/crochee/proxy/config/dynamic"
)

const ReplacedHostHeader = "X-Replaced-Host"

// replaceHost is a middleware used to replace host to an URL request.
type replaceHost struct {
	next   http.Handler
	scheme string
	host   string
	ctx    context.Context
}

// New creates a new handler.
func New(ctx context.Context, next http.Handler, host dynamic.ReplaceHost) (http.Handler, error) {
	if host.Host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}
	if host.Scheme == "" {
		host.Scheme = "http"
	}
	return &replaceHost{
		scheme: host.Scheme,
		host:   host.Host,
		next:   next,
		ctx:    ctx,
	}, nil
}

func (r *replaceHost) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	u := req.URL
	if req.RequestURI != "" {
		parsedURL, err := url.ParseRequestURI(req.RequestURI)
		if err == nil {
			u = parsedURL
		}
	}
	u.Scheme = r.scheme
	u.Host = r.host

	req.URL = u
	req.RequestURI = req.URL.RequestURI()
	req.Header.Add(ReplacedHostHeader, req.URL.Host)

	r.next.ServeHTTP(rw, req)
}
