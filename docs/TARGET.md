# cnext Target Matrix & Stability Roadmap

## Mission

**大幅降低 C/C++ 环境安装和依赖管理的痛苦。** 一行命令装好编译器，一行命令构建项目。

## 支持平台矩阵

| OS | Arch | 二进制名 | 编译器 | 状态 |
|----|------|----------|--------|------|
| Linux | x86_64 | cnext-linux-amd64 | gcc/clang | ✅ |
| Linux | aarch64 | cnext-linux-arm64 | gcc/clang | ✅ |
| macOS | x86_64 | cnext-darwin-amd64 | clang/gcc | ✅ |
| macOS | aarch64 | cnext-darwin-arm64 | clang/gcc | ✅ |
| Windows | x86_64 | cnext-windows-amd64.exe | msvc/gcc/clang | ✅ |
| Windows | aarch64 | cnext-windows-arm64.exe | msvc/gcc/clang | ✅ |

## Docker 验证矩阵

每个平台需要独立 Dockerfile 验证编译 + 运行：

| Dockerfile | 基础镜像 | 编译器 | 验证内容 |
|------------|----------|--------|----------|
| docker/Dockerfile.linux-amd64 | ubuntu:22.04 | gcc + clang | C/C++ 编译 + 运行 |
| docker/Dockerfile.linux-arm64 | arm64v8/ubuntu:22.04 | gcc + clang | C/C++ 编译 + 运行 |
| docker/Dockerfile.macos-cross | (交叉编译) | osxcross | macOS 二进制构建 |
| docker/Dockerfile.windows-cross | (交叉编译) | mingw | Windows 二进制构建 |

## 核心能力要求

### P0 - 必须实现（阻塞发布）

1. **增量编译** — ✅ 记录文件 hash，只重编变化的 .c/.cpp
2. **静态库打包** — ✅ 调用 `ar` (Unix) / `lib.exe` (MSVC) 生成 .a/.lib
3. **共享库构建** — ✅ 调用 `gcc -shared` / `clang -shared` 生成 .so/.dylib/.dll
4. **头文件依赖追踪** — ✅ 使用 `-MMD -MF` 生成 .d 文件，改 .h 自动重编依赖
5. **链接阶段传递库路径** — ✅ 链接时正确传递 `-L` 和 `-l` 参数

### P1 - 应该实现（影响体验）

6. **并行编译** — ✅ `--jobs N` 参数，goroutine 并发编译
7. **编译缓存** — ✅ 缓存 .o 文件到 `.cnext/cache/`
8. **错误信息美化** — 待实现
9. **watch 模式** — 待实现

### P2 - 可以延后

10. package registry 实际对接
11. lock 文件生成
12. `fmt`/`lint`/`doc` 命令实现

## 稳定性要求

- 所有 6 个平台的 CI 必须通过
- Docker 验证必须覆盖：init → build → run → test 完整流程
- 每次 PR 必须运行全平台测试
- Release 前必须本地验证所有 Docker 验证脚本

## 安装体验目标

**Before (痛苦):**
```bash
# 用户需要手动安装 gcc/clang/cmake，配置 PATH，处理版本冲突
sudo apt install gcc g++ clang cmake    # Linux
brew install gcc cmake                   # macOS
# Windows 需要安装 Visual Studio 或 MSYS2
```

**After (cnext):**
```bash
# 一行安装 cnext，自动处理编译器
curl -fsSL https://raw.githubusercontent.com/neko233-com/cnext/main/scripts/install.sh | bash

# 一行初始化项目
cnext init my-app && cd my-app

# 一行构建（自动检测/安装编译器）
cnext build
```
