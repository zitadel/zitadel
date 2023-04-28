package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/deviceauth"
)

const (
	DeviceAuthProjectionTable = "projections.device_authorizations"

	DeviceAuthColumnID         = "id"
	DeviceAuthColumnClientID   = "client_id"
	DeviceAuthColumnDeviceCode = "device_code"
	DeviceAuthColumnUserCode   = "user_code"
	DeviceAuthColumnExpires    = "expires"
	DeviceAuthColumnScopes     = "scopes"
	DeviceAuthColumnState      = "state"
	DeviceAuthColumnSubject    = "subject"

	DeviceAuthColumnCreationDate = "creation_date"
	DeviceAuthColumnChangeDate   = "change_date"
	DeviceAuthColumnSequence     = "sequence"
	DeviceAuthColumnInstanceID   = "instance_id"
)

type deviceAuthProjection struct {
	crdb.StatementHandler
}

func newDeviceAuthProjection(ctx context.Context, config crdb.StatementHandlerConfig) *deviceAuthProjection {
	p := new(deviceAuthProjection)
	config.ProjectionName = DeviceAuthProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(DeviceAuthColumnID, crdb.ColumnTypeText),
			crdb.NewColumn(DeviceAuthColumnClientID, crdb.ColumnTypeText),
			crdb.NewColumn(DeviceAuthColumnDeviceCode, crdb.ColumnTypeText),
			crdb.NewColumn(DeviceAuthColumnUserCode, crdb.ColumnTypeText),
			crdb.NewColumn(DeviceAuthColumnExpires, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(DeviceAuthColumnScopes, crdb.ColumnTypeTextArray),
			crdb.NewColumn(DeviceAuthColumnState, crdb.ColumnTypeEnum, crdb.Default(domain.DeviceAuthStateInitiated)),
			crdb.NewColumn(DeviceAuthColumnSubject, crdb.ColumnTypeText, crdb.Default("")),
			crdb.NewColumn(DeviceAuthColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(DeviceAuthColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(DeviceAuthColumnSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(DeviceAuthColumnInstanceID, crdb.ColumnTypeText),
		},
			crdb.NewPrimaryKey(DeviceAuthColumnInstanceID, DeviceAuthColumnID),
			crdb.WithIndex(crdb.NewIndex("user_code", []string{DeviceAuthColumnInstanceID, DeviceAuthColumnUserCode})),
			crdb.WithIndex(crdb.NewIndex("device_code", []string{DeviceAuthColumnInstanceID, DeviceAuthColumnClientID, DeviceAuthColumnDeviceCode})),
		),
	)

	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *deviceAuthProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: deviceauth.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  deviceauth.AddedEventType,
					Reduce: p.reduceAdded,
				},
				{
					Event:  deviceauth.ApprovedEventType,
					Reduce: p.reduceAppoved,
				},
				{
					Event:  deviceauth.CanceledEventType,
					Reduce: p.reduceCanceled,
				},
				{
					Event:  deviceauth.RemovedEventType,
					Reduce: p.reduceRemoved,
				},
			},
		},
	}
}

func (p *deviceAuthProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*deviceauth.AddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-chu6O", "reduce.wrong.event.type %T != %s", event, deviceauth.AddedEventType)
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(DeviceAuthColumnID, e.Aggregate().ID),
			handler.NewCol(DeviceAuthColumnClientID, e.ClientID),
			handler.NewCol(DeviceAuthColumnDeviceCode, e.DeviceCode),
			handler.NewCol(DeviceAuthColumnUserCode, e.UserCode),
			handler.NewCol(DeviceAuthColumnExpires, e.Expires),
			handler.NewCol(DeviceAuthColumnScopes, e.Scopes),
			handler.NewCol(DeviceAuthColumnCreationDate, e.CreationDate()),
			handler.NewCol(DeviceAuthColumnChangeDate, e.CreationDate()),
			handler.NewCol(DeviceAuthColumnSequence, e.Sequence()),
			handler.NewCol(DeviceAuthColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

func (p *deviceAuthProjection) reduceAppoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*deviceauth.ApprovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-kei0A", "reduce.wrong.event.type %T != %s", event, deviceauth.ApprovedEventType)
	}
	return crdb.NewUpdateStatement(e,
		[]handler.Column{
			handler.NewCol(DeviceAuthColumnState, domain.DeviceAuthStateApproved),
			handler.NewCol(DeviceAuthColumnSubject, e.Subject),
			handler.NewCol(DeviceAuthColumnChangeDate, e.CreationDate()),
			handler.NewCol(DeviceAuthColumnSequence, e.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(DeviceAuthColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(DeviceAuthColumnID, e.Aggregate().ID),
		},
	), nil
}

func (p *deviceAuthProjection) reduceCanceled(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*deviceauth.CanceledEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-eeS8d", "reduce.wrong.event.type %T != %s", event, deviceauth.CanceledEventType)
	}
	return crdb.NewUpdateStatement(e,
		[]handler.Column{
			handler.NewCol(DeviceAuthColumnState, e.Reason.State()),
			handler.NewCol(DeviceAuthColumnChangeDate, e.CreationDate()),
			handler.NewCol(DeviceAuthColumnSequence, e.Sequence()),
		},
		[]handler.Condition{
			handler.NewCond(DeviceAuthColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(DeviceAuthColumnID, e.Aggregate().ID),
		},
	), nil
}

func (p *deviceAuthProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*deviceauth.RemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-AJi1u", "reduce.wrong.event.type %T != %s", event, deviceauth.RemovedEventType)
	}
	return crdb.NewDeleteStatement(e,
		[]handler.Condition{
			handler.NewCond(DeviceAuthColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(DeviceAuthColumnID, e.Aggregate().ID),
		},
	), nil
}
