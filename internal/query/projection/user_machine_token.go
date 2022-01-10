package projection

import (
	"context"

	"github.com/caos/logging"
	"github.com/lib/pq"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/user"
)

type MachineTokenProjection struct {
	crdb.StatementHandler
}

const (
	MachineTokenProjectionTable = "zitadel.projections.machine_tokens"
)

func NewMachineTokenProjection(ctx context.Context, config crdb.StatementHandlerConfig) *MachineTokenProjection {
	p := &MachineTokenProjection{}
	config.ProjectionName = MachineTokenProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *MachineTokenProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: user.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  user.MachineTokenAddedType,
					Reduce: p.reduceMachineTokenAdded,
				},
				{
					Event:  user.MachineTokenRemovedType,
					Reduce: p.reduceMachineTokenRemoved,
				},
				{
					Event:  user.UserRemovedType,
					Reduce: p.reduceUserRemoved,
				},
			},
		},
	}
}

const (
	MachineTokenColumnID            = "id"
	MachineTokenColumnCreationDate  = "creation_date"
	MachineTokenColumnChangeDate    = "change_date"
	MachineTokenColumnResourceOwner = "resource_owner"
	MachineTokenColumnSequence      = "sequence"
	MachineTokenColumnUserID        = "user_id"
	MachineTokenColumnExpiration    = "expiration"
	MachineTokenColumnScopes        = "scopes"
)

func (p *MachineTokenProjection) reduceMachineTokenAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MachineTokenAddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Dbfg2", "seq", event.Sequence(), "expectedType", user.MachineTokenAddedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-DVgf7", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(MachineTokenColumnID, e.TokenID),
			handler.NewCol(MachineTokenColumnCreationDate, e.CreationDate()),
			handler.NewCol(MachineTokenColumnChangeDate, e.CreationDate()),
			handler.NewCol(MachineTokenColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(MachineTokenColumnSequence, e.Sequence()),
			handler.NewCol(MachineTokenColumnUserID, e.Aggregate().ID),
			handler.NewCol(MachineTokenColumnExpiration, e.Expiration),
			handler.NewCol(MachineTokenColumnScopes, pq.StringArray(e.Scopes)),
		},
	), nil
}

func (p *MachineTokenProjection) reduceMachineTokenRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.MachineTokenRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-Edf32", "seq", event.Sequence(), "expectedType", user.MachineTokenRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-g7u3F", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(MachineTokenColumnID, e.TokenID),
		},
	), nil
}

func (p *MachineTokenProjection) reduceUserRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*user.UserRemovedEvent)
	if !ok {
		logging.LogWithFields("HANDL-GEg43", "seq", event.Sequence(), "expectedType", user.UserRemovedType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-Dff3h", "reduce.wrong.event.type")
	}
	return crdb.NewDeleteStatement(
		e,
		[]handler.Condition{
			handler.NewCond(MachineTokenColumnUserID, e.Aggregate().ID),
		},
	), nil
}
