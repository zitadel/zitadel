package clean

import "strings"

func getCommand(
	databases []string,
) string {
	backupCommands := make([]string, 0)
	for _, database := range databases {
		backupCommands = append(backupCommands,
			strings.Join([]string{
				"/scripts/clean-db.sh",
				certPath,
				database,
			}, " "))
	}
	for _, database := range databases {
		backupCommands = append(backupCommands,
			strings.Join([]string{
				"/scripts/clean-user.sh",
				certPath,
				database,
			}, " "))
	}
	backupCommands = append(backupCommands,
		strings.Join([]string{
			"/scripts/clean-migration.sh",
			certPath,
		}, " "))

	return strings.Join(backupCommands, " && ")
}
