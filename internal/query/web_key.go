package query

import (
	"context"
	"database/sql"
	_ "embed"
	"encoding/json"
	"errors"
	"time"

	"github.com/go-jose/go-jose/v4"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	//go:embed web_key_by_state.sql
	webKeyByStateQuery string
	//go:embed web_key_list.sql
	webKeyListQuery string
	//go:embed web_key_public_keys.sql
	webKeyPublicKeysQuery string
)

// GetPublicWebKeyByID gets a public key by it's keyID directly from the eventstore.
func (q *Queries) GetPublicWebKeyByID(ctx context.Context, keyID string) (webKey *jose.JSONWebKey, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	model := NewWebKeyReadModel(keyID, authz.GetInstance(ctx).InstanceID())
	if err = q.eventstore.FilterToQueryReducer(ctx, model); err != nil {
		return nil, err
	}
	if model.State == domain.WebKeyStateUnspecified || model.State == domain.WebKeyStateRemoved {
		return nil, zerrors.ThrowNotFound(nil, "QUERY-AiCh0", "Errors.WebKey.NotFound")
	}
	return model.PublicKey, nil
}

// GetActiveSigningWebKey gets the current active signing key from the web_keys projection.
// The active signing key is eventual consistent.
func (q *Queries) GetActiveSigningWebKey(ctx context.Context) (webKey *jose.JSONWebKey, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	var keyValue *crypto.CryptoValue
	err = q.client.QueryRowContext(ctx, func(row *sql.Row) error {
		return row.Scan(&keyValue)
	},
		webKeyByStateQuery,
		authz.GetInstance(ctx).InstanceID(),
		domain.WebKeyStateActive,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, zerrors.ThrowInternal(err, "QUERY-Opoh7", "Errors.WebKey.NoActive")
		}
		return nil, zerrors.ThrowInternal(err, "QUERY-Shoo0", "Errors.Internal")
	}
	if err = crypto.DecryptJSON(keyValue, &webKey, q.keyEncryptionAlgorithm); err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Iuk0s", "Errors.Internal")
	}
	return webKey, nil
}

type WebKeyDetails struct {
	KeyID        string
	CreationDate time.Time
	ChangeDate   time.Time
	Sequence     int64
	State        domain.WebKeyState
	Config       crypto.WebKeyConfig
}

type WebKeyList struct {
	Keys []WebKeyDetails
}

// ListWebKeys gets a list of [WebKeyDetails] for the complete instance from the web_keys projection.
// The list is eventual consistent.
func (q *Queries) ListWebKeys(ctx context.Context) (list []WebKeyDetails, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		for rows.Next() {
			var (
				configData []byte
				configType crypto.WebKeyConfigType
			)
			var details WebKeyDetails
			if err = rows.Scan(
				&details.KeyID,
				&details.CreationDate,
				&details.ChangeDate,
				&details.Sequence,
				&details.State,
				&configData,
				&configType,
			); err != nil {
				return err
			}
			details.Config, err = crypto.UnmarshalWebKeyConfig(configData, configType)
			if err != nil {
				return err
			}
			list = append(list, details)
		}
		return rows.Err()
	},
		webKeyListQuery,
		authz.GetInstance(ctx).InstanceID(),
	)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Ohl3A", "Errors.Internal")
	}
	return list, nil
}

// GetWebKeySet gets a JSON Web Key set from the web_keys projection.
// The set contains all existing public keys for the instance.
// The set is eventual consistent.
func (q *Queries) GetWebKeySet(ctx context.Context) (_ *jose.JSONWebKeySet, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	var keys []jose.JSONWebKey

	err = q.client.QueryContext(ctx, func(rows *sql.Rows) error {
		for rows.Next() {
			var webKeyData []byte
			if err = rows.Scan(&webKeyData); err != nil {
				return err
			}
			var webKey jose.JSONWebKey
			if err = json.Unmarshal(webKeyData, &webKey); err != nil {
				return err
			}
			keys = append(keys, webKey)
		}
		return rows.Err()
	},
		webKeyPublicKeysQuery,
		authz.GetInstance(ctx).InstanceID(),
	)
	if err != nil {
		return nil, zerrors.ThrowInternal(err, "QUERY-Eeng7", "Errors.Internal")
	}
	return &jose.JSONWebKeySet{Keys: keys}, nil
}
// GetActiveWebKey gets the current active signing key with caching.
func (q *Queries) GetActiveWebKey(ctx context.Context) (webKey *jose.JSONWebKey, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	instanceID := authz.GetInstance(ctx).InstanceID()

	// Try cache first
	if q.caches != nil && q.caches.webkeyActiveSigningKey != nil {
		cacheEntry, ok := q.caches.webkeyActiveSigningKey.Get(ctx, webkeyActiveSigningKeyCacheInstanceIndex, instanceID)
		if ok {
			return cacheEntry.webKey, nil
		}
	}

	// Cache miss - fetch from database using existing method
	webKey, err = q.GetActiveSigningWebKey(ctx)
	if err != nil {
		return nil, err
	}

	// Set cache
	if q.caches != nil && q.caches.webkeyActiveSigningKey != nil {
		entry := &webkeyActiveSigningKeyCacheEntry{
			instanceID: instanceID,
			webKey:     webKey,
		}
		q.caches.webkeyActiveSigningKey.Set(ctx, entry)
	}

	return webKey, nil
}

// GetAllPublicWebKeys gets all public web keys with caching.
func (q *Queries) GetAllPublicWebKeys(ctx context.Context) (_ *jose.JSONWebKeySet, err error) {
	ctx, span := tracing.NewSpan(ctx)
	defer func() { span.EndWithError(err) }()

	instanceID := authz.GetInstance(ctx).InstanceID()

	// Try cache first
	if q.caches != nil && q.caches.webkeyPublicKeys != nil {
		cacheEntry, ok := q.caches.webkeyPublicKeys.Get(ctx, webkeyPublicKeysCacheInstanceIndex, instanceID)
		if ok {
			return cacheEntry.keySet, nil
		}
	}

	// Cache miss - fetch from database using existing method
	keySet, err := q.GetWebKeySet(ctx)
	if err != nil {
		return nil, err
	}

	// Set cache
	if q.caches != nil && q.caches.webkeyPublicKeys != nil {
		entry := &webkeyPublicKeysCacheEntry{
			instanceID: instanceID,
			keySet:     keySet,
		}
		q.caches.webkeyPublicKeys.Set(ctx, entry)
	}

	return keySet, nil
}