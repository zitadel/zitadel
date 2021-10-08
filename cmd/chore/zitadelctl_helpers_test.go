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
	"sort"
	"strings"
	"time"
)

func prefixedEnv(env string) string {
	return os.Getenv("ZITADEL_E2E_" + env)
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

type kubectlCmd func(...string) *exec.Cmd

func kubectlCmdFunc(kubectlPath string) kubectlCmd {
	return func(args ...string) *exec.Cmd {
		return exec.Command("kubectl", append([]string{"--kubeconfig", kubectlPath}, args...)...)
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

type awaitCompletedPodFromJob func(file []byte, namespace, selector string, timeout time.Duration)

func awaitCompletedPodFromJobFunc(kubectl kubectlCmd) awaitCompletedPodFromJob {
	return func(file []byte, namespace, selector string, timeout time.Duration) {
		cmd := kubectl("apply", "-f", "-")
		cmd.Stdin = strings.NewReader(os.ExpandEnv(string(file)))

		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session, 1*time.Minute).Should(gexec.Exit(0))
		Eventually(countCompletedPods(kubectl, namespace, selector), timeout, 5*time.Second).Should(Equal(int8(1)))

		cmdDel := kubectl("delete", "-f", "-")
		cmdDel.Stdin = strings.NewReader(os.ExpandEnv(string(file)))

		sessionDel, err := gexec.Start(cmdDel, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(sessionDel, 1*time.Minute).Should(gexec.Exit(0))
	}
}

type awaitCompletedPod func(namespace, selector string, timeout time.Duration)

func awaitCompletedPodFunc(kubectl kubectlCmd) awaitCompletedPod {
	return func(namespace, selector string, timeout time.Duration) {
		Eventually(countCompletedPods(kubectl, namespace, selector), timeout, 5*time.Second).Should(Equal(int8(1)))
	}
}

type awaitReadyPods func(namespace, selector string, count int, timeout time.Duration)

func awaitReadyPodsFunc(kubectl kubectlCmd) awaitReadyPods {
	return func(namespace, selector string, count int, timeout time.Duration) {
		Eventually(countReadyPods(kubectl, namespace, selector), timeout, 5*time.Second).Should(Equal(int8(count)))
	}
}

type awaitSecret func(namespace string, name string, keys []string, timeout time.Duration)

func awaitSecretFunc(kubectl kubectlCmd) awaitSecret {
	return func(namespace string, name string, keys []string, timeout time.Duration) {
		sort.Strings(keys)
		Eventually(getSecretKeysWithName(kubectl, namespace, name), timeout, 5*time.Second).Should(Equal(keys))
	}
}

type awaitCronJobScheduled func(namespace string, name string, cron string, timeout time.Duration)

func awaitCronJobScheduledFunc(kubectl kubectlCmd) awaitCronJobScheduled {
	return func(namespace string, name string, cron string, timeout time.Duration) {
		Eventually(getCronJobScheduleWithName(kubectl, namespace, name), timeout, 5*time.Second).Should(Equal(cron))
	}
}
