# Go 程序交叉编译和部署脚本

## 📋 概述

这是一个完整的 Go 程序交叉编译和自动部署脚本，支持多平台编译并通过 SCP 自动传输到指定服务器。

## 🚀 快速开始

### 1. 基本使用

```bash
# 给脚本添加执行权限
chmod +x cross-compile.sh

# 查看帮助信息
./cross-compile.sh help

# 编译主程序
./cross-compile.sh main

# 编译 Docker 演示程序
./cross-compile.sh docker_demo

# 编译所有程序
./cross-compile.sh all
```

### 2. 配置服务器信息

使用新的配置文件方式：

```bash
# 复制配置模板
cp config.env.example config.env

# 编辑配置文件，设置你的服务器信息
nano config.env
```

配置文件示例：
```bash
# 目标服务器配置
TARGET_SERVER="your-server-ip"
TARGET_USER="your-username"
TARGET_PATH="/path/to/deploy"

# 项目配置
PROJECT_NAME="test-go"
```

### 3. 运行编译

```bash
# 脚本会自动加载 config.env 配置
./cross-compile.sh main
```

## 📦 支持的平台

| 平台 | 架构 | 输出格式 |
|------|------|----------|
| Linux | AMD64 | tar.gz |
| Linux | ARM64 | tar.gz |
| Windows | AMD64 | zip |
| macOS | AMD64 | tar.gz |
| macOS | ARM64 | tar.gz |

## 🔧 功能特性

### ✅ 核心功能
- [x] 多平台交叉编译
- [x] 自动压缩打包
- [x] SCP 自动传输
- [x] 依赖检查
- [x] 连接测试
- [x] 错误处理
- [x] 彩色日志输出
- [x] 版本信息生成

### 🛡️ 安全特性
- SSH 连接测试
- 批处理模式
- 错误时自动退出
- 权限检查

### 📊 编译选项
- 优化的二进制文件 (`-ldflags="-s -w"`)
- 自动版本信息
- Git 集成
- 平台特定的文件名

## 📁 目录结构

```
.
├── cross-compile.sh      # 主编译脚本
├── deploy-config.sh      # 配置文件
├── DEPLOY_README.md      # 使用说明
├── main.go              # 主程序源码
├── docker_demo/         # Docker 演示程序
│   └── docker_demo.go   # Docker 演示源码
└── build/               # 编译输出目录
    ├── main/
    │   ├── linux-amd64/
    │   ├── linux-arm64/
    │   ├── windows-amd64/
    │   ├── darwin-amd64/
    │   └── darwin-arm64/
    └── docker_demo/
        └── ...
```

## 🔐 SSH 配置

### 方法 1: SSH 密钥认证（推荐）

```bash
# 生成 SSH 密钥
ssh-keygen -t rsa -b 4096 -C "your_email@example.com"

# 复制公钥到服务器（替换为你的服务器IP）
ssh-copy-id your-user@your-server-ip

# 测试连接
ssh your-user@your-server-ip
```

### 方法 2: 密码认证

脚本会提示输入密码，但不推荐在生产环境中使用。

## 📝 使用示例

### 示例 1: 编译并部署主程序

```bash
./cross-compile.sh main
```

输出：
```
=============================================================================
Go 程序交叉编译和部署脚本
=============================================================================
[STEP] 检查依赖...
[SUCCESS] 所有依赖检查通过
[STEP] 测试服务器连接...
[SUCCESS] 服务器连接成功
[STEP] 创建目标目录...
[SUCCESS] 目标目录创建成功
[STEP] 清理旧的编译文件...
[SUCCESS] 清理完成
[STEP] 编译程序: main
[INFO] 编译 main for linux/amd64...
[SUCCESS] 编译完成: main for linux-amd64
...
[STEP] 上传 main 到服务器...
[SUCCESS] 上传成功: linux-amd64.tar.gz
...
[SUCCESS] 所有文件上传完成
[STEP] 验证上传结果...
[SUCCESS] 验证成功，已上传 5 个文件
```

### 示例 2: 仅编译不上传

如果服务器不可用，脚本会询问是否继续编译：

```bash
./cross-compile.sh docker_demo
```

### 示例 3: 批量编译

```bash
./cross-compile.sh all
```

## 🛠️ 自定义配置

### 修改目标服务器

编辑 `cross-compile.sh` 中的配置：

```bash
TARGET_SERVER="your-server-ip"
TARGET_USER="your-username"
TARGET_PATH="/path/to/deploy"
```

### 添加新平台

在 `PLATFORMS` 数组中添加新平台：

```bash
PLATFORMS=(
    "linux-amd64:linux/amd64"
    "linux-arm64:linux/arm64"
    "freebsd-amd64:freebsd/amd64"  # 新增平台
)
```

## 🐛 故障排除

### 问题 1: SSH 连接失败

```bash
# 检查 SSH 服务（替换为你的服务器IP）
ssh -v your-user@your-server-ip

# 检查防火墙
telnet your-server-ip 22

# 检查密钥
ssh-add -l
```

### 问题 2: 编译失败

```bash
# 检查 Go 版本
go version

# 检查环境变量
echo $GOOS $GOARCH

# 清理缓存
go clean -cache
```

### 问题 3: 权限问题

```bash
# 检查脚本权限
ls -la cross-compile.sh

# 检查目标目录权限（替换为你的配置）
ssh your-user@your-server-ip "ls -la /path/to/deploy"
```

## 📊 监控和日志

### 查看编译结果

```bash
# 本地编译结果
ls -la build/

# 服务器上的文件（替换为你的配置）
ssh your-user@your-server-ip "ls -la /path/to/deploy/"
```

### 版本信息

每个编译包都包含 `version.txt` 文件，记录：
- 程序名称
- 编译时间
- Git 版本
- Go 版本
- 编译平台

## 🔄 CI/CD 集成

### GitHub Actions 示例

```yaml
name: Build and Deploy

on:
  push:
    branches: [ main ]

jobs:
  build:
    runs-on: ubuntu-latest
    steps:
    - uses: actions/checkout@v2
    
    - name: Setup Go
      uses: actions/setup-go@v2
      with:
        go-version: 1.21
    
    - name: Deploy
      env:
        SSH_KEY: ${{ secrets.SSH_KEY }}
      run: |
        echo "$SSH_KEY" > ~/.ssh/id_rsa
        chmod 600 ~/.ssh/id_rsa
        ./cross-compile.sh all
```

## 📚 相关文档

- [Go 交叉编译文档](https://golang.org/doc/install/source#environment)
- [SCP 手册](https://man.openbsd.org/scp)
- [SSH 配置指南](https://www.ssh.com/ssh/config/)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License
