package command

import (
	"fmt"
	"github.com/caos/zitadel/pkg/databases/db"
	corev1 "k8s.io/api/core/v1"
	"strings"
)

type Action string

const (
	Backup  Action = "BACKUP TO"
	Restore Action = "RESTORE FROM"
)

func GetCommand(
	bucketName string,
	backupName string,
	backupTime string,
	certsFolder string,
	accessKeyIDPath string,
	secretAccessKeyPath string,
	sessionTokenPath string,
	region string,
	endpoint string,
	dbConn db.Connection,
	action Action,
) (cmd string, pw *corev1.EnvVar) {

	parameters := []string{
		"AWS_ACCESS_KEY_ID=$(cat " + accessKeyIDPath + ")",
		"AWS_SECRET_ACCESS_KEY=$(cat " + secretAccessKeyPath + ")",
		"AWS_SESSION_TOKEN=$(cat " + sessionTokenPath + ")",
		"AWS_ENDPOINT=" + endpoint,
	}
	if region != "" {
		parameters = append(parameters, "AWS_REGION="+region)
	}

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
	dbURL = fmt.Sprintf("%s@%s:%s/defaultdb", dbURL, dbConn.Host(), dbConn.Port())

	options := dbConn.Options()
	if options != "" {
		dbURL = fmt.Sprintf("%s?options=%s", dbURL, options)
	}

	return fmt.Sprintf(
		`cockroach sql --certs-dir=%s --url=%s --execute "%s \"s3://%s/%s/%s?%s\";"`,
		certsFolder,
		dbURL,
		action,
		bucketName,
		backupName,
		backupTime,
		strings.Join(parameters, "&"),
	), pw
}

func boolPrt(b bool) *bool { return &b }
