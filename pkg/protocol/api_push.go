// Package protocol provides API parameter types for OpenClaw SDK.
//
// This file contains Push and Usage types migrated from TypeScript: src/protocol/api-params.ts
package protocol

// ============================================================================
// Push Notification Types
// ============================================================================

// PushRegisterParams parameters for registering push notification.
type PushRegisterParams struct {
	Token    string `json:"token"`
	Platform string `json:"platform"`
}

// PushRegisterResult result of registering push notification.
type PushRegisterResult struct{}

// PushUnregisterParams parameters for unregistering push notification.
type PushUnregisterParams struct {
	Token string `json:"token"`
}

// PushUnregisterResult result of unregistering push notification.
type PushUnregisterResult struct{}

// PushSendParams parameters for sending push notification.
type PushSendParams struct {
	Target string         `json:"target"`
	Title  string         `json:"title"`
	Body   string         `json:"body"`
	Data   map[string]any `json:"data,omitempty"`
}

// PushSendResult result of sending push notification.
type PushSendResult struct{}

// ============================================================================
// Usage / Billing Types
// ============================================================================

// UsageSummaryParams parameters for usage summary.
type UsageSummaryParams struct {
	Period string `json:"period,omitempty"`
}

// UsageSummaryResult result of usage summary.
type UsageSummaryResult struct {
	TotalTokens *int64         `json:"totalTokens,omitempty"`
	TotalCost   *float64       `json:"totalCost,omitempty"`
	Period      string         `json:"period,omitempty"`
	Extra       map[string]any `json:"*"` // Allows additional properties
}

// UsageDetailsParams parameters for usage details.
type UsageDetailsParams struct {
	Period  string `json:"period,omitempty"`
	AgentID string `json:"agentId,omitempty"`
}

// UsageDetailsResult result of usage details.
type UsageDetailsResult struct {
	Entries []any `json:"entries"`
}
