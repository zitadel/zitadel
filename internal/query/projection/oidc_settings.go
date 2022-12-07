package projection

import (
	"context"

	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/instance"
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

type oidcSettingsProjection struct {
	crdb.StatementHandler
}

func newOIDCSettingsProjection(ctx context.Context, config crdb.StatementHandlerConfig) *oidcSettingsProjection {
	p := new(oidcSettingsProjection)
	config.ProjectionName = OIDCSettingsProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(OIDCSettingsColumnAggregateID, crdb.ColumnTypeText),
			crdb.NewColumn(OIDCSettingsColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(OIDCSettingsColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(OIDCSettingsColumnResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(OIDCSettingsColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(OIDCSettingsColumnSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(OIDCSettingsColumnAccessTokenLifetime, crdb.ColumnTypeInt64),
			crdb.NewColumn(OIDCSettingsColumnIdTokenLifetime, crdb.ColumnTypeInt64),
			crdb.NewColumn(OIDCSettingsColumnRefreshTokenIdleExpiration, crdb.ColumnTypeInt64),
			crdb.NewColumn(OIDCSettingsColumnRefreshTokenExpiration, crdb.ColumnTypeInt64),
		},
			crdb.NewPrimaryKey(OIDCSettingsColumnInstanceID, OIDCSettingsColumnAggregateID),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *oidcSettingsProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: instance.AggregateType,
			EventRedusers: []handler.EventReducer{
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
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-f9nwf", "reduce.wrong.event.type %s", instance.OIDCSettingsAddedEventType)
	}
	return crdb.NewCreateStatement(
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
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-8JJ2d", "reduce.wrong.event.type %s", instance.OIDCSettingsChangedEventType)
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
	return crdb.NewUpdateStatement(
		e,
		columns,
		[]handler.Condition{
			handler.NewCond(OIDCSettingsColumnAggregateID, e.Aggregate().ID),
			handler.NewCond(OIDCSettingsColumnInstanceID, e.Aggregate().InstanceID),
		},
	), nil
}
