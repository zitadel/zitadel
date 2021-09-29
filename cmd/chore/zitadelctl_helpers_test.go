package chore_test

import (
	"context"
	"github.com/caos/zitadel/cmd/chore"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os"
	"os/exec"
	"time"
)

func prefixedEnv(env string) string {
	return os.Getenv("ORBOS_E2E_" + env)
}

type zitadelctlGitopsCmd func(args ...string) *exec.Cmd

func zitadelctlGitopsFunc(orbconfig string) zitadelctlGitopsCmd {
	cmdFunc, error := chore.Command(false, true, false, "")
	Expect(error).ToNot(HaveOccurred())
	return func(args ...string) *exec.Cmd {
		cmd := cmdFunc(context.Background())
		cmd.Args = append(cmd.Args, append([]string{"--disable-analytics", "--gitops", "--orbconfig", orbconfig}, args...)...)
		return cmd
	}
}

func writeRemoteFile(orbctlGitops zitadelctlGitopsCmd, remoteFile string, content []byte, env func(string) string) {
	session, err := gexec.Start(orbctlGitops("file", "patch", remoteFile, "--exact", "--value", os.Expand(string(content), env)), GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred())
	Eventually(session, 1*time.Minute).Should(gexec.Exit(0))
}

func localToRemoteFile(orbctlGitops zitadelctlGitopsCmd, remoteFile, localFile string, env func(string) string) {
	contentBytes, err := ioutil.ReadFile(localFile)
	Expect(err).ToNot(HaveOccurred())
	writeRemoteFile(orbctlGitops, remoteFile, contentBytes, env)
}
