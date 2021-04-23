package bucket

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/secret"
	"github.com/caos/orbos/pkg/labels"
	secretpkg "github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/secret/read"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/database/kinds/backups/bucket/backup"
	"github.com/caos/zitadel/operator/database/kinds/backups/bucket/clean"
	"github.com/caos/zitadel/operator/database/kinds/backups/bucket/restore"
	coreDB "github.com/caos/zitadel/operator/database/kinds/databases/core"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

const (
	secretName = "backup-serviceaccountjson"
	secretKey  = "serviceaccountjson"
)

func AdaptFunc(
	name string,
	namespace string,
	componentLabels *labels.Component,
	checkDBReady operator.EnsureFunc,
	timestamp string,
	nodeselector map[string]string,
	tolerations []corev1.Toleration,
	version string,
	features []string,
) operator.AdaptFunc {
	return func(
		monitor mntr.Monitor,
		desired *tree.Tree,
		current *tree.Tree,
	) (
		operator.QueryFunc,
		operator.DestroyFunc,
		operator.ConfigureFunc,
		map[string]*secretpkg.Secret,
		map[string]*secretpkg.Existing,
		bool,
		error,
	) {

		internalMonitor := monitor.WithField("component", "backup")

		desiredKind, err := ParseDesiredV0(desired)
		if err != nil {
			return nil, nil, nil, nil, nil, false, errors.Wrap(err, "parsing desired state failed")
		}
		desired.Parsed = desiredKind

		secrets, existing := getSecretsMap(desiredKind)

		if !monitor.IsVerbose() && desiredKind.Spec.Verbose {
			internalMonitor.Verbose()
		}

		destroyS, err := secret.AdaptFuncToDestroy(namespace, secretName)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}

		_, destroyB, err := backup.AdaptFunc(
			internalMonitor,
			name,
			namespace,
			componentLabels,
			[]string{},
			checkDBReady,
			desiredKind.Spec.Bucket,
			desiredKind.Spec.Cron,
			secretName,
			secretKey,
			timestamp,
			nodeselector,
			tolerations,
			features,
			version,
		)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}

		_, destroyR, err := restore.AdaptFunc(
			monitor,
			name,
			namespace,
			componentLabels,
			[]string{},
			desiredKind.Spec.Bucket,
			timestamp,
			nodeselector,
			tolerations,
			checkDBReady,
			secretName,
			secretKey,
			version,
		)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}

		_, destroyC, err := clean.AdaptFunc(
			monitor,
			name,
			namespace,
			componentLabels,
			[]string{},
			[]string{},
			nodeselector,
			tolerations,
			checkDBReady,
			secretName,
			secretKey,
			version,
		)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}

		destroyers := make([]operator.DestroyFunc, 0)
		for _, feature := range features {
			switch feature {
			case backup.Normal, backup.Instant:
				destroyers = append(destroyers,
					operator.ResourceDestroyToZitadelDestroy(destroyS),
					destroyB,
				)
			case clean.Instant:
				destroyers = append(destroyers,
					destroyC,
				)
			case restore.Instant:
				destroyers = append(destroyers,
					destroyR,
				)
			}
		}

		return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {

				if err := desiredKind.validateSecrets(); err != nil {
					return nil, err
				}

				currentDB, err := coreDB.ParseQueriedForDatabase(queried)
				if err != nil {
					return nil, err
				}

				databases, err := currentDB.GetListDatabasesFunc()(k8sClient)
				if err != nil {
					databases = []string{}
				}

				users, err := currentDB.GetListUsersFunc()(k8sClient)
				if err != nil {
					users = []string{}
				}

				value, err := read.GetSecretValue(k8sClient, desiredKind.Spec.ServiceAccountJSON, desiredKind.Spec.ExistingServiceAccountJSON)
				if err != nil {
					return nil, err
				}

				queryS, err := secret.AdaptFuncToEnsure(namespace, labels.MustForName(componentLabels, secretName), map[string]string{secretKey: value})
				if err != nil {
					return nil, err
				}

				queryB, _, err := backup.AdaptFunc(
					internalMonitor,
					name,
					namespace,
					componentLabels,
					databases,
					checkDBReady,
					desiredKind.Spec.Bucket,
					desiredKind.Spec.Cron,
					secretName,
					secretKey,
					timestamp,
					nodeselector,
					tolerations,
					features,
					version,
				)
				if err != nil {
					return nil, err
				}

				queryR, _, err := restore.AdaptFunc(
					monitor,
					name,
					namespace,
					componentLabels,
					databases,
					desiredKind.Spec.Bucket,
					timestamp,
					nodeselector,
					tolerations,
					checkDBReady,
					secretName,
					secretKey,
					version,
				)
				if err != nil {
					return nil, err
				}

				queryC, _, err := clean.AdaptFunc(
					monitor,
					name,
					namespace,
					componentLabels,
					databases,
					users,
					nodeselector,
					tolerations,
					checkDBReady,
					secretName,
					secretKey,
					version,
				)
				if err != nil {
					return nil, err
				}

				queriers := make([]operator.QueryFunc, 0)
				cleanupQueries := make([]operator.QueryFunc, 0)
				if databases != nil && len(databases) != 0 {
					for _, feature := range features {
						switch feature {
						case backup.Normal:
							queriers = append(queriers,
								operator.ResourceQueryToZitadelQuery(queryS),
								queryB,
							)
						case backup.Instant:
							queriers = append(queriers,
								operator.ResourceQueryToZitadelQuery(queryS),
								queryB,
							)
							cleanupQueries = append(cleanupQueries,
								operator.EnsureFuncToQueryFunc(backup.GetCleanupFunc(monitor, namespace, name)),
							)
						case clean.Instant:
							queriers = append(queriers,
								operator.ResourceQueryToZitadelQuery(queryS),
								queryC,
							)
							cleanupQueries = append(cleanupQueries,
								operator.EnsureFuncToQueryFunc(clean.GetCleanupFunc(monitor, namespace, name)),
							)
						case restore.Instant:
							queriers = append(queriers,
								operator.ResourceQueryToZitadelQuery(queryS),
								queryR,
							)
							cleanupQueries = append(cleanupQueries,
								operator.EnsureFuncToQueryFunc(restore.GetCleanupFunc(monitor, namespace, name)),
							)
						}
					}
				}

				for _, cleanup := range cleanupQueries {
					queriers = append(queriers, cleanup)
				}

				return operator.QueriersToEnsureFunc(internalMonitor, false, queriers, k8sClient, queried)
			},
			operator.DestroyersToDestroyFunc(internalMonitor, destroyers),
			func(kubernetes.ClientInt, map[string]interface{}, bool) error { return nil },
			secrets,
			existing,
			false,
			nil
	}
}
