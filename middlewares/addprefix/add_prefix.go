// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2021/1/7

// Package addprefix
package addprefix

import (
	"context"
	"fmt"
	"net/http"

	"github.com/crochee/proxy/config/dynamic"
	"github.com/crochee/proxy/logger"
)

// AddPrefix is a middleware used to add prefix to an URL request.
type addPrefix struct {
	next   http.Handler
	prefix string
	ctx    context.Context
}

// New creates a new handler.
func New(ctx context.Context, next http.Handler, prefix dynamic.AddPrefix) (http.Handler, error) {
	if prefix.Prefix == "" {
		return nil, fmt.Errorf("prefix cannot be empty")
	}
	return &addPrefix{
		prefix: prefix.Prefix,
		next:   next,
		ctx:    ctx,
	}, nil
}

func (a *addPrefix) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	log := logger.FromContext(a.ctx)

	oldURLPath := req.URL.Path
	req.URL.Path = ensureLeadingSlash(a.prefix + req.URL.Path)
	log.Debugf("URL.Path is now %s (was %s).", req.URL.Path, oldURLPath)

	if req.URL.RawPath != "" {
		oldURLRawPath := req.URL.RawPath
		req.URL.RawPath = ensureLeadingSlash(a.prefix + req.URL.RawPath)
		log.Debugf("URL.RawPath is now %s (was %s).", req.URL.RawPath, oldURLRawPath)
	}
	req.RequestURI = req.URL.RequestURI()

	a.next.ServeHTTP(rw, req)
}

func ensureLeadingSlash(str string) string {
	if str == "" {
		return str
	}

	if str[0] == '/' {
		return str
	}

	return "/" + str
}
