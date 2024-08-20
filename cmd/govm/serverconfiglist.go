// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

type ServerConfigList []*ServerConfig

// findByName finds a server by name and returns it, otherwise nil
func (list ServerConfigList) findByName(name string) *ServerConfig {
	for _, item := range list {
		if item.Name == name {
			return item
		}
	}
	return nil
}

// hasByName finds a server by name and returns true if it exists
func (list ServerConfigList) hasByName(name string) bool {
	for _, item := range list {
		if item.Name == name {
			return true
		}
	}
	return false
}
