// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/30

package config

import (
	"os"

	"gopkg.in/yaml.v3"
)

var Cfg *Config

func InitConfig(path string) {
	configYaml, err := LoadYaml(path)
	if err != nil {
		panic(err)
	}
	Cfg = configYaml
}

func LoadYaml(path string) (*Config, error) {
	file, err := os.Open(path)
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
