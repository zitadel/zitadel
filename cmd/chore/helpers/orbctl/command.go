package orbctl

import (
	"context"
	"errors"
	"fmt"
	"gopkg.in/yaml.v3"
	"os/exec"
	"runtime"
	"strings"
)

type NetworkingYML struct {
	Spec struct {
		Version string
	}
}

func GetVersion(data []byte) string {
	y := &NetworkingYML{}
	if err := yaml.Unmarshal(data, y); err != nil {
		return ""
	}
	return y.Spec.Version
}

func Command(debug, reuse, download bool, downloadTag string) (func(context.Context) *exec.Cmd, error) {

	if debug && download {
		return nil, errors.New("debug and download parameters can't both be true")
	}

	bin := orbctlPath()

	if reuse {
		return runOrbctlCmd(debug, bin), nil
	}

	if !download {
		if err := buildExecutables(debug, false); err != nil {
			return func(context.Context) *exec.Cmd { return nil }, fmt.Errorf("building executables failed: %w", err)
		}
		return runOrbctlCmd(debug, bin), nil
	}

	if err := downloadOrbctl(bin, downloadTag); err != nil {
		return nil, fmt.Errorf("downloading orbctl release failed: %w", err)
	}

	return runOrbctlCmd(debug, bin), nil
}

func runOrbctlCmd(debug bool, orbctlPath string) func(context.Context) *exec.Cmd {

	return func(ctx context.Context) *exec.Cmd {
		if debug {
			return exec.CommandContext(ctx, "dlv", "exec", "--api-version", "2", "--headless", "--listen", "127.0.0.1:2345", orbctlPath, "--")
		}
		return exec.CommandContext(ctx, orbctlPath)
	}
}

func orbctlPath() string {
	var extension string

	if runtime.GOOS == "windows" {
		extension = ".exe"
	}

	return fmt.Sprintf("./artifacts/orbctl-%s-x86_64%s", strings.ToUpper(runtime.GOOS[0:1])+runtime.GOOS[1:], extension)
}
