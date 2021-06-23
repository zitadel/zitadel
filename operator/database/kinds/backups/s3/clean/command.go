package clean

import "strings"

func getCommand(
	databases []string,
	users []string,
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
	for _, user := range users {
		backupCommands = append(backupCommands,
			strings.Join([]string{
				"/scripts/clean-user.sh",
				certPath,
				user,
			}, " "))
	}
	backupCommands = append(backupCommands,
		strings.Join([]string{
			"/scripts/clean-migration.sh",
			certPath,
		}, " "))

	return strings.Join(backupCommands, " && ")
}
