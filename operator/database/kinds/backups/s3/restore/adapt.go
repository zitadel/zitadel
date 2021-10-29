package restore

import (
	"github.com/caos/zitadel/operator/database/kinds/backups/core"
	"path/filepath"
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
	secretsPath             = "/secrets"
	internalSecretsName     = "secrets"
	internalCertsSecretName = "client-certs"
	rootSecretName          = "cockroachdb.client.root"
	timeout                 = 15 * time.Minute
)

func AdaptFunc(
	monitor mntr.Monitor,
	backupName string,
	namespace string,
	componentLabels *labels.Component,
	bucketName string,
	timestamp string,
	accessKeyIDKey string,
	secretAccessKeyKey string,
	sessionTokenKey string,
	region string,
	endpoint string,
	nodeselector map[string]string,
	tolerations []corev1.Toleration,
	checkDBReady operator.EnsureFunc,
	dbURL string,
	dbPort int32,
	image string,
) (
	queryFunc operator.QueryFunc,
	destroyFunc operator.DestroyFunc,
	err error,
) {

	jobName := core.GetRestoreJobName(backupName)
	command := getCommand(
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

	jobdef := getJob(
		namespace,
		labels.MustForName(componentLabels, jobName),
		nodeselector,
		tolerations,
		core.GetSecretName(backupName),
		image,
		command,
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
