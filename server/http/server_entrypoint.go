// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/1

package http

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"

	"github.com/crochee/proxy/config"
	"github.com/crochee/proxy/logger"
	"github.com/crochee/proxy/middlewares"
	"github.com/crochee/proxy/middlewares/forwardedheaders"
	tls2 "github.com/crochee/proxy/tls"
)

// EntryPoint is the http server.
type EntryPoint struct {
	listener     net.Listener
	switcher     *middlewares.HTTPHandlerSwitcher
	Server       *http.Server
	ctx          context.Context
	serverConfig *config.EntryPoint
}

// NewEntryPoint creates a new EntryPoint.
func NewEntryPoint(ctx context.Context, configuration *config.EntryPoint) (*EntryPoint, error) {
	listener, err := net.Listen("tcp", configuration.GetAddress())
	if err != nil {
		return nil, fmt.Errorf("error opening listener: %w", err)
	}
	httpSwitcher := middlewares.NewHandlerSwitcher(http.NotFoundHandler())
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
		Server:       srv,
		serverConfig: configuration,
	}, nil
}

func (ep *EntryPoint) Start() {
	go func() {
		if err := ep.Server.Serve(ep.listener); err != nil {
			logger.FromContext(ep.ctx).Errorf("Error while starting server: %v", err)
		}
	}()
	go func() {
		if err := ep.Server.ServeTLS(ep.listener, "", ""); err != nil {
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

	if ep.Server != nil {
		func(server *http.Server) {
			err := server.Shutdown(ctx)
			if err == nil {
				return
			}
			if errors.Is(ctx.Err(), context.DeadlineExceeded) {
				log.Debugf("Server failed to shutdown within deadline because: %s", err)
				if err = server.Close(); err != nil {
					log.Error(err.Error())
				}
				return
			}
			log.Error(err.Error())
			// We expect Close to fail again because Shutdown most likely failed when trying to close a listener.
			// We still call it however, to make sure that all connections get closed as well.
			server.Close()
		}(ep.Server)
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
