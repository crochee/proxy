// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2021/1/7

// Package replacepath
package replacepath

import (
	"context"
	"net/http"
	"net/url"

	"github.com/crochee/proxy/config/dynamic"
	"github.com/crochee/proxy/logger"
)

const (
	// ReplacedPathHeader is the default header to set the old path to.
	ReplacedPathHeader = "X-Replaced-Path"
)

// ReplacePath is a middleware used to replace the path of a URL request.
type replacePath struct {
	next http.Handler
	path string
	ctx  context.Context
}

// New creates a new replace path middleware.
func New(ctx context.Context, next http.Handler, path dynamic.ReplacePath) (http.Handler, error) {
	return &replacePath{
		next: next,
		path: path.Path,
		ctx:  ctx,
	}, nil
}

func (r *replacePath) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	if req.URL.RawPath == "" {
		req.Header.Add(ReplacedPathHeader, req.URL.Path)
	} else {
		req.Header.Add(ReplacedPathHeader, req.URL.RawPath)
	}

	req.URL.RawPath = r.path

	var err error
	req.URL.Path, err = url.PathUnescape(req.URL.RawPath)
	if err != nil {
		logger.FromContext(r.ctx).Error(err.Error())
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	req.RequestURI = req.URL.RequestURI()

	r.next.ServeHTTP(rw, req)
}
