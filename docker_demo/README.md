# Docker 容器技术演示程序

这是项目中的核心演示程序，详细展示 Docker 容器技术的核心概念和实现原理。

> 📋 **概述**: 本程序是 `docker_demo.go` 的详细使用说明。查看项目根目录的 `README.md` 了解项目整体介绍。

## 系统要求

### 推荐环境
- **操作系统**: Linux (推荐 Ubuntu 20.04+)
- **权限**: root 权限
- **Go 版本**: 1.16+

### 其他系统支持
- **macOS**: 支持交叉编译，生成 Linux 可执行文件
- **Windows**: 支持交叉编译，生成 Linux 可执行文件

## 使用方法

### 编译程序

#### 在 Linux 系统上直接编译
```bash
go build -o docker_demo docker_demo.go
```

#### 交叉编译（推荐）
在 macOS 或 Windows 上编译 Linux 版本：
```bash
# 编译 Linux AMD64 版本
GOOS=linux GOARCH=amd64 go build -o docker_demo_linux docker_demo.go

# 编译 Linux ARM64 版本
GOOS=linux GOARCH=arm64 go build -o docker_demo_linux_arm64 docker_demo.go
```

### 运行程序

#### 在 Linux 系统上（推荐）
```bash
# 直接运行的版本
sudo ./docker_demo

# 交叉编译的版本
sudo ./docker_demo_linux
```

#### 在 macOS/Windows 上
```bash
# 本地编译版本（仅理论演示）
./docker_demo

# 交叉编译版本需要传输到 Linux 系统运行
scp docker_demo_linux user@linux-server:/path/
```

## 程序输出说明

### 在非 Linux 系统上
程序会显示每个组件的理论知识和实现原理，包括：
- 各个 Namespace 的作用和实现方式
- Cgroup 的工作原理
- UnionFS 的概念
- 容器文件系统的结构

### 在 Linux 系统上（root 权限）
程序会实际执行以下操作：
- 创建各种 Namespace（部分演示）
- 设置 Cgroup 资源限制
- 创建 UnionFS 目录结构
- 构建容器根文件系统
- 自动清理所有创建的资源

## 资源管理

程序内置了资源管理器，会自动清理演示过程中创建的所有资源：
- 临时目录
- Cgroup 控制组
- 挂载点
- Namespace（进程结束时自动清理）

## 学习价值

通过这个程序，你可以学到：

1. **容器技术的本质**: 理解容器不是虚拟机，而是进程隔离
2. **Linux 内核特性**: 了解 Namespaces 和 Cgroups 的工作原理
3. **分层存储**: 理解 Docker 镜像的分层机制
4. **资源管理**: 学习如何限制和控制容器资源使用
5. **系统编程**: 了解 Linux 系统调用的使用

## 扩展学习

建议进一步学习：
- Docker 的具体实现
- Kubernetes 的容器编排
- 容器安全最佳实践
- 容器网络模型
- 容器存储方案

## 注意事项

1. **权限要求**: 在 Linux 上需要 root 权限才能看到完整效果
2. **系统兼容性**: 某些功能可能需要特定的内核版本支持
3. **资源清理**: 程序会自动清理资源，但如果异常退出可能需要手动清理
4. **教育目的**: 本程序仅用于学习演示，不建议在生产环境使用

## 故障排除

### 常见问题

1. **权限不足**
   ```
   解决方案：使用 sudo 运行程序
   ```

2. **cgroup 文件系统不存在**
   ```
   解决方案：检查内核配置和系统设置
   ```

3. **编译错误**
   ```
   解决方案：确保使用支持的 Go 版本
   ```

### 手动清理资源
如果程序异常退出，可以手动清理：
```bash
# 清理临时目录
sudo rm -rf /tmp/container-rootfs-demo
sudo rm -rf /tmp/mount-namespace-demo
sudo rm -rf /tmp/unionfs-demo

# 清理 cgroup
sudo rm -rf /sys/fs/cgroup/memory/docker-demo
```

## 许可证

本程序仅用于教育目的，可自由使用和修改。
