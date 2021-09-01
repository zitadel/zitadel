package cockroachdb

import "os/exec"

func GetBackupToGCS(
	certsFolder string,
	host string,
	port string,
	bucketName string,
	filePath string,
	serviceAccountPath string,
) *exec.Cmd {
	return exec.Command(
		"cockroach",
		"sql",
		"--certs-dir="+certsFolder,
		"--host="+host,
		"--port="+port,
		"-e",
		"\"BACKUP TO \\\"gs://"+bucketName+"/"+filePath+"?AUTH=specified&CREDENTIALS=$(cat "+serviceAccountPath+" | base64 | tr -d '\\n' )\\\";\"",
	)
}

func GetRestoreFromGCS(
	certsFolder string,
	host string,
	port string,
	bucketName string,
	filePath string,
	serviceAccountPath string,
) *exec.Cmd {
	return exec.Command(
		"cockroach",
		"sql",
		"--certs-dir="+certsFolder,
		"--host="+host,
		"--port="+port,
		"-e",
		"\"RESTORE FROM \\\"gs://"+bucketName+"/"+filePath+"?AUTH=specified&CREDENTIALS=$(cat "+serviceAccountPath+" | base64 | tr -d '\\n' )\\\";\"",
	)
}
