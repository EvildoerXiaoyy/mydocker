//go:build linux

package main

import (
	"fmt"
	"os"
	"os/exec"
	"path"
	"strconv"
	"syscall"
)

const cgroupMemoryHierarchyMount = "/sys/fs/cgroup/"

func main() {
	// 7种namespace作用
	// CLONE_NEWUTS: new nodename domainname namespace
	// CLONE_NEWIPC: new ipc资源 namespace
	// CLONE_NEWPID: new pid 视图 namespace
	// CLONE_NEWNS: new mount namespace
	// CLONE_NEWUSER: new user namespace
	// CLONE_NEWNET: new network namespace

	/* 以下为namespace 学习调试代码
	cmd := exec.Command("sh")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWIPC | syscall.CLONE_NEWPID |
			syscall.CLONE_NEWNS | syscall.CLONE_NEWUSER | syscall.CLONE_NEWNET,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	os.Stderr = os.Stderr

	err := cmd.Run
	if err != nil {
		log.Fatal(err())
	}
	os.Exit(-1)
	*/
	if os.Args[0] == "/proc/self/exe" {
		// 容器进程
		fmt.Printf("current pid %d", syscall.Getpid())
		fmt.Println()
		cmd := exec.Command("sh", "-c", `stress --vm-bytes 200m --vm-keep -m 1`)
		cmd.SysProcAttr = &syscall.SysProcAttr{}
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			fmt.Println(err)
			os.Exit(-1)
		}
	}
	cmd := exec.Command("/proc/self/exe")
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Cloneflags: syscall.CLONE_NEWUTS | syscall.CLONE_NEWPID | syscall.CLONE_NEWNS,
	}
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Start(); err != nil {
		fmt.Println("ERROR", err)
		os.Exit(-1)
	} else {
		//得到fork出来进程映射在外部命令空间的pid
		fmt.Printf("%v", cmd.Process.Pid)

		// 在系统默认创建挂载了 memory subsystem的Hierarchy上创建cgroup
		os.Mkdir(path.Join(cgroupMemoryHierarchyMount, "testmomorylimit"), 0755)
		os.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testmomorylimit", "cgroup.procs"),
			[]byte(strconv.Itoa(cmd.Process.Pid)), 0644)
		os.WriteFile(path.Join(cgroupMemoryHierarchyMount, "testmomorylimit", "memory.high"),
			[]byte("100m"), 0644)
	}
	cmd.Process.Wait()
}
