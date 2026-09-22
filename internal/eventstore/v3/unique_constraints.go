package eventstore

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	//go:embed unique_constraints_delete.sql
	deleteConstraintStmt string
	//go:embed unique_constraints_delete_placeholders.sql
	deleteConstraintPlaceholdersStmt string
	//go:embed unique_constraints_add.sql
	addConstraintStmt string
	//go:embed unique_constraints_add_without_owners.sql
	addConstraintWithoutOwnersStmt string
)

func handleUniqueConstraints(ctx context.Context, tx database.Tx, commands []eventstore.Command) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	deletePlaceholders := make([]string, 0)
	deleteArgs := make([]any, 0)

	addWithOwnersPlaceholders := make([]string, 0)
	addWithOwnersArgs := make([]any, 0)
	addWithoutOwnersPlaceholders := make([]string, 0)
	addWithoutOwnersArgs := make([]any, 0)
	addConstraints := map[string]*eventstore.UniqueConstraint{}
	deleteConstraints := map[string]*eventstore.UniqueConstraint{}

	for _, command := range commands {
		for _, constraint := range command.UniqueConstraints() {
			instanceID := command.Aggregate().InstanceID
			if constraint.IsGlobal {
				instanceID = ""
			}
			switch constraint.Action {
			case eventstore.UniqueConstraintAdd:
				constraint.UniqueField = strings.ToLower(constraint.UniqueField)
				addConstraints[fmt.Sprintf(uniqueConstraintPlaceholderFmt, instanceID, constraint.UniqueType, constraint.UniqueField)] = constraint
				if len(constraint.Owners) == 0 {
					addWithoutOwnersPlaceholders = append(addWithoutOwnersPlaceholders, fmt.Sprintf("($%d, $%d, $%d)", len(addWithoutOwnersArgs)+1, len(addWithoutOwnersArgs)+2, len(addWithoutOwnersArgs)+3))
					addWithoutOwnersArgs = append(addWithoutOwnersArgs, instanceID, constraint.UniqueType, constraint.UniqueField)
				} else {
					addWithOwnersPlaceholders = append(addWithOwnersPlaceholders, fmt.Sprintf("($%d, $%d, $%d, $%d)", len(addWithOwnersArgs)+1, len(addWithOwnersArgs)+2, len(addWithOwnersArgs)+3, len(addWithOwnersArgs)+4))
					addWithOwnersArgs = append(addWithOwnersArgs, instanceID, constraint.UniqueType, constraint.UniqueField, constraint.Owners)
				}
			case eventstore.UniqueConstraintRemove:
				deletePlaceholders = append(deletePlaceholders, fmt.Sprintf(deleteConstraintPlaceholdersStmt, len(deleteArgs)+1, len(deleteArgs)+2, len(deleteArgs)+3))
				deleteArgs = append(deleteArgs, instanceID, constraint.UniqueType, constraint.UniqueField)
				deleteConstraints[fmt.Sprintf(uniqueConstraintPlaceholderFmt, instanceID, constraint.UniqueType, constraint.UniqueField)] = constraint
			case eventstore.UniqueConstraintInstanceRemove:
				deletePlaceholders = append(deletePlaceholders, fmt.Sprintf("(instance_id = $%d)", len(deleteArgs)+1))
				deleteArgs = append(deleteArgs, instanceID)
				deleteConstraints[fmt.Sprintf(uniqueConstraintPlaceholderFmt, instanceID, constraint.UniqueType, constraint.UniqueField)] = constraint
			case eventstore.UniqueConstraintRemoveByOwner:
				if len(constraint.Owners) == 0 {
					continue
				}
				tag := constraint.Owners[0]
				deletePlaceholders = append(deletePlaceholders, fmt.Sprintf("(instance_id = $%d AND owners @> ARRAY[$%d]::text[])", len(deleteArgs)+1, len(deleteArgs)+2))
				deleteArgs = append(deleteArgs, instanceID, tag)
				deleteConstraints[fmt.Sprintf("%s:%s", instanceID, tag)] = constraint
			}
		}
	}

	if len(deletePlaceholders) > 0 {
		_, err := tx.ExecContext(ctx, fmt.Sprintf(deleteConstraintStmt, strings.Join(deletePlaceholders, " OR ")), deleteArgs...)
		if err != nil {
			logging.WithError(err).Warn("delete unique constraint failed")
			errMessage := "Errors.Internal"
			if constraint := constraintFromErr(err, deleteConstraints); constraint != nil {
				errMessage = constraint.ErrorMessage
			}
			return zerrors.ThrowInternal(err, "V3-C8l3V", errMessage)
		}
	}
	if err := execAddConstraints(ctx, tx, addConstraintWithoutOwnersStmt, addWithoutOwnersPlaceholders, addWithoutOwnersArgs, addConstraints); err != nil {
		return err
	}
	return execAddConstraints(ctx, tx, addConstraintStmt, addWithOwnersPlaceholders, addWithOwnersArgs, addConstraints)
}

func execAddConstraints(ctx context.Context, tx database.Tx, stmt string, placeholders []string, args []any, addConstraints map[string]*eventstore.UniqueConstraint) error {
	if len(placeholders) == 0 {
		return nil
	}
	_, err := tx.ExecContext(ctx, fmt.Sprintf(stmt, strings.Join(placeholders, ", ")), args...)
	if err != nil {
		logging.WithError(err).Warn("add unique constraint failed")
		errMessage := "Errors.Internal"
		if constraint := constraintFromErr(err, addConstraints); constraint != nil {
			errMessage = constraint.ErrorMessage
		}
		return zerrors.ThrowAlreadyExists(err, "V3-DKcYh", errMessage)
	}
	return nil
}

func constraintFromErr(err error, constraints map[string]*eventstore.UniqueConstraint) *eventstore.UniqueConstraint {
	pgErr := new(pgconn.PgError)
	if !errors.As(err, &pgErr) {
		return nil
	}
	for key, constraint := range constraints {
		if strings.Contains(pgErr.Detail, key) {
			return constraint
		}
	}
	return nil
}
