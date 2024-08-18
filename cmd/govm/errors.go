// Copyright (c) 2024. Sendanor <info@sendanor.fi>. All rights reserved.

package main

// Make sure these are used only once per place, e.g. IDE should report 1 usage
// for each! Create another constant for each error.

const (
	EncodingFailedError             = "encoding-failed"
	SessionGenerationFailedError    = "session-generation-failed"
	SessionAuthorizationFailedError = "session-authorization-failed"
	NotFoundError                   = "not-found"
	InternalServerError             = "internal-server-error"
	InvalidMethod                   = "invalid-method"
	UpgradeError                    = "upgrade-error"
	VncConnectError                 = "vnc-connect-error"
	UnauthorizedError               = "unauthorized"
	VncGeneratePasswordError        = "vnc-generate-password-error"
	VncSetPasswordError             = "vnc-set-password-error"
	VncGenerateTokenError           = "vnc-generate-token-error"
	VncOpenError                    = "vnc-open-error"
	BadBodyError                    = "bad-body"
	LimitParseError                 = "limit-parse-failed"
	TypeParseError                  = "type-parse-failed"
)
