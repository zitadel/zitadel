//go:build integration && !windows

package integration_test

import (
	"os/exec"
	"syscall"
	"time"
)

func configureChildCommand(cmd *exec.Cmd) {
	// Put the child in its own process group so we can kill the entire tree
	// on cleanup, preventing orphaned go-test / go-build processes.
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	cmd.Cancel = func() error {
		// Send SIGTERM to the entire process group (negative PID).
		if cmd.Process != nil {
			return syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
		}
		return nil
	}
	cmd.WaitDelay = 10 * time.Second // fallback SIGKILL after 10s
}
