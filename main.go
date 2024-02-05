//go:build linux

package main

import (
	"log"
	"os"
	"os/exec"
	"syscall"
)

func main() {
	cmd := exec.Command("sh")
	// 7种namespace作用
	// CLONE_NEWUTS: new nodename domainname namespace
	// CLONE_NEWIPC: new ipc资源 namespace
	// CLONE_NEWPID: new pid 视图 namespace
	// CLONE_NEWNS: new mount namespace
	// CLONE_NEWUSER: new user namespace
	// CLONE_NEWNET: new network namespace
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
}
