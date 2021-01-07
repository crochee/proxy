// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/6

package http

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"sync"

	"github.com/crochee/proxy/config"
	"github.com/crochee/proxy/logger"
)

type EntryPointList map[config.ServerName]*EntryPoint

// NewTCPEntryPoints creates a new TCPEntryPoints.
func NewEntryPointList(entryPointsConfig config.EntryPointList) (EntryPointList, error) {
	serverEntryPointList := make(EntryPointList, len(entryPointsConfig))
	for entryPointName, entryPoint := range entryPointsConfig {
		protocol, err := entryPoint.GetProtocol()
		if err != nil {
			return nil, fmt.Errorf("error while building entryPoint %s: %w", entryPointName, err)
		}
		if protocol != "tcp" {
			continue
		}
		ctx := logger.With(context.Background(), logger.Enable(true),
			logger.Level(strings.ToUpper("DEBUG")),
			logger.LogPath(fmt.Sprintf("./log/%s.log", entryPointName)))
		serverEntryPointList[entryPointName], err = NewEntryPoint(ctx, entryPoint)
		if err != nil {
			return nil, fmt.Errorf("error while building entryPoint %s: %w", entryPointName, err)
		}
	}
	return serverEntryPointList, nil
}

// Start the server entry points.
func (epl EntryPointList) Start() {
	for entryPointName, serverEntryPoint := range epl {
		logger.FromContext(serverEntryPoint.ctx).Debugf("start %s", entryPointName)
		go serverEntryPoint.Start()
	}
}

// Stop the server entry points.
func (epl EntryPointList) Stop() {
	var wg sync.WaitGroup

	for epn, ep := range epl {
		wg.Add(1)

		go func(entryPointName config.ServerName, entryPoint *EntryPoint) {
			defer wg.Done()
			entryPoint.Shutdown()
			logger.FromContext(entryPoint.ctx).Debugf("Entry point %s closed", entryPointName)
		}(epn, ep)
	}

	wg.Wait()
}

// Update the servers.
func (epl EntryPointList) Update(entryPointsConfig config.EntryPointList) {
	// todo 需要实现
}

// Switch the routers.
func (epl EntryPointList) Switch(routers map[config.ServerName]http.Handler) {
	for entryPointName, rt := range routers {
		epl[entryPointName].SwitchRouter(rt)
	}
}
