#!/bin/bash
set -e

echo "========================================="
echo "  cnext Cross-Platform Test Suite"
echo "========================================="

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

echo ""
echo "[1/4] Building cnext..."
go build -o cnext ./cmd/cnext
echo "  ✓ cnext built"

echo ""
echo "[2/4] Running Go unit tests..."
go test ./internal/... -v -count=1
echo "  ✓ Unit tests passed"

echo ""
echo "[3/4] Running Go integration tests..."
go test ./internal/e2e/ -v -count=1
echo "  ✓ Integration tests passed"

echo ""
echo "[4/4] Running Docker Linux verification..."
bash "$SCRIPT_DIR/test-docker.sh"
echo "  ✓ Docker tests passed"

echo ""
echo "========================================="
echo "  All tests passed!"
echo "========================================="

# Cleanup
rm -f cnext
