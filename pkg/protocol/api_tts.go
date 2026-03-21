// Package protocol provides API parameter types for OpenClaw SDK.
//
// This file contains TTS and Voice Wake types migrated from TypeScript: src/protocol/api-params.ts
package protocol

// ============================================================================
// TTS Types
// ============================================================================

// TtsSpeakParams parameters for TTS speak.
type TtsSpeakParams struct {
	Text     string  `json:"text"`
	Voice    string  `json:"voice,omitempty"`
	Language string  `json:"language,omitempty"`
	Speed    float64 `json:"speed,omitempty"`
}

// TtsSpeakResult result of TTS speak.
type TtsSpeakResult struct {
	AudioURL   string `json:"audioUrl,omitempty"`
	DurationMs int64  `json:"durationMs,omitempty"`
}

// TtsVoicesParams parameters for TTS voices.
type TtsVoicesParams struct{}

// ============================================================================
// Voice Wake Types
// ============================================================================

// VoiceWakeStartParams parameters for starting voice wake.
type VoiceWakeStartParams struct {
	Sensitivity float64  `json:"sensitivity,omitempty"`
	Keywords    []string `json:"keywords,omitempty"`
}

// VoiceWakeStopParams parameters for stopping voice wake.
type VoiceWakeStopParams struct{}

// VoiceWakeStatusParams parameters for voice wake status.
type VoiceWakeStatusParams struct{}
