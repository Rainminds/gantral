#!/bin/bash
# Gantral Developer Demo Verification Script
# Phase 4 Audit Target

set -e

echo "--- Gantral Developer Demo Audit ---"

echo "1. Validating foundations (Docker Compose)..."
if ! docker-compose ps > /dev/null 2>&1; then
    echo "ERROR: docker-compose not found or not running"
    exit 1
fi

echo "2. Starting dependencies (simulated)..."
docker-compose up -d temporal postgres

echo "3. Verifying Gantral Server entrypoint..."
if [ ! -f "cmd/server/main.go" ]; then
    echo "ERROR: cmd/server/main.go missing"
    exit 1
fi

echo "4. Checking for cloud dependencies..."
# Simple check to ensure no hard cloud imports in core
if grep -r "aws-sdk-go" core; then
    echo "ERROR: Cloud dependencies found in core logic"
    exit 1
fi

echo "--- Audit Pass ---"
