// Docker 容器技术完整演示程序
// 本程序展示了 Docker 容器技术的核心概念和实现原理
// 包括：Linux Namespaces、Cgroups、UnionFS 等关键技术
// 仅支持 Linux 平台，使用交叉编译在其他系统上构建
//go:build linux

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"syscall"
)

// isLinux 检查是否在 Linux 系统上运行
func isLinux() bool {
	return runtime.GOOS == "linux"
}

// isRoot 检查当前用户是否为 root 用户
func isRoot() bool {
	return os.Getuid() == 0
}

// ResourceManager 管理演示过程中创建的资源
type ResourceManager struct {
	cgroups      []string
	mountPoints  []string
	tempDirs     []string
	namespaces   []string // 记录创建的 namespace
}

// NewResourceManager 创建资源管理器
func NewResourceManager() *ResourceManager {
	return &ResourceManager{
		cgroups:     make([]string, 0),
		mountPoints:  make([]string, 0),
		tempDirs:     make([]string, 0),
		namespaces:   make([]string, 0),
	}
}

// AddCgroup 添加 cgroup 路径到清理列表
func (rm *ResourceManager) AddCgroup(path string) {
	rm.cgroups = append(rm.cgroups, path)
}

// AddMountPoint 添加挂载点到清理列表
func (rm *ResourceManager) AddMountPoint(path string) {
	rm.mountPoints = append(rm.mountPoints, path)
}

// AddTempDir 添加临时目录到清理列表
func (rm *ResourceManager) AddTempDir(path string) {
	rm.tempDirs = append(rm.tempDirs, path)
}

// AddNamespace 添加 namespace 到记录列表
func (rm *ResourceManager) AddNamespace(nsType string) {
	rm.namespaces = append(rm.namespaces, nsType)
}

// Cleanup 清理所有资源
func (rm *ResourceManager) Cleanup() {
	fmt.Println("🧹 开始清理演示资源...")
	
	// 清理挂载点
	for i := len(rm.mountPoints) - 1; i >= 0; i-- {
		mount := rm.mountPoints[i]
		if err := syscall.Unmount(mount, 0); err != nil {
			fmt.Printf("❌ 卸载挂载点失败 %s: %v\n", mount, err)
		} else {
			fmt.Printf("✅ 卸载挂载点: %s\n", mount)
		}
	}
	
	// 清理 cgroup
	for _, cgroup := range rm.cgroups {
		if err := os.RemoveAll(cgroup); err != nil {
			fmt.Printf("❌ 删除 cgroup 失败 %s: %v\n", cgroup, err)
		} else {
			fmt.Printf("✅ 删除 cgroup: %s\n", cgroup)
		}
	}
	
	// 清理临时目录
	for _, dir := range rm.tempDirs {
		if err := os.RemoveAll(dir); err != nil {
			fmt.Printf("❌ 删除临时目录失败 %s: %v\n", dir, err)
		} else {
			fmt.Printf("✅ 删除临时目录: %s\n", dir)
		}
	}
	
	// Namespace 会在进程结束时自动清理
	if len(rm.namespaces) > 0 {
		fmt.Printf("✅ Namespace 将在进程结束时自动清理: %v\n", rm.namespaces)
	}
	
	fmt.Println("🧹 资源清理完成")
}

