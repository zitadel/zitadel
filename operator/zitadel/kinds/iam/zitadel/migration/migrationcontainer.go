package migration

import (
	"strings"

	"github.com/caos/zitadel/operator/common"

	corev1 "k8s.io/api/core/v1"
)

func getMigrationContainer(
	dbHost string,
	dbPort string,
	migrationUser string,
	secretPasswordName string,
	users []string,
	customImageRegistry string,
) corev1.Container {

	// TODO: Parameterize
	insecure := false

	var rootCertPath string
	volumeMounts := []corev1.VolumeMount{{
		Name:      migrationConfigmap,
		MountPath: migrationsPath,
	}}
	if !insecure {
		rootCertPath = rootUserPath
		volumeMounts = append(volumeMounts, corev1.VolumeMount{
			Name:      rootUserInternal,
			MountPath: rootUserPath,
		})
	}

	return corev1.Container{
		Name:  "db-migration",
		Image: common.FlywayImage.Reference(customImageRegistry),
		Args: []string{
			"-url=" + connectionString(dbHost, dbPort, rootCertPath),
			"-locations=filesystem:" + migrationsPath,
			"migrate",
		},
		Env:                      migrationEnvVars(envMigrationUser, envMigrationPW, migrationUser, secretPasswordName, users),
		VolumeMounts:             volumeMounts,
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

func connectionString(dbHost, dbPort, rootCertPath string) string {
	location := "jdbc:postgresql://" + dbHost + ":" + dbPort + "/defaultdb"
	if rootCertPath != "" {
		location += "?sslmode=verify-full&ssl=true&sslrootcert=" + rootCertPath + "/ca.crt&sslfactory=org.postgresql.ssl.NonValidatingFactory"
	}
	return location
}
