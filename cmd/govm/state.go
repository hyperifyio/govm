// Copyright (c) 2024. Heusala Group Ltd <info@hg.fi>. All rights reserved.

package main

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"time"
)

// VirtualServerState This is the actual state of the game, which is encrypted to
// the Private field of VirtualServerStateDTO.
type VirtualServerState struct {

	// Name the name of the player
	Name string `json:"name"`

	// Created the time when game was started
	Created int64 `json:"created"`

	// Updated the time when game state was updated
	Updated int64 `json:"updated"`
}

func NewTimeNow() int64 {
	return time.Now().UnixMilli()
}

func NewVirtualServerState(
	now int64,
) *VirtualServerState {
	return &VirtualServerState{
		Name:    fmt.Sprintf("Guest%d", rand.Intn(90000)+10000),
		Created: now,
		Updated: now,
	}
}

// Encrypt serializes the VirtualServerState to JSON, then encrypts it.
func (g *VirtualServerState) Encrypt(key []byte) (string, error) {

	jsonData, err := json.Marshal(*g)
	if err != nil {
		// Not unit tested, hard to test.
		return "", fmt.Errorf("VirtualServerState.Encrypt: failed to stringify as json: %w", err)
	}

	// Encrypt the JSON string
	return encrypt(string(jsonData), key)
}

// DecryptVirtualServerState decrypts the encrypted string and deserializes the JSON back into a VirtualServerState.
func DecryptVirtualServerState(encryptedData string, key []byte) (*VirtualServerState, error) {

	// Decrypt the data to get the JSON string
	decryptedData, err := decrypt(encryptedData, key)
	if err != nil {
		return nil, fmt.Errorf("DecryptVirtualServerState: failed to decrypt json: %w", err)
	}

	var dto VirtualServerState
	err = json.Unmarshal([]byte(decryptedData), &dto)
	if err != nil {
		// Not unit tested, hard to test.
		return nil, fmt.Errorf("DecryptVirtualServerState: failed to parse json: %w", err)
	}

	return &dto, nil
}
