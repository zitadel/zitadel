package projection

import (
	"context"
	"strconv"

	internal_authz "github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/milestone"
	"github.com/zitadel/zitadel/internal/repository/project"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	MilestonesProjectionTable = "projections.milestones"

	MilestoneColumnInstanceID      = "instance_id"
	MilestoneColumnType            = "type"
	MilestoneColumnPrimaryDomain   = "primary_domain"
	MilestoneColumnReachedDate     = "reached_date"
	MilestoneColumnPushedDate      = "last_pushed_date"
	MilestoneColumnIgnoreClientIDs = "ignore_client_ids"
)

type milestoneProjection struct {
	systemUsers map[string]*internal_authz.SystemAPIUser
}

func newMilestoneProjection(ctx context.Context, config handler.Config, systemUsers map[string]*internal_authz.SystemAPIUser) *handler.Handler {
	return handler.NewHandler(ctx, &config, &milestoneProjection{systemUsers: systemUsers})
}

func (*milestoneProjection) Name() string {
	return MilestonesProjectionTable
}

func (*milestoneProjection) Init() *old_handler.Check {
	return handler.NewMultiTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(MilestoneColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(MilestoneColumnType, handler.ColumnTypeEnum),
			handler.NewColumn(MilestoneColumnReachedDate, handler.ColumnTypeTimestamp, handler.Nullable()),
			handler.NewColumn(MilestoneColumnPushedDate, handler.ColumnTypeTimestamp, handler.Nullable()),
			handler.NewColumn(MilestoneColumnPrimaryDomain, handler.ColumnTypeText, handler.Nullable()),
			handler.NewColumn(MilestoneColumnIgnoreClientIDs, handler.ColumnTypeTextArray, handler.Nullable()),
		},
			handler.NewPrimaryKey(MilestoneColumnInstanceID, MilestoneColumnType),
		),
	)
}

// Reducers implements handler.Projection.
func (p *milestoneProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
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
			EventReducers: []handler.EventReducer{
				{
					Event:  project.ProjectAddedType,
					Reduce: p.reduceProjectAdded,
				},
				{
					Event:  project.ApplicationAddedType,
					Reduce: p.reduceApplicationAdded,
				},
				{
					Event:  project.OIDCConfigAddedType,
					Reduce: p.reduceOIDCConfigAdded,
				},
				{
					Event:  project.APIConfigAddedType,
					Reduce: p.reduceAPIConfigAdded,
				},
			},
		},
		{
			Aggregate: user.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					// user.UserTokenAddedType is not emitted on creation of personal access tokens
					// PATs have no effect on milestone.AuthenticationSucceededOnApplication or milestone.AuthenticationSucceededOnInstance
					Event:  user.UserTokenAddedType,
					Reduce: p.reduceUserTokenAdded,
				},
			},
		},
		{
			Aggregate: milestone.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  milestone.PushedEventType,
					Reduce: p.reduceMilestonePushed,
				},
			},
		},
	}
}

func (p *milestoneProjection) reduceInstanceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.InstanceAddedEvent](event)
	if err != nil {
		return nil, err
	}
	allTypes := milestone.AllTypes()
	statements := make([]func(eventstore.Event) handler.Exec, 0, len(allTypes))
	for _, msType := range allTypes {
		createColumns := []handler.Column{
			handler.NewCol(MilestoneColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(MilestoneColumnType, msType),
		}
		if msType == milestone.InstanceCreated {
			createColumns = append(createColumns, handler.NewCol(MilestoneColumnReachedDate, event.CreatedAt()))
		}
		statements = append(statements, handler.AddCreateStatement(createColumns))
	}
	return handler.NewMultiStatement(e, statements...), nil
}

func (p *milestoneProjection) reduceInstanceDomainPrimarySet(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*instance.DomainPrimarySetEvent](event)
	if err != nil {
		return nil, err
	}
	return handler.NewUpdateStatement(
		e,
		[]handler.Column{
			handler.NewCol(MilestoneColumnPrimaryDomain, e.Domain),
		},
		[]handler.Condition{
			handler.NewCond(MilestoneColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewIsNullCond(MilestoneColumnPushedDate),
		},
	), nil
}

func (p *milestoneProjection) reduceProjectAdded(event eventstore.Event) (*handler.Statement, error) {
	if _, err := assertEvent[*project.ProjectAddedEvent](event); err != nil {
		return nil, err
	}
	return p.reduceReachedIfUserEventFunc(milestone.ProjectCreated)(event)
}

func (p *milestoneProjection) reduceApplicationAdded(event eventstore.Event) (*handler.Statement, error) {
	if _, err := assertEvent[*project.ApplicationAddedEvent](event); err != nil {
		return nil, err
	}
	return p.reduceReachedIfUserEventFunc(milestone.ApplicationCreated)(event)
}

func (p *milestoneProjection) reduceOIDCConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*project.OIDCConfigAddedEvent](event)
	if err != nil {
		return nil, err
	}
	return p.reduceAppConfigAdded(e, e.ClientID)
}

