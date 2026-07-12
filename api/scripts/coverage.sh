#!/usr/bin/env bash
#
# Runs the Go unit tests with coverage and fails when the total is below the
# threshold (default 80%).
#
# Scope of the gate:
#   * Infra / generated / wiring packages listed in IGNORE_PATHS are excluded
#     (they need integration tests against a real database or carry no
#     unit-test value).
#   * Of what remains, only packages that actually ship unit tests are gated.
#     Any remaining package without tests is reported as a warning so it never
#     hides silently — add tests or extend IGNORE_PATHS.
#
# Usage:
#   bash scripts/coverage.sh                       # gate at 80%
#   COVERAGE_THRESHOLD=90 bash scripts/coverage.sh # custom threshold
#
set -euo pipefail

THRESHOLD="${COVERAGE_THRESHOLD:-80}"
PROFILE="${COVERAGE_PROFILE:-coverage.out}"

# Package path fragments excluded from the coverage gate. Add more here as the
# codebase grows (each entry drops that package and everything under it).
IGNORE_PATHS="
social-engine/docs
social-engine/common/repositories
social-engine/common/apiErrors
"

# Build a regex that drops the ignored subtrees plus the root wiring package.
FILTER='^social-engine$'
for path in $IGNORE_PATHS; do
  FILTER="${FILTER}|^${path}(/|$)"
done

# Non-ignored packages, split by whether they ship unit tests.
TESTED="$(go list -f '{{if or .TestGoFiles .XTestGoFiles}}{{.ImportPath}}{{end}}' ./... | grep -vE "$FILTER" | grep . || true)"
UNTESTED="$(go list -f '{{if not (or .TestGoFiles .XTestGoFiles)}}{{.ImportPath}}{{end}}' ./... | grep -vE "$FILTER" | grep . || true)"

if [ -z "$TESTED" ]; then
  echo "❌ No tested packages selected for the coverage gate."
  exit 1
fi

echo "Packages under the coverage gate:"
echo "$TESTED" | sed 's/^/  - /'
if [ -n "$UNTESTED" ]; then
  echo
  echo "⚠️  Non-ignored packages without unit tests (not gated — add tests or ignore):"
  echo "$UNTESTED" | sed 's/^/  - /'
fi
echo

# shellcheck disable=SC2086 -- package paths are space-safe on purpose.
go test $TESTED -coverprofile="$PROFILE" -covermode=atomic

TOTAL="$(go tool cover -func="$PROFILE" | awk '/^total:/ {gsub(/%/,"",$3); print $3}')"
if [ -z "$TOTAL" ]; then
  echo "❌ Could not determine total coverage from $PROFILE."
  exit 1
fi

echo "----------------------------------------"
echo "Total coverage: ${TOTAL}%  (threshold: ${THRESHOLD}%)"

if awk -v total="$TOTAL" -v threshold="$THRESHOLD" 'BEGIN { exit !(total + 0 >= threshold + 0) }'; then
  echo "✅ Coverage meets the ${THRESHOLD}% threshold."
else
  echo "❌ Coverage ${TOTAL}% is below the ${THRESHOLD}% threshold."
  exit 1
fi
