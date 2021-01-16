// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/17

package replacehost

import (
	"context"
	"fmt"
	"github.com/crochee/proxy/config/dynamic"
	"net/http"
)

const ReplacedHostHeader = "X-Replaced-Host"

// replaceHost is a middleware used to replace host to an URL request.
type replaceHost struct {
	next http.Handler
	host string
	ctx  context.Context
}

// New creates a new handler.
func New(ctx context.Context, next http.Handler, host dynamic.ReplaceHost) (http.Handler, error) {
	if host.Host == "" {
		return nil, fmt.Errorf("host cannot be empty")
	}
	return &replaceHost{
		host: host.Host,
		next: next,
		ctx:  ctx,
	}, nil
}

func (r *replaceHost) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	req.Header.Add(ReplacedHostHeader, req.URL.Host)
	req.URL.Host = r.host
	req.RequestURI = req.URL.RequestURI()
	r.next.ServeHTTP(rw, req)
}
