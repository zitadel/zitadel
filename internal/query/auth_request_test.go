package query

import (
	"database/sql"
	"database/sql/driver"
	_ "embed"
	"fmt"
	"regexp"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/errors"
	"github.com/zitadel/zitadel/internal/query/projection"
)

func TestQueries_AuthRequestByID(t *testing.T) {
	expQuery := regexp.QuoteMeta(fmt.Sprintf(
		authRequestByIDQuery,
		asOfSystemTime,
	))

	cols := []string{
		projection.AuthRequestColumnID,
		projection.AuthRequestColumnCreationDate,
		projection.AuthRequestColumnLoginClient,
		projection.AuthRequestColumnClientID,
		projection.AuthRequestColumnScope,
		projection.AuthRequestColumnRedirectURI,
		projection.AuthRequestColumnPrompt,
		projection.AuthRequestColumnUILocales,
		projection.AuthRequestColumnLoginHint,
		projection.AuthRequestColumnMaxAge,
		projection.AuthRequestColumnHintUserID,
	}
	type args struct {
		shouldTriggerBulk bool
		id                string
		checkLoginClient  bool
	}
	tests := []struct {
		name    string
		args    args
		expect  sqlExpectation
		want    *AuthRequest
		wantErr error
	}{
		{
			name: "success, all values",
			args: args{
				shouldTriggerBulk: false,
				id:                "123",
				checkLoginClient:  true,
			},
			expect: mockQuery(expQuery, cols, []driver.Value{
				"id",
				testNow,
				"loginClient",
				"clientID",
				database.StringArray{"a", "b", "c"},
				"example.com",
				database.EnumArray[domain.Prompt]{domain.PromptLogin, domain.PromptConsent},
				database.StringArray{"en", "fi"},
				"me@example.com",
				int64(time.Minute),
				"userID",
			}, "123", "instanceID"),
			want: &AuthRequest{
				ID:           "id",
				CreationDate: testNow,
				LoginClient:  "loginClient",
				ClientID:     "clientID",
				Scope:        []string{"a", "b", "c"},
				RedirectURI:  "example.com",
				Prompt:       []domain.Prompt{domain.PromptLogin, domain.PromptConsent},
				UiLocales:    []string{"en", "fi"},
				LoginHint:    gu.Ptr("me@example.com"),
				MaxAge:       gu.Ptr(time.Minute),
				HintUserID:   gu.Ptr("userID"),
			},
		},
		{
			name: "success, null values",
			args: args{
				shouldTriggerBulk: false,
				id:                "123",
				checkLoginClient:  true,
			},
			expect: mockQuery(expQuery, cols, []driver.Value{
				"id",
				testNow,
				"loginClient",
				"clientID",
				database.StringArray{"a", "b", "c"},
				"example.com",
				database.EnumArray[domain.Prompt]{domain.PromptLogin, domain.PromptConsent},
				database.StringArray{"en", "fi"},
				sql.NullString{},
				sql.NullInt64{},
				sql.NullString{},
			}, "123", "instanceID"),
			want: &AuthRequest{
				ID:           "id",
				CreationDate: testNow,
				LoginClient:  "loginClient",
				ClientID:     "clientID",
				Scope:        []string{"a", "b", "c"},
				RedirectURI:  "example.com",
				Prompt:       []domain.Prompt{domain.PromptLogin, domain.PromptConsent},
				UiLocales:    []string{"en", "fi"},
				LoginHint:    nil,
				MaxAge:       nil,
				HintUserID:   nil,
			},
		},
		{
			name: "no rows",
			args: args{
				shouldTriggerBulk: false,
				id:                "123",
			},
			expect:  mockQuery(expQuery, cols, nil, "123", "instanceID"),
			wantErr: errors.ThrowNotFound(sql.ErrNoRows, "QUERY-Thee9", "Errors.AuthRequest.NotExisting"),
		},
		{
			name: "query error",
			args: args{
				shouldTriggerBulk: false,
				id:                "123",
			},
			expect:  mockQueryErr(expQuery, sql.ErrConnDone, "123", "instanceID"),
			wantErr: errors.ThrowInternal(sql.ErrConnDone, "QUERY-Ou8ue", "Errors.Internal"),
		},
		{
			name: "wrong login client",
			args: args{
				shouldTriggerBulk: false,
				id:                "123",
				checkLoginClient:  true,
			},
			expect: mockQuery(expQuery, cols, []driver.Value{
				"id",
				testNow,
				"wrongLoginClient",
				"clientID",
				database.StringArray{"a", "b", "c"},
				"example.com",
				database.EnumArray[domain.Prompt]{domain.PromptLogin, domain.PromptConsent},
				database.StringArray{"en", "fi"},
				sql.NullString{},
				sql.NullInt64{},
				sql.NullString{},
			}, "123", "instanceID"),
			wantErr: errors.ThrowPermissionDeniedf(nil, "OIDCv2-aL0ag", "Errors.AuthRequest.WrongLoginClient"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execMock(t, tt.expect, func(db *sql.DB) {
				q := &Queries{
					client: &database.DB{
						DB:       db,
						Database: &prepareDB{},
					},
				}
				ctx := authz.NewMockContext("instanceID", "orgID", "loginClient")

				got, err := q.AuthRequestByID(ctx, tt.args.shouldTriggerBulk, tt.args.id, tt.args.checkLoginClient)
				require.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, tt.want, got)
			})
		})
	}
}
