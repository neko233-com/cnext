# Getting Started with cnext

## Installation

### One-Click Install

**macOS / Linux:**

```bash
curl -fsSL https://raw.githubusercontent.com/neko233-com/cnext/main/scripts/install.sh | bash
```

**Windows (PowerShell):**

```powershell
irm https://raw.githubusercontent.com/neko233-com/cnext/main/scripts/install.ps1 | iex
```

### From Source

```bash
go install github.com/neko233-com/cnext/cmd/cnext@latest
```

## Create Your First Project

```bash
# Initialize a new project
cnext init my-app
cd my-app

# The project structure:
# my-app/
# ├── cnext.toml      # Configuration
# ├── src/main.cpp    # Source code
# ├── include/        # Headers
# └── tests/          # Tests

# Build the project
cnext build

# Run the executable
cnext run

# Run tests
cnext test
```

## Configuration

Edit `cnext.toml` to configure your project:

```toml
[package]
name = "my-app"
version = "0.1.0"
edition = "2024"

[build]
compiler = "auto"      # auto, gcc, clang, msvc
std = "c++20"          # C++ standard
optimization = "release"  # debug, release, size, speed
lto = true             # Link-Time Optimization
pic = true             # Position-Independent Code

[[build.executables]]
name = "my-app"
sources = ["src/**/*.cpp"]
libraries = ["mylib"]

[[build.libraries]]
name = "mylib"
type = "static"        # static or shared
sources = ["lib/**/*.cpp"]
include_dirs = ["include"]
```

## Build Commands

```bash
cnext build                # Build with default settings
cnext build --debug        # Build with debug symbols
cnext build --release      # Build with optimizations
cnext build --compiler clang  # Use specific compiler
cnext build --jobs 8       # Parallel builds
```

## Adding Dependencies

```bash
# Add a dependency
cnext add fmt@10.2.1
cnext add yaml-cpp@0.8.0

# Add a dev dependency
cnext add gtest@1.14.0 --dev

# Install all dependencies
cnext install

# Remove a dependency
cnext remove fmt
```

## Testing

cnext includes a built-in test framework:

```cpp
#include <cnext/test.h>

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

## Toolchain Management

```bash
# List installed toolchains
cnext toolchain list

# Install a toolchain
cnext toolchain install gcc@14.2.0
cnext toolchain install clang@18.1.0

# Remove a toolchain
cnext toolchain remove gcc@14.2.0
```

## Cross-Compilation

```bash
# Build for Linux x64
cnext build --target linux-x64

# Build for macOS ARM64
cnext build --target macos-arm64

# Build for Windows x64
cnext build --target windows-x64
```

## Next Steps

- [Configuration Reference](configuration.md)
- [CLI Reference](cli.md)
- [Testing Guide](testing.md)
- [Cross-Compilation](cross-compilation.md)
