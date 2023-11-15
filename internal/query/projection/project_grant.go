package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
)

const (
	ProjectGrantProjectionTable = "projections.project_grants4"

	ProjectGrantColumnGrantID       = "grant_id"
	ProjectGrantColumnCreationDate  = "creation_date"
	ProjectGrantColumnChangeDate    = "change_date"
	ProjectGrantColumnSequence      = "sequence"
	ProjectGrantColumnState         = "state"
	ProjectGrantColumnResourceOwner = "resource_owner"
	ProjectGrantColumnInstanceID    = "instance_id"
	ProjectGrantColumnProjectID     = "project_id"
	ProjectGrantColumnGrantedOrgID  = "granted_org_id"
	ProjectGrantColumnRoleKeys      = "granted_role_keys"
)

type projectGrantProjection struct{}

func newProjectGrantProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(projectGrantProjection))
}

func (*projectGrantProjection) Name() string {
	return ProjectGrantProjectionTable
}

func (*projectGrantProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(ProjectGrantColumnGrantID, handler.ColumnTypeText),
			handler.NewColumn(ProjectGrantColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(ProjectGrantColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(ProjectGrantColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(ProjectGrantColumnState, handler.ColumnTypeEnum),
			handler.NewColumn(ProjectGrantColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(ProjectGrantColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(ProjectGrantColumnProjectID, handler.ColumnTypeText),
			handler.NewColumn(ProjectGrantColumnGrantedOrgID, handler.ColumnTypeText),
			handler.NewColumn(ProjectGrantColumnRoleKeys, handler.ColumnTypeTextArray, handler.Nullable()),
		},
			handler.NewPrimaryKey(ProjectGrantColumnInstanceID, ProjectGrantColumnGrantID),
			handler.WithIndex(handler.NewIndex("resource_owner", []string{ProjectGrantColumnResourceOwner})),
			handler.WithIndex(handler.NewIndex("granted_org", []string{ProjectGrantColumnGrantedOrgID})),
		),
	)
}

func (p *projectGrantProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventReducers: []handler.EventReducer{
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
			EventReducers: []handler.EventReducer{
				{
					Event:  org.OrgRemovedEventType,
					Reduce: p.reduceOwnerRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
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
	return handler.NewCreateStatement(
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
			handler.NewCol(ProjectGrantColumnRoleKeys, database.TextArray[string](e.RoleKeys)),
		},
	), nil
}

func (p *projectGrantProjection) reduceProjectGrantChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.GrantChangedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-g0fg4", "reduce.wrong.event.type %s", project.GrantChangedType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectGrantColumnSequence, e.Sequence()),
			handler.NewCol(ProjectGrantColumnRoleKeys, database.TextArray[string](e.RoleKeys)),
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
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectGrantColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectGrantColumnSequence, e.Sequence()),
			handler.NewCol(ProjectGrantColumnRoleKeys, database.TextArray[string](e.RoleKeys)),
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
	return handler.NewUpdateStatement(
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
	return handler.NewUpdateStatement(
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
	return handler.NewDeleteStatement(
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
	return handler.NewDeleteStatement(
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

	return handler.NewMultiStatement(
		e,
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(ProjectGrantColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(ProjectGrantColumnResourceOwner, e.Aggregate().ID),
			},
		),
		handler.AddDeleteStatement(
			[]handler.Condition{
				handler.NewCond(ProjectGrantColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(ProjectGrantColumnGrantedOrgID, e.Aggregate().ID),
			},
		),
	), nil
}
