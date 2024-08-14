// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"fmt"
	"math/rand"
)

// ServerModel This is the data model of the server inside the GoVM
type ServerModel struct {

	// Name the name of the server
	Name string

	// Status the status of the server
	Status ServerStatusCode
}

func NewServerModel(
	name string,
	status ServerStatusCode,
) *ServerModel {
	if name == "" {
		name = fmt.Sprintf("Server%d", rand.Intn(90000)+10000)
	}
	return &ServerModel{
		Name:   name,
		Status: status,
	}
}

func (item *ServerModel) ToDTO() ServerDTO {
	return ServerDTO{
		Name:    item.Name,
		Status:  item.Status.String(),
		Actions: ToStatusStringList(item.Status.GetAvailableActions()),
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
