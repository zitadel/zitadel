package backup

import (
	"strconv"
	"strings"
)

func getBackupCommand(
	timestamp string,
	bucketName string,
	backupName string,
	certsFolder string,
	serviceAccountPath string,
	dbURL string,
	dbPort int32,
) string {

	backupCommands := make([]string, 0)
	if timestamp != "" {
		backupCommands = append(backupCommands, "export "+backupNameEnv+"="+timestamp)
	} else {
		backupCommands = append(backupCommands, "export "+backupNameEnv+"=$(date +%Y-%m-%dT%H:%M:%SZ)")
	}

	backupCommands = append(backupCommands, "export "+saJsonBase64Env+"=$(cat "+serviceAccountPath+" | base64 | tr -d '\n' )")

	backupCommands = append(backupCommands,
		strings.Join([]string{
			"cockroach",
			"sql",
			"--certs-dir=" + certsFolder,
			"--host=" + dbURL,
			"--port=" + strconv.Itoa(int(dbPort)),
			"-e",
			"\"BACKUP TO \\\"gs://" + bucketName + "/" + backupName + "/${" + backupNameEnv + "}?AUTH=specified&CREDENTIALS=${" + saJsonBase64Env + "}\\\";\"",
		}, " ",
		),
	)

	return strings.Join(backupCommands, " && ")
}
