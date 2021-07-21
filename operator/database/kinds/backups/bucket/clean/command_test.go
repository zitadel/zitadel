package clean

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestClean_Command1(t *testing.T) {
	databases := []string{}
	users := []string{}

	cmd := getCommand(databases, users)
	equals := "/scripts/clean-migration.sh " + certPath

	assert.Equal(t, equals, cmd)
}

func TestClean_Command2(t *testing.T) {
	databases := []string{"test"}
	users := []string{"test"}

	cmd := getCommand(databases, users)
	equals := "/scripts/clean-db.sh " + certPath + " test && /scripts/clean-user.sh " + certPath + " test && /scripts/clean-migration.sh " + certPath

	assert.Equal(t, equals, cmd)
}

func TestClean_Command3(t *testing.T) {
	databases := []string{"test1", "test2", "test3"}
	users := []string{"test1", "test2", "test3"}

	cmd := getCommand(databases, users)
	equals := "/scripts/clean-db.sh " + certPath + " test1 && /scripts/clean-db.sh " + certPath + " test2 && /scripts/clean-db.sh " + certPath + " test3 && " +
		"/scripts/clean-user.sh " + certPath + " test1 && /scripts/clean-user.sh " + certPath + " test2 && /scripts/clean-user.sh " + certPath + " test3 && " +
		"/scripts/clean-migration.sh " + certPath

	assert.Equal(t, equals, cmd)
}
