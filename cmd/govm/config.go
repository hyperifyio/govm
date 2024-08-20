// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

// Config holds the overall configuration
type Config struct {
	Servers ServerConfigList `yaml:"servers"`
}

func NewConfig(
	config ServerConfigList,
) *Config {
	return &Config{
		config,
	}
}

// AddServer adds a new server in the config and returns a new config object
func (c *Config) AddServer(name string, users UserEmailList) *Config {
	newList := append(c.Servers, NewServerConfig(name, users))
	return NewConfig(newList)
}

// ServerHasAccessToEmail checks if an email address exists in the allowed list of email addresses to have access
func (c *Config) ServerHasAccessToEmail(name, email string) bool {
	item := c.Servers.findByName(name)
	if item == nil {
		return false
	}
	return item.Users.contains(email)
}

// LoadConfig reads and parses the YAML configuration file
func LoadConfig(filename string) (*Config, error) {
	var config *Config

	// Read file
	file, err := os.Open(filename)
	if err != nil {
		return config, err
	}
	defer file.Close()

	// Decode YAML
	decoder := yaml.NewDecoder(file)
	err = decoder.Decode(&config)
	if err != nil {
		return config, err
	}

	return config, nil
}

// SaveConfig writes the Config struct to a YAML file
func (c *Config) SaveConfig(filename string) error {

	// Write to a temporary file
	tmpFilename := filename + ".tmp"
	tmpFile, err := os.Create(tmpFilename)
	if err != nil {
		return fmt.Errorf("SaveConfig: create: %s: %w", filename, err)
	}
	defer tmpFile.Close()

	encoder := yaml.NewEncoder(tmpFile)
	err = encoder.Encode(&c)
	if err != nil {
		return fmt.Errorf("SaveConfig: encoding: %s: %w", filename, err)
	}
	encoder.Close()

	// Backup the original file if it exists
	bakFilename := filename + ".bak"
	if _, err := os.Stat(filename); err == nil {
		err = os.Rename(filename, bakFilename)
		if err != nil {
			return fmt.Errorf("SaveConfig: backup: %s: %w", filename, err)
		}
	}

	// Replace the original file with the temp file
	err = os.Rename(tmpFilename, filename)
	if err != nil {
		return fmt.Errorf("SaveConfig: rename: %s: %w", filename, err)
	}

	// Remove the backup file
	err = os.Remove(bakFilename)
	if err != nil {
		return fmt.Errorf("SaveConfig: remove backup: %s: %w", bakFilename, err)
	}

	return nil
}
