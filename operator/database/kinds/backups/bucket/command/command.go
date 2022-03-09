package command

import (
	"fmt"
	"github.com/caos/zitadel/pkg/databases/db"
	corev1 "k8s.io/api/core/v1"
)

type Action string

const (
	Backup  Action = "BACKUP TO"
	Restore Action = "RESTORE FROM"
)

func GetBackupRestoreStatement(
	bucketName string,
	backupName string,
	backupTime string,
	serviceAccountPath string,
	action Action,
) string {
	return fmt.Sprintf(`%s \"gs://%s/%s/%s?AUTH=specified&CREDENTIALS=$(cat %s | base64 | tr -d '\n' )\";`,
		action,
		bucketName,
		backupName,
		backupTime,
		serviceAccountPath,
	)
}

func GetSQLCommand(
	dbConn db.Connection,
	certsFolder string,
	statements ...string,
) (cmd string, pw *corev1.EnvVar) {

	dbURL := "postgres://" + dbConn.User()

	pwSecret, pwSecretKey := dbConn.PasswordSecret()
	pwEnv := "CR_PASSWORD"
	if pwSecret != nil {
		dbURL = fmt.Sprintf("%s:${%s}", dbURL, pwEnv)
		pw = &corev1.EnvVar{
			Name: pwEnv,
			ValueFrom: &corev1.EnvVarSource{
				SecretKeyRef: &corev1.SecretKeySelector{
					LocalObjectReference: corev1.LocalObjectReference{
						Name: pwSecret.Name(),
					},
					Key:      pwSecretKey,
					Optional: boolPrt(false),
				},
			},
		}
	}
	dbURL = fmt.Sprintf("%s@%s:%s/", dbURL, dbConn.Host(), dbConn.Port())

	options := dbConn.Options()
	if options != "" {
		dbURL = fmt.Sprintf("%s?options=%s", dbURL, options)
	}

	cmd = fmt.Sprintf("cockroach sql --certs-dir=%s --url=%s", certsFolder, dbURL)
	for _, statement := range statements {
		cmd += fmt.Sprintf(` --execute "%s"`, statement)
	}
	return cmd, pw
}

func boolPrt(b bool) *bool { return &b }
