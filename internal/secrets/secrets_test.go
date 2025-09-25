package secrets

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestProcessDockerSecretsIntoViper(t *testing.T) {
	tmpDir := t.TempDir()
	
	dbPasswordFile := filepath.Join(tmpDir, "db_password")
	err := os.WriteFile(dbPasswordFile, []byte("supersecret123\n"), 0600)
	require.NoError(t, err)
	
	jwtKeyFile := filepath.Join(tmpDir, "jwt_key")
	err = os.WriteFile(jwtKeyFile, []byte("jwt-secret-key"), 0600)
	require.NoError(t, err)
	
	t.Setenv("ZITADEL_DATABASE_COCKROACH_USER_PASSWORD_FILE", dbPasswordFile)
	t.Setenv("ZITADEL_INTERNAL_AUTHZ_REPOSITORY_JWT_KEY_FILE", jwtKeyFile)
	t.Setenv("REGULAR_ENV_VAR", "regular_value")
	t.Setenv("EMPTY_FILE_VAR_FILE", "")
	
	v := viper.New()
	
	err = ProcessDockerSecretsIntoViper(v)
	require.NoError(t, err)
	
	assert.Equal(t, "supersecret123", v.GetString("database.cockroach.user.password"))
	assert.Equal(t, "jwt-secret-key", v.GetString("internal.authz.repository.jwt.key"))
	
	assert.Equal(t, "", v.GetString("regular.env.var"))
	assert.Equal(t, "", v.GetString("empty.file.var"))
}

func TestProcessDockerSecretsIntoViper_FileNotFound(t *testing.T) {
	t.Setenv("ZITADEL_DATABASE_PASSWORD_FILE", "/nonexistent/file")
	
	v := viper.New()
	
	err := ProcessDockerSecretsIntoViper(v)
	require.NoError(t, err)
	
	assert.Equal(t, "", v.GetString("database.password"))
}

func TestProcessDockerSecretsIntoViper_NoFileVars(t *testing.T) {
	t.Setenv("ZITADEL_DATABASE_PASSWORD", "regular_password")
	t.Setenv("NORMAL_VAR", "normal_value")
	
	v := viper.New()
	
	err := ProcessDockerSecretsIntoViper(v)
	require.NoError(t, err)
	
	assert.Equal(t, "", v.GetString("database.password"))
}
