package types

import "strings"

type SQL struct {
	Host     string
	Port     string
	User     string
	Password string
	Database string
	SSLmode  string
}

func (sql *SQL) ConnectionString() string {
	fields := []string{
		"host=" + sql.Host,
		"port=" + sql.Port,
		"user=" + sql.User,
		"password=" + sql.Password,
		"dbname=" + sql.Database,
		"sslmode=" + sql.SSLmode,
	}

	return strings.Join(fields, " ")
}
