package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	OIDCSettingsProjectionTable = "projections.oidc_settings2"

	OIDCSettingsColumnAggregateID                = "aggregate_id"
	OIDCSettingsColumnCreationDate               = "creation_date"
	OIDCSettingsColumnChangeDate                 = "change_date"
	OIDCSettingsColumnResourceOwner              = "resource_owner"
	OIDCSettingsColumnInstanceID                 = "instance_id"
	OIDCSettingsColumnSequence                   = "sequence"
	OIDCSettingsColumnAccessTokenLifetime        = "access_token_lifetime"
	OIDCSettingsColumnIdTokenLifetime            = "id_token_lifetime"
	OIDCSettingsColumnRefreshTokenIdleExpiration = "refresh_token_idle_expiration"
	OIDCSettingsColumnRefreshTokenExpiration     = "refresh_token_expiration"
)

type oidcSettingsProjection struct{}

func newOIDCSettingsProjection(ctx context.Context, config handler.Config) *handler.Handler {
	return handler.NewHandler(ctx, &config, new(oidcSettingsProjection))
}

func (*oidcSettingsProjection) Name() string {
	return OIDCSettingsProjectionTable
}

func (*oidcSettingsProjection) Init() *old_handler.Check {
	return handler.NewTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(OIDCSettingsColumnAggregateID, handler.ColumnTypeText),
			handler.NewColumn(OIDCSettingsColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(OIDCSettingsColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(OIDCSettingsColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(OIDCSettingsColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(OIDCSettingsColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(OIDCSettingsColumnAccessTokenLifetime, handler.ColumnTypeInt64),
			handler.NewColumn(OIDCSettingsColumnIdTokenLifetime, handler.ColumnTypeInt64),
			handler.NewColumn(OIDCSettingsColumnRefreshTokenIdleExpiration, handler.ColumnTypeInt64),
			handler.NewColumn(OIDCSettingsColumnRefreshTokenExpiration, handler.ColumnTypeInt64),
		},
			handler.NewPrimaryKey(OIDCSettingsColumnInstanceID, OIDCSettingsColumnAggregateID),
		),
	)
}

func (p *oidcSettingsProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.OIDCSettingsAddedEventType,
					Reduce: p.reduceOIDCSettingsAdded,
				},
				{
					Event:  instance.OIDCSettingsChangedEventType,
					Reduce: p.reduceOIDCSettingsChanged,
				},
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(OIDCSettingsColumnInstanceID),
				},
			},
		},
	}
}

func (p *oidcSettingsProjection) reduceOIDCSettingsAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.OIDCSettingsAddedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-f9nwf", "reduce.wrong.event.type %s", instance.OIDCSettingsAddedEventType)
	}
	return handler.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OIDCSettingsColumnAggregateID, e.Aggregate().ID),
			handler.NewCol(OIDCSettingsColumnCreationDate, e.CreationDate()),
			handler.NewCol(OIDCSettingsColumnChangeDate, e.CreationDate()),
			handler.NewCol(OIDCSettingsColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(OIDCSettingsColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(OIDCSettingsColumnSequence, e.Sequence()),
			handler.NewCol(OIDCSettingsColumnAccessTokenLifetime, e.AccessTokenLifetime),
			handler.NewCol(OIDCSettingsColumnIdTokenLifetime, e.IdTokenLifetime),
			handler.NewCol(OIDCSettingsColumnRefreshTokenIdleExpiration, e.RefreshTokenIdleExpiration),
			handler.NewCol(OIDCSettingsColumnRefreshTokenExpiration, e.RefreshTokenExpiration),
		},
	), nil
}

func (p *oidcSettingsProjection) reduceOIDCSettingsChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.OIDCSettingsChangedEvent)
	if !ok {
		return nil, zerrors.ThrowInvalidArgumentf(nil, "HANDL-8JJ2d", "reduce.wrong.event.type %s", instance.OIDCSettingsChangedEventType)
	}

	columns := make([]handler.Column, 0, 6)
	columns = append(columns,
		handler.NewCol(OIDCSettingsColumnChangeDate, e.CreationDate()),
		handler.NewCol(OIDCSettingsColumnSequence, e.Sequence()),
	)
	if e.AccessTokenLifetime != nil {
		columns = append(columns, handler.NewCol(OIDCSettingsColumnAccessTokenLifetime, *e.AccessTokenLifetime))
	}
	if e.IdTokenLifetime != nil {
		columns = append(columns, handler.NewCol(OIDCSettingsColumnIdTokenLifetime, *e.IdTokenLifetime))
	}
	if e.RefreshTokenIdleExpiration != nil {
		columns = append(columns, handler.NewCol(OIDCSettingsColumnRefreshTokenIdleExpiration, *e.RefreshTokenIdleExpiration))
	}
	if e.RefreshTokenExpiration != nil {
		columns = append(columns, handler.NewCol(OIDCSettingsColumnRefreshTokenExpiration, *e.RefreshTokenExpiration))
	}
	return handler.NewUpdateStatement(
		e,
		columns,
		[]handler.Condition{
			handler.NewCond(OIDCSettingsColumnAggregateID, e.Aggregate().ID),
			handler.NewCond(OIDCSettingsColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}
