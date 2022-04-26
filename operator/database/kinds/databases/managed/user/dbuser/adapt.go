package dbuser

import (
	"fmt"

	"github.com/caos/orbos/pkg/kubernetes/resources/secret"
	"github.com/zitadel/zitadel/operator"
	macherrs "k8s.io/apimachinery/pkg/api/errors"

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
	userName string,
	pwSecretSelectable *labels.Selectable,
	pwSecretKey string,
	pwSecretValue func(kubernetes.ClientInt) (string, error),
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	cmdSql := fmt.Sprintf("cockroach sql --certs-dir=%s", certsDir)

	createSql := func(pw string) string {
		return fmt.Sprintf(`%s --execute "CREATE USER IF NOT EXISTS %s;" --execute "GRANT admin TO %s;" --execute "ALTER USER %s WITH PASSWORD '%s';"`, cmdSql, userName, userName, userName, pw)
	}

	deleteSql := fmt.Sprintf(`%s --execute "DROP USER IF EXISTS %s;"`, cmdSql, userName)

	ensureUser := func(k8sClient kubernetes.ClientInt) error {
		pwValue, err := pwSecretValue(k8sClient)
		if err != nil {
			return err
		}
		return k8sClient.ExecInPod(namespace, podName, containerName, createSql(pwValue))
	}
	destoryUser := func(k8sClient kubernetes.ClientInt) error {
		return k8sClient.ExecInPod(namespace, podName, containerName, deleteSql)
	}

	destroyPWSecret, err := secret.AdaptFuncToDestroy(namespace, pwSecretSelectable.Name())
	if err != nil {
		return nil, nil, err
	}

	queriers := []operator.QueryFunc{
		operator.EnsureFuncToQueryFunc(ensureUser),
	}

	destroyers := []operator.DestroyFunc{
		destoryUser,
		operator.ResourceDestroyToZitadelDestroy(destroyPWSecret),
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			pwValue, err := pwSecretValue(k8sClient)
			if err != nil {
				return nil, err
			}

			existingSecret, err := k8sClient.GetSecret(namespace, pwSecretSelectable.Name())
			if err != nil && !macherrs.IsNotFound(err) {
				return nil, err
			}
			err = nil
			data := make(map[string]string)
			if existingSecret != nil {
				for k, v := range existingSecret.Data {
					data[k] = string(v)
				}
			}

			if _, ok := data[pwSecretKey]; !ok {
				data[pwSecretKey] = pwValue
			}

			queryPWSecret, err := secret.AdaptFuncToEnsure(namespace, pwSecretSelectable, data)
			if err != nil {
				return nil, err
			}

			queriers = append(queriers, operator.ResourceQueryToZitadelQuery(queryPWSecret))

			return operator.QueriersToEnsureFunc(monitor, false, queriers, k8sClient, queried)
		},
		operator.DestroyersToDestroyFunc(monitor, destroyers),
		nil
}
