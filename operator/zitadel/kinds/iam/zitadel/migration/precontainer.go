package migration

import (
	"github.com/caos/zitadel/operator/helpers"
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
	dbCertsFlyway string,
) []corev1.Container {
	certsCockroach := common.GetInitContainer(
		"cockroach",
		rootUserInternal,
		dbCertsCockroach,
		[]string{"root"},
		common.ZITADELCockroachImage.RunAsUser(),
		common.ZITADELCockroachImage.Reference(customImageRegistry, version),
	)
	certsFlyway := common.GetInitContainer(
		"flyway",
		rootUserInternal,
		dbCertsFlyway,
		[]string{"root"},
		common.ZITADELImage.RunAsUser(),
		common.ZITADELCockroachImage.Reference(customImageRegistry, version),
	)
	runAsUser := common.ZITADELCockroachImage.RunAsUser()
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
				RunAsUser:  helpers.PointerInt64(runAsUser),
				RunAsGroup: helpers.PointerInt64(runAsUser),
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
