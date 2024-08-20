// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"fmt"
	"math/rand"
)

type UserEmailList []string

// contains checks if a email address exists in the email list
func (list UserEmailList) contains(email string) bool {
	return contains(list, email)
}

// ServerModel This is the data model of the server inside the GoVM
type ServerModel struct {

	// Name the name of the server
	Name string

	// Status the status of the server
	Status ServerStatusCode

	// EnabledActions
	EnabledActions ServerActionCodeList

	Users UserEmailList
}

func NewServerModel(
	name string,
	status ServerStatusCode,
	enabledActions ServerActionCodeList,
) *ServerModel {
	if name == "" {
		name = fmt.Sprintf("Server%d", rand.Intn(90000)+10000)
	}
	return &ServerModel{
		Name:           name,
		Status:         status,
		EnabledActions: enabledActions,
	}
}

func (item *ServerModel) ToDTO() ServerDTO {
	return ServerDTO{
		Name:        item.Name,
		Status:      item.Status.String(),
		Actions:     ToStatusStringList(item.Status.GetAvailableActions(item.EnabledActions)),
		Permissions: NewServerPermissionDTOFromServerActionCodeList(item.EnabledActions),
	}
}

func ToServerListArray(
	list []*ServerModel,
) []ServerDTO {
	serverDTOList := make([]ServerDTO, len(list))
	for i, item := range list {
		serverDTOList[i] = (*item).ToDTO()
	}
	return serverDTOList
}

func ToServerListDTO(
	list []*ServerModel,
	permissions ServerPermissionDTO,
) ServerListDTO {
	return ServerListDTO{
		Payload:     ToServerListArray(list),
		Permissions: permissions,
	}
}

func ToStatusStringList(
	list []ServerActionCode,
) []string {
	ret := make([]string, len(list))
	for i, item := range list {
		ret[i] = item.String()
	}
	return ret
}

// ValidateName validates a string based on the provided rules
func ValidateName(s string) bool {
	if len(s) < 1 || len(s) > 32 {
		return false
	}
	for i, c := range s {
		if i == 0 {
			if !(c >= 'a' && c <= 'z') {
				return false
			}
		} else if i == len(s)-1 {
			if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9')) {
				return false
			}
		} else {
			if !((c >= 'a' && c <= 'z') || (c >= '0' && c <= '9') || c == '_' || c == '-') {
				return false
			}
		}
	}
	return true
}
