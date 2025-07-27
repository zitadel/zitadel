package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

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

func scanRow(ctx context.Context, querier database.Querier, builder *database.StatementBuilder, res any) error {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return err
	}

	return rows.(database.CollectableRows).CollectExactlyOneRow(res)
}

func scanRows(ctx context.Context, querier database.Querier, builder *database.StatementBuilder, res any) error {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return err
	}

	return rows.(database.CollectableRows).Collect(res)
}
