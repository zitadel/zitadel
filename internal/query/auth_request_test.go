package query

import (
	"context"
	"database/sql"
	"database/sql/driver"
	_ "embed"
	"regexp"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query/projection"
	"github.com/zitadel/zitadel/internal/zerrors"
)

// chainExpectations applies multiple sql expectations in order. sqlmock is ordered, so this
// models AuthRequestByID issuing its own query and then the private-labeling ProjectByClientID
// lookup from resolvePrivateLabelingOrgID.
func chainExpectations(exps ...sqlExpectation) sqlExpectation {
	return func(m sqlmock.Sqlmock) sqlmock.Sqlmock {
		for _, exp := range exps {
			if exp != nil {
				m = exp(m)
			}
		}
		return m
	}
}

// mockProjectByClientID expects the branding ProjectByClientID lookup and returns the given
// project row; a nil row yields no rows (ProjectByClientID -> NotFound, so branding falls back
// to the instance default org). Args are intentionally not constrained, so the test is not
// coupled to the lookup's placeholder order.
func mockProjectByClientID(row []driver.Value) sqlExpectation {
	return func(m sqlmock.Sqlmock) sqlmock.Sqlmock {
		result := m.NewRows(projectCols)
		if len(row) > 0 {
			result.AddRow(row...)
		}
		m.ExpectQuery(expectedProjectByAppQuery).WillReturnRows(result)
		return m
	}
}

// mockProjectByClientIDErr expects the ProjectByClientID lookup to fail; branding must then fall
// back to the instance default organization rather than failing the auth request.
func mockProjectByClientIDErr(err error) sqlExpectation {
	return func(m sqlmock.Sqlmock) sqlmock.Sqlmock {
		m.ExpectQuery(expectedProjectByAppQuery).WillReturnError(err)
		return m
	}
}

