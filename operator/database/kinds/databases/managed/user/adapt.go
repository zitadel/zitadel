package user

import (
	"fmt"
	"github.com/caos/zitadel/operator"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
	"github.com/caos/zitadel/operator/database/kinds/databases/managed/certificate"
)

func AdaptFunc(
	monitor mntr.Monitor,
	namespace string,
	deployName string,
	containerName string,
	certsDir string,
	userName string,
	password string,
	componentLabels *labels.Component,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	cmdSql := fmt.Sprintf("cockroach sql --certs-dir=%s", certsDir)

	createSql := fmt.Sprintf("CREATE USER IF NOT EXISTS %s ", userName)
	if password != "" {
		createSql = fmt.Sprintf("%s WITH PASSWORD %s", createSql, password)
	}

	deleteSql := fmt.Sprintf("DROP USER IF EXISTS %s", userName)

	_, _, addUserFunc, deleteUserFunc, _, err := certificate.AdaptFunc(monitor, namespace, componentLabels, "", false)
	if err != nil {
		return nil, nil, err
	}

	addUser, err := addUserFunc(userName)
	if err != nil {
		return nil, nil, err
	}
	ensureUser := func(k8sClient kubernetes.ClientInt) error {
		return k8sClient.ExecInPodOfDeployment(namespace, deployName, containerName, fmt.Sprintf("%s -e '%s;'", cmdSql, createSql))
	}

	deleteUser, err := deleteUserFunc(userName)
	if err != nil {
		return nil, nil, err
	}
	destoryUser := func(k8sClient kubernetes.ClientInt) error {
		return k8sClient.ExecInPodOfDeployment(namespace, deployName, containerName, fmt.Sprintf("%s -e '%s;'", cmdSql, deleteSql))
	}

	queriers := []operator.QueryFunc{
		addUser,
		operator.EnsureFuncToQueryFunc(ensureUser),
	}

	destroyers := []operator.DestroyFunc{
		destoryUser,
		deleteUser,
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(monitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(monitor, destroyers),
		nil
}
