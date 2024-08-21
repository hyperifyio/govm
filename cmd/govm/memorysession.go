// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"errors"
	"fmt"
)

type MemorySessionService struct {
	sessions map[string]string
}

func NewMemorySessionService() *MemorySessionService {
	return &MemorySessionService{
		sessions: make(map[string]string),
	}
}

var _ SessionService = &MemorySessionService{}

func (s *MemorySessionService) CreateSession(email string) (*Session, error) {
	token, err := generateAuthToken()
	if err != nil {
		return nil, fmt.Errorf("CreateSession: generating session: error: %v", err)
	}
	s.sessions[token] = email
	return NewSession(token, email), nil
}

func (s *MemorySessionService) ValidateSession(token string) (*Session, error) {
	email, exists := s.sessions[token]
	if exists {
		return NewSession(token, email), nil
	}
	return nil, errors.New("token invalid")
}

func (s *MemorySessionService) DeleteSession(session *Session) error {
	token := session.Token
	_, exists := s.sessions[token]
	if exists {
		delete(s.sessions, token)
	} else {
		return errors.New("no session")
	}
	return nil
}
