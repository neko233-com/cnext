# cnext — Next-Gen C/C++ Ecosystem Toolchain

## [S1] Problem

The C/C++ development ecosystem is fragmented and painful:
- CMake is verbose, has a steep learning curve, and generates poor error messages
- Package management is ad-hoc (vcpkg, conan, manual git submodules)
- No standard test framework — developers choose between gtest, catch2, doctest, etc.
- Installing compilers requires platform-specific setup (MSVC, Xcode, gcc)
- Cross-compilation is complex and error-prone
- No unified `npm install` equivalent for C/C++

## [S2] Solution Overview

cnext is a Go-based toolchain that provides a complete C/C++ development experience:

```
┌─────────────────────────────────────────────────┐
│                  cnext CLI                       │
├─────────────┬──────────────┬────────────────────┤
│  cnext build │  cnext test  │  cnext add/remove  │
├─────────────┴──────────────┴────────────────────┤
│           cnext core (Go 1.26)                  │
├─────────────┬──────────────┬────────────────────┤
│   Frontend   │   Builder    │  Package Manager   │
│  (Parser/    │  (Graph/     │  (Registry Client) │
│   AST/       │   Compiler   │                    │
│   Error)     │   Wrapping)  │                    │
├─────────────┴──────────────┴────────────────────┤
│         Bundled Toolchains                       │
│    gcc  │  clang  │  msvc (headers)             │
│    x86_64 + arm64  │  Linux/macOS/Windows       │
└─────────────────────────────────────────────────┘
         ↕                        ↕
    cnext-cloud              Third-party
    (Registry)               CMake Libraries
```

## [S3] Compiler Architecture

### [S3.1] Bundled Toolchains

cnext ships with pre-compiled toolchains for zero-install experience:

| Platform | Toolchain | Targets |
|----------|-----------|---------|
| Linux x86_64 | gcc 14 + clang 18 | x86_64, aarch64 |
| Linux arm64 | gcc 14 + clang 18 | aarch64, x86_64 (cross) |
| macOS arm64 | clang 18 (Apple fork) | arm64, x86_64 (universal) |
| Windows x86_64 | gcc 14 (MinGW) + clang 18 | x86_64, aarch64 |

Toolchains are downloaded on first use or bundled in installer:
- **Full installer**: ~200MB, includes all toolchains
- **Minimal installer**: ~5MB, downloads toolchains on demand

### [S3.2] Compiler Driver

The compiler driver wraps existing compilers with a unified interface:

```
cnext compile → analyze → resolve deps → generate build graph → invoke gcc/clang → link
```

Key features:
- **Unified flags**: Same flags work across gcc/clang/msvc
- **Better errors**: Custom error formatter with context, suggestions, and color
- **Auto-detection**: Picks best available compiler for target platform
- **Cross-compilation**: `cnext build --target aarch64-linux-gnu`

### [S3.3] Custom Frontend (Phase 2)

Enhanced C/C++ dialect with quality-of-life features:
- Better error messages with source context
- Module support (C++20 modules, precompiled headers)
- Import system for header-only libraries

## [S4] Build System

### [S4.1] Configuration: cnext.toml

```toml
[package]
name = "my-project"
version = "1.0.0"
edition = "2024"  # C++23, C23, etc.

[dependencies]
fmt = "10.2.1"
spdlog = "1.13.0"
boost = { version = "1.83.0", features = ["filesystem", "system"] }

[dev-dependencies]
gtest = "1.14.0"
benchmark = "1.8.3"

[build]
compiler = "auto"  # auto, gcc, clang, msvc
std = "c++23"
optimization = "release"  # debug, release, size, speed
lto = true
pic = true
sanitize = []  # address, undefined, thread

[build.target]
os = "linux"  # auto, linux, macos, windows
arch = "x86_64"  # auto, x86_64, aarch64

[[build.libraries]]
name = "mylib"
type = "static"  # static, shared, header-only
sources = ["src/**/*.cpp"]
include_dirs = ["include"]

[[build.executables]]
name = "myapp"
sources = ["main.cpp"]
libraries = ["mylib", "fmt"]

[[build.cmake-dependencies]]
name = "opencv"
version = "4.8.0"
# cnext auto-builds CMake dependencies

[[build.tests]]
name = "my-tests"
type = "unit"  # unit, integration, benchmark
sources = ["tests/**/*.cpp"]
framework = "cnext"  # cnext, gtest, catch2
```

