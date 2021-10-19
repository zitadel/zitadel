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
	accessKeyIDPath string,
	secretAccessKeyPath string,
	sessionTokenPath string,
	region string,
	endpoint string,
	dbURL string,
	dbPort int32,
) string {

	backupCommands := make([]string, 0)

	parameters := []string{
		"AWS_ACCESS_KEY_ID=$(cat " + accessKeyIDPath + ")",
		"AWS_SECRET_ACCESS_KEY=$(cat " + secretAccessKeyPath + ")",
		"AWS_SESSION_TOKEN=$(cat " + sessionTokenPath + ")",
		"AWS_ENDPOINT=" + endpoint,
	}
	if region != "" {
		parameters = append(parameters, "AWS_REGION="+region)
	}

	backupCommands = append(backupCommands,
		strings.Join([]string{
			"cockroach",
			"sql",
			"--certs-dir=" + certsFolder,
			"--host=" + dbURL,
			"--port=" + strconv.Itoa(int(dbPort)),
			"-e",
			"\"RESTORE FROM \\\"s3://" + bucketName + "/" + backupName + "/" + timestamp + "?" + strings.Join(parameters, "&") + "\\\";\"",
		}, " ",
		),
	)

	return strings.Join(backupCommands, " && ")
}