func TestQueries_AuthRequestByID(t *testing.T) {
	expQuery := regexp.QuoteMeta(authRequestByIDQuery)

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

	// A project row (projectCols order) whose private-labeling setting enforces the project's
	// resource-owner organization for branding.
	projectEnforceRow := []driver.Value{
		"projectID",
		testNow,
		testNow,
		"projectOrg",
		domain.ProjectStateActive,
		uint64(1),
		"projectName",
		false,
		false,
		false,
		domain.PrivateLabelingSettingEnforceProjectResourceOwnerPolicy,
	}
	// A project row whose private-labeling setting is unspecified, so branding falls back to the
	// instance default organization.
	projectUnspecifiedRow := []driver.Value{
		"projectID",
		testNow,
		testNow,
		"projectOrg",
		domain.ProjectStateActive,
		uint64(1),
		"projectName",
		false,
		false,
		false,
		domain.PrivateLabelingSettingUnspecified,
	}

	type args struct {
		shouldTriggerBulk bool
		id                string
		checkLoginClient  bool
	}
	tests := []struct {
		name            string
		args            args
		expect          sqlExpectation
		permissionCheck domain.PermissionCheck
		want            *AuthRequest
		wantErr         error
	}{
		{
			name: "success, all values",
			args: args{
				shouldTriggerBulk: false,
				id:                "123",
				checkLoginClient:  true,
			},
			expect: chainExpectations(
				mockQuery(expQuery, cols, []driver.Value{
					"id",
					testNow,
					"loginClient",
					"clientID",
					database.TextArray[string]{"a", "b", "c"},
					"example.com",
					database.NumberArray[domain.Prompt]{domain.PromptLogin, domain.PromptConsent},
					database.TextArray[string]{"en", "fi"},
					"me@example.com",
					int64(time.Minute),
					"userID",
				}, "123", "instanceID"),
				// No org scope and a client ID, so branding is resolved via the project lookup;
				// with no project found it falls back to the instance default organization.
				mockProjectByClientID(nil),
			),
			want: &AuthRequest{
				ID:                   "id",
				CreationDate:         testNow,
				LoginClient:          "loginClient",
				ClientID:             "clientID",
				Scope:                []string{"a", "b", "c"},
				RedirectURI:          "example.com",
				Prompt:               []domain.Prompt{domain.PromptLogin, domain.PromptConsent},
				UiLocales:            []string{"en", "fi"},
				LoginHint:            gu.Ptr("me@example.com"),
				MaxAge:               gu.Ptr(time.Minute),
				HintUserID:           gu.Ptr("userID"),
				PrivateLabelingOrgID: "orgID",
			},
		},
		{
			name: "success, null values",
			args: args{
				shouldTriggerBulk: false,
				id:                "123",
				checkLoginClient:  true,
			},
			expect: chainExpectations(
				mockQuery(expQuery, cols, []driver.Value{
					"id",
					testNow,
					"loginClient",
					"clientID",
					database.TextArray[string]{"a", "b", "c"},
					"example.com",
					database.NumberArray[domain.Prompt]{domain.PromptLogin, domain.PromptConsent},
					database.TextArray[string]{"en", "fi"},
					nil,
					nil,
					nil,
				}, "123", "instanceID"),
				mockProjectByClientID(nil),
			),
			want: &AuthRequest{
				ID:                   "id",
				CreationDate:         testNow,
				LoginClient:          "loginClient",
				ClientID:             "clientID",
				Scope:                []string{"a", "b", "c"},
				RedirectURI:          "example.com",
				Prompt:               []domain.Prompt{domain.PromptLogin, domain.PromptConsent},
				UiLocales:            []string{"en", "fi"},
				LoginHint:            nil,
				MaxAge:               nil,
				HintUserID:           nil,
				PrivateLabelingOrgID: "orgID",
			},
		},
		{
			name: "no rows",
			args: args{
				shouldTriggerBulk: false,
				id:                "123",
			},
			expect:  mockQueryScanErr(expQuery, cols, nil, "123", "instanceID"),
			wantErr: zerrors.ThrowNotFound(sql.ErrNoRows, "QUERY-Thee9", "Errors.AuthRequest.NotExisting"),
		},
		{
			name: "query error",
			args: args{
				shouldTriggerBulk: false,
				id:                "123",
			},
			expect:  mockQueryErr(expQuery, sql.ErrConnDone, "123", "instanceID"),
			wantErr: zerrors.ThrowInternal(sql.ErrConnDone, "QUERY-Ou8ue", "Errors.Internal"),
		},
		{
			name: "wrong login client / not permitted",
			args: args{
				shouldTriggerBulk: false,
				id:                "123",
				checkLoginClient:  true,
			},
			expect: chainExpectations(
				mockQuery(expQuery, cols, []driver.Value{
					"id",
					testNow,
					"wrongLoginClient",
					"clientID",
					database.TextArray[string]{"a", "b", "c"},
					"example.com",
					database.NumberArray[domain.Prompt]{domain.PromptLogin, domain.PromptConsent},
					database.TextArray[string]{"en", "fi"},
					nil,
					nil,
					nil,
				}, "123", "instanceID"),
				// Branding is resolved before the login-client permission check.
				mockProjectByClientID(nil),
			),
			permissionCheck: func(ctx context.Context, permission, orgID, resourceID string) (err error) {
				return zerrors.ThrowPermissionDenied(nil, "id", "not permitted")
			},
			wantErr: zerrors.ThrowPermissionDenied(nil, "id", "not permitted"),
		},
		{
			name: "other login client / permitted",
			args: args{
				shouldTriggerBulk: false,
				id:                "123",
				checkLoginClient:  true,
			},
			expect: chainExpectations(
				mockQuery(expQuery, cols, []driver.Value{
					"id",
					testNow,
					"otherLoginClient",
					"clientID",
					database.TextArray[string]{"a", "b", "c"},
					"example.com",
					database.NumberArray[domain.Prompt]{domain.PromptLogin, domain.PromptConsent},
					database.TextArray[string]{"en", "fi"},
					nil,
					nil,
					nil,
				}, "123", "instanceID"),
				mockProjectByClientID(nil),
			),
			permissionCheck: func(ctx context.Context, permission, orgID, resourceID string) (err error) {
				return nil
			},
			want: &AuthRequest{
				ID:                   "id",
				CreationDate:         testNow,
				LoginClient:          "otherLoginClient",
				ClientID:             "clientID",
				Scope:                []string{"a", "b", "c"},
				RedirectURI:          "example.com",
				Prompt:               []domain.Prompt{domain.PromptLogin, domain.PromptConsent},
				UiLocales:            []string{"en", "fi"},
				LoginHint:            nil,
				MaxAge:               nil,
				HintUserID:           nil,
				PrivateLabelingOrgID: "orgID",
			},
		},
		{
			name: "explicit org scope wins for branding",
			args: args{
				shouldTriggerBulk: false,
				id:                "123",
				checkLoginClient:  true,
			},
			// An explicit organization scope resolves branding directly; no project lookup runs.
			expect: mockQuery(expQuery, cols, []driver.Value{
				"id",
				testNow,
				"loginClient",
				"clientID",
				database.TextArray[string]{domain.OrgIDScope + "scopeOrg"},
				"example.com",
				database.NumberArray[domain.Prompt]{domain.PromptLogin},
				database.TextArray[string]{"en"},
				nil,
				nil,
				nil,
			}, "123", "instanceID"),
			want: &AuthRequest{
				ID:                   "id",
				CreationDate:         testNow,
				LoginClient:          "loginClient",
				ClientID:             "clientID",
				Scope:                []string{domain.OrgIDScope + "scopeOrg"},
				RedirectURI:          "example.com",
				Prompt:               []domain.Prompt{domain.PromptLogin},
				UiLocales:            []string{"en"},
				PrivateLabelingOrgID: "scopeOrg",
			},
		},
		{
			name: "project enforces resource owner policy for branding",
			args: args{
				shouldTriggerBulk: false,
				id:                "123",
				checkLoginClient:  true,
			},
			expect: chainExpectations(
				mockQuery(expQuery, cols, []driver.Value{
					"id",
					testNow,
					"loginClient",
					"clientID",
					database.TextArray[string]{"a", "b", "c"},
					"example.com",
					database.NumberArray[domain.Prompt]{domain.PromptLogin},
					database.TextArray[string]{"en"},
					nil,
					nil,
					nil,
				}, "123", "instanceID"),
				mockProjectByClientID(projectEnforceRow),
			),
			want: &AuthRequest{
				ID:                   "id",
				CreationDate:         testNow,
				LoginClient:          "loginClient",
				ClientID:             "clientID",
				Scope:                []string{"a", "b", "c"},
				RedirectURI:          "example.com",
				Prompt:               []domain.Prompt{domain.PromptLogin},
				UiLocales:            []string{"en"},
				PrivateLabelingOrgID: "projectOrg",
			},
		},
		{
			name: "project setting unspecified falls back to default org",
			args: args{
				shouldTriggerBulk: false,
				id:                "123",
				checkLoginClient:  true,
			},
			expect: chainExpectations(
				mockQuery(expQuery, cols, []driver.Value{
					"id",
					testNow,
					"loginClient",
					"clientID",
					database.TextArray[string]{"a", "b", "c"},
					"example.com",
					database.NumberArray[domain.Prompt]{domain.PromptLogin},
					database.TextArray[string]{"en"},
					nil,
					nil,
					nil,
				}, "123", "instanceID"),
				mockProjectByClientID(projectUnspecifiedRow),
			),
			want: &AuthRequest{
				ID:                   "id",
				CreationDate:         testNow,
				LoginClient:          "loginClient",
				ClientID:             "clientID",
				Scope:                []string{"a", "b", "c"},
				RedirectURI:          "example.com",
				Prompt:               []domain.Prompt{domain.PromptLogin},
				UiLocales:            []string{"en"},
				PrivateLabelingOrgID: "orgID",
			},
		},
		{
			name: "project lookup error falls back to default org",
			args: args{
				shouldTriggerBulk: false,
				id:                "123",
				checkLoginClient:  true,
			},
			expect: chainExpectations(
				mockQuery(expQuery, cols, []driver.Value{
					"id",
					testNow,
					"loginClient",
					"clientID",
					database.TextArray[string]{"a", "b", "c"},
					"example.com",
					database.NumberArray[domain.Prompt]{domain.PromptLogin},
					database.TextArray[string]{"en"},
					nil,
					nil,
					nil,
				}, "123", "instanceID"),
				mockProjectByClientIDErr(sql.ErrConnDone),
			),
			want: &AuthRequest{
				ID:                   "id",
				CreationDate:         testNow,
				LoginClient:          "loginClient",
				ClientID:             "clientID",
				Scope:                []string{"a", "b", "c"},
				RedirectURI:          "example.com",
				Prompt:               []domain.Prompt{domain.PromptLogin},
				UiLocales:            []string{"en"},
				PrivateLabelingOrgID: "orgID",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			execMock(t, tt.expect, func(db *sql.DB) {
				q := &Queries{
					client: &database.DB{
						DB: db,
					},
					checkPermission: tt.permissionCheck,
				}
				ctx := authz.NewMockContext("instanceID", "orgID", "loginClient")

				got, err := q.AuthRequestByID(ctx, tt.args.shouldTriggerBulk, tt.args.id, tt.args.checkLoginClient)
				require.ErrorIs(t, err, tt.wantErr)
				assert.Equal(t, tt.want, got)
			})
		})
	}
}
