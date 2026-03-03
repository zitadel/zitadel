//go:build integration && windows

package integration_test

import (
	"os/exec"
	"time"
)

func configureChildCommand(cmd *exec.Cmd) {
	cmd.Cancel = func() error {
		if cmd.Process != nil {
			return cmd.Process.Kill()
		}
		return nil
	}
	cmd.WaitDelay = 10 * time.Second
}
