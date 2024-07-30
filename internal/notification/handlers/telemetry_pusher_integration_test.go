//go:build integration

package handlers_test

import (
	"bytes"
	"context"
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/app"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/object"
	oidc_v2 "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
	"github.com/zitadel/zitadel/pkg/grpc/project"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

func TestServer_TelemetryPushMilestones(t *testing.T) {
	primaryDomain, instanceID, adminID, iamOwnerCtx := Tester.UseIsolatedInstance(t, CTX, SystemCTX)
	t.Log("testing against instance with primary domain", primaryDomain)
	awaitMilestone(t, Tester.MilestoneChan, primaryDomain, "InstanceCreated")

	projectAdded, err := Tester.Client.Mgmt.AddProject(iamOwnerCtx, &management.AddProjectRequest{Name: "integration"})
	require.NoError(t, err)
	awaitMilestone(t, Tester.MilestoneChan, primaryDomain, "ProjectCreated")

	redirectURI := "http://localhost:8888"
	application, err := Tester.Client.Mgmt.AddOIDCApp(iamOwnerCtx, &management.AddOIDCAppRequest{
		ProjectId:       projectAdded.GetId(),
		Name:            "integration",
		RedirectUris:    []string{redirectURI},
		ResponseTypes:   []app.OIDCResponseType{app.OIDCResponseType_OIDC_RESPONSE_TYPE_CODE},
		GrantTypes:      []app.OIDCGrantType{app.OIDCGrantType_OIDC_GRANT_TYPE_AUTHORIZATION_CODE},
		AppType:         app.OIDCAppType_OIDC_APP_TYPE_WEB,
		AuthMethodType:  app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE,
		DevMode:         true,
		AccessTokenType: app.OIDCTokenType_OIDC_TOKEN_TYPE_JWT,
	})
	require.NoError(t, err)
	awaitMilestone(t, Tester.MilestoneChan, primaryDomain, "ApplicationCreated")

	// create the session to be used for the authN of the clients
	sessionID, sessionToken, _, _ := Tester.CreatePasswordSession(t, iamOwnerCtx, adminID, "Password1!")

	console := consoleOIDCConfig(iamOwnerCtx, t)
	loginToClient(iamOwnerCtx, t, primaryDomain, console.GetClientId(), instanceID, console.GetRedirectUris()[0], sessionID, sessionToken)
	awaitMilestone(t, Tester.MilestoneChan, primaryDomain, "AuthenticationSucceededOnInstance")

	// make sure the client has been projected
	require.EventuallyWithT(t, func(collectT *assert.CollectT) {
		_, err := Tester.Client.Mgmt.GetAppByID(iamOwnerCtx, &management.GetAppByIDRequest{
			ProjectId: projectAdded.GetId(),
			AppId:     application.GetAppId(),
		})
		assert.NoError(collectT, err)
	}, 1*time.Minute, 100*time.Millisecond, "app not found")
	loginToClient(iamOwnerCtx, t, primaryDomain, application.GetClientId(), instanceID, redirectURI, sessionID, sessionToken)
	awaitMilestone(t, Tester.MilestoneChan, primaryDomain, "AuthenticationSucceededOnApplication")

	_, err = Tester.Client.System.RemoveInstance(SystemCTX, &system.RemoveInstanceRequest{InstanceId: instanceID})
	require.NoError(t, err)
	awaitMilestone(t, Tester.MilestoneChan, primaryDomain, "InstanceDeleted")
}

func loginToClient(iamOwnerCtx context.Context, t *testing.T, primaryDomain, clientID, instanceID, redirectURI, sessionID, sessionToken string) {
	authRequestID, err := Tester.CreateOIDCAuthRequestWithDomain(iamOwnerCtx, primaryDomain, clientID, Tester.Users.Get(instanceID, integration.IAMOwner).ID, redirectURI, "openid")
	require.NoError(t, err)
	callback, err := Tester.Client.OIDCv2.CreateCallback(iamOwnerCtx, &oidc_v2.CreateCallbackRequest{
		AuthRequestId: authRequestID,
		CallbackKind: &oidc_v2.CreateCallbackRequest_Session{Session: &oidc_v2.Session{
			SessionId:    sessionID,
			SessionToken: sessionToken,
		}},
	})
	require.NoError(t, err)
	provider, err := Tester.CreateRelyingPartyForDomain(iamOwnerCtx, primaryDomain, clientID, redirectURI)
	require.NoError(t, err)
	callbackURL, err := url.Parse(callback.GetCallbackUrl())
	require.NoError(t, err)
	code := callbackURL.Query().Get("code")
	_, err = rp.CodeExchange[*oidc.IDTokenClaims](iamOwnerCtx, code, provider, rp.WithCodeVerifier(integration.CodeVerifier))
	require.NoError(t, err)
}

func consoleOIDCConfig(iamOwnerCtx context.Context, t *testing.T) *app.OIDCConfig {
	projects, err := Tester.Client.Mgmt.ListProjects(iamOwnerCtx, &management.ListProjectsRequest{
		Queries: []*project.ProjectQuery{
			{
				Query: &project.ProjectQuery_NameQuery{
					NameQuery: &project.ProjectNameQuery{
						Name:   "ZITADEL",
						Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.Len(t, projects.GetResult(), 1)
	apps, err := Tester.Client.Mgmt.ListApps(iamOwnerCtx, &management.ListAppsRequest{
		ProjectId: projects.GetResult()[0].GetId(),
		Queries: []*app.AppQuery{
			{
				Query: &app.AppQuery_NameQuery{
					NameQuery: &app.AppNameQuery{
						Name:   "Console",
						Method: object.TextQueryMethod_TEXT_QUERY_METHOD_EQUALS,
					},
				},
			},
		},
	})
	require.NoError(t, err)
	require.Len(t, apps.GetResult(), 1)
	return apps.GetResult()[0].GetOidcConfig()
}

func awaitMilestone(t *testing.T, bodies chan []byte, primaryDomain, expectMilestoneType string) {
	for {
		select {
		case body := <-bodies:
			plain := new(bytes.Buffer)
			if err := json.Indent(plain, body, "", "  "); err != nil {
				t.Fatal(err)
			}
			t.Log("received milestone", plain.String())
			milestone := struct {
				Type          string `json:"type"`
				PrimaryDomain string `json:"primaryDomain"`
			}{}
			if err := json.Unmarshal(body, &milestone); err != nil {
				t.Error(err)
			}
			if milestone.Type == expectMilestoneType && milestone.PrimaryDomain == primaryDomain {
				return
			}
		case <-time.After(60 * time.Second):
			t.Fatalf("timed out waiting for milestone %s in domain %s", expectMilestoneType, primaryDomain)
		}
	}
}
