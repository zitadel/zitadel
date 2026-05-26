package projection

import (
	"context"
	"database/sql"

	"github.com/muhlemmer/gu"

	repoDomain "github.com/zitadel/zitadel/backend/v3/domain"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/usergrant"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (p *relationalTablesProjection) reduceAuthorizationAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*usergrant.UserGrantAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInternalf(nil, "HANDL-an0k9a", "reduce.wrong.event.type %s", usergrant.UserGrantAddedType)
	}
	return handler.NewStatement(e, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-FrDVPk", "reduce.wrong.db.pool %T", ex)
		}

		// Filter the granted role keys down to ones that actually exist in
		// project_roles for this project. Historical user_grant.added events
		// can reference roles that were since removed from the project — the
		// authorization_roles_*_role_key_fkey would otherwise fail forever.
		// Keeping only the still-valid roles preserves whatever access the user
		// can legitimately retain; if every role is gone, we still create the
		// authorization with an empty role set (the user/project/grant linkage
		// itself is preserved for audit and Changed events to update later).
		validRoles, err := filterExistingRoleKeys(ctx, tx, e.Aggregate().InstanceID, e.ProjectID, e.RoleKeys)
		if err != nil {
			return err
		}

		repo := repository.AuthorizationRepository()
		var grantID *string
		if e.ProjectGrantID != "" {
			grantID = gu.Ptr(e.ProjectGrantID)
		}
		return repo.Create(ctx, v3_sql.SQLTx(tx), &repoDomain.Authorization{
			InstanceID: e.Aggregate().InstanceID,
			ID:         e.Aggregate().ID,
			UserID:     e.UserID,
			ProjectID:  e.ProjectID,
			GrantID:    grantID,
			Roles:      validRoles,
			CreatedAt:  e.CreationDate(),
			UpdatedAt:  e.CreationDate(),
			State:      repoDomain.AuthorizationStateActive,
		})
	}), nil
}

// filterExistingRoleKeys returns the subset of the given role keys that exist
// in zitadel.project_roles for the (instance, project) pair. Used by the
// authorization and project_grant reducers to avoid FK violations when historical
// events reference roles that were later removed from a project.
func filterExistingRoleKeys(ctx context.Context, tx *sql.Tx, instanceID, projectID string, requested []string) ([]string, error) {
	if len(requested) == 0 {
		return requested, nil
	}
	rows, err := tx.QueryContext(ctx,
		`SELECT key FROM zitadel.project_roles WHERE instance_id = $1 AND project_id = $2 AND key = ANY($3)`,
		instanceID, projectID, requested,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var existing []string
	for rows.Next() {
		var k string
		if err := rows.Scan(&k); err != nil {
			return nil, err
		}
		existing = append(existing, k)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}
	return existing, nil
}

func (p *relationalTablesProjection) reduceAuthorizationChanged(event eventstore.Event) (*handler.Statement, error) {
	var roles []string
	switch e := event.(type) {
	case *usergrant.UserGrantChangedEvent:
		roles = e.RoleKeys
	case *usergrant.UserGrantCascadeChangedEvent:
		roles = e.RoleKeys
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-kg1ryL", "reduce.wrong.event.type %v", []eventstore.EventType{usergrant.UserGrantChangedType, usergrant.UserGrantCascadeChangedType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-sBwGEW", "reduce.wrong.db.pool %T", ex)
		}

		// Changed/cascade.changed events do not carry project_id in payload.
		// Resolve the existing authorization's project_id so we can filter role keys
		// against zitadel.project_roles and avoid FK violations on update.
		projectID, err := authorizationProjectID(ctx, tx, event.Aggregate().InstanceID, event.Aggregate().ID)
		if err != nil {
			return err
		}
		validRoles, err := filterExistingRoleKeys(ctx, tx, event.Aggregate().InstanceID, projectID, roles)
		if err != nil {
			return err
		}

		repo := repository.AuthorizationRepository()
		_, err = repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(event.Aggregate().InstanceID, event.Aggregate().ID),
			validRoles,
			repo.SetUpdatedAt(event.CreatedAt()),
		)
		return err
	}), nil
}

// authorizationProjectID resolves project_id for an existing authorization row.
func authorizationProjectID(ctx context.Context, tx *sql.Tx, instanceID, authorizationID string) (string, error) {
	var projectID string
	err := tx.QueryRowContext(ctx,
		`SELECT project_id FROM zitadel.authorizations WHERE instance_id = $1 AND id = $2`,
		instanceID, authorizationID,
	).Scan(&projectID)
	if err != nil {
		if err == sql.ErrNoRows {
			// No existing row to update; returning empty project keeps behavior noop-like.
			return "", nil
		}
		return "", err
	}
	return projectID, nil
}

func (p *relationalTablesProjection) reduceAuthorizationRemoved(event eventstore.Event) (*handler.Statement, error) {
	switch event.(type) {
	case *usergrant.UserGrantRemovedEvent, *usergrant.UserGrantCascadeRemovedEvent:
		// ok
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-dFtRl9", "reduce.wrong.event.type %v", []eventstore.EventType{usergrant.UserGrantRemovedType, usergrant.UserGrantCascadeRemovedType})
	}

	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-FlZ55O", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.AuthorizationRepository()
		_, err := repo.Delete(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(event.Aggregate().InstanceID, event.Aggregate().ID),
		)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceAuthorizationDeactivated(event eventstore.Event) (*handler.Statement, error) {
	_, ok := event.(*usergrant.UserGrantDeactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-dFtRl0", "reduce.wrong.event.type %s", usergrant.UserGrantDeactivatedType)
	}
	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-FlZ55P", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.AuthorizationRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(event.Aggregate().InstanceID, event.Aggregate().ID),
			nil,
			repo.SetUpdatedAt(event.CreatedAt()),
			repo.SetState(repoDomain.AuthorizationStateInactive),
		)
		return err
	}), nil
}

func (p *relationalTablesProjection) reduceAuthorizationReactivated(event eventstore.Event) (*handler.Statement, error) {
	_, ok := event.(*usergrant.UserGrantReactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-dFtRl1", "reduce.wrong.event.type %s", usergrant.UserGrantReactivatedType)
	}
	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-FlZ55Q", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.AuthorizationRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx),
			repo.PrimaryKeyCondition(event.Aggregate().InstanceID, event.Aggregate().ID),
			nil,
			repo.SetUpdatedAt(event.CreatedAt()),
			repo.SetState(repoDomain.AuthorizationStateActive),
		)
		return err
	}), nil
}
