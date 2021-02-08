package backup

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBackup_Command1(t *testing.T) {
	timestamp := ""
	databases := []string{}
	bucketName := "test"
	backupName := "test"

	cmd := getBackupCommand(timestamp, databases, bucketName, backupName)
	equals := "export " + backupNameEnv + "=$(date +%Y-%m-%dT%H:%M:%SZ)"
	assert.Equal(t, equals, cmd)
}

func TestBackup_Command2(t *testing.T) {
	timestamp := "test"
	databases := []string{}
	bucketName := "test"
	backupName := "test"

	cmd := getBackupCommand(timestamp, databases, bucketName, backupName)
	equals := "export " + backupNameEnv + "=test"
	assert.Equal(t, equals, cmd)
}

func TestBackup_Command3(t *testing.T) {
	timestamp := ""
	databases := []string{"testDb"}
	bucketName := "testBucket"
	backupName := "testBackup"

	cmd := getBackupCommand(timestamp, databases, bucketName, backupName)
	equals := "export " + backupNameEnv + "=$(date +%Y-%m-%dT%H:%M:%SZ) && /scripts/backup.sh testBackup testBucket testDb " + backupPath + " " + secretPath + " " + certPath + " ${" + backupNameEnv + "}"
	assert.Equal(t, equals, cmd)
}

func TestBackup_Command4(t *testing.T) {
	timestamp := "test"
	databases := []string{"test1", "test2", "test3"}
	bucketName := "testBucket"
	backupName := "testBackup"

	cmd := getBackupCommand(timestamp, databases, bucketName, backupName)
	equals := "export " + backupNameEnv + "=test && " +
		"/scripts/backup.sh testBackup testBucket test1 " + backupPath + " " + secretPath + " " + certPath + " ${" + backupNameEnv + "} && " +
		"/scripts/backup.sh testBackup testBucket test2 " + backupPath + " " + secretPath + " " + certPath + " ${" + backupNameEnv + "} && " +
		"/scripts/backup.sh testBackup testBucket test3 " + backupPath + " " + secretPath + " " + certPath + " ${" + backupNameEnv + "}"
	assert.Equal(t, equals, cmd)
}
