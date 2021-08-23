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
	"github.com/caos/zitadel/operator/common"
	"github.com/caos/zitadel/operator/zitadel/kinds/backups/bucket/backup"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

const (
	saSecretName     = "backup-serviceaccountjson"
	saSecretKey      = "serviceaccountjson"
	configSecretName = "backup-serviceaccountjson"
	configSecretKey  = "serviceaccountjson"
)

func AdaptFunc(
	name string,
	namespace string,
	componentLabels *labels.Component,
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
		map[string]*secretpkg.Secret,
		map[string]*secretpkg.Existing,
		bool,
		error,
	) {

		internalMonitor := monitor.WithField("component", "assetbackup")

		desiredKind, err := ParseDesiredV0(desired)
		if err != nil {
			return nil, nil, nil, nil, nil, false, errors.Wrap(err, "parsing desired state failed")
		}
		desired.Parsed = desiredKind

		secrets, existing := getSecretsMap(desiredKind)

		if !monitor.IsVerbose() && desiredKind.Spec.Verbose {
			internalMonitor.Verbose()
		}

		destroySS, err := secret.AdaptFuncToDestroy(namespace, saSecretName)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}

		destroySC, err := secret.AdaptFuncToDestroy(namespace, configSecretName)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}

		image := common.AssetBackupImage.Reference(customImageRegistry)

		_, destroyB, err := backup.AdaptFunc(
			internalMonitor,
			name,
			namespace,
			componentLabels,
			desiredKind.Spec.Bucket,
			desiredKind.Spec.Cron,
			saSecretName,
			saSecretKey,
			configSecretName,
			configSecretKey,
			timestamp,
			nodeselector,
			tolerations,
			features,
			image,
		)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}

		/*
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
				saSecretName,
				saSecretKey,
				configSecretName,
				configSecretKey,
				image,
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
				nodeselector,
				tolerations,
				saSecretName,
				saSecretKey,
				configSecretName,
				configSecretKey,
				image,
			)
			if err != nil {
				return nil, nil, nil, nil, nil, false, err
			}
		*/

		destroyers := make([]operator.DestroyFunc, 0)
		for _, feature := range features {
			switch feature {
			case backup.Normal, backup.Instant:
				destroyers = append(destroyers,
					operator.ResourceDestroyToZitadelDestroy(destroySS),
					operator.ResourceDestroyToZitadelDestroy(destroySC),
					destroyB,
				)
				/*case clean.Instant:
					destroyers = append(destroyers,
						destroyC,
					)
				case restore.Instant:
					destroyers = append(destroyers,
						destroyR,
					)*/
			}
		}

		return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {

				if err := desiredKind.validateSecrets(); err != nil {
					return nil, err
				}

				valueS, err := read.GetSecretValue(k8sClient, desiredKind.Spec.ServiceAccountJSON, desiredKind.Spec.ExistingServiceAccountJSON)
				if err != nil {
					return nil, err
				}

				querySS, err := secret.AdaptFuncToEnsure(namespace, labels.MustForName(componentLabels, saSecretName), map[string]string{saSecretKey: valueS})
				if err != nil {
					return nil, err
				}

				valueC, err := read.GetSecretValue(k8sClient, desiredKind.Spec.ServiceAccountJSON, desiredKind.Spec.ExistingServiceAccountJSON)
				if err != nil {
					return nil, err
				}
				querySC, err := secret.AdaptFuncToEnsure(namespace, labels.MustForName(componentLabels, configSecretName), map[string]string{configSecretKey: valueC})
				if err != nil {
					return nil, err
				}

				queryB, _, err := backup.AdaptFunc(
					internalMonitor,
					name,
					namespace,
					componentLabels,
					desiredKind.Spec.Bucket,
					desiredKind.Spec.Cron,
					saSecretName,
					saSecretKey,
					configSecretName,
					configSecretKey,
					timestamp,
					nodeselector,
					tolerations,
					features,
					image,
				)
				if err != nil {
					return nil, err
				}
				/*
					queryR, _, err := restore.AdaptFunc(
						monitor,
						name,
						namespace,
						componentLabels,
						desiredKind.Spec.Bucket,
						timestamp,
						nodeselector,
						tolerations,
						saSecretName,
						saSecretKey,
						configSecretName,
						configSecretKey,
						image,
					)
					if err != nil {
						return nil, err
					}

					queryC, _, err := clean.AdaptFunc(
						monitor,
						name,
						namespace,
						componentLabels,
						nodeselector,
						tolerations,
						saSecretName,
						saSecretKey,
						configSecretName,
						configSecretKey,
						image,
					)
					if err != nil {
						return nil, err
					}
				*/
				queriers := make([]operator.QueryFunc, 0)
				cleanupQueries := make([]operator.QueryFunc, 0)
				for _, feature := range features {
					switch feature {
					case backup.Normal:
						queriers = append(queriers,
							operator.ResourceQueryToZitadelQuery(querySS),
							operator.ResourceQueryToZitadelQuery(querySC),
							queryB,
						)
					case backup.Instant:
						queriers = append(queriers,
							operator.ResourceQueryToZitadelQuery(querySS),
							operator.ResourceQueryToZitadelQuery(querySC),
							queryB,
						)
						cleanupQueries = append(cleanupQueries,
							operator.EnsureFuncToQueryFunc(backup.GetCleanupFunc(monitor, namespace, name)),
						)
						/*case clean.Instant:
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
							)*/
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
