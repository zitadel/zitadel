package databases

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
)

func InstantBackup(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitClient *git.Client,
	name string,
) error {
	desired, err := api.ReadDatabaseYml(gitClient)
	if err != nil {
		monitor.Error(err)
		return err
	}
	current := &tree.Tree{}

	query, _, _, err := orbdb.AdaptFunc(name, nil, "instantbackup")(monitor, desired, current)
	if err != nil {
		monitor.Error(err)
		return err
	}

	queried := map[string]interface{}{}
	ensure, err := query(k8sClient, queried)
	if err != nil {
		monitor.Error(err)
		return err
	}

	if err := ensure(k8sClient); err != nil {
		monitor.Error(err)
		return err
	}
	return nil
}

func ListBackups(
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

	backups, err := orbdb.BackupListFunc()(monitor, desired)
	if err != nil {
		monitor.Error(err)
		return nil, err
	}

	return backups, nil
}
