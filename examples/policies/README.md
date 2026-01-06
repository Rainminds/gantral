
# Gantral Standard Policy Library

This directory contains `rego` policies that act as blueprints for common governance patterns.

## Available Policies

### 1. Basic Approval (`basic_approval.rego`)
**Use Case:** The "Hello World" of governance.
- **Logic:** Requires Human Approval if `input.workflow.materiality` is "HIGH". Otherwise auto-approves.

### 2. Multi-Step / Multi-Party (`multi_step.rego`)
**Use Case:** Critical infrastructure changes.
- **Logic:** If `input.context.category` is "CRITICAL", it demands approval from *both* "group:engineering" and "group:compliance".

### 3. Timeout (`timeout.rego`)
**Use Case:** Prevent stale executions from hanging forever.
- **Logic:** Denies execution if it has been waiting closer to 1 hour. (Requires `current_time_ns` injection).

### 4. Auto-Approve Read-Only (`auto_approve.rego`)
**Use Case:** Cost optimization / Latency reduction for safe ops.
- **Logic:** Sets `requires_human_approval = false` if `input.context.operation_type` is "READ_ONLY".

## How to Test (Using OPA CLI)

You can verify these policies locally using the `opa` CLI tool.

**Example Input (`input.json`):**
```json
{
    "workflow": {
        "materiality": "HIGH"
    },
    "context": {
        "category": "CRITICAL",
        "operation_type": "WRITE"
    }
}
```

**Run Check:**
```bash
opa eval -i input.json -d basic_approval.rego "data.gantral.policies"
```

**Expected Output:**
```json
{
  "result": [
    {
      "expressions": [
        {
          "value": {
            "allow": true,
            "requires_human_approval": true,
            "reason": "High Materiality workflow requires human approval"
          },
          "text": "data.gantral.policies",
          "location": {
            "row": 1,
            "col": 1
          }
        }
      ]
    }
  ]
}
```
