package projection

import (
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

	MilestoneColumnInstanceID     = "instance_id"
	MilestoneColumnType           = "type"
	MilestoneColumnPrimaryDomain  = "primary_domain"
	MilestoneColumnReachedDate    = "reached_date"
	MilestoneColumnPushedDate     = "last_pushed_date"
	MilestoneColumnIgnoredProject = "ignore_project"
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
			crdb.NewColumn(MilestoneColumnType, crdb.ColumnTypeEnum),
			crdb.NewColumn(MilestoneColumnReachedDate, crdb.ColumnTypeTimestamp, crdb.Nullable()),
			crdb.NewColumn(MilestoneColumnPushedDate, crdb.ColumnTypeTimestamp, crdb.Nullable()),
			crdb.NewColumn(MilestoneColumnPrimaryDomain, crdb.ColumnTypeText, crdb.Nullable()),
			crdb.NewColumn(MilestoneColumnIgnoredProject, crdb.ColumnTypeText, crdb.Nullable()),
		},
			crdb.NewPrimaryKey(MilestoneColumnInstanceID, MilestoneColumnType),
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
					Reduce: p.milestoneReached(milestone.InstanceDeleted),
				},
			},
		},
		{
			Aggregate: project.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  project.ProjectAddedType,
					Reduce: p.reduceProjectAdded,
				},
				{
					Event:  project.ApplicationAddedType,
					Reduce: p.milestoneReached(milestone.ApplicationCreated),
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
		{
			Aggregate: milestone.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  milestone.PushedEventType,
					Reduce: p.reducePushed,
				},
			},
		},
	}
}

func (p *milestoneProjection) milestoneReached(msType milestone.Type) func(event eventstore.Event) (*handler.Statement, error) {
	return func(event eventstore.Event) (*handler.Statement, error) {
		printEvent(event)
		if isSystemEvent(event) {
			return crdb.NewNoOpStatement(event), nil
		}
		return crdb.NewUpdateStatement(event, []handler.Column{
			handler.NewCol(MilestoneColumnReachedDate, event.CreationDate()),
		},
			[]handler.Condition{
				handler.NewCond(MilestoneColumnInstanceID, event.Aggregate().InstanceID),
				handler.NewCond(MilestoneColumnType, msType),
				crdb.NewIsNullCond(MilestoneColumnReachedDate),
			}), nil
	}
}

func (p *milestoneProjection) reducePushed(event eventstore.Event) (*handler.Statement, error) {
	printEvent(event)
	e, ok := event.(*milestone.PushedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-XJGXK", "reduce.wrong.event.type %s", milestone.PushedEventType)
	}
	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(MilestoneColumnPushedDate, event.CreationDate()),
		},
		[]handler.Condition{
			handler.NewCond(MilestoneColumnInstanceID, event.Aggregate().InstanceID),
			handler.NewCond(MilestoneColumnType, e.MilestoneType),
		},
	), nil
}

func (p *milestoneProjection) reduceInstanceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.InstanceAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-JbHGS", "reduce.wrong.event.type %s", instance.InstanceAddedEventType)
	}
	allTypes := milestone.AllTypes()
	statements := make([]func(eventstore.Event) crdb.Exec, 0, len(allTypes))
	for _, msType := range allTypes {
		createColumns := []handler.Column{
			handler.NewCol(MilestoneColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(MilestoneColumnType, msType),
		}
		if msType == milestone.InstanceCreated {
			createColumns = append(createColumns, handler.NewCol(MilestoneColumnReachedDate, event.CreationDate()))
		}
		statements = append(statements, crdb.AddCreateStatement(createColumns))
	}
	return crdb.NewMultiStatement(e, statements...), nil
}

func (p *milestoneProjection) reduceInstanceDomainPrimarySet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DomainPrimarySetEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Sfrgf", "reduce.wrong.event.type %s", instance.InstanceDomainPrimarySetEventType)
	}
	return crdb.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(MilestoneColumnPrimaryDomain, e.Domain),
		},
		[]handler.Condition{
			handler.NewCond(MilestoneColumnInstanceID, e.Aggregate().InstanceID),
			crdb.NewIsNullCond(MilestoneColumnPushedDate),
		},
	), nil
}

func (p *milestoneProjection) reduceProjectAdded(event eventstore.Event) (*handler.Statement, error) {
	if !isSystemEvent(event) {
		return p.milestoneReached(milestone.ProjectCreated)(event)
	}
	return crdb.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewCol(MilestoneColumnIgnoredProject, event.Aggregate().ID),
		},
		[]handler.Condition{
			handler.NewCond(MilestoneColumnInstanceID, event.Aggregate().InstanceID),
			handler.NewCond(MilestoneColumnType, milestone.AuthenticationSucceededOnApplication),
			crdb.NewIsNullCond(MilestoneColumnReachedDate),
		},
	), nil
}

func (p *milestoneProjection) reduceUserTokenAdded(event eventstore.Event) (*handler.Statement, error) {
	printEvent(event)
	e, ok := event.(*user.UserTokenAddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-3xhJ7", "reduce.wrong.event.type %s", user.UserTokenAddedType)
	}
	return crdb.NewMultiStatement(
		e,
		crdb.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(MilestoneColumnReachedDate, event.CreationDate()),
			},
			[]handler.Condition{
				handler.NewCond(MilestoneColumnInstanceID, event.Aggregate().InstanceID),
				handler.NewCond(MilestoneColumnType, milestone.AuthenticationSucceededOnInstance),
				crdb.NewIsNullCond(MilestoneColumnReachedDate),
			},
		),
		/*		crdb.AddUpdateStatement(
				[]handler.Column{
					handler.NewCol(MilestoneColumnReachedDate, event.CreationDate()),
				},
				[]handler.Condition{
					handler.NewCond(MilestoneColumnInstanceID, event.Aggregate().InstanceID),
					handler.NewCond(MilestoneColumnType, milestone.AuthenticationSucceededOnApplication),
					crdb.NewExpressionCond(func(param string) (clause string, needsParam bool) {
						return fmt.Sprintf("%s")
					}),
					crdb.NewIsNullCond(MilestoneColumnReachedDate),
				},
			),*/
	), nil
}

func isSystemEvent(event eventstore.Event) bool {
	return event.EditorUser() == "" || event.EditorService() == ""
}

// TODO: Remove
func printEvent(event eventstore.Event) {
	bytes, err := json.MarshalIndent(event, "", "    ")
	if err != nil {
		panic(err)
	}
	fmt.Println(event.Type(), string(bytes))
}
