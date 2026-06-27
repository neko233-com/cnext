# CLI Reference

## Global Flags

```
cnext [command] [flags]
```

| Flag | Description |
|------|-------------|
| `-h, --help` | Help for cnext |
| `--version` | Version for cnext |

## Project Commands

### init

Initialize a new cnext project.

```bash
cnext init [name]
```

| Flag | Description |
|------|-------------|
| `[name]` | Project name (default: "my-project") |

### build

Build the project.

```bash
cnext build [flags]
```

| Flag | Description |
|------|-------------|
| `-t, --target string` | Build a specific target |
| `-d, --debug` | Build with debug symbols |
| `-r, --release` | Build in release mode |
| `-j, --jobs int` | Number of parallel jobs |
| `-c, --compiler string` | Compiler to use (gcc, clang, etc.) |

### run

Run the built executable.

```bash
cnext run [flags] [-- args...]
```

| Flag | Description |
|------|-------------|
| `-t, --target string` | Run a specific executable |

Arguments after `--` are passed to the executable.

### test

Run project tests.

```bash
cnext test [flags]
```

| Flag | Description |
|------|-------------|
| `-t, --target string` | Run tests for a specific target |
| `-f, --filter string` | Filter tests by name pattern |
| `-v, --verbose` | Verbose test output |
| `-c, --coverage` | Generate coverage report |
| `-b, --benchmark` | Run benchmarks |
| `-w, --watch` | Watch for changes and re-run tests |

### clean

Clean build artifacts.

```bash
cnext clean [flags]
```

| Flag | Description |
|------|-------------|
| `-a, --all` | Remove all generated files including dependencies |

### info

Show project information.

```bash
cnext info
```

## Package Commands

### add

Add a dependency.

```bash
cnext add <package>[@version] [flags]
```

| Flag | Description |
|------|-------------|
| `-v, --version string` | Specific version to add |
| `--dev` | Add as a development dependency |

### remove

Remove a dependency.

```bash
cnext remove <package>
```

### install

Install all dependencies.

```bash
cnext install
```

### update

Update dependencies.

```bash
cnext update [flags]
```

| Flag | Description |
|------|-------------|
| `--dry-run` | Preview what would be updated |

## Code Quality Commands

### fmt

Format source code.

```bash
cnext fmt [flags]
```

| Flag | Description |
|------|-------------|
| `--check` | Check if formatting is needed without modifying files |
| `--exclude stringArray` | Exclude directories or patterns |

### lint

Lint source code.

```bash
cnext lint [flags]
```

| Flag | Description |
|------|-------------|
| `-f, --fix` | Automatically fix lint issues |
| `-o, --output string` | Output format (text, json, checkstyle) |

## Documentation Commands

### doc

Generate documentation.

```bash
cnext doc [flags]
```

| Flag | Description |
|------|-------------|
| `-f, --format string` | Documentation format (html, latex, man) |
| `-o, --output string` | Output directory (default: "docs") |

### publish

Publish package to registry.

```bash
cnext publish [flags]
```

| Flag | Description |
|------|-------------|
| `-r, --registry string` | Registry to publish to |
| `--dry-run` | Preview what would be published |

## Toolchain Commands

### toolchain list

List installed toolchains.

```bash
cnext toolchain list
```

### toolchain install

Install a toolchain.

```bash
cnext toolchain install <name>[@version] [flags]
```

| Flag | Description |
|------|-------------|
| `-p, --platform string` | Target platform (default: current) |

### toolchain remove

Remove a toolchain.

```bash
cnext toolchain remove <name>[@version] [flags]
```

| Flag | Description |
|------|-------------|
| `-p, --platform string` | Target platform (default: current) |

## Examples

### Basic Workflow

```bash
# Create project
cnext init my-app
cd my-app

# Add dependencies
cnext add fmt@10.2.1
cnext add gtest@1.14.0 --dev

# Build
cnext build

# Run
cnext run

# Test
cnext test
```

### Cross-Compilation

```bash
# Build for Linux
cnext build --target linux-x64

# Build for macOS ARM
cnext build --target macos-arm64

# Build for Windows
cnext build --target windows-x64
```

### Debug Build

```bash
# Build with debug symbols
cnext build --debug

# Run
cnext run

# Test with verbose output
cnext test --verbose
```

### Toolchain Management

```bash
# List available toolchains
cnext toolchain list

# Install specific versions
cnext toolchain install gcc@14.2.0
cnext toolchain install clang@18.1.0

# Build with specific compiler
cnext build --compiler gcc
```
