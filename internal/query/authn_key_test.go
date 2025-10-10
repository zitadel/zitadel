package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	prepareAuthNKeysStmt = `SELECT projections.authn_keys2.id,` +
		` projections.authn_keys2.aggregate_id,` +
		` projections.authn_keys2.creation_date,` +
		` projections.authn_keys2.change_date,` +
		` projections.authn_keys2.resource_owner,` +
		` projections.authn_keys2.sequence,` +
		` projections.authn_keys2.expiration,` +
		` projections.authn_keys2.type,` +
		` projections.authn_keys2.object_id,` +
		` COUNT(*) OVER ()` +
		` FROM projections.authn_keys2`
	prepareAuthNKeysCols = []string{
		"id",
		"aggregate_id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"expiration",
		"type",
		"object_id",
		"count",
	}

	prepareAuthNKeysDataStmt = `SELECT projections.authn_keys2.id,` +
		` projections.authn_keys2.creation_date,` +
		` projections.authn_keys2.change_date,` +
		` projections.authn_keys2.resource_owner,` +
		` projections.authn_keys2.sequence,` +
		` projections.authn_keys2.expiration,` +
		` projections.authn_keys2.type,` +
		` projections.authn_keys2.identifier,` +
		` projections.authn_keys2.public_key,` +
		` COUNT(*) OVER ()` +
		` FROM projections.authn_keys2`
	prepareAuthNKeysDataCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"expiration",
		"type",
		"identifier",
		"public_key",
		"count",
	}

	prepareAuthNKeyStmt = `SELECT projections.authn_keys2.id,` +
		` projections.authn_keys2.creation_date,` +
		` projections.authn_keys2.change_date,` +
		` projections.authn_keys2.resource_owner,` +
		` projections.authn_keys2.sequence,` +
		` projections.authn_keys2.expiration,` +
		` projections.authn_keys2.type` +
		` FROM projections.authn_keys2`
	prepareAuthNKeyCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"expiration",
		"type",
	}

	prepareAuthNKeyPublicKeyStmt = `SELECT projections.authn_keys2.public_key` +
		` FROM projections.authn_keys2`
	prepareAuthNKeyPublicKeyCols = []string{
		"public_key",
	}
)

func Test_AuthNKeyPrepares(t *testing.T) {
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
	}
	tests := []struct {
		name    string
		prepare interface{}
		want    want
		object  interface{}
	}{
		{
			name:    "prepareAuthNKeysQuery no result",
			prepare: prepareAuthNKeysQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareAuthNKeysStmt),
					nil,
					nil,
				),
			},
			object: &AuthNKeys{AuthNKeys: []*AuthNKey{}},
		},
		{
			name:    "prepareAuthNKeysQuery one result",
			prepare: prepareAuthNKeysQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareAuthNKeysStmt),
					prepareAuthNKeysCols,
					[][]driver.Value{
						{
							"id",
							"aggId",
							testNow,
							testNow,
							"ro",
							uint64(20211109),
							testNow,
							1,
							"app1",
						},
					},
				),
			},
			object: &AuthNKeys{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				AuthNKeys: []*AuthNKey{
					{
						ID:            "id",
						AggregateID:   "aggId",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						Sequence:      20211109,
						Expiration:    testNow,
						Type:          domain.AuthNKeyTypeJSON,
						ApplicationID: "app1",
					},
				},
			},
		},
		{
			name:    "prepareAuthNKeysQuery multiple result",
			prepare: prepareAuthNKeysQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareAuthNKeysStmt),
					prepareAuthNKeysCols,
					[][]driver.Value{
						{
							"id-1",
							"aggId-1",
							testNow,
							testNow,
							"ro",
							uint64(20211109),
							testNow,
							1,
							"app1",
						},
						{
							"id-2",
							"aggId-2",
							testNow,
							testNow,
							"ro",
							uint64(20211109),
							testNow,
							1,
							"app1",
						},
					},
				),
			},
			object: &AuthNKeys{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				AuthNKeys: []*AuthNKey{
					{
						ID:            "id-1",
						AggregateID:   "aggId-1",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						Sequence:      20211109,
						Expiration:    testNow,
						Type:          domain.AuthNKeyTypeJSON,
						ApplicationID: "app1",
					},
					{
						ID:            "id-2",
						AggregateID:   "aggId-2",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						Sequence:      20211109,
						Expiration:    testNow,
						Type:          domain.AuthNKeyTypeJSON,
						ApplicationID: "app1",
					},
				},
			},
		},
		{
			name:    "prepareAuthNKeysQuery sql err",
			prepare: prepareAuthNKeysQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareAuthNKeysStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: nil,
		},
		{
			name:    "prepareAuthNKeysDataQuery no result",
			prepare: prepareAuthNKeysDataQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareAuthNKeysDataStmt),
					nil,
					nil,
				),
			},
			object: &AuthNKeysData{AuthNKeysData: []*AuthNKeyData{}},
		},
		{
			name:    "prepareAuthNKeysDataQuery one result",
			prepare: prepareAuthNKeysDataQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareAuthNKeysDataStmt),
					prepareAuthNKeysDataCols,
					[][]driver.Value{
						{
							"id",
							testNow,
							testNow,
							"ro",
							uint64(20211109),
							testNow,
							1,
							"identifier",
							[]byte("public"),
						},
					},
				),
			},
			object: &AuthNKeysData{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				AuthNKeysData: []*AuthNKeyData{
					{
						ID:            "id",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						Sequence:      20211109,
						Expiration:    testNow,
						Type:          domain.AuthNKeyTypeJSON,
						Identifier:    "identifier",
						PublicKey:     []byte("public"),
					},
				},
			},
		},
		{
			name:    "prepareAuthNKeysDataQuery multiple result",
			prepare: prepareAuthNKeysDataQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareAuthNKeysDataStmt),
					prepareAuthNKeysDataCols,
					[][]driver.Value{
						{
							"id-1",
							testNow,
							testNow,
							"ro",
							uint64(20211109),
							testNow,
							1,
							"identifier1",
							[]byte("public1"),
						},
						{
							"id-2",
							testNow,
							testNow,
							"ro",
							uint64(20211109),
							testNow,
							1,
							"identifier2",
							[]byte("public2"),
						},
					},
				),
			},
			object: &AuthNKeysData{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				AuthNKeysData: []*AuthNKeyData{
					{
						ID:            "id-1",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						Sequence:      20211109,
						Expiration:    testNow,
						Type:          domain.AuthNKeyTypeJSON,
						Identifier:    "identifier1",
						PublicKey:     []byte("public1"),
					},
					{
						ID:            "id-2",
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "ro",
						Sequence:      20211109,
						Expiration:    testNow,
						Type:          domain.AuthNKeyTypeJSON,
						Identifier:    "identifier2",
						PublicKey:     []byte("public2"),
					},
				},
			},
		},
		{
			name:    "prepareAuthNKeysDataQuery sql err",
			prepare: prepareAuthNKeysDataQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareAuthNKeysDataStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*AuthNKey)(nil),
		},
		{
			name:    "prepareAuthNKeyQuery no result",
			prepare: prepareAuthNKeyQuery,
			want: want{
				sqlExpectations: mockQueriesScanErr(
					regexp.QuoteMeta(prepareAuthNKeyStmt),
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !zerrors.IsNotFound(err) {
						return fmt.Errorf("err should be zitadel.NotFoundError got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*AuthNKey)(nil),
		},
		{
			name:    "prepareAuthNKeyQuery found",
			prepare: prepareAuthNKeyQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(prepareAuthNKeyStmt),
					prepareAuthNKeyCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						"ro",
						uint64(20211109),
						testNow,
						1,
					},
				),
			},
			object: &AuthNKey{
				ID:            "id",
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "ro",
				Sequence:      20211109,
				Expiration:    testNow,
				Type:          domain.AuthNKeyTypeJSON,
			},
		},
		{
			name:    "prepareAuthNKeyQuery sql err",
			prepare: prepareAuthNKeyQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareAuthNKeyStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*AuthNKey)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}

