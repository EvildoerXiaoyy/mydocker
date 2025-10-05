# Docker 容器技术学习项目

## 🐳 项目概述

这是一个专注于 Docker 容器技术学习和演示的 Go 项目，通过实际的代码示例来展示 Docker 容器技术的核心概念和实现原理。

## 📁 项目结构

```
docker-learning/
├── docker_demo/              # Docker 技术演示程序
│   └── docker_demo.go        # 完整的 Docker 容器技术演示
├── main.go                   # 简单的 Hello World 程序
├── cross-compile.sh          # 交叉编译和部署脚本
├── config.env.example        # 配置文件模板
├── DEPLOY_README.md          # 详细部署说明
├── PROJECT_SUMMARY.md        # 项目总结
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
# 编译 Docker 演示程序
./cross-compile.sh docker_demo

# 编译所有程序
./cross-compile.sh all

# 查看帮助
./cross-compile.sh help
```

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

### docker_demo.go

这个程序展示了 Docker 容器技术的核心概念：

1. **Namespace 演示**
   - 展示各种 Namespace 的作用和实现方式
   - 说明容器隔离的原理

2. **Cgroup 演示**
   - 展示资源限制的实现
   - 说明 Cgroup 的工作原理

3. **UnionFS 演示**
   - 展示分层文件系统
   - 说明写时复制机制

4. **网络配置演示**
   - 展示容器网络配置
   - 说明网络模式

5. **安全机制演示**
   - 展示容器安全特性
   - 说明安全防护措施

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
