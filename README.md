# cnext

A next-generation build tool and package manager for C/C++ projects.

## Features

- **One-click project initialization** — scaffold a complete project in seconds
- **Declarative builds** — define targets, dependencies, and options in `cnext.toml`
- **Package manager** — add, remove, and update C/C++ dependencies
- **Built-in test framework** — discover and run tests with `cnext test`
- **CMake bridging** — integrate CMake-based dependencies seamlessly
- **Cross-compilation** — target different OS and architectures from a single config
- **Multi-compiler support** — auto-detect or explicitly choose gcc, clang, or MSVC
- **Toolchain management** — install, switch, and manage compiler toolchains

## Installation

```bash
# From source (requires Go 1.25+)
go install github.com/neko233-com/cnext/cmd/cnext@latest

# Or build from source
git clone https://github.com/neko233-com/cnext.git
cd cnext
go build -o cnext ./cmd/cnext
```

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

## License

MIT License — see [LICENSE](LICENSE) for details.
