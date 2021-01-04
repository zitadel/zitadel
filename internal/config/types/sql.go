package types

import (
	"database/sql"
	"strings"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
)

const (
	sslDisabledMode = "disable"
)

type SQL struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSL      *ssl
}

type ssl struct {
	// type of connection security
	Mode string
	// RootCert Path to the CA certificate
	RootCert string
	// Cert Path to the client certificate
	Cert string
	// Key Path to the client private key
	Key string
}

func (s *SQL) connectionString() string {
	fields := []string{
		"host=" + s.Host,
		"port=" + s.Port,
		"user=" + s.User,
		"dbname=" + s.Database,
		"application_name=zitadel",
		"sslmode=" + s.SSL.Mode,
	}
	if s.Password != "" {
		fields = append(fields, "password="+s.Password)
	}

	if s.SSL.Mode != sslDisabledMode {
		fields = append(fields, []string{
			"sslrootcert=" + s.SSL.RootCert,
			"sslcert=" + s.SSL.Cert,
			"sslkey=" + s.SSL.Key,
		}...)
	}

	return strings.Join(fields, " ")
}

func (s *SQL) Start() (*sql.DB, error) {
	s.checkSSL()
	client, err := sql.Open("postgres", s.connectionString())
	if err != nil {
		return nil, errors.ThrowPreconditionFailed(err, "TYPES-9qBtr", "unable to open database connection")
	}
	// as we open many sql clients we set the max
	// open cons deep. now 3(maxconn) * 8(clients) = max 24 conns per pod
	client.SetMaxOpenConns(3)
	client.SetMaxIdleConns(3)
	return client, nil
}

func (s *SQL) checkSSL() {
	if s.SSL == nil || s.SSL.Mode == sslDisabledMode || s.SSL.Mode == "" {
		s.SSL = &ssl{Mode: sslDisabledMode}
		return
	}
	if s.SSL.Cert == "" || s.SSL.Key == "" || s.SSL.RootCert == "" {
		logging.LogWithFields("TYPES-LFdzP",
			"cert set", s.SSL.Cert != "",
			"key set", s.SSL.Key != "",
			"rootCert set", s.SSL.RootCert != "",
		).Fatal("fields for secure connection missing")
	}
}
