package cockroachdb

import (
	"os/exec"
	"strings"
)

func GetBackupToS3(
	certsFolder string,
	host string,
	port string,
	bucketName string,
	filePath string,
	accessKeyIDPath string,
	secretAccessKeyPath string,
	sessionTokenPath string,
	endpoint string,
	region string,
) *exec.Cmd {
	parameters := []string{
		"AWS_ACCESS_KEY_ID=$(cat " + accessKeyIDPath + ")",
		"AWS_SECRET_ACCESS_KEY=$(cat " + secretAccessKeyPath + ")",
		"AWS_SESSION_TOKEN=$(cat " + sessionTokenPath + ")",
		"AWS_ENDPOINT=" + endpoint,
	}
	if region != "" {
		parameters = append(parameters, "AWS_REGION="+region)
	}

	return exec.Command(
		"cockroach",
		"sql",
		"--certs-dir="+certsFolder,
		"--host="+host,
		"--port="+port,
		"-e",
		"\"BACKUP TO \\\"s3://"+bucketName+"/"+filePath+"?"+strings.Join(parameters, "&")+"\\\";\"",
	)
}

func GetRestoreFromS3(
	certsFolder string,
	host string,
	port string,
	bucketName string,
	filePath string,
	accessKeyIDPath string,
	secretAccessKeyPath string,
	sessionTokenPath string,
	endpoint string,
	region string,
) *exec.Cmd {
	parameters := []string{
		"AWS_ACCESS_KEY_ID=$(cat " + accessKeyIDPath + ")",
		"AWS_SECRET_ACCESS_KEY=$(cat " + secretAccessKeyPath + ")",
		"AWS_SESSION_TOKEN=$(cat " + sessionTokenPath + ")",
		"AWS_ENDPOINT=" + endpoint,
	}
	if region != "" {
		parameters = append(parameters, "AWS_REGION="+region)
	}

	return exec.Command(
		"cockroach",
		"sql",
		"--certs-dir="+certsFolder,
		"--host="+host,
		"--port="+port,
		"-e",
		"\"RESTORE FROM \\\"s3://"+bucketName+"/"+filePath+"?"+strings.Join(parameters, "&")+"\\\";\"",
	)
}
