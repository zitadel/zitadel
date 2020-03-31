package view

import (
	"fmt"
	"strings"
)

type ViewConfig struct {
	Address *address
}

type address struct {
	Host     string
	Port     string
	User     string
	Password string
	Dbname   string
	Sslmode  string
}

func (c *ViewConfig) GetConnectionString() string {
	return c.Address.join()
}

func (a *address) join() string {
	fields := make([]string, 0, 6)
	fields = append(fields, joinSQLField("host", a.Host))
	fields = append(fields, joinSQLField("port", a.Port))
	fields = append(fields, joinSQLField("user", a.User))
	fields = append(fields, joinSQLField("password", a.Password))
	fields = append(fields, joinSQLField("dbname", a.Dbname))
	fields = append(fields, joinSQLField("sslmode", a.Sslmode))

	return strings.Join(fields, " ")
}

func joinSQLField(key, value string) string {
	return fmt.Sprintf("%s=%s", key, value)
}
