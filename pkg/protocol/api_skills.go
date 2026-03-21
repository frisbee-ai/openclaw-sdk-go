// Package protocol provides API parameter types for OpenClaw SDK.
//
// This file contains Skills and Tools types migrated from TypeScript: src/protocol/api-params.ts
package protocol

// ============================================================================
// Skills Types
// ============================================================================

// SkillsStatusParams parameters for skills status.
type SkillsStatusParams struct {
	SkillID string `json:"skillId,omitempty"`
}

// ToolsCatalogParams parameters for tools catalog.
type ToolsCatalogParams struct{}

// ToolsCatalogResult result of tools catalog.
type ToolsCatalogResult struct {
	Tools []any `json:"tools"`
}

// SkillsBinsParams parameters for skills bins.
type SkillsBinsParams struct{}

// SkillsBinsResult result of skills bins.
type SkillsBinsResult struct {
	Bins []any `json:"bins"`
}

// SkillsInstallParams parameters for installing a skill.
type SkillsInstallParams struct {
	SkillID string `json:"skillId"`
}

// SkillsUpdateParams parameters for updating a skill.
type SkillsUpdateParams struct {
	SkillID string `json:"skillId"`
}
