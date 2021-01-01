// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/31

package service

import (
	"context"
	"errors"
	"io"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"time"

	"github.com/crochee/proxy/logger"
)

// StatusClientClosedRequest non-standard HTTP status code for client disconnection.
const StatusClientClosedRequest = 499

// StatusClientClosedRequestText non-standard HTTP status for client disconnection.
const StatusClientClosedRequestText = "Client Closed Request"

func buildProxy(flushInterval time.Duration, roundTripper http.RoundTripper) (http.Handler, error) {
	proxy := &httputil.ReverseProxy{
		Director:      Director,
		Transport:     roundTripper,
		FlushInterval: flushInterval,
		BufferPool:    newBufferPool(),
		ErrorHandler:  ErrorHandler,
	}

	return proxy, nil
}

func statusText(statusCode int) string {
	if statusCode == StatusClientClosedRequest {
		return StatusClientClosedRequestText
	}
	return http.StatusText(statusCode)
}

func ErrorHandler(w http.ResponseWriter, request *http.Request, err error) {
	statusCode := http.StatusInternalServerError

	switch {
	case errors.Is(err, io.EOF):
		statusCode = http.StatusBadGateway
	case errors.Is(err, context.Canceled):
		statusCode = StatusClientClosedRequest
	default:
		var netErr net.Error
		if errors.As(err, &netErr) {
			if netErr.Timeout() {
				statusCode = http.StatusGatewayTimeout
			} else {
				statusCode = http.StatusBadGateway
			}
		}
	}

	logger.Debugf("'%d %s' caused by: %v", statusCode, statusText(statusCode), err)
	w.WriteHeader(statusCode)
	if _, err = w.Write([]byte(statusText(statusCode))); err != nil {
		logger.Errorf("Error while writing status code", err)
	}
}

func Director(request *http.Request) {
	u := request.URL
	if request.RequestURI != "" {
		parsedURL, err := url.ParseRequestURI(request.RequestURI)
		if err == nil {
			u = parsedURL
		}
	}
	request.URL.Path = u.Path
	request.URL.RawPath = u.RawPath
	request.URL.RawQuery = u.RawQuery
	request.RequestURI = "" // Outgoing request should not have RequestURI

	if _, ok := request.Header["User-Agent"]; !ok {
		request.Header.Set("User-Agent", "")
	}
}

var DefaultTransport = &http.Transport{
	Proxy: http.ProxyFromEnvironment,
	DialContext: (&net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
		DualStack: true,
	}).DialContext,
	ForceAttemptHTTP2:     true,
	MaxIdleConns:          100,
	IdleConnTimeout:       90 * time.Second,
	TLSHandshakeTimeout:   10 * time.Second,
	ExpectContinueTimeout: 1 * time.Second,
}
