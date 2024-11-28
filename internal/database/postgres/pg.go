package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/jackc/pgx/v5/stdlib"
	"github.com/mitchellh/mapstructure"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database/dialect"
)

func init() {
	config := new(Config)
	dialect.Register(config, config, false)
}

const (
	sslDisabledMode = "disable"
	sslRequireMode  = "require"
	sslAllowMode    = "allow"
	sslPreferMode   = "prefer"
)

type Config struct {
	Host               string
	Port               int32
	Database           string
	EventPushConnRatio float64
	MaxOpenConns       uint32
	MaxIdleConns       uint32
	MaxConnLifetime    time.Duration
	MaxConnIdleTime    time.Duration
	User               User
	Admin              AdminUser
	// Additional options to be appended as options=<Options>
	// The value will be taken as is. Multiple options are space separated.
	Options string
}

func (c *Config) MatchName(name string) bool {
	for _, key := range []string{"pg", "postgres"} {
		if strings.TrimSpace(strings.ToLower(name)) == key {
			return true
		}
	}
	return false
}

func (_ *Config) Decode(configs []interface{}) (dialect.Connector, error) {
	connector := new(Config)
	decoder, err := mapstructure.NewDecoder(&mapstructure.DecoderConfig{
		DecodeHook:       mapstructure.StringToTimeDurationHookFunc(),
		WeaklyTypedInput: true,
		Result:           connector,
	})
	if err != nil {
		return nil, err
	}

	for _, config := range configs {
		if err = decoder.Decode(config); err != nil {
			return nil, err
		}
	}

	return connector, nil
}

func (c *Config) Connect(useAdmin bool, pusherRatio, spoolerRatio float64, purpose dialect.DBPurpose) (*sql.DB, *pgxpool.Pool, error) {
	connConfig, err := dialect.NewConnectionConfig(c.MaxOpenConns, c.MaxIdleConns, pusherRatio, spoolerRatio, purpose)
	if err != nil {
		return nil, nil, err
	}

	config, err := pgxpool.ParseConfig(c.String(useAdmin, purpose.AppName()))
	if err != nil {
		return nil, nil, err
	}

	config.AfterConnect = func(ctx context.Context, conn *pgx.Conn) error {
		for _, f := range connConfig.AfterConnect {
			if err := f(ctx, conn); err != nil {
				return err
			}
		}
		return nil
	}

	// For the pusher we set the app name with the instance ID
	if purpose == dialect.DBPurposeEventPusher {
		config.BeforeAcquire = func(ctx context.Context, conn *pgx.Conn) bool {
			return setAppNameWithID(ctx, conn, purpose, authz.GetInstance(ctx).InstanceID())
		}
		config.AfterRelease = func(conn *pgx.Conn) bool {
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			return setAppNameWithID(ctx, conn, purpose, "IDLE")
		}
	}

	if connConfig.MaxOpenConns != 0 {
		config.MaxConns = int32(connConfig.MaxOpenConns)
	}

	config.MaxConnLifetime = c.MaxConnLifetime
	config.MaxConnIdleTime = c.MaxConnIdleTime

	pool, err := pgxpool.NewWithConfig(
		context.Background(),
		config,
	)
	if err != nil {
		return nil, nil, err
	}

	if err := pool.Ping(context.Background()); err != nil {
		return nil, nil, err
	}

	return stdlib.OpenDBFromPool(pool), pool, nil
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
	return "postgres"
}

func (c *Config) Timetravel(time.Duration) string {
	return ""
}

type User struct {
	Username string
	Password string
	SSL      SSL
}

type AdminUser struct {
	// ExistingDatabase is the database to connect to before the ZITADEL database exists
	ExistingDatabase string
	User             `mapstructure:",squash"`
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

func (s *Config) checkSSL(user User) {
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
		user = c.Admin.User
	}
	c.checkSSL(user)
	fields := []string{
		"host=" + c.Host,
		"port=" + strconv.Itoa(int(c.Port)),
		"user=" + user.Username,
		"application_name=" + appName,
		"sslmode=" + user.SSL.Mode,
	}
	if c.Options != "" {
		fields = append(fields, "options="+c.Options)
	}
	if user.Password != "" {
		fields = append(fields, "password="+user.Password)
	}
	if !useAdmin {
		fields = append(fields, "dbname="+c.Database)
	} else {
		defaultDB := c.Admin.ExistingDatabase
		if defaultDB == "" {
			defaultDB = "postgres"
		}
		fields = append(fields, "dbname="+defaultDB)
	}
	if user.SSL.Mode != sslDisabledMode {
		if user.SSL.RootCert != "" {
			fields = append(fields, "sslrootcert="+user.SSL.RootCert)
		}
		if user.SSL.Cert != "" {
			fields = append(fields, "sslcert="+user.SSL.Cert)
		}
		if user.SSL.Key != "" {
			fields = append(fields, "sslkey="+user.SSL.Key)
		}
	}

	return strings.Join(fields, " ")
}

func setAppNameWithID(ctx context.Context, conn *pgx.Conn, purpose dialect.DBPurpose, id string) bool {
	// needs to be set like this because psql complains about parameters in the SET statement
	query := fmt.Sprintf("SET application_name = '%s_%s'", purpose.AppName(), id)
	_, err := conn.Exec(ctx, query)
	logging.OnError(err).Warn("failed to set application name")
	return err == nil
}
