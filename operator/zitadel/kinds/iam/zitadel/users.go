package zitadel

import (
	"sort"

	"github.com/caos/orbos/pkg/secret/read"

	"github.com/caos/orbos/pkg/kubernetes"
	"github.com/caos/orbos/pkg/secret"

	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/configuration"
)

const (
	migrationUser = "flyway"
	mgmtUser      = "management"
	adminUser     = "adminapi"
	authUser      = "auth"
	authzUser     = "authz"
	notUser       = "notification"
	esUser        = "eventstore"
	queriesUser   = "queries"
)

func getUserListWithoutPasswords(desired *DesiredV0) []string {
	userpw, _ := getAllUsers(nil, desired)
	users := make([]string, 0)
	for user := range userpw {
		users = append(users, user)
	}

	sort.Slice(users, func(i, j int) bool {
		return users[i] < users[j]
	})
	return users
}

func getAllUsers(k8sClient kubernetes.ClientInt, desired *DesiredV0) (map[string]string, error) {
	passwords := &configuration.Passwords{}
	if desired != nil && desired.Spec != nil && desired.Spec.Configuration != nil && desired.Spec.Configuration.Passwords != nil {
		passwords = desired.Spec.Configuration.Passwords
	}
	users := make(map[string]string, 0)

	if err := fillInUserPassword(k8sClient, migrationUser, passwords.Migration, passwords.ExistingMigration, users); err != nil {
		return nil, err
	}
	if err := fillInUserPassword(k8sClient, mgmtUser, passwords.Management, passwords.ExistingManagement, users); err != nil {
		return nil, err
	}
	if err := fillInUserPassword(k8sClient, adminUser, passwords.Adminapi, passwords.ExistingAdminapi, users); err != nil {
		return nil, err
	}
	if err := fillInUserPassword(k8sClient, authUser, passwords.Auth, passwords.ExistingAuth, users); err != nil {
		return nil, err
	}
	if err := fillInUserPassword(k8sClient, authzUser, passwords.Authz, passwords.ExistingAuthz, users); err != nil {
		return nil, err
	}
	if err := fillInUserPassword(k8sClient, notUser, passwords.Notification, passwords.ExistingNotification, users); err != nil {
		return nil, err
	}
	if err := fillInUserPassword(k8sClient, esUser, passwords.Eventstore, passwords.ExistingEventstore, users); err != nil {
		return nil, err
	}
	if err := fillInUserPassword(k8sClient, queriesUser, passwords.Queries, passwords.ExistingQueries, users); err != nil {
		return nil, err
	}

	return users, nil
}

func fillInUserPassword(
	k8sClient kubernetes.ClientInt,
	user string,
	secret *secret.Secret,
	existing *secret.Existing,
	userpw map[string]string,
) error {
	if k8sClient == nil {
		userpw[user] = user
		return nil
	}

	pw, err := read.GetSecretValue(k8sClient, secret, existing)
	if err != nil {
		return err
	}
	if pw != "" {
		userpw[user] = pw
	} else {
		userpw[user] = user
	}

	return nil
}

func getZitadelUserList(k8sClient kubernetes.ClientInt, desired *DesiredV0) (map[string]string, error) {
	allUsersMap, err := getAllUsers(k8sClient, desired)
	if err != nil {
		return nil, err
	}

	allZitadelUsers := make(map[string]string, 0)
	for k, v := range allUsersMap {
		if k != migrationUser {
			allZitadelUsers[k] = v
		}
	}

	return allZitadelUsers, nil
}
