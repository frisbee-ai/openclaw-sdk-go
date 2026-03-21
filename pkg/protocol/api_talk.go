// Package protocol provides API parameter types for OpenClaw SDK.
//
// This file contains Talk and Channels types migrated from TypeScript: src/protocol/api-params.ts
package protocol

// ============================================================================
// Talk Types
// ============================================================================

// TalkConfigParams parameters for talk config.
type TalkConfigParams struct{}

// TalkConfigResult result of talk config.
type TalkConfigResult struct {
	Enabled bool           `json:"enabled"`
	Extra   map[string]any `json:"*"` // Allows additional properties
}

// TalkModeParams parameters for setting talk mode.
type TalkModeParams struct {
	Enabled bool `json:"enabled"`
}

// TalkStartParams parameters for starting a talk.
type TalkStartParams struct {
	Language string `json:"language,omitempty"`
}

// TalkStartResult result of starting a talk.
type TalkStartResult struct {
	SessionID string `json:"sessionId"`
}

// TalkStopParams parameters for stopping a talk.
type TalkStopParams struct {
	SessionID string `json:"sessionId"`
}

// TalkStopResult result of stopping a talk.
type TalkStopResult struct{}

// ============================================================================
// Channels Types
// ============================================================================

// ChannelsStatusParams parameters for channels status.
type ChannelsStatusParams struct{}

// ChannelsStatusResult result of channels status.
type ChannelsStatusResult struct {
	Channels []any `json:"channels"`
}

// ChannelsLogoutParams parameters for logging out of a channel.
type ChannelsLogoutParams struct {
	ChannelID string `json:"channelId"`
}
