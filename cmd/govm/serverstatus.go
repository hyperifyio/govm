// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"fmt"
	"strings"
)

// Enum for server statuses
type ServerStatusCode int

const UninitializedServerStatus string = "uninitialized"
const DeployingServerStatus string = "deploying"
const StoppedServerStatus string = "stopped"
const StartingServerStatus string = "starting"
const StoppingServerStatus string = "stopping"
const StartedServerStatus string = "started"
const DeletingServerStatus string = "deleting"

const (
	UninitializedServerStatusCode ServerStatusCode = iota
	DeployingServerStatusCode
	StoppedServerStatusCode
	StartingServerStatusCode
	StoppingServerStatusCode
	StartedServerStatusCode
	DeletingServerStatusCode
)

// String method to get the name of the server status code
func (d ServerStatusCode) String() string {
	return [...]string{
		UninitializedServerStatus,
		DeployingServerStatus,
		StoppedServerStatus,
		StartingServerStatus,
		StoppingServerStatus,
		StartedServerStatus,
		DeletingServerStatus,
	}[d]
}

func (d ServerStatusCode) GetAvailableActions() []ServerActionCode {
	switch d {
	case UninitializedServerStatusCode:
		return []ServerActionCode{DeployServerActionCode, DeleteServerActionCode}
	case DeployingServerStatusCode:
		return []ServerActionCode{}
	case StoppedServerStatusCode:
		return []ServerActionCode{StartServerActionCode, DeleteServerActionCode}
	case StartingServerStatusCode:
		return []ServerActionCode{}
	case StoppingServerStatusCode:
		return []ServerActionCode{}
	case StartedServerStatusCode:
		return []ServerActionCode{StopServerActionCode, RestartServerActionCode}
	case DeletingServerStatusCode:
		return []ServerActionCode{}
	}
	return []ServerActionCode{}
}

// ParseServerStatusCode parses a string to ServerStatusCode
func ParseServerStatusCode(name string) (ServerStatusCode, error) {
	switch strings.ToLower(name) {
	case UninitializedServerStatus:
		return UninitializedServerStatusCode, nil
	case DeployingServerStatus:
		return DeployingServerStatusCode, nil
	case StoppedServerStatus:
		return StoppedServerStatusCode, nil
	case StartingServerStatus:
		return StartingServerStatusCode, nil
	case StoppingServerStatus:
		return StoppingServerStatusCode, nil
	case StartedServerStatus:
		return StartedServerStatusCode, nil
	case DeletingServerStatus:
		return DeletingServerStatusCode, nil
	default:
		return -1, fmt.Errorf("unknown server status code: %s", name)
	}
}
