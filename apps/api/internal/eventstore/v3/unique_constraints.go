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
)

func handleUniqueConstraints(ctx context.Context, tx database.Tx, commands []eventstore.Command) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	deletePlaceholders := make([]string, 0)
	deleteArgs := make([]any, 0)

	addPlaceholders := make([]string, 0)
	addArgs := make([]any, 0)
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
				addPlaceholders = append(addPlaceholders, fmt.Sprintf("($%d, $%d, $%d)", len(addArgs)+1, len(addArgs)+2, len(addArgs)+3))
				addArgs = append(addArgs, instanceID, constraint.UniqueType, constraint.UniqueField)
				addConstraints[fmt.Sprintf(uniqueConstraintPlaceholderFmt, instanceID, constraint.UniqueType, constraint.UniqueField)] = constraint
			case eventstore.UniqueConstraintRemove:
				deletePlaceholders = append(deletePlaceholders, fmt.Sprintf(deleteConstraintPlaceholdersStmt, len(deleteArgs)+1, len(deleteArgs)+2, len(deleteArgs)+3))
				deleteArgs = append(deleteArgs, instanceID, constraint.UniqueType, constraint.UniqueField)
				deleteConstraints[fmt.Sprintf(uniqueConstraintPlaceholderFmt, instanceID, constraint.UniqueType, constraint.UniqueField)] = constraint
			case eventstore.UniqueConstraintInstanceRemove:
				deletePlaceholders = append(deletePlaceholders, fmt.Sprintf("(instance_id = $%d)", len(deleteArgs)+1))
				deleteArgs = append(deleteArgs, instanceID)
				deleteConstraints[fmt.Sprintf(uniqueConstraintPlaceholderFmt, instanceID, constraint.UniqueType, constraint.UniqueField)] = constraint
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
	if len(addPlaceholders) > 0 {
		_, err := tx.ExecContext(ctx, fmt.Sprintf(addConstraintStmt, strings.Join(addPlaceholders, ", ")), addArgs...)
		if err != nil {
			logging.WithError(err).Warn("add unique constraint failed")
			errMessage := "Errors.Internal"
			if constraint := constraintFromErr(err, addConstraints); constraint != nil {
				errMessage = constraint.ErrorMessage
			}
			return zerrors.ThrowAlreadyExists(err, "V3-DKcYh", errMessage)
		}
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
