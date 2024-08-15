// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

type ServerService interface {
	Start() error
	Stop() error
	AddServer(name string) (*ServerModel, error)
	GetServerList() ([]*ServerModel, error)
	FindServer(name string) (*ServerModel, error)
	DeployServer(name string) (*ServerModel, error)
	StartServer(name string) (*ServerModel, error)
	StopServer(name string) (*ServerModel, error)
	RestartServer(name string) (*ServerModel, error)
	DeleteServer(name string) (*ServerModel, error)
	GetVNC(name string) (string, error)
	SetVNCPassword(name, password string) error
}
