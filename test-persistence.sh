#!/bin/bash
set -e

echo "1. Create instance..."
# Capture ID from response (assuming JSON response has instance_id)
RESPONSE=$(curl -s -X POST http://localhost:8080/instances \
  -H "Content-Type: application/json" \
  -d '{"workflow_id": "test", "trigger_context": {}}')

echo "Response: $RESPONSE"
ID=$(echo $RESPONSE | jq -r .instance_id)

if [ "$ID" == "null" ] || [ -z "$ID" ]; then
    echo "Failed to capture Instance ID"
    exit 1
fi

echo "Created Instance ID: $ID"

echo "2. Restart services..."
make dev-down
sleep 2
make dev &
SERVER_PID=$!
sleep 5 # Wait for server to come up

echo "3. Verify instance survived..."
curl -s http://localhost:8080/instances/$ID | jq .

# Cleanup
kill $SERVER_PID
