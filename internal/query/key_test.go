package query

import (
	"crypto/rsa"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"math/big"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	errs "github.com/zitadel/zitadel/internal/errors"
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
					regexp.QuoteMeta(`SELECT projections.keys.id,`+
						` projections.keys.creation_date,`+
						` projections.keys.change_date,`+
						` projections.keys.sequence,`+
						` projections.keys.resource_owner,`+
						` projections.keys.algorithm,`+
						` projections.keys.use,`+
						` projections.keys_public.expiry,`+
						` projections.keys_public.key,`+
						` COUNT(*) OVER ()`+
						` FROM projections.keys`+
						` LEFT JOIN projections.keys_public ON projections.keys.id = projections.keys_public.id`),
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !errs.IsNotFound(err) {
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
					regexp.QuoteMeta(`SELECT projections.keys.id,`+
						` projections.keys.creation_date,`+
						` projections.keys.change_date,`+
						` projections.keys.sequence,`+
						` projections.keys.resource_owner,`+
						` projections.keys.algorithm,`+
						` projections.keys.use,`+
						` projections.keys_public.expiry,`+
						` projections.keys_public.key,`+
						` COUNT(*) OVER ()`+
						` FROM projections.keys`+
						` LEFT JOIN projections.keys_public ON projections.keys.id = projections.keys_public.id`),
					[]string{
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
					},
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
					regexp.QuoteMeta(`SELECT projections.keys.id,`+
						` projections.keys.creation_date,`+
						` projections.keys.change_date,`+
						` projections.keys.sequence,`+
						` projections.keys.resource_owner,`+
						` projections.keys.algorithm,`+
						` projections.keys.use,`+
						` projections.keys_public.expiry,`+
						` projections.keys_public.key,`+
						` COUNT(*) OVER ()`+
						` FROM projections.keys`+
						` LEFT JOIN projections.keys_public ON projections.keys.id = projections.keys_public.id`),
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
			name:    "preparePrivateKeysQuery no result",
			prepare: preparePrivateKeysQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(`SELECT projections.keys.id,`+
						` projections.keys.creation_date,`+
						` projections.keys.change_date,`+
						` projections.keys.sequence,`+
						` projections.keys.resource_owner,`+
						` projections.keys.algorithm,`+
						` projections.keys.use,`+
						` projections.keys_private.expiry,`+
						` projections.keys_private.key,`+
						` COUNT(*) OVER ()`+
						` FROM projections.keys`+
						` LEFT JOIN projections.keys_private ON projections.keys.id = projections.keys_private.id`),
					nil,
					nil,
				),
				err: func(err error) (error, bool) {
					if !errs.IsNotFound(err) {
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
					regexp.QuoteMeta(`SELECT projections.keys.id,`+
						` projections.keys.creation_date,`+
						` projections.keys.change_date,`+
						` projections.keys.sequence,`+
						` projections.keys.resource_owner,`+
						` projections.keys.algorithm,`+
						` projections.keys.use,`+
						` projections.keys_private.expiry,`+
						` projections.keys_private.key,`+
						` COUNT(*) OVER ()`+
						` FROM projections.keys`+
						` LEFT JOIN projections.keys_private ON projections.keys.id = projections.keys_private.id`),
					[]string{
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
					},
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
					regexp.QuoteMeta(`SELECT projections.keys.id,`+
						` projections.keys.creation_date,`+
						` projections.keys.change_date,`+
						` projections.keys.sequence,`+
						` projections.keys.resource_owner,`+
						` projections.keys.algorithm,`+
						` projections.keys.use,`+
						` projections.keys_private.expiry,`+
						` projections.keys_private.key,`+
						` COUNT(*) OVER ()`+
						` FROM projections.keys`+
						` LEFT JOIN projections.keys_private ON projections.keys.id = projections.keys_private.id`),
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
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
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
