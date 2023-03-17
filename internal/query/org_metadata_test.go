package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	errs "github.com/zitadel/zitadel/internal/errors"
)

var (
	orgMetadataQuery = `SELECT projections.org_metadata2.creation_date,` +
		` projections.org_metadata2.change_date,` +
		` projections.org_metadata2.resource_owner,` +
		` projections.org_metadata2.sequence,` +
		` projections.org_metadata2.key,` +
		` projections.org_metadata2.value` +
		` FROM projections.org_metadata2` +
		` AS OF SYSTEM TIME '-1 ms'`
	orgMetadataCols = []string{
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"key",
		"value",
	}
	orgMetadataListQuery = `SELECT projections.org_metadata2.creation_date,` +
		` projections.org_metadata2.change_date,` +
		` projections.org_metadata2.resource_owner,` +
		` projections.org_metadata2.sequence,` +
		` projections.org_metadata2.key,` +
		` projections.org_metadata2.value,` +
		` COUNT(*) OVER ()` +
		` FROM projections.org_metadata2` +
		` AS OF SYSTEM TIME '-1 ms'`
	orgMetadataListCols = []string{
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"key",
		"value",
		"count",
	}
)

func Test_OrgMetadataPrepares(t *testing.T) {
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
			name:    "prepareOrgMetadataQuery no result",
			prepare: prepareOrgMetadataQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(orgMetadataQuery),
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
			object: (*OrgMetadata)(nil),
		},
		{
			name:    "prepareOrgMetadataQuery found",
			prepare: prepareOrgMetadataQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(orgMetadataQuery),
					orgMetadataCols,
					[]driver.Value{
						testNow,
						testNow,
						"resource_owner",
						uint64(20211108),
						"key",
						[]byte("value"),
					},
				),
			},
			object: &OrgMetadata{
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "resource_owner",
				Sequence:      20211108,
				Key:           "key",
				Value:         []byte("value"),
			},
		},
		{
			name:    "prepareOrgMetadataQuery sql err",
			prepare: prepareOrgMetadataQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(orgMetadataQuery),
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
			name:    "prepareOrgMetadataListQuery no result",
			prepare: prepareOrgMetadataListQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(orgMetadataListQuery),
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
			object: &OrgMetadataList{Metadata: []*OrgMetadata{}},
		},
		{
			name:    "prepareOrgMetadataListQuery one result",
			prepare: prepareOrgMetadataListQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(orgMetadataListQuery),
					orgMetadataListCols,
					[][]driver.Value{
						{
							testNow,
							testNow,
							"resource_owner",
							uint64(20211108),
							"key",
							[]byte("value"),
						},
					},
				),
			},
			object: &OrgMetadataList{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Metadata: []*OrgMetadata{
					{
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "resource_owner",
						Sequence:      20211108,
						Key:           "key",
						Value:         []byte("value"),
					},
				},
			},
		},
		{
			name:    "prepareOrgMetadataListQuery multiple results",
			prepare: prepareOrgMetadataListQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(orgMetadataListQuery),
					orgMetadataListCols,
					[][]driver.Value{
						{
							testNow,
							testNow,
							"resource_owner",
							uint64(20211108),
							"key",
							[]byte("value"),
						},
						{
							testNow,
							testNow,
							"resource_owner",
							uint64(20211108),
							"key2",
							[]byte("value2"),
						},
					},
				),
			},
			object: &OrgMetadataList{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				Metadata: []*OrgMetadata{
					{
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "resource_owner",
						Sequence:      20211108,
						Key:           "key",
						Value:         []byte("value"),
					},
					{
						CreationDate:  testNow,
						ChangeDate:    testNow,
						ResourceOwner: "resource_owner",
						Sequence:      20211108,
						Key:           "key2",
						Value:         []byte("value2"),
					},
				},
			},
		},
		{
			name:    "prepareOrgMetadataListQuery sql err",
			prepare: prepareOrgMetadataListQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(orgMetadataListQuery),
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
