package backup

import (
	"path/filepath"
	"strings"
)

func getBackupCommand(
	timestamp string,
	bucketName string,
	backupName string,
) string {

	backupCommands := make([]string, 0)
	if timestamp != "" {
		backupCommands = append(backupCommands, "export "+backupNameEnv+"="+timestamp)
	} else {
		backupCommands = append(backupCommands, "export "+backupNameEnv+"=$(date +%Y-%m-%dT%H:%M:%SZ)")
	}

	backupCommands = append(backupCommands,
		strings.Join([]string{
			"rclone",
			"--no-check-certificate",
			"--config",
			configSecretPath,
			"sync",
			sourceName + ":" + bucketName,
			destinationName + ":" + filepath.Join(bucketName, cronJobNamePrefix+"-${"+backupName+"}"),
		}, " "))

	return strings.Join(backupCommands, " && ")
}
