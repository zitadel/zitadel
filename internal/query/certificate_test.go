package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	errs "github.com/zitadel/zitadel/internal/errors"
)

var (
	prepareCertificateStmt = `SELECT projections.keys4.id,` +
		` projections.keys4.creation_date,` +
		` projections.keys4.change_date,` +
		` projections.keys4.sequence,` +
		` projections.keys4.resource_owner,` +
		` projections.keys4.algorithm,` +
		` projections.keys4.use,` +
		` projections.keys4_certificate.expiry,` +
		` projections.keys4_certificate.certificate,` +
		` projections.keys4_private.key,` +
		` COUNT(*) OVER ()` +
		` FROM projections.keys4` +
		` LEFT JOIN projections.keys4_certificate ON projections.keys4.id = projections.keys4_certificate.id AND projections.keys4.instance_id = projections.keys4_certificate.instance_id` +
		` LEFT JOIN projections.keys4_private ON projections.keys4.id = projections.keys4_private.id AND projections.keys4.instance_id = projections.keys4_private.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`
	prepareCertificateCols = []string{
		"id",
		"creation_date",
		"change_date",
		"sequence",
		"resource_owner",
		"algorithm",
		"use",
		"expiry",
		"certificate",
		"key",
		"count",
	}
)

func Test_CertificatePrepares(t *testing.T) {
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
			name:    "prepareCertificateQuery no result",
			prepare: prepareCertificateQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareCertificateStmt),
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
			object: &Certificates{Certificates: []Certificate{}},
		},
		{
			name:    "prepareCertificateQuery found",
			prepare: prepareCertificateQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareCertificateStmt),
					prepareCertificateCols,
					[][]driver.Value{
						{
							"key-id",
							testNow,
							testNow,
							uint64(20211109),
							"ro",
							"",
							1,
							testNow,
							[]byte(`privateKey`),
							[]byte(`{"Algorithm": "enc", "Crypted": "cHJpdmF0ZUtleQ==", "CryptoType": 0, "KeyID": "id"}`),
						},
					},
				),
			},
			object: &Certificates{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Certificates: []Certificate{
					&rsaCertificate{
						key: key{
							id:            "key-id",
							creationDate:  testNow,
							changeDate:    testNow,
							sequence:      20211109,
							resourceOwner: "ro",
							algorithm:     "",
							use:           domain.KeyUsageSAMLMetadataSigning,
						},
						expiry:      testNow,
						certificate: []byte("privateKey"),
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
			name:    "prepareCertificateQuery sql err",
			prepare: prepareCertificateQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareCertificateStmt),
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
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}
