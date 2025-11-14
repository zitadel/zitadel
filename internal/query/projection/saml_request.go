package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/samlrequest"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	SamlRequestsProjectionTable = "projections.saml_requests"

	SamlRequestColumnID            = "id"
	SamlRequestColumnCreationDate  = "creation_date"
	SamlRequestColumnChangeDate    = "change_date"
	SamlRequestColumnSequence      = "sequence"
	SamlRequestColumnResourceOwner = "resource_owner"
	SamlRequestColumnInstanceID    = "instance_id"
	SamlRequestColumnLoginClient   = "login_client"
	SamlRequestColumnIssuer        = "issuer"
	SamlRequestColumnACS           = "acs"
	SamlRequestColumnRelayState    = "relay_state"
	SamlRequestColumnBinding       = "binding"
)

type samlRequestProjection struct{}

// Name implements handler.Projection.
func (*samlRequestProjection) Name() string {
	return SamlRequestsProjectionTable
}

func newSamlRequestProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(samlRequestProjection))
}

func (*samlRequestProjection) Init() *old_handler.Check {
	return handler.NewMultiTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(SamlRequestColumnID, handler.ColumnTypeText),
			handler.NewColumn(SamlRequestColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(SamlRequestColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(SamlRequestColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(SamlRequestColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(SamlRequestColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(SamlRequestColumnLoginClient, handler.ColumnTypeText),
			handler.NewColumn(SamlRequestColumnIssuer, handler.ColumnTypeText),
			handler.NewColumn(SamlRequestColumnACS, handler.ColumnTypeText),
			handler.NewColumn(SamlRequestColumnRelayState, handler.ColumnTypeText),
			handler.NewColumn(SamlRequestColumnBinding, handler.ColumnTypeText),
		},
			handler.NewPrimaryKey(SamlRequestColumnInstanceID, SamlRequestColumnID),
		),
	)
}

func (p *samlRequestProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: samlrequest.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  samlrequest.AddedType,
					Reduce: p.reduceSamlRequestAdded,
				},
				{
					Event:  samlrequest.SucceededType,
					Reduce: p.reduceSamlRequestEnded,
				},
				{
					Event:  samlrequest.FailedType,
					Reduce: p.reduceSamlRequestEnded,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(SamlRequestColumnInstanceID),
				},
			},
		},
	}
}

func (p *samlRequestProjection) reduceSamlRequestAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*samlrequest.AddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-Sfwfa", "reduce.wrong.event.type %s", samlrequest.AddedType)
	}

	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(SamlRequestColumnID, e.Aggregate().ID),
			handler.NewCol(SamlRequestColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(SamlRequestColumnCreationDate, e.CreationDate()),
			handler.NewCol(SamlRequestColumnChangeDate, e.CreationDate()),
			handler.NewCol(SamlRequestColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(SamlRequestColumnSequence, e.Sequence()),
			handler.NewCol(SamlRequestColumnLoginClient, e.LoginClient),
			handler.NewCol(SamlRequestColumnIssuer, e.Issuer),
			handler.NewCol(SamlRequestColumnACS, e.ACSURL),
			handler.NewCol(SamlRequestColumnRelayState, e.RelayState),
			handler.NewCol(SamlRequestColumnBinding, e.Binding),
		},
	), nil
}

func (p *samlRequestProjection) reduceSamlRequestEnded(event eventstore.Event) (*handler.Statement, error) {
	switch event.(type) {
	case *samlrequest.SucceededEvent,
		*samlrequest.FailedEvent:
		break
	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-ASF3h", "reduce.wrong.event.type %s", []eventstore.EventType{samlrequest.SucceededType, samlrequest.FailedType})
	}

	return handler.NewDeleteStatement(
		event,
		[]handler.Condition{
			handler.NewCond(SamlRequestColumnID, event.Aggregate().ID),
			handler.NewCond(SamlRequestColumnInstanceID, event.Aggregate().InstanceID),
		},
	), nil
}
