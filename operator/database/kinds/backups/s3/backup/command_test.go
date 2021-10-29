package backup

import (
	"github.com/stretchr/testify/assert"
	"path/filepath"
	"testing"
)

func TestBackup_Command1(t *testing.T) {
	timestamp := ""
	bucketName := "test"
	backupName := "test"
	dbURL := "testDB"
	dbPort := int32(80)
	region := "region"
	endpoint := "endpoint"
	akidPath := filepath.Join(secretsPath, "akid")
	sakPath := filepath.Join(secretsPath, "sak")
	stPath := filepath.Join(secretsPath, "st")

	cmd := getBackupCommand(
		timestamp,
		bucketName,
		backupName,
		certPath,
		akidPath,
		sakPath,
		stPath,
		region,
		endpoint,
		dbURL,
		dbPort,
	)
	equals := "export " + backupNameEnv + "=$(date +%Y-%m-%dT%H:%M:%SZ) && cockroach sql --certs-dir=" + certPath + " --host=testDB --port=80 -e \"BACKUP TO \\\"s3://test/test/${BACKUP_NAME}?AWS_ACCESS_KEY_ID=$(cat " + akidPath + ")&AWS_SECRET_ACCESS_KEY=$(cat " + sakPath + ")&AWS_SESSION_TOKEN=$(cat " + stPath + ")&AWS_ENDPOINT=endpoint&AWS_REGION=region\\\";\""
	assert.Equal(t, equals, cmd)
}

func TestBackup_Command2(t *testing.T) {
	timestamp := "test"
	bucketName := "test"
	backupName := "test"
	dbURL := "testDB"
	dbPort := int32(80)
	region := "region"
	endpoint := "endpoint"
	akidPath := filepath.Join(secretsPath, "akid")
	sakPath := filepath.Join(secretsPath, "sak")
	stPath := filepath.Join(secretsPath, "st")

	cmd := getBackupCommand(
		timestamp,
		bucketName,
		backupName,
		certPath,
		akidPath,
		sakPath,
		stPath,
		region,
		endpoint,
		dbURL,
		dbPort,
	)
	equals := "export " + backupNameEnv + "=test && cockroach sql --certs-dir=" + certPath + " --host=testDB --port=80 -e \"BACKUP TO \\\"s3://test/test/${BACKUP_NAME}?AWS_ACCESS_KEY_ID=$(cat " + akidPath + ")&AWS_SECRET_ACCESS_KEY=$(cat " + sakPath + ")&AWS_SESSION_TOKEN=$(cat " + stPath + ")&AWS_ENDPOINT=endpoint&AWS_REGION=region\\\";\""
	assert.Equal(t, equals, cmd)
}
