package helpers_test

import (
	"context"
	k8s_test "github.com/caos/zitadel/cmd/chore/helpers/k8s"
	"github.com/caos/zitadel/cmd/chore/helpers/orbctl"
	"github.com/caos/zitadel/cmd/chore/helpers/zitadelctl"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"github.com/onsi/gomega/gexec"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
	"time"
)

func PrefixedEnv(env string) string {
	return os.Getenv("ZITADEL_E2E_" + env)
}

type ZitadelctlGitopsCmd func(args ...string) *exec.Cmd

func ZitadelctlGitopsFunc(orbconfig string, tag string) ZitadelctlGitopsCmd {
	cmdFunc, error := zitadelctl.Command(false, true, false, tag)
	Expect(error).ToNot(HaveOccurred())
	return func(args ...string) *exec.Cmd {
		cmd := cmdFunc(context.Background())
		cmd.Args = append(cmd.Args, append([]string{"--disable-analytics", "--gitops", "--orbconfig", orbconfig}, args...)...)
		return cmd
	}
}

type OrbctlGitopsCmd func(args ...string) *exec.Cmd

func OrbctlGitopsFunc(orbconfig string, orbctlVersion string) OrbctlGitopsCmd {
	cmdFunc, error := orbctl.Command(false, true, false, orbctlVersion)
	Expect(error).ToNot(HaveOccurred())
	return func(args ...string) *exec.Cmd {
		cmd := cmdFunc(context.Background())
		cmd.Args = append(cmd.Args, append([]string{"--disable-analytics", "--gitops", "--orbconfig", orbconfig}, args...)...)
		return cmd
	}
}

func writeRemoteFile(orbctlGitops ZitadelctlGitopsCmd, remoteFile string, content []byte, env func(string) string) {
	session, err := gexec.Start(orbctlGitops("file", "patch", remoteFile, "--exact", "--value", os.Expand(string(content), env)), GinkgoWriter, GinkgoWriter)
	Expect(err).ToNot(HaveOccurred())
	Eventually(session, 1*time.Minute).Should(gexec.Exit(0))
}

func LocalToRemoteFile(orbctlGitops ZitadelctlGitopsCmd, remoteFile, localFile string, env func(string) string) {
	contentBytes, err := ioutil.ReadFile(localFile)
	Expect(err).ToNot(HaveOccurred())
	writeRemoteFile(orbctlGitops, remoteFile, contentBytes, env)
}

type AwaitReadyNodes func(count int, timeout time.Duration)

func AwaitReadyNodesFunc(kubectl k8s_test.KubectlCmd) AwaitReadyNodes {
	return func(count int, timeout time.Duration) {
		getReadyNodes := func() int {
			cmd := kubectl("get", "nodes")
			session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
			Expect(err).ToNot(HaveOccurred())
			Eventually(session, 1*time.Minute, 5*time.Second).Should(gexec.Exit(0))

			content := session.Out.Contents()
			lines := strings.Split(string(content), "\n")
			readyNodes := 0
			for _, line := range lines {
				if strings.Contains(line, "Ready") && !strings.Contains(line, "NotReady") {
					readyNodes++
				}
			}
			return readyNodes
		}
		Eventually(getReadyNodes(), timeout).Should(Equal(count))
	}
}

type ApplyFile func(file []byte)

func ApplyFileFunc(kubectl k8s_test.KubectlCmd) ApplyFile {
	return func(file []byte) {
		cmd := kubectl("apply", "-f", "-")
		cmd.Stdin = strings.NewReader(os.ExpandEnv(string(file)))

		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session, 1*time.Minute).Should(gexec.Exit(0))
	}
}

type DeleteFile func(file []byte)

func DeleteFileFunc(kubectl k8s_test.KubectlCmd) DeleteFile {
	return func(file []byte) {
		cmd := kubectl("delete", "-f", "-")
		cmd.Stdin = strings.NewReader(os.ExpandEnv(string(file)))

		session, err := gexec.Start(cmd, GinkgoWriter, GinkgoWriter)
		Expect(err).ToNot(HaveOccurred())
		Eventually(session, 1*time.Minute).Should(gexec.Exit(0))
	}
}

