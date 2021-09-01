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
	enpoint := "testEndpoint"
	prefix := "testPrefix"

	cmd := getBackupCommand(
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
	equals := "export " + backupNameEnv + "=$(date +%Y-%m-%dT%H:%M:%SZ) && backupctl backup gcs --backupname=test --backupnameenv=BACKUP_NAME --asset-endpoint=testEndpoint --asset-akid=$(cat /secrets/akid) --asset-sak=$(cat /secrets/sak) --asset-prefix=testPrefix --host=testDB --port=80 --destination-sajsonpath=/secrets/sa.json --destination-buckettest --certs-dir=/cockroach/cockroach-certs"
	assert.Equal(t, equals, cmd)
}

func TestBackup_Command2(t *testing.T) {
	timestamp := "test"
	bucketName := "test"
	backupName := "test"
	dbURL := "testDB"
	dbPort := int32(80)
	enpoint := "testEndpoint"
	prefix := "testPrefix"

	cmd := getBackupCommand(
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
	equals := "export " + backupNameEnv + "=test && backupctl backup gcs --backupname=test --backupnameenv=BACKUP_NAME --asset-endpoint=testEndpoint --asset-akid=$(cat /secrets/akid) --asset-sak=$(cat /secrets/sak) --asset-prefix=testPrefix --host=testDB --port=80 --destination-sajsonpath=/secrets/sa.json --destination-buckettest --certs-dir=/cockroach/cockroach-certs"
	assert.Equal(t, equals, cmd)
}
