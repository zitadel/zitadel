package user

import (
	"fmt"
	"github.com/caos/orbos/pkg/kubernetes/resources/secret"
	"github.com/caos/zitadel/operator"

	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/labels"
)

func AdaptFunc(
	monitor mntr.Monitor,
	namespace,
	podName,
	containerName,
	certsDir,
	userName,
	password,
	certsSecretName,
	userCrtFilename,
	userKeyFilename string,
	pwSecretSelectable *labels.Selectable,
	pwSecretKey string,
	componentLabels *labels.Component,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	cmdSql := fmt.Sprintf("cockroach sql --certs-dir=%s", certsDir)
	createSql := fmt.Sprintf(`%s --execute "CREATE USER IF NOT EXISTS %s;" --execute "GRANT admin TO %s;"`, cmdSql, userName, userName)

	if password != "" {
		createSql += fmt.Sprintf(` --execute "ALTER USER %s WITH PASSWORD '%s';"`, userName, password)
	}

	deleteSql := fmt.Sprintf(`%s --execute "DROP USER IF EXISTS %s;"`, cmdSql, userName)

	ensureUser := func(k8sClient kubernetes.ClientInt) error {
		return k8sClient.ExecInPod(namespace, podName, containerName, createSql)
	}
	destoryUser := func(k8sClient kubernetes.ClientInt) error {
		return k8sClient.ExecInPod(namespace, podName, containerName, deleteSql)
	}

	queryPWSecret, err := secret.AdaptFuncToEnsure(namespace, pwSecretSelectable, map[string]string{pwSecretKey: password})
	if err != nil {
		return nil, nil, err
	}
	destroyPWSecret, err := secret.AdaptFuncToDestroy(namespace, pwSecretSelectable.Name())
	if err != nil {
		return nil, nil, err
	}

	queriers := []operator.QueryFunc{
		operator.EnsureFuncToQueryFunc(ensureUser),
		operator.ResourceQueryToZitadelQuery(queryPWSecret),
	}

	destroyers := []operator.DestroyFunc{
		destoryUser,
		operator.ResourceDestroyToZitadelDestroy(destroyPWSecret),
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			return operator.QueriersToEnsureFunc(monitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(monitor, destroyers),
		nil
}
