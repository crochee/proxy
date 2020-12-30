// Copyright 2020, The Go Authors. All rights reserved.
// Author: crochee
// Date: 2020/12/30

package main

import (
	"context"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/crochee/proxy/config"
	"github.com/crochee/proxy/logger"
)

func main() {
	logger.InitLogger()
	config.InitConfig()
	srv := &http.Server{
		Addr:    ":80",
		Handler: nil,
	}
	go func() {
		logger.Info("proxy running...")
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatal(err.Error())
		}
	}()
	quit := make(chan os.Signal)
	// kill (no param) default send syscall.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall.SIGKILL but can't be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	// The context is used to inform the server it has 5 seconds to finish
	// the request it is currently handling
	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown:%v", err)
	}
	logger.Info("proxy server exit!")
}
