package repository

import (
	"context"
	"errors"

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

func scan(ctx context.Context, querier database.Querier, builder *database.StatementBuilder, res any) error {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return err
	}

	if err := rows.(database.CollectableRows).CollectExactlyOneRow(res); err != nil {
		if err.Error() == "no rows in result set" {
			return ErrResourceDoesNotExist
		}
		return err
	}
	return nil
}

func scanMultiple(ctx context.Context, querier database.Querier, builder *database.StatementBuilder, res any) error {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return err
	}

	if err := rows.(database.CollectableRows).Collect(res); err != nil {
		// if no results returned, this is not a error
		// it just means the organization was not found
		// the caller should check if the returned organization is nil
		if err.Error() == "no rows in result set" {
			return nil
		}
		return err
	}
	return nil
}
