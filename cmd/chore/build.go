package chore

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

func BuildExecutables(debug, hostBinsOnly bool) error {

	cmd := exec.Command("git", "rev-parse", "--abbrev-ref", "HEAD")
	outBuf := new(bytes.Buffer)
	cmd.Stdout = outBuf
	if err := run(cmd); err != nil {
		return err
	}

	version := strings.TrimSpace(strings.Replace(outBuf.String(), "heads/", "", 1)) + "-dev"

	cmd = exec.Command("git", "rev-parse", "HEAD")
	outBuf = new(bytes.Buffer)
	cmd.Stdout = outBuf
	if err := run(cmd); err != nil {
		return err
	}

	commit := strings.TrimSpace(outBuf.String())

	files, err := filepath.Glob("./cmd/chore/gen-executables/*.go")
	if err != nil {
		panic(err)
	}
	args := []string{"run", "-race"}
	args = append(args, files...)
	args = append(args,
		"--version", version,
		"--commit", commit,
		"--githubclientid", os.Getenv("GITHUBOAUTHCLIENTID"),
		"--githubclientsecret", os.Getenv("GITHUBOAUTHCLIENTSECRET"),
		"--orbctl", "./artifacts",
		"--dev",
	)
	if debug {
		args = append(args, "--debug")
	}
	if hostBinsOnly {
		args = append(args, "--host-bins-only")
	}
	cmd = exec.Command("go", args...)
	cmd.Stdout = os.Stderr
	// gen-executables
	if err := run(cmd); err != nil || hostBinsOnly {
		// error contains --githubclientid and --githubclientsecret values
		return errors.New("building executables failed")
	}

	files, err = filepath.Glob("./cmd/chore/gen-charts/*.go")
	if err != nil {
		panic(err)
	}
	args = []string{"build", "-o", "./artifacts/gen-charts"}
	args = append(args, files...)
	cmd = exec.Command("go", args...)
	cmd.Stdout = os.Stderr
	cmd.Env = []string{"CGO_ENABLED=0", "GOOS=linux"}
	// gen-charts
	return run(cmd)
}

func run(cmd *exec.Cmd) error {
	cmd.Stderr = os.Stderr
	cmd.Env = append(cmd.Env, os.Environ()...)
	if err := cmd.Run(); err != nil {
		return fmt.Errorf("executing %s failed: %s", strings.Join(cmd.Args, " "), err.Error())
	}
	return nil
}
