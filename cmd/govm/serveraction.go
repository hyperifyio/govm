// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"fmt"
	"strings"
)

// Enum for server action strings
type ServerAction string

// Enum for server actions
type ServerActionCode int

const (
	CreateServerAction  = "create"
	DeployServerAction  = "deploy"
	StartServerAction   = "start"
	StopServerAction    = "stop"
	RestartServerAction = "restart"
	DeleteServerAction  = "delete"
	ConsoleServerAction = "console"
)

const (
	CreateServerActionCode ServerActionCode = iota
	DeployServerActionCode
	StartServerActionCode
	StopServerActionCode
	RestartServerActionCode
	DeleteServerActionCode
	ConsoleServerActionCode
)

func AllServerActionCodes() []ServerActionCode {
	return []ServerActionCode{
		CreateServerActionCode,
		DeployServerActionCode,
		StartServerActionCode,
		StopServerActionCode,
		RestartServerActionCode,
		DeleteServerActionCode,
		ConsoleServerActionCode,
	}
}

// String method to get the name of the server action code
func (d ServerActionCode) String() string {
	return [...]string{
		CreateServerAction,
		DeployServerAction,
		StartServerAction,
		StopServerAction,
		RestartServerAction,
		DeleteServerAction,
		ConsoleServerAction,
	}[d]
}

// ServerAction method to get the name of the server action code
func (d ServerActionCode) ServerAction() ServerAction {
	return [...]ServerAction{
		CreateServerAction,
		DeployServerAction,
		StartServerAction,
		StopServerAction,
		RestartServerAction,
		DeleteServerAction,
		ConsoleServerAction,
	}[d]
}

// ParseServerActionCode parses a string to ServerActionCode
func ParseServerActionCode(name string) (ServerActionCode, error) {
	switch strings.ToLower(name) {
	case CreateServerAction:
		return CreateServerActionCode, nil
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
	case ConsoleServerAction:
		return ConsoleServerActionCode, nil
	default:
		return -1, fmt.Errorf("unknown server action code: %s", name)
	}
}

func HasServerAction(list []ServerAction, search ServerAction) bool {
	for _, action := range list {
		if action == search {
			return true
		}
	}
	return false
}

func HasServerActionCode(list []ServerActionCode, search ServerActionCode) bool {
	for _, action := range list {
		if action == search {
			return true
		}
	}
	return false
}

type ServerActionCodeList []ServerActionCode

func (list ServerActionCodeList) ToServerAction() []ServerAction {
	return ToServerActionList(list)
}

func ToServerActionList(list []ServerActionCode) []ServerAction {
	toList := make([]ServerAction, len(list))
	for i, item := range list {
		toList[i] = item.ServerAction()
	}
	return toList
}

func ParseServerActionCodeList(inputList []string) ([]ServerActionCode, error) {
	list := make([]ServerActionCode, len(inputList))
	for i, item := range inputList {
		code, err := ParseServerActionCode(item)
		if err != nil {
			return nil, fmt.Errorf("ParseServerActionCodeList: failed: '%s': %w", item, err)
		}
		list[i] = code
	}
	return list, nil
}
