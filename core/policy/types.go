package policy

// MaterialityLevel defines the risk level of an operation.
type MaterialityLevel string

const (
	MaterialityLow    MaterialityLevel = "LOW"
	MaterialityMedium MaterialityLevel = "MEDIUM"
	MaterialityHigh   MaterialityLevel = "HIGH"
)

// Policy defines the governance rules for execution.
type Policy struct {
	ID                    string           `json:"id"`
	Materiality           MaterialityLevel `json:"materiality"`
	RequiresHumanApproval bool             `json:"requires_human_approval"`
}

// EvaluationResult captures the decision made by the Policy Engine.
type EvaluationResult struct {
	ShouldPause bool   `json:"should_pause"`
	NextState   string `json:"next_state"` // e.g., "RUNNING", "WAITING_FOR_HUMAN"
	Reason      string `json:"reason"`
}
