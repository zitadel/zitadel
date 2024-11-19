//go:build integration

package oidc_test

import (
	"context"
	"net/url"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2beta"
	session "github.com/zitadel/zitadel/pkg/grpc/session/v2beta"
)

var (
	CTX      context.Context
	Instance *integration.Instance
	Client   oidc_pb.OIDCServiceClient
)

const (
	redirectURI         = "oidcintegrationtest://callback"
	redirectURIImplicit = "http://localhost:9999/callback"
	logoutRedirectURI   = "oidcintegrationtest://logged-out"
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)
		Client = Instance.Client.OIDCv2beta

		CTX = Instance.WithAuthorization(ctx, integration.UserTypeOrgOwner)
		return m.Run()
	}())
}

func TestServer_GetAuthRequest(t *testing.T) {
	project, err := Instance.CreateProject(CTX)
	require.NoError(t, err)
	client, err := Instance.CreateOIDCNativeClient(CTX, redirectURI, logoutRedirectURI, project.GetId(), false)
	require.NoError(t, err)
	authRequestID, err := Instance.CreateOIDCAuthRequest(CTX, client.GetClientId(), Instance.Users.Get(integration.UserTypeOrgOwner).ID, redirectURI)
	require.NoError(t, err)
	now := time.Now()

	tests := []struct {
		name          string
		AuthRequestID string
		want          *oidc_pb.GetAuthRequestResponse
		wantErr       bool
	}{
		{
			name:          "Not found",
			AuthRequestID: "123",
			wantErr:       true,
		},
		{
			name:          "success",
			AuthRequestID: authRequestID,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.GetAuthRequest(CTX, &oidc_pb.GetAuthRequestRequest{
				AuthRequestId: tt.AuthRequestID,
			})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			authRequest := got.GetAuthRequest()
			assert.NotNil(t, authRequest)
			assert.Equal(t, authRequestID, authRequest.GetId())
			assert.WithinRange(t, authRequest.GetCreationDate().AsTime(), now.Add(-time.Second), now.Add(time.Second))
			assert.Contains(t, authRequest.GetScope(), "openid")
		})
	}
}

func TestServer_CreateCallback(t *testing.T) {
	project, err := Instance.CreateProject(CTX)
	require.NoError(t, err)
	client, err := Instance.CreateOIDCNativeClient(CTX, redirectURI, logoutRedirectURI, project.GetId(), false)
	require.NoError(t, err)
	sessionResp, err := Instance.Client.SessionV2beta.CreateSession(CTX, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{
					UserId: Instance.Users.Get(integration.UserTypeOrgOwner).ID,
				},
			},
		},
	})
	require.NoError(t, err)

	tests := []struct {
		name      string
		req       *oidc_pb.CreateCallbackRequest
		AuthError string
		want      *oidc_pb.CreateCallbackResponse
		wantURL   *url.URL
		wantErr   bool
	}{
		{
			name: "Not found",
			req: &oidc_pb.CreateCallbackRequest{
				AuthRequestId: "123",
				CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
					Session: &oidc_pb.Session{
						SessionId:    sessionResp.GetSessionId(),
						SessionToken: sessionResp.GetSessionToken(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "session not found",
			req: &oidc_pb.CreateCallbackRequest{
				AuthRequestId: func() string {
					authRequestID, err := Instance.CreateOIDCAuthRequest(CTX, client.GetClientId(), Instance.Users.Get(integration.UserTypeOrgOwner).ID, redirectURI)
					require.NoError(t, err)
					return authRequestID
				}(),
				CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
					Session: &oidc_pb.Session{
						SessionId:    "foo",
						SessionToken: "bar",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "session token invalid",
			req: &oidc_pb.CreateCallbackRequest{
				AuthRequestId: func() string {
					authRequestID, err := Instance.CreateOIDCAuthRequest(CTX, client.GetClientId(), Instance.Users.Get(integration.UserTypeOrgOwner).ID, redirectURI)
					require.NoError(t, err)
					return authRequestID
				}(),
				CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
					Session: &oidc_pb.Session{
						SessionId:    sessionResp.GetSessionId(),
						SessionToken: "bar",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fail callback",
			req: &oidc_pb.CreateCallbackRequest{
				AuthRequestId: func() string {
					authRequestID, err := Instance.CreateOIDCAuthRequest(CTX, client.GetClientId(), Instance.Users.Get(integration.UserTypeOrgOwner).ID, redirectURI)
					require.NoError(t, err)
					return authRequestID
				}(),
				CallbackKind: &oidc_pb.CreateCallbackRequest_Error{
					Error: &oidc_pb.AuthorizationError{
						Error:            oidc_pb.ErrorReason_ERROR_REASON_ACCESS_DENIED,
						ErrorDescription: gu.Ptr("nope"),
						ErrorUri:         gu.Ptr("https://example.com/docs"),
					},
				},
			},
			want: &oidc_pb.CreateCallbackResponse{
				CallbackUrl: regexp.QuoteMeta(`oidcintegrationtest://callback?error=access_denied&error_description=nope&error_uri=https%3A%2F%2Fexample.com%2Fdocs&state=state`),
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
			wantErr: false,
		},
		{
			name: "code callback",
			req: &oidc_pb.CreateCallbackRequest{
				AuthRequestId: func() string {
					authRequestID, err := Instance.CreateOIDCAuthRequest(CTX, client.GetClientId(), Instance.Users.Get(integration.UserTypeOrgOwner).ID, redirectURI)
					require.NoError(t, err)
					return authRequestID
				}(),
				CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
					Session: &oidc_pb.Session{
						SessionId:    sessionResp.GetSessionId(),
						SessionToken: sessionResp.GetSessionToken(),
					},
				},
			},
			want: &oidc_pb.CreateCallbackResponse{
				CallbackUrl: `oidcintegrationtest:\/\/callback\?code=(.*)&state=state`,
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
			wantErr: false,
		},
		{
			name: "implicit",
			req: &oidc_pb.CreateCallbackRequest{
				AuthRequestId: func() string {
					client, err := Instance.CreateOIDCImplicitFlowClient(CTX, redirectURIImplicit)
					require.NoError(t, err)
					authRequestID, err := Instance.CreateOIDCAuthRequestImplicit(CTX, client.GetClientId(), Instance.Users.Get(integration.UserTypeOrgOwner).ID, redirectURIImplicit)
					require.NoError(t, err)
					return authRequestID
				}(),
				CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
					Session: &oidc_pb.Session{
						SessionId:    sessionResp.GetSessionId(),
						SessionToken: sessionResp.GetSessionToken(),
					},
				},
			},
			want: &oidc_pb.CreateCallbackResponse{
				CallbackUrl: `http:\/\/localhost:9999\/callback#access_token=(.*)&expires_in=(.*)&id_token=(.*)&state=state&token_type=Bearer`,
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.CreateCallback(CTX, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
			if tt.want != nil {
				assert.Regexp(t, regexp.MustCompile(tt.want.CallbackUrl), got.GetCallbackUrl())
			}
		})
	}
}
