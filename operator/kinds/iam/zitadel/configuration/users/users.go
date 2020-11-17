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
	repoURL string,
	repoKey string,
) (
	operator.QueryFunc,
	operator.DestroyFunc,
	error,
) {
	internalMonitor := monitor.WithField("component", "db-users")
	destroyers := make([]operator.DestroyFunc, 0)

	destroyers = append(destroyers, func(k8sClient kubernetes.ClientInt) error {
		list, err := database.ListUsers(internalMonitor, k8sClient, repoURL, repoKey)
		if err != nil {
			return err
		}
		for _, listedUser := range list {
			if err := database.DeleteUser(internalMonitor, listedUser, k8sClient, repoURL, repoKey); err != nil {
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
			list, err := database.ListUsers(internalMonitor, k8sClient, repoURL, repoKey)
			if err != nil {
				return nil, err
			}

			for _, username := range usernames {
				queriers = append(queriers, createIfNecessary(monitor, username, list, repoURL, repoKey))
			}
			for _, listedUser := range list {
				queriers = append(queriers, deleteIfNotRequired(monitor, listedUser, usernames, repoURL, repoKey))
			}

			return operator.QueriersToEnsureFunc(internalMonitor, false, queriers, k8sClient, queried)
		}, operator.DestroyersToDestroyFunc(internalMonitor, destroyers),
		nil
}

func createIfNecessary(monitor mntr.Monitor, user string, list []string, repoURL, repoKey string) operator.QueryFunc {
	addUser := func(k8sClient kubernetes.ClientInt) error {
		existing := false
		for _, listedUser := range list {
			if listedUser == user {
				existing = true
			}
		}
		if !existing {
			return database.AddUser(monitor, user, k8sClient, repoURL, repoKey)
		}
		return nil
	}
	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
		return addUser, nil
	}
}

func deleteIfNotRequired(monitor mntr.Monitor, listedUser string, list []string, repoURL, repoKey string) operator.QueryFunc {
	deleteUser := func(k8sClient kubernetes.ClientInt) error {
		required := false
		for _, user := range list {
			if user == listedUser {
				required = true
			}
		}
		if !required {
			return database.DeleteUser(monitor, listedUser, k8sClient, repoURL, repoKey)
		}
		return nil
	}
	return func(k8sClient kubernetes.ClientInt, queried map[string]interface{}) (operator.EnsureFunc, error) {
		return deleteUser, nil
	}
}
