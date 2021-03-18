package orb

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/namespace"
	"github.com/caos/orbos/pkg/orb"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam"
	zitadeldb "github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/database"
	"github.com/pkg/errors"
)

const (
	namespaceName = "caos-zitadel"
)

func AdaptFunc(
	orbconfig *orb.Orb,
	action string,
	binaryVersion *string,
	gitops bool,
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
		allExisting map[string]*secret.Existing,
		err error,
	) {
		defer func() {
			err = errors.Wrapf(err, "building %s failed", desiredTree.Common.Kind)
		}()

		allSecrets = make(map[string]*secret.Secret)
		allExisting = make(map[string]*secret.Existing)

		orbMonitor := monitor.WithField("kind", "orb")

		desiredKind, err := ParseDesiredV0(desiredTree)
		if err != nil {
			return nil, nil, nil, nil, errors.Wrap(err, "parsing desired state failed")
		}
		desiredTree.Parsed = desiredKind
		currentTree = &tree.Tree{}

		if desiredKind.Spec.Verbose && !orbMonitor.IsVerbose() {
			orbMonitor = orbMonitor.Verbose()
		}

		var dbClient zitadeldb.Client
		if gitops {
			dbClientT, err := zitadeldb.NewGitOpsClient(monitor, orbconfig.URL, orbconfig.Repokey)
			if err != nil {
				monitor.Error(err)
				return nil, nil, nil, nil, err
			}
			dbClient = dbClientT
		} else {
			dbClient = zitadeldb.NewCrdClient(monitor)
		}

		operatorLabels := mustZITADELOperator(binaryVersion)

		queryNS, err := namespace.AdaptFuncToEnsure(namespaceName)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		/*destroyNS, err := namespace.AdaptFuncToDestroy(namespaceName)
		if err != nil {
			return nil, nil, allSecrets, err
		}*/

		iamCurrent := &tree.Tree{}
		queryIAM, destroyIAM, zitadelSecrets, zitadelExisting, err := iam.GetQueryAndDestroyFuncs(
			orbMonitor,
			operatorLabels,
			desiredKind.IAM,
			iamCurrent,
			desiredKind.Spec.NodeSelector,
			desiredKind.Spec.Tolerations,
			dbClient,
			namespaceName,
			action,
			&desiredKind.Spec.Version,
			features,
		)
		if err != nil {
			return nil, nil, nil, nil, err
		}
		secret.AppendSecrets("", allSecrets, zitadelSecrets, allExisting, zitadelExisting)

		destroyers := make([]operator.DestroyFunc, 0)
		queriers := make([]operator.QueryFunc, 0)
		for _, feature := range features {
			switch feature {
			case "iam", "migration", "scaleup", "scaledown":
				queriers = append(queriers,
					operator.ResourceQueryToZitadelQuery(queryNS),
					queryIAM,
				)
				destroyers = append(destroyers, destroyIAM)
			case "operator":
				queriers = append(queriers,
					operator.ResourceQueryToZitadelQuery(queryNS),
					operator.EnsureFuncToQueryFunc(Reconcile(monitor, desiredKind.Spec)),
				)
			}
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
			allExisting,
			nil
	}
}
