package repository

import (
	"fmt"
	"strings"

	es_stor "github.com/caos/eventstore-lib/pkg/storage"
)

type Config struct {
	Dialect string
	Address address
}

type address struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
	Sslmode  string
}

func (c *Config) New() es_stor.Storage {
	sql := new(SQL)
	sql.dialect = c.Dialect
	sql.address = c.Address.connectionString()
	return sql
}

func (a *address) connectionString() string {
	fields := make([]string, 0, 6)
	fields = append(fields, joinSQLField("host", a.Host),
		joinSQLField("port", a.Port),
		joinSQLField("user", a.User),
		joinSQLField("password", a.Password),
		joinSQLField("dbname", a.Dbname),
		joinSQLField("sslmode", a.Sslmode))

	return strings.Join(fields, " ")
}

func joinSQLField(key, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}
