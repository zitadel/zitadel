package databases

import (
	"errors"
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api/database"
	"github.com/caos/zitadel/operator/api/zitadel"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
)

func GitOpsListBackups(
	monitor mntr.Monitor,
	gitClient *git.Client,
	k8sClient kubernetes.ClientInt,
) ([]string, error) {
	return listBackups(monitor, k8sClient, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.DatabaseFile)
	}, func() (*tree.Tree, error) {
		return gitClient.ReadTree(git.ZitadelFile)
	})
}

func CrdListBackups(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
) ([]string, error) {
	return listBackups(monitor, k8sClient,
		func() (*tree.Tree, error) {
			return database.ReadCrd(k8sClient)
		}, func() (*tree.Tree, error) {
			return zitadel.ReadCrd(k8sClient)
		})
}

func listBackups(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	databaseTree func() (*tree.Tree, error),
	zitadelTree func() (*tree.Tree, error),
) (
	[]string,
	error,
) {

	dbTree, err := databaseTree()
	if err != nil {
		return nil, err
	}
	if dbTree == nil || dbTree.Original == nil {
		return nil, errors.New("backups and restores are only supported for managed databases, but found no specs")
	}

	backups, err := orbdb.BackupListFunc()(monitor, k8sClient, dbTree)
	if err != nil {
		return nil, err
	}

	return backups, nil
}
