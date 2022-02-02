package migration

import (
	"fmt"
	"strings"

	"github.com/caos/zitadel/pkg/databases/db"

	"github.com/caos/zitadel/operator/common"

	corev1 "k8s.io/api/core/v1"
)

func getMigrationContainer(
	dbConn db.Connection,
	customImageRegistry string,
	certsVolumeMount corev1.VolumeMount,
) corev1.Container {

	pwSecretName, pwSecretKey := dbConn.PasswordSecret()

	return corev1.Container{
		Name:  "db-migration",
		Image: common.FlywayImage.Reference(customImageRegistry),
		Args: []string{
			fmt.Sprintf("-url=%s", connectionURL(dbConn, certsVolumeMount.MountPath)),
			fmt.Sprintf("-locations=filesystem:%s", migrationsPath),
			"migrate",
		},
		Env: migrationEnvVars(envMigrationUser, envMigrationPW, dbConn.User(), pwSecretName, pwSecretKey),
		VolumeMounts: []corev1.VolumeMount{certsVolumeMount, {
			Name:      migrationConfigmap,
			MountPath: migrationsPath,
		}},
		TerminationMessagePath:   corev1.TerminationMessagePathDefault,
		TerminationMessagePolicy: "File",
		ImagePullPolicy:          "IfNotPresent",
	}
}

func connectionURL(conn db.Connection, certsDir string) string {

	url := fmt.Sprintf("jdbc:postgresql://%s:%s/defaultdb?%s", conn.Host(), conn.Port(), sslParams(conn.SSL(), certsDir))

	options := conn.Options()
	if options != "" {
		url += "&options=" + options
	}

	return url

}

func sslParams(ssl *db.SSL, certsDir string) string {

	if ssl == nil {
		return "sslmode=disable"
	}

	params := "sslmode=verify-full"

	if ssl.RootCert {
		params += fmt.Sprintf("&sslrootcert=%s/%s", certsDir, db.RootCert)
	}

	if ssl.UserCertAndKey {
		params += fmt.Sprintf("&sslcert=%s/%s&sslkey=%s%s", certsDir, db.UserCert, certsDir, db.UserKey)
	}

	return params
}

func migrationEnvVars(envMigrationUser, envMigrationPW, migrationUser, userPasswordsSecret, userPasswordKey string) []corev1.EnvVar {
	envVars := baseEnvVars(envMigrationUser, envMigrationPW, migrationUser, userPasswordsSecret, userPasswordKey)

	vars := make([]corev1.EnvVar, 0)
	for _, v := range envVars {
		vars = append(vars, v)
	}

	deprecatedUsers := []string{
		"management",
		"adminapi",
		"auth",
		"authz",
		"notification",
		"eventstore",
		"queries",
	}
	for _, user := range deprecatedUsers {
		vars = append(vars, corev1.EnvVar{
			Name: "FLYWAY_PLACEHOLDERS_" + strings.ToUpper(user) + "PASSWORD",
			// TODO: Drop users in a new migration
			Value: "'to-be-deleted'",
			/*			ValueFrom: &corev1.EnvVarSource{
						SecretKeyRef: &corev1.SecretKeySelector{
							LocalObjectReference: corev1.LocalObjectReference{Name: userPasswordsSecret},
							Key:                  user,
						},
					},*/
		})
	}
	return vars
}
