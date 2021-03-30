package databases

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api"
	"github.com/caos/zitadel/operator/api/database"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
)

func GitOpsListBackups(
	monitor mntr.Monitor,
	gitClient *git.Client,
) (
	[]string,
	error,
) {
	desired, err := api.ReadDatabaseYml(gitClient)
	if err != nil {
		monitor.Error(err)
		return nil, err
	}

	return listBackups(monitor, desired)
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

	return listBackups(monitor, desired)
}

func listBackups(
	monitor mntr.Monitor,
	desired *tree.Tree,
) (
	[]string,
	error,
) {
	backups, err := orbdb.BackupListFunc()(monitor, desired)
	if err != nil {
		monitor.Error(err)
		return nil, err
	}

	return backups, nil
}
