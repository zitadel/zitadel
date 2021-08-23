package backup

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/cronjob"
	"github.com/caos/orbos/pkg/kubernetes/resources/job"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator"
	corev1 "k8s.io/api/core/v1"
	"time"
)

const (
	defaultMode              int32 = 256
	saInternalSecretName           = "sa-json"
	saSecretPath                   = "/secrets/sa.json"
	configInternalSecretName       = "rconfig"
	configSecretPath               = "/secrets/rconfig"
	cronJobNamePrefix              = "backup-"
	backupNameEnv                  = "BACKUP_NAME"
	timeout                        = 15 * time.Minute
	sourceName                     = "minio"
	destinationName                = "bucket"
	Normal                         = "assetbackup"
	Instant                        = "assetinstantbackup"
)

func AdaptFunc(
	monitor mntr.Monitor,
	backupName string,
	namespace string,
	componentLabels *labels.Component,
	bucketName string,
	cron string,
	saSecretName string,
	saSecretKey string,
	configSecretName string,
	configSecretKey string,
	timestamp string,
	nodeselector map[string]string,
	tolerations []corev1.Toleration,
	features []string,
	image string,
) (
	queryFunc operator.QueryFunc,
	destroyFunc operator.DestroyFunc,
	err error,
) {

	command := getBackupCommand(
		timestamp,
		bucketName,
		backupName,
	)

	jobSpecDef := getJobSpecDef(
		nodeselector,
		tolerations,
		saSecretName,
		saSecretKey,
		configSecretName,
		configSecretKey,
		backupName,
		command,
		image,
	)

	destroyers := []operator.DestroyFunc{}
	queriers := []operator.QueryFunc{}

	cronJobDef := getCronJob(
		namespace,
		labels.MustForName(componentLabels, GetJobName(backupName)),
		cron,
		jobSpecDef,
	)

	destroyCJ, err := cronjob.AdaptFuncToDestroy(cronJobDef.Namespace, cronJobDef.Name)
	if err != nil {
		return nil, nil, err
	}

	queryCJ, err := cronjob.AdaptFuncToEnsure(cronJobDef)
	if err != nil {
		return nil, nil, err
	}

	jobDef := getJob(
		namespace,
		labels.MustForName(componentLabels, cronJobNamePrefix+backupName),
		jobSpecDef,
	)

	destroyJ, err := job.AdaptFuncToDestroy(jobDef.Namespace, jobDef.Name)
	if err != nil {
		return nil, nil, err
	}

	queryJ, err := job.AdaptFuncToEnsure(jobDef)
	if err != nil {
		return nil, nil, err
	}

	for _, feature := range features {
		switch feature {
		case Normal:
			destroyers = append(destroyers,
				operator.ResourceDestroyToZitadelDestroy(destroyCJ),
			)
			queriers = append(queriers,
				operator.ResourceQueryToZitadelQuery(queryCJ),
			)
		case Instant:
			destroyers = append(destroyers,
				operator.ResourceDestroyToZitadelDestroy(destroyJ),
			)
			queriers = append(queriers,
				operator.ResourceQueryToZitadelQuery(queryJ),
			)
		}
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(monitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(monitor, destroyers),
		nil
}

func GetJobName(backupName string) string {
	return cronJobNamePrefix + backupName
}
