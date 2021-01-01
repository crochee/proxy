// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/1

package main

import (
	"context"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/crochee/proxy/logger"
	"github.com/crochee/proxy/tls/generate"
)

var TlsFlags = []cli.Flag{
	&cli.StringFlag{
		Name:    "cert",
		Aliases: []string{"c"},
		Usage:   "cert path",
		EnvVars: []string{"cert_path"},
	},
	&cli.StringFlag{
		Name:    "key",
		Aliases: []string{"k"},
		Usage:   "key path",
		EnvVars: []string{"key_path"},
	},
}

func Certificate(c *cli.Context) error {
	ctx := logger.With(context.Background(),
		logger.Enable(c.Bool("enable-log")),
		logger.Level(strings.ToUpper(c.String("log-level"))),
		logger.LogPath(c.String("log-path")),
	)
	logger.FromContext(ctx).Info("generates random TLS certificates start!")
	_, err := generate.DefaultCertificate(c.String("cert"), c.String("key"))
	if err != nil {
		logger.FromContext(ctx).Error(err.Error())
	}
	return nil
}
