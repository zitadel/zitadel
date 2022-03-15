package migration

import (
	"fmt"
	"strings"

	"github.com/caos/zitadel/pkg/databases/db"

	"github.com/caos/zitadel/operator/common"

	corev1 "k8s.io/api/core/v1"
)

func getReadyPreContainer(
	dbConn db.Connection,
	customImageRegistry string,
) corev1.Container {
	return corev1.Container{
		Name:  "check-db-ready",
		Image: common.PostgresImage.Reference(customImageRegistry),
		Command: []string{
			"sh",
			"-c",
			"until pg_isready -h " + dbConn.Host() + " -p " + dbConn.Port() + "; do echo waiting for database; sleep 2; done;",
		},
		TerminationMessagePath:   corev1.TerminationMessagePathDefault,
		TerminationMessagePolicy: "File",
		ImagePullPolicy:          "IfNotPresent",
	}
}

func getFlywayUserPreContainer(
	dbConn db.Connection,
	customImageRegistry string,
	migrationUser string,
	secretPasswordName string,
	certsVolumemount corev1.VolumeMount,
) corev1.Container {

	migrationUserPasswordSecret, migrationUserPasswordSecretKey := dbConn.PasswordSecret()
	var migrationUserPasswordSecretName string
	if migrationUserPasswordSecret != nil {
		migrationUserPasswordSecretName = migrationUserPasswordSecret.Name()
	}

	return corev1.Container{
		Name:         "create-flyway-user",
		Image:        common.CockroachImage.Reference(customImageRegistry),
		Env:          baseEnvVars(envMigrationUser, envMigrationPW, migrationUser, migrationUserPasswordSecretName, migrationUserPasswordSecretKey),
		VolumeMounts: []corev1.VolumeMount{certsVolumemount},
		Command:      []string{"/bin/bash", "-c", "--"},
		Args: []string{
			strings.Join([]string{
				createUserCommand(envMigrationUser, envMigrationPW, createFile),
				grantUserCommand(envMigrationUser, grantFile),
				"cockroach.sh sql --url=" + connectionURL(dbConn, certsVolumemount.MountPath) + " -e \"$(cat " + createFile + ")\" -e \"$(cat " + grantFile + ")\";",
			},
				";"),
		},
		TerminationMessagePath:   corev1.TerminationMessagePathDefault,
		TerminationMessagePolicy: "File",
		ImagePullPolicy:          "IfNotPresent",
	}
}

func createUserCommand(user, pw, file string) string {
	if user == "" || file == "" {
		return ""
	}

	createUser := strings.Join([]string{
		"echo -n 'CREATE USER IF NOT EXISTS ' > " + file,
		"echo -n ${" + user + "} >> " + file,
		"echo -n ';' >> " + file,
	}, ";")

	if pw != "" {
		createUser = strings.Join([]string{
			createUser,
			"echo -n 'ALTER USER ' >> " + file,
			"echo -n ${" + user + "} >> " + file,
			"echo -n ' WITH PASSWORD ' >> " + file,
			"echo -n ${" + pw + "} >> " + file,
			"echo -n ';' >> " + file,
		}, ";")
	}

	return createUser
}

func grantUserCommand(user, file string) string {
	return strings.Join([]string{
		"echo -n 'GRANT admin TO ' > " + file,
		"echo -n ${" + user + "} >> " + file,
		"echo -n ' WITH ADMIN OPTION;'  >> " + file,
	}, ";")
}

func connectionURL(conn db.Connection, certsDir string) string {

	url := fmt.Sprintf("jdbc:postgresql://%s:%s/defaultdb?%s", conn.Host(), conn.Port(), sslParams(conn.SSL(), conn.User(), certsDir))

	options := conn.Options()
	if options != "" {
		url += "&options=" + options
	}

	return url

}

func sslParams(ssl *db.SSL, user, certsDir string) string {

	if ssl == nil {
		return "sslmode=disable"
	}

	params := "sslmode=verify-full&ssl=true"

	if ssl.RootCert {
		params += fmt.Sprintf("&sslrootcert=%s/%s", certsDir, db.CACert)
	}

	if ssl.UserCertAndKey {
		params += fmt.Sprintf("&sslcert=%s/%s&sslkey=%s/%s&sslfactory=org.postgresql.ssl.NonValidatingFactory", certsDir, db.UserCert(user), certsDir, db.UserKey(user))
	}

	return params
}
