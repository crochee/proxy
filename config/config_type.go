// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/30

package config

type Config struct {
	List []*ProxyHost   `yaml:"list,omitempty"`
	Spec EntryPointList `yaml:"spec,omitempty"`
	//Transport *ServersTransport `json:"transport" yaml:"transport"`
}

type ProxyHost struct {
	Origin []string `yaml:"origin,omitempty"`
	Target []string `yaml:"target,omitempty"`
}
