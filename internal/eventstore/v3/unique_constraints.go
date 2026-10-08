package eventstore

import (
	"context"
	_ "embed"
	"errors"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v5/pgconn"
	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
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

func handleUniqueConstraints(ctx context.Context, tx database.Transaction, commands []eventstore.Command) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	keyedPlaceholders := make([]string, 0)
	keyedArgs := make([]any, 0)
	ownerPlaceholders := make([]string, 0)
	ownerArgs := make([]any, 0)

	addWithOwnersPlaceholders := make([]string, 0)
	addWithOwnersArgs := make([]any, 0)
	addWithoutOwnersPlaceholders := make([]string, 0)
	addWithoutOwnersArgs := make([]any, 0)
	addConstraints := map[string]*eventstore.UniqueConstraint{}
	keyedDeleteConstraints := map[string]*eventstore.UniqueConstraint{}
	ownerDeleteConstraints := map[string]*eventstore.UniqueConstraint{}

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
				keyedPlaceholders = append(keyedPlaceholders, fmt.Sprintf(deleteConstraintPlaceholdersStmt, len(keyedArgs)+1, len(keyedArgs)+2, len(keyedArgs)+3))
				keyedArgs = append(keyedArgs, instanceID, constraint.UniqueType, constraint.UniqueField)
				keyedDeleteConstraints[fmt.Sprintf(uniqueConstraintPlaceholderFmt, instanceID, constraint.UniqueType, constraint.UniqueField)] = constraint
			case eventstore.UniqueConstraintInstanceRemove:
				keyedPlaceholders = append(keyedPlaceholders, fmt.Sprintf("(instance_id = $%d)", len(keyedArgs)+1))
				keyedArgs = append(keyedArgs, instanceID)
				keyedDeleteConstraints[fmt.Sprintf(uniqueConstraintPlaceholderFmt, instanceID, constraint.UniqueType, constraint.UniqueField)] = constraint
			case eventstore.UniqueConstraintRemoveByOwner:
				if len(constraint.Owners) == 0 {
					continue
				}
				tag := constraint.Owners[0]
				ownerPlaceholders = append(ownerPlaceholders, fmt.Sprintf("(instance_id = $%d AND owners @> ARRAY[$%d]::text[])", len(ownerArgs)+1, len(ownerArgs)+2))
				ownerArgs = append(ownerArgs, instanceID, tag)
				ownerDeleteConstraints[fmt.Sprintf("%s:%s", instanceID, tag)] = constraint
			}
		}
	}

	if err := execDeleteConstraints(ctx, tx, keyedPlaceholders, keyedArgs, keyedDeleteConstraints); err != nil {
		return err
	}
	if err := execDeleteConstraints(ctx, tx, ownerPlaceholders, ownerArgs, ownerDeleteConstraints); err != nil {
		return err
	}
	if err := execAddConstraints(ctx, tx, addConstraintWithoutOwnersStmt, addWithoutOwnersPlaceholders, addWithoutOwnersArgs, addConstraints); err != nil {
		return err
	}
	return execAddConstraints(ctx, tx, addConstraintStmt, addWithOwnersPlaceholders, addWithOwnersArgs, addConstraints)
}

func execDeleteConstraints(ctx context.Context, tx database.Transaction, placeholders []string, args []any, deleteConstraints map[string]*eventstore.UniqueConstraint) error {
	if len(placeholders) == 0 {
		return nil
	}
	_, err := tx.Exec(ctx, fmt.Sprintf(deleteConstraintStmt, strings.Join(placeholders, " OR ")), args...)
	if err != nil {
		logging.WithError(err).Warn("delete unique constraint failed")
		errMessage := "Errors.Internal"
		if constraint := constraintFromErr(err, deleteConstraints); constraint != nil {
			errMessage = constraint.ErrorMessage
		}
		return zerrors.ThrowInternal(err, "V3-C8l3V", errMessage)
	}
	return nil
}

func execAddConstraints(ctx context.Context, tx database.Transaction, stmt string, placeholders []string, args []any, addConstraints map[string]*eventstore.UniqueConstraint) error {
	if len(placeholders) == 0 {
		return nil
	}
	_, err := tx.Exec(ctx, fmt.Sprintf(stmt, strings.Join(placeholders, ", ")), args...)
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
