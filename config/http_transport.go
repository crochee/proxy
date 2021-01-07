// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2021/1/6

package config

import (
	"time"

	"github.com/crochee/proxy/tls"
)

type ServerName string

// ServersTransport options to configure communication between Traefik and the servers.
type ServersTransport struct {
	ServerName         ServerName          `yaml:"serverName,omitempty"`
	InsecureSkipVerify bool                `yaml:"insecureSkipVerify,omitempty"`
	RootCAs            []tls.FileOrContent `yaml:"rootCAs,omitempty"`
	Certificates       tls.Certificates    `yaml:"certificates,omitempty"`
	MaxIdleConnPerHost int                 `yaml:"maxIdleConnPerHost,omitempty"`
	ForwardingTimeouts *ForwardingTimeouts `yaml:"forwardingTimeouts,omitempty"`
}

// ForwardingTimeouts contains timeout configurations for forwarding requests to the backend servers.
type ForwardingTimeouts struct {
	DialTimeout           time.Duration `yaml:"dialTimeout,omitempty"`
	ResponseHeaderTimeout time.Duration `yaml:"responseHeaderTimeout,omitempty"`
	IdleConnTimeout       time.Duration `yaml:"idleConnTimeout,omitempty"`
}
