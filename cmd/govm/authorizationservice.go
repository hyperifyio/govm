// Copyright (c) 2024. Hangover Games <info@hangover.games>. All rights reserved.

package main

type AuthorizationService interface {

	// ValidateCredentials returns true if correct credentials has been provided
	ValidateCredentials(email, password string) (bool, error)
}
