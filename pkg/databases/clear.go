package databases

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed"
	orbdb "github.com/caos/zitadel/operator/database/kinds/orb"
)

func Clear(
	monitor mntr.Monitor,
	k8sClient kubernetes.ClientInt,
	gitClient *git.Client,
) error {
	desired, err := gitClient.ReadTree(git.DatabaseFile)
	if err != nil {
		monitor.Error(err)
		return err
	}
	current := &tree.Tree{}

	query, _, _, _, _, _, err := orbdb.AdaptFunc("", nil, false, managed.Clean)(monitor, desired, current)
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
