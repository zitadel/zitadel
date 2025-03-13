package query

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"testing"

	sq "github.com/Masterminds/squirrel"

	"github.com/zitadel/zitadel/internal/domain"
)

var (
	loginPolicyIDPLinksQuery = regexp.QuoteMeta(`SELECT projections.idp_login_policy_links5.idp_id,` +
		` projections.idp_templates6.name,` +
		` projections.idp_templates6.type,` +
		` projections.idp_templates6.owner_type,` +
		` projections.idp_templates6.is_creation_allowed,` +
		` projections.idp_templates6.is_linking_allowed,` +
		` projections.idp_templates6.is_auto_creation,` +
		` projections.idp_templates6.is_auto_update,` +
		` projections.idp_templates6.auto_linking,` +
		` COUNT(*) OVER ()` +
		` FROM projections.idp_login_policy_links5` +
		` LEFT JOIN projections.idp_templates6 ON projections.idp_login_policy_links5.idp_id = projections.idp_templates6.id AND projections.idp_login_policy_links5.instance_id = projections.idp_templates6.instance_id` +
		` RIGHT JOIN (SELECT login_policy_owner.aggregate_id, login_policy_owner.instance_id, login_policy_owner.owner_removed FROM projections.login_policies5 AS login_policy_owner` +
		` WHERE (login_policy_owner.instance_id = $1 AND (login_policy_owner.aggregate_id = $2 OR login_policy_owner.aggregate_id = $3)) ORDER BY login_policy_owner.is_default LIMIT 1) AS login_policy_owner` +
		` ON login_policy_owner.aggregate_id = projections.idp_login_policy_links5.resource_owner AND login_policy_owner.instance_id = projections.idp_login_policy_links5.instance_id`)
	loginPolicyIDPLinksCols = []string{
		"idp_id",
		"name",
		"type",
		"owner_type",
		"is_creation_allowed",
		"is_linking_allowed",
		"is_auto_creation",
		"is_auto_update",
		"auto_linking",
		"count",
	}
)

func Test_IDPLoginPolicyLinkPrepares(t *testing.T) {
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
			name: "prepareIDPsQuery found",
			prepare: func(ctx context.Context) (sq.SelectBuilder, func(*sql.Rows) (*IDPLoginPolicyLinks, error)) {
				return prepareIDPLoginPolicyLinksQuery(ctx, "resourceOwner")
			},
			want: want{
				sqlExpectations: mockQueries(
					loginPolicyIDPLinksQuery,
					loginPolicyIDPLinksCols,
					[][]driver.Value{
						{
							"idp-id",
							"idp-name",
							domain.IDPTypeJWT,
							domain.IdentityProviderTypeSystem,
							true,
							true,
							true,
							true,
							domain.AutoLinkingOptionUsername,
						},
					},
				),
			},
			object: &IDPLoginPolicyLinks{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Links: []*IDPLoginPolicyLink{
					{
						IDPID:             "idp-id",
						IDPName:           "idp-name",
						IDPType:           domain.IDPTypeJWT,
						OwnerType:         domain.IdentityProviderTypeSystem,
						IsCreationAllowed: true,
						IsLinkingAllowed:  true,
						IsAutoCreation:    true,
						IsAutoUpdate:      true,
						AutoLinking:       domain.AutoLinkingOptionUsername,
					},
				},
			},
		},
		{
			name: "prepareIDPsQuery no idp",
			prepare: func(ctx context.Context) (sq.SelectBuilder, func(*sql.Rows) (*IDPLoginPolicyLinks, error)) {
				return prepareIDPLoginPolicyLinksQuery(ctx, "resourceOwner")
			},
			want: want{
				sqlExpectations: mockQueries(
					loginPolicyIDPLinksQuery,
					loginPolicyIDPLinksCols,
					[][]driver.Value{
						{
							"idp-id",
							nil,
							nil,
							nil,
							false,
							false,
							false,
							false,
							0,
						},
					},
				),
			},
			object: &IDPLoginPolicyLinks{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Links: []*IDPLoginPolicyLink{
					{
						IDPID:             "idp-id",
						IDPName:           "",
						IDPType:           domain.IDPTypeUnspecified,
						IsCreationAllowed: false,
						IsLinkingAllowed:  false,
						IsAutoCreation:    false,
						IsAutoUpdate:      false,
						AutoLinking:       domain.AutoLinkingOptionUnspecified,
					},
				},
			},
		},
		{
			name: "prepareIDPsQuery sql err",
			prepare: func(ctx context.Context) (sq.SelectBuilder, func(*sql.Rows) (*IDPLoginPolicyLinks, error)) {
				return prepareIDPLoginPolicyLinksQuery(ctx, "resourceOwner")
			},
			want: want{
				sqlExpectations: mockQueryErr(
					loginPolicyIDPLinksQuery,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*IDPs)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, reflect.ValueOf(context.Background()))
		})
	}
}
