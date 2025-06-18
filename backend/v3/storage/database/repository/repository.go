package repository

import "github.com/zitadel/zitadel/backend/v3/storage/database"

type repository struct {
	client database.QueryExecutor
}

func (r repository) writeCondition(
	builder *database.StatementBuilder,
	condition database.Condition,
) {
	if condition == nil {
		return
	}
	builder.WriteString(" WHERE ")
	condition.Write(builder)
}
