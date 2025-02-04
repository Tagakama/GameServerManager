package utils

import (
	"fmt"
	"syscall"
)

// KillProcess kills a process by its PID.
func KillProcess(pid int) error {
	handle, err := syscall.OpenProcess(syscall.PROCESS_TERMINATE, false, uint32(pid))
	if err != nil {
		return fmt.Errorf("failed to get process handle: %w", err)
	} else {
		fmt.Printf("Server %d stopped\n", pid)
	}
	defer syscall.CloseHandle(handle)

	return syscall.TerminateProcess(handle, 0)
}
