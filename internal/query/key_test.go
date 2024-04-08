package query

import (
	"context"
	"crypto/rsa"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"io"
	"math/big"
	"regexp"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.uber.org/mock/gomock"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	key_repo "github.com/zitadel/zitadel/internal/repository/keypair"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	preparePublicKeysStmt = `SELECT projections.keys4.id,` +
		` projections.keys4.creation_date,` +
		` projections.keys4.change_date,` +
		` projections.keys4.sequence,` +
		` projections.keys4.resource_owner,` +
		` projections.keys4.algorithm,` +
		` projections.keys4.use,` +
		` projections.keys4_public.expiry,` +
		` projections.keys4_public.key,` +
		` COUNT(*) OVER ()` +
		` FROM projections.keys4` +
		` LEFT JOIN projections.keys4_public ON projections.keys4.id = projections.keys4_public.id AND projections.keys4.instance_id = projections.keys4_public.instance_id` +
		` AS OF SYSTEM TIME '-1 ms' `
	preparePublicKeysCols = []string{
		"id",
		"creation_date",
		"change_date",
		"sequence",
		"resource_owner",
		"algorithm",
		"use",
		"expiry",
		"key",
		"count",
	}

	preparePrivateKeysStmt = `SELECT projections.keys4.id,` +
		` projections.keys4.creation_date,` +
		` projections.keys4.change_date,` +
		` projections.keys4.sequence,` +
		` projections.keys4.resource_owner,` +
		` projections.keys4.algorithm,` +
		` projections.keys4.use,` +
		` projections.keys4_private.expiry,` +
		` projections.keys4_private.key,` +
		` COUNT(*) OVER ()` +
		` FROM projections.keys4` +
		` LEFT JOIN projections.keys4_private ON projections.keys4.id = projections.keys4_private.id AND projections.keys4.instance_id = projections.keys4_private.instance_id` +
		` AS OF SYSTEM TIME '-1 ms' `
)

func Test_KeyPrepares(t *testing.T) {
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
			name:    "preparePublicKeysQuery no result",
			prepare: preparePublicKeysQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(preparePublicKeysStmt),
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
			object: &PublicKeys{Keys: []PublicKey{}},
		},
		{
			name:    "preparePublicKeysQuery found",
			prepare: preparePublicKeysQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(preparePublicKeysStmt),
					preparePublicKeysCols,
					[][]driver.Value{
						{
							"key-id",
							testNow,
							testNow,
							uint64(20211109),
							"ro",
							"RS256",
							0,
							testNow,
							[]byte("-----BEGIN RSA PUBLIC KEY-----\nMIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAsvX9P58JFxEs5C+L+H7W\nduFSWL5EPzber7C2m94klrSV6q0bAcrYQnGwFOlveThsY200hRbadKaKjHD7qIKH\nDEe0IY2PSRht33Jye52AwhkRw+M3xuQH/7R8LydnsNFk2KHpr5X2SBv42e37LjkE\nslKSaMRgJW+v0KZ30piY8QsdFRKKaVg5/Ajt1YToM1YVsdHXJ3vmXFMtypLdxwUD\ndIaLEX6pFUkU75KSuEQ/E2luT61Q3ta9kOWm9+0zvi7OMcbdekJT7mzcVnh93R1c\n13ZhQCLbh9A7si8jKFtaMWevjayrvqQABEcTN9N4Hoxcyg6l4neZtRDk75OMYcqm\nDQIDAQAB\n-----END RSA PUBLIC KEY-----\n"),
						},
					},
				),
			},
			object: &PublicKeys{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Keys: []PublicKey{
					&rsaPublicKey{
						key: key{
							id:            "key-id",
							creationDate:  testNow,
							changeDate:    testNow,
							sequence:      20211109,
							resourceOwner: "ro",
							algorithm:     "RS256",
							use:           domain.KeyUsageSigning,
						},
						expiry: testNow,
						publicKey: &rsa.PublicKey{
							E: 65537,
							N: fromBase16("b2f5fd3f9f0917112ce42f8bf87ed676e15258be443f36deafb0b69bde2496b495eaad1b01cad84271b014e96f79386c636d348516da74a68a8c70fba882870c47b4218d8f49186ddf72727b9d80c21911c3e337c6e407ffb47c2f2767b0d164d8a1e9af95f6481bf8d9edfb2e3904b2529268c460256fafd0a677d29898f10b1d15128a695839fc08edd584e8335615b1d1d7277be65c532dca92ddc7050374868b117ea9154914ef9292b8443f13696e4fad50ded6bd90e5a6f7ed33be2ece31c6dd7a4253ee6cdc56787ddd1d5cd776614022db87d03bb22f23285b5a3167af8dacabbea40004471337d3781e8c5cca0ea5e27799b510e4ef938c61caa60d"),
						},
					},
				},
			},
		},
		{
			name:    "preparePublicKeysQuery sql err",
			prepare: preparePublicKeysQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(preparePublicKeysStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*PublicKeys)(nil),
		},
		{
			name:    "preparePrivateKeysQuery no result",
			prepare: preparePrivateKeysQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(preparePrivateKeysStmt),
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
			object: &PrivateKeys{Keys: []PrivateKey{}},
		},
		{
			name:    "preparePrivateKeysQuery found",
			prepare: preparePrivateKeysQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(preparePrivateKeysStmt),
					preparePublicKeysCols,
					[][]driver.Value{
						{
							"key-id",
							testNow,
							testNow,
							uint64(20211109),
							"ro",
							"RS256",
							0,
							testNow,
							[]byte(`{"Algorithm": "enc", "Crypted": "cHJpdmF0ZUtleQ==", "CryptoType": 0, "KeyID": "id"}`),
						},
					},
				),
			},
			object: &PrivateKeys{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Keys: []PrivateKey{
					&privateKey{
						key: key{
							id:            "key-id",
							creationDate:  testNow,
							changeDate:    testNow,
							sequence:      20211109,
							resourceOwner: "ro",
							algorithm:     "RS256",
							use:           domain.KeyUsageSigning,
						},
						expiry: testNow,
						privateKey: &crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "enc",
							KeyID:      "id",
							Crypted:    []byte("privateKey"),
						},
					},
				},
			},
		},
		{
			name:    "preparePrivateKeysQuery sql err",
			prepare: preparePrivateKeysQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(preparePrivateKeysStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*PrivateKeys)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}

