// Package protocol provides API parameter types for OpenClaw SDK.
//
// This file contains Logs and ExecApprovals types migrated from TypeScript: src/protocol/api-params.ts
package protocol

// ============================================================================
// Logs Types
// ============================================================================

// LogsTailParams parameters for tailing logs.
type LogsTailParams struct {
	Lines int64 `json:"lines,omitempty"`
}

// LogsTailResult result of tailing logs.
type LogsTailResult struct {
	Logs []string `json:"logs"`
}

// ============================================================================
// ExecApprovals Types
// ============================================================================

// ExecApprovalsGetParams parameters for getting exec approvals.
type ExecApprovalsGetParams struct{}

// ExecApprovalsSetParams parameters for setting exec approvals.
type ExecApprovalsSetParams struct {
	Enabled bool `json:"enabled"`
}

// ExecApprovalsSnapshot represents exec approvals snapshot.
type ExecApprovalsSnapshot struct {
	Approvals []any `json:"approvals"`
}
