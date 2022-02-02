package migration

import (
	"strings"

	"github.com/caos/zitadel/pkg/databases/db"

	"github.com/caos/zitadel/operator/common"

	corev1 "k8s.io/api/core/v1"
)

func getPreContainer(
	dbConn db.Connection,
	customImageRegistry string,
) []corev1.Container {

	return []corev1.Container{
		{
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
		}, /*{
			Name:  "create-flyway-user",
			Image: common.CockroachImage.Reference(customImageRegistry),
			Env:   baseEnvVars(envMigrationUser, envMigrationPW, dbConn.User(), secretPasswordName),
			VolumeMounts: []corev1.VolumeMount{{
				Name:      rootUserInternal,
				MountPath: certsDir,
			}},
			Command: []string{"/bin/bash", "-c", "--"},
			Args: []string{
				strings.Join([]string{
					//					createUserCommand(envMigrationUser, envMigrationPW, createFile),
					grantUserCommand(envMigrationUser, grantFile),
					//					"cockroach.sh sql --certs-dir=/certificates --host=" + dbHost + ":" + dbPort + " -e \"$(cat " + createFile + ")\" -e \"$(cat " + grantFile + ")\";",
					fmt.Sprintf(`cockroach.sh sql --url='%s' -e "$(cat %s)" -e "$(cat %s)";`, dbConn.URL("/certificates"), createFile, grantFile),
				},
					";"),
			},
			TerminationMessagePath:   corev1.TerminationMessagePathDefault,
			TerminationMessagePolicy: "File",
			ImagePullPolicy:          "IfNotPresent",
		},*/
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
		"echo -n 'CREATE ROLE IF NOT EXISTS can_create_db WITH CREATEDB;' >> " + file,
		"echo -n 'GRANT can_create_db TO ' >> " + file,
		"echo -n ${" + user + "} >> " + file,
		"echo -n ' WITH ADMIN OPTION;'  >> " + file,
	}, ";")
}
