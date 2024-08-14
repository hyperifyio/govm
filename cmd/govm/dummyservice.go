// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"fmt"
	"log"
	"time"
)

type DummyService struct {
	servers []*ServerModel
}

func NewDummyService() *DummyService {
	return &DummyService{}
}

func (s *DummyService) Start() error {
	log.Printf("Starting dummy service")
	return nil
}

func (s *DummyService) Stop() error {
	log.Printf("Stopping dummy service")
	return nil
}

func (s *DummyService) AddServer(
	name string,
) (*ServerModel, error) {
	item := NewServerModel(name, UninitializedServerStatusCode)
	s.servers = append(s.servers, item)
	return item, nil
}

func (s *DummyService) GetServerList() ([]*ServerModel, error) {
	return s.servers, nil
}

func (s *DummyService) FindServer(name string) (*ServerModel, error) {
	for _, state := range s.servers {
		if state.Name == name {
			return state, nil
		}
	}
	return nil, nil
}

func (s *DummyService) DeployServer(name string) (*ServerModel, error) {
	server, err := s.FindServer(name)
	if err != nil {
		return nil, fmt.Errorf("DeployServer: failed to find the server: error: %v", err)
	}
	if server == nil {
		return nil, fmt.Errorf("DeployServer: failed to find the server: not found")
	}
	if server.Status == UninitializedServerStatusCode {
		server.Status = DeployingServerStatusCode
		time.AfterFunc(3*time.Second, func() {
			server.Status = StoppedServerStatusCode
		})
	}
	return server, nil
}

func (s *DummyService) StartServer(name string) (*ServerModel, error) {
	server, err := s.FindServer(name)
	if err != nil {
		return nil, fmt.Errorf("StartServer: failed to find the server: error: %v", err)
	}
	if server == nil {
		return nil, fmt.Errorf("StartServer: failed to find the server: not found")
	}
	if server.Status == StoppedServerStatusCode {
		server.Status = StartingServerStatusCode
		time.AfterFunc(3*time.Second, func() {
			server.Status = StartedServerStatusCode
		})
	}
	return server, nil
}

func (s *DummyService) StopServer(name string) (*ServerModel, error) {
	server, err := s.FindServer(name)
	if err != nil {
		return nil, fmt.Errorf("StopServer: failed to find the server: error: %v", err)
	}
	if server == nil {
		return nil, fmt.Errorf("StopServer: failed to find the server: not found")
	}
	if server.Status == StartedServerStatusCode {
		server.Status = StoppingServerStatusCode
		time.AfterFunc(3*time.Second, func() {
			server.Status = StoppedServerStatusCode
		})
	}
	return server, nil
}

func (s *DummyService) RestartServer(name string) (*ServerModel, error) {
	server, err := s.FindServer(name)
	if err != nil {
		return nil, fmt.Errorf("RestartServer: failed to find the server: error: %v", err)
	}
	if server == nil {
		return nil, fmt.Errorf("RestartServer: failed to find the server: not found")
	}
	if server.Status == StoppedServerStatusCode {
		server.Status = StartingServerStatusCode
		time.AfterFunc(3*time.Second, func() {
			server.Status = StartedServerStatusCode
		})
	} else if server.Status == StartedServerStatusCode {
		server.Status = StoppingServerStatusCode
		time.AfterFunc(3*time.Second, func() {
			server.Status = StartingServerStatusCode
			time.AfterFunc(3*time.Second, func() {
				server.Status = StartedServerStatusCode
			})
		})
	}
	return server, nil
}

func (s *DummyService) DeleteServer(name string) (*ServerModel, error) {
	server, err := s.FindServer(name)
	if err != nil {
		return nil, fmt.Errorf("DeleteServer: failed to find the server: error: %v", err)
	}
	if server == nil {
		return nil, fmt.Errorf("DeleteServer: failed to find the server: not found")
	}
	if server.Status == StoppedServerStatusCode || server.Status == UninitializedServerStatusCode {
		server.Status = DeletingServerStatusCode
		time.AfterFunc(3*time.Second, func() {
			err := s.removeServer(name)
			if err != nil {
				log.Printf("ERROR: failed to remove server: %s: %v", name, err)
			}
		})
	}
	return server, nil
}

func (s *DummyService) removeServer(name string) error {
	for i, server := range s.servers {
		if server.Name == name {
			s.servers = append(s.servers[:i], s.servers[i+1:]...)
			return nil
		}
	}
	return nil
}

var _ ServerService = &DummyService{}
