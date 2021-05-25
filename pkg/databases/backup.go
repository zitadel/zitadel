package databases

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api/database"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
)

func GitOpsInstantBackup(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitClient *git.Client,
	name string,
) error {
	desired, err := gitClient.ReadTree(git.DatabaseFile)
	if err != nil {
		return err
	}
	return instantBackup(monitor, k8sClient, desired, name)
}

func CrdInstantBackup(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	name string,
) error {
	desired, err := database.ReadCrd(k8sClient)
	if err != nil {
		monitor.Error(err)
		return err
	}
	return instantBackup(monitor, k8sClient, desired, name)
}

func instantBackup(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	desired *tree.Tree,
	name string,
) error {
	current := &tree.Tree{}

	query, _, _, _, _, _, err := orbdb.AdaptFunc(name, nil, false, "instantbackup")(monitor, desired, current)
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
