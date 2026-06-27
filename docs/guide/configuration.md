# Configuration Reference

cnext uses `cnext.toml` for project configuration.

## Package Section

```toml
[package]
name = "my-app"         # Project name
version = "1.0.0"       # Semantic version
edition = "2024"        # Edition year
```

## Build Section

```toml
[build]
compiler = "auto"       # Compiler: auto, gcc, clang, msvc
std = "c++20"           # C/C++ standard: c11, c17, c++17, c++20, c++23
optimization = "release"  # Optimization: debug, release, size, speed
lto = true              # Link-Time Optimization
pic = true              # Position-Independent Code
sanitize = []           # Sanitizers: address, undefined, thread
```

## Target Configuration

```toml
[build.target]
os = "auto"             # Target OS: auto, linux, windows, darwin
arch = "auto"           # Target arch: auto, x86_64, aarch64, arm
```

### Common Targets

| Target | OS | Arch |
|--------|-----|------|
| `linux-x64` | linux | x86_64 |
| `linux-arm64` | linux | aarch64 |
| `windows-x64` | windows | x86_64 |
| `macos-arm64` | darwin | aarch64 |
| `macos-x64` | darwin | x86_64 |

## Executables

```toml
[[build.executables]]
name = "my-app"                    # Executable name
sources = ["src/**/*.cpp"]         # Source file patterns
libraries = ["mylib", "utils"]     # Library dependencies
```

## Libraries

```toml
[[build.libraries]]
name = "mylib"                     # Library name
type = "static"                    # Type: static or shared
sources = ["lib/**/*.cpp"]         # Source file patterns
include_dirs = ["include"]         # Include directories
```

## CMake Dependencies

```toml
[[build.cmake-dependencies]]
name = "opencv"                    # Dependency name
version = "4.8.0"                  # Version
source = "git"                     # Source: git, url, local
git = "https://github.com/opencv/opencv.git"
tag = "4.8.0"
cmake_options = ["-DBUILD_SHARED_LIBS=ON"]
```

## Test Targets

```toml
[[build.tests]]
name = "unit-tests"                # Test name
type = "unit"                      # Test type
sources = ["tests/**/*.cpp"]       # Test source patterns
framework = "cnext-test"           # Test framework
```

## Dependencies

```toml
[dependencies]
fmt = "10.2.1"                     # Simple version
yaml-cpp = { version = "0.8.0", features = ["shared"] }  # With features

[dev-dependencies]
gtest = "1.14.0"
benchmark = { version = "1.8.3", features = ["tools"] }
```

### Version Constraints

| Constraint | Example | Description |
|------------|---------|-------------|
| Exact | `1.0.0` | Exact version |
| Caret | `^1.0.0` | Compatible with 1.0.0 (>=1.0.0, <2.0.0) |
| Tilde | `~1.0.0` | Patch updates (>=1.0.0, <1.1.0) |
| Range | `>=1.0.0` | Minimum version |
| Range | `<2.0.0` | Maximum version |

## Complete Example

```toml
[package]
name = "my-project"
version = "1.0.0"
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
include_dirs = ["include"]

[[build.tests]]
name = "unit-tests"
sources = ["tests/**/*.cpp"]
framework = "cnext-test"

[dependencies]
fmt = "10.2.1"
spdlog = "1.13.0"

[dev-dependencies]
gtest = "1.14.0"
```
