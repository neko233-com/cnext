#!/bin/bash
set -e

echo "========================================="
echo "  cnext Docker Linux Verification"
echo "========================================="

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

echo ""
echo "[1/4] Building cnext for Linux (amd64)..."
GOOS=linux GOARCH=amd64 go build -o cnext-linux ./cmd/cnext
echo "  ✓ Linux binary built"

echo ""
echo "[2/4] Creating test project..."
rm -rf test-project
mkdir -p test-project/src test-project/include

cat > test-project/cnext.toml << 'EOF'
[package]
name = "docker-test"
version = "0.1.0"
edition = "2024"

[build]
compiler = "auto"
std = "c++20"
optimization = "release"

[[build.executables]]
name = "docker-test"
sources = ["src/main.cpp"]
EOF

cat > test-project/src/main.cpp << 'EOF'
#include <iostream>
#include <cstdlib>

int main() {
    std::cout << "Hello from Docker Linux!" << std::endl;
    std::cout << "Platform: " <<
#ifdef __linux__
        "Linux"
#elif defined(_WIN32)
        "Windows"
#elif defined(__APPLE__)
        "macOS"
#else
        "Unknown"
#endif
    << std::endl;
    return 0;
}
EOF
echo "  ✓ Test project created"

echo ""
echo "[3/4] Building Docker image..."
docker build -t cnext-test .
echo "  ✓ Docker image built"

echo ""
echo "[4/4] Running tests in Docker..."
echo "  Testing with GCC..."
docker run --rm -e CC=gcc -e CXX=g++ cnext-test sh -c "cnext build && ./docker-test"
echo "  ✓ GCC test passed"

echo "  Testing with Clang..."
docker run --rm -e CC=clang -e CXX=clang++ cnext-test sh -c "cnext build && ./docker-test"
echo "  ✓ Clang test passed"

echo ""
echo "========================================="
echo "  All Docker tests passed!"
echo "========================================="

# Cleanup
echo ""
echo "Cleaning up..."
rm -f cnext-linux
rm -rf test-project
docker rmi cnext-test 2>/dev/null || true
echo "✓ Cleanup complete"
