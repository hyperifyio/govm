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
const BlockedServerStatus string = "blocked"
const PausedServerStatus string = "paused"
const CrashedServerStatus string = "crashed"
const SuspendedServerStatus string = "suspended"
const UnknownServerStatus string = "unknown"
const DeletingServerStatus string = "deleting"
const DeletedServerStatus string = "deleted"

const (
	UninitializedServerStatusCode ServerStatusCode = iota
	DeployingServerStatusCode
	StoppedServerStatusCode
	StartingServerStatusCode
	StoppingServerStatusCode
	StartedServerStatusCode
	BlockedServerStatusCode
	PausedServerStatusCode
	CrashedServerStatusCode
	SuspendedServerStatusCode
	UnknownServerStatusCode
	DeletingServerStatusCode
	DeletedServerStatusCode
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
		BlockedServerStatus,
		PausedServerStatus,
		CrashedServerStatus,
		SuspendedServerStatus,
		UnknownServerStatus,
		DeletingServerStatus,
		DeletedServerStatus,
	}[d]
}

func (d ServerStatusCode) GetAvailableActions() []ServerActionCode {
	switch d {
	case UninitializedServerStatusCode:
		return []ServerActionCode{DeployServerActionCode, DeleteServerActionCode}
	case StoppedServerStatusCode:
		return []ServerActionCode{StartServerActionCode, DeleteServerActionCode}
	case StartedServerStatusCode:
		return []ServerActionCode{StopServerActionCode, RestartServerActionCode}
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
	case BlockedServerStatus:
		return BlockedServerStatusCode, nil
	case PausedServerStatus:
		return PausedServerStatusCode, nil
	case CrashedServerStatus:
		return CrashedServerStatusCode, nil
	case SuspendedServerStatus:
		return SuspendedServerStatusCode, nil
	case DeletedServerStatus:
		return DeletedServerStatusCode, nil
	case UnknownServerStatus:
		return UnknownServerStatusCode, nil
	default:
		return -1, fmt.Errorf("unknown server status code: %s", name)
	}
}
