//go:build integration && windows

package integration_test

import (
	"os/exec"
	"time"
)

func configureChildCommand(cmd *exec.Cmd) {
	// Windows does not support Unix-style process groups in the same way, so
	// we fall back to terminating the child process directly.
	cmd.Cancel = func() error {
		if cmd.Process != nil {
			return cmd.Process.Kill()
		}
		return nil
	}
	cmd.WaitDelay = 10 * time.Second
}
