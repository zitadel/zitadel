package projection

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/crdb"
	"github.com/zitadel/zitadel/internal/repository/keypair"
)

const (
	KeyProjectionTable = "projections.keys"
	KeyPrivateTable    = KeyProjectionTable + "_" + privateKeyTableSuffix
	KeyPublicTable     = KeyProjectionTable + "_" + publicKeyTableSuffix

	KeyColumnID            = "id"
	KeyColumnCreationDate  = "creation_date"
	KeyColumnChangeDate    = "change_date"
	KeyColumnResourceOwner = "resource_owner"
	KeyColumnInstanceID    = "instance_id"
	KeyColumnSequence      = "sequence"
	KeyColumnAlgorithm     = "algorithm"
	KeyColumnUse           = "use"

	privateKeyTableSuffix      = "private"
	KeyPrivateColumnID         = "id"
	KeyPrivateColumnInstanceID = "instance_id"
	KeyPrivateColumnExpiry     = "expiry"
	KeyPrivateColumnKey        = "key"

	publicKeyTableSuffix      = "public"
	KeyPublicColumnID         = "id"
	KeyPublicColumnInstanceID = "instance_id"
	KeyPublicColumnExpiry     = "expiry"
	KeyPublicColumnKey        = "key"
)

type KeyProjection struct {
	crdb.StatementHandler
	encryptionAlgorithm crypto.EncryptionAlgorithm
}

func NewKeyProjection(ctx context.Context, config crdb.StatementHandlerConfig, keyEncryptionAlgorithm crypto.EncryptionAlgorithm) *KeyProjection {
	p := new(KeyProjection)
	config.ProjectionName = KeyProjectionTable
	config.Reducers = p.reducers()
	config.InitCheck = crdb.NewMultiTableCheck(
		crdb.NewTable([]*crdb.Column{
			crdb.NewColumn(KeyColumnID, crdb.ColumnTypeText),
			crdb.NewColumn(KeyColumnCreationDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(KeyColumnChangeDate, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(KeyColumnResourceOwner, crdb.ColumnTypeText),
			crdb.NewColumn(KeyColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(KeyColumnSequence, crdb.ColumnTypeInt64),
			crdb.NewColumn(KeyColumnAlgorithm, crdb.ColumnTypeText, crdb.Default("")),
			crdb.NewColumn(KeyColumnUse, crdb.ColumnTypeEnum, crdb.Default(0)),
		},
			crdb.NewPrimaryKey(KeyColumnID, KeyColumnInstanceID),
			crdb.WithConstraint(crdb.NewConstraint("id_unique", []string{KeyColumnID})),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(KeyPrivateColumnID, crdb.ColumnTypeText),
			crdb.NewColumn(KeyPrivateColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(KeyPrivateColumnExpiry, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(KeyPrivateColumnKey, crdb.ColumnTypeJSONB),
		},
			crdb.NewPrimaryKey(KeyPrivateColumnID, KeyPrivateColumnInstanceID),
			privateKeyTableSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys("fk_private_ref_keys")),
		),
		crdb.NewSuffixedTable([]*crdb.Column{
			crdb.NewColumn(KeyPublicColumnID, crdb.ColumnTypeText),
			crdb.NewColumn(KeyPublicColumnInstanceID, crdb.ColumnTypeText),
			crdb.NewColumn(KeyPublicColumnExpiry, crdb.ColumnTypeTimestamp),
			crdb.NewColumn(KeyPublicColumnKey, crdb.ColumnTypeBytes),
		},
			crdb.NewPrimaryKey(KeyPublicColumnID, KeyPublicColumnInstanceID),
			publicKeyTableSuffix,
			crdb.WithForeignKey(crdb.NewForeignKeyOfPublicKeys("fk_public_ref_keys")),
		),
	)
	p.StatementHandler = crdb.NewStatementHandler(ctx, config)
	p.encryptionAlgorithm = keyEncryptionAlgorithm

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

func (p *KeyProjection) reduceKeyPairAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*keypair.AddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-SAbr2", "reduce.wrong.event.type %s", keypair.AddedEventType)
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
				handler.NewCol(KeyColumnInstanceID, e.Aggregate().InstanceID),
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
				handler.NewCol(KeyPrivateColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(KeyPrivateColumnExpiry, e.PrivateKey.Expiry),
				handler.NewCol(KeyPrivateColumnKey, e.PrivateKey.Key),
			},
			crdb.WithTableSuffix(privateKeyTableSuffix),
		))
	}
	if e.PublicKey.Expiry.After(time.Now()) {
		publicKey, err := crypto.Decrypt(e.PublicKey.Key, p.encryptionAlgorithm)
		if err != nil {
			return nil, errors.ThrowInternal(err, "HANDL-DAg2f", "cannot decrypt public key")
		}
		creates = append(creates, crdb.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(KeyPublicColumnID, e.Aggregate().ID),
				handler.NewCol(KeyPublicColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(KeyPublicColumnExpiry, e.PublicKey.Expiry),
				handler.NewCol(KeyPublicColumnKey, publicKey),
			},
			crdb.WithTableSuffix(publicKeyTableSuffix),
		))
	}
	return crdb.NewMultiStatement(e, creates...), nil
}
