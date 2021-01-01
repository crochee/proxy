// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/1

package main

import (
	"context"
	"log"
	"os"
	"strings"

	"github.com/urfave/cli/v2"

	"github.com/crochee/proxy/cmd"
	"github.com/crochee/proxy/config"
	"github.com/crochee/proxy/logger"
	"github.com/crochee/proxy/service"
)

func main() {
	app := cli.NewApp()
	app.Name = "proxy"
	app.Version = cmd.Version()
	app.Usage = "Generates proxy"

	app.Commands = cli.Commands{
		{
			Name:        "proxy",
			Aliases:     []string{"p"},
			Usage:       "proxy server",
			After:       nil,
			Action:      Run,
			Subcommands: nil,
			Flags:       BeforeFlags,
		},
		{
			Name:    "tls",
			Aliases: []string{"t"},
			Usage:   "generates random TLS certificates",
			Action:  Certificate,
			Flags:   append(BeforeFlags, TlsFlags...),
		},
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

var BeforeFlags = []cli.Flag{
	&cli.BoolFlag{
		Name:    "enable-log",
		Usage:   "enable log switch",
		EnvVars: []string{"enable_log"},
	},
	&cli.StringFlag{
		Name:    "log-path",
		Usage:   "log path",
		EnvVars: []string{"log_path"},
	},
	&cli.StringFlag{
		Name:    "log-level",
		Usage:   "log level",
		EnvVars: []string{"log_level"},
	},
}

func Run(c *cli.Context) error {
	ctx := logger.With(context.Background(),
		logger.Enable(c.Bool("enable-log")),
		logger.Level(strings.ToUpper(c.String("log-level"))),
		logger.LogPath(c.String("log-path")),
	)
	Test(ctx)
	return nil
}

func setup(ctx context.Context, cfg *config.Config) error {
	roundTripperManager := service.NewRoundTripperManager()

	roundTripperManager.Update(map[string]*config.ServersTransport{})
	return nil
}

func Test(ctx context.Context) {
	logger.FromContext(ctx).Info("test")
}
