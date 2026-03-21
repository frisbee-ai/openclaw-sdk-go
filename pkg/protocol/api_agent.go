// Package protocol provides API parameter types for OpenClaw SDK.
//
// This file contains Agent, Node Pairing, and Device Pairing types
// migrated from TypeScript: src/protocol/api-params.ts
package protocol

// ============================================================================
// Agent Types
// ============================================================================

// AgentIdentityParams parameters for agent identity verification.
type AgentIdentityParams struct {
	AgentID string `json:"agentId"`
}

// AgentIdentityResult result of agent identity verification.
type AgentIdentityResult struct {
	ID      string        `json:"id"`
	Summary *AgentSummary `json:"summary,omitempty"`
}

// AgentWaitParams parameters for waiting on agent.
type AgentWaitParams struct {
	AgentID   string `json:"agentId"`
	TimeoutMs int64  `json:"timeoutMs,omitempty"`
}

// AgentsFileEntry represents a file entry for agent file operations.
type AgentsFileEntry struct {
	Path    string `json:"path"`
	Content string `json:"content"`
}

// AgentsCreateParams parameters for creating an agent.
type AgentsCreateParams struct {
	AgentID string            `json:"agentId"`
	Files   []AgentsFileEntry `json:"files"`
}

// AgentsCreateResult result of creating an agent.
type AgentsCreateResult struct {
	AgentID string `json:"agentId"`
}

// AgentsUpdateParams parameters for updating an agent.
type AgentsUpdateParams struct {
	AgentID string            `json:"agentId"`
	Files   []AgentsFileEntry `json:"files"`
}

// AgentsUpdateResult result of updating an agent.
type AgentsUpdateResult struct {
	AgentID string `json:"agentId"`
}

// AgentsDeleteParams parameters for deleting an agent.
type AgentsDeleteParams struct {
	AgentID string `json:"agentId"`
}

// AgentsDeleteResult result of deleting an agent.
type AgentsDeleteResult struct {
	AgentID string `json:"agentId"`
}

// AgentsFilesListParams parameters for listing agent files.
type AgentsFilesListParams struct {
	AgentID string `json:"agentId"`
}

// AgentsFilesListResult result of listing agent files.
type AgentsFilesListResult struct {
	Files []string `json:"files"`
}

// AgentsFilesGetParams parameters for getting an agent file.
type AgentsFilesGetParams struct {
	AgentID string `json:"agentId"`
	Path    string `json:"path"`
}

// AgentsFilesGetResult result of getting an agent file.
type AgentsFilesGetResult struct {
	Content string `json:"content"`
}

// AgentsFilesSetParams parameters for setting an agent file.
type AgentsFilesSetParams struct {
	AgentID string `json:"agentId"`
	Path    string `json:"path"`
	Content string `json:"content"`
}

// AgentsFilesSetResult result of setting an agent file.
type AgentsFilesSetResult struct{}

// AgentsListParams parameters for listing agents.
type AgentsListParams struct{}

// AgentsListResult result of listing agents.
type AgentsListResult struct {
	Agents []AgentSummary `json:"agents"`
}

// ============================================================================
// Node Pairing Types
// ============================================================================

// NodePairRequestParams parameters for requesting node pairing.
type NodePairRequestParams struct {
	NodeID string `json:"nodeId"`
	TtlSec int64  `json:"ttlSec,omitempty"`
}

// NodePairListParams parameters for listing node pairings.
type NodePairListParams struct {
	NodeID string `json:"nodeId"`
}

// NodePairApproveParams parameters for approving node pairing.
type NodePairApproveParams struct {
	NodeID    string `json:"nodeId"`
	PairingID string `json:"pairingId"`
}

// NodePairRejectParams parameters for rejecting node pairing.
type NodePairRejectParams struct {
	NodeID    string `json:"nodeId"`
	PairingID string `json:"pairingId"`
}

// NodePairVerifyParams parameters for verifying node pairing.
type NodePairVerifyParams struct {
	NodeID    string `json:"nodeId"`
	PairingID string `json:"pairingId"`
	Code      string `json:"code"`
}

// ============================================================================
// Device Pairing Types
// ============================================================================

// DevicePairListParams parameters for listing device pairings.
type DevicePairListParams struct{}

// DevicePairApproveParams parameters for approving device pairing.
type DevicePairApproveParams struct {
	PairingID string `json:"pairingId"`
}

// DevicePairRejectParams parameters for rejecting device pairing.
type DevicePairRejectParams struct {
	PairingID string `json:"pairingId"`
}
