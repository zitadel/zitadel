package zitadel

import (
	"errors"

	"github.com/caos/zitadel/operator"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
)

func Takeoff(
	monitor mntr.Monitor,
	gitClient *git.Client,
	adapt operator.AdaptFunc,
	k8sClient *kubernetes.Client,
) func() error {
	return func() error {
		internalMonitor := monitor.WithField("operator", "zitadel")
		internalMonitor.Info("Takeoff")
		treeDesired, err := operator.Parse(gitClient, "zitadel.yml")
		if err != nil {
			monitor.Error(err)
			return err
		}
		treeCurrent := &tree.Tree{}

		if !k8sClient.Available() {
			internalMonitor.Error(errors.New("kubeclient is not available"))
			return err
		}

		query, _, _, _, _, err := adapt(internalMonitor, treeDesired, treeCurrent)
		if err != nil {
			internalMonitor.Error(err)
			return err
		}

		ensure, err := query(k8sClient, map[string]interface{}{})
		if err != nil {
			internalMonitor.Error(err)
			return err
		}

		if err := ensure(k8sClient); err != nil {
			internalMonitor.Error(err)
			return err
		}
		return nil
	}
}
