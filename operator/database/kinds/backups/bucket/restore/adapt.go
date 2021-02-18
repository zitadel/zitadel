package restore

import (
	"github.com/caos/zitadel/operator"
	"time"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/job"
	"github.com/caos/orbos/pkg/labels"
	corev1 "k8s.io/api/core/v1"
)

const (
	Instant                          = "restore"
	defaultMode                      = int32(256)
	certPath                         = "/cockroach/cockroach-certs"
	secretPath                       = "/secrets/sa.json"
	jobPrefix                        = "backup-"
	jobSuffix                        = "-restore"
	image                            = "ghcr.io/caos/zitadel-crbackup"
	internalSecretName               = "client-certs"
	rootSecretName                   = "cockroachdb.client.root"
	timeout            time.Duration = 1200
)

func AdaptFunc(
	monitor mntr.Monitor,
	backupName string,
	namespace string,
	componentLabels *labels.Component,
	databases []string,
	bucketName string,
	timestamp string,
	nodeselector map[string]string,
	tolerations []corev1.Toleration,
	checkDBReady operator.EnsureFunc,
	secretName string,
	secretKey string,
	version string,
) (
	queryFunc operator.QueryFunc,
	destroyFunc operator.DestroyFunc,
	err error,
) {

	jobName := jobPrefix + backupName + jobSuffix
	command := getCommand(
		timestamp,
		databases,
		bucketName,
		backupName,
	)

	jobdef := getJob(
		namespace,
		labels.MustForName(componentLabels, GetJobName(backupName)),
		nodeselector,
		tolerations,
		secretName,
		secretKey,
		version,
		command)

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
