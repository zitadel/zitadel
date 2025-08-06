package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/org"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	ProjectMetadataProjectionTable = "projections.project_metadata"

	ProjectMetadataColumnProjectID     = "project_id"
	ProjectMetadataColumnCreationDate  = "creation_date"
	ProjectMetadataColumnChangeDate    = "change_date"
	ProjectMetadataColumnSequence      = "sequence"
	ProjectMetadataColumnResourceOwner = "resource_owner"
	ProjectMetadataColumnInstanceID    = "instance_id"
	ProjectMetadataColumnKey           = "key"
	ProjectMetadataColumnValue         = "value"
)

type projectMetadataProjection struct{}

func (*projectMetadataProjection) Name() string {
	return ProjectMetadataProjectionTable
}

func (*projectMetadataProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(ProjectMetadataColumnProjectID, handler.ColumnTypeText),
			handler.NewColumn(ProjectMetadataColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(ProjectMetadataColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(ProjectMetadataColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(ProjectMetadataColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(ProjectMetadataColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(ProjectMetadataColumnKey, handler.ColumnTypeText),
			handler.NewColumn(ProjectMetadataColumnValue, handler.ColumnTypeBytes, handler.Nullable()),
		},
			handler.NewPrimaryKey(ProjectMetadataColumnInstanceID, ProjectMetadataColumnProjectID, ProjectMetadataColumnKey),
			handler.WithIndex(handler.NewIndex(ProjectMetadataColumnValue, []string{ProjectMetadataColumnValue})),
		),
	)
}

func newProjectMetadataProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(projectMetadataProjection))
}

func (p *projectMetadataProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  project.MetadataSetType,
					Reduce: p.reduceMetadataSet,
				},
				{
					Event:  project.MetadataRemovedType,
					Reduce: p.reduceMetadataRemoved,
				},
			},
		},
		{
			Aggregate: project.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  project.ProjectRemovedType,
					Reduce: p.reduceProjectRemoved,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
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
					Reduce: reduceInstanceRemovedHelper(ProjectMetadataColumnInstanceID),
				},
			},
		},
	}
}

func (p *projectMetadataProjection) reduceMetadataSet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.MetadataSetEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-BxvwG3", "reduce.wrong.event.type %s", project.MetadataSetType)
	}
	return handler.NewUpsertStatement(
		e,
		[]handler.Column{
			handler.NewCol(ProjectMetadataColumnInstanceID, nil),
			handler.NewCol(ProjectMetadataColumnProjectID, nil),
			handler.NewCol(ProjectMetadataColumnKey, e.Key),
		},
		[]handler.Column{
			handler.NewCol(ProjectMetadataColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(ProjectMetadataColumnProjectID, e.Aggregate().ID),
			handler.NewCol(ProjectMetadataColumnKey, e.Key),
			handler.NewCol(ProjectMetadataColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(ProjectMetadataColumnCreationDate, handler.OnlySetValueOnInsert(ProjectMetadataProjectionTable, e.CreationDate())),
			handler.NewCol(ProjectMetadataColumnChangeDate, e.CreationDate()),
			handler.NewCol(ProjectMetadataColumnSequence, e.Sequence()),
			handler.NewCol(ProjectMetadataColumnValue, e.Value),
		},
	), nil
}

func (p *projectMetadataProjection) reduceMetadataRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.MetadataRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-R20XC8", "reduce.wrong.event.type %s", project.MetadataRemovedType)
	}
	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectMetadataColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(ProjectMetadataColumnProjectID, e.Aggregate().ID),
			handler.NewCond(ProjectMetadataColumnKey, e.Key),
			handler.NewCond(ProjectMetadataColumnResourceOwner, e.Aggregate().ResourceOwner),
		},
	), nil
}

func (p *projectMetadataProjection) reduceOwnerRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*org.OrgRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-zdKjXb", "reduce.wrong.event.type %s", org.OrgRemovedEventType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectMetadataColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(ProjectMetadataColumnResourceOwner, e.Aggregate().ID),
		},
	), nil
}

func (p *projectMetadataProjection) reduceProjectRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*project.ProjectRemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-ydmC4w", "reduce.wrong.event.type %s", project.ProjectRemovedType)
	}

	return handler.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(ProjectMetadataColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(ProjectMetadataColumnProjectID, e.Aggregate().ID),
		},
	), nil
}
