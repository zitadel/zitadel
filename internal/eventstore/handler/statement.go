package handler

import (
	"context"
)

type Executer interface {
	Exec(context.Context, string, ...interface{}) (int64, error)
}
