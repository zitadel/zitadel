package orb

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/orb"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/kinds/iam"
	"github.com/pkg/errors"
)

func AdaptFunc(orbconfig *orb.Orb, action string, features []string) operator.AdaptFunc {
	return func(monitor mntr.Monitor, desiredTree *tree.Tree, currentTree *tree.Tree) (queryFunc operator.QueryFunc, destroyFunc operator.DestroyFunc, err error) {
		defer func() {
			err = errors.Wrapf(err, "building %s failed", desiredTree.Common.Kind)
		}()

		orbMonitor := monitor.WithField("kind", "orb")

		desiredKind, err := parseDesiredV0(desiredTree)
		if err != nil {
			return nil, nil, errors.Wrap(err, "parsing desired state failed")
		}
		desiredTree.Parsed = desiredKind
		currentTree = &tree.Tree{}

		if desiredKind.Spec.Verbose && !orbMonitor.IsVerbose() {
			orbMonitor = orbMonitor.Verbose()
		}

		/* TODO: self-reconciling
		query := operator.EnsureFuncToQueryFunc(func(k8sClient *kubernetes.Client) error {
			if err := kubernetes.EnsureZitadelArtifacts(monitor, k8sClient, desiredKind.Spec.Version, desiredKind.Spec.NodeSelector, desiredKind.Spec.Tolerations); err != nil {
				monitor.Error(errors.Wrap(err, "Failed to deploy zitadel-operator into k8s-cluster"))
				return err
			}
			return nil
		})*/

		iamCurrent := &tree.Tree{}
		queryIAM, destroyIAM, err := iam.GetQueryAndDestroyFuncs(
			orbMonitor,
			desiredKind.IAM,
			iamCurrent,
			desiredKind.Spec.NodeSelector,
			desiredKind.Spec.Tolerations,
			orbconfig,
			action,
			features,
		)
		if err != nil {
			return nil, nil, err
		}

		queriers := []operator.QueryFunc{
			//query,
			queryIAM,
		}

		destroyers := []operator.DestroyFunc{
			destroyIAM,
		}

		currentTree.Parsed = &DesiredV0{
			Common: &tree.Common{
				Kind:    "zitadel.caos.ch/Orb",
				Version: "v0",
			},
			IAM: iamCurrent,
		}

		return func(k8sClient *kubernetes.Client, _ map[string]interface{}) (operator.EnsureFunc, error) {
				queried := map[string]interface{}{}
				monitor.WithField("queriers", len(queriers)).Info("Querying")
				return operator.QueriersToEnsureFunc(monitor, true, queriers, k8sClient, queried)
			},
			func(k8sClient *kubernetes.Client) error {
				monitor.WithField("destroyers", len(queriers)).Info("Destroy")
				return operator.DestroyersToDestroyFunc(monitor, destroyers)(k8sClient)
			},
			nil
	}
}
