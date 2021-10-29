package s3

import (
	"fmt"
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/secret"
	"github.com/caos/orbos/pkg/labels"
	secretpkg "github.com/caos/orbos/pkg/secret"
	"github.com/caos/orbos/pkg/secret/read"
	"github.com/caos/orbos/pkg/tree"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/common"
	"github.com/caos/zitadel/operator/database/kinds/backups/core"
	"github.com/caos/zitadel/operator/database/kinds/backups/s3/backup"
	"github.com/caos/zitadel/operator/database/kinds/backups/s3/restore"
	corev1 "k8s.io/api/core/v1"
)

const (
	accessKeyIDKey     = "accessaccountkey"
	secretAccessKeyKey = "secretaccesskey"
	sessionTokenKey    = "sessiontoken"
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

		internalMonitor := monitor.WithField("component", "backup")

		desiredKind, err := ParseDesiredV0(desired)
		if err != nil {
			return nil, nil, nil, nil, nil, false, fmt.Errorf("parsing desired state failed: %s", err)
		}
		desired.Parsed = desiredKind

		secrets, existing := getSecretsMap(desiredKind)

		if !monitor.IsVerbose() && desiredKind.Spec.Verbose {
			internalMonitor.Verbose()
		}

		destroySecret, err := secret.AdaptFuncToDestroy(namespace, core.GetSecretName(name))
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}

		image := common.BackupImage.Reference(customImageRegistry, version)

		_, destroyB, err := backup.AdaptFunc(
			internalMonitor,
			name,
			namespace,
			componentLabels,
			checkDBReady,
			desiredKind.Spec.Bucket,
			desiredKind.Spec.Cron,
			accessKeyIDKey,
			secretAccessKeyKey,
			sessionTokenKey,
			desiredKind.Spec.Region,
			desiredKind.Spec.Endpoint,
			timestamp,
			nodeselector,
			tolerations,
			dbURL,
			dbPort,
			features,
			image,
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
			accessKeyIDKey,
			secretAccessKeyKey,
			sessionTokenKey,
			desiredKind.Spec.Region,
			desiredKind.Spec.Endpoint,
			nodeselector,
			tolerations,
			checkDBReady,
			dbURL,
			dbPort,
			image,
		)
		if err != nil {
			return nil, nil, nil, nil, nil, false, err
		}

		destroyers := make([]operator.DestroyFunc, 0)
		for _, feature := range features {
			switch feature {
			case backup.Normal, backup.Instant:
				destroyers = append(destroyers,
					operator.ResourceDestroyToZitadelDestroy(destroySecret),
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

				secrets := map[string]string{}
				valueAKI, err := read.GetSecretValue(k8sClient, desiredKind.Spec.AccessKeyID, desiredKind.Spec.ExistingAccessKeyID)
				if err != nil {
					return nil, err
				}
				if valueAKI != "" {
					secrets[accessKeyIDKey] = valueAKI
				}

				valueSAK, err := read.GetSecretValue(k8sClient, desiredKind.Spec.SecretAccessKey, desiredKind.Spec.ExistingSecretAccessKey)
				if err != nil {
					return nil, err
				}
				if valueAKI != "" {
					secrets[secretAccessKeyKey] = valueSAK
				}

				valueST, err := read.GetSecretValue(k8sClient, desiredKind.Spec.SessionToken, desiredKind.Spec.ExistingSessionToken)
				if err != nil {
					return nil, err
				}
				if valueAKI != "" {
					secrets[sessionTokenKey] = valueST
				}

				querySecret, err := secret.AdaptFuncToEnsure(namespace,
					labels.MustForName(componentLabels, core.GetSecretName(name)),
					secrets,
				)
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
					accessKeyIDKey,
					secretAccessKeyKey,
					sessionTokenKey,
					desiredKind.Spec.Region,
					desiredKind.Spec.Endpoint,
					timestamp,
					nodeselector,
					tolerations,
					dbURL,
					dbPort,
					features,
					image,
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
					accessKeyIDKey,
					secretAccessKeyKey,
					sessionTokenKey,
					desiredKind.Spec.Region,
					desiredKind.Spec.Endpoint,
					nodeselector,
					tolerations,
					checkDBReady,
					dbURL,
					dbPort,
					image,
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
							operator.ResourceQueryToZitadelQuery(querySecret),
							queryB,
						)
					case backup.Instant:
						queriers = append(queriers,
							operator.ResourceQueryToZitadelQuery(querySecret),
							queryB,
						)
						cleanupQueries = append(cleanupQueries,
							operator.EnsureFuncToQueryFunc(backup.GetCleanupFunc(monitor, namespace, name)),
						)
					case restore.Instant:
						queriers = append(queriers,
							operator.ResourceQueryToZitadelQuery(querySecret),
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
