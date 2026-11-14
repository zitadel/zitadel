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
	//go:embed unique_constraints_add.sql
	addConstraintStmt string
)

func handleUniqueConstraints(ctx context.Context, tx database.Transaction, commands []eventstore.Command) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	var deleteInstances, deleteTypes, deleteFields, removeInstances []string
	var addInstances, addTypes, addFields, claimantTypes, claimantIDs []string
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
				addInstances = append(addInstances, instanceID)
				addTypes = append(addTypes, constraint.UniqueType)
				addFields = append(addFields, constraint.UniqueField)
				claimantTypes = append(claimantTypes, string(constraint.ClaimantType))
				claimantIDs = append(claimantIDs, constraint.ClaimantID)
				addConstraints[fmt.Sprintf(uniqueConstraintPlaceholderFmt, instanceID, constraint.UniqueType, constraint.UniqueField)] = constraint
			case eventstore.UniqueConstraintRemove:
				if constraint.ClaimantID != "" {
					if err := releaseClaim(ctx, tx, instanceID, constraint); err != nil {
						return err
					}
					continue
				}
				deleteInstances = append(deleteInstances, instanceID)
				deleteTypes = append(deleteTypes, constraint.UniqueType)
				deleteFields = append(deleteFields, constraint.UniqueField)
				deleteConstraints[fmt.Sprintf(uniqueConstraintPlaceholderFmt, instanceID, constraint.UniqueType, constraint.UniqueField)] = constraint
			case eventstore.UniqueConstraintInstanceRemove:
				removeInstances = append(removeInstances, instanceID)
				deleteConstraints[fmt.Sprintf(uniqueConstraintPlaceholderFmt, instanceID, constraint.UniqueType, constraint.UniqueField)] = constraint
			}
		}
	}

	if len(deleteFields) > 0 {
		_, err := tx.Exec(ctx, deleteConstraintStmt, deleteInstances, deleteTypes, deleteFields)
		if err != nil {
			logging.WithError(err).Warn("delete unique constraint failed")
			errMessage := "Errors.Internal"
			if constraint := constraintFromErr(err, deleteConstraints); constraint != nil {
				errMessage = constraint.ErrorMessage
			}
			return zerrors.ThrowInternal(err, "V3-C8l3V", errMessage)
		}
	}
	if len(removeInstances) > 0 {
		_, err := tx.Exec(ctx, "DELETE FROM eventstore.unique_constraints WHERE instance_id = ANY($1::text[])", removeInstances)
		if err != nil {
			logging.WithError(err).Warn("delete instance unique constraints failed")
			errMessage := "Errors.Internal"
			if constraint := constraintFromErr(err, deleteConstraints); constraint != nil {
				errMessage = constraint.ErrorMessage
			}
			return zerrors.ThrowInternal(err, "V3-C8l3V", errMessage)
		}
	}
	if len(addFields) > 0 {
		_, err := tx.Exec(ctx, addConstraintStmt, addInstances, addTypes, addFields, claimantTypes, claimantIDs)
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
