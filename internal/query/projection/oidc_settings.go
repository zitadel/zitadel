package projection

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/instance"
	"github.com/caos/zitadel/internal/repository/project"
)

type OIDCSettingsProjection struct {
	crdb.StatementHandler
}

const (
	OIDCSettingsProjectionTable = "zitadel.projections.oidc_settings"
)

func NewOIDCSettingsProjection(ctx context.Context, config crdb.StatementHandlerConfig) *OIDCSettingsProjection {
	p := new(OIDCSettingsProjection)
	config.ProjectionName = OIDCSettingsProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *OIDCSettingsProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  instance.OIDCSettingsAddedEventType,
					Reduce: p.reduceOIDCSettingsAdded,
				},
				{
					Event:  instance.OIDCSettingsChangedEventType,
					Reduce: p.reduceOIDCSettingsChanged,
				},
			},
		},
	}
}

const (
	OIDCSettingsColumnAggregateID                = "aggregate_id"
	OIDCSettingsColumnCreationDate               = "creation_date"
	OIDCSettingsColumnChangeDate                 = "change_date"
	OIDCSettingsColumnResourceOwner              = "resource_owner"
	OIDCSettingsColumnSequence                   = "sequence"
	OIDCSettingsColumnAccessTokenLifetime        = "access_token_lifetime"
	OIDCSettingsColumnIdTokenLifetime            = "id_token_lifetime"
	OIDCSettingsColumnRefreshTokenIdleExpiration = "refresh_token_idle_expiration"
	OIDCSettingsColumnRefreshTokenExpiration     = "refresh_token_expiration"
)

func (p *OIDCSettingsProjection) reduceOIDCSettingsAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.OIDCSettingsAddedEvent)
	if !ok {
		logging.WithFields("seq", event.Sequence(), "expectedType", instance.OIDCSettingsAddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-f9nwf", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OIDCSettingsColumnAggregateID, e.Aggregate().ID),
			handler.NewCol(OIDCSettingsColumnCreationDate, e.CreationDate()),
			handler.NewCol(OIDCSettingsColumnChangeDate, e.CreationDate()),
			handler.NewCol(OIDCSettingsColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(OIDCSettingsColumnSequence, e.Sequence()),
			handler.NewCol(OIDCSettingsColumnAccessTokenLifetime, e.AccessTokenLifetime),
			handler.NewCol(OIDCSettingsColumnIdTokenLifetime, e.IdTokenLifetime),
			handler.NewCol(OIDCSettingsColumnRefreshTokenIdleExpiration, e.RefreshTokenIdleExpiration),
			handler.NewCol(OIDCSettingsColumnRefreshTokenExpiration, e.RefreshTokenExpiration),
		},
	), nil
}

func (p *OIDCSettingsProjection) reduceOIDCSettingsChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*instance.OIDCSettingsChangedEvent)
	if !ok {
		logging.WithFields("seq", event.Sequence(), "expected", instance.OIDCSettingsChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-8JJ2d", "reduce.wrong.event.type")
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
		},
	), nil
}
