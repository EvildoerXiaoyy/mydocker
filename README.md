# Docker 容器技术学习项目

## 🐳 项目概述

这是一个专注于 Docker 容器技术学习和演示的 Go 项目，通过实际的代码示例来展示 Docker 容器技术的核心概念和实现原理。

## 📁 项目结构

```
test-go/                      # 项目根目录
├── docker_demo/              # Docker 技术演示程序
│   ├── docker_demo.go        # Docker 容器技术演示源码
│   └── README.md             # Docker 演示程序说明
├── cross-compile.sh          # 交叉编译和部署脚本
├── config.env                # 部署配置文件
├── config.env.example        # 配置文件模板
├── DEPLOY_README.md          # 详细部署说明
├── PROJECT_SUMMARY.md        # 项目总结
├── README.md                 # 项目说明文档
├── go.mod                    # Go 模块文件
└── launch.json               # VS Code 调试配置
```

## 🚀 快速开始

### 运行 Docker 演示程序

```bash
# 进入 Docker 演示目录
cd docker_demo

# 运行演示程序
go run docker_demo.go
```

### 交叉编译和部署

```bash
# 配置部署信息（首次使用）
cp config.env.example config.env
# 编辑 config.env 文件，填入实际的服务器信息

# 编译 Docker 演示程序
./cross-compile.sh docker_demo

# 编译所有程序
./cross-compile.sh all

# 查看帮助
./cross-compile.sh help
```

#### 脚本功能特性

- ✅ **多平台编译**: 支持 Linux、Windows、macOS (Intel/Apple Silicon)
- ✅ **自动部署**: 编译完成后自动上传到指定服务器
- ✅ **配置管理**: 通过配置文件管理部署参数
- ✅ **智能清理**: 上传完成后自动清理本地打包文件
- ✅ **错误处理**: 完善的错误检查和日志输出
- ✅ **版本信息**: 自动生成包含 Git 版本和编译信息的文件

#### 支持的目标平台

| 平台 | 架构 | 输出格式 |
|------|------|----------|
| Linux | amd64, arm64 | tar.gz |
| Windows | amd64 | zip |
| macOS | amd64, arm64 | tar.gz |

## 📚 学习内容

### Docker 容器技术核心概念

1. **Linux Namespaces**
   - PID Namespace: 进程隔离
   - Network Namespace: 网络隔离
   - Mount Namespace: 文件系统隔离
   - UTS Namespace: 主机名隔离
   - IPC Namespace: 进程间通信隔离
   - User Namespace: 用户隔离

2. **Cgroups**
   - CPU 限制
   - 内存限制
   - 进程数量限制
   - 磁盘 I/O 限制

3. **UnionFS (OverlayFS)**
   - 分层文件系统
   - 写时复制机制
   - 镜像分层存储

4. **容器网络**
   - Bridge 模式
   - Host 模式
   - None 模式
   - Container 模式

5. **容器安全**
   - 能力限制
   - Seccomp 过滤
   - 网络安全

## 🛠️ 技术特性

- ✅ 完整的 Docker 技术演示
- ✅ 多平台交叉编译支持
- ✅ 自动化部署流程
- ✅ 详细的代码注释
- ✅ 跨平台兼容性

## 📖 演示程序功能

项目包含一个完整的 Docker 容器技术演示程序 (`docker_demo/docker_demo.go`)，展示了以下核心概念：

1. **Namespace 演示** - 进程隔离、网络隔离、文件系统隔离等
2. **Cgroup 演示** - CPU、内存、磁盘 I/O 等资源限制
3. **UnionFS 演示** - 分层文件系统和写时复制机制
4. **容器根文件系统** - 完整的容器环境创建
5. **安全机制演示** - 容器安全特性展示

> 💡 **详细说明**: 查看 `docker_demo/README.md` 了解演示程序的详细使用方法、系统要求和故障排除指南。

## 🔧 开发环境

- Go 1.21+
- 支持 Linux/macOS/Windows
- 需要 Git、SSH、SCP 工具

## 📦 部署

项目支持自动化部署到指定服务器：

```bash
# 配置服务器信息（首次需要编辑 config.env 文件）
cp config.env.example config.env  # 如果有模板文件
# 编辑 config.env 文件，设置服务器信息

# 编译并部署
./cross-compile.sh docker_demo
```

### 配置文件说明

创建 `config.env` 文件来配置部署参数：

```bash
# 目标服务器配置
TARGET_SERVER="your-server-ip"
TARGET_USER="your-username"
TARGET_PATH="/path/to/deploy"

# 项目配置
PROJECT_NAME="your-project-name"
```

## 🎯 学习目标

通过这个项目，你将学习到：

1. Docker 容器技术的核心原理
2. Linux 系统编程基础
3. Go 语言系统调用
4. 跨平台编译技术
5. 自动化部署流程

## 📄 许可证

MIT License

## 🤝 贡献

欢迎提交 Issue 和 Pull Request！

---

**注意**: 某些功能（如 Namespace 创建）需要在 Linux 环境下运行才能看到完整效果。在 macOS 和 Windows 上，程序会显示相应的说明信息。
