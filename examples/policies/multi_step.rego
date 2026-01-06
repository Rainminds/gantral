
package gantral.policies

default allow = true
default requires_human_approval = false
default approvers = []

# Critical Actions require dual approval
requires_human_approval {
    input.context.category == "CRITICAL"
}

approvers = ["group:engineering", "group:compliance"] {
    input.context.category == "CRITICAL"
}

reason = "Critical actions require Engineering and Compliance approval" {
    input.context.category == "CRITICAL"
}
