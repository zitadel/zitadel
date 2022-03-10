package orb

import (
	"fmt"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/namespace"
	"github.com/caos/orbos/pkg/orb"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"

	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam"
	zitadeldb "github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/database"
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
		configureFunc operator.ConfigureFunc,
		allSecrets map[string]*secret.Secret,
		allExisting map[string]*secret.Existing,
		migrate bool,
		err error,
	) {
		defer func() {
			if err != nil {
				err = fmt.Errorf("building %s failed: %w", desiredTree.Common.Kind, err)
			}
		}()

		allSecrets = make(map[string]*secret.Secret)
		allExisting = make(map[string]*secret.Existing)

		orbMonitor := monitor.WithField("kind", "orb")

		desiredKind, err := ParseDesiredV0(desiredTree)
		if err != nil {
			return nil, nil, nil, nil, nil, false, fmt.Errorf("parsing desired state failed: %w", err)
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
				return nil, nil, nil, nil, nil, false, err
			}
			dbClient = dbClientT
		} else {
			dbClient = zitadeldb.NewCrdClient(monitor)
		}

		operatorLabels := mustZITADELOperator(binaryVersion)

		queryNS, err := namespace.AdaptFuncToEnsure(namespaceName)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}
		/*destroyNS, err := namespace.AdaptFuncToDestroy(namespaceName)
		if err != nil {
			return nil, nil, allSecrets, err
		}*/

		iamCurrent := &tree.Tree{}
		queryIAM, destroyIAM, configureIAM, zitadelSecrets, zitadelExisting, migrateIAM, err := iam.Adapt(
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
			desiredKind.Spec.CustomImageRegistry,
		)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}
		migrate = migrate || migrateIAM
		secret.AppendSecrets("", allSecrets, zitadelSecrets, allExisting, zitadelExisting)

		rec, _ := Reconcile(monitor, desiredKind.Spec, gitops)

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
					operator.EnsureFuncToQueryFunc(rec),
				)
			}
		}

		currentTree.Parsed = &DesiredV0{
			Common: tree.NewCommon("zitadel.caos.ch/Orb", "v0", false),
			IAM:    iamCurrent,
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
			func(k8sClient kubernetes.ClientInt, queried map[string]interface{}, gitops bool) error {
				return configureIAM(k8sClient, queried, gitops)
			},
			allSecrets,
			allExisting,
			migrate,
			nil
	}
}
