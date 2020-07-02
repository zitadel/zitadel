package types

import (
	"database/sql"
	"strings"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
)

const (
	sslDisabledMode = "disabled"
)

type SQL struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	// type of connection security
	SSLMode string
	// RootCert Path to the CA certificate
	SSLRootCert string
	// Cert Path to the client certificate
	SSLCert string
	// Key Path to the client private key
	SSLKey string
}

func (s *SQL) ConnectionString() string {
	fields := []string{
		"host=" + s.Host,
		"port=" + s.Port,
		"user=" + s.User,
		"password=" + s.Password,
		"dbname=" + s.Database,
		"sslmode=" + s.SSLMode,
	}
	if s.SSLMode != sslDisabledMode {
		fields = append(fields, []string{
			"ssl=true",
			"sslrootcert=" + s.SSLRootCert,
			"sslcert=" + s.SSLCert,
			"sslkey=" + s.SSLKey,
		}...)
	}

	return strings.Join(fields, " ")
}

func (s *SQL) Start() (*sql.DB, error) {
	s.checkSSL()
	client, err := sql.Open("postgres", s.ConnectionString())
	if err != nil {
		return nil, errors.ThrowPreconditionFailed(err, "TYPES-9qBtr", "unable to open database connection")
	}
	return client, nil
}

func (s *SQL) checkSSL() {
	if s.SSLMode != sslDisabledMode && (s.SSLCert == "" || s.SSLKey == "" || s.SSLRootCert == "") {
		logging.LogWithFields("TYPES-LFdzP", "mode",
			s.SSLMode, "cert",
			s.SSLCert, "key",
			s.SSLKey, "rootCert",
			s.SSLRootCert).
			Fatal("wrong SSL config for database")
	}
}
