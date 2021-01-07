// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2021/1/7

// Package replacepathregex
package replacepathregex

import (
	"context"
	"fmt"
	"github.com/crochee/proxy/config/dynamic"
	"github.com/crochee/proxy/logger"
	"github.com/crochee/proxy/middlewares/replacepath"
	"net/http"
	"net/url"
	"regexp"
	"strings"
)

// ReplacePathRegex is a middleware used to replace the path of a URL request with a regular expression.
type replacePathRegex struct {
	next        http.Handler
	regexp      *regexp.Regexp
	replacement string
	ctx         context.Context
}

// New creates a new replace path regex middleware.
func New(ctx context.Context, next http.Handler, config dynamic.ReplacePathRegex) (http.Handler, error) {

	exp, err := regexp.Compile(strings.TrimSpace(config.Regex))
	if err != nil {
		return nil, fmt.Errorf("error compiling regular expression %s: %w", config.Regex, err)
	}

	return &replacePathRegex{
		regexp:      exp,
		replacement: strings.TrimSpace(config.Replacement),
		next:        next,
		ctx:         ctx,
	}, nil
}

func (rp *replacePathRegex) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	var currentPath string
	if req.URL.RawPath == "" {
		currentPath = req.URL.Path
	} else {
		currentPath = req.URL.RawPath
	}

	if rp.regexp != nil && len(rp.replacement) > 0 && rp.regexp.MatchString(currentPath) {
		req.Header.Add(replacepath.ReplacedPathHeader, currentPath)

		req.URL.RawPath = rp.regexp.ReplaceAllString(currentPath, rp.replacement)

		// as replacement can introduce escaped characters
		// Path must remain an unescaped version of RawPath
		// Doesn't handle multiple times encoded replacement (`/` => `%2F` => `%252F` => ...)
		var err error
		req.URL.Path, err = url.PathUnescape(req.URL.RawPath)
		if err != nil {
			logger.FromContext(rp.ctx).Error(err.Error())
			http.Error(rw, err.Error(), http.StatusInternalServerError)
			return
		}

		req.RequestURI = req.URL.RequestURI()
	}

	rp.next.ServeHTTP(rw, req)
}
