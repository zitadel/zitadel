package types

import (
	"database/sql"
	"strings"

	"github.com/caos/zitadel/internal/errors"
)

type SQL struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLmode  string
}

func (s *SQL) ConnectionString() string {
	fields := []string{
		"host=" + s.Host,
		"port=" + s.Port,
		"user=" + s.User,
		"password=" + s.Password,
		"dbname=" + s.Database,
		"sslmode=" + s.SSLmode,
	}

	return strings.Join(fields, " ")
}

func (s *SQL) Start() (*sql.DB, error) {
	client, err := sql.Open("postgres", s.ConnectionString())
	if err != nil {
		return nil, errors.ThrowPreconditionFailed(err, "TYPES-9qBtr", "unable to open database connection")
	}
	return client, nil
}
