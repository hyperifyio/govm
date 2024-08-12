// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"fmt"
	"math/rand"
)

// ServerState This is the actual state of the game, which is encrypted to
// the Private field of VirtualServerStateDTO.
type ServerState struct {

	// Name the name of the server
	Name string

	// Status the status of the server
	Status ServerStatusCode
}

func NewServerState(
	name string,
	status ServerStatusCode,
) *ServerState {
	if name == "" {
		name = fmt.Sprintf("Domain%d", rand.Intn(90000)+10000)
	}
	return &ServerState{
		Name:   name,
		Status: status,
	}
}

func (item *ServerState) ToDTO() ServerDTO {
	return ServerDTO{
		Name:    item.Name,
		Status:  item.Status.String(),
		Actions: ToStatusStringList(item.Status.GetAvailableActions()),
	}
}

func ToServerListArray(
	list []*ServerState,
) []ServerDTO {
	serverDTOList := make([]ServerDTO, len(list))
	for i, item := range list {
		serverDTOList[i] = (*item).ToDTO()
	}
	return serverDTOList
}

func ToServerListDTO(
	list []*ServerState,
) ServerListDTO {
	return ServerListDTO{
		Payload: ToServerListArray(list),
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

func FindServerStateByName(states []*ServerState, name string) (*ServerState, bool) {
	for _, state := range states {
		if state.Name == name {
			return state, true
		}
	}
	// Return an empty ServerState and false if not found
	return nil, false
}
