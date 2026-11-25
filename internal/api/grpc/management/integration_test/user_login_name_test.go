//go:build integration

package management_test

import (
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/object"
	"github.com/zitadel/zitadel/pkg/grpc/user"
)

// TestServer_ListUsers_LoginNamePattern lists the organization's users by a
// login name pattern, as the console's user search does.
func TestServer_ListUsers_LoginNamePattern(t *testing.T) {
	t.Parallel()

	// Login names are the usernames, which are lowercase.
	token := integration.RandString(10)
	prefix := "ln" + token
	alice := prefix + "-alice@example.com"
	bob := prefix + "-bob@example.com"
	for _, username := range []string{alice, bob, "carol" + integration.RandString(10) + "@example.com"} {
		Instance.CreateHumanUserVerified(IAMOwnerCTX, Instance.DefaultOrg.GetId(), username, integration.Phone())
	}
	// A matching user of another organization is not listed.
	otherOrg := Instance.CreateOrganization(IAMOwnerCTX, integration.OrganizationName(), integration.Email())
	Instance.CreateHumanUserVerified(IAMOwnerCTX, otherOrg.GetOrganizationId(), prefix+"-dave@example.com", integration.Phone())

	tests := []struct {
		name      string
		loginName string
		method    object.TextQueryMethod
		want      []string
	}{
		{
			name:      "contains ignore case",
			loginName: strings.ToUpper(prefix),
			method:    object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS_IGNORE_CASE,
			want:      []string{alice, bob},
		},
		{
			name:      "contains",
			loginName: prefix,
			method:    object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS,
			want:      []string{alice, bob},
		},
		{
			name:      "contains is case-sensitive",
			loginName: strings.ToUpper(prefix),
			method:    object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS,
		},
		{
			name:      "contains escapes wildcards",
			loginName: prefix + "_alice",
			method:    object.TextQueryMethod_TEXT_QUERY_METHOD_CONTAINS,
		},
		{
			name:      "starts with",
			loginName: prefix + "-a",
			method:    object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH,
			want:      []string{alice},
		},
		{
			name:      "starts with ignore case",
			loginName: strings.ToUpper(prefix + "-b"),
			method:    object.TextQueryMethod_TEXT_QUERY_METHOD_STARTS_WITH_IGNORE_CASE,
			want:      []string{bob},
		},
		{
			name:      "ends with",
			loginName: token + "-bob@example.com",
			method:    object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH,
			want:      []string{bob},
		},
		{
			name:      "ends with ignore case",
			loginName: strings.ToUpper(token + "-alice@example.com"),
			method:    object.TextQueryMethod_TEXT_QUERY_METHOD_ENDS_WITH_IGNORE_CASE,
			want:      []string{alice},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := &management.ListUsersRequest{
				Queries: []*user.SearchQuery{{
					Query: &user.SearchQuery_LoginNameQuery{
						LoginNameQuery: &user.LoginNameQuery{
							LoginName: tt.loginName,
							Method:    tt.method,
						},
					},
				}},
			}
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(OrgCTX, 20*time.Second)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.ListUsers(OrgCTX, req)
				require.NoError(ttt, err)

				usernames := make([]string, 0, len(got.GetResult()))
				for _, u := range got.GetResult() {
					usernames = append(usernames, u.GetUserName())
				}
				assert.ElementsMatch(ttt, tt.want, usernames)
				assert.EqualValues(ttt, len(tt.want), got.GetDetails().GetTotalResult())
			}, retryDuration, tick, "timeout waiting for expected user result")
		})
	}
}
