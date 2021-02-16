package database

import (
	"fmt"
	"github.com/caos/zitadel/operator"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
)

func AdaptFunc(
	monitor mntr.Monitor,
	namespace string,
	deployName string,
	containerName string,
	certsDir string,
	userName string,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	cmdSql := fmt.Sprintf("cockroach sql --certs-dir=%s", certsDir)

	createSql := fmt.Sprintf("CREATE DATABASE IF NOT EXISTS %s ", userName)

	deleteSql := fmt.Sprintf("DROP DATABASE IF EXISTS %s", userName)

	ensureDatabase := func(k8sClient kubernetes.ClientInt) error {
		return k8sClient.ExecInPodOfDeployment(namespace, deployName, containerName, fmt.Sprintf("%s -e '%s;'", cmdSql, createSql))
	}

	destroyDatabase := func(k8sClient kubernetes.ClientInt) error {
		return k8sClient.ExecInPodOfDeployment(namespace, deployName, containerName, fmt.Sprintf("%s -e '%s;'", cmdSql, deleteSql))
	}

	queriers := []operator.QueryFunc{
		operator.EnsureFuncToQueryFunc(ensureDatabase),
	}

	destroyers := []operator.DestroyFunc{
		destroyDatabase,
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(monitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(monitor, destroyers),
		nil
}
