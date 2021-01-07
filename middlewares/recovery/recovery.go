// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2021/1/7

// Package recovery
package recovery

import (
	"context"
	"net/http"
	"runtime/debug"

	"github.com/crochee/proxy/logger"
)

type recovery struct {
	next http.Handler
	ctx  context.Context
}

// New creates recovery middleware.
func New(ctx context.Context, next http.Handler) (http.Handler, error) {
	return &recovery{
		next: next,
		ctx:  ctx,
	}, nil
}

func (re *recovery) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	defer func(ctx context.Context, rw http.ResponseWriter, r *http.Request) {
		if err := recover(); err != nil {
			log := logger.FromContext(ctx)
			if err == http.ErrAbortHandler {
				log.Debugf("Request has been aborted [%s - %s]: %v", r.RemoteAddr, r.URL, err)
				return
			}

			log.Errorf("Recovered from panic in HTTP handler [%s - %s]: %+v", r.RemoteAddr, r.URL, err)

			log.Errorf("Stack: %s", debug.Stack())

			http.Error(rw, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		}
	}(re.ctx, rw, req)
	re.next.ServeHTTP(rw, req)
}
