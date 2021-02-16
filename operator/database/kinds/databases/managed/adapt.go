package managed

import (
	"github.com/caos/zitadel/operator"
	"strconv"
	"strings"

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

func AdaptFunc(
	componentLabels *labels.Component,
	namespace string,
	timestamp string,
	nodeselector map[string]string,
	tolerations []corev1.Toleration,
	version string,
	features []string,
) func(
	monitor mntr.Monitor,
	desired *tree.Tree,
	current *tree.Tree,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	map[string]*secret.Secret,
	error,
) {

	return func(
		monitor mntr.Monitor,
		desired *tree.Tree,
		current *tree.Tree,
	) (
		operator.QueryFunc,
		operator.DestroyFunc,
		map[string]*secret.Secret,
		error,
	) {
		internalMonitor := monitor.WithField("kind", "cockroachdb")
		allSecrets := map[string]*secret.Secret{}

		desiredKind, err := parseDesiredV0(desired)
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "parsing desired state failed")
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
			return nil, nil, nil, err
		}
		addRoot, err := addUser("root")
		if err != nil {
			return nil, nil, nil, err
		}
		destroyRoot, err := deleteUser("root")
		if err != nil {
			return nil, nil, nil, err
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
			image,
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
			return nil, nil, nil, err
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

		//externalName := "cockroachdb-public." + namespaceStr + ".svc.cluster.local"
		//queryES, destroyES, err := service.AdaptFunc("cockroachdb-public", "default", labels, []service.Port{}, "ExternalName", map[string]string{}, false, "", externalName)
		//if err != nil {
		//	return nil, nil, err
		//}

		queryPDB, err := pdb.AdaptFuncToEnsure(namespace, labels.MustForName(componentLabels, pdbName), cockroachSelector, "1")
		if err != nil {
			return nil, nil, nil, err
		}

		destroyPDB, err := pdb.AdaptFuncToDestroy(namespace, pdbName)
		if err != nil {
			return nil, nil, nil, err
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

		queriers := make([]operator.QueryFunc, 0)
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

		destroyers := make([]operator.DestroyFunc, 0)
		if isFeatureDatabase {
			destroyers = append(destroyers,
				operator.ResourceDestroyToZitadelDestroy(destroyPDB),
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
					queryB, destroyB, secrets, err := backups.GetQueryAndDestroyFuncs(
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
					)
					if err != nil {
						return nil, nil, nil, err
					}

					secret.AppendSecrets(backupName, allSecrets, secrets)
					destroyers = append(destroyers, destroyB)
					queriers = append(queriers, queryB)
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
			allSecrets,
			nil
	}
}
