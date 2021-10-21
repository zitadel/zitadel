package zitadelctl

import (
	"context"
	"errors"
	"fmt"
	"os/exec"
	"runtime"
)

func Command(debug, reuse, download bool, downloadTag string) (func(context.Context) *exec.Cmd, error) {

	if debug && download {
		return nil, errors.New("debug and download parameters can't both be true")
	}

	bin := zitadelctlPath()

	if reuse {
		return runCmd(debug, bin), nil
	}

	if !download {
		if err := buildExecutables(debug); err != nil {
			return func(context.Context) *exec.Cmd { return nil }, fmt.Errorf("building executables failed: %w", err)
		}
		return runCmd(debug, bin), nil
	}

	if err := downloadZitadelctl(bin, downloadTag); err != nil {
		return nil, fmt.Errorf("downloading zitadelctl release failed: %w", err)
	}

	return runCmd(debug, bin), nil
}

func runCmd(debug bool, zitadelctlPath string) func(context.Context) *exec.Cmd {

	return func(ctx context.Context) *exec.Cmd {
		if debug {
			return exec.CommandContext(ctx, "dlv", "exec", "--api-version", "2", "--headless", "--listen", "127.0.0.1:2345", zitadelctlPath, "--")
		}
		return exec.CommandContext(ctx, zitadelctlPath)
	}
}

func zitadelctlPath() string {
	var extension string

	if runtime.GOOS == "windows" {
		extension = ".exe"
	}

	return fmt.Sprintf("./artifacts/zitadelctl-%s-%s%s", runtime.GOOS, runtime.GOARCH, extension)
}
