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
) corev1.Container {

	pwSecretName, pwSecretKey := dbConn.PasswordSecret()

	return corev1.Container{
		Name:  "db-migration",
		Image: common.FlywayImage.Reference(customImageRegistry),
		Args: []string{
			fmt.Sprintf("-url=jdbc:postgresql://%s:%s/defaultdb?%s", dbConn.Host(), dbConn.Port(), dbConn.ConnectionParams(chownedCertsDir)),
			fmt.Sprintf("-locations=filesystem:%s", migrationsPath),
			"migrate",
		},
		Env: migrationEnvVars(envMigrationUser, envMigrationPW, dbConn.User(), pwSecretName, pwSecretKey),
		VolumeMounts: []corev1.VolumeMount{{
			Name:      migrationConfigmap,
			MountPath: migrationsPath,
		}, {
			Name:      chownedCertsVolumeName,
			MountPath: chownedCertsDir,
		}},
		TerminationMessagePath:   corev1.TerminationMessagePathDefault,
		TerminationMessagePolicy: "File",
		ImagePullPolicy:          "IfNotPresent",
	}
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
			// TODO: Delete users in a new migration
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
