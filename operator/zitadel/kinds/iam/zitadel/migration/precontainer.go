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
	dbCertsCockroach string,
	runAsUserCockroach int64,
	dbCertsFlyway string,
	runAsUserFlyway int64,
) []corev1.Container {
	certsCockroach := deployment.GetInitContainer(
		"cockroach",
		rootUserInternal,
		dbCertsCockroach,
		[]string{"root"},
		runAsUserCockroach,
		customImageRegistry,
		version,
	)
	certsFlyway := deployment.GetInitContainer(
		"flyway",
		rootUserInternal,
		dbCertsFlyway,
		[]string{"root"},
		runAsUserFlyway,
		customImageRegistry,
		version,
	)
	return []corev1.Container{
		certsCockroach,
		certsFlyway,
		{
			Name:  "create-flyway-user",
			Image: common.ZITADELCockroachImage.Reference(customImageRegistry, version),
			Env:   baseEnvVars(envMigrationUser, envMigrationPW, migrationUser, secretPasswordName),
			VolumeMounts: []corev1.VolumeMount{
				{
					Name:      dbCertsCockroach,
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
				RunAsUser:  helpers.PointerInt64(runAsUserCockroach),
				RunAsGroup: helpers.PointerInt64(runAsUserCockroach),
			},
			TerminationMessagePath:   corev1.TerminationMessagePathDefault,
			TerminationMessagePolicy: corev1.TerminationMessageReadFile,
			ImagePullPolicy:          corev1.PullIfNotPresent,
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

	createUser = strings.Join([]string{
		createUser,
		"chmod +xr " + file,
	}, ";")

	return createUser
}

func grantUserCommand(user, file string) string {
	return strings.Join([]string{
		"echo -n 'GRANT admin TO ' > " + file,
		"echo -n ${" + user + "} >> " + file,
		"echo -n ' WITH ADMIN OPTION;'  >> " + file,
		"chmod +xr " + file,
	}, ";")
}
