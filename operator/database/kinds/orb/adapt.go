package orb

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/namespace"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/orbos/pkg/treelabels"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/database/kinds/backups/bucket/backup"
	"github.com/caos/zitadel/operator/database/kinds/backups/bucket/clean"
	"github.com/caos/zitadel/operator/database/kinds/backups/bucket/restore"
	"github.com/caos/zitadel/operator/database/kinds/databases"
	"github.com/pkg/errors"
)

const (
	NamespaceStr = "caos-zitadel"
)

func OperatorSelector() *labels.Selector {
	return labels.OpenOperatorSelector("ZITADEL", "database.caos.ch")
}

func AdaptFunc(timestamp string, binaryVersion *string, gitops bool, features ...string) operator.AdaptFunc {

	return func(
		monitor mntr.Monitor,
		orbDesiredTree *tree.Tree,
		currentTree *tree.Tree,
	) (
		queryFunc operator.QueryFunc,
		destroyFunc operator.DestroyFunc,
		secrets map[string]*secret.Secret,
		existing map[string]*secret.Existing,
		migrate bool,
		err error,
	) {
		defer func() {
			err = errors.Wrapf(err, "building %s failed", orbDesiredTree.Common.Kind)
		}()

		orbMonitor := monitor.WithField("kind", "orb")

		desiredKind, err := ParseDesiredV0(orbDesiredTree)
		if err != nil {
			return nil, nil, nil, nil, migrate, errors.Wrap(err, "parsing desired state failed")
		}
		orbDesiredTree.Parsed = desiredKind
		currentTree = &tree.Tree{}

		if desiredKind.Spec.Verbose && !orbMonitor.IsVerbose() {
			orbMonitor = orbMonitor.Verbose()
		}

		queryNS, err := namespace.AdaptFuncToEnsure(NamespaceStr)
		if err != nil {
			return nil, nil, nil, nil, migrate, err
		}
		/*destroyNS, err := namespace.AdaptFuncToDestroy(NamespaceStr)
		if err != nil {
			return nil, nil, nil, err
		}*/

		databaseCurrent := &tree.Tree{}

		operatorLabels := mustDatabaseOperator(binaryVersion)

		queryDB, destroyDB, secrets, existing, migrate, err := databases.GetQueryAndDestroyFuncs(
			orbMonitor,
			desiredKind.Database,
			databaseCurrent,
			NamespaceStr,
			treelabels.MustForAPI(desiredKind.Database, operatorLabels),
			timestamp,
			desiredKind.Spec.NodeSelector,
			desiredKind.Spec.Tolerations,
			desiredKind.Spec.Version,
			features,
		)
		if err != nil {
			return nil, nil, nil, nil, migrate, err
		}

		destroyers := make([]operator.DestroyFunc, 0)
		queriers := make([]operator.QueryFunc, 0)
		for _, feature := range features {
			switch feature {
			case "database", backup.Instant, backup.Normal, restore.Instant, clean.Instant:
				queriers = append(queriers,
					operator.ResourceQueryToZitadelQuery(queryNS),
					queryDB,
				)
				destroyers = append(destroyers,
					destroyDB,
				)
			case "operator":
				queriers = append(queriers,
					operator.ResourceQueryToZitadelQuery(queryNS),
					operator.EnsureFuncToQueryFunc(Reconcile(monitor, desiredKind.Spec, gitops)),
				)
			}
		}

		currentTree.Parsed = &DesiredV0{
			Common: &tree.Common{
				Kind:    "databases.caos.ch/Orb",
				Version: "v0",
			},
			Database: databaseCurrent,
		}

		return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
				if queried == nil {
					queried = map[string]interface{}{}
				}
				monitor.WithField("queriers", len(queriers)).Info("Querying")
				return operator.QueriersToEnsureFunc(monitor, true, queriers, k8sClient, queried)
			},
			func(k8sClient kubernetes.ClientInt) error {
				monitor.WithField("destroyers", len(queriers)).Info("Destroy")
				return operator.DestroyersToDestroyFunc(monitor, destroyers)(k8sClient)
			},
			secrets,
			existing,
			migrate,
			nil
	}
}
