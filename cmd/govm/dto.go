// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

// AuthenticateEmailDTO struct defines the structure of the body to authenticate a session
type AuthenticateEmailDTO struct {

	// Email is the email address
	Email string `json:"email"`

	// Password is the password
	Password string `json:"password,omitempty"`
}

// EmailTokenDTO struct defines the structure of the authentication session DTO returned from the server
type EmailTokenDTO struct {

	// Token is the token string
	Token string `json:"token"`

	// Email is the email address
	Email string `json:"email"`

	// Verified is a boolean which defines if this session was authenticated
	Verified bool `json:"verified,omitempty"`
}

// IndexDTO struct defines the structure of the response DTO returned from the API index
type IndexDTO struct {

	// IsAuthenticated returns true if user has been authenticated
	IsAuthenticated bool `json:"isAuthenticated"`

	// Email address of the authenticated user
	Email string `json:"email"`
}

// VirtualServerDTO struct defines the structure of the response DTO returned from the server
type VirtualServerDTO struct {

	// Name is the name of the virtual server
	Name string `json:"name"`

	// Status is the status of the virtual server
	Status string `json:"status"`
}

// CreateVirtualServerDTO defines the structure of the request body to deploy a new server
type CreateVirtualServerDTO struct {

	// Name Optional name of the virtual server
	Name *string `json:"name,omitempty"`
}

// VirtualServerListDTO struct defines the structure of the response DTO returned from the server
type VirtualServerListDTO struct {

	// Payload is
	Payload []VirtualServerDTO `json:"payload"`
}
