package migration

import (
	"github.com/caos/zitadel/operator/helpers"
	"github.com/caos/zitadel/operator/zitadel/kinds/iam/zitadel/deployment"
	"strings"

	"github.com/caos/zitadel/operator/common"

	corev1 "k8s.io/api/core/v1"
)

const certTempMountPath = "/tmp/certs"

func getPreContainer(
	dbHost string,
	dbPort string,
	migrationUser string,
	secretPasswordName string,
	customImageRegistry string,
	version string,
	dbCerts string,
) []corev1.Container {

	return []corev1.Container{
		deployment.GetInitContainer(
			rootUserInternal,
			dbCerts,
			[]string{"root"},
			1000,
			customImageRegistry,
			version,
		),
		/*
			{
				Name:  "check-db-ready",
				Image: common.PostgresImage.Reference(customImageRegistry),
				Command: []string{
					"sh",
					"-c",
					"until pg_isready -h " + dbHost + " -p " + dbPort + "; do echo waiting for database; sleep 2; done;",
				},
				SecurityContext: &corev1.SecurityContext{
					RunAsUser:    helpers.PointerInt64(70),
					RunAsGroup:   helpers.PointerInt64(70),
					RunAsNonRoot: helpers.PointerBool(true),
				},
				TerminationMessagePath:   corev1.TerminationMessagePathDefault,
				TerminationMessagePolicy: "File",
				ImagePullPolicy:          "IfNotPresent",
			},*/
		{
			Name:  "create-flyway-user",
			Image: common.BackupImage.Reference(customImageRegistry, version),
			Env:   baseEnvVars(envMigrationUser, envMigrationPW, migrationUser, secretPasswordName),
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      dbCerts,
					MountPath: certTempMountPath,
				},
			},
			Command: []string{"/bin/bash", "-c", "--"},
			Args: []string{
				strings.Join([]string{
					createUserCommand(envMigrationUser, envMigrationPW, createFile),
					grantUserCommand(envMigrationUser, grantFile),
					"cockroach.sh sql --certs-dir=" + certTempMountPath + " --host=" + dbHost + ":" + dbPort + " -e \"$(cat " + createFile + ")\" -e \"$(cat " + grantFile + ")\";",
				},
					";"),
			},
			SecurityContext: &corev1.SecurityContext{
				RunAsUser:    helpers.PointerInt64(1000),
				RunAsGroup:   helpers.PointerInt64(1000),
				RunAsNonRoot: helpers.PointerBool(true),
			},
			TerminationMessagePath:   corev1.TerminationMessagePathDefault,
			TerminationMessagePolicy: "File",
			ImagePullPolicy:          "IfNotPresent",
		},
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
