// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/30

package config

import "github.com/crochee/proxy/config/dynamic"

type Config struct {
	List       []*ProxyHost        `yaml:"list,omitempty"`
	Spec       EntryPointList      `yaml:"spec,omitempty"`
	Transport  *ServersTransport   `yaml:"transport,omitempty"`
	Middleware *dynamic.Middleware `yaml:"middleware,omitempty"`
}

type ProxyHost struct {
	Origin []string `yaml:"origin,omitempty"`
	Target []string `yaml:"target,omitempty"`
}
