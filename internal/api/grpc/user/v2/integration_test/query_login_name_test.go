//go:build integration

package user_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_ListUsers_LoginNamePattern(t *testing.T) {
	iamOwnerCtx := InstancePermissionV2.WithAuthorizationToken(OrgCTX, integration.UserTypeIAMOwner)
	orgResp := InstancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())
	otherOrgResp := InstancePermissionV2.CreateOrganization(iamOwnerCtx, integration.OrganizationName(), integration.Email())

	// Login names are the usernames, which are lowercase.
	token := integration.RandString(10)
	prefix := "ln" + token
	alice := createUserWithUserName(iamOwnerCtx, InstancePermissionV2, prefix+"-alice@example.com", orgResp.OrganizationId, false)
	bob := createUserWithUserName(iamOwnerCtx, InstancePermissionV2, prefix+"-bob@example.com", orgResp.OrganizationId, false)
	carol := createUserWithUserName(iamOwnerCtx, InstancePermissionV2, "carol"+integration.RandString(10)+"@example.com", orgResp.OrganizationId, false)
	dave := createUserWithUserName(iamOwnerCtx, InstancePermissionV2, prefix+"-dave@example.com", otherOrgResp.OrganizationId, false)
	orgQuery := OrganizationIdQuery(orgResp.OrganizationId)

	tests := []struct {
		name          string
		queries       []*user.SearchQuery
		wantUsernames []string
	}{
		{
			name: "org and contains ignore case",
			queries: []*user.SearchQuery{
				orgQuery,
				LoginNameQuery(strings.ToUpper(prefix), object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE),
			},
			wantUsernames: []string{alice.Username, bob.Username},
		},
		{
			name: "org and contains is case-sensitive",
			queries: []*user.SearchQuery{
				orgQuery,
				LoginNameQuery(strings.ToUpper(prefix), object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS),
			},
		},
		{
			name: "contains in all organizations",
			queries: []*user.SearchQuery{
				LoginNameQuery(prefix, object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS),
			},
			wantUsernames: []string{alice.Username, bob.Username, dave.Username},
		},
		{
			name: "org and starts with",
			queries: []*user.SearchQuery{
				orgQuery,
				LoginNameQuery(prefix+"-b", object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH),
			},
			wantUsernames: []string{bob.Username},
		},
		{
			name: "org and ends with ignore case",
			queries: []*user.SearchQuery{
				orgQuery,
				LoginNameQuery(strings.ToUpper(token+"-alice@example.com"), object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE),
			},
			wantUsernames: []string{alice.Username},
		},
		{
			name: "org and login name or username",
			queries: []*user.SearchQuery{
				orgQuery,
				OrQuery([]*user.SearchQuery{
					LoginNameQuery(strings.ToUpper(prefix+"-alice"), object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE),
					UsernameQuery(carol.Username),
				}),
			},
			wantUsernames: []string{alice.Username, carol.Username},
		},
		{
			name: "not login name",
			queries: []*user.SearchQuery{
				InUserIDsQuery([]string{alice.UserID, bob.UserID, carol.UserID, dave.UserID}),
				NotQuery(LoginNameQuery(prefix, object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS)),
			},
			wantUsernames: []string{carol.Username},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &user.ListUsersRequest{Queries: tt.queries}
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(iamOwnerCtx, 20*time.Second)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := InstancePermissionV2.Client.UserV2.ListUsers(iamOwnerCtx, req)
				require.NoError(ttt, err)

				gotUsernames := make([]string, 0, len(got.GetResult()))
				for _, u := range got.GetResult() {
					gotUsernames = append(gotUsernames, u.GetUsername())
				}
				assert.ElementsMatch(ttt, tt.wantUsernames, gotUsernames)
				assert.EqualValues(ttt, len(tt.wantUsernames), got.GetDetails().GetTotalResult())
			}, retryDuration, tick, "timeout waiting for expected user result")
		})
	}
}

func LoginNameQuery(loginName string, method object.TextQueryMethod) *user.SearchQuery {
	return &user.SearchQuery{
		Query: &user.SearchQuery_LoginNameQuery{
			LoginNameQuery: &user.LoginNameQuery{
				LoginName: loginName,
				Method:    method,
			},
		},
	}
}

func NotQuery(query *user.SearchQuery) *user.SearchQuery {
	return &user.SearchQuery{
		Query: &user.SearchQuery_NotQuery{
			NotQuery: &user.NotQuery{
				Query: query,
			},
		},
	}
}