func fromBase16(base16 string) *big.Int {
	i, ok := new(big.Int).SetString(base16, 16)
	if !ok {
		panic("bad number: " + base16)
	}
	return i
}

const pubKey = `-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAs38btwb3c7r0tMaQpGvB
mY+mPwMU/LpfuPoC0k2t4RsKp0fv40SMl50CRrHgk395wch8PMPYbl3+8TtYAJuy
rFALIj3Ff1UcKIk0hOH5DDsfh7/q2wFuncTmS6bifYo8CfSq2vDGnM7nZnEvxY/M
fSydZdcmIqlkUpfQmtzExw9+tSe5Dxq6gn5JtlGgLgZGt69r5iMMrTEGhhVAXzNu
MZbmlCoBru+rC8ITlTX/0V1ZcsSbL8tYWhthyu9x6yjo1bH85wiVI4gs0MhU8f2a
+kjL/KGZbR14Ua2eo6tonBZLC5DHWM2TkYXgRCDPufjcgmzN0Lm91E4P8KvBcvly
6QIDAQAB
-----END PUBLIC KEY-----
`

func TestQueries_GetPublicKeyByID(t *testing.T) {
	now := time.Now()
	future := now.Add(time.Hour)

	tests := []struct {
		name       string
		eventstore func(*testing.T) *eventstore.Eventstore
		encryption func(*testing.T) *crypto.MockEncryptionAlgorithm
		want       *rsaPublicKey
		wantErr    error
	}{
		{
			name: "filter error",
			eventstore: expectEventstore(
				expectFilterError(io.ErrClosedPipe),
			),
			wantErr: io.ErrClosedPipe,
		},
		{
			name: "not found error",
			eventstore: expectEventstore(
				expectFilter(),
			),
			wantErr: zerrors.ThrowNotFound(nil, "QUERY-Ahf7x", "Errors.Key.NotFound"),
		},
		{
			name: "decrypt error",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(key_repo.NewAddedEvent(context.Background(),
						&eventstore.Aggregate{
							ID:            "keyID",
							Type:          key_repo.AggregateType,
							ResourceOwner: "instanceID",
							InstanceID:    "instanceID",
							Version:       key_repo.AggregateVersion,
						},
						domain.KeyUsageSigning, "alg",
						&crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "alg",
							KeyID:      "keyID",
							Crypted:    []byte("private"),
						},
						&crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "alg",
							KeyID:      "keyID",
							Crypted:    []byte("public"),
						},
						future,
						future,
					)),
				),
			),
			encryption: func(t *testing.T) *crypto.MockEncryptionAlgorithm {
				encryption := crypto.NewMockEncryptionAlgorithm(gomock.NewController(t))
				expect := encryption.EXPECT()
				expect.Algorithm().Return("alg")
				expect.DecryptionKeyIDs().Return([]string{})
				return encryption
			},
			wantErr: zerrors.ThrowInternal(nil, "QUERY-Ie4oh", "Errors.Internal"),
		},
		{
			name: "parse error",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(key_repo.NewAddedEvent(context.Background(),
						&eventstore.Aggregate{
							ID:            "keyID",
							Type:          key_repo.AggregateType,
							ResourceOwner: "instanceID",
							InstanceID:    "instanceID",
							Version:       key_repo.AggregateVersion,
						},
						domain.KeyUsageSigning, "alg",
						&crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "alg",
							KeyID:      "keyID",
							Crypted:    []byte("private"),
						},
						&crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "alg",
							KeyID:      "keyID",
							Crypted:    []byte("public"),
						},
						future,
						future,
					)),
				),
			),
			encryption: func(t *testing.T) *crypto.MockEncryptionAlgorithm {
				encryption := crypto.NewMockEncryptionAlgorithm(gomock.NewController(t))
				expect := encryption.EXPECT()
				expect.Algorithm().Return("alg")
				expect.DecryptionKeyIDs().Return([]string{"keyID"})
				expect.Decrypt([]byte("public"), "keyID").Return([]byte("foo"), nil)
				return encryption
			},
			wantErr: zerrors.ThrowInternal(nil, "QUERY-Kai2Z", "Errors.Internal"),
		},
		{
			name: "success",
			eventstore: expectEventstore(
				expectFilter(
					eventFromEventPusher(key_repo.NewAddedEvent(context.Background(),
						&eventstore.Aggregate{
							ID:            "keyID",
							Type:          key_repo.AggregateType,
							ResourceOwner: "instanceID",
							InstanceID:    "instanceID",
							Version:       key_repo.AggregateVersion,
						},
						domain.KeyUsageSigning, "alg",
						&crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "alg",
							KeyID:      "keyID",
							Crypted:    []byte("private"),
						},
						&crypto.CryptoValue{
							CryptoType: crypto.TypeEncryption,
							Algorithm:  "alg",
							KeyID:      "keyID",
							Crypted:    []byte("public"),
						},
						future,
						future,
					)),
				),
			),
			encryption: func(t *testing.T) *crypto.MockEncryptionAlgorithm {
				encryption := crypto.NewMockEncryptionAlgorithm(gomock.NewController(t))
				expect := encryption.EXPECT()
				expect.Algorithm().Return("alg")
				expect.DecryptionKeyIDs().Return([]string{"keyID"})
				expect.Decrypt([]byte("public"), "keyID").Return([]byte(pubKey), nil)
				return encryption
			},
			want: &rsaPublicKey{
				key: key{
					id:            "keyID",
					resourceOwner: "instanceID",
					algorithm:     "alg",
					use:           domain.KeyUsageSigning,
				},
				expiry: future,
				publicKey: func() *rsa.PublicKey {
					publicKey, err := crypto.BytesToPublicKey([]byte(pubKey))
					if err != nil {
						panic(err)
					}
					return publicKey
				}(),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			q := &Queries{
				eventstore: tt.eventstore(t),
			}
			if tt.encryption != nil {
				q.keyEncryptionAlgorithm = tt.encryption(t)
			}
			ctx := authz.NewMockContext("instanceID", "orgID", "loginClient")
			key, err := q.GetPublicKeyByID(ctx, "keyID")
			if tt.wantErr != nil {
				require.ErrorIs(t, err, tt.wantErr)
				return
			}
			require.NoError(t, err)
			require.NotNil(t, key)

			got := key.(*rsaPublicKey)
			assert.WithinDuration(t, tt.want.expiry, got.expiry, time.Second)
			tt.want.expiry = time.Time{}
			got.expiry = time.Time{}
			assert.Equal(t, tt.want, got)
		})
	}
}
