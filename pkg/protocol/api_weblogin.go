// Package protocol provides API parameter types for OpenClaw SDK.
//
// This file contains WebLogin types migrated from TypeScript: src/protocol/api-params.ts
package protocol

// ============================================================================
// WebLogin Types
// ============================================================================

// WebLoginStartParams parameters for starting web login.
type WebLoginStartParams struct {
	ReturnURL string `json:"returnUrl,omitempty"`
}

// WebLoginWaitParams parameters for waiting for web login.
type WebLoginWaitParams struct {
	Token     string `json:"token"`
	TimeoutMs int64  `json:"timeoutMs,omitempty"`
}

// WebLoginStartResult result of starting web login.
type WebLoginStartResult struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}

// WebLoginWaitResult result of waiting for web login.
type WebLoginWaitResult struct {
	Success bool   `json:"success"`
	UserID  string `json:"userId,omitempty"`
}

// WebLoginCancelParams parameters for cancelling web login.
type WebLoginCancelParams struct {
	Token string `json:"token"`
}

// WebLoginCancelResult result of cancelling web login.
type WebLoginCancelResult struct{}
