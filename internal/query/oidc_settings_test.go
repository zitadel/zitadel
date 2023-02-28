package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"
	"time"

	errs "github.com/zitadel/zitadel/internal/errors"
)

var (
	prepareOIDCSettingsStmt = `SELECT projections.oidc_settings2.aggregate_id,` +
		` projections.oidc_settings2.creation_date,` +
		` projections.oidc_settings2.change_date,` +
		` projections.oidc_settings2.resource_owner,` +
		` projections.oidc_settings2.sequence,` +
		` projections.oidc_settings2.access_token_lifetime,` +
		` projections.oidc_settings2.id_token_lifetime,` +
		` projections.oidc_settings2.refresh_token_idle_expiration,` +
		` projections.oidc_settings2.refresh_token_expiration` +
		` FROM projections.oidc_settings2` +
		` AS OF SYSTEM TIME '-1 ms'`
	prepareOIDCSettingsCols = []string{
		"aggregate_id",
		"creation_date",
		"change_date",
		"resource_owner",
		"sequence",
		"access_token_lifetime",
		"id_token_lifetime",
		"refresh_token_idle_expiration",
		"refresh_token_expiration",
	}
)

func Test_OIDCConfigsPrepares(t *testing.T) {
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
			name:    "prepareOIDCSettingsQuery no result",
			prepare: prepareOIDCSettingsQuery,
			want: want{
				sqlExpectations: mockQueries(
					prepareOIDCSettingsStmt,
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
			object: (*OIDCSettings)(nil),
		},
		{
			name:    "prepareOIDCSettingsQuery found",
			prepare: prepareOIDCSettingsQuery,
			want: want{
				sqlExpectations: mockQuery(
					regexp.QuoteMeta(prepareOIDCSettingsStmt),
					prepareOIDCSettingsCols,
					[]driver.Value{
						"agg-id",
						testNow,
						testNow,
						"ro",
						uint64(20211108),
						time.Minute * 1,
						time.Minute * 2,
						time.Minute * 3,
						time.Minute * 4,
					},
				),
			},
			object: &OIDCSettings{
				AggregateID:                "agg-id",
				CreationDate:               testNow,
				ChangeDate:                 testNow,
				ResourceOwner:              "ro",
				Sequence:                   20211108,
				AccessTokenLifetime:        time.Minute * 1,
				IdTokenLifetime:            time.Minute * 2,
				RefreshTokenIdleExpiration: time.Minute * 3,
				RefreshTokenExpiration:     time.Minute * 4,
			},
		},
		{
			name:    "prepareOIDCSettingsQuery sql err",
			prepare: prepareOIDCSettingsQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					regexp.QuoteMeta(prepareOIDCSettingsStmt),
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
