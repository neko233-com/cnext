#!/bin/bash
set -e

# cnext - macOS/Linux installer
# Usage: curl -fsSL https://raw.githubusercontent.com/neko233-com/cnext/main/scripts/install.sh | bash
# Or:    curl -fsSL .../install.sh | bash -s -- v1.0.0

VERSION="${1:-latest}"
BINARY_NAME="cnext"
REPO="neko233-com/cnext"

detect_os() {
    case "$(uname -s)" in
        Linux*)     echo "linux" ;;
        Darwin*)    echo "darwin" ;;
        *)          echo "unsupported" ;;
    esac
}

detect_arch() {
    case "$(uname -m)" in
        x86_64|amd64)   echo "amd64" ;;
        aarch64|arm64)  echo "arm64" ;;
        *)              echo "amd64" ;;
    esac
}

normalize_version() {
    local v="$1"
    v="${v#v}"
    v="${v#V}"
    echo "$v"
}

check_compiler() {
    local has_gcc=false
    local has_clang=false

    if command -v g++ &>/dev/null || command -v gcc &>/dev/null; then
        has_gcc=true
    fi
    if command -v clang++ &>/dev/null || command -v clang &>/dev/null; then
        has_clang=true
    fi

    if $has_gcc || $has_clang; then
        return 0
    fi

    echo ""
    echo "⚠ No C/C++ compiler found!"
    echo ""
    echo "cnext requires a compiler (gcc or clang). Install one:"
    echo ""
    echo "  Ubuntu/Debian:  sudo apt install gcc g++ clang"
    echo "  CentOS/RHEL:    sudo yum install gcc gcc-c++ clang"
    echo "  macOS:          xcode-select --install"
    echo "  Arch:           sudo pacman -S gcc clang"
    echo ""
    echo "After installing, run: cnext --help"
    echo ""
    return 1
}

install_binary() {
    local os="$1"
    local arch="$2"
    local ver="$3"

    local asset="${BINARY_NAME}-${os}-${arch}"
    local url
    if [ "$ver" = "latest" ]; then
        url="https://github.com/${REPO}/releases/latest/download/${asset}"
    else
        url="https://github.com/${REPO}/releases/download/v${ver}/${asset}"
    fi
    local install_dir="/usr/local/bin"
    local target="${BINARY_NAME}"

    echo "Downloading ${url}..."
    TMPDIR=$(mktemp -d)
    curl -fsSL "$url" -o "${TMPDIR}/${target}"

    if [ -w "$install_dir" ]; then
        mv -f "${TMPDIR}/${target}" "${install_dir}/${target}"
    else
        sudo mv -f "${TMPDIR}/${target}" "${install_dir}/${target}"
    fi

    chmod +x "${install_dir}/${target}"
    rm -rf "$TMPDIR"

    echo "✓ Installed to ${install_dir}/${target}"
}

main() {
    OS=$(detect_os)
    ARCH=$(detect_arch)

    if [ "$OS" = "unsupported" ]; then
        echo "Unsupported operating system."
        echo "Windows users: run install.ps1 in PowerShell or CMD."
        echo "  irm https://raw.githubusercontent.com/${REPO}/main/scripts/install.ps1 | iex"
        exit 1
    fi

    if [ "$VERSION" != "latest" ] && [ -n "$VERSION" ]; then
        VERSION=$(normalize_version "$VERSION")
    else
        VERSION="latest"
    fi

    echo "========================================="
    echo "  cnext installer"
    echo "========================================="
    echo ""
    echo "Detected: ${OS}/${ARCH}"
    echo "Version:  ${VERSION}"
    echo ""

    install_binary "$OS" "$ARCH" "$VERSION"

    echo ""
    echo "========================================="
    echo "  Verifying installation..."
    echo "========================================="

    if command -v ${BINARY_NAME} &>/dev/null; then
        echo "✓ ${BINARY_NAME} is in PATH"
        ${BINARY_NAME} --help | head -5
    else
        echo ""
        echo "⚠ ${BINARY_NAME} installed but not in PATH."
        echo "Add to your shell profile:"
        echo ""
        echo "  export PATH=\"/usr/local/bin:\$PATH\""
        echo ""
    fi

    check_compiler || true

    echo ""
    echo "========================================="
    echo "  Quick Start"
    echo "========================================="
    echo ""
    echo "  cnext init my-app    # Create new project"
    echo "  cd my-app"
    echo "  cnext build          # Build project"
    echo "  cnext run            # Run executable"
    echo ""
}

main "$@"
