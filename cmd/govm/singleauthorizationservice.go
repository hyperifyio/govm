// Copyright (c) 2024. Hangover Games <info@hangover.games>. All rights reserved.

package main

type SingleMemoryAuthorizationService struct {
	email    string
	password string
}

func NewSingleMemoryAuthorizationService(
	email, password string,
) *SingleMemoryAuthorizationService {
	return &SingleMemoryAuthorizationService{
		email:    email,
		password: password,
	}
}

var _ AuthorizationService = &SingleMemoryAuthorizationService{}

func (s *SingleMemoryAuthorizationService) ValidateCredentials(email, password string) (bool, error) {
	if s.email == "" || s.password == "" || email == "" || password == "" {
		return false, nil
	}
	return email == s.email && password == s.password, nil
}
