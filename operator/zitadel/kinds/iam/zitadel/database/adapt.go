package database

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator"
)

func AdaptFunc(
	monitor mntr.Monitor,
	dbClient Client,
) (
	operator.QueryFunc,
	error,
) {

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {

		dbHost, dbPort, err := dbClient.GetConnectionInfo(monitor, k8sClient)
		if err != nil {
			return nil, err
		}

		users, err := dbClient.ListUsers(monitor, k8sClient)
		if err != nil {
			return nil, err
		}

		curr := &Current{
			Host:  dbHost,
			Port:  dbPort,
			Users: users,
		}

		SetDatabaseInQueried(queried, curr)

		return func(k8sClient kubernetes.ClientInt) error {
			return nil
		}, nil
	}, nil
}
