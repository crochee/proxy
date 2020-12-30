// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/30

package config

type Config struct {
	List []*ProxyHost `json:"list" yaml:"list"`
}

type ProxyHost struct {
	Origin []string `json:"origin" yaml:"origin"`
	Target []string `json:"target" yaml:"target"`
}
