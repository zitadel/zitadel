package zitadel

import (
	"sort"

	"github.com/caos/zitadel/operator/kinds/iam/zitadel/configuration"
)

const migrationUser = "flyway"

func getAllUsers(desired *DesiredV0) map[string]string {
	passwords := &configuration.Passwords{}
	if desired != nil && desired.Spec != nil && desired.Spec.Configuration != nil && desired.Spec.Configuration.Passwords != nil {
		passwords = desired.Spec.Configuration.Passwords
	}
	users := make(map[string]string, 0)

	migrationPassword := migrationUser
	if passwords.Migration != nil {
		migrationPassword = passwords.Migration.Value
	}
	users[migrationUser] = migrationPassword

	mgmtUser := "management"
	mgmtPassword := mgmtUser
	if passwords != nil && passwords.Management != nil {
		mgmtPassword = passwords.Management.Value
	}
	users[mgmtUser] = mgmtPassword

	adminUser := "adminapi"
	adminPassword := adminUser
	if passwords != nil && passwords.Adminapi != nil {
		adminPassword = passwords.Adminapi.Value
	}
	users[adminUser] = adminPassword

	authUser := "auth"
	authPassword := authUser
	if passwords != nil && passwords.Auth != nil {
		authPassword = passwords.Auth.Value
	}
	users[authUser] = authPassword

	authzUser := "authz"
	authzPassword := authzUser
	if passwords != nil && passwords.Authz != nil {
		authzPassword = passwords.Authz.Value
	}
	users[authzUser] = authzPassword

	notUser := "notification"
	notPassword := notUser
	if passwords != nil && passwords.Notification != nil {
		notPassword = passwords.Notification.Value
	}
	users[notUser] = notPassword

	esUser := "eventstore"
	esPassword := esUser
	if passwords != nil && passwords.Eventstore != nil {
		esPassword = passwords.Eventstore.Value
	}
	users[esUser] = esPassword

	return users
}

func getZitadelUserList() []string {
	allUsersMap := getAllUsers(nil)

	allZitadelUsers := make([]string, 0)
	for k := range allUsersMap {
		if k != migrationUser {
			allZitadelUsers = append(allZitadelUsers, k)
		}
	}
	sort.Slice(allZitadelUsers, func(i, j int) bool {
		return allZitadelUsers[i] < allZitadelUsers[j]
	})

	return allZitadelUsers
}
