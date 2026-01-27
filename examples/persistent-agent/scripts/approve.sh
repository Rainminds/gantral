
#!/bin/bash
if [ -z "$1" ]; then
    echo "Usage: ./approve.sh <execution_id>"
    exit 1
fi
echo "Approving execution $1..."

# Generate Token using Python
TARGET_DIR="$(dirname "$0")/../.venv/bin/python3"
if [ ! -f "$TARGET_DIR" ]; then
    TARGET_DIR="python3"
fi

TOKEN=$($TARGET_DIR -c 'import jwt, datetime; payload={"sub":"admin-user","roles":["admin"],"iat":datetime.datetime.utcnow(),"exp":datetime.datetime.utcnow()+datetime.timedelta(minutes=10)}; print(jwt.encode(payload, "dev-secret-key", algorithm="HS256"))')

curl -X POST http://localhost:8080/instances/$1/decisions \
  -H "Authorization: Bearer $TOKEN" \
  -H "Content-Type: application/json" \
  -d "{\"type\": \"APPROVE\", \"actor_id\": \"human-operator\", \"justification\": \"cli-approval\"}"
echo ""
