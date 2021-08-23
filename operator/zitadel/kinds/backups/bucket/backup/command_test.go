package backup

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBackup_Command1(t *testing.T) {
	timestamp := ""
	bucketName := "test"
	backupName := "test"

	cmd := getBackupCommand(timestamp, bucketName, backupName)
	equals := "export BACKUP_NAME=$(date +%Y-%m-%dT%H:%M:%SZ) && rclone --no-check-certificate --config /secrets/rconfig sync minio:test bucket:test/backup--${test}"
	assert.Equal(t, equals, cmd)
}

func TestBackup_Command2(t *testing.T) {
	timestamp := "test"
	bucketName := "test"
	backupName := "test"

	cmd := getBackupCommand(timestamp, bucketName, backupName)
	equals := "export BACKUP_NAME=test && rclone --no-check-certificate --config /secrets/rconfig sync minio:test bucket:test/backup--${test}"
	assert.Equal(t, equals, cmd)
}

func TestBackup_Command3(t *testing.T) {
	timestamp := ""
	bucketName := "testBucket2"
	backupName := "testBackup2"

	cmd := getBackupCommand(timestamp, bucketName, backupName)
	equals := "export BACKUP_NAME=$(date +%Y-%m-%dT%H:%M:%SZ) && rclone --no-check-certificate --config /secrets/rconfig sync minio:testBucket2 bucket:testBucket2/backup--${testBackup2}"
	assert.Equal(t, equals, cmd)
}

func TestBackup_Command4(t *testing.T) {
	timestamp := "test2"
	bucketName := "testBucket2"
	backupName := "testBackup2"

	cmd := getBackupCommand(timestamp, bucketName, backupName)
	equals := "export BACKUP_NAME=test2 && rclone --no-check-certificate --config /secrets/rconfig sync minio:testBucket2 bucket:testBucket2/backup--${testBackup2}"
	assert.Equal(t, equals, cmd)
}
