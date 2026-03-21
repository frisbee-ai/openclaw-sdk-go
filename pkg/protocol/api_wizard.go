// Package protocol provides API parameter types for OpenClaw SDK.
//
// This file contains Wizard types migrated from TypeScript: src/protocol/api-params.ts
package protocol

// ============================================================================
// Wizard Types
// ============================================================================

// WizardStartParams parameters for starting a wizard.
type WizardStartParams struct {
	WizardID string `json:"wizardId"`
	Input    any    `json:"input,omitempty"`
}

// WizardNextParams parameters for wizard next step.
type WizardNextParams struct {
	WizardID string `json:"wizardId"`
	Input    any    `json:"input,omitempty"`
}

// WizardCancelParams parameters for cancelling a wizard.
type WizardCancelParams struct {
	WizardID string `json:"wizardId"`
}

// WizardStatusParams parameters for wizard status.
type WizardStatusParams struct {
	WizardID string `json:"wizardId"`
}

// WizardNextResult result of wizard next step.
type WizardNextResult struct {
	Step     WizardStep `json:"step"`
	Complete bool       `json:"complete"`
}

// WizardStartResult result of wizard start.
type WizardStartResult = WizardNextResult

// WizardStatusResult result of wizard status.
type WizardStatusResult struct {
	CurrentStep WizardStep `json:"currentStep"`
	Complete    bool       `json:"complete"`
}
