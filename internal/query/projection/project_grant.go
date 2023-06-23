package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
)

const (
	ProjectGrantProjectionTable = "projections.project_grants3"

	ProjectGrantColumnGrantID           = "grant_id"
	ProjectGrantColumnCreationDate      = "creation_date"
	ProjectGrantColumnChangeDate        = "change_date"
	ProjectGrantColumnSequence          = "sequence"
	ProjectGrantColumnState             = "state"
	ProjectGrantColumnResourceOwner     = "resource_owner"
	ProjectGrantColumnInstanceID        = "instance_id"
	ProjectGrantColumnProjectID         = "project_id"
	ProjectGrantColumnGrantedOrgID      = "granted_org_id"
	ProjectGrantColumnRoleKeys          = "granted_role_keys"
	ProjectGrantColumnOwnerRemoved      = "owner_removed"
	ProjectGrantColumnGrantedOrgRemoved = "granted_org_removed"
)

type projectGrantProjection struct {
	crdb.StatementHandler
}

func newProjectGrantProjection(ctx context.Context, config crdb.StatementHandlerConfig) *projectGrantProjection {
	p := new(projectGrantProjection)
	config.ProjectionName = ProjectGrantProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(ProjectGrantColumnGrantID, crdb.ColumnTypeText),
			crdb.NewColumn(ProjectGrantColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(ProjectGrantColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(ProjectGrantColumnSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(ProjectGrantColumnState, crdb.ColumnTypeEnum),
			crdb.NewColumn(ProjectGrantColumnResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(ProjectGrantColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(ProjectGrantColumnProjectID, crdb.ColumnTypeText),
			crdb.NewColumn(ProjectGrantColumnGrantedOrgID, crdb.ColumnTypeText),
			crdb.NewColumn(ProjectGrantColumnRoleKeys, crdb.ColumnTypeTextArray, crdb.Nullable()),
			crdb.NewColumn(ProjectGrantColumnOwnerRemoved, crdb.ColumnTypeBool, crdb.Default(false)),
			crdb.NewColumn(ProjectGrantColumnGrantedOrgRemoved, crdb.ColumnTypeBool, crdb.Default(false)),
		},
			crdb.NewPrimaryKey(ProjectGrantColumnInstanceID, ProjectGrantColumnGrantID),
			crdb.WithIndex(crdb.NewIndex("resource_owner", []string{ProjectGrantColumnResourceOwner})),
			crdb.WithIndex(crdb.NewIndex("granted_org", []string{ProjectGrantColumnGrantedOrgID})),
			crdb.WithIndex(crdb.NewIndex("owner_removed", []string{ProjectGrantColumnOwnerRemoved})),
			crdb.WithIndex(crdb.NewIndex("granted_org_removed", []string{ProjectGrantColumnGrantedOrgRemoved})),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *projectGrantProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  project.GrantAddedType,
					Reduce: p.reduceProjectGrantAdded,
				},
				{
					Event:  project.GrantChangedType,
					Reduce: p.reduceProjectGrantChanged,
				},
				{
					Event:  project.GrantCascadeChangedType,
					Reduce: p.reduceProjectGrantCascadeChanged,
				},
				{
					Event:  project.GrantDeactivatedType,
					Reduce: p.reduceProjectGrantDeactivated,
				},
				{
					Event:  project.GrantReactivatedType,
					Reduce: p.reduceProjectGrantReactivated,
				},
				{
					Event:  project.GrantRemovedType,
					Reduce: p.reduceProjectGrantRemoved,
				},
				{
					Event:  project.ProjectRemovedType,
					Reduce: p.reduceProjectRemoved,
				},
			},
		},
		{
			Aggregate: org.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(ProjectGrantColumnInstanceID),
				},
			},
		},
	}
}

