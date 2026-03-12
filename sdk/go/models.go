package sdk

// MaterialityLevel defines the risk level of an operation.
type MaterialityLevel string

const (
	MaterialityLow    MaterialityLevel = "LOW"
	MaterialityMedium MaterialityLevel = "MEDIUM"
	MaterialityHigh   MaterialityLevel = "HIGH"
)

// Policy defines the governance rules for execution.
type Policy struct {
	ID                     string           `json:"id"`
	Materiality            MaterialityLevel `json:"materiality"`
	RequiresHumanApproval  bool             `json:"requires_human_approval,omitempty"`
	ApprovalTimeoutSeconds int64            `json:"approval_timeout_seconds,omitempty"`
}

// DecisionType defines the type of decision made.
type DecisionType string

const (
	DecisionApprove  DecisionType = "APPROVE"
	DecisionReject   DecisionType = "REJECT"
	DecisionOverride DecisionType = "OVERRIDE"
)

// State represents the state of a Gantral execution instance.
type State string

const (
	StatePending         State = "PENDING"
	StateRunning         State = "RUNNING"
	StateWaitingForHuman State = "WAITING_FOR_HUMAN"
	StateApproved        State = "APPROVED"
	StateRejected        State = "REJECTED"
	StateOverridden      State = "OVERRIDDEN"
	StateCompleted       State = "COMPLETED"
	StateFailed          State = "FAILED"
)

// Instance represents a workflow execution instance.
type Instance struct {
	ID              string                 `json:"id"`
	WorkflowID      string                 `json:"workflow_id"`
	OwningTeamID    string                 `json:"owning_team_id"`
	State           State                  `json:"state"`
	TriggerContext  map[string]interface{} `json:"trigger_context,omitempty"`
	Policy          Policy                 `json:"policy"`
	ArtifactHash    string                 `json:"artifact_hash,omitempty"`
	ArtifactVersion int                    `json:"artifact_version,omitempty"`
	CreatedAt       string                 `json:"created_at"`
	UpdatedAt       string                 `json:"updated_at"`
}

// CommitmentArtifact represents the immutable proof of an execution authority transition.
// This is a mirror of the internal model for use in SDK clients and auditors.
type CommitmentArtifact struct {
	ArtifactVersion     string        `json:"artifact_version"`
	ArtifactID          string        `json:"artifact_id"`
	InstanceID          string        `json:"instance_id"`
	WorkflowVersionID   string        `json:"workflow_version_id"`
	PolicyVersionID     string        `json:"policy_version_id"`
	AuthorityState      string        `json:"authority_state"`
	ContextSnapshotHash string        `json:"context_snapshot_hash"`
	HumanActorID        string        `json:"human_actor_id"`
	Justification       string        `json:"justification"`
	PrevArtifactHash    interface{}   `json:"prev_artifact_hash"`
	ArtifactHash        string        `json:"artifact_hash"`
	CryptoProfile       string        `json:"crypto_profile"`
	ArtifactSignature   string        `json:"artifact_signature"`
	SignatureAlgorithm  string        `json:"signature_algorithm"`
	TimestampToken      string        `json:"timestamp_token"`
	TimestampAlgorithm  string        `json:"timestamp_algorithm"`
	Attestations        []Attestation `json:"attestations"`
}

// Attestation represents a cryptographic extension to an artifact.
type Attestation struct {
	Type    string      `json:"type"`
	Payload interface{} `json:"payload"`
}
