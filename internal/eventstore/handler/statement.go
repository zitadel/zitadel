package handler

import (
	"context"
	"database/sql"
)

type Executer interface {
	Exec(string, ...interface{}) (sql.Result, error)
	ExecContext(context.Context, string, ...interface{}) (sql.Result, error)
}
