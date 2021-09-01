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
	assetEndpoint string,
	assetPrefix string,
) string {

	backupCommands := make([]string, 0)

	backupCommands = append(backupCommands,
		strings.Join([]string{
			"backupctl",
			"restore",
			"gcs",
			"--backupname=" + backupName,
			"--backupnameenv=" + backupNameEnv,
			"--asset-endpoint=" + assetEndpoint,
			"--asset-akid=$(cat " + akidSecretPath + ")",
			"--asset-sak=$(cat " + sakSecretPath + ")",
			"--host=" + dbURL,
			"--port=" + strconv.Itoa(int(dbPort)),
			"--source-sajsonpath=" + serviceAccountPath,
			"--source-bucket" + bucketName,
			"--certs-dir=" + certsFolder,
		}, " ",
		),
	)

	return strings.Join(backupCommands, " && ")
}
