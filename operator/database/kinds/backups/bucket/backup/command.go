package backup

import "strings"

func getBackupCommand(
	timestamp string,
	databases []string,
	bucketName string,
	backupName string,
) string {

	backupCommands := make([]string, 0)
	if timestamp != "" {
		backupCommands = append(backupCommands, "export "+backupNameEnv+"="+timestamp)
	} else {
		backupCommands = append(backupCommands, "export "+backupNameEnv+"=$(date +%Y-%m-%dT%H:%M:%SZ)")
	}

	for _, database := range databases {
		backupCommands = append(backupCommands,
			strings.Join([]string{
				"/scripts/backup.sh",
				backupName,
				bucketName,
				database,
				backupPath,
				secretPath,
				certPath,
				"${" + backupNameEnv + "}",
			}, " "))
	}
	return strings.Join(backupCommands, " && ")
}
