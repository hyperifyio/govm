// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

import (
	"net/http"
)

// Make sure these are used only once per place, e.g. IDE should report 1 usage
// for each! Create another constant for each error.

const (
	EncodingFailedError          = "encoding-failed"
	SessionGenerationFailedError = "session-generation-failed"
	UnauthorizedError            = "unauthorized"
	BadBodyError                 = "bad-body"
	LimitParseError              = "limit-parse-failed"
	TypeParseError               = "type-parse-failed"
)

func sendHttpError(w http.ResponseWriter, code string, status int) {
	recordFailedOperationMetric(code)
	// Handle error if limit query parameter is not an integer
	http.Error(w, code, status)
}
