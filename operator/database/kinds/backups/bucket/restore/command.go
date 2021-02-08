package restore

import "strings"

func getCommand(
	timestamp string,
	databases []string,
	bucketName string,
	backupName string,

) string {

	backupCommands := make([]string, 0)
	for _, database := range databases {
		backupCommands = append(backupCommands,
			strings.Join([]string{
				"/scripts/restore.sh",
				bucketName,
				backupName,
				timestamp,
				database,
				secretPath,
				certPath,
			}, " "))
	}

	return strings.Join(backupCommands, " && ")
}
