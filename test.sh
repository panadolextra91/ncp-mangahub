#!/bin/bash
set -e

echo "==================================="
echo "🚀 Running MangaHub Tests..."
echo "==================================="

# Run tests and output coverage profile
export PATH="/opt/homebrew/bin:$PATH"
# Scoping to internal and pkg modules
go test -v -coverprofile=coverage.out ./internal/... ./pkg/... ./config/...

# Exclude generated code from coverage analysis
grep -v ".pb.go" coverage.out > coverage_filtered.out

echo ""
echo "📊 Coverage Report (excluding generated code):"
echo "-----------------------------------"
go tool cover -func=coverage_filtered.out
echo "-----------------------------------"

# Check if coverage is above the required 80%
TOTAL_COVERAGE=$(go tool cover -func=coverage_filtered.out | grep total | awk '{print substr($3, 1, length($3)-1)}' | head -n 1)

echo "Total Coverage: ${TOTAL_COVERAGE}%"

# Simple bash float comparison workaround
COVERAGE_INT=$(echo $TOTAL_COVERAGE | cut -d'.' -f1)

if [ "$COVERAGE_INT" -lt 80 ]; then
  echo "❌ FAILED: Coverage ${TOTAL_COVERAGE}% is strictly below the Definition of Done (80%)."
  exit 1
else
  echo "✅ PASSED: Coverage ${TOTAL_COVERAGE}% satisfies the DoD!"
fi
