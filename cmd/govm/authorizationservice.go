// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

type AuthorizationService interface {

	// ValidateCredentials returns true if correct credentials has been provided
	ValidateCredentials(email, password string) (bool, error)
}
