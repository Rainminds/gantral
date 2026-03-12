#!/bin/bash
set -e

THRESHOLD=60
COVERAGE_FILE="coverage.out"
RAW_FILE="coverage.raw.out"

echo "📊 Starting CI-Parity Coverage Check..."

# 1. Run Tests with Coverage
# We use -coverpkg=./... to ensure we capture cross-package coverage (e.g. integration tests hitting engine)
go test -tags=integration -coverpkg=./... -coverprofile=$RAW_FILE ./...

# 2. Filter Coverage
# Exclude infra, generated code, cmd entrypoints, and boilerplate as per CI config
grep -v -E "infra/|cmd/|pkg/logger/|pkg/config/|sdk/|github.com/Rainminds/gantral/main.go" $RAW_FILE > $COVERAGE_FILE

# 3. Calculate Total
TOTAL=$(go tool cover -func=$COVERAGE_FILE | grep total | awk '{print $3}' | sed 's/%//')

if [ -z "$TOTAL" ]; then
    echo "❌ Error: Could not calculate coverage total."
    exit 1
fi

echo "✅ Total Filtered Coverage: $TOTAL%"

# 4. Enforce Threshold
if (( $(echo "$TOTAL < $THRESHOLD" | bc -l) )); then
    echo "❌ Coverage $TOTAL% is below the required $THRESHOLD% threshold."
    exit 1
fi

echo "🚀 Coverage Check Passed!"
rm -f $RAW_FILE $COVERAGE_FILE
