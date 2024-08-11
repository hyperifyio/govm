// Copyright (c) 2024. Heusala Group Ltd <info@hg.fi>. All rights reserved.

package main

// VirtualServerDTO struct defines the structure of the response DTO
type VirtualServerDTO struct {

	// Name is the name of the player
	Name string `json:"name"`
}

// CreateVirtualServerDTO defines the structure of the request body to the game server
type CreateVirtualServerDTO struct {

	// Name Optional name of the player
	Name *string `json:"name,omitempty"`
}
