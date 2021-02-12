package bucket

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/helper"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/secret"
	"github.com/caos/orbos/pkg/labels"
	secretpkg "github.com/caos/orbos/pkg/secret"
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
	return func(monitor mntr.Monitor, desired *tree.Tree, current *tree.Tree) (queryFunc operator.QueryFunc, destroyFunc operator.DestroyFunc, secrets map[string]*secretpkg.Secret, err error) {

		internalMonitor := monitor.WithField("component", "backup")

		desiredKind, err := ParseDesiredV0(desired)
		if err != nil {
			return nil, nil, nil, errors.Wrap(err, "parsing desired state failed")
		}
		desired.Parsed = desiredKind

		if !monitor.IsVerbose() && desiredKind.Spec.Verbose {
			internalMonitor.Verbose()
		}

		destroyS, err := secret.AdaptFuncToDestroy(namespace, secretName)
		if err != nil {
			return nil, nil, nil, err
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
			return nil, nil, nil, err
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
			return nil, nil, nil, err
		}

		_, destroyC, err := clean.AdaptFunc(
			monitor,
			name,
			namespace,
			componentLabels,
			[]string{},
			nodeselector,
			tolerations,
			checkDBReady,
			secretName,
			secretKey,
			version,
		)
		if err != nil {
			return nil, nil, nil, err
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
				currentDB, err := coreDB.ParseQueriedForDatabase(queried)
				if err != nil {
					return nil, err
				}

				databases, err := currentDB.GetListDatabasesFunc()(k8sClient)
				if err != nil {
					databases = []string{}
				}

				value, err := helper.GetSecretValue(k8sClient, desiredKind.Spec.ServiceAccountJSON, desiredKind.Spec.ExistingServiceAccountJSON)
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
				if databases != nil && len(databases) != 0 {
					for _, feature := range features {
						switch feature {
						case backup.Normal, backup.Instant:
							queriers = append(queriers,
								operator.ResourceQueryToZitadelQuery(queryS),
								queryB,
							)
						case clean.Instant:
							queriers = append(queriers,
								queryC,
							)
						case restore.Instant:
							queriers = append(queriers,
								queryR,
							)
						}
					}
				}

				return operator.QueriersToEnsureFunc(internalMonitor, false, queriers, k8sClient, queried)
			},
			operator.DestroyersToDestroyFunc(internalMonitor, destroyers),
			getSecretsMap(desiredKind),
			nil
	}
}
