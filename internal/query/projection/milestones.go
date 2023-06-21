package projection

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"

	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"

	"github.com/zitadel/zitadel/internal/repository/milestone"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
)

const (
	MilestonesProjectionTable = "projections.milestones"

	MilestoneColumnInstanceID    = "instance_id"
	MilestoneColumnMilestoneType = "milestone_type"
	MilestoneColumnReachedAt     = "reached_at"
	MilestoneColumnPushedAt      = "pushed_at"
	MilestoneColumnPrimaryDomain = "primary_domain"
)

type milestoneProjection struct {
	crdb.StatementHandler
}

func newMilestoneInstanceProjection(ctx context.Context, config crdb.StatementHandlerConfig) *milestoneProjection {
	p := new(milestoneProjection)
	config.ProjectionName = MilestonesProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewMultiTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(MilestoneColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(MilestoneColumnMilestoneType, crdb.ColumnTypeEnum),
			crdb.NewColumn(MilestoneColumnReachedAt, crdb.ColumnTypeTimestamp, crdb.Nullable()),
			crdb.NewColumn(MilestoneColumnPushedAt, crdb.ColumnTypeTimestamp, crdb.Nullable()),
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
					Reduce: p.reduceInstanceRemoved,
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
					Reduce: p.reduceApplicationAdded,
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
					Reduce: p.milestonePushed,
				},
			},
		},
	}
}

func (p *milestoneProjection) reduceInstanceDomainPrimarySet(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.DomainPrimarySetEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-Sfrgf", "reduce.wrong.event.type %s", instance.InstanceDomainPrimarySetEventType)
	}

	var statements []func(eventstore.Event) crdb.Exec
	for _, ms := range milestone.All() {
		statements = append(statements, crdb.AddUpsertStatement(
			[]handler.Column{
				handler.NewCol(MilestoneColumnInstanceID, nil),
				handler.NewCol(MilestoneColumnMilestoneType, nil),
			},
			[]handler.Column{
				handler.NewCol(MilestoneColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(MilestoneColumnMilestoneType, ms),
				handler.NewCol(MilestoneColumnPrimaryDomain, e.Domain),
			},
		))
	}

	return crdb.NewMultiStatement(e, statements...), nil
}

func (p *milestoneProjection) reduceInstanceAdded(event eventstore.Event) (*handler.Statement, error) {
	printEvent(event)

	return crdb.NewNoOpStatement(event), nil
}

func (p *milestoneProjection) reduceProjectAdded(event eventstore.Event) (*handler.Statement, error) {
	printEvent(event)
	// ignore instance.ProjectSetEventType
	return crdb.NewNoOpStatement(event), nil
}

func (p *milestoneProjection) reduceApplicationAdded(event eventstore.Event) (*handler.Statement, error) {
	printEvent(event)
	return crdb.NewNoOpStatement(event), nil
}

func (p *milestoneProjection) reduceUserTokenAdded(event eventstore.Event) (*handler.Statement, error) {
	printEvent(event)
	return crdb.NewNoOpStatement(event), nil
}

func (p *milestoneProjection) reduceInstanceRemoved(event eventstore.Event) (*handler.Statement, error) {
	printEvent(event)
	return crdb.NewNoOpStatement(event), nil
}

func (p *milestoneProjection) milestonePushed(event eventstore.Event) (*handler.Statement, error) {
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
