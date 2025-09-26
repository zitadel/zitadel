package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

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

func checkRestrictingColumns(
	condition database.Condition,
	requiredColumns ...database.Column,
) error {
	for _, col := range requiredColumns {
		if !condition.IsRestrictingColumn(col) {
			return database.NewMissingConditionError(col)
		}
	}
	return nil
}

func getOne[Target any](ctx context.Context, querier database.Querier, builder *database.StatementBuilder) (*Target, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	var target Target
	if err := rows.(database.CollectableRows).CollectExactlyOneRow(&target); err != nil {
		return nil, err
	}
	return &target, nil
}

func getMany[Target any](ctx context.Context, querier database.Querier, builder *database.StatementBuilder) ([]*Target, error) {
	rows, err := querier.Query(ctx, builder.String(), builder.Args()...)
	if err != nil {
		return nil, err
	}
	var targets []*Target
	if err := rows.(database.CollectableRows).Collect(&targets); err != nil {
		return nil, err
	}
	return targets, nil
}
