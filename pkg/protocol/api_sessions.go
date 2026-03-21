// Package protocol provides API parameter types for OpenClaw SDK.
//
// This file contains Sessions types migrated from TypeScript: src/protocol/api-params.ts
package protocol

// ============================================================================
// Sessions Types
// ============================================================================

// SessionsListParams parameters for listing sessions.
type SessionsListParams struct{}

// SessionsListResult result of listing sessions.
type SessionsListResult struct {
	Sessions []SessionInfo `json:"sessions"`
}

// SessionInfo represents a session in the list.
type SessionInfo struct {
	ID     string         `json:"id"`
	Status string         `json:"status,omitempty"`
	Extra  map[string]any `json:"*"` // Allows additional properties
}

// SessionsPreviewParams parameters for previewing a session.
type SessionsPreviewParams struct {
	SessionID string `json:"sessionId"`
}

// SessionsPreviewResult result of previewing a session.
type SessionsPreviewResult struct {
	Preview string `json:"preview"`
}

// SessionsResolveParams parameters for resolving a session.
type SessionsResolveParams struct {
	SessionID string `json:"sessionId"`
}

// SessionsPatchParams parameters for patching a session.
type SessionsPatchParams struct {
	SessionID string `json:"sessionId"`
	Patch     any    `json:"patch"`
}

// SessionsPatchResult result of patching a session.
type SessionsPatchResult struct{}

// SessionsResetParams parameters for resetting a session.
type SessionsResetParams struct {
	SessionID string `json:"sessionId"`
}

// SessionsDeleteParams parameters for deleting a session.
type SessionsDeleteParams struct {
	SessionID string `json:"sessionId"`
}

// SessionsCompactParams parameters for compacting sessions.
type SessionsCompactParams struct{}

// SessionsUsageParams parameters for session usage.
type SessionsUsageParams struct{}

// SessionsUsageResult result of session usage.
type SessionsUsageResult struct {
	Usage map[string]any `json:"usage"`
}
