// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/31

package service

import (
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"net"
	"net/http"
	"reflect"
	"sync"
	"time"

	"github.com/crochee/proxy/config"
	"github.com/crochee/proxy/logger"
	tls2 "github.com/crochee/proxy/tls"
)

// NewRoundTripperManager creates a new RoundTripperManager.
func NewRoundTripperManager() *RoundTripperManager {
	return &RoundTripperManager{
		roundTrippers: make(map[string]http.RoundTripper),
		configs:       make(map[string]*config.ServersTransport),
	}
}

// RoundTripperManager handles roundtripper for the reverse proxy.
type RoundTripperManager struct {
	rtLock        sync.RWMutex
	roundTrippers map[string]http.RoundTripper
	configs       map[string]*config.ServersTransport
}

// Update updates the roundtrippers configurations.
func (r *RoundTripperManager) Update(newConfigs map[string]*config.ServersTransport) {
	r.rtLock.Lock()
	defer r.rtLock.Unlock()

	var err error
	// update it have
	for configName, serversTransport := range r.configs {
		newConfig, ok := newConfigs[configName]
		if !ok {
			delete(r.configs, configName)
			delete(r.roundTrippers, configName)
			continue
		}

		if reflect.DeepEqual(newConfig, serversTransport) {
			continue
		}

		r.roundTrippers[configName], err = createRoundTripper(newConfig)
		if err != nil {
			logger.Errorf("Could not configure HTTP Transport %s, fallback on default transport: %v", configName, err)
			r.roundTrippers[configName] = http.DefaultTransport
		}
	}
	// add new
	for newConfigName, newConfig := range newConfigs {
		if _, ok := r.configs[newConfigName]; ok {
			continue
		}

		r.roundTrippers[newConfigName], err = createRoundTripper(newConfig)
		if err != nil {
			logger.Errorf("Could not configure HTTP Transport %s, fallback on default transport: %v", newConfigName, err)
			r.roundTrippers[newConfigName] = http.DefaultTransport
		}
	}

	r.configs = newConfigs
}

// Get get a roundtripper by name.
func (r *RoundTripperManager) Get(name string) (http.RoundTripper, error) {
	if len(name) == 0 {
		name = "default@internal"
	}

	r.rtLock.RLock()
	defer r.rtLock.RUnlock()

	if rt, ok := r.roundTrippers[name]; ok {
		return rt, nil
	}

	return nil, fmt.Errorf("servers transport not found %s", name)
}

// createRoundTripper creates an http.RoundTripper configured with the Transport configuration settings.
// For the settings that can't be configured in Traefik it uses the default http.Transport settings.
// An exception to this is the MaxIdleConns setting as we only provide the option MaxIdleConnsPerHostin Traefik at this point in time.
// Setting this value to the default of 100 could lead to confusing behavior and backwards compatibility issues.
func createRoundTripper(cfg *config.ServersTransport) (http.RoundTripper, error) {
	if cfg == nil {
		return nil, errors.New("no transport configuration given")
	}

	dialer := &net.Dialer{
		Timeout:   30 * time.Second,
		KeepAlive: 30 * time.Second,
	}

	if cfg.ForwardingTimeouts != nil {
		dialer.Timeout = cfg.ForwardingTimeouts.DialTimeout
	}

	transport := &http.Transport{
		Proxy:                 http.ProxyFromEnvironment,
		DialContext:           dialer.DialContext,
		MaxIdleConnsPerHost:   cfg.MaxIdleConnPerHost,
		IdleConnTimeout:       90 * time.Second,
		TLSHandshakeTimeout:   10 * time.Second,
		ExpectContinueTimeout: 1 * time.Second,
	}

	if cfg.ForwardingTimeouts != nil {
		transport.ResponseHeaderTimeout = cfg.ForwardingTimeouts.ResponseHeaderTimeout
		transport.IdleConnTimeout = cfg.ForwardingTimeouts.IdleConnTimeout
	}

	if cfg.InsecureSkipVerify || len(cfg.RootCAs) > 0 || len(cfg.ServerName) > 0 || len(cfg.Certificates) > 0 {
		transport.TLSClientConfig = &tls.Config{
			ServerName:         cfg.ServerName,
			InsecureSkipVerify: cfg.InsecureSkipVerify,
			RootCAs:            createRootCACertPool(cfg.RootCAs),
			Certificates:       cfg.Certificates.GetCertificates(),
		}
	}

	return newSmartRoundTripper(transport)
}

func createRootCACertPool(rootCAs []tls2.FileOrContent) *x509.CertPool {
	if len(rootCAs) == 0 {
		return nil
	}

	roots := x509.NewCertPool()

	for _, cert := range rootCAs {
		certContent, err := cert.Read()
		if err != nil {
			logger.Errorf("Error while read RootCAs,%w", err)
			continue
		}
		roots.AppendCertsFromPEM(certContent)
	}

	return roots
}
