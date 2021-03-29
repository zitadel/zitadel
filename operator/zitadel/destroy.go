package zitadel

import (
	"errors"
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/git"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
)

func Destroy(
	monitor mntr.Monitor,
	gitClient *git.Client,
	adapt operator.AdaptFunc,
	k8sClient *kubernetes.Client,
) func() error {
	return func() error {
		internalMonitor := monitor.WithField("operator", "zitadel")
		internalMonitor.Info("Destroy")
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

		_, destroy, _, _, _, err := adapt(internalMonitor, treeDesired, treeCurrent)
		if err != nil {
			internalMonitor.Error(err)
			return err
		}

		if err := destroy(k8sClient); err != nil {
			internalMonitor.Error(err)
			return err
		}
		return nil
	}
}
