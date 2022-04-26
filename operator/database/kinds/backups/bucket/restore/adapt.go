package restore

import (
	"time"

	"github.com/zitadel/zitadel/operator/database/kinds/backups/bucket/command"
	"github.com/zitadel/zitadel/pkg/databases/db"

	"github.com/zitadel/zitadel/operator"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/kubernetes/resources/job"
	"github.com/caos/orbos/pkg/labels"
	corev1 "k8s.io/api/core/v1"
)

const (
	Instant            = "restore"
	defaultMode        = int32(256)
	certPath           = "/cockroach/cockroach-certs"
	secretPath         = "/secrets/sa.json"
	jobPrefix          = "backup-"
	jobSuffix          = "-restore"
	internalSecretName = "client-certs"
	rootSecretName     = "cockroachdb.node"
	timeout            = 45 * time.Minute
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
	secretName string,
	secretKey string,
	dbConn db.Connection,
	image string,
) (
	queryFunc operator.QueryFunc,
	destroyFunc operator.DestroyFunc,
	err error,
) {

	jobName := jobPrefix + backupName + jobSuffix

	restoreCmd, env := command.GetSQLCommand(
		dbConn,
		certPath,
		command.GetBackupRestoreStatement(
			bucketName,
			backupName,
			timestamp,
			secretPath,
			command.Restore,
		),
	)

	jobdef := getJob(
		namespace,
		labels.MustForName(componentLabels, GetJobName(backupName)),
		nodeselector,
		tolerations,
		secretName,
		secretKey,
		restoreCmd,
		image,
		env,
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
