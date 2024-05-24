package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"runtime"
	"syscall"
)

const (
	// https://filippo.io/linux-syscall-table/
	PRCTL_SYSCALL = 157

	// see /usr/include/linux/prctl.h
	PR_SET_THP_DISABLE = 41
	PR_GET_THP_DISABLE = 42
)

func disableTHP() {
	_, _, errno := syscall.RawSyscall6(uintptr(PRCTL_SYSCALL), uintptr(PR_SET_THP_DISABLE), uintptr(1), 0, 0, 0, 0)
	if errno != 0 {
		log.Fatalf("failed to disable THP: %v", errno)
	}
}
func isTHPDisabled() bool {
	s, _, errno := syscall.RawSyscall6(uintptr(PRCTL_SYSCALL), uintptr(PR_GET_THP_DISABLE), 0, 0, 0, 0, 0)
	if errno != 0 {
		log.Fatalf("failed get THP disable status: %v", errno)
	}
	return s == 1
}

func main() {
	if os.Getenv("X_IN_CHILD_PROC") == "" {
		runtime.LockOSThread()
		disableTHP()
		cmd := exec.Command(os.Args[0], os.Args...)
		cmd.Env = append(os.Environ(), fmt.Sprintf("X_IN_CHILD_PROC=yes"))
		cmd.Stdin, cmd.Stdout, cmd.Stderr = os.Stdin, os.Stdout, os.Stderr
		_ = cmd.Run() // err discarded
		os.Exit(cmd.ProcessState.ExitCode())
	}

	// code below has THP disabled

	if !isTHPDisabled() { // sanity check
		log.Fatal("THP is somehow not disabled")
		
	}

	fmt.Println("THP is disabled")
	fmt.Println("hello world")
}
