package chore

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func BuildExecutables(debug bool) error {

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

	args := []string{"build", "-a"}
	args = append(args,
		"-installsuffix cgo",
		"-ldflags \"-extldflags -static -X main.Version="+version+" -X main.githubClientID="+os.Getenv("GITHUBOAUTHCLIENTID")+" -X main.githubClientSecret="+os.Getenv("GITHUBOAUTHCLIENTSECRET")+"}\"",
		"-o", "./artifacts/zitadelctl",
		"./cmd/zitadelctl/main.go",
	)
	if debug {
		args = append(args, "--debug")
	}

	cmd = exec.Command("go", args...)
	cmd.Stdout = os.Stderr
	// gen-executables
	if err := run(cmd); err != nil {
		// error contains --githubclientid and --githubclientsecret values
		return errors.New("building executables failed")
	}

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
