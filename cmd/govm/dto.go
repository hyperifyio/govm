// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

// ErrorDTO struct defines the structure of the body for API errors
type ErrorDTO struct {

	// Error is the error message
	Error string `json:"error"`

	// Code is the error code
	Code int `json:"code,omitempty"`
}

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

	// Permissions is permissions available to the user
	Permissions ServerPermissionDTO `json:"permissions"`
}

// ServerDTO struct defines the structure of the response DTO returned from the server
type ServerDTO struct {

	// Name is the name of the virtual server
	Name string `json:"name"`

	// Status is the status of the virtual server
	Status string `json:"status"`

	// Actions which are available to perform on the server
	Actions []string `json:"actions"`

	// Permissions is permissions available to the user
	Permissions ServerPermissionDTO `json:"permissions"`
}

// CreateServerDTO defines the structure of the request body to deploy a new server
type CreateServerDTO struct {

	// Name Optional name of the server
	Name *string `json:"name,omitempty"`
}

// ServerActionDTO defines the structure of the request body to perform an action on the server
type ServerActionDTO struct {

	// Action to perform
	Action string `json:"action"`
}

// ServerListDTO struct defines the structure of the response DTO returned from the server
type ServerListDTO struct {
	Payload     []ServerDTO         `json:"payload"`
	Permissions ServerPermissionDTO `json:"permissions"`
}

// ServerVncDTO defines an response to open a VNC console
type ServerVncDTO struct {

	// URL is the URL to bundled novnc service
	URL string `json:"url"`

	// Path is the path for websocket interface
	WS string `json:"ws"`

	// Password is the VNC password
	Password string `json:"password"`

	// Token is the VNC token
	Token string `json:"token"`
}

type ServerPermissionDTO struct {
	EnabledActions []ServerAction `json:"enabledActions"`
	CreateEnabled  bool           `json:"createEnabled"`
	DeployEnabled  bool           `json:"deployEnabled"`
	StartEnabled   bool           `json:"startEnabled"`
	StopEnabled    bool           `json:"stopEnabled"`
	RestartEnabled bool           `json:"restartEnabled"`
	DeleteEnabled  bool           `json:"deleteEnabled"`
	ConsoleEnabled bool           `json:"consoleEnabled"`
}

func NewServerPermissionDTOFromServerActionList(
	enabledActions []ServerAction,
) ServerPermissionDTO {
	return ServerPermissionDTO{
		EnabledActions: enabledActions,
		CreateEnabled:  HasServerAction(enabledActions, CreateServerAction),
		DeployEnabled:  HasServerAction(enabledActions, DeployServerAction),
		StartEnabled:   HasServerAction(enabledActions, StartServerAction),
		StopEnabled:    HasServerAction(enabledActions, StopServerAction),
		RestartEnabled: HasServerAction(enabledActions, RestartServerAction),
		DeleteEnabled:  HasServerAction(enabledActions, DeleteServerAction),
		ConsoleEnabled: HasServerAction(enabledActions, ConsoleServerAction),
	}
}

func NewServerPermissionDTOFromServerActionCodeList(
	enabledActions ServerActionCodeList,
) ServerPermissionDTO {
	return ServerPermissionDTO{
		EnabledActions: enabledActions.ToServerAction(),
		CreateEnabled:  HasServerActionCode(enabledActions, CreateServerActionCode),
		DeployEnabled:  HasServerActionCode(enabledActions, DeployServerActionCode),
		StartEnabled:   HasServerActionCode(enabledActions, StartServerActionCode),
		StopEnabled:    HasServerActionCode(enabledActions, StopServerActionCode),
		RestartEnabled: HasServerActionCode(enabledActions, RestartServerActionCode),
		DeleteEnabled:  HasServerActionCode(enabledActions, DeleteServerActionCode),
		ConsoleEnabled: HasServerActionCode(enabledActions, ConsoleServerActionCode),
	}
}
