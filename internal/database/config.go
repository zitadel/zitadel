package database

import (
	"strings"
	"time"

	"github.com/caos/logging"
)

const (
	sslDisabledMode = "disable"
)

type Config struct {
	Host            string
	Port            string
	Database        string
	MaxOpenConns    uint32
	MaxConnLifetime time.Duration
	MaxConnIdleTime time.Duration

	//Additional options to be appended as options=<Options>
	//The value will be taken as is. Multiple options are space separated.
	Options string
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

func (s *Config) checkSSL() {
	if s.SSL.Mode == sslDisabledMode || s.SSL.Mode == "" {
		s.SSL = SSL{Mode: sslDisabledMode}
		return
	}
	if s.SSL.RootCert == "" {
		logging.WithFields(
			"cert set", s.SSL.Cert != "",
			"key set", s.SSL.Key != "",
			"rootCert set", s.SSL.RootCert != "",
		).Fatal("at least ssl root cert has to be set")
	}
}

func (c Config) String() string {
	c.checkSSL()
	fields := []string{
		"host=" + c.Host,
		"port=" + c.Port,
		"user=" + c.Username,
		"dbname=" + c.Database,
		"application_name=zitadel",
		"sslmode=" + c.SSL.Mode,
	}
	if c.Options != "" {
		fields = append(fields, "options="+c.Options)
	}
	if c.Password != "" {
		fields = append(fields, "password="+c.Password)
	}
	if c.SSL.Mode != sslDisabledMode {
		fields = append(fields, "sslrootcert="+c.SSL.RootCert)
		if c.SSL.Cert != "" {
			fields = append(fields, "sslcert="+c.SSL.Cert)
		}
		if c.SSL.Key != "" {
			fields = append(fields, "sslkey="+c.SSL.Key)
		}
	}

	return strings.Join(fields, " ")
}
