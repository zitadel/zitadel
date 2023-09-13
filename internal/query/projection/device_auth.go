package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
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

type deviceAuthProjection struct{}

func newDeviceAuthProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(deviceAuthProjection))
}

func (*deviceAuthProjection) Name() string {
	return DeviceAuthProjectionTable
}

func (*deviceAuthProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(DeviceAuthColumnID, handler.ColumnTypeText),
			handler.NewColumn(DeviceAuthColumnClientID, handler.ColumnTypeText),
			handler.NewColumn(DeviceAuthColumnDeviceCode, handler.ColumnTypeText),
			handler.NewColumn(DeviceAuthColumnUserCode, handler.ColumnTypeText),
			handler.NewColumn(DeviceAuthColumnExpires, handler.ColumnTypeTimestamp),
			handler.NewColumn(DeviceAuthColumnScopes, handler.ColumnTypeTextArray),
			handler.NewColumn(DeviceAuthColumnState, handler.ColumnTypeEnum, handler.Default(domain.DeviceAuthStateInitiated)),
			handler.NewColumn(DeviceAuthColumnSubject, handler.ColumnTypeText, handler.Default("")),
			handler.NewColumn(DeviceAuthColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(DeviceAuthColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(DeviceAuthColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(DeviceAuthColumnInstanceID, handler.ColumnTypeText),
		},
			handler.NewPrimaryKey(DeviceAuthColumnInstanceID, DeviceAuthColumnID),
			handler.WithIndex(handler.NewIndex("user_code", []string{DeviceAuthColumnInstanceID, DeviceAuthColumnUserCode})),
			handler.WithIndex(handler.NewIndex("device_code", []string{DeviceAuthColumnInstanceID, DeviceAuthColumnClientID, DeviceAuthColumnDeviceCode})),
		),
	)
}

func (p *deviceAuthProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: deviceauth.AggregateType,
			EventReducers: []handler.EventReducer{
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
	return handler.NewCreateStatement(
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
	return handler.NewUpdateStatement(e,
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
	return handler.NewUpdateStatement(e,
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
	return handler.NewDeleteStatement(e,
		[]handler.Condition{
			handler.NewCond(DeviceAuthColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(DeviceAuthColumnID, e.Aggregate().ID),
		},
	), nil
}
