package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/webkey"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	WebKeyTable = "projections.web_keys1"

	WebKeyInstanceIDCol   = "instance_id"
	WebKeyKeyIDCol        = "key_id"
	WebKeyCreationDateCol = "creation_date"
	WebKeyChangeDateCol   = "change_date"
	WebKeySequenceCol     = "sequence"
	WebKeyStateCol        = "state"
	WebKeyPrivateKeyCol   = "private_key"
	WebKeyPublicKeyCol    = "public_key"
	WebKeyConfigCol       = "config"
	WebKeyConfigTypeCol   = "config_type"
)

type webKeyProjection struct{}

func newWebKeyProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(webKeyProjection))
}

func (*webKeyProjection) Name() string {
	return WebKeyTable
}

func (*webKeyProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable(
			[]*handler.InitColumn{
				handler.NewColumn(WebKeyInstanceIDCol, handler.ColumnTypeText),
				handler.NewColumn(WebKeyKeyIDCol, handler.ColumnTypeText),
				handler.NewColumn(WebKeyCreationDateCol, handler.ColumnTypeTimestamp),
				handler.NewColumn(WebKeyChangeDateCol, handler.ColumnTypeTimestamp),
				handler.NewColumn(WebKeySequenceCol, handler.ColumnTypeInt64),
				handler.NewColumn(WebKeyStateCol, handler.ColumnTypeInt64),
				handler.NewColumn(WebKeyPrivateKeyCol, handler.ColumnTypeJSONB),
				handler.NewColumn(WebKeyPublicKeyCol, handler.ColumnTypeJSONB),
				handler.NewColumn(WebKeyConfigCol, handler.ColumnTypeJSONB),
				handler.NewColumn(WebKeyConfigTypeCol, handler.ColumnTypeInt64),
			},
			handler.NewPrimaryKey(WebKeyInstanceIDCol, WebKeyKeyIDCol),

			// index to find the current active private key for an instance.
			handler.WithIndex(handler.NewIndex(
				"web_key_state",
				[]string{WebKeyInstanceIDCol, WebKeyStateCol},
			)),
		),
	)
}

func (p *webKeyProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{{
		Aggregate: webkey.AggregateType,
		EventReducers: []handler.EventReducer{
			{
				Event:  webkey.AddedEventType,
				Reduce: p.reduceWebKeyAdded,
			},
			{
				Event:  webkey.ActivatedEventType,
				Reduce: p.reduceWebKeyActivated,
			},
			{
				Event:  webkey.DeactivatedEventType,
				Reduce: p.reduceWebKeyDeactivated,
			},
			{
				Event:  webkey.RemovedEventType,
				Reduce: p.reduceWebKeyRemoved,
			},
			{
				Event:  instance.InstanceRemovedEventType,
				Reduce: reduceInstanceRemovedHelper(WebKeyInstanceIDCol),
			},
		},
	}}
}

func (p *webKeyProjection) reduceWebKeyAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*webkey.AddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-jei2K", "reduce.wrong.event.type %s", webkey.AddedEventType)
	}
	return handler.NewCreateStatement(e,
		[]handler.Column{
			handler.NewCol(WebKeyInstanceIDCol, e.Agg.InstanceID),
			handler.NewCol(WebKeyKeyIDCol, e.Agg.ID),
			handler.NewCol(WebKeyCreationDateCol, e.CreationDate()),
			handler.NewCol(WebKeyChangeDateCol, e.CreationDate()),
			handler.NewCol(WebKeySequenceCol, e.Sequence()),
			handler.NewCol(WebKeyStateCol, domain.WebKeyStateInitial),
			handler.NewCol(WebKeyPrivateKeyCol, e.PrivateKey),
			handler.NewCol(WebKeyPublicKeyCol, e.PublicKey),
			handler.NewCol(WebKeyConfigCol, e.Config),
			handler.NewCol(WebKeyConfigTypeCol, e.ConfigType),
		},
	), nil
}

func (p *webKeyProjection) reduceWebKeyActivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*webkey.ActivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-iiQu2", "reduce.wrong.event.type %s", webkey.ActivatedEventType)
	}
	return handler.NewUpdateStatement(e,
		[]handler.Column{
			handler.NewCol(WebKeyChangeDateCol, e.CreationDate()),
			handler.NewCol(WebKeySequenceCol, e.Sequence()),
			handler.NewCol(WebKeyStateCol, domain.WebKeyStateActive),
		},
		[]handler.Condition{
			handler.NewCond(WebKeyInstanceIDCol, e.Agg.InstanceID),
			handler.NewCond(WebKeyKeyIDCol, e.Agg.ID),
		},
	), nil
}

func (p *webKeyProjection) reduceWebKeyDeactivated(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*webkey.DeactivatedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-zei3E", "reduce.wrong.event.type %s", webkey.DeactivatedEventType)
	}
	return handler.NewUpdateStatement(e,
		[]handler.Column{
			handler.NewCol(WebKeyChangeDateCol, e.CreationDate()),
			handler.NewCol(WebKeySequenceCol, e.Sequence()),
			handler.NewCol(WebKeyStateCol, domain.WebKeyStateInactive),
		},
		[]handler.Condition{
			handler.NewCond(WebKeyInstanceIDCol, e.Agg.InstanceID),
			handler.NewCond(WebKeyKeyIDCol, e.Agg.ID),
		},
	), nil
}

func (p *webKeyProjection) reduceWebKeyRemoved(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*webkey.RemovedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "PROJE-Zei6f", "reduce.wrong.event.type %s", webkey.RemovedEventType)
	}
	return handler.NewDeleteStatement(e,
		[]handler.Condition{
			handler.NewCond(WebKeyInstanceIDCol, e.Agg.InstanceID),
			handler.NewCond(WebKeyKeyIDCol, e.Agg.ID),
		},
	), nil
}
