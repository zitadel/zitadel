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
	loginPolicyIDPLinksQuery = regexp.QuoteMeta(`SELECT projections.idp_login_policy_links4.idp_id,` +
		` projections.idp_templates3.name,` +
		` projections.idp_templates3.type,` +
		` projections.idp_templates3.owner_type,` +
		` COUNT(*) OVER ()` +
		` FROM projections.idp_login_policy_links4` +
		` LEFT JOIN projections.idp_templates3 ON projections.idp_login_policy_links4.idp_id = projections.idp_templates3.id AND projections.idp_login_policy_links4.instance_id = projections.idp_templates3.instance_id` +
		` AS OF SYSTEM TIME '-1 ms'`)
	loginPolicyIDPLinksCols = []string{
		"idp_id",
		"name",
		"type",
		"owner_type",
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
		prepare interface{}
		want    want
		object  interface{}
	}{
		{
			name:    "prepareIDPsQuery found",
			prepare: prepareIDPLoginPolicyLinksQuery,
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
						IDPID:     "idp-id",
						IDPName:   "idp-name",
						IDPType:   domain.IDPTypeJWT,
						OwnerType: domain.IdentityProviderTypeSystem,
					},
				},
			},
		},
		{
			name:    "prepareIDPsQuery no idp",
			prepare: prepareIDPLoginPolicyLinksQuery,
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
						IDPID:   "idp-id",
						IDPName: "",
						IDPType: domain.IDPTypeUnspecified,
					},
				},
			},
		},
		{
			name:    "prepareIDPsQuery sql err",
			prepare: prepareIDPLoginPolicyLinksQuery,
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
			object: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			assertPrepare(t, tt.prepare, tt.object, tt.want.sqlExpectations, tt.want.err, defaultPrepareArgs...)
		})
	}
}
