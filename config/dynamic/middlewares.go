// Copyright (c) Huawei Technologies Co., Ltd. 2020-2020. All rights reserved.
// Description:
// Author: l30002214
// Create: 2021/1/7

// Package dynamic
package dynamic

import "time"

// Middleware holds the Middleware configuration.
type Middleware struct {
	AddPrefix        *AddPrefix        `yaml:"addPrefix,omitempty"`
	ReplaceHost      *ReplaceHost      `yaml:"replaceHost,omitempty"`
	ReplacePath      *ReplacePath      `yaml:"replacePath,omitempty"`
	ReplacePathRegex *ReplacePathRegex `yaml:"replacePathRegex,omitempty"`
	RateLimit        *RateLimit        `yaml:"rateLimit,omitempty"`
	CircuitBreaker   *CircuitBreaker   `yaml:"circuitBreaker,omitempty"`
	Retry            *Retry            `yaml:"retry,omitempty"`
}

// AddPrefix holds the AddPrefix configuration.
type AddPrefix struct {
	Prefix string `yaml:"prefix,omitempty"`
}

type ReplaceHost struct {
	Scheme string `yaml:"scheme,omitempty"`
	Host   string `yaml:"host,omitempty"`
}

// ReplacePath holds the ReplacePath configuration.
type ReplacePath struct {
	Path string `yaml:"path,omitempty"`
}

// ReplacePathRegex holds the ReplacePathRegex configuration.
type ReplacePathRegex struct {
	Regex       string `yaml:"regex,omitempty"`
	Replacement string `yaml:"replacement,omitempty"`
}

type RateLimit struct {
	Every time.Duration `yaml:"every,omitempty"`
	Burst int           `yaml:"burst,omitempty"`
}

// CircuitBreaker holds the circuit breaker configuration.
type CircuitBreaker struct {
	Expression string `yaml:"expression,omitempty"`
}

// Retry holds the retry configuration.
type Retry struct {
	Attempts        int           `yaml:"attempts,omitempty"`
	InitialInterval time.Duration `yaml:"initialInterval,omitempty"`
}
