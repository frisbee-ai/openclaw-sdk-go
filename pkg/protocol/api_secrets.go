// Package protocol provides API parameter types for OpenClaw SDK.
//
// This file contains Secrets Management types migrated from TypeScript: src/protocol/api-params.ts
package protocol

// ============================================================================
// Secrets Management Types
// ============================================================================

// SecretsListParams parameters for listing secrets.
type SecretsListParams struct{}

// SecretsGetParams parameters for getting a secret.
type SecretsGetParams struct {
	Key string `json:"key"`
}

// SecretsGetResult result of getting a secret.
type SecretsGetResult struct {
	Value string `json:"value"`
}

// SecretsSetParams parameters for setting a secret.
type SecretsSetParams struct {
	Key   string `json:"key"`
	Value string `json:"value"`
}

// SecretsSetResult result of setting a secret.
type SecretsSetResult struct{}

// SecretsDeleteParams parameters for deleting a secret.
type SecretsDeleteParams struct {
	Key string `json:"key"`
}

// SecretsDeleteResult result of deleting a secret.
type SecretsDeleteResult struct{}
