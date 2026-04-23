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
	groupMetadataQuery = `SELECT projections.group_metadata.creation_date,` +
		` projections.group_metadata.change_date,` +
		` projections.group_metadata.resource_owner,` +
		` projections.group_metadata.sequence,` +
		` projections.group_metadata.key,` +
		` projections.group_metadata.value` +
		` FROM projections.group_metadata`
	groupMetadataCols = []string{
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"key",
		"value",
	}
	groupMetadataListQuery = `SELECT projections.group_metadata.creation_date,` +
		` projections.group_metadata.change_date,` +
		` projections.group_metadata.group_id,` +
		` projections.group_metadata.resource_owner,` +
		` projections.group_metadata.sequence,` +
		` projections.group_metadata.key,` +
		` projections.group_metadata.value,` +
		` COUNT(*) OVER ()` +
		` FROM projections.group_metadata`
	groupMetadataListCols = []string{
		"creation_date",
		"change_date",
		"group_id",
		"resource_owner",
		"sequence",
		"key",
		"value",
		"count",
	}
)

func Test_GroupMetadataPrepares(t *testing.T) {
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
			name:    "prepareGroupMetadataQuery no result",
			prepare: prepareGroupMetadataQuery,
			want: want{
				sqlExpectations: mockQueryScanErr(
					regexp.QuoteMeta(groupMetadataQuery),
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
			object: (*GroupMetadata)(nil),
		},
		{
			name:    "prepareGroupMetadataQuery found",
			prepare: prepareGroupMetadataQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(groupMetadataQuery),
					groupMetadataCols,
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
			object: &GroupMetadata{
				CreationDate:  testNow,
				ChangeDate:    testNow,
				ResourceOwner: "resource_owner",
				Sequence:      20211108,
				Key:           "key",
				Value:         []byte("value"),
			},
		},
		{
			name:    "prepareGroupMetadataQuery sql err",
			prepare: prepareGroupMetadataQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(groupMetadataQuery),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*GroupMetadata)(nil),
		},
		{
			name:    "prepareGroupMetadataListQuery no result",
			prepare: prepareGroupMetadataListQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(groupMetadataListQuery),
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
			object: &GroupMetadataList{Metadata: []*GroupMetadata{}},
		},
		{
			name:    "prepareGroupMetadataListQuery one result",
			prepare: prepareGroupMetadataListQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(groupMetadataListQuery),
					groupMetadataListCols,
					[][]driver.Value{
						{
							testNow,
							testNow,
							"group-1",
							"resource_owner",
							uint64(20211108),
							"key",
							[]byte("value"),
						},
					},
				),
			},
			object: &GroupMetadataList{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Metadata: []*GroupMetadata{
					{
						CreationDate:  testNow,
						ChangeDate:    testNow,
						GroupID:       "group-1",
						ResourceOwner: "resource_owner",
						Sequence:      20211108,
						Key:           "key",
						Value:         []byte("value"),
					},
				},
			},
		},
		{
			name:    "prepareGroupMetadataListQuery multiple results",
			prepare: prepareGroupMetadataListQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(groupMetadataListQuery),
					groupMetadataListCols,
					[][]driver.Value{
						{
							testNow,
							testNow,
							"group-1",
							"resource_owner",
							uint64(20211108),
							"key",
							[]byte("value"),
						},
						{
							testNow,
							testNow,
							"group-2",
							"resource_owner",
							uint64(20211108),
							"key2",
							[]byte("value2"),
						},
					},
				),
			},
			object: &GroupMetadataList{
				SearchResponse: SearchResponse{
					Count: 2,
				},
				Metadata: []*GroupMetadata{
					{
						CreationDate:  testNow,
						ChangeDate:    testNow,
						GroupID:       "group-1",
						ResourceOwner: "resource_owner",
						Sequence:      20211108,
						Key:           "key",
						Value:         []byte("value"),
					},
					{
						CreationDate:  testNow,
						ChangeDate:    testNow,
						GroupID:       "group-2",
						ResourceOwner: "resource_owner",
						Sequence:      20211108,
						Key:           "key2",
						Value:         []byte("value2"),
					},
				},
			},
		},
		{
			name:    "prepareGroupMetadataListQuery sql err",
			prepare: prepareGroupMetadataListQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(groupMetadataListQuery),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*GroupMetadataList)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
