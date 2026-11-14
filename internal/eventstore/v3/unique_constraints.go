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
	//go:embed unique_constraints_release_claim.sql
	releaseClaimStmt string
)

func handleUniqueConstraints(ctx context.Context, tx database.Transaction, commands []eventstore.Command) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	var deleteInstances, deleteTypes, deleteFields, removeInstances []string
	var releaseInstances, releaseTypes, releaseFields, releaseClaimantTypes, releaseClaimantIDs []string
	var adds []*pendingAdd
	pendingAdds := map[string]*pendingAdd{}
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
				key := fmt.Sprintf(uniqueConstraintPlaceholderFmt, instanceID, constraint.UniqueType, constraint.UniqueField)
				add := &pendingAdd{instanceID: instanceID, constraint: constraint}
				adds = append(adds, add)
				pendingAdds[key] = add
				addConstraints[key] = constraint
			case eventstore.UniqueConstraintRemove:
				if constraint.ClaimantID != "" {
					// Adds are batched after all releases, so a release of a key added
					// earlier in the same push cancels the pending add. The stored row is
					// still released so a conflicting or stale claim is not silently kept.
					key := fmt.Sprintf(uniqueConstraintPlaceholderFmt, instanceID, constraint.UniqueType, strings.ToLower(constraint.UniqueField))
					if add := pendingAdds[key]; add != nil {
						if add.constraint.ClaimantType != constraint.ClaimantType || add.constraint.ClaimantID != constraint.ClaimantID {
							return eventstore.ErrConcurrentModification
						}
						add.released = true
						delete(pendingAdds, key)
					}
					releaseInstances = append(releaseInstances, instanceID)
					releaseTypes = append(releaseTypes, constraint.UniqueType)
					releaseFields = append(releaseFields, constraint.UniqueField)
					releaseClaimantTypes = append(releaseClaimantTypes, string(constraint.ClaimantType))
					releaseClaimantIDs = append(releaseClaimantIDs, constraint.ClaimantID)
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

	if len(releaseFields) > 0 {
		// One statement keeps large cascades from costing a round trip per claim.
		var released bool
		err := tx.QueryRow(ctx, releaseClaimStmt,
			releaseInstances, releaseTypes, releaseFields, releaseClaimantTypes, releaseClaimantIDs).Scan(&released)
		if err != nil {
			return zerrors.ThrowInternal(err, "V3-release-claim", "Errors.Internal")
		}
		if !released {
			return eventstore.ErrConcurrentModification
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
	var addInstances, addTypes, addFields, claimantTypes, claimantIDs []string
	for _, add := range adds {
		if add.released {
			continue
		}
		addInstances = append(addInstances, add.instanceID)
		addTypes = append(addTypes, add.constraint.UniqueType)
		addFields = append(addFields, add.constraint.UniqueField)
		claimantTypes = append(claimantTypes, string(add.constraint.ClaimantType))
		claimantIDs = append(claimantIDs, add.constraint.ClaimantID)
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

type pendingAdd struct {
	instanceID string
	constraint *eventstore.UniqueConstraint
	released   bool
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