type AwaitCompletedPodFromJob func(file []byte, namespace, selector string, timeout time.Duration)

func AwaitCompletedPodFromJobFunc(kubectl k8s_test.KubectlCmd) AwaitCompletedPodFromJob {
	return func(file []byte, namespace, selector string, timeout time.Duration) {
		ApplyFileFunc(kubectl)(file)

		Eventually(k8s_test.CountCompletedPods(kubectl, namespace, selector), timeout, 5*time.Second).Should(Equal(int8(1)))

		DeleteFileFunc(kubectl)(file)
	}
}

type AwaitCompletedPod func(namespace, selector string, timeout time.Duration)

func AwaitCompletedPodFunc(kubectl k8s_test.KubectlCmd) AwaitCompletedPod {
	return func(namespace, selector string, timeout time.Duration) {
		Eventually(k8s_test.CountCompletedPods(kubectl, namespace, selector), timeout, 5*time.Second).Should(Equal(int8(1)))
	}
}

type AwaitReadyPods func(namespace, selector string, count int, timeout time.Duration)

func AwaitReadyPodsFunc(kubectl k8s_test.KubectlCmd) AwaitReadyPods {
	return func(namespace, selector string, count int, timeout time.Duration) {
		Eventually(k8s_test.CountReadyPods(kubectl, namespace, selector), timeout, 5*time.Second).Should(Equal(int8(count)))
	}
}

type AwaitSecret func(namespace string, name string, keys []string, timeout time.Duration)

func AwaitSecretFunc(kubectl k8s_test.KubectlCmd) AwaitSecret {
	return func(namespace string, name string, keys []string, timeout time.Duration) {
		sort.Strings(keys)
		Eventually(k8s_test.GetSecretKeysWithName(kubectl, namespace, name), timeout, 5*time.Second).Should(Equal(keys))
	}
}

type AwaitCronJobScheduled func(namespace string, name string, cron string, timeout time.Duration)

func AwaitCronJobScheduledFunc(kubectl k8s_test.KubectlCmd) AwaitCronJobScheduled {
	return func(namespace string, name string, cron string, timeout time.Duration) {
		Eventually(k8s_test.GetCronJobScheduleWithName(kubectl, namespace, name), timeout, 5*time.Second).Should(Equal(cron))
	}
}

type GetLogsOfPod func(namespace string, selector string) string

func GetLogsOfPodFunc(kubectl k8s_test.KubectlCmd) GetLogsOfPod {
	return func(namespace string, selector string) string {
		cmd := kubectl("logs", "-n", namespace, "--selector", selector)
		out, err := cmd.CombinedOutput()
		Expect(err).ToNot(HaveOccurred())

		return string(out)
	}
}

func LastVersionOfMigrations(folder string) int {
	absFolder, err := filepath.Abs(folder)
	if err != nil {
		return 0
	}

	highest := 0
	err = filepath.Walk(absFolder, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() && strings.HasPrefix(info.Name(), "V") {
			parts := strings.Split(info.Name(), "__")
			versionParts := strings.Split(parts[0], ".")
			version, err := strconv.Atoi(versionParts[1])
			if err != nil {
				return err
			}
			if highest < version {
				highest = version
			}
		}
		return nil
	})
	if err != nil {
		return 0
	}
	return highest
}

type DeleteResource func(resource string, name string)

func DeleteResourceFunc(kubectl k8s_test.KubectlCmd) DeleteResource {
	return func(resource string, name string) {
		cmd := kubectl("delete", resource, name)
		out, err := cmd.CombinedOutput()
		if strings.Contains(string(out), "\""+name+"\" not found") {
			err = nil
		}
		Expect(err).ToNot(HaveOccurred())
	}
}

type DeleteNamespacedResource func(resource string, namespace string, name string)

func DeleteNamespacedResourceFunc(kubectl k8s_test.KubectlCmd) DeleteNamespacedResource {
	return func(resource string, namespace string, name string) {
		cmd := kubectl("delete", resource, "-n", namespace, name)
		out, err := cmd.CombinedOutput()
		if strings.Contains(string(out), "\""+name+"\" not found") {
			err = nil
		}
		Expect(err).ToNot(HaveOccurred())
	}
}
