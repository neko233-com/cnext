#!/bin/bash
set -e

# cnext Cross-Platform Docker Verification
# Usage: bash docker/verify-all.sh

SCRIPT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
PROJECT_DIR="$(dirname "$SCRIPT_DIR")"

cd "$PROJECT_DIR"

echo "========================================="
echo "  cnext Cross-Platform Verification"
echo "========================================="

# Step 1: Build cnext for all platforms
echo ""
echo "[1/5] Building cnext binaries..."
GOOS=linux GOARCH=amd64 go build -o docker/cnext-linux-amd64 ./cmd/cnext
GOOS=linux GOARCH=arm64 go build -o docker/cnext-linux-arm64 ./cmd/cnext
echo "  ✓ Linux binaries built"

# Step 2: Create test project
echo ""
echo "[2/5] Creating test project..."
rm -rf docker/test-project
mkdir -p docker/test-project/src docker/test-project/include

cat > docker/test-project/cnext.toml << 'EOF'
[package]
name = "platform-test"
version = "1.0.0"
edition = "2024"

[build]
compiler = "auto"
std = "c++20"
optimization = "release"

[[build.executables]]
name = "platform-test"
sources = ["src/main.cpp"]
EOF

cat > docker/test-project/src/main.cpp << 'EOF'
#include <iostream>
#include <string>

std::string getPlatform() {
#ifdef __linux__
    #ifdef __x86_64__
        return "Linux x86_64";
    #elif __aarch64__
        return "Linux ARM64";
    #else
        return "Linux Unknown";
    #endif
#elif defined(__APPLE__)
    #ifdef __x86_64__
        return "macOS x86_64";
    #elif __aarch64__
        return "macOS ARM64";
    #else
        return "macOS Unknown";
    #endif
#elif defined(_WIN32)
    #ifdef _M_X64
        return "Windows x86_64";
    #elif _M_ARM64
        return "Windows ARM64";
    #else
        return "Windows Unknown";
    #endif
#else
    return "Unknown Platform";
#endif
}

int main() {
    std::cout << "Platform: " << getPlatform() << std::endl;
    std::cout << "Compiler: " <<
#ifdef __GNUC__
        "GCC " + std::to_string(__GNUC__) + "." + std::to_string(__GNUC_MINOR__)
#elif defined(__clang__)
        "Clang " + std::to_string(__clang_major__) + "." + std::to_string(__clang_minor__)
#elif defined(_MSC_VER)
        "MSVC " + std::to_string(_MSC_VER)
#else
        "Unknown"
#endif
    << std::endl;
    return 0;
}
EOF
echo "  ✓ Test project created"

# Step 3: Test Linux AMD64
echo ""
echo "[3/5] Testing Linux AMD64 (Docker)..."
docker build -f docker/Dockerfile.linux-amd64 -t cnext-test-amd64 docker/
docker run --rm cnext-test-amd64 sh -c "cnext build && ./platform-test"
echo "  ✓ Linux AMD64 passed"
docker rmi cnext-test-amd64

# Step 4: Test Linux ARM64 (if QEMU available)
echo ""
echo "[4/5] Testing Linux ARM64 (Docker)..."
if docker buildx ls 2>/dev/null | grep -q "linux/arm64"; then
    docker buildx build --platform linux/arm64 -f docker/Dockerfile.linux-arm64 -t cnext-test-arm64 docker/ --load
    docker run --rm --platform linux/arm64 cnext-test-arm64 sh -c "cnext build && ./platform-test"
    echo "  ✓ Linux ARM64 passed"
    docker rmi cnext-test-arm64
else
    echo "  ⚠ ARM64 emulation not available, skipping"
fi

# Step 5: Test Windows cross-compile
echo ""
echo "[5/5] Testing Windows cross-compile..."
GOOS=windows GOARCH=amd64 go build -o docker/cnext-windows-amd64.exe ./cmd/cnext
echo "  ✓ Windows binary built"

# Cleanup
echo ""
echo "Cleaning up..."
rm -f docker/cnext-linux-amd64 docker/cnext-linux-arm64 docker/cnext-windows-amd64.exe
rm -rf docker/test-project

echo ""
echo "========================================="
echo "  All platform verifications passed!"
echo "========================================="
