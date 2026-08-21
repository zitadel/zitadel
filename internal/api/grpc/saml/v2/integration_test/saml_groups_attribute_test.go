//go:build integration

package saml_test

import (
	"encoding/base64"
	"encoding/xml"
	"testing"
	"time"

	"github.com/crewjam/saml"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/filter/v2"
	group_v2 "github.com/zitadel/zitadel/pkg/grpc/group/v2"
)

// TestServer_CreateResponse_GroupsAttribute_SpecialChars proves group names
// containing XML-special characters are escaped in the assertion and survive
// the round-trip through encoding/xml — appendCustomAttribute and the SAML
// serializer must not emit them raw.
func TestServer_CreateResponse_GroupsAttribute_SpecialChars(t *testing.T) {
	idpMetadata, err := Instance.GetSAMLIDPMetadata()
	require.NoError(t, err)
	acsPost := idpMetadata.IDPSSODescriptors[0].SingleSignOnServices[1]
	_, _, spMiddleware := createSAMLApplication(CTX, t, idpMetadata, saml.HTTPPostBinding, false, false)

	orgID := Instance.DefaultOrg.GetId()
	subject := Instance.CreateHumanUserVerified(CTX, orgID, integration.Email(), integration.Phone())
	groupSuffix := integration.GroupName()
	groupName := "Ops & Dev <" + groupSuffix + ">"
	group := Instance.CreateGroup(CTX, t, orgID, groupName)
	Instance.AddUsersToGroup(CTX, t, group.GetId(), []string{subject.GetUserId()})

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 3*time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		resp, err := Instance.Client.GroupV2.ListGroupUsers(CTX, &group_v2.ListGroupUsersRequest{
			Filters: []*group_v2.GroupUsersSearchFilter{{
				Filter: &group_v2.GroupUsersSearchFilter_GroupIds{
					GroupIds: &filter.InIDsFilter{Ids: []string{group.GetId()}},
				},
			}},
		})
		require.NoError(ttt, err)
		require.Len(ttt, resp.GetGroupUsers(), 1)
	}, retryDuration, tick, "timeout waiting for group membership projection")

	req := createSessionAndSmlRequestForCallback(LoginCTX, t, spMiddleware, Instance.Users[integration.UserTypeLogin].ID, acsPost, subject.GetUserId(), saml.HTTPPostBinding)
	resp, err := Client.CreateResponse(LoginCTX, req)
	require.NoError(t, err)
	samlResponse := resp.GetPost().GetSamlResponse()
	require.NotEmpty(t, samlResponse)

	decoded, err := base64.StdEncoding.DecodeString(samlResponse)
	require.NoError(t, err)

	// raw XML must hold the escaped form; the unescaped substring would be a parse error
	escaped := "Ops &amp; Dev &lt;" + groupSuffix + "&gt;"
	assert.Contains(t, string(decoded), escaped, "group name must be XML-escaped in the assertion")
	assert.NotContains(t, string(decoded), groupName, "raw group name must not appear unescaped")

	parsed := new(saml.Response)
	require.NoError(t, xml.Unmarshal(decoded, parsed))
	require.NotNil(t, parsed.Assertion, "assertion must be present in plain text")

	var groupValues []string
	for _, statement := range parsed.Assertion.AttributeStatements {
		for _, attribute := range statement.Attributes {
			if attribute.Name != "groups" {
				continue
			}
			for _, value := range attribute.Values {
				groupValues = append(groupValues, value.Value)
			}
		}
	}
	assert.Equal(t, []string{groupName}, groupValues, "the parser must unescape back to the original name")
}

// TestServer_CreateResponse_GroupsAttribute proves group memberships surface
// as a "groups" attribute in the assertion of a real SAML response.
func TestServer_CreateResponse_GroupsAttribute(t *testing.T) {
	idpMetadata, err := Instance.GetSAMLIDPMetadata()
	require.NoError(t, err)
	acsPost := idpMetadata.IDPSSODescriptors[0].SingleSignOnServices[1]
	_, _, spMiddleware := createSAMLApplication(CTX, t, idpMetadata, saml.HTTPPostBinding, false, false)

	orgID := Instance.DefaultOrg.GetId()
	subject := Instance.CreateHumanUserVerified(CTX, orgID, integration.Email(), integration.Phone())
	groupName := integration.GroupName()
	group := Instance.CreateGroup(CTX, t, orgID, groupName)
	Instance.AddUsersToGroup(CTX, t, group.GetId(), []string{subject.GetUserId()})

	// the SAML storage reads the same projection; wait for the membership first
	// since the auth request is consumed by a single CreateResponse call
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, 3*time.Minute)
	require.EventuallyWithT(t, func(ttt *assert.CollectT) {
		resp, err := Instance.Client.GroupV2.ListGroupUsers(CTX, &group_v2.ListGroupUsersRequest{
			Filters: []*group_v2.GroupUsersSearchFilter{{
				Filter: &group_v2.GroupUsersSearchFilter_GroupIds{
					GroupIds: &filter.InIDsFilter{Ids: []string{group.GetId()}},
				},
			}},
		})
		require.NoError(ttt, err)
		require.Len(ttt, resp.GetGroupUsers(), 1)
	}, retryDuration, tick, "timeout waiting for group membership projection")

	req := createSessionAndSmlRequestForCallback(LoginCTX, t, spMiddleware, Instance.Users[integration.UserTypeLogin].ID, acsPost, subject.GetUserId(), saml.HTTPPostBinding)
	resp, err := Client.CreateResponse(LoginCTX, req)
	require.NoError(t, err)
	samlResponse := resp.GetPost().GetSamlResponse()
	require.NotEmpty(t, samlResponse)

	decoded, err := base64.StdEncoding.DecodeString(samlResponse)
	require.NoError(t, err)
	parsed := new(saml.Response)
	require.NoError(t, xml.Unmarshal(decoded, parsed))
	require.NotNil(t, parsed.Assertion, "assertion must be present in plain text")

	var groupValues []string
	for _, statement := range parsed.Assertion.AttributeStatements {
		for _, attribute := range statement.Attributes {
			if attribute.Name != "groups" {
				continue
			}
			for _, value := range attribute.Values {
				groupValues = append(groupValues, value.Value)
			}
		}
	}
	assert.Equal(t, []string{groupName}, groupValues)
}
