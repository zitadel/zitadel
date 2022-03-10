package restore

import (
	"strconv"
	"strings"
)

func getCommand(
	timestamp string,
	bucketName string,
	backupName string,
	certsFolder string,
	serviceAccountPath string,
	dbURL string,
	dbPort int32,
) string {

	backupCommands := make([]string, 0)

	backupCommands = append(backupCommands, "export "+saJsonBase64Env+"=$(cat "+serviceAccountPath+" | base64 | tr -d '\n' )")

	backupCommands = append(backupCommands,
		strings.Join([]string{
			"cockroach",
			"sql",
			"--certs-dir=" + certsFolder,
			"--host=" + dbURL,
			"--port=" + strconv.Itoa(int(dbPort)),
			"-e",
			"\"RESTORE FROM \\\"gs://" + bucketName + "/" + backupName + "/" + timestamp + "?AUTH=specified&CREDENTIALS=${" + saJsonBase64Env + "}\\\";\"",
		}, " ",
		),
	)

	return strings.Join(backupCommands, " && ")
}
