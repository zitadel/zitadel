package eventstore

import (
	"context"
	_ "embed"
	"encoding/json"
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
	//go:embed unique_constraints_delete_by_owner.sql
	deleteConstraintByOwnerStmt string
	//go:embed unique_constraints_add.sql
	addConstraintStmt string
	//go:embed unique_constraints_add_without_owners.sql
	addConstraintWithoutOwnersStmt string
	//go:embed unique_constraints_release_claim.sql
	releaseClaimStmt string
)

func handleUniqueConstraints(ctx context.Context, tx database.Transaction, commands []eventstore.Command) (err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	var deleteInstances, deleteTypes, deleteFields, removeInstances []string
	var ownerInstances, ownerTags []string
	var releaseInstances, releaseTypes, releaseFields, releaseClaimants []string
	var adds []*pendingAdd
	pendingAdds := map[string]*pendingAdd{}
	addConstraints := map[string]*eventstore.UniqueConstraint{}
	deleteConstraints := map[string]*eventstore.UniqueConstraint{}
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
				key := fmt.Sprintf(uniqueConstraintPlaceholderFmt, instanceID, constraint.UniqueType, constraint.UniqueField)
				add := &pendingAdd{instanceID: instanceID, constraint: constraint}
				adds = append(adds, add)
				pendingAdds[key] = add
				addConstraints[key] = constraint
			case eventstore.UniqueConstraintRemove:
				if constraint.Claimant != "" {
					// Adds are batched after all releases, so a release of a key added
					// earlier in the same push cancels the pending add. The stored row is
					// still released so a conflicting or stale claim is not silently kept.
					key := fmt.Sprintf(uniqueConstraintPlaceholderFmt, instanceID, constraint.UniqueType, strings.ToLower(constraint.UniqueField))
					if add := pendingAdds[key]; add != nil {
						if add.constraint.Claimant != constraint.Claimant {
							return eventstore.ErrConcurrentModification
						}
						add.released = true
						delete(pendingAdds, key)
					}
					releaseInstances = append(releaseInstances, instanceID)
					releaseTypes = append(releaseTypes, constraint.UniqueType)
					releaseFields = append(releaseFields, constraint.UniqueField)
					releaseClaimants = append(releaseClaimants, constraint.Claimant)
					continue
				}
				deleteInstances = append(deleteInstances, instanceID)
				deleteTypes = append(deleteTypes, constraint.UniqueType)
				deleteFields = append(deleteFields, constraint.UniqueField)
				deleteConstraints[fmt.Sprintf(uniqueConstraintPlaceholderFmt, instanceID, constraint.UniqueType, constraint.UniqueField)] = constraint
			case eventstore.UniqueConstraintInstanceRemove:
				removeInstances = append(removeInstances, instanceID)
				deleteConstraints[fmt.Sprintf(uniqueConstraintPlaceholderFmt, instanceID, constraint.UniqueType, constraint.UniqueField)] = constraint
			case eventstore.UniqueConstraintRemoveByOwner:
				if len(constraint.Owners) == 0 {
					continue
				}
				ownerInstances = append(ownerInstances, instanceID)
				ownerTags = append(ownerTags, constraint.Owners[0])
				ownerDeleteConstraints[fmt.Sprintf("%s:%s", instanceID, constraint.Owners[0])] = constraint
			}
		}
	}

	if len(releaseFields) > 0 {
		// One statement keeps large cascades from costing a round trip per claim.
		var released bool
		err := tx.QueryRow(ctx, releaseClaimStmt, releaseInstances, releaseTypes, releaseFields, releaseClaimants).Scan(&released)
		if err != nil {
			return zerrors.ThrowInternal(err, "V3-release-claim", "Errors.Internal")
		}
		if !released {
			return eventstore.ErrConcurrentModification
		}
	}
	if len(deleteFields) > 0 {
		if err := execDeleteConstraints(ctx, tx, deleteConstraintStmt, deleteConstraints, deleteInstances, deleteTypes, deleteFields); err != nil {
			return err
		}
	}
	if len(removeInstances) > 0 {
		if err := execDeleteConstraints(ctx, tx, "DELETE FROM eventstore.unique_constraints WHERE instance_id = ANY($1::text[])", deleteConstraints, removeInstances); err != nil {
			return err
		}
	}
	if len(ownerTags) > 0 {
		if err := execDeleteConstraints(ctx, tx, deleteConstraintByOwnerStmt, ownerDeleteConstraints, ownerInstances, ownerTags); err != nil {
			return err
		}
	}
	// Adds without owners do not name the owners column, so they also work before setup step 78.
	var plain, owned constraintAdds
	for _, add := range adds {
		if add.released {
			continue
		}
		if len(add.constraint.Owners) == 0 {
			plain.append(add, "")
			continue
		}
		encoded, err := json.Marshal(add.constraint.Owners)
		if err != nil {
			return zerrors.ThrowInternal(err, "V3-owners", "Errors.Internal")
		}
		owned.append(add, string(encoded))
	}
	if err := execAddConstraints(ctx, tx, addConstraintWithoutOwnersStmt, addConstraints, plain.instances, plain.types, plain.fields); err != nil {
		return err
	}
	return execAddConstraints(ctx, tx, addConstraintStmt, addConstraints, owned.instances, owned.types, owned.fields, owned.owners)
}

type constraintAdds struct {
	instances, types, fields, owners []string
}

func (a *constraintAdds) append(add *pendingAdd, owners string) {
	a.instances = append(a.instances, add.instanceID)
	a.types = append(a.types, add.constraint.UniqueType)
	a.fields = append(a.fields, add.constraint.UniqueField)
	a.owners = append(a.owners, owners)
}

func execAddConstraints(ctx context.Context, tx database.Transaction, stmt string, constraints map[string]*eventstore.UniqueConstraint, instances []string, args ...any) error {
	if len(instances) == 0 {
		return nil
	}
	_, err := tx.Exec(ctx, stmt, append([]any{instances}, args...)...)
	if err != nil {
		logging.WithError(err).Warn("add unique constraint failed")
		errMessage := "Errors.Internal"
		if constraint := constraintFromErr(err, constraints); constraint != nil {
			errMessage = constraint.ErrorMessage
		}
		return zerrors.ThrowAlreadyExists(err, "V3-DKcYh", errMessage)
	}
	return nil
}

type pendingAdd struct {
	instanceID string
	constraint *eventstore.UniqueConstraint
	released   bool
}

func execDeleteConstraints(ctx context.Context, tx database.Transaction, stmt string, constraints map[string]*eventstore.UniqueConstraint, args ...any) error {
	_, err := tx.Exec(ctx, stmt, args...)
	if err != nil {
		logging.WithError(err).Warn("delete unique constraint failed")
		errMessage := "Errors.Internal"
		if constraint := constraintFromErr(err, constraints); constraint != nil {
			errMessage = constraint.ErrorMessage
		}
		return zerrors.ThrowInternal(err, "V3-C8l3V", errMessage)
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
