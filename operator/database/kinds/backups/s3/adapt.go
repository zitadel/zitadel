package s3

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/secret"
	"github.com/caos/orbos/pkg/labels"
	secretpkg "github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/secret/read"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/database/kinds/backups/s3/backup"
	"github.com/caos/zitadel/operator/database/kinds/backups/s3/restore"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
)

const (
	accessKeyIDName     = "backup-accessaccountkey"
	accessKeyIDKey      = "accessaccountkey"
	secretAccessKeyName = "backup-secretaccesskey"
	secretAccessKeyKey  = "secretaccesskey"
	sessionTokenName    = "backup-sessiontoken"
	sessionTokenKey     = "sessiontoken"
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
	dbURL string,
	dbPort int32,
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

		destroySAKI, err := secret.AdaptFuncToDestroy(namespace, accessKeyIDName)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}

		destroySSAK, err := secret.AdaptFuncToDestroy(namespace, secretAccessKeyName)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}

		destroySSTK, err := secret.AdaptFuncToDestroy(namespace, sessionTokenName)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}

		_, destroyB, err := backup.AdaptFunc(
			internalMonitor,
			name,
			namespace,
			componentLabels,
			checkDBReady,
			desiredKind.Spec.Bucket,
			desiredKind.Spec.Cron,
			accessKeyIDName,
			accessKeyIDKey,
			secretAccessKeyName,
			secretAccessKeyKey,
			sessionTokenName,
			sessionTokenKey,
			desiredKind.Spec.Region,
			desiredKind.Spec.Endpoint,
			timestamp,
			nodeselector,
			tolerations,
			dbURL,
			dbPort,
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
			desiredKind.Spec.Bucket,
			timestamp,
			accessKeyIDName,
			accessKeyIDKey,
			secretAccessKeyName,
			secretAccessKeyKey,
			sessionTokenName,
			sessionTokenKey,
			desiredKind.Spec.Region,
			desiredKind.Spec.Endpoint,
			nodeselector,
			tolerations,
			checkDBReady,
			dbURL,
			dbPort,
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
					operator.ResourceDestroyToZitadelDestroy(destroySSAK),
					operator.ResourceDestroyToZitadelDestroy(destroySAKI),
					operator.ResourceDestroyToZitadelDestroy(destroySSTK),
					destroyB,
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

				valueAKI, err := read.GetSecretValue(k8sClient, desiredKind.Spec.AccessKeyID, desiredKind.Spec.ExistingAccessKeyID)
				if err != nil {
					return nil, err
				}

				querySAKI, err := secret.AdaptFuncToEnsure(namespace, labels.MustForName(componentLabels, accessKeyIDName), map[string]string{accessKeyIDKey: valueAKI})
				if err != nil {
					return nil, err
				}

				valueSAK, err := read.GetSecretValue(k8sClient, desiredKind.Spec.SecretAccessKey, desiredKind.Spec.ExistingSecretAccessKey)
				if err != nil {
					return nil, err
				}

				querySSAK, err := secret.AdaptFuncToEnsure(namespace, labels.MustForName(componentLabels, secretAccessKeyName), map[string]string{secretAccessKeyKey: valueSAK})
				if err != nil {
					return nil, err
				}

				valueST, err := read.GetSecretValue(k8sClient, desiredKind.Spec.SessionToken, desiredKind.Spec.ExistingSessionToken)
				if err != nil {
					return nil, err
				}

				querySST, err := secret.AdaptFuncToEnsure(namespace, labels.MustForName(componentLabels, sessionTokenName), map[string]string{sessionTokenKey: valueST})
				if err != nil {
					return nil, err
				}

				queryB, _, err := backup.AdaptFunc(
					internalMonitor,
					name,
					namespace,
					componentLabels,
					checkDBReady,
					desiredKind.Spec.Bucket,
					desiredKind.Spec.Cron,
					accessKeyIDName,
					accessKeyIDKey,
					secretAccessKeyName,
					secretAccessKeyKey,
					sessionTokenName,
					sessionTokenKey,
					desiredKind.Spec.Region,
					desiredKind.Spec.Endpoint,
					timestamp,
					nodeselector,
					tolerations,
					dbURL,
					dbPort,
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
					desiredKind.Spec.Bucket,
					timestamp,
					accessKeyIDName,
					accessKeyIDKey,
					secretAccessKeyName,
					secretAccessKeyKey,
					sessionTokenName,
					sessionTokenKey,
					desiredKind.Spec.Region,
					desiredKind.Spec.Endpoint,
					nodeselector,
					tolerations,
					checkDBReady,
					dbURL,
					dbPort,
					version,
				)
				if err != nil {
					return nil, err
				}

				queriers := make([]operator.QueryFunc, 0)
				cleanupQueries := make([]operator.QueryFunc, 0)
				for _, feature := range features {
					switch feature {
					case backup.Normal:
						queriers = append(queriers,
							operator.ResourceQueryToZitadelQuery(querySAKI),
							operator.ResourceQueryToZitadelQuery(querySSAK),
							operator.ResourceQueryToZitadelQuery(querySST),
							queryB,
						)
					case backup.Instant:
						queriers = append(queriers,
							operator.ResourceQueryToZitadelQuery(querySAKI),
							operator.ResourceQueryToZitadelQuery(querySSAK),
							operator.ResourceQueryToZitadelQuery(querySST),
							queryB,
						)
						cleanupQueries = append(cleanupQueries,
							operator.EnsureFuncToQueryFunc(backup.GetCleanupFunc(monitor, namespace, name)),
						)
					case restore.Instant:
						queriers = append(queriers,
							operator.ResourceQueryToZitadelQuery(querySAKI),
							operator.ResourceQueryToZitadelQuery(querySSAK),
							operator.ResourceQueryToZitadelQuery(querySST),
							queryR,
						)
						cleanupQueries = append(cleanupQueries,
							operator.EnsureFuncToQueryFunc(restore.GetCleanupFunc(monitor, namespace, name)),
						)
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
