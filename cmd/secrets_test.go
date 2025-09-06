package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadSecretsFromFiles(t *testing.T) {
	tempDir, err := os.MkdirTemp("", "zitadel_secrets_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}
	defer os.RemoveAll(tempDir)

	dbPasswordFile := filepath.Join(tempDir, "db_password")
	err = os.WriteFile(dbPasswordFile, []byte("supersecret123\n"), 0600)
	require.NoError(t, err)

	jwtKeyFile := filepath.Join(tempDir, "jwt_key")
	err = os.WriteFile(jwtKeyFile, []byte("jwt-secret-key"), 0600)
	require.NoError(t, err)

	t.Setenv("ZITADEL_DATABASE_COCKROACH_USER_PASSWORD_FILE", dbPasswordFile)
	t.Setenv("ZITADEL_INTERNAL_AUTHZ_REPOSITORY_JWT_KEY_FILE", jwtKeyFile)

	v := viper.New()

	LoadSecretsFromFiles(v)

	assert.Equal(t, "supersecret123", v.GetString("database.cockroach.user.password"))
	assert.Equal(t, "jwt-secret-key", v.GetString("internal.authz.repository.jwt.key"))
}

func TestLoadSecretsFromFiles_FileNotFound(t *testing.T) {
	t.Setenv("ZITADEL_DATABASE_PASSWORD_FILE", "/nonexistent/file")

	v := viper.New()

	LoadSecretsFromFiles(v)

	assert.Equal(t, "", v.GetString("database.password"))
}