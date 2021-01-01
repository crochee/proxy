// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/1

package server

import (
	"context"
	"errors"
	"time"

	"github.com/crochee/proxy/logger"
	"github.com/crochee/proxy/safe"
)

type Watcher interface {
}

// NewServer returns an initialized Server.
func NewServer(ctx context.Context, routinesPool *safe.Pool) *Server {
	return &Server{
		ctx:            ctx,
		routinesPool:   routinesPool,
		watcher:        nil,
		tcpEntryPoints: nil,
		stopChan:       make(chan bool, 1),
	}
}

type Server struct {
	ctx            context.Context
	routinesPool   *safe.Pool
	watcher        Watcher
	tcpEntryPoints map[string]interface{}
	stopChan       chan bool
}

func (s *Server) Start() {
	go func() {
		<-s.ctx.Done()
		log := logger.FromContext(s.ctx)
		log.Info("I have to go...")
		log.Info("Stopping server gracefully")
		s.Stop()
	}()
}

// Wait blocks until the server shutdown.
func (s *Server) Wait() {
	<-s.stopChan
}

// Stop stops the server.
func (s *Server) Stop() {
	s.stopChan <- true
	logger.FromContext(s.ctx).Info("Server stopped")
}

// Close destroys the server.
func (s *Server) Close() {
	ctx, cancel := context.WithTimeout(s.ctx, 10*time.Second)

	go func(ctx context.Context) {
		<-ctx.Done()
		if errors.Is(ctx.Err(), context.Canceled) {
			return
		} else if errors.Is(ctx.Err(), context.DeadlineExceeded) {
			panic("Timeout while stopping proxy, killing instance ✝")
		}
	}(ctx)

	s.routinesPool.Stop()

	close(s.stopChan)
	cancel()
}
