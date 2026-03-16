package projection

import (
	"context"
	"database/sql"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	v3_sql "github.com/zitadel/zitadel/backend/v3/storage/database/dialect/sql"
	"github.com/zitadel/zitadel/backend/v3/storage/database/repository"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func (p *relationalTablesProjection) reduceInstanceAdminAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.MemberAddedEvent](event)
	if err != nil {
		return nil, err
	}
	return p.createAdministratorStatement(e, &domain.Administrator{
		InstanceID: e.Aggregate().InstanceID,
		UserID:     e.UserID,
		Scope:      domain.AdministratorScopeInstance,
		Roles:      e.Roles,
		CreatedAt:  e.CreationDate(),
		UpdatedAt:  e.CreationDate(),
	})
}

func (p *relationalTablesProjection) reduceInstanceAdminChanged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.MemberChangedEvent](event)
	if err != nil {
		return nil, err
	}
	return p.updateAdministratorStatement(e,
		repository.AdministratorRepository().InstanceAdministratorCondition(e.Aggregate().InstanceID, e.UserID),
		repository.AdministratorRepository().SetUpdatedAt(e.CreationDate()),
		repository.AdministratorRepository().SetRoles(e.Roles),
	)
}

func (p *relationalTablesProjection) reduceInstanceAdminRemoved(event eventstore.Event) (*handler.Statement, error) {
	switch e := event.(type) {
	case *instance.MemberRemovedEvent:
		return p.removeAdministratorStatement(e,
			repository.AdministratorRepository().InstanceAdministratorCondition(e.Aggregate().InstanceID, e.UserID),
		)
	case *instance.MemberCascadeRemovedEvent:
		return p.removeAdministratorStatement(e,
			repository.AdministratorRepository().InstanceAdministratorCondition(e.Aggregate().InstanceID, e.UserID),
		)
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-v3ADM03", "reduce.wrong.event.type %v", []eventstore.EventType{instance.MemberRemovedEventType, instance.MemberCascadeRemovedEventType})
	}
}

func (p *relationalTablesProjection) reduceOrganizationAdminAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*org.MemberAddedEvent](event)
	if err != nil {
		return nil, err
	}
	return p.createAdministratorStatement(e, &domain.Administrator{
		InstanceID:     e.Aggregate().InstanceID,
		UserID:         e.UserID,
		Scope:          domain.AdministratorScopeOrganization,
		OrganizationID: stringPtr(e.Aggregate().ID),
		Roles:          e.Roles,
		CreatedAt:      e.CreationDate(),
		UpdatedAt:      e.CreationDate(),
	})
}

func (p *relationalTablesProjection) reduceOrganizationAdminChanged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*org.MemberChangedEvent](event)
	if err != nil {
		return nil, err
	}
	return p.updateAdministratorStatement(e,
		repository.AdministratorRepository().OrganizationAdministratorCondition(e.Aggregate().InstanceID, e.Aggregate().ID, e.UserID),
		repository.AdministratorRepository().SetUpdatedAt(e.CreationDate()),
		repository.AdministratorRepository().SetRoles(e.Roles),
	)
}

func (p *relationalTablesProjection) reduceOrganizationAdminRemoved(event eventstore.Event) (*handler.Statement, error) {
	switch e := event.(type) {
	case *org.MemberRemovedEvent:
		return p.removeAdministratorStatement(e,
			repository.AdministratorRepository().OrganizationAdministratorCondition(e.Aggregate().InstanceID, e.Aggregate().ID, e.UserID),
		)
	case *org.MemberCascadeRemovedEvent:
		return p.removeAdministratorStatement(e,
			repository.AdministratorRepository().OrganizationAdministratorCondition(e.Aggregate().InstanceID, e.Aggregate().ID, e.UserID),
		)
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-v3ADM06", "reduce.wrong.event.type %v", []eventstore.EventType{org.MemberRemovedEventType, org.MemberCascadeRemovedEventType})
	}
}

func (p *relationalTablesProjection) reduceProjectAdminAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*project.MemberAddedEvent](event)
	if err != nil {
		return nil, err
	}
	return p.createAdministratorStatement(e, &domain.Administrator{
		InstanceID: e.Aggregate().InstanceID,
		UserID:     e.UserID,
		Scope:      domain.AdministratorScopeProject,
		ProjectID:  stringPtr(e.Aggregate().ID),
		Roles:      e.Roles,
		CreatedAt:  e.CreationDate(),
		UpdatedAt:  e.CreationDate(),
	})
}

