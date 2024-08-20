// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"fmt"
	"sync"
)

// ConfigManager handles the config and background writing
type ConfigManager struct {
	config      *Config
	queue       chan string
	configMutex sync.Mutex
	filename    string
}

// NewConfigManager creates a new ConfigManager
func NewConfigManager(filename string, config *Config) *ConfigManager {
	manager := &ConfigManager{
		config:   config,
		queue:    make(chan string, ConfigManagerBufferSize),
		filename: filename,
	}
	go manager.runWorker()
	return manager
}

func (m *ConfigManager) GetConfig() *Config {
	m.configMutex.Lock()
	ret := m.config
	m.configMutex.Unlock()
	return ret
}

// AddServerConfig adds a server to the config and queues a write operation
func (m *ConfigManager) AddServerConfig(name string, users UserEmailList) {

	m.configMutex.Lock()
	m.config = m.config.AddServer(name, users)
	m.configMutex.Unlock()

	// Queue the write operation
	m.queue <- name
}

// runWorker processes the queue in the background
func (m *ConfigManager) runWorker() {
	for range m.queue {

		m.configMutex.Lock()
		err := m.saveConfig()
		m.configMutex.Unlock()

		if err != nil {
			fmt.Println("ConfigManager: Error saving config:", err)
		} else {
			fmt.Println("ConfigManager: Config saved successfully")
		}
	}
}

// saveConfig writes the Config struct to a YAML file safely
func (m *ConfigManager) saveConfig() error {
	return m.config.SaveConfig(m.filename)
}
