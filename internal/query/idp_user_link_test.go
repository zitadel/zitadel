package query

import (
	"database/sql"
	"database/sql/driver"
	"errors"
	"fmt"
	"regexp"
	"testing"

	"github.com/zitadel/zitadel/internal/domain"
)

var (
	idpUserLinksQuery = regexp.QuoteMeta(`SELECT projections.idp_user_links3.idp_id,` +
		` projections.idp_user_links3.user_id,` +
		` projections.idp_templates3.name,` +
		` projections.idp_user_links3.external_user_id,` +
		` projections.idp_user_links3.display_name,` +
		` projections.idp_templates3.type,` +
		` projections.idp_user_links3.resource_owner,` +
		` COUNT(*) OVER ()` +
		` FROM projections.idp_user_links3` +
		` LEFT JOIN projections.idp_templates3 ON projections.idp_user_links3.idp_id = projections.idp_templates3.id AND projections.idp_user_links3.instance_id = projections.idp_templates3.instance_id` +
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
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}