### [S4.2] Build Graph Generation

The build system:
1. Parses `cnext.toml`
2. Resolves dependencies (local + registry)
3. Generates dependency graph
4. Detects CMake dependencies, runs cmake automatically
5. Parallelizes compilation
6. Links final artifacts

### [S4.3] CMake Bridging

When a dependency uses CMake:
```toml
[[build.cmake-dependencies]]
name = "opencv"
version = "4.8.0"
source = "registry"  # registry, git, local
git = "https://github.com/opencv/opencv.git"
tag = "4.8.0"
cmake_options = ["-DBUILD_SHARED_LIBS=ON", "-DBUILD_TESTS=OFF"]
```

cnext automatically:
1. Clones/downloads the CMake project
2. Runs `cmake` with specified options
3. Builds with the bundled toolchain
4. Generates a `cnext`-compatible package from the output
5. Caches the result

### [S4.4] Incremental Builds

- File-level dependency tracking
- Hash-based change detection
- Parallel compilation (auto-detect CPU cores)
- Distributed compilation support (Phase 3)

## [S5] Test Framework

### [S5.1] Go-Style Testing

```cpp
#include <cnext/test.h>

// Basic test
TEST(MyTest, Addition) {
    EXPECT_EQ(2 + 2, 4);
}

// Table-driven test
TEST(MyTest, StringTrim) {
    struct TestCase {
        std::string input;
        std::string expected;
    };
    
    TestCase cases[] = {
        {"  hello  ", "hello"},
        {"hello", "hello"},
        {"  ", ""},
    };
    
    for (const auto& tc : cases) {
        EXPECT_EQ(trim(tc.input), tc.expected);
    }
}

// Benchmark
BENCHMARK(MyBenchmark, SortLargeArray) {
    std::vector<int> data = generate_large_array();
    BENCHMARK_START();
    std::sort(data.begin(), data.end());
    BENCHMARK_END();
}
```

### [S5.2] Test Runner

```
cnext test
running 3 tests...
  PASS: MyTest::Addition (0.1ms)
  PASS: MyTest::StringTrim (0.2ms)
  PASS: MyBenchmark::SortLargeArray (12.3ms)

2 passed, 0 failed
coverage: 87.3%
```

Features:
- Parallel test execution
- Color output
- Filter by name: `cnext test MyTest`
- Watch mode: `cnext test --watch`
- Coverage report: `cnext test --coverage`
- Benchmark mode: `cnext test --benchmark`

### [S5.3] Assertion Library

```cpp
#include <cnext/assert.h>

// Equality
EXPECT_EQ(a, b);
ASSERT_EQ(a, b);  // stops test on failure

// Comparison
EXPECT_GT(a, b);
EXPECT_LT(a, b);
EXPECT_GE(a, b);
EXPECT_LE(a, b);

// Boolean
EXPECT_TRUE(condition);
EXPECT_FALSE(condition);

// Floating point
EXPECT_NEAR(a, b, epsilon);

// Exceptions
EXPECT_THROW(expr, ExceptionType);
EXPECT_NO_THROW(expr);

// String
EXPECT_STR_EQ(a, b);
EXPECT_STR_CONTAINS(haystack, needle);
```

## [S6] Package Manager

### [S6.1] Module Structure

```
github.com/user/mylib/
├── cnext.toml          # package config
├── include/
│   └── mylib.h
├── src/
│   └── mylib.cpp
├── tests/
│   └── mylib_test.cpp
└── README.md
```

### [S6.2] Package Operations

```bash
# Add dependency
cnext add fmt@10.2.1
cnext add boost --features filesystem,system

# Remove dependency
cnext remove fmt

# Update dependencies
cnext update
cnext update fmt

# Search packages
cnext search json

# Install all dependencies
cnext install

# Publish package
cnext publish --registry https://cnext-cloud.dev
```

### [S6.3] Lock File (cnext.lock)

```toml
[[packages]]
name = "fmt"
version = "10.2.1"
hash = "sha256:abc123..."
source = "registry"

[[packages]]
name = "mylib"
version = "0.1.0"
hash = "git:abc123"
source = "git+https://github.com/user/mylib.git"
```