// demonstratePIDNamespace 演示 PID Namespace
func demonstratePIDNamespace(rm *ResourceManager) {
	fmt.Println("=== PID Namespace 演示 ===")
	
	if !isLinux() {
		fmt.Println("❌ PID Namespace 需要 Linux 系统")
		return
	}
	
	if !isRoot() {
		fmt.Println("❌ 需要 root 权限")
		return
	}
	
	fmt.Printf("🔍 当前进程 PID: %d\n", os.Getpid())
	
	// 创建新的 PID namespace
	fmt.Println("🚀 创建 PID Namespace...")
	if err := syscall.Unshare(syscall.CLONE_NEWPID); err != nil {
		fmt.Printf("❌ 创建失败: %v\n", err)
		return
	}
	
	rm.AddNamespace("PID")
	fmt.Println("✅ PID Namespace 创建成功")
	
	// 在新 namespace 中创建子进程
	fmt.Println("🔄 在新 PID Namespace 中创建子进程...")
	
	cmd := exec.Command("sleep", "1")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID,
	}
	
	if err := cmd.Start(); err != nil {
		fmt.Printf("❌ 创建子进程失败: %v\n", err)
		return
	}
	
	pid := cmd.Process.Pid
	fmt.Printf("👨 父进程看到子进程 PID: %d\n", pid)
	
	// 等待子进程完成
	if err := cmd.Wait(); err != nil {
		fmt.Printf("❌ 等待子进程失败: %v\n", err)
	}
	
	// 创建另一个子进程来验证 PID 1
	fmt.Println("🔄 验证新 PID Namespace 中的进程...")
	
	verifyCmd := exec.Command("sh", "-c", "echo '👶 子进程 PID:' $$ && ps -ef")
	verifyCmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWPID,
	}
	verifyCmd.Stdout = os.Stdout
	verifyCmd.Stderr = os.Stderr
	
	if err := verifyCmd.Run(); err != nil {
		fmt.Printf("❌ 验证子进程失败: %v\n", err)
	}
	
	fmt.Println("💡 PID Namespace 效果：子进程在新 namespace 中获得 PID 1")
}

