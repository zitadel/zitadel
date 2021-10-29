package core

import "strings"

const (
	jobNamePrefix        = "backup-"
	jobNameReleaseSuffix = "-restore"
)

func GetBackupJobName(backupName string) string {
	return jobNamePrefix + backupName
}

func TrimBackupJobName(name string) string {
	return strings.TrimPrefix(name, jobNamePrefix)
}

func GetRestoreJobName(backupName string) string {
	return jobNamePrefix + backupName + jobNameReleaseSuffix
}

func GetSecretName(backupName string) string {
	return jobNamePrefix + backupName
}
