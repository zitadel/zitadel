package projection

import (
	"context"
	"time"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/crypto"
	"github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore"
	"github.com/caos/zitadel/internal/eventstore/handler"
	"github.com/caos/zitadel/internal/eventstore/handler/crdb"
	"github.com/caos/zitadel/internal/repository/keypair"
)

type KeyProjection struct {
	crdb.StatementHandler
	encryptionAlgorithm crypto.EncryptionAlgorithm
}

const (
	KeyProjectionTable = "zitadel.projections.keys"
	KeyPrivateTable    = KeyProjectionTable + "_" + privateKeyTableSuffix
	KeyPublicTable     = KeyProjectionTable + "_" + publicKeyTableSuffix
)

func NewKeyProjection(ctx context.Context, config crdb.StatementHandlerConfig, keyConfig systemdefaults.KeyConfig) (_ *KeyProjection, err error) {
	p := &KeyProjection{}
	config.ProjectionName = KeyProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	p.encryptionAlgorithm, err = crypto.NewAESCrypto(keyConfig.EncryptionConfig)
	if err != nil {
		return nil, err
	}
	return p, nil
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
	KeyColumnCreationDate  = "creation_date"
	KeyColumnChangeDate    = "change_date"
	KeyColumnResourceOwner = "resource_owner"
	KeyColumnSequence      = "sequence"
	KeyColumnAlgorithm     = "algorithm"
	KeyColumnUse           = "use"

	privateKeyTableSuffix  = "private"
	KeyPrivateColumnID     = "id"
	KeyPrivateColumnExpiry = "expiry"
	KeyPrivateColumnKey    = "key"

	publicKeyTableSuffix  = "public"
	KeyPublicColumnID     = "id"
	KeyPublicColumnExpiry = "expiry"
	KeyPublicColumnKey    = "key"
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
	publicKey, err := crypto.Decrypt(e.PublicKey.Key, p.encryptionAlgorithm)
	if err != nil {
		logging.LogWithFields("HANDL-SDfw2", "seq", event.Sequence()).Error("cannot decrypt public key")
		return nil, errors.ThrowInternal(err, "HANDL-DAg2f", "cannot decrypt public key")
	}

	return crdb.NewMultiStatement(e,
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(KeyColumnID, e.Aggregate().ID),
				handler.NewCol(KeyColumnCreationDate, e.CreationDate()),
				handler.NewCol(KeyColumnChangeDate, e.CreationDate()),
				handler.NewCol(KeyColumnResourceOwner, e.Aggregate().ResourceOwner),
				handler.NewCol(KeyColumnSequence, e.Sequence()),
				handler.NewCol(KeyColumnAlgorithm, e.Algorithm),
				handler.NewCol(KeyColumnUse, e.Usage),
			},
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(KeyPrivateColumnID, e.Aggregate().ID),
				handler.NewCol(KeyPrivateColumnExpiry, e.PrivateKey.Expiry),
				handler.NewCol(KeyPrivateColumnKey, e.PrivateKey.Key),
			},
			crdb.WithTableSuffix(privateKeyTableSuffix),
		),
		crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(KeyPublicColumnID, e.Aggregate().ID),
				handler.NewCol(KeyPublicColumnExpiry, e.PublicKey.Expiry),
				handler.NewCol(KeyPublicColumnKey, publicKey),
			},
			crdb.WithTableSuffix(publicKeyTableSuffix),
		),
	), nil
}
