package cockroach

import (
	"database/sql"
	"strconv"
	"strings"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/mitchellh/mapstructure"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database/dialect"
)

func init() {
	config := &Config{}
	dialect.Register(config, config, true)
}

const (
	sslDisabledMode = "disable"
	sslRequireMode  = "require"
	sslAllowMode    = "allow"
	sslPreferMode   = "prefer"
)

type Config struct {
	Host            string
	Port            uint16
	Database        string
	MaxOpenConns    uint32
	MaxIdleConns    uint32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	User            User
	Admin           User
	// Additional options to be appended as options=<Options>
	// The value will be taken as is. Multiple options are space separated.
	Options string
}

func (c *Config) MatchName(name string) bool {
	for _, key := range []string{"crdb", "cockroach"} {
		if strings.TrimSpace(strings.ToLower(name)) == key {
			return true
		}
	}
	return false
}

func (c *Config) Decode(configs []interface{}) (dialect.Connector, error) {
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
		WeaklyTypedInput: true,
		Result:           c,
	})
	if err != nil {
		return nil, err
	}

	for _, config := range configs {
		if err = decoder.Decode(config); err != nil {
			return nil, err
		}
	}

	return c, nil
}

func (c *Config) Connect(useAdmin bool, pusherRatio, spoolerRatio float64, purpose dialect.DBPurpose) (*sql.DB, error) {
	client, err := sql.Open("pgx", c.String(useAdmin, purpose.AppName()))
	if err != nil {
		return nil, err
	}

	connConfig, err := dialect.NewConnectionConfig(c.MaxOpenConns, c.MaxIdleConns, spoolerRatio, pusherRatio, purpose)
	if err != nil {
		return nil, err
	}

	client.SetMaxOpenConns(int(connConfig.MaxOpenConns))
	client.SetMaxIdleConns(int(connConfig.MaxIdleConns))
	client.SetConnMaxLifetime(c.MaxConnLifetime)
	client.SetConnMaxIdleTime(c.MaxConnIdleTime)

	return client, nil
}

func (c *Config) DatabaseName() string {
	return c.Database
}

func (c *Config) Username() string {
	return c.User.Username
}

func (c *Config) Password() string {
	return c.User.Password
}

func (c *Config) Type() string {
	return "cockroach"
}

func (c *Config) Timetravel(d time.Duration) string {
	return ""
}

type User struct {
	Username string
	Password string
	SSL      SSL
}

type SSL struct {
	// type of connection security
	Mode string
	// RootCert Path to the CA certificate
	RootCert string
	// Cert Path to the client certificate
	Cert string
	// Key Path to the client private key
	Key string
}

func (c *Config) checkSSL(user User) {
	if user.SSL.Mode == sslDisabledMode || user.SSL.Mode == "" {
		user.SSL = SSL{Mode: sslDisabledMode}
		return
	}

	if user.SSL.Mode == sslRequireMode || user.SSL.Mode == sslAllowMode || user.SSL.Mode == sslPreferMode {
		return
	}

	if user.SSL.RootCert == "" {
		logging.WithFields(
			"cert set", user.SSL.Cert != "",
			"key set", user.SSL.Key != "",
			"rootCert set", user.SSL.RootCert != "",
		).Fatal("at least ssl root cert has to be set")
	}
}

func (c Config) String(useAdmin bool, appName string) string {
	user := c.User
	if useAdmin {
		user = c.Admin
	}
	c.checkSSL(user)
	fields := []string{
		"host=" + c.Host,
		"port=" + strconv.Itoa(int(c.Port)),
		"user=" + user.Username,
		"dbname=" + c.Database,
		"application_name=" + appName,
		"sslmode=" + user.SSL.Mode,
	}
	if c.Options != "" {
		fields = append(fields, "options="+c.Options)
	}
	if !useAdmin {
		fields = append(fields, "dbname="+c.Database)
	}
	if user.Password != "" {
		fields = append(fields, "password="+user.Password)
	}
	if user.SSL.Mode != sslDisabledMode {
		fields = append(fields, "sslrootcert="+user.SSL.RootCert)
		if user.SSL.Cert != "" {
			fields = append(fields, "sslcert="+user.SSL.Cert)
		}
		if user.SSL.Key != "" {
			fields = append(fields, "sslkey="+user.SSL.Key)
		}
	}

	return strings.Join(fields, " ")
}
