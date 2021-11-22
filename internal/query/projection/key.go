package projection

import (
	"context"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/keypair"
)

type KeyProjection struct {
	crdb.StatementHandler
}

const KeyProjectionTable = "zitadel.projections.keys"

func NewKeyProjection(ctx context.Context, config crdb.StatementHandlerConfig) *KeyProjection {
	p := &KeyProjection{}
	config.ProjectionName = KeyProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	return p
}

func (p *KeyProjection) reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: keypair.AggregateType,
			EventRedusers: []handler.EventReducer{
				{
					Event:  keypair.AddedEventType,
					Reduce: p.reduceKeyPairAdded,
				},
			},
		},
	}
}

const (
	KeyColumnID            = "id"
	KeyColumnIsPrivate     = "is_private"
	KeyColumnCreationDate  = "creation_date"
	KeyColumnChangeDate    = "change_date"
	KeyColumnResourceOwner = "resource_owner"
	KeyColumnSequence      = "sequence"
	KeyColumnAlgorithm     = "algorithm"
	KeyColumnUse           = "use"
	KeyColumnExpiry        = "expiry"
	KeyColumnKey           = "key"
)

func (p *KeyProjection) reduceKeyPairAdded(event eventstore.EventReader) (*handler.Statement, error) {
	e, ok := event.(*keypair.AddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-GEdg3", "seq", event.Sequence(), "expectedType", keypair.AddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-SAbr2", "reduce.wrong.event.type")
	}
	if e.PrivateKey.Expiry.Before(time.Now()) && e.PublicKey.Expiry.Before(time.Now()) {
		return crdb.NewNoOpStatement(e), nil
	}

	return crdb.NewMultiStatement(e,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(KeyColumnID, e.Aggregate().ID),
				handler.NewCol(KeyColumnIsPrivate, true),
				handler.NewCol(KeyColumnCreationDate, e.CreationDate()),
				handler.NewCol(KeyColumnChangeDate, e.CreationDate()),
				handler.NewCol(KeyColumnResourceOwner, e.Aggregate().ResourceOwner),
				handler.NewCol(KeyColumnSequence, e.Sequence()),
				handler.NewCol(KeyColumnAlgorithm, e.Algorithm),
				handler.NewCol(KeyColumnUse, e.Usage),
				handler.NewCol(KeyColumnExpiry, e.PrivateKey.Expiry),
				handler.NewCol(KeyColumnKey, e.PrivateKey.Key),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(KeyColumnID, e.Aggregate().ID),
				handler.NewCol(KeyColumnIsPrivate, false),
				handler.NewCol(KeyColumnCreationDate, e.CreationDate()),
				handler.NewCol(KeyColumnChangeDate, e.CreationDate()),
				handler.NewCol(KeyColumnResourceOwner, e.Aggregate().ResourceOwner),
				handler.NewCol(KeyColumnSequence, e.Sequence()),
				handler.NewCol(KeyColumnAlgorithm, e.Algorithm),
				handler.NewCol(KeyColumnUse, e.Usage),
				handler.NewCol(KeyColumnExpiry, e.PublicKey.Expiry),
				handler.NewCol(KeyColumnKey, e.PublicKey.Key),
			},
		),
	), nil
}
