// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/30

package config

import (
	"os"
	"testing"
	"time"

	"gopkg.in/yaml.v3"
)

func TestLoadYaml(t *testing.T) {
	cf := &Config{
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
