package migration

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMigration_CreateUserCommand(t *testing.T) {
	user := "test"
	pw := "test"
	file := "test"
	equals := "echo -n 'CREATE USER IF NOT EXISTS ' > " + file + ";echo -n ${" + user + "} >> " + file + ";echo -n ';' >> " + file + ";echo -n 'ALTER USER ' >> " + file + ";echo -n ${" + user + "} >> " + file + ";echo -n ' WITH PASSWORD ' >> " + file + ";echo -n ${" + pw + "} >> " + file + ";echo -n ';' >> " + file

	cmd := createUserCommand(user, pw, file)
	assert.Equal(t, cmd, equals)

	user = "test1"
	pw = "test1"
	file = "test1"
	equals = "echo -n 'CREATE USER IF NOT EXISTS ' > " + file + ";echo -n ${" + user + "} >> " + file + ";echo -n ';' >> " + file + ";echo -n 'ALTER USER ' >> " + file + ";echo -n ${" + user + "} >> " + file + ";echo -n ' WITH PASSWORD ' >> " + file + ";echo -n ${" + pw + "} >> " + file + ";echo -n ';' >> " + file

	cmd = createUserCommand(user, pw, file)
	assert.Equal(t, cmd, equals)

	user = "test2"
	pw = "test2"
	file = "test2"
	equals = "echo -n 'CREATE USER IF NOT EXISTS ' > " + file + ";echo -n ${" + user + "} >> " + file + ";echo -n ';' >> " + file + ";echo -n 'ALTER USER ' >> " + file + ";echo -n ${" + user + "} >> " + file + ";echo -n ' WITH PASSWORD ' >> " + file + ";echo -n ${" + pw + "} >> " + file + ";echo -n ';' >> " + file

	cmd = createUserCommand(user, pw, file)
	assert.Equal(t, cmd, equals)

	user = "test"
	pw = ""
	file = "test"
	equals = "echo -n 'CREATE USER IF NOT EXISTS ' > " + file + ";echo -n ${" + user + "} >> " + file + ";echo -n ';' >> " + file

	cmd = createUserCommand(user, pw, file)
	assert.Equal(t, cmd, equals)

	user = "test2"
	pw = ""
	file = "test2"
	equals = "echo -n 'CREATE USER IF NOT EXISTS ' > " + file + ";echo -n ${" + user + "} >> " + file + ";echo -n ';' >> " + file

	cmd = createUserCommand(user, pw, file)
	assert.Equal(t, cmd, equals)

	user = "test2"
	pw = ""
	file = ""

	cmd = createUserCommand(user, pw, file)
	assert.Equal(t, cmd, "")

	user = ""
	pw = ""
	file = "test"

	cmd = createUserCommand(user, pw, file)
	assert.Equal(t, cmd, "")

}

func TestMigration_GrantUserCommand(t *testing.T) {
	user := "test"
	file := "test"
	equals := strings.Join([]string{
		"echo -n 'GRANT admin TO ' > " + file,
		"echo -n ${" + user + "} >> " + file,
		"echo -n ' WITH ADMIN OPTION;'  >> " + file,
	}, ";")

	cmd := grantUserCommand(user, file)
	assert.Equal(t, cmd, equals)
}

func TestMigration_Image_Default(t *testing.T) {

}
