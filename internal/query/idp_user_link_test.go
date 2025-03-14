package query

import (
	"context"
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/api/authz"
	"github.com/zitadel/zitadel/internal/domain"
)

func TestUser_idpLinksCheckPermission(t *testing.T) {
	type want struct {
		links []*IDPUserLink
	}
	type args struct {
		user  string
		links *IDPUserLinks
	}
	tests := []struct {
		name        string
		args        args
		want        want
		permissions []string
	}{
		{
			"permissions for all users",
			args{
				"none",
				&IDPUserLinks{
					Links: []*IDPUserLink{
						{UserID: "first"}, {UserID: "second"}, {UserID: "third"},
					},
				},
			},
			want{
				links: []*IDPUserLink{
					{UserID: "first"}, {UserID: "second"}, {UserID: "third"},
				},
			},
			[]string{"first", "second", "third"},
		},
		{
			"permissions for one user, first",
			args{
				"none",
				&IDPUserLinks{
					Links: []*IDPUserLink{
						{UserID: "first"}, {UserID: "second"}, {UserID: "third"},
					},
				},
			},
			want{
				links: []*IDPUserLink{
					{UserID: "first"},
				},
			},
			[]string{"first"},
		},
		{
			"permissions for one user, second",
			args{
				"none",
				&IDPUserLinks{
					Links: []*IDPUserLink{
						{UserID: "first"}, {UserID: "second"}, {UserID: "third"},
					},
				},
			},
			want{
				links: []*IDPUserLink{
					{UserID: "second"},
				},
			},
			[]string{"second"},
		},
		{
			"permissions for one user, third",
			args{
				"none",
				&IDPUserLinks{
					Links: []*IDPUserLink{
						{UserID: "first"}, {UserID: "second"}, {UserID: "third"},
					},
				},
			},
			want{
				links: []*IDPUserLink{
					{UserID: "third"},
				},
			},
			[]string{"third"},
		},
		{
			"permissions for two users, first",
			args{
				"none",
				&IDPUserLinks{
					Links: []*IDPUserLink{
						{UserID: "first"}, {UserID: "second"}, {UserID: "third"},
					},
				},
			},
			want{
				links: []*IDPUserLink{
					{UserID: "first"}, {UserID: "third"},
				},
			},
			[]string{"first", "third"},
		},
		{
			"permissions for two users, second",
			args{
				"none",
				&IDPUserLinks{
					Links: []*IDPUserLink{
						{UserID: "first"}, {UserID: "second"}, {UserID: "third"},
					},
				},
			},
			want{
				links: []*IDPUserLink{
					{UserID: "second"}, {UserID: "third"},
				},
			},
			[]string{"second", "third"},
		},
		{
			"no permissions",
			args{
				"none",
				&IDPUserLinks{
					Links: []*IDPUserLink{
						{UserID: "first"}, {UserID: "second"}, {UserID: "third"},
					},
				},
			},
			want{
				links: []*IDPUserLink{},
			},
			[]string{},
		},
		{
			"no permissions, self",
			args{
				"second",
				&IDPUserLinks{
					Links: []*IDPUserLink{
						{UserID: "first"}, {UserID: "second"}, {UserID: "third"},
					},
				},
			},
			want{
				links: []*IDPUserLink{{UserID: "second"}},
			},
			[]string{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			checkPermission := func(ctx context.Context, permission, orgID, resourceID string) (err error) {
				for _, perm := range tt.permissions {
					if resourceID == perm {
						return nil
					}
				}
				return errors.New("failed")
			}
			idpLinksCheckPermission(authz.SetCtxData(context.Background(), authz.CtxData{UserID: tt.args.user}), tt.args.links, checkPermission)
			require.Equal(t, tt.want.links, tt.args.links.Links)
		})
	}
}

var (
	idpUserLinksQuery = regexp.QuoteMeta(`SELECT projections.idp_user_links3.idp_id,` +
		` projections.idp_user_links3.user_id,` +
		` projections.idp_templates6.name,` +
		` projections.idp_user_links3.external_user_id,` +
		` projections.idp_user_links3.display_name,` +
		` projections.idp_templates6.type,` +
		` projections.idp_user_links3.resource_owner,` +
		` COUNT(*) OVER ()` +
		` FROM projections.idp_user_links3` +
		` LEFT JOIN projections.idp_templates6 ON projections.idp_user_links3.idp_id = projections.idp_templates6.id AND projections.idp_user_links3.instance_id = projections.idp_templates6.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`)
	idpUserLinksCols = []string{
		"idp_id",
		"user_id",
		"name",
		"external_user_id",
		"display_name",
		"type",
		"resource_owner",
		"count",
	}
)

func Test_IDPUserLinkPrepares(t *testing.T) {
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
			name:    "prepareIDPsQuery found",
			prepare: prepareIDPUserLinksQuery,
			want: want{
				sqlExpectations: mockQueries(
					idpUserLinksQuery,
					idpUserLinksCols,
					[][]driver.Value{
						{
							"idp-id",
							"user-id",
							"idp-name",
							"external-user-id",
							"display-name",
							domain.IDPTypeJWT,
							"ro",
						},
					},
				),
			},
			object: &IDPUserLinks{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Links: []*IDPUserLink{
					{
						IDPID:            "idp-id",
						UserID:           "user-id",
						IDPName:          "idp-name",
						ProvidedUserID:   "external-user-id",
						ProvidedUsername: "display-name",
						IDPType:          domain.IDPTypeJWT,
						ResourceOwner:    "ro",
					},
				},
			},
		},
		{
			name:    "prepareIDPsQuery no idp",
			prepare: prepareIDPUserLinksQuery,
			want: want{
				sqlExpectations: mockQueries(
					idpUserLinksQuery,
					idpUserLinksCols,
					[][]driver.Value{
						{
							"idp-id",
							"user-id",
							nil,
							"external-user-id",
							"display-name",
							nil,
							"ro",
						},
					},
				),
			},
			object: &IDPUserLinks{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Links: []*IDPUserLink{
					{
						IDPID:            "idp-id",
						UserID:           "user-id",
						IDPName:          "",
						ProvidedUserID:   "external-user-id",
						ProvidedUsername: "display-name",
						IDPType:          domain.IDPTypeUnspecified,
						ResourceOwner:    "ro",
					},
				},
			},
		},
		{
			name:    "prepareIDPsQuery sql err",
			prepare: prepareIDPUserLinksQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					idpUserLinksQuery,
					sql.ErrConnDone,
				),
				err: func(err error) (error, bool) {
					if !errors.Is(err, sql.ErrConnDone) {
						return fmt.Errorf("err should be sql.ErrConnDone got: %w", err), false
					}
					return nil, true
				},
			},
			object: (*IDPUserLinks)(nil),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}
