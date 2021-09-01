package backup

import (
	"os/exec"
	"syscall"
)

func runCommand(cmd *exec.Cmd) error {
	errs := make(chan error)
	go func() {
		if err := cmd.Start(); err != nil {
			errs <- err
			return
		}

		if err := cmd.Wait(); err != nil {
			if exiterr, ok := err.(*exec.ExitError); ok {
				if _, ok := exiterr.Sys().(syscall.WaitStatus); ok {
					errs <- err
					return
				}
			} else {
				errs <- err
				return
			}
		}
		errs <- nil
	}()

	return <-errs
}
