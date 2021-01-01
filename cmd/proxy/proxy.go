// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/1

package main

import (
	"log"
	"os"

	"github.com/urfave/cli/v2"

	"github.com/crochee/proxy/cmd"
	"github.com/crochee/proxy/logger"
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
			Before:      Before,
			After:       nil,
			Action:      Run,
			Subcommands: nil,
			Flags:       BeforeFlags,
		},
		{
			Name:    "tls",
			Aliases: []string{"t"},
			Usage:   "generates random TLS certificates",
			Before:  Before,
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

func Before(ctx *cli.Context) error {
	if !ctx.Bool("enable-log") {
		return nil
	}
	logger.InitLogger(ctx.String("log-path"), ctx.String("log-level"))
	return nil
}

func Run(ctx *cli.Context) error {
	logger.Info("test")
	return nil
}
