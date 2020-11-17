package migration

import (
	corev1 "k8s.io/api/core/v1"
	"strings"
)

func getMigrationContainer(
	dbHost string,
	dbPort string,
	migrationUser string,
	secretPasswordName string,
	users []string,
) corev1.Container {
	return corev1.Container{
		Name:  "db-migration",
		Image: "flyway/flyway:6.5.0",
		Args: []string{
			"-url=jdbc:postgresql://" + dbHost + ":" + dbPort + "/defaultdb?&sslmode=verify-full&ssl=true&sslrootcert=" + rootUserPath + "/ca.crt&sslfactory=org.postgresql.ssl.NonValidatingFactory",
			"-locations=filesystem:" + migrationsPath,
			"migrate",
		},
		Env: migrationEnvVars(envMigrationUser, envMigrationPW, migrationUser, secretPasswordName, users),
		VolumeMounts: []corev1.VolumeMount{{
			Name:      migrationConfigmap,
			MountPath: migrationsPath,
		}, {
			Name:      rootUserInternal,
			MountPath: rootUserPath,
		}},
		TerminationMessagePath:   corev1.TerminationMessagePathDefault,
		TerminationMessagePolicy: "File",
		ImagePullPolicy:          "IfNotPresent",
	}
}

func migrationEnvVars(envMigrationUser, envMigrationPW, migrationUser, userPasswordsSecret string, users []string) []corev1.EnvVar {
	envVars := baseEnvVars(envMigrationUser, envMigrationPW, migrationUser, userPasswordsSecret)

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
