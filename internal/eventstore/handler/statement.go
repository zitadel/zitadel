package handler

import (
	"database/sql"
)

type Executer interface {
	Exec(string, ...any) (sql.Result, error)
}
