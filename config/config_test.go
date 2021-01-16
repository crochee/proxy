// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/30

package config

import (
	"os"
	"testing"
	"time"

	"gopkg.in/yaml.v3"

	"github.com/crochee/proxy/config/dynamic"
)

func TestLoadYaml(t *testing.T) {
	cf := &Config{
		List: nil,
		Spec: EntryPointList{
			"proxy": &EntryPoint{
				Port:     8085,
				Protocol: "TCP",
				Transport: &EntryPointsTransport{
					LifeCycle: &LifeCycle{
						RequestAcceptGraceTimeout: 1 * time.Second,
						GraceTimeOut:              5 * time.Second,
					},
					RespondingTimeouts: &RespondingTimeouts{
						ReadTimeout:  0,
						WriteTimeout: 0,
						IdleTimeout:  3 * time.Minute,
					},
				},
				ForwardedHeaders: &ForwardedHeaders{
					Insecure:   true,
					TrustedIPs: []string{},
				},
			},
		},
		Transport: &ServersTransport{
			ServerName:         "default",
			InsecureSkipVerify: true,
			RootCAs:            nil,
			Certificates:       nil,
			MaxIdleConnPerHost: 100,
			ForwardingTimeouts: nil,
		},
		Middleware: &dynamic.Middleware{
			ReplaceHost: &dynamic.ReplaceHost{
				Scheme: "http",
				Host:   "127.0.0.1:8150",
			},
		},
	}
	configPath, ok := os.LookupEnv("config_path")
	if !ok {
		configPath = "D:/project/proxy/conf/config.yml"
	}
	file, err := os.Create(configPath)
	if err != nil {
		t.Fatal(err)
	}
	defer file.Close()
	if err = yaml.NewEncoder(file).Encode(cf); err != nil {
		t.Fatal(err)
	}
}
