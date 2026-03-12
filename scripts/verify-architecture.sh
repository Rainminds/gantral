#!/bin/bash
# verify-architecture.sh
# description: Enforces dependency boundaries for the 'examples' directory.
# The `examples/` directory must only use the public SDK or Specs to interact with Gantral.
# It is strictly prohibited for external examples to import internal engine/core/pkg implementations.

set -e

echo "Running Architecture Verification..."

# Detect if `examples/...` transitively imports internal packages.
# `go list -deps` will list all dependencies (direct and transitive).
# We search for forbidden internal paths.
VIOLATIONS=$(go list -deps ./examples/... | grep -E 'github.com/Rainminds/gantral/(core|internal|adapters|pkg)' || true)

if [ -n "$VIOLATIONS" ]; then
    echo "❌ ARCHITECTURE VIOLATION DETECTED!"
    echo "The following internal packages leaked into the 'examples' directory:"
    echo "$VIOLATIONS"
    echo ""
    echo "Rule Validation Failed: The examples directory must ONLY interact with Gantral via the SDK (github.com/Rainminds/gantral/sdk)."
    exit 1
fi

echo "✅ Architecture Verification Passed. No internal leaks detected."
exit 0
