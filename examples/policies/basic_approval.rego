package gantral.policies

# Standard Decision Signals (TRD §8)
# ALLOW | REQUIRE_HUMAN | DENY

default decision = "ALLOW"
default reason = "Default allow"

decision = "REQUIRE_HUMAN" {
    input.workflow.materiality == "HIGH"
}

reason = "High Materiality workflow requires human approval" {
    input.workflow.materiality == "HIGH"
}

# Mapping for backwards compatibility with previous examples if needed
allow = true
requires_human_approval { decision == "REQUIRE_HUMAN" }
deny { decision == "DENY" }
