package postgres

import (
	"testing"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// parseDSN is a test helper that parses a DSN and returns a Config
// with parsedDSN populated, mimicking what Decode() does.
func parseDSN(t *testing.T, dsn string) Config {
	t.Helper()
	parsed, err := pgxpool.ParseConfig(dsn)
	require.NoError(t, err)
	return Config{DSN: dsn, parsedDSN: parsed}
}

func TestConfig_parseDSN_parsesDSN(t *testing.T) {
	c := parseDSN(t, "postgresql://myuser:mypass@dbhost:5433/mydb?sslmode=require")
	require.NotNil(t, c.parsedDSN)
	assert.Equal(t, "dbhost", c.parsedDSN.ConnConfig.Host)
	assert.Equal(t, uint16(5433), c.parsedDSN.ConnConfig.Port)
	assert.Equal(t, "mydb", c.parsedDSN.ConnConfig.Database)
	assert.Equal(t, "myuser", c.parsedDSN.ConnConfig.User)
	assert.Equal(t, "mypass", c.parsedDSN.ConnConfig.Password)
}

func TestConfig_DatabaseName(t *testing.T) {
	t.Run("from DSN", func(t *testing.T) {
		c := parseDSN(t, "postgresql://u:p@h:5432/fromdsn?sslmode=disable")
		assert.Equal(t, "fromdsn", c.DatabaseName())
	})

	t.Run("from legacy field", func(t *testing.T) {
		c := Config{Database: "fromlegacy"}
		assert.Equal(t, "fromlegacy", c.DatabaseName())
	})
}

func TestConfig_Username(t *testing.T) {
	t.Run("from DSN", func(t *testing.T) {
		c := parseDSN(t, "postgresql://dsnuser:p@h:5432/db?sslmode=disable")
		assert.Equal(t, "dsnuser", c.Username())
	})

	t.Run("from legacy field", func(t *testing.T) {
		c := Config{User: User{Username: "legacyuser"}}
		assert.Equal(t, "legacyuser", c.Username())
	})
}

func TestConfig_Password(t *testing.T) {
	t.Run("from DSN", func(t *testing.T) {
		c := parseDSN(t, "postgresql://u:dsnpass@h:5432/db?sslmode=disable")
		assert.Equal(t, "dsnpass", c.Password())
	})

	t.Run("from legacy field", func(t *testing.T) {
		c := Config{User: User{Password: "legacypass"}}
		assert.Equal(t, "legacypass", c.Password())
	})
}

func TestConfig_Connect_DSN_with_useAdmin_errors(t *testing.T) {
	c := parseDSN(t, "postgresql://u:p@h:5432/db?sslmode=disable")
	_, _, err := c.Connect(true)
	assert.ErrorIs(t, err, ErrDSNWithAdminConnect)
}

func TestConfig_String(t *testing.T) {
	t.Run("returns DSN when set", func(t *testing.T) {
		dsn := "postgresql://u:p@h:5432/db?sslmode=disable"
		c := Config{DSN: dsn}
		assert.Equal(t, dsn, c.String(false))
	})

	t.Run("builds legacy string", func(t *testing.T) {
		c := Config{
			Host:     "myhost",
			Port:     5433,
			Database: "mydb",
			User: User{
				Username: "myuser",
				Password: "mypass",
				SSL:      SSL{Mode: "disable"},
			},
		}
		s := c.String(false)
		assert.Contains(t, s, "host=myhost")
		assert.Contains(t, s, "port=5433")
		assert.Contains(t, s, "user=myuser")
		assert.Contains(t, s, "password=mypass")
		assert.Contains(t, s, "dbname=mydb")
		assert.Contains(t, s, "sslmode=disable")
	})
}

func TestConfig_DSN_fullSpec(t *testing.T) {
	tests := []struct {
		name     string
		dsn      string
		wantHost string
		wantPort uint16
		wantDB   string
		wantUser string
		wantPass string
	}{
		{
			name:     "standard URL format",
			dsn:      "postgresql://myuser:mypass@myhost:5433/mydb?sslmode=disable",
			wantHost: "myhost",
			wantPort: 5433,
			wantDB:   "mydb",
			wantUser: "myuser",
			wantPass: "mypass",
		},
		{
			name:     "postgres:// scheme",
			dsn:      "postgres://u:p@h:5432/db?sslmode=disable",
			wantHost: "h",
			wantPort: 5432,
			wantDB:   "db",
			wantUser: "u",
			wantPass: "p",
		},
		{
			name:     "URL-encoded password",
			dsn:      "postgresql://user:p%40ss%23word@host:5432/db?sslmode=disable",
			wantHost: "host",
			wantPort: 5432,
			wantDB:   "db",
			wantUser: "user",
			wantPass: "p@ss#word",
		},
		{
			name:     "with SSL parameters",
			dsn:      "postgresql://user:pass@host:5432/db?sslmode=require",
			wantHost: "host",
			wantPort: 5432,
			wantDB:   "db",
			wantUser: "user",
			wantPass: "pass",
		},
		{
			name:     "with application_name",
			dsn:      "postgresql://user:pass@host:5432/db?sslmode=disable&application_name=myapp",
			wantHost: "host",
			wantPort: 5432,
			wantDB:   "db",
			wantUser: "user",
			wantPass: "pass",
		},
		{
			name:     "with options parameter",
			dsn:      "postgresql://user:pass@host:5432/db?sslmode=disable&options=-c%20search_path%3Dpublic",
			wantHost: "host",
			wantPort: 5432,
			wantDB:   "db",
			wantUser: "user",
			wantPass: "pass",
		},
		{
			name:     "key-value format",
			dsn:      "host=kvhost port=5434 user=kvuser password=kvpass dbname=kvdb sslmode=disable",
			wantHost: "kvhost",
			wantPort: 5434,
			wantDB:   "kvdb",
			wantUser: "kvuser",
			wantPass: "kvpass",
		},
		{
			name:     "default port (omitted)",
			dsn:      "postgresql://user:pass@host/db?sslmode=disable",
			wantHost: "host",
			wantPort: 5432,
			wantDB:   "db",
			wantUser: "user",
			wantPass: "pass",
		},
		{
			name:     "without password",
			dsn:      "postgresql://user@host:5432/db?sslmode=disable",
			wantHost: "host",
			wantPort: 5432,
			wantDB:   "db",
			wantUser: "user",
			wantPass: "",
		},
		{
			name:     "without database",
			dsn:      "postgresql://user:pass@host:5432?sslmode=disable",
			wantHost: "host",
			wantPort: 5432,
			wantDB:   "",
			wantUser: "user",
			wantPass: "pass",
		},
		{
			name:     "IPv6 host",
			dsn:      "postgresql://user:pass@[::1]:5432/db?sslmode=disable",
			wantHost: "::1",
			wantPort: 5432,
			wantDB:   "db",
			wantUser: "user",
			wantPass: "pass",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := parseDSN(t, tt.dsn)

			assert.Equal(t, tt.wantDB, c.DatabaseName())
			assert.Equal(t, tt.wantUser, c.Username())
			assert.Equal(t, tt.wantPass, c.Password())

			assert.Equal(t, tt.wantHost, c.parsedDSN.ConnConfig.Host)
			assert.Equal(t, tt.wantPort, c.parsedDSN.ConnConfig.Port)
		})
	}
}

func TestConfig_DSN_ignores_legacy_fields(t *testing.T) {
	// When DSN is set and parsedDSN is populated, accessors return
	// values from the DSN, not from the legacy fields.
	c := parseDSN(t, "postgresql://dsnuser:dsnpass@dsnhost:5433/dsndb?sslmode=disable")
	c.Host = "legacyhost"
	c.Port = 9999
	c.Database = "legacydb"
	c.User = User{Username: "legacyuser", Password: "legacypass"}

	assert.Equal(t, "dsndb", c.DatabaseName())
	assert.Equal(t, "dsnuser", c.Username())
	assert.Equal(t, "dsnpass", c.Password())
}
