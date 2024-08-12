// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

const InitializingServerStatus string = "initializing"
const StoppingServerStatus string = "stopping"
const StartingServerStatus string = "starting"
const StoppedServerStatus string = "stopped"
const StartedServerStatus string = "started"
const DeletingServerStatus string = "deleting"

// ServerState This is the actual state of the game, which is encrypted to
// the Private field of VirtualServerStateDTO.
type ServerState struct {

	// Name the name of the virtual server
	Name string `json:"name"`

	// Status the status of the virtual server
	Status string `json:"status"`
}

func NewTimeNow() int64 {
	return time.Now().UnixMilli()
}

func NewServerState(
	name string,
	status string,
) *ServerState {
	if name == "" {
		name = fmt.Sprintf("Domain%d", rand.Intn(90000)+10000)
	}
	if status == "" {
		status = StoppedServerStatus
	}
	return &ServerState{
		Name:   name,
		Status: status,
	}
}

// Encrypt serializes the ServerState to JSON, then encrypts it.
func (g *ServerState) Encrypt(key []byte) (string, error) {

	jsonData, err := json.Marshal(*g)
	if err != nil {
		// Not unit tested, hard to test.
		return "", fmt.Errorf("ServerState.Encrypt: failed to stringify as json: %w", err)
	}

	// Encrypt the JSON string
	return encrypt(string(jsonData), key)
}

// DecryptVirtualServerState decrypts the encrypted string and deserializes the JSON back into a ServerState.
func DecryptVirtualServerState(encryptedData string, key []byte) (*ServerState, error) {

	// Decrypt the data to get the JSON string
	decryptedData, err := decrypt(encryptedData, key)
	if err != nil {
		return nil, fmt.Errorf("DecryptVirtualServerState: failed to decrypt json: %w", err)
	}

	var dto ServerState
	err = json.Unmarshal([]byte(decryptedData), &dto)
	if err != nil {
		// Not unit tested, hard to test.
		return nil, fmt.Errorf("DecryptVirtualServerState: failed to parse json: %w", err)
	}

	return &dto, nil
}
