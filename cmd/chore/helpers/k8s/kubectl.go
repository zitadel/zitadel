package k8s_test

import "os/exec"

type KubectlCmd func(...string) *exec.Cmd

func KubectlCmdFunc(kubectlPath string) KubectlCmd {
	return func(args ...string) *exec.Cmd {
		return exec.Command("kubectl", append([]string{"--kubeconfig", kubectlPath}, args...)...)
	}
}