func (p *relationalTablesProjection) reduceProjectAdminChanged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*project.MemberChangedEvent](event)
	if err != nil {
		return nil, err
	}
	return p.updateAdministratorStatement(e,
		repository.AdministratorRepository().ProjectAdministratorCondition(e.Aggregate().InstanceID, e.Aggregate().ID, e.UserID),
		repository.AdministratorRepository().SetUpdatedAt(e.CreationDate()),
		repository.AdministratorRepository().SetRoles(e.Roles),
	)
}

func (p *relationalTablesProjection) reduceProjectAdminRemoved(event eventstore.Event) (*handler.Statement, error) {
	switch e := event.(type) {
	case *project.MemberRemovedEvent:
		return p.removeAdministratorStatement(e,
			repository.AdministratorRepository().ProjectAdministratorCondition(e.Aggregate().InstanceID, e.Aggregate().ID, e.UserID),
		)
	case *project.MemberCascadeRemovedEvent:
		return p.removeAdministratorStatement(e,
			repository.AdministratorRepository().ProjectAdministratorCondition(e.Aggregate().InstanceID, e.Aggregate().ID, e.UserID),
		)
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-v3ADM09", "reduce.wrong.event.type %v", []eventstore.EventType{project.MemberRemovedEventType, project.MemberCascadeRemovedEventType})
	}
}

func (p *relationalTablesProjection) reduceProjectGrantAdminAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*project.GrantMemberAddedEvent](event)
	if err != nil {
		return nil, err
	}
	return p.createAdministratorStatement(e, &domain.Administrator{
		InstanceID:     e.Aggregate().InstanceID,
		UserID:         e.UserID,
		Scope:          domain.AdministratorScopeProjectGrant,
		ProjectGrantID: stringPtr(e.GrantID),
		Roles:          e.Roles,
		CreatedAt:      e.CreationDate(),
		UpdatedAt:      e.CreationDate(),
	})
}

func (p *relationalTablesProjection) reduceProjectGrantAdminChanged(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*project.GrantMemberChangedEvent](event)
	if err != nil {
		return nil, err
	}
	return p.updateAdministratorStatement(e,
		repository.AdministratorRepository().ProjectGrantAdministratorCondition(e.Aggregate().InstanceID, e.GrantID, e.UserID),
		repository.AdministratorRepository().SetUpdatedAt(e.CreationDate()),
		repository.AdministratorRepository().SetRoles(e.Roles),
	)
}

func (p *relationalTablesProjection) reduceProjectGrantAdminRemoved(event eventstore.Event) (*handler.Statement, error) {
	switch e := event.(type) {
	case *project.GrantMemberRemovedEvent:
		return p.removeAdministratorStatement(e,
			repository.AdministratorRepository().ProjectGrantAdministratorCondition(e.Aggregate().InstanceID, e.GrantID, e.UserID),
		)
	case *project.GrantMemberCascadeRemovedEvent:
		return p.removeAdministratorStatement(e,
			repository.AdministratorRepository().ProjectGrantAdministratorCondition(e.Aggregate().InstanceID, e.GrantID, e.UserID),
		)
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-v3ADM12", "reduce.wrong.event.type %v", []eventstore.EventType{project.GrantMemberRemovedType, project.GrantMemberCascadeRemovedType})
	}
}

func (p *relationalTablesProjection) createAdministratorStatement(event eventstore.Event, administrator *domain.Administrator) (*handler.Statement, error) {
	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-v3ADM13", "reduce.wrong.db.pool %T", ex)
		}
		return repository.AdministratorRepository().Create(ctx, v3_sql.SQLTx(tx), administrator)
	}), nil
}

func (p *relationalTablesProjection) updateAdministratorStatement(event eventstore.Event, condition database.Condition, changes ...database.Change) (*handler.Statement, error) {
	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-v3ADM14", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.AdministratorRepository()
		_, err := repo.Update(ctx, v3_sql.SQLTx(tx), condition, changes...)
		return err
	}), nil
}

func (p *relationalTablesProjection) removeAdministratorStatement(event eventstore.Event, condition database.Condition) (*handler.Statement, error) {
	return handler.NewStatement(event, func(ctx context.Context, ex handler.Executer, _ string) error {
		tx, ok := ex.(*sql.Tx)
		if !ok {
			return zerrors.ThrowInvalidArgumentf(nil, "HANDL-v3ADM15", "reduce.wrong.db.pool %T", ex)
		}
		repo := repository.AdministratorRepository()
		_, err := repo.Delete(ctx, v3_sql.SQLTx(tx), condition)
		return err
	}), nil
}

func stringPtr(value string) *string {
	return &value
}
