package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/caos/zitadel/internal/domain"
)

var (
	userIDPLinksQuery = regexp.QuoteMeta(`SELECT zitadel.projections.idp_user_links.idp_id,` +
		` zitadel.projections.idp_user_links.user_id,` +
		` zitadel.projections.idps.name,` +
		` zitadel.projections.idp_user_links.external_user_id,` +
		` zitadel.projections.idp_user_links.display_name,` +
		` zitadel.projections.idps.type,` +
		` COUNT(*) OVER ()` +
		` FROM zitadel.projections.idp_user_links` +
		` LEFT JOIN zitadel.projections.idps ON zitadel.projections.idp_user_links.idp_id = zitadel.projections.idps.id`)
	userIDPLinksCols = []string{
		"idp_id",
		"user_id",
		"name",
		"external_user_id",
		"display_name",
		"type",
		"count",
	}
)

func Test_UserIDPLinkPrepares(t *testing.T) {
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
			prepare: prepareUserIDPLinksQuery,
			want: want{
				sqlExpectations: mockQueries(
					userIDPLinksQuery,
					userIDPLinksCols,
					[][]driver.Value{
						{
							"idp-id",
							"user-id",
							"idp-name",
							"external-user-id",
							"display-name",
							domain.IDPConfigTypeJWT,
						},
					},
				),
			},
			object: &UserIDPLinks{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Links: []*UserIDPLink{
					{
						IDPID:            "idp-id",
						UserID:           "user-id",
						IDPName:          "idp-name",
						ProvidedUserID:   "external-user-id",
						ProvidedUsername: "display-name",
						IDPType:          domain.IDPConfigTypeJWT,
					},
				},
			},
		},
		{
			name:    "prepareIDPsQuery no idp",
			prepare: prepareUserIDPLinksQuery,
			want: want{
				sqlExpectations: mockQueries(
					userIDPLinksQuery,
					userIDPLinksCols,
					[][]driver.Value{
						{
							"idp-id",
							"user-id",
							nil,
							"external-user-id",
							"display-name",
							nil,
						},
					},
				),
			},
			object: &UserIDPLinks{
				SearchResponse: SearchResponse{
					Count: 1,
				},
				Links: []*UserIDPLink{
					{
						IDPID:            "idp-id",
						UserID:           "user-id",
						IDPName:          "",
						ProvidedUserID:   "external-user-id",
						ProvidedUsername: "display-name",
						IDPType:          domain.IDPConfigTypeUnspecified,
					},
				},
			},
		},
		{
			name:    "prepareIDPsQuery sql err",
			prepare: prepareUserIDPLinksQuery,
			want: want{
				sqlExpectations: mockQueryErr(
					userIDPLinksQuery,
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
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err)
		})
	}
}
