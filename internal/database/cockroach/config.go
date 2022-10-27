package cockroach

import (
	"database/sql"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mitchellh/mapstructure"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database/dialect"
)

const (
	sslDisabledMode = "disable"
)

type Config struct {
	Host            string
	Port            uint16
	Database        string
	MaxOpenConns    uint32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration
	User            User
	Admin           User

	//Additional options to be appended as options=<Options>
	//The value will be taken as is. Multiple options are space separated.
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

func (c *Config) Connect(useAdmin bool) (*sql.DB, error) {
	client, err := sql.Open("pgx", c.String(useAdmin))
	if err != nil {
		return nil, err
	}

	client.SetMaxOpenConns(int(c.MaxOpenConns))
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
	// RootCertValue is the CA certificate in plain text
	RootCertValue     string
	rootCertValueFile string
	// CertValue is the client certificate in plain text
	CertValue     string
	certValueFile string
	// KeyValue is the client private key in plain text
	KeyValue     string
	keyValueFile string
}

func (s *SSL) ensureFiles() {
	s.rootCertValueFile = ensureFile(s.RootCert, s.RootCertValue, "./zitadel-generated-from-config-root.crt")
	s.certValueFile = ensureFile(s.Cert, s.CertValue, "./zitadel-generated-from-config.crt")
	s.keyValueFile = ensureFile(s.Key, s.KeyValue, "./zitadel-generated-from-config.key")
}

func ensureFile(pathConfig, valueConfig, createFileForValue string) string {
	if pathConfig != "" {
		return ""
	}
	file, err := os.OpenFile(createFileForValue, os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0400)
	logging.OnError(err).Fatalf("creating file %s for certificate failed", createFileForValue)

	_, err = io.Copy(file, strings.NewReader(valueConfig))
	logging.OnError(err).Fatalf("copying certificate to file %s failed", file.Name())
	return file.Name()
}

func (c *Config) checkSSL(user User) {
	if user.SSL.Mode == sslDisabledMode || user.SSL.Mode == "" {
		user.SSL = SSL{Mode: sslDisabledMode}
		return
	}
	if user.SSL.RootCert == "" && user.SSL.RootCertValue == "" {
		logging.WithFields(
			"cert path set", user.SSL.Cert != "",
			"cert value set", user.SSL.CertValue != "",
			"key path set", user.SSL.Key != "",
			"key value set", user.SSL.KeyValue != "",
			"rootCert path set", user.SSL.RootCert != "",
			"rootCert value set", user.SSL.RootCertValue != "",
		).Fatal("at least ssl root cert has to be set")
	}
}

func (c Config) String(useAdmin bool) string {
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
		"application_name=zitadel",
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
		rootCertPath := user.SSL.rootCertValueFile
		if rootCertPath == "" {
			rootCertPath = user.SSL.RootCert
		}
		fields = append(fields, "sslrootcert="+rootCertPath)
		certPath := user.SSL.certValueFile
		if certPath == "" {
			certPath = user.SSL.Cert
		}
		if user.SSL.Cert != "" {
			fields = append(fields, "sslcert="+certPath)
		}
		keyPath := user.SSL.keyValueFile
		if keyPath == "" {
			keyPath = user.SSL.Key
		}
		if user.SSL.Key != "" {
			fields = append(fields, "sslkey="+keyPath)
		}
	}

	return strings.Join(fields, " ")
}
