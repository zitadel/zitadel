package projection

import (
	"context"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/member"
	"github.com/caos/zitadel/internal/repository/project"
)

type ProjectGrantMemberProjection struct {
	crdb.StatementHandler
}

const (
	ProjectGrantMemberProjectionTable = "zitadel.projections.project_grant_members"
)

func NewProjectGrantMemberProjection(ctx context.Context, config crdb.StatementHandlerConfig) *ProjectGrantMemberProjection {
	p := &ProjectGrantMemberProjection{}
	config.ProjectionName = ProjectGrantMemberProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *ProjectGrantMemberProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  project.GrantMemberAddedType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  project.GrantMemberChangedType,
					Reduce: p.reduceChanged,
				},
				{
					Event:  project.GrantMemberCascadeRemovedType,
					Reduce: p.reduceCascadeRemoved,
				},
				{
					Event:  project.GrantMemberRemovedType,
					Reduce: p.reduceRemoved,
				},
			},
		},
	}
}

type ProjectGrantMemberColumn string

const (
	ProjectGrantMemberProjectIDCol = "project_id"
	ProjectGrantMemberGrantIDCol   = "grant_id"
)

func (p *ProjectGrantMemberProjection) reduceAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.GrantMemberAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-csr8B", "seq", event.Sequence(), "expectedType", project.GrantMemberAddedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-0EBQf", "reduce.wrong.event.type")
	}
	return reduceMemberAdded(
		*member.NewMemberAddedEvent(&e.BaseEvent, e.UserID, e.Roles...),
		withMemberCol(ProjectGrantMemberProjectIDCol, e.Aggregate().ID),
		withMemberCol(ProjectGrantMemberGrantIDCol, e.GrantID),
	)
}

func (p *ProjectGrantMemberProjection) reduceChanged(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.GrantMemberChangedEvent)
	if !ok {
		logging.LogWithFields("HANDL-ZubbI", "seq", event.Sequence(), "expectedType", project.GrantMemberChangedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-YX5Tk", "reduce.wrong.event.type")
	}
	return reduceMemberChanged(
		*member.NewMemberChangedEvent(&e.BaseEvent, e.UserID, e.Roles...),
		withMemberCond(ProjectGrantMemberProjectIDCol, e.Aggregate().ID),
		withMemberCond(ProjectGrantMemberGrantIDCol, e.GrantID),
	)
}

func (p *ProjectGrantMemberProjection) reduceCascadeRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.GrantMemberCascadeRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-azx7K", "seq", event.Sequence(), "expectedType", project.GrantMemberCascadeRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-adnHG", "reduce.wrong.event.type")
	}
	return reduceMemberCascadeRemoved(
		*member.NewCascadeRemovedEvent(&e.BaseEvent, e.UserID),
		withMemberCond(ProjectGrantMemberProjectIDCol, e.Aggregate().ID),
		withMemberCond(ProjectGrantMemberGrantIDCol, e.GrantID),
	)
}

func (p *ProjectGrantMemberProjection) reduceRemoved(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*project.GrantMemberRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-6Z4dH", "seq", event.Sequence(), "expectedType", project.GrantMemberRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-MGNnA", "reduce.wrong.event.type")
	}
	return reduceMemberRemoved(
		*member.NewRemovedEvent(&e.BaseEvent, e.UserID),
		withMemberCond(ProjectGrantMemberProjectIDCol, e.Aggregate().ID),
		withMemberCond(ProjectGrantMemberGrantIDCol, e.GrantID),
	)
}
