# Policy Engine

Gantral uses a policy engine to evaluate "Materiality" and enforcement rules. Policies are declarative content.

## Schema Example

```yaml
apiVersion: gantral.io/v1alpha1
kind: Policy
metadata:
  name: "financial-approval-policy"
spec:
  rules:
    - name: "high-value-transaction"
      condition: "transaction.amount > 10000"
      action: "REQUIRE_APPROVAL"
      approvers:
        - "role:finance-admin"
      timeout: "24h"
    - name: "low-value-transaction"
      condition: "transaction.amount <= 10000"
      action: "AUTO_APPROVE"
```

## Concepts
- **Materiality:** The assessment of risk. High/Low/Critical.
- **Rules:** Condition -> Action mappings.
- **Approvers:** Roles or users required to sign off.