func (p *projectGrantProjection) reduceProjectGrantAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-g92Fg", "reduce.wrong.event.type %s", project.GrantAddedType)
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectGrantColumnGrantID, e.GrantID),
			handler.NewCol(ProjectGrantColumnProjectID, e.Aggregate().ID),
			handler.NewCol(ProjectGrantColumnCreationDate, e.CreationDate()),
			handler.NewCol(ProjectGrantColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectGrantColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(ProjectGrantColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(ProjectGrantColumnState, domain.ProjectGrantStateActive),
			handler.NewCol(ProjectGrantColumnSequence, e.Sequence()),
			handler.NewCol(ProjectGrantColumnGrantedOrgID, e.GrantedOrgID),
			handler.NewCol(ProjectGrantColumnRoleKeys, database.StringArray(e.RoleKeys)),
		},
	), nil
}

func (p *projectGrantProjection) reduceProjectGrantChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-g0fg4", "reduce.wrong.event.type %s", project.GrantChangedType)
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectGrantColumnSequence, e.Sequence()),
			handler.NewCol(ProjectGrantColumnRoleKeys, database.StringArray(e.RoleKeys)),
		},
		[]handler.Condition{
			handler.NewCond(ProjectGrantColumnGrantID, e.GrantID),
			handler.NewCond(ProjectGrantColumnProjectID, e.Aggregate().ID),
			handler.NewCond(ProjectGrantColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *projectGrantProjection) reduceProjectGrantCascadeChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantCascadeChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-ll9Ts", "reduce.wrong.event.type %s", project.GrantCascadeChangedType)
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectGrantColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectGrantColumnSequence, e.Sequence()),
			handler.NewCol(ProjectGrantColumnRoleKeys, database.StringArray(e.RoleKeys)),
		},
		[]handler.Condition{
			handler.NewCond(ProjectGrantColumnGrantID, e.GrantID),
			handler.NewCond(ProjectGrantColumnProjectID, e.Aggregate().ID),
			handler.NewCond(ProjectGrantColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *projectGrantProjection) reduceProjectGrantDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantDeactivateEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-0fj2f", "reduce.wrong.event.type %s", project.GrantDeactivatedType)
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectGrantColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectGrantColumnSequence, e.Sequence()),
			handler.NewCol(ProjectGrantColumnState, domain.ProjectGrantStateInactive),
		},
		[]handler.Condition{
			handler.NewCond(ProjectGrantColumnGrantID, e.GrantID),
			handler.NewCond(ProjectGrantColumnProjectID, e.Aggregate().ID),
			handler.NewCond(ProjectGrantColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *projectGrantProjection) reduceProjectGrantReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantReactivatedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-2M0ve", "reduce.wrong.event.type %s", project.GrantReactivatedType)
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectGrantColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectGrantColumnSequence, e.Sequence()),
			handler.NewCol(ProjectGrantColumnState, domain.ProjectGrantStateActive),
		},
		[]handler.Condition{
			handler.NewCond(ProjectGrantColumnGrantID, e.GrantID),
			handler.NewCond(ProjectGrantColumnProjectID, e.Aggregate().ID),
			handler.NewCond(ProjectGrantColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *projectGrantProjection) reduceProjectGrantRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-o0w4f", "reduce.wrong.event.type %s", project.GrantRemovedType)
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectGrantColumnGrantID, e.GrantID),
			handler.NewCond(ProjectGrantColumnProjectID, e.Aggregate().ID),
			handler.NewCond(ProjectGrantColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *projectGrantProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-gn9rw", "reduce.wrong.event.type %s", project.ProjectRemovedType)
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectGrantColumnProjectID, e.Aggregate().ID),
			handler.NewCond(ProjectGrantColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *projectGrantProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "PROJE-HDgW3", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return crdb.NewMultiStatement(
		e,
		crdb.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(ProjectGrantColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(ProjectGrantColumnResourceOwner, e.Aggregate().ID),
			},
		),
		crdb.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(ProjectGrantColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(ProjectGrantColumnGrantedOrgID, e.Aggregate().ID),
			},
		),
	), nil
}
