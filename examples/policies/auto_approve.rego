
package gantral.policies

default allow = true
default requires_human_approval = true # Global default safe

# Override: Auto-approve read-only actions
requires_human_approval = false {
    input.context.operation_type == "READ_ONLY"
}

reason = "Auto-approved for Read-Only operation" {
    input.context.operation_type == "READ_ONLY"
}
