package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/deviceauth"
	"github.com/zitadel/zitadel/internal/repository/user"
)

const (
	DeviceAuthProjectionTable = "projections.device_authorizations"

	DeviceAuthColumnID         = "id"
	DeviceAuthColumnClientID   = "client_id"
	DeviceAuthColumnDeviceCode = "device_code"
	DeviceAuthColumnUserCode   = "user_code"
	DeviceAuthColumnExpires    = "expires"
	DeviceAuthColumnState      = "state"
	DeviceAuthColumnSubject    = "subject"

	DeviceAuthCreationDate = "creation_date"
	DeviceAuthChangeDate   = "change_date"
	DeviceAuthSequence     = "sequence"
	DeviceAuthInstanceID   = "instance_id"
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
			crdb.NewColumn(DeviceAuthColumnState, crdb.ColumnTypeEnum, crdb.Default(domain.DeviceAuthStateInitiated)),
			crdb.NewColumn(DeviceAuthColumnSubject, crdb.ColumnTypeText, crdb.Default("")),
			crdb.NewColumn(DeviceAuthCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(DeviceAuthChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(DeviceAuthSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(DeviceAuthInstanceID, crdb.ColumnTypeText),
		},
			crdb.NewPrimaryKey(DeviceAuthInstanceID, DeviceAuthColumnID),
			crdb.WithIndex(crdb.NewIndex("user_code", []string{DeviceAuthInstanceID, DeviceAuthColumnUserCode})),
			crdb.WithIndex(crdb.NewIndex("device_code", []string{DeviceAuthInstanceID, DeviceAuthColumnClientID, DeviceAuthColumnDeviceCode})),
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
					Event:  deviceauth.DeniedEventType,
					Reduce: p.reduceDenied,
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
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-< TODO: CODE >", "reduce.wrong.event.type %s", user.PersonalAccessTokenAddedType)
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(DeviceAuthInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(DeviceAuthColumnID, e.Aggregate().ID),
			handler.NewCol(DeviceAuthColumnClientID, e.ClientID),
			handler.NewCol(DeviceAuthColumnDeviceCode, e.DeviceCode),
			handler.NewCol(DeviceAuthColumnUserCode, e.UserCode),
			handler.NewCol(DeviceAuthColumnExpires, e.Expires),
		},
	), nil
}

func (p *deviceAuthProjection) reduceAppoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*deviceauth.ApprovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-< TODO: CODE >", "reduce.wrong.event.type %s", user.PersonalAccessTokenAddedType)
	}
	return crdb.NewUpdateStatement(e,
		[]handler.Column{
			handler.NewCol(DeviceAuthColumnSubject, e.Subject),
		},
		[]handler.Condition{
			handler.NewCond(DeviceAuthInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(DeviceAuthColumnID, e.Aggregate().ID),
		},
	), nil
}

func (p *deviceAuthProjection) reduceDenied(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*deviceauth.ApprovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-< TODO: CODE >", "reduce.wrong.event.type %s", user.PersonalAccessTokenAddedType)
	}
	return crdb.NewUpdateStatement(e,
		[]handler.Column{
			handler.NewCol(DeviceAuthColumnState, domain.DeviceAuthStateUserDenied),
			handler.NewCol(DeviceAuthColumnSubject, e.Subject),
		},
		[]handler.Condition{
			handler.NewCond(DeviceAuthInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(DeviceAuthColumnID, e.Aggregate().ID),
		},
	), nil
}

func (p *deviceAuthProjection) reduceRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*deviceauth.RemovedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-< TODO: CODE >", "reduce.wrong.event.type %s", user.PersonalAccessTokenAddedType)
	}
	return crdb.NewDeleteStatement(e,
		[]handler.Condition{
			handler.NewCond(DeviceAuthInstanceID, e.Aggregate().InstanceID),
			handler.NewCond(DeviceAuthColumnID, e.Aggregate().ID),
		},
	), nil
}
