package repository

import (
	"errors"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

var ErrResourceDoesNotExist = errors.New("resource does not exist")

type repository struct {
	client database.QueryExecutor
}

func writeCondition(
	builder *database.StatementBuilder,
	condition database.Condition,
) {
	if condition == nil {
		return
	}
	builder.WriteString(" WHERE ")
	condition.Write(builder)
}
