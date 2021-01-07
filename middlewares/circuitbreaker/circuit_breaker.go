// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2021/1/7

// Package circuitbreaker
package circuitbreaker

import (
	"context"
	"net/http"

	"github.com/afex/hystrix-go/hystrix"
	"github.com/crochee/proxy/config/dynamic"
)

type circuitBreaker struct {
	name string
	next http.Handler
	ctx  context.Context
}

// New creates a new circuit breaker middleware.
func New(ctx context.Context, next http.Handler, breaker dynamic.CircuitBreaker, name string) (http.Handler, error) {
	hystrix.ConfigureCommand(name, hystrix.CommandConfig{
		Timeout:                0,
		MaxConcurrentRequests:  0,
		RequestVolumeThreshold: 0,
		SleepWindow:            0,
		ErrorPercentThreshold:  0,
	})
	return &circuitBreaker{
		name: name,
		next: next,
		ctx:  ctx,
	}, nil
}

// todo 熔断器待实现
func (c *circuitBreaker) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c.next.ServeHTTP(rw, req)
}
