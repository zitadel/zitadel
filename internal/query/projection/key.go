package projection

import (
	"context"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/keypair"
)

type keyProjection struct {
	crdb.StatementHandler
	encryptionAlgorithm crypto.EncryptionAlgorithm
	keyChan             chan<- interface{}
}

const (
	KeyProjectionTable = "zitadel.projections.keys"
	KeyPrivateTable    = KeyProjectionTable + "_" + privateKeyTableSuffix
	KeyPublicTable     = KeyProjectionTable + "_" + publicKeyTableSuffix
)

func newKeyProjection(ctx context.Context, config crdb.StatementHandlerConfig, keyConfig systemdefaults.KeyConfig, keyChan chan<- interface{}) (_ *keyProjection, err error) {
	p := &keyProjection{}
	config.ProjectionName = KeyProjectionTable
	config.Reducers = p.reducers()
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	p.keyChan = keyChan
	p.encryptionAlgorithm, err = crypto.NewAESCrypto(keyConfig.EncryptionConfig)
	if err != nil {
		return nil, err
	}
	return p, nil
}

func (p *keyProjection) reducers() []handler.AggregateReducer {
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

func (p *keyProjection) reduceKeyPairAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*keypair.AddedEvent)
	if !ok {
		logging.LogWithFields("HANDL-GEdg3", "seq", event.Sequence(), "expectedType", keypair.AddedEventType).Error("wrong event type")
		return nil, errors.ThrowInvalidArgument(nil, "HANDL-SAbr2", "reduce.wrong.event.type")
	}
	if e.PrivateKey.Expiry.Before(time.Now()) && e.PublicKey.Expiry.Before(time.Now()) {
		return crdb.NewNoOpStatement(e), nil
	}
	creates := []func(eventstore.Event) crdb.Exec{
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
	}
	if e.PrivateKey.Expiry.After(time.Now()) {
		creates = append(creates, crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(KeyPrivateColumnID, e.Aggregate().ID),
				handler.NewCol(KeyPrivateColumnExpiry, e.PrivateKey.Expiry),
				handler.NewCol(KeyPrivateColumnKey, e.PrivateKey.Key),
			},
			crdb.WithTableSuffix(privateKeyTableSuffix),
		))
		if p.keyChan != nil {
			p.keyChan <- true
		}
	}
	if e.PublicKey.Expiry.After(time.Now()) {
		publicKey, err := crypto.Decrypt(e.PublicKey.Key, p.encryptionAlgorithm)
		if err != nil {
			logging.LogWithFields("HANDL-SDfw2", "seq", event.Sequence()).Error("cannot decrypt public key")
			return nil, errors.ThrowInternal(err, "HANDL-DAg2f", "cannot decrypt public key")
		}
		creates = append(creates, crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(KeyPublicColumnID, e.Aggregate().ID),
				handler.NewCol(KeyPublicColumnExpiry, e.PublicKey.Expiry),
				handler.NewCol(KeyPublicColumnKey, publicKey),
			},
			crdb.WithTableSuffix(publicKeyTableSuffix),
		))
	}
	return crdb.NewMultiStatement(e, creates...), nil
}