func TestQueries_GetAuthNKeyUser(t *testing.T) {
	expQuery := regexp.QuoteMeta(authNKeyUserQuery)
	cols := []string{"user_id", "resource_owner", "username", "access_token_type", "public_key"}
	pubkey := []byte(`-----BEGIN RSA PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEA2ufAL1b72bIy1ar+Ws6b
GohJJQFB7dfRapDqeqM8Ukp6CVdPzq/pOz1viAq50yzWZJryF+2wshFAKGF9A2/B
2Yf9bJXPZ/KbkFrYT3NTvYDkvlaSTl9mMnzrU29s48F1PTWKfB+C3aMsOEG1BufV
s63qF4nrEPjSbhljIco9FZq4XppIzhMQ0fDdA/+XygCJqvuaL0LibM1KrlUdnu71
YekhSJjEPnvOisXIk4IXywoGIOwtjxkDvNItQvaMVldr4/kb6uvbgdWwq5EwBZXq
low2kyJov38V4Uk2I8kuXpLcnrpw5Tio2ooiUE27b0vHZqBKOei9Uo88qCrn3EKx
6QIDAQAB
-----END RSA PUBLIC KEY-----`)

	tests := []struct {
		name    string
		mock    sqlExpectation
		want    *AuthNKeyUser
		wantErr error
	}{
		{
			name:    "no rows",
			mock:    mockQueryErr(expQuery, sql.ErrNoRows, "instanceID", "keyID", "userID"),
			wantErr: zerrors.ThrowNotFound(sql.ErrNoRows, "QUERY-Tha6f", "Errors.AuthNKey.NotFound"),
		},
		{
			name:    "internal error",
			mock:    mockQueryErr(expQuery, sql.ErrConnDone, "instanceID", "keyID", "userID"),
			wantErr: zerrors.ThrowInternal(sql.ErrConnDone, "QUERY-aen2A", "Errors.Internal"),
		},
		{
			name: "success",
			mock: mockQuery(expQuery, cols,
				[]driver.Value{"userID", "orgID", "username", domain.OIDCTokenTypeJWT, pubkey},
				"instanceID", "keyID", "userID",
			),
			want: &AuthNKeyUser{
				UserID:        "userID",
				ResourceOwner: "orgID",
				Username:      "username",
				TokenType:     domain.OIDCTokenTypeJWT,
				PublicKey:     pubkey,
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
				ctx := authz.NewMockContext("instanceID", "orgID", "userID")
				got, err := q.GetAuthNKeyUser(ctx, "keyID", "userID")
				require.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, tt.want, got)
			})
		})
	}
}
