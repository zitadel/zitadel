package query

import (
	"context"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/x509"
	"database/sql"
	"database/sql/driver"
	"encoding/json"
	"io"
	"regexp"
	"strconv"
	"testing"
	"time"

	"github.com/go-jose/go-jose/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/repository/webkey"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestQueries_GetPublicWebKeyByID(t *testing.T) {
	ctx := authz.NewMockContextWithPermissions("instance1", "org1", "user1", nil)
	key, err := ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	require.NoError(t, err)

	type fields struct {
		eventstore func(*testing.T) *eventstore.Eventstore
	}
	type args struct {
		keyID string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *jose.JSONWebKey
		wantErr error
	}{
		{
			name: "filter error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilterError(io.ErrClosedPipe),
				),
			},
			args:    args{"key1"},
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(),
				),
			},
			args:    args{"key1"},
			wantErr: zerrors.ThrowNotFound(nil, "QUERY-AiCh0", "Errors.WebKey.NotFound"),
		},
		{
			name: "removed, not found error",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(mustNewWebkeyAddedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
							&jose.JSONWebKey{
								Key:       &key.PublicKey,
								KeyID:     "key1",
								Algorithm: string(jose.ES384),
								Use:       crypto.KeyUsageSigning.String(),
							},
							&crypto.WebKeyECDSAConfig{
								Curve: crypto.EllipticCurveP384,
							},
						)),
						eventFromEventPusher(webkey.NewRemovedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
						)),
					),
				),
			},
			args:    args{"key1"},
			wantErr: zerrors.ThrowNotFound(nil, "QUERY-AiCh0", "Errors.WebKey.NotFound"),
		},
		{
			name: "ok",
			fields: fields{
				eventstore: expectEventstore(
					expectFilter(
						eventFromEventPusher(mustNewWebkeyAddedEvent(ctx,
							webkey.NewAggregate("key1", "instance1"),
							&crypto.CryptoValue{
								CryptoType: crypto.TypeEncryption,
								Algorithm:  "alg",
								KeyID:      "encKey",
								Crypted:    []byte("crypted"),
							},
							&jose.JSONWebKey{
								Key:       &key.PublicKey,
								KeyID:     "key1",
								Algorithm: string(jose.ES384),
								Use:       crypto.KeyUsageSigning.String(),
							},
							&crypto.WebKeyECDSAConfig{
								Curve: crypto.EllipticCurveP384,
							},
						)),
					),
				),
			},
			args: args{"key1"},
			want: &jose.JSONWebKey{
				Key:                         &key.PublicKey,
				KeyID:                       "key1",
				Algorithm:                   string(jose.ES384),
				Use:                         crypto.KeyUsageSigning.String(),
				Certificates:                []*x509.Certificate{},
				CertificateThumbprintSHA1:   []byte{},
				CertificateThumbprintSHA256: []byte{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queries{
				eventstore: tt.fields.eventstore(t),
			}
			got, err := q.GetPublicWebKeyByID(ctx, tt.args.keyID)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func mustNewWebkeyAddedEvent(
	ctx context.Context,
	aggregate *eventstore.Aggregate,
	privateKey *crypto.CryptoValue,
	publicKey *jose.JSONWebKey,
	config crypto.WebKeyConfig) *webkey.AddedEvent {
	event, err := webkey.NewAddedEvent(ctx, aggregate, privateKey, publicKey, config)
	if err != nil {
		panic(err)
	}
	return event
}

func TestQueries_GetActiveSigningWebKey(t *testing.T) {
	ctx := authz.NewMockContextWithPermissions("instance1", "org1", "user1", nil)
	expQuery := regexp.QuoteMeta(webKeyByStateQuery)
	queryArgs := []driver.Value{"instance1", domain.WebKeyStateActive}
	cols := []string{"private_key"}

	alg := crypto.CreateMockEncryptionAlg(gomock.NewController(t))
	encryptedPrivate, _, err := crypto.GenerateEncryptedWebKey("key1", alg, &crypto.WebKeyED25519Config{})
	require.NoError(t, err)

	var expectedWebKey *jose.JSONWebKey
	err = crypto.DecryptJSON(encryptedPrivate, &expectedWebKey, alg)
	require.NoError(t, err)

	tests := []struct {
		name    string
		mock    sqlExpectation
		want    *jose.JSONWebKey
		wantErr error
	}{
		{
			name:    "no active error",
			mock:    mockQueryErr(expQuery, sql.ErrNoRows, queryArgs...),
			wantErr: zerrors.ThrowInternal(sql.ErrNoRows, "QUERY-Opoh7", "Errors.WebKey.NoActive"),
		},
		{
			name:    "internal error",
			mock:    mockQueryErr(expQuery, sql.ErrConnDone, queryArgs...),
			wantErr: zerrors.ThrowInternal(sql.ErrConnDone, "QUERY-Shoo0", "Errors.Internal"),
		},
		{
			name:    "invalid crypto value error",
			mock:    mockQuery(expQuery, cols, []driver.Value{&crypto.CryptoValue{}}, queryArgs...),
			wantErr: zerrors.ThrowInvalidArgument(nil, "CRYPT-Nx7XlT", "value was encrypted with a different key"),
		},
		{
			name: "found, ok",
			mock: mockQuery(expQuery, cols, []driver.Value{encryptedPrivate}, queryArgs...),
			want: expectedWebKey,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execMock(t, tt.mock, func(db *sql.DB) {
				q := &Queries{
					client: &database.DB{
						DB: db,
					},
					keyEncryptionAlgorithm: alg,
				}
				got, err := q.GetActiveSigningWebKey(ctx)
				require.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, tt.want, got)
			})
		})
	}
}

func TestQueries_ListWebKeys(t *testing.T) {
	ctx := authz.NewMockContextWithPermissions("instance1", "org1", "user1", nil)
	expQuery := regexp.QuoteMeta(webKeyListQuery)
	queryArgs := []driver.Value{"instance1"}
	cols := []string{"key_id", "creation_date", "change_date", "sequence", "state", "config", "config_type"}

	webKeyConfig := &crypto.WebKeyRSAConfig{
		Bits:   crypto.RSABits4096,
		Hasher: crypto.RSAHasherSHA512,
	}
	webKeyConfigJSON, err := json.Marshal(webKeyConfig)
	require.NoError(t, err)

	tests := []struct {
		name    string
		mock    sqlExpectation
		want    []WebKeyDetails
		wantErr error
	}{
		{
			name:    "internal error",
			mock:    mockQueryErr(expQuery, sql.ErrConnDone, queryArgs...),
			wantErr: zerrors.ThrowInternal(sql.ErrConnDone, "QUERY-Ohl3A", "Errors.Internal"),
		},
		{
			name: "invalid json error",
			mock: mockQueriesScanErr(expQuery, cols, [][]driver.Value{
				{
					"key1",
					time.Unix(1, 2),
					time.Unix(3, 4),
					1,
					domain.WebKeyStateActive,
					"~~~~~",
					crypto.WebKeyConfigTypeRSA,
				},
			}, queryArgs...),
			wantErr: zerrors.ThrowInternal(err, "QUERY-Ohl3A", "Errors.Internal"),
		},
		{
			name: "ok",
			mock: mockQueries(expQuery, cols, [][]driver.Value{
				{
					"key1",
					time.Unix(1, 2),
					time.Unix(3, 4),
					1,
					domain.WebKeyStateActive,
					webKeyConfigJSON,
					crypto.WebKeyConfigTypeRSA,
				},
				{
					"key2",
					time.Unix(5, 6),
					time.Unix(7, 8),
					2,
					domain.WebKeyStateInitial,
					webKeyConfigJSON,
					crypto.WebKeyConfigTypeRSA,
				},
			}, queryArgs...),
			want: []WebKeyDetails{
				{
					KeyID:        "key1",
					CreationDate: time.Unix(1, 2),
					ChangeDate:   time.Unix(3, 4),
					Sequence:     1,
					State:        domain.WebKeyStateActive,
					Config:       webKeyConfig,
				},
				{
					KeyID:        "key2",
					CreationDate: time.Unix(5, 6),
					ChangeDate:   time.Unix(7, 8),
					Sequence:     2,
					State:        domain.WebKeyStateInitial,
					Config:       webKeyConfig,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execMock(t, tt.mock, func(db *sql.DB) {
				q := &Queries{
					client: &database.DB{
						DB: db,
					},
				}
				got, err := q.ListWebKeys(ctx)
				require.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, tt.want, got)
			})
		})
	}
}

func TestQueries_GetWebKeySet(t *testing.T) {
	ctx := authz.NewMockContextWithPermissions("instance1", "org1", "user1", nil)
	expQuery := regexp.QuoteMeta(webKeyPublicKeysQuery)
	queryArgs := []driver.Value{"instance1"}
	cols := []string{"public_key"}

	alg := crypto.CreateMockEncryptionAlg(gomock.NewController(t))
	conf := &crypto.WebKeyED25519Config{}
	expectedKeySet := &jose.JSONWebKeySet{
		Keys: make([]jose.JSONWebKey, 3),
	}
	expectedRows := make([][]driver.Value, 3)

	for i := 0; i < 3; i++ {
		_, pubKey, err := crypto.GenerateEncryptedWebKey(strconv.Itoa(i), alg, conf)
		require.NoError(t, err)
		pubKeyJSON, err := json.Marshal(pubKey)
		require.NoError(t, err)
		err = json.Unmarshal(pubKeyJSON, &expectedKeySet.Keys[i])
		require.NoError(t, err)
		expectedRows[i] = []driver.Value{pubKeyJSON}
	}

	tests := []struct {
		name    string
		mock    sqlExpectation
		want    *jose.JSONWebKeySet
		wantErr error
	}{
		{
			name:    "internal error",
			mock:    mockQueryErr(expQuery, sql.ErrConnDone, queryArgs...),
			wantErr: zerrors.ThrowInternal(sql.ErrConnDone, "QUERY-Eeng7", "Errors.Internal"),
		},
		{
			name:    "invalid json error",
			mock:    mockQueriesScanErr(expQuery, cols, [][]driver.Value{{"~~~"}}, queryArgs...),
			wantErr: zerrors.ThrowInternal(nil, "QUERY-Eeng7", "Errors.Internal"),
		},
		{
			name: "ok",
			mock: mockQueries(expQuery, cols, expectedRows, queryArgs...),
			want: expectedKeySet,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execMock(t, tt.mock, func(db *sql.DB) {
				q := &Queries{
					client: &database.DB{
						DB: db,
					},
				}
				got, err := q.GetWebKeySet(ctx)
				require.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, tt.want, got)
			})
		})
	}
}
