package migration

import (
	corev1 "k8s.io/api/core/v1"
	"strings"
)

func getPostContainers(
	dbHost string,
	dbPort string,
	migrationUser string,
	secretPasswordName string,
) []corev1.Container {
	return []corev1.Container{
		{
			Name:    "delete-flyway-user",
			Image:   "cockroachdb/cockroach:v20.2.3",
			Command: []string{"/bin/bash", "-c", "--"},
			Args: []string{
				strings.Join([]string{
					deleteUserCommand(envMigrationUser, deleteFile),
					"cockroach.sh sql --certs-dir=/certificates --host=" + dbHost + ":" + dbPort + " -e \"$(cat " + deleteFile + ")\";",
				}, ";"),
			},
			Env: baseEnvVars(envMigrationUser, envMigrationPW, migrationUser, secretPasswordName),
			VolumeMounts: []corev1.VolumeMount{{
				Name:      rootUserInternal,
				MountPath: rootUserPath,
			}},
			TerminationMessagePath:   corev1.TerminationMessagePathDefault,
			TerminationMessagePolicy: "File",
			ImagePullPolicy:          "IfNotPresent",
		},
	}
}

func deleteUserCommand(user, file string) string {
	return strings.Join([]string{
		"echo -n 'DROP USER IF EXISTS ' > " + file,
		"echo -n ${" + user + "} >> " + file,
		"echo -n ';' >> " + file,
	}, ";")
}
