// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

// ServerConfig represents a server and the users that have access to it
type ServerConfig struct {
	Name  string        `yaml:"name"`
	Users UserEmailList `yaml:"users"`
}

func NewServerConfig(
	name string,
	users UserEmailList,
) *ServerConfig {
	return &ServerConfig{
		name,
		users,
	}
}

// emailHasAccess checks if an email address exists in the allowed list of email addresses to have access
func (item *ServerConfig) emailHasAccess(email string) bool {
	return item.Users.contains(email)
}
