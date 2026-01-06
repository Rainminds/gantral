
#!/bin/bash
if [ -z "$1" ]; then
    echo "Usage: ./approve.sh <execution_id>"
    exit 1
fi
echo "Approving execution $1..."
curl -X POST http://localhost:8080/instances/$1/decisions \
  -H "Content-Type: application/json" \
  -d "{\"type\": \"APPROVE\", \"actor_id\": \"human-operator\", \"justification\": \"cli-approval\"}"
echo ""
