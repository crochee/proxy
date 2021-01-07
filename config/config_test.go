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
		List: []*ProxyHost{
			{
				Origin: []string{"localhost:80", "127.0.0.1:80"},
				Target: []string{"localhost:8080", "127.0.0.1:8080"},
			},
			{
				Origin: []string{"localhost:81", "127.0.0.1:81"},
				Target: []string{"localhost:8081", "127.0.0.1:8081"},
			},
			{
				Origin: []string{"localhost:8082", "127.0.0.1:8083"},
				Target: []string{"localhost:8082", "127.0.0.1:8083"},
			},
		},
		Spec: EntryPointList{
			"sever1": &EntryPoint{
				Address: "127.0.0.1",
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
