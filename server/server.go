// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/1

package server

import (
	"context"
	"os"
)

type Watcher interface {
}

type Server struct {
	watcher        Watcher
	tcpEntryPoints map[string]interface{}
	signals        chan os.Signal
}

func (s *Server) Start(ctx context.Context) {

}
