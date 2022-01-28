package migration

import (
	"fmt"
	"github.com/caos/zitadel/pkg/databases/db"

	"github.com/caos/zitadel/operator/common"

	corev1 "k8s.io/api/core/v1"
)

func getMigrationContainer(
	dbConn db.Connection,
	customImageRegistry string,
) corev1.Container {

	pwSecretName, pwSecretKey := dbConn.PasswordSecret()

	/*	var pwSecretEnv string
		if pwSecretName != "" {
			pwSecretEnv = envMigrationPW
		}

	*/

	//			"-url=jdbc:postgresql://" + dbHost + ":" + dbPort + "/defaultdb?&sslmode=verify-full&ssl=true&sslrootcert=" + certsDir + "/ca.crt&sslfactory=org.postgresql.ssl.NonValidatingFactory"},
	/*		Args: []string{
			"-url=jdbc:" + dbConn.URL(certsDir, pwSecretEnv),
			//			"-url=jdbc:postgresql://" + dbHost + ":" + dbPort + "/defaultdb?&sslmode=verify-full&ssl=true&sslrootcert=" + certsDir + "/ca.crt&sslfactory=org.postgresql.ssl.NonValidatingFactory",
			"-locations=filesystem:" + migrationsPath,
			"migrate",
		},*/

	return corev1.Container{
		Name:  "db-migration",
		Image: common.FlywayImage.Reference(customImageRegistry),
		Args: []string{
			fmt.Sprintf("-url=jdbc:postgresql://%s:%s/defaultdb?%s", dbConn.Host(), dbConn.Port(), dbConn.ConnectionParams(chownedCertsDir)),
			fmt.Sprintf("-locations=filesystem:%s", migrationsPath),
			"migrate",
		},
		Env: baseEnvVars(envMigrationUser, envMigrationPW, dbConn.User(), pwSecretName, pwSecretKey),
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

/*
func migrationEnvVars(envMigrationUser, envMigrationPW, migrationUser, userPasswordsSecret string, users []string) []corev1.EnvVar {
	envVars := baseEnvVars(envMigrationUser, envMigrationPW, migrationUser, userPasswordsSecret, migrationUser)

	migrationEnvVars := make([]corev1.EnvVar, 0)
	for _, v := range envVars {
		migrationEnvVars = append(migrationEnvVars, v)
	}
	for _, user := range users {
		migrationEnvVars = append(migrationEnvVars, corev1.EnvVar{
			Name: "FLYWAY_PLACEHOLDERS_" + strings.ToUpper(user) + "PASSWORD",
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{Name: userPasswordsSecret},
					Key:                  user,
				},
			},
		})
	}
	return migrationEnvVars
}
*/
