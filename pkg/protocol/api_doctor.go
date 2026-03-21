// Package protocol provides API parameter types for OpenClaw SDK.
//
// This file contains Doctor and Diagnostics types migrated from TypeScript: src/protocol/api-params.ts
package protocol

// ============================================================================
// Doctor / Diagnostics Types
// ============================================================================

// DoctorCheckParams parameters for doctor check.
type DoctorCheckParams struct{}

// DoctorFixParams parameters for doctor fix.
type DoctorFixParams struct {
	CheckName string `json:"checkName,omitempty"`
}

// DoctorFixResult result of doctor fix.
type DoctorFixResult struct {
	Fixed  []string `json:"fixed"`
	Failed []string `json:"failed"`
}

// DiagnosticsSnapshotParams parameters for diagnostics snapshot.
type DiagnosticsSnapshotParams struct{}

// DiagnosticsSnapshotResult result of diagnostics snapshot.
type DiagnosticsSnapshotResult struct {
	Snapshot any `json:"snapshot"`
}
