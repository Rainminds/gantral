
#!/bin/bash
if [ -z "$1" ]; then
    echo "Usage: ./status.sh <execution_id>"
    exit 1
fi
curl -X GET http://localhost:8080/instances/$1
echo ""
