// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/1

package main

import (
	"github.com/crochee/proxy/logger"
	"github.com/urfave/cli/v2"

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

func Certificate(ctx *cli.Context) error {
	logger.Info("generates random TLS certificates start!")
	_, err := generate.DefaultCertificate(ctx.String("cert"), ctx.String("key"))
	return err
}
