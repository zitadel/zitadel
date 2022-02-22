package projection

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/iam"
	"github.com/caos/zitadel/internal/repository/project"
)

type OIDCConfigProjection struct {
	crdb.StatementHandler
}

const (
	OIDCConfigProjectionTable = "zitadel.projections.oidc_configs"
)

func NewOIDCConfigProjection(ctx context.Context, config crdb.StatementHandlerConfig) *OIDCConfigProjection {
	p := &OIDCConfigProjection{}
	config.ProjectionName = OIDCConfigProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *OIDCConfigProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: project.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  iam.OIDCConfigAddedEventType,
					Reduce: p.reduceOIDCConfigAdded,
				},
				{
					Event:  iam.OIDCConfigChangedEventType,
					Reduce: p.reduceOIDCConfigChanged,
				},
			},
		},
	}
}

const (
	OIDCConfigColumnAggregateID                = "aggregate_id"
	OIDCConfigColumnCreationDate               = "creation_date"
	OIDCConfigColumnChangeDate                 = "change_date"
	OIDCConfigColumnResourceOwner              = "resource_owner"
	OIDCConfigColumnSequence                   = "sequence"
	OIDCConfigColumnAccessTokenLifetime        = "access_token_lifetime"
	OIDCConfigColumnIdTokenLifetime            = "id_token_lifetime"
	OIDCConfigColumnRefreshTokenIdleExpiration = "refresh_token_idle_expiration"
	OIDCConfigColumnRefreshTokenExpiration     = "refresh_token_expiration"
)

func (p *OIDCConfigProjection) reduceOIDCConfigAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.OIDCConfigAddedEvent)
	if !ok {
		logging.WithFields("seq", event.Sequence(), "expectedType", iam.OIDCConfigAddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-f9nwf", "reduce.wrong.event.type")
	}
	return crdb.NewCreateStatement(
		e,
		[]handler.Column{
			handler.NewCol(OIDCConfigColumnAggregateID, e.Aggregate().ID),
			handler.NewCol(OIDCConfigColumnCreationDate, e.CreationDate()),
			handler.NewCol(OIDCConfigColumnChangeDate, e.CreationDate()),
			handler.NewCol(OIDCConfigColumnResourceOwner, e.Aggregate().ResourceOwner),
			handler.NewCol(OIDCConfigColumnSequence, e.Sequence()),
			handler.NewCol(OIDCConfigColumnAccessTokenLifetime, e.AccessTokenLifetime),
			handler.NewCol(OIDCConfigColumnIdTokenLifetime, e.IdTokenLifetime),
			handler.NewCol(OIDCConfigColumnRefreshTokenIdleExpiration, e.RefreshTokenIdleExpiration),
			handler.NewCol(OIDCConfigColumnRefreshTokenExpiration, e.RefreshTokenExpiration),
		},
	), nil
}

func (p *OIDCConfigProjection) reduceOIDCConfigChanged(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*iam.OIDCConfigChangedEvent)
	if !ok {
		logging.WithFields("seq", event.Sequence(), "expected", iam.OIDCConfigChangedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-8JJ2d", "reduce.wrong.event.type")
	}

	columns := make([]handler.Column, 0, 6)
	columns = append(columns, handler.NewCol(OIDCConfigColumnChangeDate, e.CreationDate()),
		handler.NewCol(OIDCConfigColumnSequence, e.Sequence()))
	if e.AccessTokenLifetime != nil {
		columns = append(columns, handler.NewCol(OIDCConfigColumnAccessTokenLifetime, *e.AccessTokenLifetime))
	}
	if e.IdTokenLifetime != nil {
		columns = append(columns, handler.NewCol(OIDCConfigColumnIdTokenLifetime, *e.IdTokenLifetime))
	}
	if e.RefreshTokenIdleExpiration != nil {
		columns = append(columns, handler.NewCol(OIDCConfigColumnRefreshTokenIdleExpiration, *e.RefreshTokenIdleExpiration))
	}
	if e.RefreshTokenExpiration != nil {
		columns = append(columns, handler.NewCol(OIDCConfigColumnRefreshTokenExpiration, *e.RefreshTokenExpiration))
	}
	return crdb.NewUpdateStatement(
		e,
		columns,
		[]handler.Condition{
			handler.NewCond(OIDCConfigColumnAggregateID, e.Aggregate().ID),
		},
	), nil
}
