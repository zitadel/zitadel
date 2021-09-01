package restore

import (
	"time"

	"github.com/caos/zitadel/operator"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/job"
	"github.com/caos/orbos/pkg/labels"
	corev1 "k8s.io/api/core/v1"
)

const (
	Instant                 = "restore"
	defaultMode             = int32(256)
	certPath                = "/cockroach/cockroach-certs"
	saInternalSecretName    = "sa-json"
	saSecretPath            = "/secrets/sa.json"
	akidInternalSecretName  = "akid"
	akidSecretPath          = "/secrets/akid"
	sakInternalSecretName   = "sak"
	sakSecretPath           = "/secrets/sak"
	jobPrefix               = "backup-"
	jobSuffix               = "-restore"
	certsInternalSecretName = "client-certs"
	rootSecretName          = "cockroachdb.client.root"
	timeout                 = 45 * time.Minute
	backupNameEnv           = "BACKUP_NAME"
)

func AdaptFunc(
	monitor mntr.Monitor,
	backupName string,
	namespace string,
	componentLabels *labels.Component,
	bucketName string,
	timestamp string,
	nodeselector map[string]string,
	tolerations []corev1.Toleration,
	checkDBReady operator.EnsureFunc,
	backupSecretName string,
	saSecretKey string,
	assetAKIDKey string,
	assetSAKKey string,
	dbURL string,
	dbPort int32,
	image string,
	assetEndpoint string,
	assetPrefix string,
) (
	queryFunc operator.QueryFunc,
	destroyFunc operator.DestroyFunc,
	err error,
) {

	jobName := jobPrefix + backupName + jobSuffix
	command := getCommand(
		timestamp,
		bucketName,
		backupName,
		certPath,
		saSecretPath,
		dbURL,
		dbPort,
		assetEndpoint,
		assetPrefix,
	)

	jobdef := getJob(
		namespace,
		labels.MustForName(componentLabels, GetJobName(backupName)),
		nodeselector,
		tolerations,
		backupSecretName,
		saSecretKey,
		assetAKIDKey,
		assetSAKKey,
		command,
		image,
	)

	destroyJ, err := job.AdaptFuncToDestroy(jobName, namespace)
	if err != nil {
		return nil, nil, err
	}

	destroyers := []operator.DestroyFunc{
		operator.ResourceDestroyToZitadelDestroy(destroyJ),
	}

	queryJ, err := job.AdaptFuncToEnsure(jobdef)
	if err != nil {
		return nil, nil, err
	}

	queriers := []operator.QueryFunc{
		operator.EnsureFuncToQueryFunc(checkDBReady),
		operator.ResourceQueryToZitadelQuery(queryJ),
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(monitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(monitor, destroyers),

		nil
}

func GetJobName(backupName string) string {
	return jobPrefix + backupName + jobSuffix
}
