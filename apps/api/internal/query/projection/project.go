package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	ProjectProjectionTable = "projections.projects4"

	ProjectColumnID                     = "id"
	ProjectColumnCreationDate           = "creation_date"
	ProjectColumnChangeDate             = "change_date"
	ProjectColumnSequence               = "sequence"
	ProjectColumnState                  = "state"
	ProjectColumnResourceOwner          = "resource_owner"
	ProjectColumnInstanceID             = "instance_id"
	ProjectColumnName                   = "name"
	ProjectColumnProjectRoleAssertion   = "project_role_assertion"
	ProjectColumnProjectRoleCheck       = "project_role_check"
	ProjectColumnHasProjectCheck        = "has_project_check"
	ProjectColumnPrivateLabelingSetting = "private_labeling_setting"
)

type projectProjection struct{}

func newProjectProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(projectProjection))
}

func (*projectProjection) Name() string {
	return ProjectProjectionTable
}

func (*projectProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(ProjectColumnID, handler.ColumnTypeText),
			handler.NewColumn(ProjectColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(ProjectColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(ProjectColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(ProjectColumnState, handler.ColumnTypeEnum),
			handler.NewColumn(ProjectColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(ProjectColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(ProjectColumnName, handler.ColumnTypeText),
			handler.NewColumn(ProjectColumnProjectRoleAssertion, handler.ColumnTypeBool),
			handler.NewColumn(ProjectColumnProjectRoleCheck, handler.ColumnTypeBool),
			handler.NewColumn(ProjectColumnHasProjectCheck, handler.ColumnTypeBool),
			handler.NewColumn(ProjectColumnPrivateLabelingSetting, handler.ColumnTypeEnum),
		},
			handler.NewPrimaryKey(ProjectColumnInstanceID, ProjectColumnID),
			handler.WithIndex(handler.NewIndex("resource_owner", []string{ProjectColumnResourceOwner})),
		),
	)
}

func (p *projectProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  project.ProjectAddedType,
					Reduce: p.reduceProjectAdded,
				},
				{
					Event:  project.ProjectChangedType,
					Reduce: p.reduceProjectChanged,
				},
				{
					Event:  project.ProjectDeactivatedType,
					Reduce: p.reduceProjectDeactivated,
				},
				{
					Event:  project.ProjectReactivatedType,
					Reduce: p.reduceProjectReactivated,
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
					Reduce: reduceInstanceRemovedHelper(ProjectColumnInstanceID),
				},
			},
		},
	}
}

func (p *projectProjection) reduceProjectAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-l000S", "reduce.wrong.event.type %s", project.ProjectAddedType)
	}
	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectColumnID, e.Aggregate().ID),
			handler.NewCol(ProjectColumnCreationDate, e.CreationDate()),
			handler.NewCol(ProjectColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(ProjectColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(ProjectColumnSequence, e.Sequence()),
			handler.NewCol(ProjectColumnName, e.Name),
			handler.NewCol(ProjectColumnProjectRoleAssertion, e.ProjectRoleAssertion),
			handler.NewCol(ProjectColumnProjectRoleCheck, e.ProjectRoleCheck),
			handler.NewCol(ProjectColumnHasProjectCheck, e.HasProjectCheck),
			handler.NewCol(ProjectColumnPrivateLabelingSetting, e.PrivateLabelingSetting),
			handler.NewCol(ProjectColumnState, domain.ProjectStateActive),
		},
	), nil
}

func (p *projectProjection) reduceProjectChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectChangeEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-s00Fs", "reduce.wrong.event.type %s", project.ProjectChangedType)
	}
	if e.Name == nil && e.HasProjectCheck == nil && e.ProjectRoleAssertion == nil && e.ProjectRoleCheck == nil && e.PrivateLabelingSetting == nil {
		return handler.NewNoOpStatement(e), nil
	}

	columns := make([]handler.Column, 0, 7)
	columns = append(columns, handler.NewCol(ProjectColumnChangeDate, e.CreationDate()),
		handler.NewCol(ProjectColumnSequence, e.Sequence()))
	if e.Name != nil {
		columns = append(columns, handler.NewCol(ProjectColumnName, *e.Name))
	}
	if e.ProjectRoleAssertion != nil {
		columns = append(columns, handler.NewCol(ProjectColumnProjectRoleAssertion, *e.ProjectRoleAssertion))
	}
	if e.ProjectRoleCheck != nil {
		columns = append(columns, handler.NewCol(ProjectColumnProjectRoleCheck, *e.ProjectRoleCheck))
	}
	if e.HasProjectCheck != nil {
		columns = append(columns, handler.NewCol(ProjectColumnHasProjectCheck, *e.HasProjectCheck))
	}
	if e.PrivateLabelingSetting != nil {
		columns = append(columns, handler.NewCol(ProjectColumnPrivateLabelingSetting, *e.PrivateLabelingSetting))
	}
	return handler.NewUpdateStatement(
		e,
		columns,
		[]handler.Condition{
			handler.NewCond(ProjectColumnID, e.Aggregate().ID),
			handler.NewCond(ProjectColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *projectProjection) reduceProjectDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectDeactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-LLp0f", "reduce.wrong.event.type %s", project.ProjectDeactivatedType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectColumnSequence, e.Sequence()),
			handler.NewCol(ProjectColumnState, domain.ProjectStateInactive),
		},
		[]handler.Condition{
			handler.NewCond(ProjectColumnID, e.Aggregate().ID),
			handler.NewCond(ProjectColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *projectProjection) reduceProjectReactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectReactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-9J98f", "reduce.wrong.event.type %s", project.ProjectReactivatedType)
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectColumnSequence, e.Sequence()),
			handler.NewCol(ProjectColumnState, domain.ProjectStateActive),
		},
		[]handler.Condition{
			handler.NewCond(ProjectColumnID, e.Aggregate().ID),
			handler.NewCond(ProjectColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *projectProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-5N9fs", "reduce.wrong.event.type %s", project.ProjectRemovedType)
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectColumnID, e.Aggregate().ID),
			handler.NewCond(ProjectColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *projectProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-sbgru", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(ProjectColumnResourceOwner, e.Aggregate().ID),
		},
	), nil
}
