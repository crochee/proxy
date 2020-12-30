// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/30

package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var Cfg *Config

func InitConfig() {
	configYaml, err := loadYaml()
	if err != nil {
		panic(err)
	}
	Cfg = configYaml
}

func loadYaml() (*Config, error) {
	configPath, ok := os.LookupEnv("config_path")
	if !ok {
		configPath = "D:/project/obs/conf/config.yml"
	}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	var config Config
	if err = yaml.NewDecoder(file).Decode(&config); err != nil {
		return nil, err
	}
	return &config, nil
}
