package restore

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestBackup_Command1(t *testing.T) {
	timestamp := ""
	databases := []string{}
	bucketName := "testBucket"
	backupName := "testBackup"

	cmd := getCommand(timestamp, databases, bucketName, backupName)
	equals := ""
	assert.Equal(t, equals, cmd)
}

func TestBackup_Command2(t *testing.T) {
	timestamp := ""
	databases := []string{"testDb"}
	bucketName := "testBucket"
	backupName := "testBackup"

	cmd := getCommand(timestamp, databases, bucketName, backupName)
	equals := "/scripts/restore.sh testBucket testBackup  testDb /secrets/sa.json /cockroach/cockroach-certs"
	assert.Equal(t, equals, cmd)
}

func TestBackup_Command3(t *testing.T) {
	timestamp := "test"
	databases := []string{"testDb"}
	bucketName := "testBucket"
	backupName := "testBackup"

	cmd := getCommand(timestamp, databases, bucketName, backupName)
	equals := "/scripts/restore.sh testBucket testBackup test testDb /secrets/sa.json /cockroach/cockroach-certs"
	assert.Equal(t, equals, cmd)
}

func TestBackup_Command4(t *testing.T) {
	timestamp := ""
	databases := []string{}
	bucketName := "test"
	backupName := "test"

	cmd := getCommand(timestamp, databases, bucketName, backupName)
	equals := ""
	assert.Equal(t, equals, cmd)
}
