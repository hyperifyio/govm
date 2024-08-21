// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

type Session struct {
	Token string
	Email string
}

func NewSession(token, email string) *Session {
	return &Session{Token: token, Email: email}
}

type SessionService interface {

	// CreateSession creates a new session
	CreateSession(email string) (*Session, error)

	// ValidateSession validates a token and returns the session if it is valid, otherwise nil
	ValidateSession(token string) (*Session, error)

	// DeleteSession invalidates the session, otherwise an error
	DeleteSession(session *Session) error
}
