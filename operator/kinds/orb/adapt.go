package orb

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/orb"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/kinds/iam"
	"github.com/pkg/errors"
)

func AdaptFunc(
	orbconfig *orb.Orb,
	action string,
	migrationsPath string,
	binaryVersion *string,
	features []string,
) operator.AdaptFunc {
	return func(
		monitor mntr.Monitor,
		desiredTree *tree.Tree,
		currentTree *tree.Tree,
	) (
		queryFunc operator.QueryFunc,
		destroyFunc operator.DestroyFunc,
		allSecrets map[string]*secret.Secret,
		err error,
	) {
		defer func() {
			err = errors.Wrapf(err, "building %s failed", desiredTree.Common.Kind)
		}()

		allSecrets = make(map[string]*secret.Secret)

		orbMonitor := monitor.WithField("kind", "orb")

		desiredKind, err := parseDesiredV0(desiredTree)
		if err != nil {
			return nil, nil, allSecrets, errors.Wrap(err, "parsing desired state failed")
		}
		desiredTree.Parsed = desiredKind
		currentTree = &tree.Tree{}

		if desiredKind.Spec.Verbose && !orbMonitor.IsVerbose() {
			orbMonitor = orbMonitor.Verbose()
		}

		operatorLabels := mustDatabaseOperator(binaryVersion)

		iamCurrent := &tree.Tree{}
		queryIAM, destroyIAM, zitadelSecrets, err := iam.GetQueryAndDestroyFuncs(
			orbMonitor,
			operatorLabels,
			desiredKind.IAM,
			iamCurrent,
			desiredKind.Spec.NodeSelector,
			desiredKind.Spec.Tolerations,
			orbconfig,
			action,
			migrationsPath,
			binaryVersion,
			features,
		)
		if err != nil {
			return nil, nil, allSecrets, err
		}
		secret.AppendSecrets("", allSecrets, zitadelSecrets)

		queriers := []operator.QueryFunc{
			queryIAM,
		}
		if desiredKind.Spec.SelfReconciling {
			queriers = append(queriers, operator.EnsureFuncToQueryFunc(Reconcile(monitor, desiredTree, false)))
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

		return func(k8sClient kubernetes.ClientInt, _ map[string]interface{}) (operator.EnsureFunc, error) {
				queried := map[string]interface{}{}
				monitor.WithField("queriers", len(queriers)).Info("Querying")
				return operator.QueriersToEnsureFunc(monitor, true, queriers, k8sClient, queried)
			},
			func(k8sClient kubernetes.ClientInt) error {
				monitor.WithField("destroyers", len(queriers)).Info("Destroy")
				return operator.DestroyersToDestroyFunc(monitor, destroyers)(k8sClient)
			},
			allSecrets,
			nil
	}
}
