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
	userMetadataQuery = `SELECT projections.user_metadata4.creation_date,` +
		` projections.user_metadata4.change_date,` +
		` projections.user_metadata4.resource_owner,` +
		` projections.user_metadata4.sequence,` +
		` projections.user_metadata4.key,` +
		` projections.user_metadata4.value` +
		` FROM projections.user_metadata4` +
		` AS OF SYSTEM TIME '-1 ms'`
	userMetadataCols = []string{
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"key",
		"value",
	}
	userMetadataListQuery = `SELECT projections.user_metadata4.creation_date,` +
		` projections.user_metadata4.change_date,` +
		` projections.user_metadata4.resource_owner,` +
		` projections.user_metadata4.sequence,` +
		` projections.user_metadata4.key,` +
		` projections.user_metadata4.value,` +
		` COUNT(*) OVER ()` +
		` FROM projections.user_metadata4`
	userMetadataListCols = []string{
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"key",
		"value",
		"count",
	}
)

func Test_UserMetadataPrepares(t *testing.T) {
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
			name:    "prepareUserMetadataQuery no result",
			prepare: prepareUserMetadataQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(userMetadataQuery),
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
			object: (*UserMetadata)(nil),
		},
		{
			name:    "prepareUserMetadataQuery found",
			prepare: prepareUserMetadataQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(userMetadataQuery),
					userMetadataCols,
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
			object: &UserMetadata{
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "resource_owner",
				Sequence:      20211108,
				Key:           "key",
				Value:         []byte("value"),
			},
		},
		{
			name:    "prepareUserMetadataQuery sql err",
			prepare: prepareUserMetadataQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(userMetadataQuery),
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
			name:    "prepareUserMetadataListQuery no result",
			prepare: prepareUserMetadataListQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(userMetadataListQuery),
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
			object: &UserMetadataList{Metadata: []*UserMetadata{}},
		},
		{
			name:    "prepareUserMetadataListQuery one result",
			prepare: prepareUserMetadataListQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(userMetadataListQuery),
					userMetadataListCols,
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
			object: &UserMetadataList{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Metadata: []*UserMetadata{
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
			name:    "prepareUserMetadataListQuery multiple results",
			prepare: prepareUserMetadataListQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(userMetadataListQuery),
					userMetadataListCols,
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
			object: &UserMetadataList{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				Metadata: []*UserMetadata{
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
			name:    "prepareUserMetadataListQuery sql err",
			prepare: prepareUserMetadataListQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(userMetadataListQuery),
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
