package backup

import (
	"github.com/caos/zitadel/operator/database/kinds/backups/core"
	"path/filepath"
	"time"

	"github.com/caos/zitadel/operator"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/cronjob"
	"github.com/caos/orbos/pkg/kubernetes/resources/job"
	"github.com/caos/orbos/pkg/labels"
	corev1 "k8s.io/api/core/v1"
)

const (
	defaultMode             int32 = 256
	certPath                      = "/cockroach/cockroach-certs"
	secretsPath                   = "/secrets"
	internalSecretName            = "secrets"
	backupNameEnv                 = "BACKUP_NAME"
	cronJobNamePrefix             = "backup-"
	internalCertsSecretName       = "client-certs"
	rootSecretName                = "cockroachdb.client.root"
	timeout                       = 15 * time.Minute
	Normal                        = "backup"
	Instant                       = "instantbackup"
)

func AdaptFunc(
	monitor mntr.Monitor,
	backupName string,
	namespace string,
	componentLabels *labels.Component,
	checkDBReady operator.EnsureFunc,
	bucketName string,
	cron string,
	accessKeyIDKey string,
	secretAccessKeyKey string,
	sessionTokenKey string,
	region string,
	endpoint string,
	timestamp string,
	nodeselector map[string]string,
	tolerations []corev1.Toleration,
	dbURL string,
	dbPort int32,
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
		certPath,
		filepath.Join(secretsPath, accessKeyIDKey),
		filepath.Join(secretsPath, secretAccessKeyKey),
		filepath.Join(secretsPath, sessionTokenKey),
		region,
		endpoint,
		dbURL,
		dbPort,
	)

	jobSpecDef := getJobSpecDef(
		nodeselector,
		tolerations,
		core.GetSecretName(backupName),
		backupName,
		image,
		command,
	)

	destroyers := []operator.DestroyFunc{}
	queriers := []operator.QueryFunc{}

	cronJobDef := getCronJob(
		namespace,
		labels.MustForName(componentLabels, core.GetBackupJobName(backupName)),
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
				operator.EnsureFuncToQueryFunc(checkDBReady),
				operator.ResourceQueryToZitadelQuery(queryCJ),
			)
		case Instant:
			destroyers = append(destroyers,
				operator.ResourceDestroyToZitadelDestroy(destroyJ),
			)
			queriers = append(queriers,
				operator.EnsureFuncToQueryFunc(checkDBReady),
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
