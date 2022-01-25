package orb

import (
	"errors"
	"fmt"

	cockroachdb "github.com/caos/zitadel/operator/zitadel/kinds/dbconn"

	"github.com/caos/zitadel/pkg/databases/db"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/namespace"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"

	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam"
)

const namespaceName = "caos-zitadel"

var ErrUndefinedDBConn = errors.New("desired state for database connection is undefined")

func AdaptFunc(
	action string,
	binaryVersion *string,
	gitops bool,
	features []string,
	dbClient db.Client,
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

		iamCurrent := &tree.Tree{}
		databaseConnectionCurrent := &tree.Tree{}

		currentTree.Parsed = &DesiredV0{
			Common:             tree.NewCommon("zitadel.caos.ch/Orb", "v0", false),
			IAM:                iamCurrent,
			DatabaseConnection: databaseConnectionCurrent,
		}

		if desiredKind.Spec.Verbose && !orbMonitor.IsVerbose() {
			orbMonitor = orbMonitor.Verbose()
		}

		operatorLabels := mustZITADELOperator(binaryVersion)

		queriers := make([]operator.QueryFunc, 0)

		queryNS, err := namespace.AdaptFuncToEnsure(namespaceName)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}
		/*destroyNS, err := namespace.AdaptFuncToDestroy(namespaceName)
		if err != nil {
			return nil, nil, allSecrets, err
		}*/

		configurers := make([]operator.ConfigureFunc, 0)

		queryDBConn, destroyDBConn, configureDBConn, dbConnEncryptedSecrets, dbConnExistingSecrets, migrateDBConn, err := cockroachdb.Adapt(orbMonitor, operatorLabels, desiredKind.DatabaseConnection, databaseConnectionCurrent)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}

		configurers = append(configurers, configureDBConn)
		secret.AppendSecrets("dbconn", allSecrets, dbConnEncryptedSecrets, allExisting, dbConnExistingSecrets)

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
		configurers = append(configurers, configureIAM)
		secret.AppendSecrets("", allSecrets, zitadelSecrets, allExisting, zitadelExisting)

		rec, _ := Reconcile(monitor, desiredKind.Spec, gitops)

		destroyers := make([]operator.DestroyFunc, 0)
		for _, feature := range features {
			switch feature {
			case "iam", "migration", "scaleup", "scaledown":
				queriers = append(queriers,
					operator.ResourceQueryToZitadelQuery(queryNS),
					queryIAM,
					queryDBConn,
				)
				destroyers = append(destroyers, destroyIAM, destroyDBConn)
			case "operator":
				queriers = append(queriers,
					operator.ResourceQueryToZitadelQuery(queryNS),
					operator.EnsureFuncToQueryFunc(rec),
				)
			case "dbconnection":
				queriers = append(queriers, queryDBConn)
				destroyers = append(destroyers, destroyDBConn)
			}
		}

		return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
				monitor.WithField("queriers", len(queriers)).Info("Querying")
				return operator.QueriersToEnsureFunc(monitor, true, queriers, k8sClient, queried)
			},
			func(k8sClient kubernetes.ClientInt) error {
				monitor.WithField("destroyers", len(queriers)).Info("Destroy")
				return operator.DestroyersToDestroyFunc(monitor, destroyers)(k8sClient)
			},
			func(k8sClient kubernetes.ClientInt, queried map[string]interface{}, gitops bool) error {

				for _, configurer := range configurers {
					if err := configurer(k8sClient, queried, gitops); err != nil {
						return err
					}
				}

				return nil
			},
			allSecrets,
			allExisting,
			migrateIAM || migrateDBConn,
			nil
	}
}
