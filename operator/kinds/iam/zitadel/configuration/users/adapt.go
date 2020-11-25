package users

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/kinds/iam/zitadel/database"
)

func AdaptFunc(
	monitor mntr.Monitor,
	users map[string]string,
	dbClient database.ClientInt,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("component", "db-users")
	destroyers := make([]operator.DestroyFunc, 0)

	destroyers = append(destroyers, func(k8sClient kubernetes.ClientInt) error {
		list, err := dbClient.ListUsers(internalMonitor, k8sClient)
		if err != nil {
			return err
		}
		for _, listedUser := range list {
			if err := dbClient.DeleteUser(internalMonitor, listedUser, k8sClient); err != nil {
				return err
			}
		}
		return nil
	})

	usernames := []string{}
	for username := range users {
		usernames = append(usernames, username)
	}

	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
			queriers := make([]operator.QueryFunc, 0)
			db, err := database.GetDatabaseInQueried(queried)
			if err != nil {
				return nil, err
			}

			for _, username := range usernames {
				ensure := createIfNecessary(monitor, username, db.Users, dbClient)
				if ensure != nil {
					queriers = append(queriers, operator.EnsureFuncToQueryFunc(ensure))
				}
			}
			for _, listedUser := range db.Users {
				ensure := deleteIfNotRequired(monitor, listedUser, usernames, dbClient)
				if ensure != nil {
					queriers = append(queriers, operator.EnsureFuncToQueryFunc(ensure))
				}
			}

			if queriers == nil || len(queriers) == 0 {
				return func(k8sClient kubernetes.ClientInt) error { return nil }, nil
			}
			return operator.QueriersToEnsureFunc(internalMonitor, false, queriers, k8sClient, queried)
		}, operator.DestroyersToDestroyFunc(internalMonitor, destroyers),
		nil
}
