package databases

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api/database"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
)

func GitOpsListBackups(
	monitor mntr.Monitor,
	gitClient *git.Client,
	k8sClient kubernetes.ClientInt,
) (
	[]string,
	error,
) {
	desired, err := gitClient.ReadTree(git.DatabaseFile)
	if err != nil {
		monitor.Error(err)
		return nil, err
	}

	return listBackups(monitor, k8sClient, desired)
}

func CrdListBackups(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
) (
	[]string,
	error,
) {
	desired, err := database.ReadCrd(k8sClient)
	if err != nil {
		monitor.Error(err)
		return nil, err
	}

	return listBackups(monitor, k8sClient, desired)
}

func listBackups(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	desired *tree.Tree,
) (
	[]string,
	error,
) {
	backups, err := orbdb.BackupListFunc()(monitor, k8sClient, desired)
	if err != nil {
		monitor.Error(err)
		return nil, err
	}

	return backups, nil
}
