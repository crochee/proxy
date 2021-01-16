// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/1

package http

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"github.com/crochee/proxy/middlewares/replacehost"
	"net"
	"net/http"
	"time"

	"github.com/crochee/proxy/config"
	"github.com/crochee/proxy/logger"
	"github.com/crochee/proxy/middlewares"
	"github.com/crochee/proxy/middlewares/forwardedheaders"
	"github.com/crochee/proxy/server/service"
	tls2 "github.com/crochee/proxy/tls"
)

// EntryPoint is the http server.
type EntryPoint struct {
	listener     net.Listener
	switcher     *middlewares.HTTPHandlerSwitcher
	server       *http.Server
	ctx          context.Context
	serverConfig *config.EntryPoint
}

// NewEntryPoint creates a new EntryPoint.
func NewEntryPoint(ctx context.Context, configuration *config.EntryPoint) (*EntryPoint, error) {
	listener, err := net.Listen("tcp", configuration.GetPort())
	if err != nil {
		return nil, fmt.Errorf("error opening listener: %w", err)
	}
	var rt http.RoundTripper
	if rt, err = service.CreateRoundTripper(config.Cfg.Transport); err != nil {
		return nil, err
	}
	var proxyRoute http.Handler
	if proxyRoute, err = service.BuildProxy(30*time.Second, rt); err != nil {
		return nil, err
	}
	//httpSwitcher := middlewares.NewHandlerSwitcher(http.NotFoundHandler())
	var route http.Handler
	if route, err = replacehost.New(ctx, proxyRoute, *config.Cfg.Middleware.ReplaceHost); err != nil {
		return nil, err
	}
	httpSwitcher := middlewares.NewHandlerSwitcher(route)
	var handler http.Handler
	if handler, err = forwardedheaders.NewXForwarded(
		configuration.ForwardedHeaders.Insecure,
		configuration.ForwardedHeaders.TrustedIPs,
		httpSwitcher); err != nil {
		return nil, err
	}
	// todo 修改
	cfs := new(tls2.Certificates)
	var tlsConfig *tls.Config
	if tlsConfig, err = cfs.CreateTLSConfig("default"); err != nil {
		return nil, err
	}
	srv := &http.Server{
		Handler:      handler,
		TLSConfig:    tlsConfig,
		ReadTimeout:  configuration.Transport.RespondingTimeouts.ReadTimeout,
		WriteTimeout: configuration.Transport.RespondingTimeouts.WriteTimeout,
		IdleTimeout:  configuration.Transport.RespondingTimeouts.IdleTimeout,
	}

	return &EntryPoint{
		listener:     listener,
		switcher:     httpSwitcher,
		ctx:          ctx,
		server:       srv,
		serverConfig: configuration,
	}, nil
}

func (ep *EntryPoint) Start() {
	go func() {
		if err := ep.server.Serve(ep.listener); err != nil {
			logger.FromContext(ep.ctx).Errorf("Error while starting server: %v", err)
		}
	}()
	go func() {
		if err := ep.server.ServeTLS(ep.listener, "", ""); err != nil {
			logger.FromContext(ep.ctx).Errorf("Error while starting server: %v", err)
		}
	}()
}

// Shutdown stops the http connections.
func (ep *EntryPoint) Shutdown() {
	log := logger.FromContext(ep.ctx)

	reqAcceptGraceTimeOut := ep.serverConfig.Transport.LifeCycle.RequestAcceptGraceTimeout
	if reqAcceptGraceTimeOut > 0 {
		log.Infof("Waiting %s for incoming requests to cease", reqAcceptGraceTimeOut)
		time.Sleep(reqAcceptGraceTimeOut)
	}

	graceTimeOut := ep.serverConfig.Transport.LifeCycle.GraceTimeOut
	ctx, cancel := context.WithTimeout(ep.ctx, graceTimeOut)
	log.Debugf("Waiting %s seconds before killing connections.", graceTimeOut)

	if ep.server != nil {
		func(server *http.Server) {
			err := server.Shutdown(ctx)
			if err == nil {
				return
			}
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				log.Debugf("server failed to shutdown within deadline because: %s", err)
				if err = server.Close(); err != nil {
					log.Error(err.Error())
				}
				return
			}
			log.Error(err.Error())
			// We expect Close to fail again because Shutdown most likely failed when trying to close a listener.
			// We still call it however, to make sure that all connections get closed as well.
			server.Close()
		}(ep.server)
	}
	cancel()
}

// SwitchRouter switches the http router handler.
func (ep *EntryPoint) SwitchRouter(handler http.Handler) {
	if handler == nil {
		return
	}
	ep.switcher.UpdateHandler(handler)
}
