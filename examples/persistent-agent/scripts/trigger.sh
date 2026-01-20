#!/bin/bash
echo "Triggering new execution via authenticated script..."
# Ensure we use the venv python if it exists
if [ -f "$(dirname "$0")/../.venv/bin/python3" ]; then
    PYTHON_CMD="$(dirname "$0")/../.venv/bin/python3"
else
    PYTHON_CMD="python3"
fi

$PYTHON_CMD $(dirname "$0")/trigger.py
