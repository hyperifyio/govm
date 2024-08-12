// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

// Make sure these are used only once per place, e.g. IDE should report 1 usage
// for each! Create another constant for each error.

const (
	EncodingFailedError          = "encoding-failed"
	SessionGenerationFailedError = "session-generation-failed"
	NotFoundError                = "not-found"
	UnauthorizedError            = "unauthorized"
	BadBodyError                 = "bad-body"
	LimitParseError              = "limit-parse-failed"
	TypeParseError               = "type-parse-failed"
)
