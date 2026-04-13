package repository

import (
	"context"

	"github.com/zitadel/zitadel/backend/v3/domain"
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

// checkPKCondition checks if the Primary Key columns are part of the condition.
// This can ensure only a single row is affected by updates and deletes.
func checkPKCondition(
	repo domain.Repository,
	condition database.Condition,
) error {
	return checkRestrictingColumns(
		condition,
		repo.PrimaryKeyColumns()...,
	)
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

type updatable interface {
	PrimaryKeyColumns() []database.Column
	UpdatedAtColumn() database.Column
	qualifiedTableName() string
}

func updateOne[Target updatable](ctx context.Context, client database.QueryExecutor, target Target, condition database.Condition, changes ...database.Change) (int64, error) {
	if len(changes) == 0 {
		return 0, database.ErrNoChanges
	}
	if err := checkPKCondition(target, condition); err != nil {
		return 0, err
	}
	if !database.Changes(changes).IsOnColumn(target.UpdatedAtColumn()) {
		changes = append(changes, database.NewChange(target.UpdatedAtColumn(), database.NullInstruction))
	}
	builder := database.NewStatementBuilder("UPDATE ")
	builder.WriteString(target.qualifiedTableName())
	builder.WriteString(" SET ")
	if err := database.Changes(changes).Write(builder); err != nil {
		return 0, err
	}
	writeCondition(builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}

type deletable interface {
	PrimaryKeyColumns() []database.Column
	qualifiedTableName() string
}

func deleteOne[Target deletable](ctx context.Context, client database.QueryExecutor, target Target, condition database.Condition) (int64, error) {
	if err := checkPKCondition(target, condition); err != nil {
		return 0, err
	}

	builder := database.NewStatementBuilder("DELETE FROM ")
	builder.WriteString(target.qualifiedTableName())
	builder.WriteRune(' ')
	writeCondition(builder, condition)

	return client.Exec(ctx, builder.String(), builder.Args()...)
}
