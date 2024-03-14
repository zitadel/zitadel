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
	prepareTargetsStmt = `SELECT projections.targets.id,` +
		` projections.targets.change_date,` +
		` projections.targets.resource_owner,` +
		` projections.targets.sequence,` +
		` projections.targets.name,` +
		` projections.targets.target_type,` +
		` projections.targets.timeout,` +
		` projections.targets.url,` +
		` projections.targets.async,` +
		` projections.targets.interrupt_on_error,` +
		` COUNT(*) OVER ()` +
		` FROM projections.targets`
	prepareTargetsCols = []string{
		"id",
		"change_date",
		"resource_owner",
		"sequence",
		"name",
		"target_type",
		"timeout",
		"url",
		"async",
		"interrupt_on_error",
		"count",
	}

	prepareTargetStmt = `SELECT projections.targets.id,` +
		` projections.targets.change_date,` +
		` projections.targets.resource_owner,` +
		` projections.targets.sequence,` +
		` projections.targets.name,` +
		` projections.targets.target_type,` +
		` projections.targets.timeout,` +
		` projections.targets.url,` +
		` projections.targets.async,` +
		` projections.targets.interrupt_on_error` +
		` FROM projections.targets`
	prepareTargetCols = []string{
		"id",
		"change_date",
		"resource_owner",
		"sequence",
		"name",
		"target_type",
		"timeout",
		"url",
		"async",
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
						URL:              "https://example.com",
						Async:            true,
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
							false,
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
							true,
						},
					},
				),
			},
			object: &Targets{
				SearchResponse: SearchResponse{
					Count: 2,
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
						URL:              "https://example.com",
						Async:            true,
						InterruptOnError: false,
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
						URL:              "https://example.com",
						Async:            false,
						InterruptOnError: true,
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
						false,
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
				URL:              "https://example.com",
				Async:            true,
				InterruptOnError: false,
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
