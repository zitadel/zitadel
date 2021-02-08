package users

import (
	"github.com/caos/orbos/mntr"
	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/zitadel/operator"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/database"
)

func createIfNecessary(monitor mntr.Monitor, user string, list []string, dbClient database.ClientInt) operator.EnsureFunc {
	existing := false
	for _, listedUser := range list {
		if listedUser == user {
			existing = true
		}
	}
	if !existing {
		return func(k8sClient kubernetes.ClientInt) error {
			return dbClient.AddUser(monitor, user, k8sClient)
		}
	}

	return nil
}

func deleteIfNotRequired(monitor mntr.Monitor, listedUser string, list []string, dbClient database.ClientInt) operator.EnsureFunc {
	required := false
	for _, user := range list {
		if user == listedUser {
			required = true
		}
	}
	if !required {
		return func(k8sClient kubernetes.ClientInt) error {
			return dbClient.DeleteUser(monitor, listedUser, k8sClient)
		}
	}

	return nil
}
