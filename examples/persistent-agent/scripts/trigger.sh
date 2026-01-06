
#!/bin/bash
echo "Triggering new execution..."
curl -X POST http://localhost:8080/instances \
  -H "Content-Type: application/json" \
  -d '{"workflow_id": "demo-agent-flow", "trigger_context": {"tier": "prod"}, "policy": {"materiality": "HIGH"}}'
echo ""
