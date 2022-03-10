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

	cmd := getCommand(
		timestamp,
		bucketName,
		backupName,
		certPath,
		secretPath,
		dbURL,
		dbPort,
	)

	equals := "export SAJSON=$(cat /secrets/sa.json | base64 | tr -d '\n' ) && cockroach sql --certs-dir=/cockroach/cockroach-certs --host=testDB --port=80 -e \"RESTORE FROM \\\"gs://testBucket/testBackup/test1?AUTH=specified&CREDENTIALS=${SAJSON}\\\";\""
	assert.Equal(t, equals, cmd)
}

func TestBackup_Command2(t *testing.T) {
	timestamp := "test2"
	bucketName := "testBucket"
	backupName := "testBackup"
	dbURL := "testDB2"
	dbPort := int32(81)

	cmd := getCommand(
		timestamp,
		bucketName,
		backupName,
		certPath,
		secretPath,
		dbURL,
		dbPort,
	)
	equals := "export SAJSON=$(cat /secrets/sa.json | base64 | tr -d '\n' ) && cockroach sql --certs-dir=/cockroach/cockroach-certs --host=testDB2 --port=81 -e \"RESTORE FROM \\\"gs://testBucket/testBackup/test2?AUTH=specified&CREDENTIALS=${SAJSON}\\\";\""
	assert.Equal(t, equals, cmd)
}
