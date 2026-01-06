
#!/bin/bash
echo "Triggering new Split Agent execution..."
curl -X POST http://localhost:8080/instances \
  -H "Content-Type: application/json" \
  -d '{"workflow_id": "split-flow", "trigger_context": {"mode": "split"}, "policy": {"materiality": "HIGH"}}'
echo ""
