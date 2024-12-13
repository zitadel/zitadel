//go:build integration

package handlers_test

import (
	"bytes"
	"encoding/json"
	"net/url"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/oidc/v3/pkg/client/rp"
	"github.com/zitadel/oidc/v3/pkg/oidc"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/integration/sink"
	"github.com/zitadel/zitadel/internal/repository/milestone"
	"github.com/zitadel/zitadel/pkg/grpc/app"
	"github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/object"
	oidc_v2 "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
	"github.com/zitadel/zitadel/pkg/grpc/project"
	"github.com/zitadel/zitadel/pkg/grpc/system"
)

func TestServer_TelemetryPushMilestones(t *testing.T) {
	sub := sink.Subscribe(CTX, sink.ChannelMilestone)
	defer sub.Close()

	instance := integration.NewInstance(CTX)
	iamOwnerCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)
	t.Log("testing against instance", instance.ID())
	awaitMilestone(t, sub, instance.ID(), milestone.InstanceCreated)

	projectAdded, err := instance.Client.Mgmt.AddProject(iamOwnerCtx, &management.AddProjectRequest{Name: "integration"})
	require.NoError(t, err)
	awaitMilestone(t, sub, instance.ID(), milestone.ProjectCreated)

	redirectURI := "http://localhost:8888"
	application, err := instance.Client.Mgmt.AddOIDCApp(iamOwnerCtx, &management.AddOIDCAppRequest{
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
	awaitMilestone(t, sub, instance.ID(), milestone.ApplicationCreated)

	// create the session to be used for the authN of the clients
	sessionID, sessionToken, _, _ := instance.CreatePasswordSession(t, iamOwnerCtx, instance.AdminUserID, "Password1!")

	console := consoleOIDCConfig(t, instance)
	loginToClient(t, instance, console.GetClientId(), console.GetRedirectUris()[0], sessionID, sessionToken)
	awaitMilestone(t, sub, instance.ID(), milestone.AuthenticationSucceededOnInstance)

	// make sure the client has been projected
	require.EventuallyWithT(t, func(collectT *assert.CollectT) {
		_, err := instance.Client.Mgmt.GetAppByID(iamOwnerCtx, &management.GetAppByIDRequest{
			ProjectId: projectAdded.GetId(),
			AppId:     application.GetAppId(),
		})
		assert.NoError(collectT, err)
	}, time.Minute, time.Second, "app not found")
	loginToClient(t, instance, application.GetClientId(), redirectURI, sessionID, sessionToken)
	awaitMilestone(t, sub, instance.ID(), milestone.AuthenticationSucceededOnApplication)

	_, err = integration.SystemClient().RemoveInstance(CTX, &system.RemoveInstanceRequest{InstanceId: instance.ID()})
	require.NoError(t, err)
	awaitMilestone(t, sub, instance.ID(), milestone.InstanceDeleted)
}

func loginToClient(t *testing.T, instance *integration.Instance, clientID, redirectURI, sessionID, sessionToken string) {
	iamOwnerCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	authRequestID, err := instance.CreateOIDCAuthRequestWithDomain(iamOwnerCtx, instance.Domain, clientID, instance.Users.Get(integration.UserTypeIAMOwner).ID, redirectURI, "openid")
	require.NoError(t, err)
	callback, err := instance.Client.OIDCv2.CreateCallback(iamOwnerCtx, &oidc_v2.CreateCallbackRequest{
		AuthRequestId: authRequestID,
		CallbackKind: &oidc_v2.CreateCallbackRequest_Session{Session: &oidc_v2.Session{
			SessionId:    sessionID,
			SessionToken: sessionToken,
		}},
	})
	require.NoError(t, err)
	provider, err := instance.CreateRelyingPartyForDomain(iamOwnerCtx, instance.Domain, clientID, redirectURI)
	require.NoError(t, err)
	callbackURL, err := url.Parse(callback.GetCallbackUrl())
	require.NoError(t, err)
	code := callbackURL.Query().Get("code")
	_, err = rp.CodeExchange[*oidc.IDTokenClaims](iamOwnerCtx, code, provider, rp.WithCodeVerifier(integration.CodeVerifier))
	require.NoError(t, err)
}

func consoleOIDCConfig(t *testing.T, instance *integration.Instance) *app.OIDCConfig {
	iamOwnerCtx := instance.WithAuthorization(CTX, integration.UserTypeIAMOwner)

	projects, err := instance.Client.Mgmt.ListProjects(iamOwnerCtx, &management.ListProjectsRequest{
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
	apps, err := instance.Client.Mgmt.ListApps(iamOwnerCtx, &management.ListAppsRequest{
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

func awaitMilestone(t *testing.T, sub *sink.Subscription, instanceID string, expectMilestoneType milestone.Type) {
	for {
		select {
		case req := <-sub.Recv():
			plain := new(bytes.Buffer)
			if err := json.Indent(plain, req.Body, "", "  "); err != nil {
				t.Fatal(err)
			}
			t.Log("received milestone", plain.String())
			milestone := struct {
				InstanceID string         `json:"instanceId"`
				Type       milestone.Type `json:"type"`
			}{}
			if err := json.Unmarshal(req.Body, &milestone); err != nil {
				t.Error(err)
			}
			if milestone.Type == expectMilestoneType && milestone.InstanceID == instanceID {
				return
			}
		case <-time.After(20 * time.Second):
			t.Fatalf("timed out waiting for milestone %s for instance %s", expectMilestoneType, instanceID)
		}
	}
}
