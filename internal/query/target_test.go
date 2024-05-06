package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

var (
	prepareTargetsStmt = `SELECT projections.targets1.id,` +
		` projections.targets1.change_date,` +
		` projections.targets1.resource_owner,` +
		` projections.targets1.sequence,` +
		` projections.targets1.name,` +
		` projections.targets1.target_type,` +
		` projections.targets1.timeout,` +
		` projections.targets1.endpoint,` +
		` projections.targets1.interrupt_on_error,` +
		` COUNT(*) OVER ()` +
		` FROM projections.targets1`
	prepareTargetsCols = []string{
		"id",
		"change_date",
		"resource_owner",
		"sequence",
		"name",
		"target_type",
		"timeout",
		"endpoint",
		"interrupt_on_error",
		"count",
	}

	prepareTargetStmt = `SELECT projections.targets1.id,` +
		` projections.targets1.change_date,` +
		` projections.targets1.resource_owner,` +
		` projections.targets1.sequence,` +
		` projections.targets1.name,` +
		` projections.targets1.target_type,` +
		` projections.targets1.timeout,` +
		` projections.targets1.endpoint,` +
		` projections.targets1.interrupt_on_error` +
		` FROM projections.targets1`
	prepareTargetCols = []string{
		"id",
		"change_date",
		"resource_owner",
		"sequence",
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
							"ro",
							uint64(20211109),
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
						ID: "id",
						ObjectDetails: domain.ObjectDetails{
							EventDate:     testNow,
							ResourceOwner: "ro",
							Sequence:      20211109,
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
							"ro",
							uint64(20211109),
							"target-name1",
							domain.TargetTypeWebhook,
							1 * time.Second,
							"https://example.com",
							true,
						},
						{
							"id-2",
							testNow,
							"ro",
							uint64(20211110),
							"target-name2",
							domain.TargetTypeWebhook,
							1 * time.Second,
							"https://example.com",
							false,
						},
						{
							"id-3",
							testNow,
							"ro",
							uint64(20211110),
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
						ID: "id-1",
						ObjectDetails: domain.ObjectDetails{
							EventDate:     testNow,
							ResourceOwner: "ro",
							Sequence:      20211109,
						},
						Name:             "target-name1",
						TargetType:       domain.TargetTypeWebhook,
						Timeout:          1 * time.Second,
						Endpoint:         "https://example.com",
						InterruptOnError: true,
					},
					{
						ID: "id-2",
						ObjectDetails: domain.ObjectDetails{
							EventDate:     testNow,
							ResourceOwner: "ro",
							Sequence:      20211110,
						},
						Name:             "target-name2",
						TargetType:       domain.TargetTypeWebhook,
						Timeout:          1 * time.Second,
						Endpoint:         "https://example.com",
						InterruptOnError: false,
					},
					{
						ID: "id-3",
						ObjectDetails: domain.ObjectDetails{
							EventDate:     testNow,
							ResourceOwner: "ro",
							Sequence:      20211110,
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
						"ro",
						uint64(20211109),
						"target-name",
						domain.TargetTypeWebhook,
						1 * time.Second,
						"https://example.com",
						true,
					},
				),
			},
			object: &Target{
				ID: "id",
				ObjectDetails: domain.ObjectDetails{
					EventDate:     testNow,
					ResourceOwner: "ro",
					Sequence:      20211109,
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