// demonstrateNetworkNamespace 演示 Network Namespace
func demonstrateNetworkNamespace(rm *ResourceManager) {
	fmt.Println("=== Network Namespace 演示 ===")
	
	if !isLinux() {
		fmt.Println("❌ Network Namespace 需要 Linux 系统")
		return
	}
	
	if !isRoot() {
		fmt.Println("❌ 需要 root 权限")
		return
	}
	
	// 显示当前网络接口
	fmt.Println("🔍 当前网络接口：")
	if cmd := exec.Command("ip", "addr", "show"); cmd.Run() == nil {
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else if cmd := exec.Command("ifconfig"); cmd.Run() == nil {
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	
	// 创建新的 Network namespace
	fmt.Println("🚀 创建 Network Namespace...")
	if err := syscall.Unshare(syscall.CLONE_NEWNET); err != nil {
		fmt.Printf("❌ 创建失败: %v\n", err)
		return
	}
	
	rm.AddNamespace("Network")
	fmt.Println("✅ Network Namespace 创建成功")
	
	// 显示新 namespace 中的网络接口
	fmt.Println("🔍 新 Network Namespace 中的网络接口：")
	if cmd := exec.Command("ip", "addr", "show"); cmd.Run() == nil {
		cmd.Stdout = os.Stdout
		cmd.Run()
	} else if cmd := exec.Command("ifconfig"); cmd.Run() == nil {
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	
	fmt.Println("💡 Network Namespace 效果：只有 lo 回环接口")
}

// demonstrateMountNamespace 演示 Mount Namespace
func demonstrateMountNamespace(rm *ResourceManager) {
	fmt.Println("=== Mount Namespace 演示 ===")
	
	if !isLinux() {
		fmt.Println("❌ Mount Namespace 需要 Linux 系统")
		return
	}
	
	if !isRoot() {
		fmt.Println("❌ 需要 root 权限")
		return
	}
	
	// 创建临时目录
	tempDir := "/tmp/mount-namespace-demo"
	os.MkdirAll(tempDir, 0755)
	rm.AddTempDir(tempDir)
	
	// 创建新的 Mount namespace
	fmt.Println("🚀 创建 Mount Namespace...")
	if err := syscall.Unshare(syscall.CLONE_NEWNS); err != nil {
		fmt.Printf("❌ 创建失败: %v\n", err)
		return
	}
	
	rm.AddNamespace("Mount")
	fmt.Println("✅ Mount Namespace 创建成功")
	
	// 在新 namespace 中挂载 tmpfs
	tmpfsPath := filepath.Join(tempDir, "tmpfs")
	os.MkdirAll(tmpfsPath, 0755)
	
	if err := syscall.Mount("tmpfs", tmpfsPath, "tmpfs", 0, ""); err != nil {
		fmt.Printf("❌ 挂载 tmpfs 失败: %v\n", err)
		return
	}
	
	rm.AddMountPoint(tmpfsPath)
	
	// 创建测试文件
	testFile := filepath.Join(tmpfsPath, "test.txt")
	ioutil.WriteFile(testFile, []byte("Hello from Mount Namespace!"), 0644)
	
	fmt.Printf("✅ 在新 Mount Namespace 中创建文件: %s\n", testFile)
	
	// 显示挂载信息
	fmt.Println("🔍 新 Mount Namespace 中的挂载点：")
	cmd := exec.Command("mount")
	cmd.Stdout = os.Stdout
	cmd.Run()
	
	fmt.Println("💡 Mount Namespace 效果：挂载只在当前 namespace 中可见")
}

// demonstrateUTSNamespace 演示 UTS Namespace
func demonstrateUTSNamespace(rm *ResourceManager) {
	fmt.Println("=== UTS Namespace 演示 ===")
	
	if !isLinux() {
		fmt.Println("❌ UTS Namespace 需要 Linux 系统")
		return
	}
	
	if !isRoot() {
		fmt.Println("❌ 需要 root 权限")
		return
	}
	
	// 显示当前主机名
	hostname, _ := os.Hostname()
	fmt.Printf("🔍 当前主机名: %s\n", hostname)
	
	// 创建新的 UTS namespace
	fmt.Println("🚀 创建 UTS Namespace...")
	if err := syscall.Unshare(syscall.CLONE_NEWUTS); err != nil {
		fmt.Printf("❌ 创建失败: %v\n", err)
		return
	}
	
	rm.AddNamespace("UTS")
	fmt.Println("✅ UTS Namespace 创建成功")
	
	// 设置新主机名
	newHostname := "container-demo"
	if err := syscall.Sethostname([]byte(newHostname)); err != nil {
		fmt.Printf("❌ 设置主机名失败: %v\n", err)
		return
	}
	
	fmt.Printf("✅ 设置新主机名: %s\n", newHostname)
	
	// 验证主机名设置
	currentHostname, _ := os.Hostname()
	fmt.Printf("🔍 验证当前主机名: %s\n", currentHostname)
	
	fmt.Println("💡 UTS Namespace 效果：主机名变更只在当前 namespace 中有效")
}

// demonstrateIPCNamespace 演示 IPC Namespace
func demonstrateIPCNamespace(rm *ResourceManager) {
	fmt.Println("=== IPC Namespace 演示 ===")
	
	if !isLinux() {
		fmt.Println("❌ IPC Namespace 需要 Linux 系统")
		return
	}
	
	if !isRoot() {
		fmt.Println("❌ 需要 root 权限")
		return
	}
	
	// 显示当前 IPC 资源
	fmt.Println("🔍 当前 IPC 资源：")
	if cmd := exec.Command("ipcs"); cmd.Run() == nil {
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	
	// 创建新的 IPC namespace
	fmt.Println("🚀 创建 IPC Namespace...")
	if err := syscall.Unshare(syscall.CLONE_NEWIPC); err != nil {
		fmt.Printf("❌ 创建失败: %v\n", err)
		return
	}
	
	rm.AddNamespace("IPC")
	fmt.Println("✅ IPC Namespace 创建成功")
	
	// 显示新 namespace 中的 IPC 资源
	fmt.Println("🔍 新 IPC Namespace 中的资源：")
	if cmd := exec.Command("ipcs"); cmd.Run() == nil {
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	
	fmt.Println("💡 IPC Namespace 效果：新的 namespace 中没有继承原有的 IPC 资源")
}

// demonstrateUserNamespace 演示 User Namespace
func demonstrateUserNamespace(rm *ResourceManager) {
	fmt.Println("=== User Namespace 演示 ===")
	
	if !isLinux() {
		fmt.Println("❌ User Namespace 需要 Linux 系统")
		return
	}
	
	if !isRoot() {
		fmt.Println("❌ 需要 root 权限")
		return
	}
	
	// 显示当前用户信息
	fmt.Printf("🔍 当前用户 UID: %d, GID: %d\n", os.Getuid(), os.Getgid())
	
	// 创建新的 User namespace
	fmt.Println("🚀 创建 User Namespace...")
	if err := syscall.Unshare(syscall.CLONE_NEWUSER); err != nil {
		fmt.Printf("❌ 创建失败: %v\n", err)
		fmt.Println("💡 User Namespace 可能需要特殊配置")
		return
	}
	
	rm.AddNamespace("User")
	fmt.Println("✅ User Namespace 创建成功")
	
	fmt.Printf("🔍 新 User Namespace 中的 UID: %d, GID: %d\n", os.Getuid(), os.Getgid())
	fmt.Println("💡 User Namespace 效果：可以重新映射用户 ID")
}

// demonstrateCgroup 演示 Cgroup
func demonstrateCgroup(rm *ResourceManager) {
	fmt.Println("=== Cgroup 演示 ===")
	
	if !isLinux() {
		fmt.Println("❌ Cgroup 需要 Linux 系统")
		return
	}
	
	if !isRoot() {
		fmt.Println("❌ 需要 root 权限")
		return
	}
	
	// 检查 cgroup 文件系统
	if _, err := os.Stat("/sys/fs/cgroup"); os.IsNotExist(err) {
		fmt.Println("❌ Cgroup 文件系统不存在")
		return
	}
	
	fmt.Println("✅ Cgroup 文件系统已挂载")
	
	// 显示当前进程的 cgroup 信息
	if content, err := ioutil.ReadFile("/proc/self/cgroup"); err == nil {
		fmt.Printf("📍 当前进程的 Cgroup 信息：\n%s\n", string(content))
	}
	
	// 创建演示 cgroup
	cgroupName := "docker-demo"
	cgroupPath := fmt.Sprintf("/sys/fs/cgroup/memory/%s", cgroupName)
	
	fmt.Printf("🚀 创建 Cgroup: %s\n", cgroupPath)
	if err := os.MkdirAll(cgroupPath, 0755); err != nil {
		fmt.Printf("❌ 创建失败: %v\n", err)
		return
	}
	
	rm.AddCgroup(cgroupPath)
	
	// 设置内存限制
	memoryLimit := "100000000" // 100MB
	memoryLimitFile := filepath.Join(cgroupPath, "memory.limit_in_bytes")
	if err := ioutil.WriteFile(memoryLimitFile, []byte(memoryLimit), 0644); err != nil {
		fmt.Printf("❌ 设置内存限制失败: %v\n", err)
	} else {
		fmt.Printf("✅ 设置内存限制: %s bytes (%.1f MB)\n", memoryLimit, float64(100000000)/1024/1024)
	}
	
	// 读取当前内存使用
	memoryUsageFile := filepath.Join(cgroupPath, "memory.usage_in_bytes")
	if usage, err := ioutil.ReadFile(memoryUsageFile); err == nil {
		fmt.Printf("📊 当前内存使用: %s bytes\n", string(usage))
	}
	
	// 将当前进程加入 cgroup
	pid := os.Getpid()
	tasksFile := filepath.Join(cgroupPath, "cgroup.procs")
	if err := ioutil.WriteFile(tasksFile, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		fmt.Printf("❌ 将进程加入 Cgroup 失败: %v\n", err)
	} else {
		fmt.Printf("✅ 将进程 %d 加入 Cgroup\n", pid)
	}
	
	// 验证进程是否在 cgroup 中
	if content, err := ioutil.ReadFile(tasksFile); err == nil {
		fmt.Printf("📍 Cgroup 中的进程: %s\n", string(content))
	}
	
	fmt.Println("💡 Cgroup 效果：进程的资源使用被限制在设定范围内")
}

// demonstrateUnionFS 演示 UnionFS
func demonstrateUnionFS(rm *ResourceManager) {
	fmt.Println("=== UnionFS (OverlayFS) 演示 ===")
	
	if !isLinux() {
		fmt.Println("❌ UnionFS/OverlayFS 需要 Linux 系统")
		return
	}
	
	if !isRoot() {
		fmt.Println("❌ 需要 root 权限")
		return
	}
	
	// 创建演示目录结构
	baseDir := "/tmp/unionfs-demo"
	lowerDir := filepath.Join(baseDir, "lower")
	upperDir := filepath.Join(baseDir, "upper")
	workDir := filepath.Join(baseDir, "work")
	mergedDir := filepath.Join(baseDir, "merged")
	
	// 清理之前的演示
	os.RemoveAll(baseDir)
	rm.AddTempDir(baseDir)
	
	// 创建目录
	for _, dir := range []string{lowerDir, upperDir, workDir, mergedDir} {
		if err := os.MkdirAll(dir, 0755); err != nil {
			fmt.Printf("❌ 创建目录失败: %v\n", err)
			return
		}
	}
	
	// 在 lower 层创建文件
	lowerFiles := map[string]string{
		"base.txt":   "基础层文件\n",
		"config.txt": "基础配置\n",
	}
	
	fmt.Println("📁 在基础层创建文件：")
	for file, content := range lowerFiles {
		fullPath := filepath.Join(lowerDir, file)
		if err := ioutil.WriteFile(fullPath, []byte(content), 0644); err != nil {
			fmt.Printf("❌ 创建文件失败: %v\n", err)
			return
		}
		fmt.Printf("  ✅ %s\n", file)
	}
	
	// 在 upper 层创建文件（覆盖和新增）
	upperFiles := map[string]string{
		"config.txt": "修改后的配置\n", // 覆盖 lower 层的文件
		"new.txt":    "新增文件\n",  // 只在 upper 层的文件
	}
	
	fmt.Println("📁 在上层创建文件：")
	for file, content := range upperFiles {
		fullPath := filepath.Join(upperDir, file)
		if err := ioutil.WriteFile(fullPath, []byte(content), 0644); err != nil {
			fmt.Printf("❌ 创建文件失败: %v\n", err)
			return
		}
		fmt.Printf("  ✅ %s %s\n", file, map[bool]string{true: "(覆盖)", false: "(新增)"}[lowerFiles[file] != ""])
	}
	
	// 挂载 overlayfs
	fmt.Println("🚀 挂载 OverlayFS...")
	mountCmd := exec.Command("mount", "-t", "overlay", "overlay",
		"-o", fmt.Sprintf("lowerdir=%s,upperdir=%s,workdir=%s", lowerDir, upperDir, workDir),
		mergedDir)
	
	if err := mountCmd.Run(); err != nil {
		fmt.Printf("❌ 挂载失败: %v\n", err)
		return
	}
	
	rm.AddMountPoint(mergedDir)
	fmt.Println("✅ OverlayFS 挂载成功")
	
	// 显示合并后的文件
	fmt.Println("🔍 合并后的文件：")
	filepath.Walk(mergedDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}
		relPath, _ := filepath.Rel(mergedDir, path)
		content, _ := ioutil.ReadFile(path)
		fmt.Printf("📄 %s: %s", relPath, string(content))
		return nil
	})
	
	fmt.Println("💡 UnionFS 效果：基础层和上层文件合并，上层覆盖基础层同名文件")
}

// createContainerRootfs 创建容器根文件系统
func createContainerRootfs(rm *ResourceManager) {
	fmt.Println("=== 创建容器根文件系统演示 ===")
	
	if !isLinux() {
		fmt.Println("❌ 需要在 Linux 系统上创建容器文件系统")
		return
	}
	
	rootfs := "/tmp/container-rootfs-demo"
	os.RemoveAll(rootfs)
	rm.AddTempDir(rootfs)
	
	// 创建基本目录结构
	dirs := []string{"bin", "etc", "proc", "sys", "tmp", "dev", "home", "root", "usr", "var"}
	fmt.Println("📁 创建容器目录结构：")
	for _, dir := range dirs {
		dirPath := filepath.Join(rootfs, dir)
		if err := os.MkdirAll(dirPath, 0755); err != nil {
			fmt.Printf("❌ 创建目录失败: %v\n", err)
			return
		}
		fmt.Printf("  ✅ %s/\n", dir)
	}
	
	// 创建基本文件
	files := map[string]string{
		"etc/passwd": `root:x:0:0:root:/root:/bin/sh
daemon:x:1:1:daemon:/usr/sbin:/usr/sbin/nologin
nobody:x:65534:65534:nobody:/nonexistent:/usr/sbin/nologin
`,
		"etc/group": `root:x:0:
daemon:x:1:
nobody:x:65534:
`,
		"etc/hostname": "container-demo\n",
		"etc/hosts":    "127.0.0.1 localhost\n127.0.0.1 container-demo\n",
		"etc/resolv.conf": "nameserver 8.8.8.8\nnameserver 8.8.4.4\n",
		"etc/fstab": `proc /proc proc defaults 0 0
sysfs /sys sysfs defaults 0 0
tmpfs /tmp tmpfs defaults 0 0
`,
		"bin/hello": `#!/bin/sh
echo "================================"
echo "Hello from Container!"
echo "================================"
echo "Current user: $(whoami)"
echo "Current directory: $(pwd)"
echo "Current hostname: $(hostname)"
echo "Available commands in /bin:"
ls /bin
echo "Environment variables:"
env
echo "================================"
`,
		"root/.bashrc": `export PS1='container-demo:\w\$ '
alias ll='ls -la'
`,
	}
	
	fmt.Println("📄 创建容器配置文件：")
	for filePath, content := range files {
		fullPath := filepath.Join(rootfs, filePath)
		if err := ioutil.WriteFile(fullPath, []byte(content), 0644); err != nil {
			fmt.Printf("❌ 创建文件失败: %v\n", err)
			return
		}
		fmt.Printf("  ✅ %s\n", filePath)
	}
	
	// 设置可执行权限
	os.Chmod(filepath.Join(rootfs, "bin/hello"), 0755)
	
	// 创建设备文件
	fmt.Println("🔧 创建设备文件：")
	devices := map[string]uint32{
		"dev/null":  syscall.S_IFCHR | 0666,
		"dev/zero":  syscall.S_IFCHR | 0666,
		"dev/tty":   syscall.S_IFCHR | 0666,
		"dev/random": syscall.S_IFCHR | 0666,
		"dev/urandom": syscall.S_IFCHR | 0666,
	}
	
	for devicePath, mode := range devices {
		fullPath := filepath.Join(rootfs, devicePath)
		if err := syscall.Mknod(fullPath, mode, 0); err == nil {
			fmt.Printf("  ✅ %s\n", devicePath)
		} else {
			fmt.Printf("  ❌ %s: %v\n", devicePath, err)
		}
	}
	
	fmt.Printf("✅ 容器根文件系统创建完成: %s\n", rootfs)
	
	// 显示文件结构
	fmt.Println("📁 文件结构：")
	filepath.Walk(rootfs, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		relPath, _ := filepath.Rel(rootfs, path)
		if relPath == "." {
			return nil
		}
		
		if info.IsDir() {
			fmt.Printf("  📁 %s/\n", relPath)
		} else {
			fmt.Printf("  📄 %s\n", relPath)
		}
		return nil
	})
	
	fmt.Println("💡 容器根文件系统包含了基本的系统文件和配置")
}

// 主演示函数
func demonstrateDockerFeatures() {
	fmt.Println("=== Docker 容器技术完整演示 ===")
	fmt.Printf("🖥️  运行环境: %s\n", runtime.GOOS)
	fmt.Printf("👤 当前用户: %s\n", func() string {
		if isRoot() {
			return "root"
		}
		return "普通用户"
	}())
	fmt.Println()
	
	if !isLinux() {
		fmt.Printf("❌ 当前系统 %s 不支持，请使用 Linux 系统\n", runtime.GOOS)
		fmt.Println("💡 可以使用交叉编译：GOOS=linux go build -o docker_demo_linux docker_demo.go")
		return
	}
	
	if !isRoot() {
		fmt.Println("⚠️  注意：需要 root 权限才能看到完整演示")
		fmt.Println("💡 建议使用: sudo ./docker_demo")
		fmt.Println()
	}
	
	// 创建资源管理器
	rm := NewResourceManager()
	
	// 确保清理资源
	defer rm.Cleanup()
	
	// 演示各个组件
	demonstratePIDNamespace(rm)
	fmt.Println()
	
	demonstrateNetworkNamespace(rm)
	fmt.Println()
	
	demonstrateMountNamespace(rm)
	fmt.Println()
	
	demonstrateUTSNamespace(rm)
	fmt.Println()
	
	demonstrateIPCNamespace(rm)
	fmt.Println()
	
	demonstrateUserNamespace(rm)
	fmt.Println()
	
	demonstrateCgroup(rm)
	fmt.Println()
	
	demonstrateUnionFS(rm)
	fmt.Println()
	
	createContainerRootfs(rm)
	fmt.Println()
	
	// 总结
	fmt.Println("=== 总结 ===")
	fmt.Println("Docker 容器技术的核心组件：")
	fmt.Println("1. Linux Namespaces - 资源隔离")
	fmt.Println("   - PID Namespace: 进程隔离")
	fmt.Println("   - Network Namespace: 网络隔离")
	fmt.Println("   - Mount Namespace: 文件系统隔离")
	fmt.Println("   - UTS Namespace: 主机名隔离")
	fmt.Println("   - IPC Namespace: 进程间通信隔离")
	fmt.Println("   - User Namespace: 用户隔离")
	fmt.Println()
	fmt.Println("2. Cgroups - 资源限制")
	fmt.Println("   - CPU 限制")
	fmt.Println("   - 内存限制")
	fmt.Println("   - 磁盘 I/O 限制")
	fmt.Println("   - 网络带宽限制")
	fmt.Println()
	fmt.Println("3. UnionFS - 分层文件系统")
	fmt.Println("   - 基础镜像层共享")
	fmt.Println("   - 写时复制")
	fmt.Println("   - 镜像分层存储")
	fmt.Println()
	fmt.Println("4. 安全机制 - 多层防护")
	fmt.Println("   - 能力限制")
	fmt.Println("   - Seccomp 过滤")
	fmt.Println("   - AppArmor/SELinux")
	fmt.Println()
	fmt.Println("这些技术共同实现了轻量级虚拟化！")
}

func main() {
	fmt.Println("Docker 容器技术完整演示程序")
	fmt.Println("================================")
	fmt.Println("🐳 本程序展示了 Docker 容器技术的核心概念和实现原理")
	fmt.Println("⚠️  仅支持 Linux 系统，需要 root 权限")
	fmt.Println()
	
	demonstrateDockerFeatures()
	
	fmt.Println("\n演示完成！")
	fmt.Println("💡 所有创建的资源已自动清理")
	fmt.Println("💡 这个程序展示了 Docker 容器技术的核心原理")
	fmt.Println()
	fmt.Println("📚 学习建议：")
	fmt.Println("1. 了解每个技术的具体实现方式")
	fmt.Println("2. 思考这些技术如何组合实现容器化")
	fmt.Println("3. 探索 Docker 在这些基础技术之上的创新")
	fmt.Println("4. 学习容器安全和最佳实践")
}
