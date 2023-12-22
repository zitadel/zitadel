package query

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"time"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/api/call"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/repository/keypair"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type Key interface {
	ID() string
	Algorithm() string
	Use() domain.KeyUsage
	Sequence() uint64
}

type PrivateKey interface {
	Key
	Expiry() time.Time
	Key() *crypto.CryptoValue
}

type PublicKey interface {
	Key
	Expiry() time.Time
	Key() interface{}
}

type PrivateKeys struct {
	SearchResponse
	Keys []PrivateKey
}

type PublicKeys struct {
	SearchResponse
	Keys []PublicKey
}

type key struct {
	id            string
	creationDate  time.Time
	changeDate    time.Time
	sequence      uint64
	resourceOwner string
	algorithm     string
	use           domain.KeyUsage
}

func (k *key) ID() string {
	return k.id
}

func (k *key) Algorithm() string {
	return k.algorithm
}

func (k *key) Use() domain.KeyUsage {
	return k.use
}

func (k *key) Sequence() uint64 {
	return k.sequence
}

type privateKey struct {
	key
	expiry     time.Time
	privateKey *crypto.CryptoValue
}

func (k *privateKey) Expiry() time.Time {
	return k.expiry
}

func (k *privateKey) Key() *crypto.CryptoValue {
	return k.privateKey
}

type rsaPublicKey struct {
	key
	expiry    time.Time
	publicKey *rsa.PublicKey
}

func (r *rsaPublicKey) Expiry() time.Time {
	return r.expiry
}

func (r *rsaPublicKey) Key() interface{} {
	return r.publicKey
}

var (
	keyTable = table{
		name:          projection.KeyProjectionTable,
		instanceIDCol: projection.KeyColumnInstanceID,
	}
	KeyColID = Column{
		name:  projection.KeyColumnID,
		table: keyTable,
	}
	KeyColCreationDate = Column{
		name:  projection.KeyColumnCreationDate,
		table: keyTable,
	}
	KeyColChangeDate = Column{
		name:  projection.KeyColumnChangeDate,
		table: keyTable,
	}
	KeyColResourceOwner = Column{
		name:  projection.KeyColumnResourceOwner,
		table: keyTable,
	}
	KeyColInstanceID = Column{
		name:  projection.KeyColumnInstanceID,
		table: keyTable,
	}
	KeyColSequence = Column{
		name:  projection.KeyColumnSequence,
		table: keyTable,
	}
	KeyColAlgorithm = Column{
		name:  projection.KeyColumnAlgorithm,
		table: keyTable,
	}
	KeyColUse = Column{
		name:  projection.KeyColumnUse,
		table: keyTable,
	}
)

var (
	keyPrivateTable = table{
		name:          projection.KeyPrivateTable,
		instanceIDCol: projection.KeyPrivateColumnInstanceID,
	}
	KeyPrivateColID = Column{
		name:  projection.KeyPrivateColumnID,
		table: keyPrivateTable,
	}
	KeyPrivateColExpiry = Column{
		name:  projection.KeyPrivateColumnExpiry,
		table: keyPrivateTable,
	}
	KeyPrivateColKey = Column{
		name:  projection.KeyPrivateColumnKey,
		table: keyPrivateTable,
	}
)

var (
	keyPublicTable = table{
		name:          projection.KeyPublicTable,
		instanceIDCol: projection.KeyPrivateColumnInstanceID,
	}
	KeyPublicColID = Column{
		name:  projection.KeyPublicColumnID,
		table: keyPublicTable,
	}
	KeyPublicColExpiry = Column{
		name:  projection.KeyPublicColumnExpiry,
		table: keyPublicTable,
	}
	KeyPublicColKey = Column{
		name:  projection.KeyPublicColumnKey,
		table: keyPublicTable,
	}
)

func (q *Queries) ActivePublicKeys(ctx context.Context, t time.Time) (keys *PublicKeys, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	query, scan := preparePublicKeysQuery(ctx, q.client)
	if t.IsZero() {
		t = time.Now()
	}
	stmt, args, err := query.Where(
		sq.And{
			sq.Eq{KeyColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID()},
			sq.Gt{KeyPublicColExpiry.identifier(): t},
		}).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-SDFfg", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		keys, err = scan(rows)
		return err
	}, stmt, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Sghn4", "Errors.Internal")
	}

	keys.State, err = q.latestState(ctx, keyTable)
	if !zerrors.IsNotFound(err) {
		return keys, err
	}
	return keys, nil
}

func (q *Queries) ActivePrivateSigningKey(ctx context.Context, t time.Time) (keys *PrivateKeys, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	stmt, scan := preparePrivateKeysQuery(ctx, q.client)
	if t.IsZero() {
		t = time.Now()
	}
	query, args, err := stmt.Where(
		sq.And{
			sq.Eq{
				KeyColUse.identifier():        domain.KeyUsageSigning,
				KeyColInstanceID.identifier(): authz.GetInstance(ctx).InstanceID(),
			},
			sq.Gt{KeyPrivateColExpiry.identifier(): t},
		}).OrderBy(KeyPrivateColExpiry.identifier()).ToSql()
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-SDff2", "Errors.Query.SQLStatement")
	}

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		keys, err = scan(rows)
		return err
	}, query, args...)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-WRFG4", "Errors.Internal")
	}
	keys.State, err = q.latestState(ctx, keyTable)
	if !zerrors.IsNotFound(err) {
		return keys, err
	}
	return keys, nil
}

func preparePublicKeysQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*PublicKeys, error)) {
	return sq.Select(
			KeyColID.identifier(),
			KeyColCreationDate.identifier(),
			KeyColChangeDate.identifier(),
			KeyColSequence.identifier(),
			KeyColResourceOwner.identifier(),
			KeyColAlgorithm.identifier(),
			KeyColUse.identifier(),
			KeyPublicColExpiry.identifier(),
			KeyPublicColKey.identifier(),
			countColumn.identifier(),
		).From(keyTable.identifier()).
			LeftJoin(join(KeyPublicColID, KeyColID) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*PublicKeys, error) {
			keys := make([]PublicKey, 0)
			var count uint64
			for rows.Next() {
				k := new(rsaPublicKey)
				var keyValue []byte
				err := rows.Scan(
					&k.id,
					&k.creationDate,
					&k.changeDate,
					&k.sequence,
					&k.resourceOwner,
					&k.algorithm,
					&k.use,
					&k.expiry,
					&keyValue,
					&count,
				)
				if err != nil {
					return nil, err
				}
				k.publicKey, err = crypto.BytesToPublicKey(keyValue)
				if err != nil {
					return nil, err
				}
				keys = append(keys, k)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-rKd6k", "Errors.Query.CloseRows")
			}

			return &PublicKeys{
				Keys: keys,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

func preparePrivateKeysQuery(ctx context.Context, db prepareDatabase) (sq.SelectBuilder, func(*sql.Rows) (*PrivateKeys, error)) {
	return sq.Select(
			KeyColID.identifier(),
			KeyColCreationDate.identifier(),
			KeyColChangeDate.identifier(),
			KeyColSequence.identifier(),
			KeyColResourceOwner.identifier(),
			KeyColAlgorithm.identifier(),
			KeyColUse.identifier(),
			KeyPrivateColExpiry.identifier(),
			KeyPrivateColKey.identifier(),
			countColumn.identifier(),
		).From(keyTable.identifier()).
			LeftJoin(join(KeyPrivateColID, KeyColID) + db.Timetravel(call.Took(ctx))).
			PlaceholderFormat(sq.Dollar),
		func(rows *sql.Rows) (*PrivateKeys, error) {
			keys := make([]PrivateKey, 0)
			var count uint64
			for rows.Next() {
				k := new(privateKey)
				err := rows.Scan(
					&k.id,
					&k.creationDate,
					&k.changeDate,
					&k.sequence,
					&k.resourceOwner,
					&k.algorithm,
					&k.use,
					&k.expiry,
					&k.privateKey,
					&count,
				)
				if err != nil {
					return nil, err
				}
				keys = append(keys, k)
			}

			if err := rows.Close(); err != nil {
				return nil, zerrors.ThrowInternal(err, "QUERY-rKd6k", "Errors.Query.CloseRows")
			}

			return &PrivateKeys{
				Keys: keys,
				SearchResponse: SearchResponse{
					Count: count,
				},
			}, nil
		}
}

type PublicKeyReadModel struct {
	eventstore.ReadModel

	Algorithm string
	Key       *crypto.CryptoValue
	Expiry    time.Time
	Usage     domain.KeyUsage
}

func NewPublicKeyReadModel(keyID, resourceOwner string) *PublicKeyReadModel {
	return &PublicKeyReadModel{
		ReadModel: eventstore.ReadModel{
			AggregateID:   keyID,
			ResourceOwner: resourceOwner,
		},
	}
}

func (wm *PublicKeyReadModel) AppendEvents(events ...eventstore.Event) {
	wm.ReadModel.AppendEvents(events...)
}

func (wm *PublicKeyReadModel) Reduce() error {
	for _, event := range wm.Events {
		switch e := event.(type) {
		case *keypair.AddedEvent:
			wm.Algorithm = e.Algorithm
			wm.Key = e.PublicKey.Key
			wm.Expiry = e.PublicKey.Expiry
			wm.Usage = e.Usage
		default:
		}
	}
	return wm.ReadModel.Reduce()
}

func (wm *PublicKeyReadModel) Query() *eventstore.SearchQueryBuilder {
	return eventstore.NewSearchQueryBuilder(eventstore.ColumnsEvent).
		AwaitOpenTransactions().
		ResourceOwner(wm.ResourceOwner).
		AddQuery().
		AggregateTypes(keypair.AggregateType).
		AggregateIDs(wm.AggregateID).
		EventTypes(keypair.AddedEventType).
		Builder()
}

func (q *Queries) GetActivePublicKeyByID(ctx context.Context, keyID string, current time.Time) (_ PublicKey, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	model := NewPublicKeyReadModel(keyID, authz.GetInstance(ctx).InstanceID())
	if err := q.eventstore.FilterToQueryReducer(ctx, model); err != nil {
		return nil, err
	}
	if model.Algorithm == "" || model.Key == nil {
		return nil, zerrors.ThrowNotFound(err, "QUERY-Ahf7x", "Errors.Key.NotFound")
	}
	if model.Expiry.Before(current) {
		return nil, zerrors.ThrowInvalidArgument(err, "QUERY-ciF4k", "Errors.Key.ExpireBeforeNow")
	}
	keyValue, err := crypto.Decrypt(model.Key, q.keyEncryptionAlgorithm)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Ie4oh", "Errors.Internal")
	}
	publicKey, err := crypto.BytesToPublicKey(keyValue)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Kai2Z", "Errors.Internal")
	}

	return &rsaPublicKey{
		key: key{
			id:            model.AggregateID,
			creationDate:  model.CreationDate,
			changeDate:    model.ChangeDate,
			sequence:      model.ProcessedSequence,
			resourceOwner: model.ResourceOwner,
			algorithm:     model.Algorithm,
			use:           model.Usage,
		},
		expiry:    model.Expiry,
		publicKey: publicKey,
	}, nil
}
