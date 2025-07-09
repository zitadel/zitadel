package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	projectMetadataQuery = `SELECT projections.project_metadata.creation_date,` +
		` projections.project_metadata.change_date,` +
		` projections.project_metadata.resource_owner,` +
		` projections.project_metadata.sequence,` +
		` projections.project_metadata.key,` +
		` projections.project_metadata.value` +
		` FROM projections.project_metadata`
	projectMetadataCols = []string{
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"key",
		"value",
	}
	projectMetadataListQuery = `SELECT projections.project_metadata.creation_date,` +
		` projections.project_metadata.change_date,` +
		` projections.project_metadata.resource_owner,` +
		` projections.project_metadata.sequence,` +
		` projections.project_metadata.key,` +
		` projections.project_metadata.value,` +
		` COUNT(*) OVER ()` +
		` FROM projections.project_metadata`
	projectMetadataListCols = []string{
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"key",
		"value",
		"count",
	}
)

func Test_ProjectMetadataPrepares(t *testing.T) {
	type want struct {
		sqlExpectations sqlExpectation
		err             checkErr
	}
	tests := []struct {
		name    string
		prepare any
		want    want
		object  any
	}{
		{
			name:    "prepareProjectMetadataQuery no result",
			prepare: prepareProjectMetadataQuery,
			want: want{
				sqlExpectations: mockQueryScanErr(
					regexp.QuoteMeta(projectMetadataQuery),
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
			object: (*ProjectMetadata)(nil),
		},
		{
			name:    "prepareProjectMetadataQuery found",
			prepare: prepareProjectMetadataQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(projectMetadataQuery),
					projectMetadataCols,
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
			object: &ProjectMetadata{
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "resource_owner",
				Sequence:      20211108,
				Key:           "key",
				Value:         []byte("value"),
			},
		},
		{
			name:    "prepareProjectMetadataQuery sql err",
			prepare: prepareProjectMetadataQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(projectMetadataQuery),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*ProjectMetadata)(nil),
		},
		{
			name:    "prepareProjectMetadataListQuery no result",
			prepare: prepareProjectMetadataListQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(projectMetadataListQuery),
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
			object: &ProjectMetadataList{Metadata: []*ProjectMetadata{}},
		},
		{
			name:    "prepareProjectMetadataListQuery one result",
			prepare: prepareProjectMetadataListQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(projectMetadataListQuery),
					projectMetadataListCols,
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
			object: &ProjectMetadataList{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Metadata: []*ProjectMetadata{
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
			name:    "prepareProjectMetadataListQuery multiple results",
			prepare: prepareProjectMetadataListQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(projectMetadataListQuery),
					projectMetadataListCols,
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
			object: &ProjectMetadataList{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				Metadata: []*ProjectMetadata{
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
			name:    "prepareProjectMetadataListQuery sql err",
			prepare: prepareProjectMetadataListQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(projectMetadataListQuery),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*ProjectMetadataList)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
