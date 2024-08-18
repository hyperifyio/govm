// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"errors"
	"fmt"
)

type SingleMemorySessionService struct {
	authenticatedSessionEmail string
	authenticatedSessionToken string
}

func NewSingleMemorySessionService() *SingleMemorySessionService {
	return &SingleMemorySessionService{
		authenticatedSessionToken: "",
		authenticatedSessionEmail: "",
	}
}

var _ SessionService = &SingleMemorySessionService{}

func (s *SingleMemorySessionService) CreateSession(email string) (*Session, error) {
	token, err := generateAuthToken()
	if err != nil {
		return nil, fmt.Errorf("CreateSession: generating session: error: %v", err)
	}
	s.authenticatedSessionEmail = email
	s.authenticatedSessionToken = token
	return NewSession(token, email), nil
}

func (s *SingleMemorySessionService) ValidateSession(token string) (*Session, error) {
	if s.authenticatedSessionToken == "" {
		return nil, errors.New("no token provided")
	}
	if s.authenticatedSessionToken != token {
		return nil, errors.New("token invalid")
	}
	return NewSession(token, s.authenticatedSessionEmail), nil
}