func (p *milestoneProjection) reduceAPIConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*project.APIConfigAddedEvent](event)
	if err != nil {
		return nil, err
	}
	return p.reduceAppConfigAdded(e, e.ClientID)
}

func (p *milestoneProjection) reduceUserTokenAdded(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*user.UserTokenAddedEvent](event)
	if err != nil {
		return nil, err
	}
	statements := []func(eventstore.Event) handler.Exec{
		handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(MilestoneColumnReachedDate, event.CreatedAt()),
			},
			[]handler.Condition{
				handler.NewCond(MilestoneColumnInstanceID, event.Aggregate().InstanceID),
				handler.NewCond(MilestoneColumnType, milestone.AuthenticationSucceededOnInstance),
				handler.NewIsNullCond(MilestoneColumnReachedDate),
			},
		),
	}
	// We ignore authentications without app, for example JWT profile or PAT
	if e.ApplicationID != "" {
		statements = append(statements, handler.AddUpdateStatement(
			[]handler.Column{
				handler.NewCol(MilestoneColumnReachedDate, event.CreatedAt()),
			},
			[]handler.Condition{
				handler.NewCond(MilestoneColumnInstanceID, event.Aggregate().InstanceID),
				handler.NewCond(MilestoneColumnType, milestone.AuthenticationSucceededOnApplication),
				handler.Not(handler.NewTextArrayContainsCond(MilestoneColumnIgnoreClientIDs, e.ApplicationID)),
				handler.NewIsNullCond(MilestoneColumnReachedDate),
			},
		))
	}
	return handler.NewMultiStatement(e, statements...), nil
}

func (p *milestoneProjection) reduceInstanceRemoved(event eventstore.Event) (*handler.Statement, error) {
	if _, err := assertEvent[*instance.InstanceRemovedEvent](event); err != nil {
		return nil, err
	}
	return p.reduceReachedFunc(milestone.InstanceDeleted)(event)
}

func (p *milestoneProjection) reduceMilestonePushed(event eventstore.Event) (*handler.Statement, error) {
	e, err := assertEvent[*milestone.PushedEvent](event)
	if err != nil {
		return nil, err
	}
	if e.MilestoneType != milestone.InstanceDeleted {
		return handler.NewUpdateStatement(
			event,
			[]handler.Column{
				handler.NewCol(MilestoneColumnPushedDate, event.CreatedAt()),
			},
			[]handler.Condition{
				handler.NewCond(MilestoneColumnInstanceID, event.Aggregate().InstanceID),
				handler.NewCond(MilestoneColumnType, e.MilestoneType),
			},
		), nil
	}
	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(MilestoneColumnInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}

func (p *milestoneProjection) reduceReachedIfUserEventFunc(msType milestone.Type) func(event eventstore.Event) (*handler.Statement, error) {
	return func(event eventstore.Event) (*handler.Statement, error) {
		if p.isSystemEvent(event) {
			return handler.NewNoOpStatement(event), nil
		}
		return p.reduceReachedFunc(msType)(event)
	}
}

func (p *milestoneProjection) reduceReachedFunc(msType milestone.Type) func(event eventstore.Event) (*handler.Statement, error) {
	return func(event eventstore.Event) (*handler.Statement, error) {
		return handler.NewUpdateStatement(event, []handler.Column{
			handler.NewCol(MilestoneColumnReachedDate, event.CreatedAt()),
		},
			[]handler.Condition{
				handler.NewCond(MilestoneColumnInstanceID, event.Aggregate().InstanceID),
				handler.NewCond(MilestoneColumnType, msType),
				handler.NewIsNullCond(MilestoneColumnReachedDate),
			}), nil
	}
}

func (p *milestoneProjection) reduceAppConfigAdded(event eventstore.Event, clientID string) (*handler.Statement, error) {
	if !p.isSystemEvent(event) {
		return handler.NewNoOpStatement(event), nil
	}
	return handler.NewUpdateStatement(
		event,
		[]handler.Column{
			handler.NewArrayAppendCol(MilestoneColumnIgnoreClientIDs, clientID),
		},
		[]handler.Condition{
			handler.NewCond(MilestoneColumnInstanceID, event.Aggregate().InstanceID),
			handler.NewCond(MilestoneColumnType, milestone.AuthenticationSucceededOnApplication),
			handler.NewIsNullCond(MilestoneColumnReachedDate),
		},
	), nil
}

func (p *milestoneProjection) isSystemEvent(event eventstore.Event) bool {
	if userId, err := strconv.Atoi(event.Creator()); err == nil && userId > 0 {
		return false
	}

	// check if it is a hard coded event creator
	for _, creator := range []string{"", "system", "OIDC", "LOGIN", "SYSTEM"} {
		if creator == event.Creator() {
			return true
		}
	}

	_, ok := p.systemUsers[event.Creator()]
	return ok
}