## [S7] CLI Commands

```
cnext init              # Initialize new project
cnext build             # Build project
cnext build --release   # Build in release mode
cnext test              # Run tests
cnext test --coverage   # Run tests with coverage
cnext run               # Run built executable
cnext run -- args       # Run with arguments
cnext add <pkg>         # Add dependency
cnext remove <pkg>      # Remove dependency
cnext fmt               # Format code
cnext lint              # Lint code
cnext doc               # Generate documentation
cnext publish           # Publish package
cnext clean             # Clean build artifacts
cnext update            # Update dependencies
cnext info              # Show project info
cnext toolchain         # Manage toolchains
```

## [S8] cnext-cloud Registry

### [S8.1] REST API

```
GET    /api/v1/packages              # List packages
GET    /api/v1/packages/:name        # Get package info
GET    /api/v1/packages/:name/:ver   # Get version info
POST   /api/v1/packages/:name        # Create package
PUT    /api/v1/packages/:name/:ver   # Upload version
GET    /api/v1/packages/:name/:ver/download  # Download
DELETE /api/v1/packages/:name/:ver   # Delete version
GET    /api/v1/search?q=fmt          # Search packages
```

### [S8.2] Server Implementation

- Go + Gin/Fiber web framework
- SQLite/PostgreSQL for metadata
- Object storage (S3/minio) for package files
- JWT authentication
- Rate limiting
- CDN integration for downloads

## [S9] Project Structure

```
cnext/
├── cmd/
│   └── cnext/
│       └── main.go
├── internal/
│   ├── compiler/        # Compiler driver
│   │   ├── driver.go
│   │   ├── gcc.go
│   │   ├── clang.go
│   │   └── msvc.go
│   ├── frontend/        # Parser, AST, errors
│   │   ├── parser.go
│   │   ├── ast.go
│   │   └── errors.go
│   ├── build/           # Build graph, scheduling
│   │   ├── graph.go
│   │   ├── scheduler.go
│   │   └── incremental.go
│   ├── cmake/           # CMake bridging
│   │   └── bridge.go
│   ├── test/            # Test runner
│   │   ├── runner.go
│   │   └── framework.go
│   ├── package/         # Package manager
│   │   ├── resolver.go
│   │   ├── lock.go
│   │   └── registry.go
│   └── toolchain/       # Bundled toolchains
│       ├── manager.go
│       └── download.go
├── pkg/
│   └── cnext/           # Public API
├── cnext-test/          # Test framework headers
│   ├── include/
│   │   ├── cnext/
│   │   │   ├── test.h
│   │   │   └── assert.h
│   │   └── cnext.test/
│   │       └── runner.h
│   └── src/
│       └── runner.cpp
├── docs/
├── go.mod
└── go.sum

cnext-cloud/
├── cmd/
│   └── server/
│       └── main.go
├── internal/
│   ├── api/             # REST handlers
│   ├── models/          # Database models
│   ├── storage/         # Object storage
│   └── auth/            # Authentication
├── migrations/
├── go.mod
└── go.sum
```

## [S10] Implementation Phases

### Phase 1: Foundation (Weeks 1-2)
- [ ] CLI scaffold with cobra
- [ ] cnext.toml parser
- [ ] Basic compiler driver (wrapping gcc/clang)
- [ ] Simple build command
- [ ] Project initialization (`cnext init`)

### Phase 2: Core Features (Weeks 3-4)
- [ ] Dependency resolution
- [ ] Package manager (local + git)
- [ ] Test framework (basic assertions)
- [ ] Test runner
- [ ] CMake bridging

### Phase 3: Registry (Week 5)
- [ ] cnext-cloud server scaffold
- [ ] Package upload/download
- [ ] Version management
- [ ] Search API

### Phase 4: Polish (Weeks 6-7)
- [ ] Code formatter (`cnext fmt`)
- [ ] Linter integration (`cnext lint`)
- [ ] Documentation generator
- [ ] Incremental builds
- [ ] Cross-compilation

### Phase 5: Advanced (Phase 2)
- [ ] Custom frontend (better errors)
- [ ] Distributed compilation
- [ ] IDE integration (LSP)
- [ ] Plugin system
