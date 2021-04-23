package databases

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/api/database"
	"github.com/caos/zitadel/operator/database/kinds/databases/core"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
)

func GitOpsClear(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitClient *git.Client,
	databases []string,
	users []string,
) error {
	desired, err := gitClient.ReadTree(git.DatabaseFile)
	if err != nil {
		monitor.Error(err)
		return err
	}

	return clear(monitor, k8sClient, databases, users, desired)
}

func CrdClear(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	databases []string,
	users []string,
) error {
	desired, err := database.ReadCrd(k8sClient)
	if err != nil {
		return err
	}

	return clear(monitor, k8sClient, databases, users, desired)
}

func clear(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	databases []string,
	users []string,
	desired *tree.Tree,
) error {
	current := &tree.Tree{}

	query, _, _, _, _, _, err := orbdb.AdaptFunc("", nil, false, "clean")(monitor, desired, current)
	if err != nil {
		monitor.Error(err)
		return err
	}
	queried := map[string]interface{}{}
	core.SetQueriedForDatabaseDBList(queried, databases, users)

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
