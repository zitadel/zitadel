package query

import (
	"context"
	"database/sql/driver"
	"fmt"
	"strings"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/database"
	db_mock "github.com/zitadel/zitadel/internal/database/mock"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/feature"
	"github.com/zitadel/zitadel/internal/zerrors"
)

func TestQueries_UserGrants_permissionCheck(t *testing.T) {
	// authz.WithFeatures changes the instance of the context it is given, so each case gets its own context
	newCtx := func(features feature.Features) context.Context {
		ctx := authz.WithInstanceID(t.Context(), "instance-id")
		ctx = authz.WithFeatures(ctx, features)
		return authz.SetCtxData(ctx, authz.CtxData{UserID: "caller"})
	}

	const (
		listStmt        = "FROM projections.user_grants5"
		latestStateStmt = "FROM projections.current_states"
	)

	tests := []struct {
		name              string
		ctx               context.Context
		rows              [][]driver.Value
		wantStmt          []string
		wantNotInStmt     []string
		wantIDs           []string
		wantCount         uint64
		wantPermissionChk int
	}{
		{
			name: "without permission check v2, the grants of the page are filtered after they are read",
			ctx:  newCtx(feature.Features{}),
			rows: [][]driver.Value{
				userGrantsPermissionRow("project-grant", "", 2),
				userGrantsPermissionRow("denied-grant", "grant-id", 2),
			},
			wantNotInStmt:     []string{"eventstore.permitted_projects"},
			wantIDs:           []string{"project-grant"},
			wantCount:         2,
			wantPermissionChk: 2,
		},
		{
			name: "with permission check v2, the database filters and counts the permitted grants",
			ctx:  newCtx(feature.Features{PermissionCheckV2: true}),
			rows: [][]driver.Value{
				userGrantsPermissionRow("project-grant", "", 1),
			},
			wantStmt: []string{
				"INNER JOIN eventstore.permitted_projects(",
				"(projections.user_grants5.grant_id = '' AND projections.user_grants5.project_id = ANY(permissions.project_ids))",
				"projections.user_grants5.grant_id = ANY(permissions.project_grant_ids)",
				"projections.user_grants5.user_id = $",
			},
			wantIDs:           []string{"project-grant"},
			wantCount:         1,
			wantPermissionChk: 0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var gotStmt string
			client, mock, err := sqlmock.New(
				sqlmock.QueryMatcherOption(sqlmock.QueryMatcherFunc(func(expected, actual string) error {
					if !strings.Contains(actual, expected) {
						return fmt.Errorf("%q does not contain %q", actual, expected)
					}
					if expected == listStmt {
						gotStmt = actual
					}
					return nil
				})),
				sqlmock.ValueConverterOption(new(db_mock.TypeConverter)),
			)
			require.NoError(t, err)
			defer client.Close()

			mock.ExpectQuery(latestStateStmt).WillReturnRows(
				sqlmock.NewRows([]string{"event_date", "position", "last_updated"}),
			)
			rows := sqlmock.NewRows(userGrantsCols)
			for _, row := range tt.rows {
				rows.AddRow(row...)
			}
			mock.ExpectQuery(listStmt).WillReturnRows(rows)

			var permissionChecks int
			permissionCheck := func(_ context.Context, _, _, aggregateID string) error {
				permissionChecks++
				if aggregateID == "grant-id" {
					return zerrors.ThrowPermissionDenied(nil, "TEST-Ug1Pc", "denied")
				}
				return nil
			}

			q := &Queries{client: &database.DB{DB: client}}
			got, err := q.UserGrants(tt.ctx, &UserGrantsQueries{SearchRequest: SearchRequest{Limit: 10}}, false, permissionCheck)
			require.NoError(t, err)
			require.NoError(t, mock.ExpectationsWereMet())

			for _, want := range tt.wantStmt {
				assert.Contains(t, gotStmt, want)
			}
			for _, notWant := range tt.wantNotInStmt {
				assert.NotContains(t, gotStmt, notWant)
			}
			gotIDs := make([]string, len(got.UserGrants))
			for i, grant := range got.UserGrants {
				gotIDs[i] = grant.ID
			}
			assert.Equal(t, tt.wantIDs, gotIDs)
			assert.Equal(t, tt.wantCount, got.Count)
			assert.Equal(t, tt.wantPermissionChk, permissionChecks)
		})
	}
}

func userGrantsPermissionRow(id, grantID string, count uint64) []driver.Value {
	return []driver.Value{
		id,
		testNow,
		testNow,
		20211111,
		grantID,
		database.TextArray[string]{"role-key"},
		nil,
		domain.UserGrantStateActive,
		"user-id",
		"username",
		domain.UserTypeHuman,
		"user-resource-owner",
		"user-resource-owner-name",
		"user-resource-owner-domain",
		"first-name",
		"last-name",
		"email",
		"display-name",
		"avatar-key",
		"login-name",
		"ro",
		"org-name",
		"primary-domain",
		"project-id",
		"project-name",
		"project-resource-owner",
		"granted-org-id",
		"granted-org-name",
		"granted-org-domain",
		count,
	}
}
