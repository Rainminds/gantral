
package gantral.policies

# Default: Execution Allowed, No Pause
default allow = true
default requires_human_approval = false
default reason = "Default allow"

# Rule 1: High Materiality requires Human Approval
requires_human_approval {
    input.workflow.materiality == "HIGH"
}

reason = "High Materiality workflow requires human approval" {
    input.workflow.materiality == "HIGH"
}
