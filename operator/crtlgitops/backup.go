package crtlgitops

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/pkg/databases"
)

func Restore(
	monitor mntr.Monitor,
	gitClient *git.Client,
	k8sClient *kubernetes.Client,
	backup string,
) error {
	if err := databases.GitOpsClear(monitor, k8sClient, gitClient); err != nil {
		return err
	}

	if err := databases.GitOpsRestore(
		monitor,
		k8sClient,
		gitClient,
		backup,
	); err != nil {
		return err
	}
	return nil
}
