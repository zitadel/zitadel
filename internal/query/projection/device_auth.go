package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/deviceauth"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	DeviceAuthRequestProjectionTable = "projections.device_auth_requests2"

	DeviceAuthRequestColumnClientID     = "client_id"
	DeviceAuthRequestColumnDeviceCode   = "device_code"
	DeviceAuthRequestColumnUserCode     = "user_code"
	DeviceAuthRequestColumnScopes       = "scopes"
	DeviceAuthRequestColumnAudience     = "audience"
	DeviceAuthRequestColumnCreationDate = "creation_date"
	DeviceAuthRequestColumnChangeDate   = "change_date"
	DeviceAuthRequestColumnSequence     = "sequence"
	DeviceAuthRequestColumnInstanceID   = "instance_id"
)

// deviceAuthRequestProjection holds device authorization requests
// and makes them search-able by User Code.
// In principle the projected data is only needed during user login.
// Device Token logic uses the eventstore directly.
type deviceAuthRequestProjection struct{}

func newDeviceAuthProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(deviceAuthRequestProjection))
}

func (*deviceAuthRequestProjection) Name() string {
	return DeviceAuthRequestProjectionTable
}

func (*deviceAuthRequestProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(DeviceAuthRequestColumnClientID, handler.ColumnTypeText),
			handler.NewColumn(DeviceAuthRequestColumnDeviceCode, handler.ColumnTypeText),
			handler.NewColumn(DeviceAuthRequestColumnUserCode, handler.ColumnTypeText),
			handler.NewColumn(DeviceAuthRequestColumnScopes, handler.ColumnTypeTextArray, handler.Nullable()),
			handler.NewColumn(DeviceAuthRequestColumnAudience, handler.ColumnTypeTextArray, handler.Nullable()),
			handler.NewColumn(DeviceAuthRequestColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(DeviceAuthRequestColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(DeviceAuthRequestColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(DeviceAuthRequestColumnInstanceID, handler.ColumnTypeText),
		},
			handler.NewPrimaryKey(DeviceAuthRequestColumnInstanceID, DeviceAuthRequestColumnDeviceCode),
			handler.WithIndex(handler.NewIndex("user_code", []string{DeviceAuthRequestColumnInstanceID, DeviceAuthRequestColumnUserCode})),
		),
	)
}

func (p *deviceAuthRequestProjection) Reducers() []handler.AggregateReducer {
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
					Reduce: p.reduceDoneEvents,
				},
				{
					Event:  deviceauth.CanceledEventType,
					Reduce: p.reduceDoneEvents,
				},
			},
		},
	}
}

func (p *deviceAuthRequestProjection) reduceAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*deviceauth.AddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-chu6O", "reduce.wrong.event.type %T != %s", event, deviceauth.AddedEventType)
	}
	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(DeviceAuthRequestColumnClientID, e.ClientID),
			handler.NewCol(DeviceAuthRequestColumnDeviceCode, e.DeviceCode),
			handler.NewCol(DeviceAuthRequestColumnUserCode, e.UserCode),
			handler.NewCol(DeviceAuthRequestColumnScopes, e.Scopes),
			handler.NewCol(DeviceAuthRequestColumnAudience, e.Audience),
			handler.NewCol(DeviceAuthRequestColumnCreationDate, e.CreationDate()),
			handler.NewCol(DeviceAuthRequestColumnChangeDate, e.CreationDate()),
			handler.NewCol(DeviceAuthRequestColumnSequence, e.Sequence()),
			handler.NewCol(DeviceAuthRequestColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}

// reduceDoneEvents removes the device auth request from the projection.
func (p *deviceAuthRequestProjection) reduceDoneEvents(event eventstore.Event) (*handler.Statement, error) {
	switch event.(type) {
	case *deviceauth.ApprovedEvent, *deviceauth.CanceledEvent:
		return handler.NewDeleteStatement(event,
			[]handler.Condition{
				handler.NewCond(DeviceAuthRequestColumnInstanceID, event.Aggregate().InstanceID),
				handler.NewCond(DeviceAuthRequestColumnDeviceCode, event.Aggregate().ID),
			},
		), nil

	default:
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-eeS8d", "reduce.wrong.event.type %T", event)
	}
}
