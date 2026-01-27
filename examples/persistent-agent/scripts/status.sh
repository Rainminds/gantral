
#!/bin/bash
if [ -z "$1" ]; then
    echo "Usage: ./status.sh <execution_id>"
    exit 1
fi

# Generate Token using Python (same logic as trigger.py)
TARGET_DIR="$(dirname "$0")/../.venv/bin/python3"
if [ ! -f "$TARGET_DIR" ]; then
    TARGET_DIR="python3"
fi

TOKEN=$($TARGET_DIR -c 'import jwt, datetime; payload={"sub":"admin-user","roles":["admin"],"iat":datetime.datetime.utcnow(),"exp":datetime.datetime.utcnow()+datetime.timedelta(minutes=10)}; print(jwt.encode(payload, "dev-secret-key", algorithm="HS256"))')

curl -X GET http://localhost:8080/instances/$1 \
  -H "Authorization: Bearer $TOKEN"
echo ""
