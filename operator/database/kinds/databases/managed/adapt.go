package managed

import (
	"strconv"
	"strings"

	"github.com/caos/zitadel/operator/common"

	"github.com/caos/zitadel/operator"

	"github.com/caos/orbos/pkg/labels"

	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/certificate"

	corev1 "k8s.io/api/core/v1"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/pdb"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator/database/kinds/backups"
	"github.com/caos/zitadel/operator/database/kinds/databases/core"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/rbac"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/services"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/statefulset"
	"github.com/pkg/errors"
)

const (
	SfsName            = "cockroachdb"
	pdbName            = SfsName + "-budget"
	serviceAccountName = SfsName
	PublicServiceName  = SfsName + "-public"
	privateServiceName = SfsName
	cockroachPort      = int32(26257)
	cockroachHTTPPort  = int32(8080)
	image              = "cockroachdb/cockroach:v20.2.3"
)

func Adapter(
	componentLabels *labels.Component,
	namespace string,
	timestamp string,
	nodeselector map[string]string,
	tolerations []corev1.Toleration,
	version string,
	features []string,
	customImageRegistry string,
) operator.AdaptFunc {

	return func(
		monitor mntr.Monitor,
		desired *tree.Tree,
		current *tree.Tree,
	) (
		operator.QueryFunc,
		operator.DestroyFunc,
		operator.ConfigureFunc,
		map[string]*secret.Secret,
		map[string]*secret.Existing,
		bool,
		error,
	) {

		var (
			internalMonitor = monitor.WithField("kind", "cockroachdb")
			allSecrets      = make(map[string]*secret.Secret)
			allExisting     = make(map[string]*secret.Existing)
			migrate         bool
		)

		desiredKind, err := parseDesiredV0(desired)
		if err != nil {
			return nil, nil, nil, nil, nil, false, errors.Wrap(err, "parsing desired state failed")
		}
		desired.Parsed = desiredKind

		if !monitor.IsVerbose() && desiredKind.Spec.Verbose {
			internalMonitor.Verbose()
		}

		var (
			isFeatureDatabase bool
			isFeatureRestore  bool
		)
		for _, feature := range features {
			switch feature {
			case "database":
				isFeatureDatabase = true
			case "restore":
				isFeatureRestore = true
			}
		}

		queryCert, destroyCert, addUser, deleteUser, listUsers, err := certificate.AdaptFunc(internalMonitor, namespace, componentLabels, desiredKind.Spec.ClusterDns, isFeatureDatabase)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}
		addRoot, err := addUser("root")
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}
		destroyRoot, err := deleteUser("root")
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}

		queryRBAC, destroyRBAC, err := rbac.AdaptFunc(internalMonitor, namespace, labels.MustForName(componentLabels, serviceAccountName))

		cockroachNameLabels := labels.MustForName(componentLabels, SfsName)
		cockroachSelector := labels.DeriveNameSelector(cockroachNameLabels, false)
		cockroachSelectabel := labels.AsSelectable(cockroachNameLabels)
		querySFS, destroySFS, ensureInit, checkDBReady, listDatabases, err := statefulset.AdaptFunc(
			internalMonitor,
			cockroachSelectabel,
			cockroachSelector,
			desiredKind.Spec.Force,
			namespace,
			common.CockroachImage.Reference(customImageRegistry),
			serviceAccountName,
			desiredKind.Spec.ReplicaCount,
			desiredKind.Spec.StorageCapacity,
			cockroachPort,
			cockroachHTTPPort,
			desiredKind.Spec.StorageClass,
			desiredKind.Spec.NodeSelector,
			desiredKind.Spec.Tolerations,
			desiredKind.Spec.Resources,
		)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}

		queryS, destroyS, err := services.AdaptFunc(
			internalMonitor,
			namespace,
			labels.MustForName(componentLabels, PublicServiceName),
			labels.MustForName(componentLabels, privateServiceName),
			cockroachSelector,
			cockroachPort,
			cockroachHTTPPort,
		)

		queryPDB, err := pdb.AdaptFuncToEnsure(namespace, labels.MustForName(componentLabels, pdbName), cockroachSelector, "1")
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}

		currentDB := &Current{
			Common: &tree.Common{
				Kind:    "databases.caos.ch/CockroachDB",
				Version: "v0",
			},
			Current: &CurrentDB{
				CA: &certificate.Current{},
			},
		}
		current.Parsed = currentDB

		var (
			queriers    = make([]operator.QueryFunc, 0)
			destroyers  = make([]operator.DestroyFunc, 0)
			configurers = make([]operator.ConfigureFunc, 0)
		)
		if isFeatureDatabase {
			queriers = append(queriers,
				queryRBAC,
				queryCert,
				addRoot,
				operator.ResourceQueryToZitadelQuery(querySFS),
				operator.ResourceQueryToZitadelQuery(queryPDB),
				queryS,
				operator.EnsureFuncToQueryFunc(ensureInit),
			)
		}

		if isFeatureDatabase {
			destroyers = append(destroyers,
				destroyS,
				operator.ResourceDestroyToZitadelDestroy(destroySFS),
				destroyRBAC,
				destroyCert,
				destroyRoot,
			)
		}

		if desiredKind.Spec.Backups != nil {

			oneBackup := false
			for backupName := range desiredKind.Spec.Backups {
				if timestamp != "" && strings.HasPrefix(timestamp, backupName) {
					oneBackup = true
				}
			}

			for backupName, desiredBackup := range desiredKind.Spec.Backups {
				currentBackup := &tree.Tree{}
				if timestamp == "" || !oneBackup || (timestamp != "" && strings.HasPrefix(timestamp, backupName)) {
					queryB, destroyB, configureB, secrets, existing, migrateB, err := backups.Adapt(
						internalMonitor,
						desiredBackup,
						currentBackup,
						backupName,
						namespace,
						componentLabels,
						checkDBReady,
						strings.TrimPrefix(timestamp, backupName+"."),
						nodeselector,
						tolerations,
						version,
						features,
						customImageRegistry,
					)
					if err != nil {
						return nil, nil, nil, nil, nil, false, err
					}

					migrate = migrate || migrateB

					secret.AppendSecrets(backupName, allSecrets, secrets, allExisting, existing)
					destroyers = append(destroyers, destroyB)
					queriers = append(queriers, queryB)
					configurers = append(configurers, configureB)
				}
			}
		}

		return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
				if !isFeatureRestore {
					queriedCurrentDB, err := core.ParseQueriedForDatabase(queried)
					if err != nil || queriedCurrentDB == nil {
						// TODO: query system state
						currentDB.Current.Port = strconv.Itoa(int(cockroachPort))
						currentDB.Current.URL = PublicServiceName
						currentDB.Current.ReadyFunc = checkDBReady
						currentDB.Current.AddUserFunc = addUser
						currentDB.Current.DeleteUserFunc = deleteUser
						currentDB.Current.ListUsersFunc = listUsers
						currentDB.Current.ListDatabasesFunc = listDatabases

						core.SetQueriedForDatabase(queried, current)
						internalMonitor.Info("set current state of managed database")
					}
				}

				ensure, err := operator.QueriersToEnsureFunc(internalMonitor, true, queriers, k8sClient, queried)
				return ensure, err
			},
			operator.DestroyersToDestroyFunc(internalMonitor, destroyers),
			func(k8sClient kubernetes.ClientInt, queried map[string]interface{}, gitops bool) error {
				for i := range configurers {
					if err := configurers[i](k8sClient, queried, gitops); err != nil {
						return err
					}
				}
				return nil
			},
			allSecrets,
			allExisting,
			migrate,
			nil
	}
}
