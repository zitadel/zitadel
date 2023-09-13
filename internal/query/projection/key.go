package projection

import (
	"context"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/eventstore"
	old_handler "github.com/zitadel/zitadel/internal/eventstore/handler"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	"github.com/zitadel/zitadel/internal/repository/instance"
	"github.com/zitadel/zitadel/internal/repository/keypair"
)

const (
	KeyProjectionTable = "projections.keys4"
	KeyPrivateTable    = KeyProjectionTable + "_" + privateKeyTableSuffix
	KeyPublicTable     = KeyProjectionTable + "_" + publicKeyTableSuffix
	CertificateTable   = KeyProjectionTable + "_" + certificateTableSuffix

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

	certificateTableSuffix       = "certificate"
	CertificateColumnID          = "id"
	CertificateColumnInstanceID  = "instance_id"
	CertificateColumnExpiry      = "expiry"
	CertificateColumnCertificate = "certificate"
)

type keyProjection struct {
	encryptionAlgorithm     crypto.EncryptionAlgorithm
	certEncryptionAlgorithm crypto.EncryptionAlgorithm
}

func newKeyProjection(ctx context.Context, config handler.Config, keyEncryptionAlgorithm, certEncryptionAlgorithm crypto.EncryptionAlgorithm) *handler.Handler {
	p := &keyProjection{
		encryptionAlgorithm:     keyEncryptionAlgorithm,
		certEncryptionAlgorithm: certEncryptionAlgorithm,
	}
	return handler.NewHandler(ctx, &config, p)
}

func (*keyProjection) Name() string {
	return KeyProjectionTable
}

func (*keyProjection) Init() *old_handler.Check {
	return handler.NewMultiTableCheck(
		handler.NewTable([]*handler.InitColumn{
			handler.NewColumn(KeyColumnID, handler.ColumnTypeText),
			handler.NewColumn(KeyColumnCreationDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(KeyColumnChangeDate, handler.ColumnTypeTimestamp),
			handler.NewColumn(KeyColumnResourceOwner, handler.ColumnTypeText),
			handler.NewColumn(KeyColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(KeyColumnSequence, handler.ColumnTypeInt64),
			handler.NewColumn(KeyColumnAlgorithm, handler.ColumnTypeText, handler.Default("")),
			handler.NewColumn(KeyColumnUse, handler.ColumnTypeEnum, handler.Default(0)),
		},
			handler.NewPrimaryKey(KeyColumnInstanceID, KeyColumnID),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(KeyPrivateColumnID, handler.ColumnTypeText),
			handler.NewColumn(KeyPrivateColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(KeyPrivateColumnExpiry, handler.ColumnTypeTimestamp),
			handler.NewColumn(KeyPrivateColumnKey, handler.ColumnTypeJSONB),
		},
			handler.NewPrimaryKey(KeyPrivateColumnInstanceID, KeyPrivateColumnID),
			privateKeyTableSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(KeyPublicColumnID, handler.ColumnTypeText),
			handler.NewColumn(KeyPublicColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(KeyPublicColumnExpiry, handler.ColumnTypeTimestamp),
			handler.NewColumn(KeyPublicColumnKey, handler.ColumnTypeBytes),
		},
			handler.NewPrimaryKey(KeyPublicColumnInstanceID, KeyPublicColumnID),
			publicKeyTableSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
		handler.NewSuffixedTable([]*handler.InitColumn{
			handler.NewColumn(CertificateColumnID, handler.ColumnTypeText),
			handler.NewColumn(CertificateColumnInstanceID, handler.ColumnTypeText),
			handler.NewColumn(CertificateColumnExpiry, handler.ColumnTypeTimestamp),
			handler.NewColumn(CertificateColumnCertificate, handler.ColumnTypeBytes),
		},
			handler.NewPrimaryKey(CertificateColumnInstanceID, CertificateColumnID),
			certificateTableSuffix,
			handler.WithForeignKey(handler.NewForeignKeyOfPublicKeys()),
		),
	)
}

func (p *keyProjection) Reducers() []handler.AggregateReducer {
	return []handler.AggregateReducer{
		{
			Aggregate: keypair.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  keypair.AddedEventType,
					Reduce: p.reduceKeyPairAdded,
				},
				{
					Event:  keypair.AddedCertificateEventType,
					Reduce: p.reduceCertificateAdded,
				},
			},
		},
		{
			Aggregate: instance.AggregateType,
			EventReducers: []handler.EventReducer{
				{
					Event:  instance.InstanceRemovedEventType,
					Reduce: reduceInstanceRemovedHelper(KeyColumnInstanceID),
				},
			},
		},
	}
}

func (p *keyProjection) reduceKeyPairAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*keypair.AddedEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-SAbr2", "reduce.wrong.event.type %s", keypair.AddedEventType)
	}
	if e.PrivateKey.Expiry.Before(time.Now()) && e.PublicKey.Expiry.Before(time.Now()) {
		return handler.NewNoOpStatement(e), nil
	}
	creates := []func(eventstore.Event) handler.Exec{
		handler.AddCreateStatement(
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
		creates = append(creates, handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(KeyPrivateColumnID, e.Aggregate().ID),
				handler.NewCol(KeyPrivateColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(KeyPrivateColumnExpiry, e.PrivateKey.Expiry),
				handler.NewCol(KeyPrivateColumnKey, e.PrivateKey.Key),
			},
			handler.WithTableSuffix(privateKeyTableSuffix),
		))
	}
	if e.PublicKey.Expiry.After(time.Now()) {
		publicKey, err := crypto.Decrypt(e.PublicKey.Key, p.encryptionAlgorithm)
		if err != nil {
			return nil, errors.ThrowInternal(err, "HANDL-DAg2f", "cannot decrypt public key")
		}
		creates = append(creates, handler.AddCreateStatement(
			[]handler.Column{
				handler.NewCol(KeyPublicColumnID, e.Aggregate().ID),
				handler.NewCol(KeyPublicColumnInstanceID, e.Aggregate().InstanceID),
				handler.NewCol(KeyPublicColumnExpiry, e.PublicKey.Expiry),
				handler.NewCol(KeyPublicColumnKey, publicKey),
			},
			handler.WithTableSuffix(publicKeyTableSuffix),
		))
	}

	return handler.NewMultiStatement(e, creates...), nil
}

func (p *keyProjection) reduceCertificateAdded(event eventstore.Event) (*handler.Statement, error) {
	e, ok := event.(*keypair.AddedCertificateEvent)
	if !ok {
		return nil, errors.ThrowInvalidArgumentf(nil, "HANDL-SAbr09", "reduce.wrong.event.type %s", keypair.AddedCertificateEventType)
	}

	if e.Certificate.Expiry.Before(time.Now()) {
		return handler.NewNoOpStatement(e), nil
	}

	certificate, err := crypto.Decrypt(e.Certificate.Key, p.certEncryptionAlgorithm)
	if err != nil {
		return nil, errors.ThrowInternal(err, "HANDL-Dajwig2f", "cannot decrypt certificate")
	}

	creates := []func(eventstore.Event) handler.Exec{handler.AddCreateStatement(
		[]handler.Column{
			handler.NewCol(CertificateColumnID, e.Aggregate().ID),
			handler.NewCol(CertificateColumnInstanceID, e.Aggregate().InstanceID),
			handler.NewCol(CertificateColumnExpiry, e.Certificate.Expiry),
			handler.NewCol(CertificateColumnCertificate, certificate),
		},
		handler.WithTableSuffix(certificateTableSuffix),
	)}

	return handler.NewMultiStatement(e, creates...), nil
}
