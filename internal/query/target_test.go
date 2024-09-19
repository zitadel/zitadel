package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/zitadel/zitadel/v2/internal/domain"
	"github.com/zitadel/zitadel/v2/internal/zerrors"
)

var (
	prepareTargetsStmt = `SELECT projections.targets1.id,` +
		` projections.targets1.creation_date,` +
		` projections.targets1.change_date,` +
		` projections.targets1.resource_owner,` +
		` projections.targets1.name,` +
		` projections.targets1.target_type,` +
		` projections.targets1.timeout,` +
		` projections.targets1.endpoint,` +
		` projections.targets1.interrupt_on_error,` +
		` COUNT(*) OVER ()` +
		` FROM projections.targets1`
	prepareTargetsCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"name",
		"target_type",
		"timeout",
		"endpoint",
		"interrupt_on_error",
		"count",
	}

	prepareTargetStmt = `SELECT projections.targets1.id,` +
		` projections.targets1.creation_date,` +
		` projections.targets1.change_date,` +
		` projections.targets1.resource_owner,` +
		` projections.targets1.name,` +
		` projections.targets1.target_type,` +
		` projections.targets1.timeout,` +
		` projections.targets1.endpoint,` +
		` projections.targets1.interrupt_on_error` +
		` FROM projections.targets1`
	prepareTargetCols = []string{
		"id",
		"creation_date",
		"change_date",
		"resource_owner",
		"name",
		"target_type",
		"timeout",
		"endpoint",
		"interrupt_on_error",
	}
)

func Test_TargetPrepares(t *testing.T) {
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
			name:    "prepareTargetsQuery no result",
			prepare: prepareTargetsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareTargetsStmt),
					nil,
					nil,
				),
			},
			object: &Targets{Targets: []*Target{}},
		},
		{
			name:    "prepareTargetsQuery one result",
			prepare: prepareTargetsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareTargetsStmt),
					prepareTargetsCols,
					[][]driver.Value{
						{
							"id",
							testNow,
							testNow,
							"ro",
							"target-name",
							domain.TargetTypeWebhook,
							1 * time.Second,
							"https://example.com",
							true,
						},
					},
				),
			},
			object: &Targets{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Targets: []*Target{
					{
						ObjectDetails: domain.ObjectDetails{
							ID:            "id",
							EventDate:     testNow,
							CreationDate:  testNow,
							ResourceOwner: "ro",
						},
						Name:             "target-name",
						TargetType:       domain.TargetTypeWebhook,
						Timeout:          1 * time.Second,
						Endpoint:         "https://example.com",
						InterruptOnError: true,
					},
				},
			},
		},
		{
			name:    "prepareTargetsQuery multiple result",
			prepare: prepareTargetsQuery,
			want: want{
				sqlExpectations: mockQueries(
					regexp.QuoteMeta(prepareTargetsStmt),
					prepareTargetsCols,
					[][]driver.Value{
						{
							"id-1",
							testNow,
							testNow,
							"ro",
							"target-name1",
							domain.TargetTypeWebhook,
							1 * time.Second,
							"https://example.com",
							true,
						},
						{
							"id-2",
							testNow,
							testNow,
							"ro",
							"target-name2",
							domain.TargetTypeWebhook,
							1 * time.Second,
							"https://example.com",
							false,
						},
						{
							"id-3",
							testNow,
							testNow,
							"ro",
							"target-name3",
							domain.TargetTypeAsync,
							1 * time.Second,
							"https://example.com",
							false,
						},
					},
				),
			},
			object: &Targets{
				SearchResponse: SearchResponse{
					Count: 3,
				},
				Targets: []*Target{
					{
						ObjectDetails: domain.ObjectDetails{
							ID:            "id-1",
							EventDate:     testNow,
							CreationDate:  testNow,
							ResourceOwner: "ro",
						},
						Name:             "target-name1",
						TargetType:       domain.TargetTypeWebhook,
						Timeout:          1 * time.Second,
						Endpoint:         "https://example.com",
						InterruptOnError: true,
					},
					{
						ObjectDetails: domain.ObjectDetails{
							ID:            "id-2",
							EventDate:     testNow,
							CreationDate:  testNow,
							ResourceOwner: "ro",
						},
						Name:             "target-name2",
						TargetType:       domain.TargetTypeWebhook,
						Timeout:          1 * time.Second,
						Endpoint:         "https://example.com",
						InterruptOnError: false,
					},
					{
						ObjectDetails: domain.ObjectDetails{
							ID:            "id-3",
							EventDate:     testNow,
							CreationDate:  testNow,
							ResourceOwner: "ro",
						},
						Name:             "target-name3",
						TargetType:       domain.TargetTypeAsync,
						Timeout:          1 * time.Second,
						Endpoint:         "https://example.com",
						InterruptOnError: false,
					},
				},
			},
		},
		{
			name:    "prepareTargetsQuery sql err",
			prepare: prepareTargetsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareTargetsStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Target)(nil),
		},
		{
			name:    "prepareTargetQuery no result",
			prepare: prepareTargetQuery,
			want: want{
				sqlExpectations: mockQueriesScanErr(
					regexp.QuoteMeta(prepareTargetStmt),
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
			object: (*Target)(nil),
		},
		{
			name:    "prepareTargetQuery found",
			prepare: prepareTargetQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(prepareTargetStmt),
					prepareTargetCols,
					[]driver.Value{
						"id",
						testNow,
						testNow,
						"ro",
						"target-name",
						domain.TargetTypeWebhook,
						1 * time.Second,
						"https://example.com",
						true,
					},
				),
			},
			object: &Target{
				ObjectDetails: domain.ObjectDetails{
					ID:            "id",
					EventDate:     testNow,
					CreationDate:  testNow,
					ResourceOwner: "ro",
				},
				Name:             "target-name",
				TargetType:       domain.TargetTypeWebhook,
				Timeout:          1 * time.Second,
				Endpoint:         "https://example.com",
				InterruptOnError: true,
			},
		},
		{
			name:    "prepareTargetQuery sql err",
			prepare: prepareTargetQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareTargetStmt),
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*Target)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}
