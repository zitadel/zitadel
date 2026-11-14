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
)

func handleUniqueConstraints(ctx context.Context, tx database.Transaction, commands []eventstore.Command) (err error) {
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
			if (constraint.ClaimantType == "") != (constraint.ClaimantID == "") {
				return zerrors.ThrowInvalidArgument(nil, "V3-claimant", "Reservation claimant requires both type and ID")
			}
			instanceID := command.Aggregate().InstanceID
			if constraint.IsGlobal {
				instanceID = ""
			}
			switch constraint.Action {
			case eventstore.UniqueConstraintAdd:
				constraint.UniqueField = strings.ToLower(constraint.UniqueField)
				addPlaceholders = append(addPlaceholders, fmt.Sprintf("($%d, $%d, $%d, NULLIF($%d, ''), NULLIF($%d, ''))", len(addArgs)+1, len(addArgs)+2, len(addArgs)+3, len(addArgs)+4, len(addArgs)+5))
				addArgs = append(addArgs, instanceID, constraint.UniqueType, constraint.UniqueField, string(constraint.ClaimantType), constraint.ClaimantID)
				addConstraints[fmt.Sprintf(uniqueConstraintPlaceholderFmt, instanceID, constraint.UniqueType, constraint.UniqueField)] = constraint
			case eventstore.UniqueConstraintRemove:
				if constraint.ClaimantID != "" {
					if err := releaseClaim(ctx, tx, instanceID, constraint); err != nil {
						return err
					}
					continue
				}
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
		_, err := tx.Exec(ctx, fmt.Sprintf(deleteConstraintStmt, strings.Join(deletePlaceholders, " OR ")), deleteArgs...)
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
		_, err := tx.Exec(ctx, fmt.Sprintf(addConstraintStmt, strings.Join(addPlaceholders, ", ")), addArgs...)
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

func releaseClaim(ctx context.Context, tx database.Transaction, instanceID string, constraint *eventstore.UniqueConstraint) error {
	count, err := tx.Exec(ctx, `DELETE FROM eventstore.unique_constraints
		WHERE instance_id = $1 AND unique_type = $2 AND unique_field = LOWER($3)
		AND claimant_type = $4 AND claimant_id = $5`,
		instanceID, constraint.UniqueType, constraint.UniqueField, string(constraint.ClaimantType), constraint.ClaimantID)
	if err != nil {
		return zerrors.ThrowInternal(err, "V3-release-claim", "Errors.Internal")
	}
	if count != 1 {
		return eventstore.ErrConcurrentModification
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
