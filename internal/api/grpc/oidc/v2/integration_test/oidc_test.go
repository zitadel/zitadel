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
	"github.com/zitadel/zitadel/pkg/grpc/app"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

var (
	CTX            context.Context
	CTXLoginClient context.Context
	Instance       *integration.Instance
	Client         oidc_pb.OIDCServiceClient
	loginV2        = &app.LoginVersion{Version: &app.LoginVersion_LoginV2{LoginV2: &app.LoginV2{BaseUri: nil}}}
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
		Client = Instance.Client.OIDCv2

		CTX = Instance.WithAuthorization(ctx, integration.UserTypeOrgOwner)
		CTXLoginClient = Instance.WithAuthorization(ctx, integration.UserTypeLogin)
		return m.Run()
	}())
}

func TestServer_GetAuthRequest(t *testing.T) {
	project, err := Instance.CreateProject(CTX)
	require.NoError(t, err)
	client, err := Instance.CreateOIDCNativeClient(CTX, redirectURI, logoutRedirectURI, project.GetId(), false)
	require.NoError(t, err)

	now := time.Now()

	tests := []struct {
		name          string
		AuthRequestID string
		ctx           context.Context
		want          *oidc_pb.GetAuthRequestResponse
		wantErr       bool
	}{
		{
			name:          "Not found",
			AuthRequestID: "123",
			ctx:           CTX,
			wantErr:       true,
		},
		{
			name: "success",
			AuthRequestID: func() string {
				authRequestID, err := Instance.CreateOIDCAuthRequest(CTX, client.GetClientId(), Instance.Users[integration.UserTypeOrgOwner].ID, redirectURI)
				require.NoError(t, err)
				return authRequestID
			}(),
			ctx: CTX,
		},
		{
			name: "without login client, no permission",
			AuthRequestID: func() string {
				client, err := Instance.CreateOIDCClientLoginVersion(CTX, redirectURI, logoutRedirectURI, project.GetId(), app.OIDCAppType_OIDC_APP_TYPE_NATIVE, app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE, false, loginV2)
				require.NoError(t, err)
				authRequestID, err := Instance.CreateOIDCAuthRequestWithoutLoginClientHeader(CTX, client.GetClientId(), redirectURI, "")
				require.NoError(t, err)
				return authRequestID
			}(),
			ctx:     CTX,
			wantErr: true,
		},
		{
			name: "without login client, with permission",
			AuthRequestID: func() string {
				client, err := Instance.CreateOIDCClientLoginVersion(CTX, redirectURI, logoutRedirectURI, project.GetId(), app.OIDCAppType_OIDC_APP_TYPE_NATIVE, app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE, false, loginV2)
				require.NoError(t, err)
				authRequestID, err := Instance.CreateOIDCAuthRequestWithoutLoginClientHeader(CTX, client.GetClientId(), redirectURI, "")
				require.NoError(t, err)
				return authRequestID
			}(),
			ctx: CTXLoginClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.GetAuthRequest(tt.ctx, &oidc_pb.GetAuthRequestRequest{
				AuthRequestId: tt.AuthRequestID,
			})
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			authRequest := got.GetAuthRequest()
			assert.NotNil(t, authRequest)
			assert.Equal(t, tt.AuthRequestID, authRequest.GetId())
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
	clientV2, err := Instance.CreateOIDCClientLoginVersion(CTX, redirectURI, logoutRedirectURI, project.GetId(), app.OIDCAppType_OIDC_APP_TYPE_NATIVE, app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE, false, loginV2)
	require.NoError(t, err)
	sessionResp, err := Instance.Client.SessionV2.CreateSession(CTX, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{
					UserId: Instance.Users[integration.UserTypeOrgOwner].ID,
				},
			},
		},
	})
	require.NoError(t, err)

	tests := []struct {
		name      string
		ctx       context.Context
		req       *oidc_pb.CreateCallbackRequest
		AuthError string
		want      *oidc_pb.CreateCallbackResponse
		wantURL   *url.URL
		wantErr   bool
	}{
		{
			name: "Not found",
			ctx:  CTX,
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
			ctx:  CTX,
			req: &oidc_pb.CreateCallbackRequest{
				AuthRequestId: func() string {
					authRequestID, err := Instance.CreateOIDCAuthRequest(CTX, client.GetClientId(), Instance.Users[integration.UserTypeOrgOwner].ID, redirectURI)
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
			ctx:  CTX,
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
			ctx:  CTX,
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
			name: "fail callback, no login client header",
			ctx:  CTXLoginClient,
			req: &oidc_pb.CreateCallbackRequest{
				AuthRequestId: func() string {
					authRequestID, err := Instance.CreateOIDCAuthRequestWithoutLoginClientHeader(CTX, clientV2.GetClientId(), redirectURI, "")
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
			ctx:  CTX,
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
			name: "code callback, no login client header, no permission, error",
			ctx:  CTX,
			req: &oidc_pb.CreateCallbackRequest{
				AuthRequestId: func() string {
					authRequestID, err := Instance.CreateOIDCAuthRequestWithoutLoginClientHeader(CTX, clientV2.GetClientId(), redirectURI, "")
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
			wantErr: true,
		},
		{
			name: "code callback, no login client header, with permission",
			ctx:  CTXLoginClient,
			req: &oidc_pb.CreateCallbackRequest{
				AuthRequestId: func() string {
					authRequestID, err := Instance.CreateOIDCAuthRequestWithoutLoginClientHeader(CTX, clientV2.GetClientId(), redirectURI, "")
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
			ctx:  CTX,
			req: &oidc_pb.CreateCallbackRequest{
				AuthRequestId: func() string {
					client, err := Instance.CreateOIDCImplicitFlowClient(CTX, redirectURIImplicit, nil)
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
		{
			name: "implicit, no login client header",
			ctx:  CTXLoginClient,
			req: &oidc_pb.CreateCallbackRequest{
				AuthRequestId: func() string {
					clientV2, err := Instance.CreateOIDCImplicitFlowClient(CTX, redirectURIImplicit, loginV2)
					require.NoError(t, err)
					authRequestID, err := Instance.CreateOIDCAuthRequestImplicitWithoutLoginClientHeader(CTX, clientV2.GetClientId(), redirectURIImplicit)
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
			got, err := Client.CreateCallback(tt.ctx, tt.req)
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
