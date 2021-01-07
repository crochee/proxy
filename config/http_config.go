// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/31

package config

import (
	"fmt"
	"strings"
	"time"
)

const (

	// prior to shutting down.
	DefaultGraceTimeout = 10 * time.Second

	// DefaultIdleTimeout before closing an idle connection.
	DefaultIdleTimeout = 180 * time.Second
)

type EntryPointList map[ServerName]*EntryPoint

// EntryPoint holds the entry point configuration.
type EntryPoint struct {
	Port             int                   `yaml:"port,omitempty"`
	Protocol         string                `yaml:"protocol,omitempty"`
	Transport        *EntryPointsTransport `yaml:"transport,omitempty"`
	ForwardedHeaders *ForwardedHeaders     `yaml:"forwardedHeaders,omitempty"`
}

// GetAddress strips any potential protocol part of the address field of the
// entry point, in order to return the actual address.
func (ep EntryPoint) GetPort() string {
	return fmt.Sprintf(":%d", ep.Port)
}

// GetProtocol returns the protocol part of the address field of the entry point.
// If none is specified, it defaults to "tcp".
func (ep EntryPoint) GetProtocol() (string, error) {
	protocol := strings.ToLower(ep.Protocol)
	if protocol == "tcp" || protocol == "udp" {
		return protocol, nil
	}

	return "", fmt.Errorf("invalid protocol: %s", ep.Protocol)
}

// SetDefaults sets the default values.
func (ep *EntryPoint) SetDefaults() {
	ep.ForwardedHeaders = &ForwardedHeaders{}
}

// EntryPointsTransport configures communication between clients and Traefik.
type EntryPointsTransport struct {
	LifeCycle          *LifeCycle          `yaml:"lifeCycle,omitempty"`
	RespondingTimeouts *RespondingTimeouts `yaml:"respondingTimeouts,omitempty"`
}

// SetDefaults sets the default values.
func (t *EntryPointsTransport) SetDefaults() {
	t.LifeCycle = &LifeCycle{}
	t.LifeCycle.SetDefaults()
	t.RespondingTimeouts = &RespondingTimeouts{}
	t.RespondingTimeouts.SetDefaults()
}

// ForwardedHeaders Trust client forwarding headers.
type ForwardedHeaders struct {
	Insecure   bool     `yaml:"insecure,omitempty"`
	TrustedIPs []string `yaml:"trustedIPs,omitempty"`
}

// LifeCycle contains configurations relevant to the lifecycle (such as the shutdown phase) of Traefik.
type LifeCycle struct {
	RequestAcceptGraceTimeout time.Duration `yaml:"requestAcceptGraceTimeout,omitempty"`
	GraceTimeOut              time.Duration `yaml:"graceTimeOut,omitempty"`
}

// SetDefaults sets the default values.
func (a *LifeCycle) SetDefaults() {
	a.GraceTimeOut = DefaultGraceTimeout
}

// RespondingTimeouts contains timeout configurations for incoming requests to the Traefik instance.
type RespondingTimeouts struct {
	ReadTimeout  time.Duration `yaml:"readTimeout,omitempty"`
	WriteTimeout time.Duration `yaml:"writeTimeout,omitempty"`
	IdleTimeout  time.Duration `yaml:"idleTimeout,omitempty"`
}

// SetDefaults sets the default values.
func (a *RespondingTimeouts) SetDefaults() {
	a.IdleTimeout = DefaultIdleTimeout
}
