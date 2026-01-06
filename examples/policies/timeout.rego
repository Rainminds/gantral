
package gantral.policies

default allow = true

# Deny if pending for too long (e.g. 1 hour = 3600000000000 ns)
deny {
    input.state == "WAITING_FOR_HUMAN"
    input.current_time_ns - input.state_start_time_ns > 3600000000000
}

reason = "Execution timed out waiting for approval" {
    input.state == "WAITING_FOR_HUMAN"
    input.current_time_ns - input.state_start_time_ns > 3600000000000
}
