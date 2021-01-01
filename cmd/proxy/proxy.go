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
		Name:    "enableLog",
		Aliases: []string{"el"},
		Usage:   "enable log switch",
		EnvVars: []string{"enable_log"},
	},
	&cli.StringFlag{
		Name:    "logPath",
		Aliases: []string{"lp"},
		Usage:   "log path",
		EnvVars: []string{"log_path"},
	},
	&cli.StringFlag{
		Name:    "logLevel",
		Aliases: []string{"ll"},
		Usage:   "log level",
		EnvVars: []string{"log_level"},
	},
}

func Before(ctx *cli.Context) error {
	if !ctx.Bool("enableLog") {
		return nil
	}
	logger.InitLogger(ctx.String("logPath"), ctx.String("logLevel"))
	return nil
}

func Run(ctx *cli.Context) error {
	logger.Info("test")
	return nil
}
