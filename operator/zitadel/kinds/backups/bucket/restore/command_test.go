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
	enpoint := "testEndpoint"
	prefix := "testPrefix"

	cmd := getCommand(
		timestamp,
		bucketName,
		backupName,
		certPath,
		saSecretPath,
		dbURL,
		dbPort,
		enpoint,
		prefix,
	)

	equals := "backupctl restore gcs --backupname=testBackup --backupnameenv=BACKUP_NAME --asset-endpoint=testEndpoint --asset-akid=$(cat /secrets/akid) --asset-sak=$(cat /secrets/sak) --host=testDB --port=80 --source-sajsonpath=/secrets/sa.json --source-buckettestBucket --certs-dir=/cockroach/cockroach-certs"
	assert.Equal(t, equals, cmd)
}

func TestBackup_Command2(t *testing.T) {
	timestamp := "test2"
	bucketName := "testBucket"
	backupName := "testBackup"
	dbURL := "testDB2"
	dbPort := int32(81)
	enpoint := "testEndpoint"
	prefix := "testPrefix"

	cmd := getCommand(
		timestamp,
		bucketName,
		backupName,
		certPath,
		saSecretPath,
		dbURL,
		dbPort,
		enpoint,
		prefix,
	)
	equals := "backupctl restore gcs --backupname=testBackup --backupnameenv=BACKUP_NAME --asset-endpoint=testEndpoint --asset-akid=$(cat /secrets/akid) --asset-sak=$(cat /secrets/sak) --host=testDB2 --port=81 --source-sajsonpath=/secrets/sa.json --source-buckettestBucket --certs-dir=/cockroach/cockroach-certs"
	assert.Equal(t, equals, cmd)
}
