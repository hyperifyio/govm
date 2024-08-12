// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"fmt"
	"strings"
)

// Enum for server actions
type ServerActionCode int

const DeployServerAction string = "deploy"
const StartServerAction string = "start"
const StopServerAction string = "stop"
const RestartServerAction string = "restart"
const DeleteServerAction string = "delete"

const (
	DeployServerActionCode ServerActionCode = iota
	StartServerActionCode
	StopServerActionCode
	RestartServerActionCode
	DeleteServerActionCode
)

// String method to get the name of the server action code
func (d ServerActionCode) String() string {
	return [...]string{
		DeployServerAction,
		StartServerAction,
		StopServerAction,
		RestartServerAction,
		DeleteServerAction,
	}[d]
}

// ParseServerActionCode parses a string to ServerActionCode
func ParseServerActionCode(name string) (ServerActionCode, error) {
	switch strings.ToLower(name) {
	case DeployServerAction:
		return DeployServerActionCode, nil
	case StartServerAction:
		return StartServerActionCode, nil
	case StopServerAction:
		return StopServerActionCode, nil
	case RestartServerAction:
		return RestartServerActionCode, nil
	case DeleteServerAction:
		return DeleteServerActionCode, nil
	default:
		return -1, fmt.Errorf("unknown server action code: %s", name)
	}
}
