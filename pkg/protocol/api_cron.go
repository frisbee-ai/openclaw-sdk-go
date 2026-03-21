// Package protocol provides API parameter types for OpenClaw SDK.
//
// This file contains Cron types migrated from TypeScript: src/protocol/api-params.ts
package protocol

// ============================================================================
// Cron Types
// ============================================================================

// CronListParams parameters for listing cron jobs.
type CronListParams struct{}

// CronStatusParams parameters for cron job status.
type CronStatusParams struct {
	JobID string `json:"jobId"`
}

// CronAddParams parameters for adding a cron job.
type CronAddParams struct {
	Cron   string `json:"cron"`
	Prompt string `json:"prompt"`
}

// CronUpdateParams parameters for updating a cron job.
type CronUpdateParams struct {
	JobID  string `json:"jobId"`
	Cron   string `json:"cron,omitempty"`
	Prompt string `json:"prompt,omitempty"`
}

// CronRemoveParams parameters for removing a cron job.
type CronRemoveParams struct {
	JobID string `json:"jobId"`
}

// CronRunParams parameters for running a cron job.
type CronRunParams struct {
	JobID string `json:"jobId"`
}

// CronRunsParams parameters for listing cron runs.
type CronRunsParams struct{}
