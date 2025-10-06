# Docker 容器技术演示程序 - 交叉编译和部署指南

## 📋 概述

这是一个专门为 Docker 容器技术演示程序设计的交叉编译和自动部署脚本。该程序展示了 Docker 容器技术的核心概念和实现原理，包括 Linux Namespaces、Cgroups、UnionFS 等关键技术。

## 🐳 程序特性

### 核心演示功能
- **PID Namespace**: 进程 ID 空间隔离演示
- **Network Namespace**: 网络设备、IP 地址、路由表隔离演示
- **Mount Namespace**: 文件系统挂载点隔离演示
- **UTS Namespace**: 主机名和域名隔离演示
- **IPC Namespace**: 进程间通信资源隔离演示
- **User Namespace**: 用户和组 ID 隔离演示
- **Cgroups**: 资源限制演示（CPU、内存等）
- **UnionFS**: 分层文件系统演示
- **容器根文件系统**: 完整的容器环境创建

### 系统要求
- **推荐环境**: Linux 系统 + root 权限
- **交叉编译**: 支持 macOS/Windows 编译 Linux 版本
- **Go 版本**: 1.16+

## 🚀 快速开始

### 1. 基本使用

```bash
# 给脚本添加执行权限
chmod +x cross-compile.sh

# 查看帮助信息
./cross-compile.sh help

# 编译 Docker 演示程序
./cross-compile.sh docker_demo

# 编译所有支持的平台
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
./cross-compile.sh docker_demo
```

## 🎯 Docker 演示程序特定说明

### 编译和部署流程

1. **交叉编译**: 在 macOS/Windows 上编译 Linux 版本
2. **自动部署**: 通过 SCP 传输到 Linux 服务器
3. **权限运行**: 在目标服务器上以 root 权限运行

### 推荐部署方式

```bash
# 方式 1: 本地交叉编译 + 自动部署
./cross-compile.sh docker_demo

# 方式 2: 仅本地编译（用于测试）
GOOS=linux GOARCH=amd64 go build -o docker_demo_linux docker_demo/docker_demo.go

# 方式 3: 在目标 Linux 服务器上直接编译
git clone <repository>
cd test-go/docker_demo
go build -o docker_demo docker_demo.go
sudo ./docker_demo
```

### 运行要求

**在 Linux 服务器上运行：**
```bash
# 解压编译好的程序
tar -xzf linux-amd64.tar.gz
cd linux-amd64/

# 需要 root 权限才能看到完整演示
sudo ./docker_demo
```

**预期输出示例：**
```
Docker 容器技术完整演示程序
🐳 本程序展示了 Docker 容器技术的核心概念和实现原理
⚠️  仅支持 Linux 系统，需要 root 权限

=== Docker 容器技术完整演示 ===
🖥️  运行环境: linux
👤 当前用户: root

=== PID Namespace 演示 ===
🔍 当前进程 PID: 1234
🚀 创建 PID Namespace...
✅ PID Namespace 创建成功
🔄 在新 PID Namespace 中创建子进程...
👨 父进程看到子进程 PID: 1235
💡 PID Namespace 效果：子进程在新 namespace 中获得 PID 1
```

## 🏗️ 支持的平台

| 平台 | 架构 | 输出格式 |
|------|------|----------|
| Linux | AMD64 | tar.gz |
| Linux | ARM64 | tar.gz |
| Windows | AMD64 | zip |
| macOS | AMD64 | tar.gz |
| macOS | ARM64 | tar.gz |

## ⚙️ 功能特性

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
├── config.env            # 配置文件
├── DEPLOY_README.md      # 使用说明
├── docker_demo/          # Docker 演示程序
│   └── docker_demo.go    # Docker 演示源码
└── build/                # 编译输出目录
    └── docker_demo/
        ├── linux-amd64/
        ├── linux-arm64/
        ├── windows-amd64/
        ├── darwin-amd64/
        └── darwin-arm64/
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

### 示例 1: 编译并部署 Docker 演示程序

```bash
./cross-compile.sh docker_demo
```

输出：
```
==============================================================================
Go 程序交叉编译和部署脚本
==============================================================================
[STEP] 检查依赖...
[SUCCESS] 所有依赖检查通过
[STEP] 测试服务器连接...
[SUCCESS] 服务器连接成功
[STEP] 创建目标目录...
[SUCCESS] 目标目录创建成功
[STEP] 清理旧的编译文件...
[SUCCESS] 清理完成
[STEP] 编译程序: docker_demo
[INFO] 编译 docker_demo for linux/amd64...
[SUCCESS] 编译完成: docker_demo for linux-amd64
...
[STEP] 上传 docker_demo 到服务器...
[SUCCESS] 上传成功: linux-amd64.tar.gz
...
[SUCCESS] 所有文件上传完成
[STEP] 验证上传结果...
[SUCCESS] 验证成功，已上传 2 个文件
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

编辑 `config.env` 中的配置：

```bash
TARGET_SERVER="your-server-ip"
TARGET_USER="your-username"
TARGET_PATH="/path/to/deploy"
```

### 添加新平台

在 `cross-compile.sh` 的 `PLATFORMS` 数组中添加新平台：

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

### 🔧 Docker 演示程序特定问题

#### 问题 4: 在非 Linux 系统运行

**症状**: 提示 "当前系统不支持，请使用 Linux 系统"

**解决方案**:
```bash
# 方式 1: 使用交叉编译
GOOS=linux GOARCH=amd64 go build -o docker_demo_linux docker_demo/docker_demo.go

# 方式 2: 在 Linux 服务器上运行
scp docker_demo_linux user@linux-server:/tmp/
ssh user@linux-server
sudo /tmp/docker_demo_linux
```

#### 问题 5: 权限不足

**症状**: 提示 "需要 root 权限"

**解决方案**:
```bash
# 使用 sudo 运行
sudo ./docker_demo

# 或者切换到 root 用户
su - root
./docker_demo
```

#### 问题 6: Namespace 创建失败

**症状**: 提示 "创建失败: operation not permitted"

**解决方案**:
```bash
# 检查内核版本（需要 3.8+）
uname -r

# 检查是否在容器中运行
grep container /proc/1/cgroup

# 某些 Namespace 在容器中不可用，这是正常现象
```

#### 问题 7: Cgroup 文件系统不存在

**症状**: 提示 "Cgroup 文件系统不存在"

**解决方案**:
```bash
# 检查 cgroup 挂载
mount | grep cgroup

# 手动挂载 cgroup（需要 root）
sudo mount -t tmpfs cgroup /sys/fs/cgroup
sudo mkdir -p /sys/fs/cgroup/memory
sudo mount -t cgroup -o memory cgroup /sys/fs/cgroup/memory
```

#### 问题 8: OverlayFS 挂载失败

**症状**: 提示 "挂载失败: operation not permitted"

**解决方案**:
```bash
# 检查内核是否支持 overlay
grep overlay /proc/filesystems

# 检查是否在容器中运行
# 在容器中通常无法挂载新的文件系统
```

#### 问题 9: 资源清理失败

**症状**: 程序异常退出后资源未清理

**解决方案**:
```bash
# 手动清理临时目录
sudo rm -rf /tmp/container-rootfs-demo
sudo rm -rf /tmp/mount-namespace-demo
sudo rm -rf /tmp/unionfs-demo

# 手动清理 cgroup
sudo rm -rf /sys/fs/cgroup/memory/docker-demo

# 检查挂载点
mount | grep tmpfs
sudo umount /tmp/mount-namespace-demo/tmpfs
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
- [Docker 容器技术原理](https://docs.docker.com/engine/overview/)

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

## 📄 许可证

MIT License
