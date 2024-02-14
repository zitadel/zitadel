package handler

import (
	"database/sql"
)

type Executer interface {
	Exec(string, ...interface{}) (sql.Result, error)
}
