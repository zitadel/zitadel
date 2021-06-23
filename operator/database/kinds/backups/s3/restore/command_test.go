package restore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBackup_Command1(t *testing.T) {
	timestamp := "test1"
	bucketName := "testBucket"
	backupName := "testBackup"
	dbURL := "testDB"
	dbPort := int32(80)
	region := "region"
	endpoint := "endpoint"

	cmd := getCommand(
		timestamp,
		bucketName,
		backupName,
		certPath,
		accessKeyIDPath,
		secretAccessKeyPath,
		sessionTokenPath,
		region,
		endpoint,
		dbURL,
		dbPort,
	)

	equals := "cockroach sql --certs-dir=" + certPath + " --host=testDB --port=80 -e \"RESTORE FROM \\\"s3://testBucket/testBackup/test1?AWS_ACCESS_KEY_ID=$(cat " + accessKeyIDPath + ")&AWS_SECRET_ACCESS_KEY=$(cat " + secretAccessKeyPath + ")&AWS_SESSION_TOKEN=$(cat " + sessionTokenPath + ")&AWS_ENDPOINT=endpoint&AWS_REGION=region\\\";\""
	assert.Equal(t, equals, cmd)
}

func TestBackup_Command2(t *testing.T) {
	timestamp := "test2"
	bucketName := "testBucket"
	backupName := "testBackup"
	dbURL := "testDB2"
	dbPort := int32(81)
	region := "region2"
	endpoint := "endpoint2"

	cmd := getCommand(
		timestamp,
		bucketName,
		backupName,
		certPath,
		accessKeyIDPath,
		secretAccessKeyPath,
		sessionTokenPath,
		region,
		endpoint,
		dbURL,
		dbPort,
	)
	equals := "cockroach sql --certs-dir=" + certPath + " --host=testDB2 --port=81 -e \"RESTORE FROM \\\"s3://testBucket/testBackup/test2?AWS_ACCESS_KEY_ID=$(cat " + accessKeyIDPath + ")&AWS_SECRET_ACCESS_KEY=$(cat " + secretAccessKeyPath + ")&AWS_SESSION_TOKEN=$(cat " + sessionTokenPath + ")&AWS_ENDPOINT=endpoint2&AWS_REGION=region2\\\";\""
	assert.Equal(t, equals, cmd)
}
