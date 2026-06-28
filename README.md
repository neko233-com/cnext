# cnext

[![CI](https://github.com/neko233-com/cnext/actions/workflows/ci.yaml/badge.svg)](https://github.com/neko233-com/cnext/actions/workflows/ci.yaml)
[![Release](https://img.shields.io/github/v/release/neko233-com/cnext)](https://github.com/neko233-com/cnext/releases)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)

**一行命令装好 C/C++ 开发环境。** 自动检测/安装编译器，跨平台构建，增量编译，依赖管理。

## 安装

### 一键安装 (推荐)

**macOS / Linux:**

```bash
curl -fsSL https://raw.githubusercontent.com/neko233-com/cnext/main/scripts/install.sh | bash
```

**Windows (PowerShell):**

```powershell
irm https://raw.githubusercontent.com/neko233-com/cnext/main/scripts/install.ps1 | iex
```

**Windows (CMD):**

```cmd
powershell -NoProfile -ExecutionPolicy Bypass -Command "irm https://raw.githubusercontent.com/neko233-com/cnext/main/scripts/install.ps1 | iex"
```

### 指定版本安装

```bash
curl -fsSL https://raw.githubusercontent.com/neko233-com/cnext/main/scripts/install.sh | bash -s -- v0.1.0
```

### 从源码编译

```bash
# 需要 Go 1.25+
go install github.com/neko233-com/cnext/cmd/cnext@latest
```

### 预编译二进制

从 [GitHub Releases](https://github.com/neko233-com/cnext/releases) 下载

## 快速开始

### One-Click Install (Recommended)

**macOS / Linux** (`install.sh`)

```bash
curl -fsSL https://raw.githubusercontent.com/neko233-com/cnext/main/scripts/install.sh | bash
```

**With specific version**

```bash
curl -fsSL https://raw.githubusercontent.com/neko233-com/cnext/main/scripts/install.sh | bash -s -- v1.0.0
```

**Windows** (`install.ps1` — do not use `.sh` on Windows)

PowerShell:

```powershell
irm https://raw.githubusercontent.com/neko233-com/cnext/main/scripts/install.ps1 | iex
```

CMD:

```cmd
powershell -NoProfile -ExecutionPolicy Bypass -Command "irm https://raw.githubusercontent.com/neko233-com/cnext/main/scripts/install.ps1 | iex"
```

**With specific version (Windows)**

```powershell
irm https://raw.githubusercontent.com/neko233-com/cnext/main/scripts/install.ps1 -OutFile $env:TEMP\cnext-install.ps1
& $env:TEMP\cnext-install.ps1 v1.0.0
```

### From Source

```bash
# Requires Go 1.25+
go install github.com/neko233-com/cnext/cmd/cnext@latest

# Or build from source
git clone https://github.com/neko233-com/cnext.git
cd cnext
go build -o cnext ./cmd/cnext
```

### Pre-built Binaries

Download from [GitHub Releases](https://github.com/neko233-com/cnext/releases)

## Quick Start

```bash
# Create a new project
cnext init my-app
cd my-app

# Build the project
cnext build

# Run it
cnext run

# Run tests
cnext test
```

## Features

- **One-click project initialization** — scaffold a complete project in seconds
- **Declarative builds** — define targets, dependencies, and options in `cnext.toml`
- **Package manager** — add, remove, and update C/C++ dependencies
- **Built-in test framework** — discover and run tests with `cnext test`
- **CMake bridging** — integrate CMake-based dependencies seamlessly
- **Cross-compilation** — target different OS and architectures from a single config
- **Multi-compiler support** — auto-detect or explicitly choose gcc, clang, or MSVC
- **Toolchain management** — install, switch, and manage compiler toolchains

## Configuration

cnext uses a TOML configuration file (`cnext.toml`) at the project root:

```toml
[package]
name = "my-app"
version = "0.1.0"
edition = "2024"

[build]
compiler = "auto"
std = "c++20"
optimization = "release"
lto = true
pic = true

[build.target]
os = "auto"
arch = "auto"

[[build.executables]]
name = "my-app"
sources = ["src/**/*.cpp"]
libraries = ["mylib"]

[[build.libraries]]
name = "mylib"
type = "static"
sources = ["lib/**/*.cpp"]

[[build.tests]]
name = "unit-tests"
sources = ["tests/**/*.cpp"]
framework = "cnext-test"

[dependencies]
fmt = "10.2.1"
yaml-cpp = { version = "0.8.0", features = ["shared"] }

[dev-dependencies]
gtest = "1.14.0"
```

## CLI Commands

| Command | Description |
|---------|-------------|
| `cnext init [name]` | Initialize a new project |
| `cnext build` | Build all targets |
| `cnext run` | Build and run the default executable |
| `cnext test` | Run project tests |
| `cnext add <pkg>` | Add a dependency |
| `cnext remove <pkg>` | Remove a dependency |
| `cnext install` | Install all dependencies |
| `cnext update` | Update dependencies |
| `cnext fmt` | Format source code |
| `cnext lint` | Lint source code |
| `cnext clean` | Remove build artifacts |
| `cnext info` | Show project information |
| `cnext doc` | Generate documentation |
| `cnext publish` | Publish package to registry |
| `cnext toolchain list` | List installed toolchains |
| `cnext toolchain install <name>` | Install a toolchain |

## Testing

cnext includes a built-in test framework with rich assertions:

```cpp
#include <cnext/test.h>
#include <cnext/assert.h>

TEST(MathSuite, Addition) {
    EXPECT_EQ(1 + 1, 2);
}

TEST(MathSuite, Comparison) {
    EXPECT_GT(5, 3);
    EXPECT_NEAR(3.14, 3.14159, 0.01);
}

int main() {
    return cnext::test::TestRunner().Run();
}
```

Run tests:

```bash
cnext test
cnext test --filter "MathSuite"
cnext test --verbose
```

## Build Options

```bash
cnext build --release          # Optimized build
cnext build --debug            # Debug symbols
cnext build --compiler clang   # Use specific compiler
cnext build --target x86       # Cross-compile
cnext build --jobs 8           # Parallel builds
```

## Project Structure

```
my-app/
├── cnext.toml          # Project configuration
├── src/                # Source files
│   └── main.cpp
├── include/            # Public headers
├── tests/              # Test files
├── build/              # Build artifacts (generated)
└── .gitignore
```

## Cross-Platform Support

| Platform | amd64 | arm64 |
|----------|-------|-------|
| Windows  | ✅    | N/A   |
| Linux    | ✅    | ✅    |
| macOS    | ✅    | ✅    |

## License

MIT License — see [LICENSE](LICENSE) for details.

## Contributing

Contributions are welcome! Please open an issue or submit a PR.

## Links

- GitHub: [https://github.com/neko233-com/cnext](https://github.com/neko233-com/cnext)
- Issues: [https://github.com/neko233-com/cnext/issues](https://github.com/neko233-com/cnext/issues)
