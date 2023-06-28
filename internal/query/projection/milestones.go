package projection

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/zitadel/zitadel/internal/repository/milestone"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	MilestonesProjectionTable = "projections.milestones"

	MilestoneColumnInstanceID    = "instance_id"
	MilestoneColumnMilestoneType = "milestone_type"
	MilestoneColumnPrimaryDomain = "primary_domain"
	MilestoneColumnReachedDate   = "reached_date"
	MilestoneColumnPushedDate    = "pushed_date"
)

type milestoneProjection struct {
	crdb.StatementHandler
}

func newMilestoneProjection(ctx context.Context, config crdb.StatementHandlerConfig) *milestoneProjection {
	p := new(milestoneProjection)
	config.ProjectionName = MilestonesProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewMultiTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(MilestoneColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(MilestoneColumnMilestoneType, crdb.ColumnTypeText),
			crdb.NewColumn(MilestoneColumnReachedDate, crdb.ColumnTypeTimestamp, crdb.Nullable()),
			crdb.NewColumn(MilestoneColumnPushedDate, crdb.ColumnTypeTimestamp, crdb.Nullable()),
			crdb.NewColumn(MilestoneColumnPrimaryDomain, crdb.ColumnTypeText, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(MilestoneColumnInstanceID, MilestoneColumnMilestoneType),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *milestoneProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.InstanceAddedEventType,
					Reduce: p.reduceInstanceAdded,
				},
				{
					Event:  instance.InstanceDomainPrimarySetEventType,
					Reduce: p.reduceInstanceDomainPrimarySet,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: p.milestoneReached(milestone.PushedInstanceDeletedEventType),
				},
			},
		},
		{
			Aggregate: project.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  project.ProjectAddedType,
					Reduce: p.milestoneReached(milestone.PushedProjectCreatedEventType),
				},
				{
					Event:  project.ApplicationAddedType,
					Reduce: p.milestoneReached(milestone.PushedApplicationCreatedEventType),
				},
			},
		},
		{
			Aggregate: user.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  user.UserTokenAddedType,
					Reduce: p.reduceUserTokenAdded,
				},
			},
		},
	}
}

func (p *milestoneProjection) reduceInstanceAdded(event eventstore.Event) (*handler.Statement, error) {
	printEvent(event)
	e, ok := event.(*instance.InstanceAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-JbHGS", "reduce.wrong.event.type %s", instance.InstanceAddedEventType)
	}
	allTypes := milestone.PushedEventTypes()
	statements := make([]func(eventstore.Event) crdb.Exec, 0, len(allTypes))
	for _, ms := range allTypes {
		createColumns := []handler.Column{
			handler.NewCol(MilestoneColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(MilestoneColumnMilestoneType, ms),
		}
		if ms == milestone.PushedInstanceCreatedEventType {
			createColumns = append(createColumns, handler.NewCol(MilestoneColumnReachedDate, event.CreationDate()))
		}
		statements = append(statements, crdb.AddCreateStatement(createColumns))
	}
	return crdb.NewMultiStatement(e, statements...), nil
}

func (p *milestoneProjection) reduceInstanceDomainPrimarySet(event eventstore.Event) (*handler.Statement, error) {
	printEvent(event)
	e, ok := event.(*instance.DomainPrimarySetEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Sfrgf", "reduce.wrong.event.type %s", instance.InstanceDomainPrimarySetEventType)
	}
	allTypes := milestone.PushedEventTypes()
	statements := make([]func(eventstore.Event) crdb.Exec, 0, len(allTypes))
	for _, ms := range allTypes {
		statements = append(statements, crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(MilestoneColumnPrimaryDomain, e.Domain),
			},
			[]handler.Condition{
				handler.NewCond(MilestoneColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCond(MilestoneColumnMilestoneType, ms),
				crdb.NewIsNullCond(MilestoneColumnPushedDate),
			},
		))
	}
	return crdb.NewMultiStatement(e, statements...), nil
}

func (p *milestoneProjection) milestoneReached(eventType milestone.PushedEventType) func(event eventstore.Event) (*handler.Statement, error) {
	return func(event eventstore.Event) (*handler.Statement, error) {
		printEvent(event)
		if event.EditorUser() == "" || event.EditorService() == "" {
			return crdb.NewNoOpStatement(event), nil
		}
		return crdb.NewUpdateStatement(
			event,
			[]handler.Column{
				handler.NewCol(MilestoneColumnReachedDate, event.CreationDate()),
			},
			[]handler.Condition{
				handler.NewCond(MilestoneColumnInstanceID, event.Aggregate().InstanceID),
				handler.NewCond(MilestoneColumnMilestoneType, eventType),
				crdb.NewIsNullCond(MilestoneColumnReachedDate),
				crdb.NewIsNullCond(MilestoneColumnPushedDate),
			},
		), nil
	}
}

func (p *milestoneProjection) reduceUserTokenAdded(event eventstore.Event) (*handler.Statement, error) {
	printEvent(event)
	return crdb.NewNoOpStatement(event), nil
}

func (p *milestoneProjection) reduceInstanceRemoved(event eventstore.Event) (*handler.Statement, error) {
	printEvent(event)
	return crdb.NewNoOpStatement(event), nil
}

func printEvent(event eventstore.Event) {
	var pretty bytes.Buffer
	if err := json.Indent(&pretty, event.DataAsBytes(), "", "    "); err != nil {
		panic(err)
	}
	fmt.Println(event.Type(), pretty.String())
}
