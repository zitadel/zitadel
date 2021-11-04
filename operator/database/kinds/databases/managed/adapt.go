package managed

import (
	"fmt"
	"github.com/caos/orbos/pkg/kubernetes/resources/cronjob"
	"strconv"
	"strings"

	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/pdb"
	secretK8s "github.com/caos/orbos/pkg/kubernetes/resources/secret"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/tree"

	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/common"
	"github.com/caos/zitadel/operator/database/kinds/backups"
	coreBackup "github.com/caos/zitadel/operator/database/kinds/backups/core"
	"github.com/caos/zitadel/operator/database/kinds/databases/core"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/certificate"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/rbac"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/services"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/statefulset"
)

const (
	SfsName            = "cockroachdb"
	pdbName            = SfsName + "-budget"
	serviceAccountName = SfsName
	PublicServiceName  = SfsName + "-public"
	privateServiceName = SfsName
	cockroachPort      = int32(26257)
	cockroachHTTPPort  = int32(8080)
	Clean              = "clean"
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
		_ operator.QueryFunc,
		_ operator.DestroyFunc,
		_ operator.ConfigureFunc,
		_ map[string]*secret.Secret,
		_ map[string]*secret.Existing,
		migrate bool,
		err error,
	) {

		defer func() {
			if err != nil {
				err = fmt.Errorf("adapting managed database failed: %w", err)
			}
		}()

		var (
			internalMonitor = monitor.WithField("kind", "cockroachdb")
			allSecrets      = make(map[string]*secret.Secret)
			allExisting     = make(map[string]*secret.Existing)
		)

		desiredKind, err := parseDesiredV0(desired)
		if err != nil {
			return nil, nil, nil, nil, nil, false, fmt.Errorf("parsing desired state failed: %w", err)
		}
		desired.Parsed = desiredKind

		storageCapacity, err := resource.ParseQuantity("5G")
		if err != nil {
			return nil, nil, nil, nil, nil, false, mntr.ToUserError(fmt.Errorf("parsing storage capacity format failed: %w", err))
		}
		if desiredKind.Spec.StorageCapacity != "" {
			storageCapacityT, err := resource.ParseQuantity(desiredKind.Spec.StorageCapacity)
			if err != nil {
				return nil, nil, nil, nil, nil, false, mntr.ToUserError(fmt.Errorf("parsing storage capacity format failed: %w", err))
			}
			storageCapacity = storageCapacityT
		}

		if !monitor.IsVerbose() && desiredKind.Spec.Verbose {
			internalMonitor.Verbose()
		}

		var (
			isFeatureDatabase bool
			isFeatureClean    bool
		)
		for _, feature := range features {
			switch feature {
			case "database":
				isFeatureDatabase = true
			case Clean:
				isFeatureClean = true
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
			storageCapacity,
			cockroachPort,
			cockroachHTTPPort,
			desiredKind.Spec.StorageClass,
			desiredKind.Spec.NodeSelector,
			desiredKind.Spec.Tolerations,
			desiredKind.Spec.Resources,
			desiredKind.Spec.Cache,
			desiredKind.Spec.MaxSQLMemory,
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
			Common: tree.NewCommon("databases.caos.ch/CockroachDB", "v0", false),
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
			destroyers = append(destroyers,
				destroyS,
				operator.ResourceDestroyToZitadelDestroy(destroySFS),
				destroyRBAC,
				destroyCert,
				destroyRoot,
			)
		}

		if isFeatureClean {
			queriers = append(queriers,
				operator.ResourceQueryToZitadelQuery(
					statefulset.CleanPVCs(
						monitor,
						namespace,
						cockroachSelectabel,
						desiredKind.Spec.ReplicaCount,
					),
				),
				operator.EnsureFuncToQueryFunc(ensureInit),
				operator.EnsureFuncToQueryFunc(checkDBReady),
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
						PublicServiceName,
						cockroachPort,
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

				backupDefs := map[string]*tree.Tree{}
				if desiredKind.Spec.Backups != nil {
					backupDefs = desiredKind.Spec.Backups
				}
				cleanupDestroy, err := cleanup(
					monitor,
					backupDefs,
					k8sClient,
					namespace,
					getBackupLabels(),
				)
				if err != nil {
					return nil, err
				}
				queriers = append(queriers, func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
					return func(k8sClient kubernetes.ClientInt) error {
						return cleanupDestroy(k8sClient)
					}, nil
				})

				ensure, err := operator.QueriersToEnsureFunc(internalMonitor, true, queriers, k8sClient, queried)
				return ensure, err
			},
			func(k8sClient kubernetes.ClientInt) error {
				cleanupDestroy, err := cleanup(
					monitor,
					map[string]*tree.Tree{},
					k8sClient,
					namespace,
					getBackupLabels(),
				)
				if err != nil {
					return err
				}
				destroyers = append(destroyers, cleanupDestroy)
				return operator.DestroyersToDestroyFunc(internalMonitor, destroyers)(k8sClient)
			},
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

func cleanup(
	monitor mntr.Monitor,
	backupDefs map[string]*tree.Tree,
	k8sClient kubernetes.ClientInt,
	namespace string,
	joblabels map[string]string,
) (
	operator.DestroyFunc,
	error,
) {
	names := make([]string, 0)
	for name := range backupDefs {
		names = append(names, name)
	}

	list, err := k8sClient.ListCronJobs(namespace, joblabels)
	if err != nil {
		return nil, err
	}

	destroyers := make([]operator.DestroyFunc, 0)
	for _, cj := range list.Items {
		backupName := coreBackup.TrimBackupJobName(cj.Name)
		found := false
		for _, name := range names {
			if name == backupName {
				found = true
			}
		}
		if found {
			continue
		}

		destroyCJ, err := cronjob.AdaptFuncToDestroy(namespace, cj.Name)
		if err != nil {
			return nil, err
		}

		destroySecret, err := secretK8s.AdaptFuncToDestroy(namespace, coreBackup.GetSecretName(backupName))
		if err != nil {
			return nil, err
		}

		destroyers = append(destroyers,
			operator.ResourceDestroyToZitadelDestroy(destroyCJ),
			operator.ResourceDestroyToZitadelDestroy(destroySecret),
		)
	}

	return operator.DestroyersToDestroyFunc(monitor, destroyers), nil
}

func getBackupLabels() map[string]string {
	return map[string]string{
		"app.kubernetes.io/component":  backups.Component,
		"app.kubernetes.io/managed-by": "database.caos.ch",
		"app.kubernetes.io/part-of":    "ZITADEL",
	}
}
