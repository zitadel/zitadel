package backup

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBackup_Command1(t *testing.T) {
	timestamp := ""
	bucketName := "test"
	backupName := "test"
	dbURL := "testDB"
	dbPort := int32(80)

	cmd := getBackupCommand(
		timestamp,
		bucketName,
		backupName,
		certPath,
		secretPath,
		dbURL,
		dbPort,
	)
	equals := "export " + backupNameEnv + "=$(date +%Y-%m-%dT%H:%M:%SZ) && export SAJSON=$(cat /secrets/sa.json | base64 | tr -d '\n' ) && cockroach sql --certs-dir=/cockroach/cockroach-certs --host=testDB --port=80 -e \"BACKUP TO \\\"gs://test/test/${BACKUP_NAME}?AUTH=specified&CREDENTIALS=${SAJSON}\\\";\""
	assert.Equal(t, equals, cmd)
}

func TestBackup_Command2(t *testing.T) {
	timestamp := "test"
	bucketName := "test"
	backupName := "test"
	dbURL := "testDB"
	dbPort := int32(80)

	cmd := getBackupCommand(
		timestamp,
		bucketName,
		backupName,
		certPath,
		secretPath,
		dbURL,
		dbPort,
	)
	equals := "export " + backupNameEnv + "=test && export SAJSON=$(cat /secrets/sa.json | base64 | tr -d '\n' ) && cockroach sql --certs-dir=/cockroach/cockroach-certs --host=testDB --port=80 -e \"BACKUP TO \\\"gs://test/test/${BACKUP_NAME}?AUTH=specified&CREDENTIALS=${SAJSON}\\\";\""
	assert.Equal(t, equals, cmd)
}
