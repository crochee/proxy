// Copyright 2020, The Go Authors. All rights reserved.
// Author: OnlyOneFace
// Date: 2020/12/30

package router

type Balancer interface {
	Servers() []string
	RemoveServer(host string) error
	UpsertServer(host string, options ...ServerOption) error
	Next() string
}

type ServerOption func(*server) error

type server struct {
	host   string
	weight int
	status ServerStatus
}

type ServerStatus uint8

const (
	Down ServerStatus = 0
	Up   ServerStatus = 1
)

type Robin struct {
	serverList []*server
}

func (r *Robin) Servers() []string {
	list := make([]string, 0, len(r.serverList))
	for _, server := range r.serverList {
		if server.weight != 0 && server.status != Down {
			list = append(list, server.host)
		}
	}
	return list
}

func (r Robin) RemoveServer(host string) error {
	for _, server := range r.serverList {
		if server.host == host {
			server.status = Down
		}
	}
	return nil
}

func (r Robin) UpsertServer(host string, options ...ServerOption) error {
	for _, option := range options {
		if err := option(&server{
			host:   host,
			weight: 1,
			status: Up,
		}); err != nil {
			return err
		}
	}
	return nil
}

func (r *Robin) Next() string {
	if len(r.serverList) != 0 {
		return r.serverList[0].host
	}
	return ""
}
