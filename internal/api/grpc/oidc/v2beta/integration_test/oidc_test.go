//go:build integration

package oidc_test

import (
	"context"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/app"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2beta"
	session "github.com/zitadel/zitadel/pkg/grpc/session/v2beta"
)

func TestServer_GetAuthRequest(t *testing.T) {
	project, err := Instance.CreateProject(CTX)
	require.NoError(t, err)
	client, err := Instance.CreateOIDCNativeClient(CTX, redirectURI, logoutRedirectURI, project.GetId(), false)
	require.NoError(t, err)

	tests := []struct {
		name    string
		dep     func() (time.Time, string, error)
		ctx     context.Context
		want    *oidc_pb.GetAuthRequestResponse
		wantErr bool
	}{
		{
			name: "Not found",
			dep: func() (time.Time, string, error) {
				return time.Now(), "123", nil
			},
			ctx:     CTX,
			wantErr: true,
		},
		{
			name: "success",
			dep: func() (time.Time, string, error) {
				return Instance.CreateOIDCAuthRequest(CTX, client.GetClientId(), Instance.Users[integration.UserTypeOrgOwner].ID, redirectURI)
			},
			ctx: CTX,
		},
		{
			name: "without login client, no permission",
			dep: func() (time.Time, string, error) {
				client, err := Instance.CreateOIDCClientLoginVersion(CTX, redirectURI, logoutRedirectURI, project.GetId(), app.OIDCAppType_OIDC_APP_TYPE_NATIVE, app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE, false, loginV2)
				require.NoError(t, err)
				return Instance.CreateOIDCAuthRequestWithoutLoginClientHeader(CTX, client.GetClientId(), redirectURI, "")
			},
			ctx:     CTX,
			wantErr: true,
		},
		{
			name: "without login client, with permission",
			dep: func() (time.Time, string, error) {
				client, err := Instance.CreateOIDCClientLoginVersion(CTX, redirectURI, logoutRedirectURI, project.GetId(), app.OIDCAppType_OIDC_APP_TYPE_NATIVE, app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE, false, loginV2)
				require.NoError(t, err)
				return Instance.CreateOIDCAuthRequestWithoutLoginClientHeader(CTX, client.GetClientId(), redirectURI, "")

			},
			ctx: CTXLoginClient,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now, authRequestID, err := tt.dep()
			require.NoError(t, err)

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.GetAuthRequest(tt.ctx, &oidc_pb.GetAuthRequestRequest{
					AuthRequestId: authRequestID,
				})
				if tt.wantErr {
					assert.Error(ttt, err)
					return
				}
				assert.NoError(ttt, err)
				authRequest := got.GetAuthRequest()
				assert.NotNil(ttt, authRequest)
				assert.Equal(ttt, authRequestID, authRequest.GetId())
				assert.WithinRange(ttt, authRequest.GetCreationDate().AsTime(), now.Add(-time.Second), now.Add(time.Second))
				assert.Contains(ttt, authRequest.GetScope(), "openid")
			}, retryDuration, tick)
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
	sessionResp := createSession(t, CTX, Instance.Users[integration.UserTypeOrgOwner].ID)

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
					_, authRequestID, err := Instance.CreateOIDCAuthRequest(CTX, client.GetClientId(), Instance.Users[integration.UserTypeOrgOwner].ID, redirectURI)
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
					_, authRequestID, err := Instance.CreateOIDCAuthRequest(CTX, client.GetClientId(), Instance.Users.Get(integration.UserTypeOrgOwner).ID, redirectURI)
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
					_, authRequestID, err := Instance.CreateOIDCAuthRequest(CTX, client.GetClientId(), Instance.Users.Get(integration.UserTypeOrgOwner).ID, redirectURI)
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
					_, authRequestID, err := Instance.CreateOIDCAuthRequestWithoutLoginClientHeader(CTX, clientV2.GetClientId(), redirectURI, "")
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
					_, authRequestID, err := Instance.CreateOIDCAuthRequest(CTX, client.GetClientId(), Instance.Users.Get(integration.UserTypeOrgOwner).ID, redirectURI)
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
					_, authRequestID, err := Instance.CreateOIDCAuthRequestWithoutLoginClientHeader(CTX, clientV2.GetClientId(), redirectURI, "")
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
					_, authRequestID, err := Instance.CreateOIDCAuthRequestWithoutLoginClientHeader(CTX, clientV2.GetClientId(), redirectURI, "")
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

func TestServer_CreateCallback_Permission(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		dep     func(ctx context.Context, t *testing.T) *oidc_pb.CreateCallbackRequest
		want    *oidc_pb.CreateCallbackResponse
		wantURL *url.URL
		wantErr bool
	}{
		{
			name: "usergrant to project and different resourceowner with different project grant",
			ctx:  CTX,
			dep: func(ctx context.Context, t *testing.T) *oidc_pb.CreateCallbackRequest {
				projectID, clientID := createOIDCApplication(ctx, t, true, true)
				projectID2, _ := createOIDCApplication(ctx, t, true, true)

				orgResp := Instance.CreateOrganization(ctx, "oidc-permission-"+gofakeit.AppName(), gofakeit.Email())
				Instance.CreateProjectGrant(ctx, projectID2, orgResp.GetOrganizationId())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), gofakeit.Email(), gofakeit.Phone())
				Instance.CreateProjectUserGrant(t, ctx, projectID, user.GetUserId())

				return createSessionAndAuthRequestForCallback(ctx, t, clientID, Instance.Users.Get(integration.UserTypeOrgOwner).ID, user.GetUserId())
			},
			wantErr: true,
		},
		{
			name: "usergrant to project and different resourceowner with project grant",
			ctx:  CTX,
			dep: func(ctx context.Context, t *testing.T) *oidc_pb.CreateCallbackRequest {
				projectID, clientID := createOIDCApplication(ctx, t, true, true)

				orgResp := Instance.CreateOrganization(ctx, "oidc-permission-"+gofakeit.AppName(), gofakeit.Email())
				Instance.CreateProjectGrant(ctx, projectID, orgResp.GetOrganizationId())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), gofakeit.Email(), gofakeit.Phone())
				Instance.CreateProjectUserGrant(t, ctx, projectID, user.GetUserId())

				return createSessionAndAuthRequestForCallback(ctx, t, clientID, Instance.Users.Get(integration.UserTypeOrgOwner).ID, user.GetUserId())
			},
			want: &oidc_pb.CreateCallbackResponse{
				CallbackUrl: `oidcintegrationtest:\/\/callback\?code=(.*)&state=state`,
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
		},
		{
			name: "usergrant to project grant and different resourceowner with project grant",
			ctx:  CTX,
			dep: func(ctx context.Context, t *testing.T) *oidc_pb.CreateCallbackRequest {
				projectID, clientID := createOIDCApplication(ctx, t, true, true)

				orgResp := Instance.CreateOrganization(ctx, "oidc-permission-"+gofakeit.AppName(), gofakeit.Email())
				projectGrantResp := Instance.CreateProjectGrant(ctx, projectID, orgResp.GetOrganizationId())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), gofakeit.Email(), gofakeit.Phone())
				Instance.CreateProjectGrantUserGrant(ctx, orgResp.GetOrganizationId(), projectID, projectGrantResp.GetGrantId(), user.GetUserId())

				return createSessionAndAuthRequestForCallback(ctx, t, clientID, Instance.Users.Get(integration.UserTypeOrgOwner).ID, user.GetUserId())
			},
			want: &oidc_pb.CreateCallbackResponse{
				CallbackUrl: `oidcintegrationtest:\/\/callback\?code=(.*)&state=state`,
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
		},
		{
			name: "no usergrant and different resourceowner",
			ctx:  CTX,
			dep: func(ctx context.Context, t *testing.T) *oidc_pb.CreateCallbackRequest {
				_, clientID := createOIDCApplication(ctx, t, true, true)

				orgResp := Instance.CreateOrganization(ctx, "oidc-permission-"+gofakeit.AppName(), gofakeit.Email())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), gofakeit.Email(), gofakeit.Phone())

				return createSessionAndAuthRequestForCallback(ctx, t, clientID, Instance.Users.Get(integration.UserTypeOrgOwner).ID, user.GetUserId())
			},
			wantErr: true,
		},
		{
			name: "no usergrant and same resourceowner",
			ctx:  CTX,
			dep: func(ctx context.Context, t *testing.T) *oidc_pb.CreateCallbackRequest {
				_, clientID := createOIDCApplication(ctx, t, true, true)
				user := Instance.CreateHumanUser(ctx)

				return createSessionAndAuthRequestForCallback(ctx, t, clientID, Instance.Users.Get(integration.UserTypeOrgOwner).ID, user.GetUserId())
			},
			wantErr: true,
		},
		{
			name: "usergrant and different resourceowner",
			ctx:  CTX,
			dep: func(ctx context.Context, t *testing.T) *oidc_pb.CreateCallbackRequest {
				projectID, clientID := createOIDCApplication(ctx, t, true, true)

				orgResp := Instance.CreateOrganization(ctx, "oidc-permission-"+gofakeit.AppName(), gofakeit.Email())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), gofakeit.Email(), gofakeit.Phone())
				Instance.CreateProjectUserGrant(t, ctx, projectID, user.GetUserId())

				return createSessionAndAuthRequestForCallback(ctx, t, clientID, Instance.Users.Get(integration.UserTypeOrgOwner).ID, user.GetUserId())
			},
			wantErr: true,
		},
		{
			name: "usergrant and same resourceowner",
			ctx:  CTX,
			dep: func(ctx context.Context, t *testing.T) *oidc_pb.CreateCallbackRequest {
				projectID, clientID := createOIDCApplication(ctx, t, true, true)
				user := Instance.CreateHumanUser(ctx)
				Instance.CreateProjectUserGrant(t, ctx, projectID, user.GetUserId())

				return createSessionAndAuthRequestForCallback(ctx, t, clientID, Instance.Users.Get(integration.UserTypeOrgOwner).ID, user.GetUserId())
			},
			want: &oidc_pb.CreateCallbackResponse{
				CallbackUrl: `oidcintegrationtest:\/\/callback\?code=(.*)&state=state`,
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
		},
		{
			name: "projectRoleCheck, usergrant and same resourceowner",
			ctx:  CTX,
			dep: func(ctx context.Context, t *testing.T) *oidc_pb.CreateCallbackRequest {
				projectID, clientID := createOIDCApplication(ctx, t, true, false)
				user := Instance.CreateHumanUser(ctx)
				Instance.CreateProjectUserGrant(t, ctx, projectID, user.GetUserId())

				return createSessionAndAuthRequestForCallback(ctx, t, clientID, Instance.Users.Get(integration.UserTypeOrgOwner).ID, user.GetUserId())
			},
			want: &oidc_pb.CreateCallbackResponse{
				CallbackUrl: `oidcintegrationtest:\/\/callback\?code=(.*)&state=state`,
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
		},
		{
			name: "projectRoleCheck, no usergrant and same resourceowner",
			ctx:  CTX,
			dep: func(ctx context.Context, t *testing.T) *oidc_pb.CreateCallbackRequest {
				_, clientID := createOIDCApplication(ctx, t, true, false)
				user := Instance.CreateHumanUser(ctx)

				return createSessionAndAuthRequestForCallback(ctx, t, clientID, Instance.Users.Get(integration.UserTypeOrgOwner).ID, user.GetUserId())
			},
			wantErr: true,
		},
		{
			name: "projectRoleCheck, usergrant and different resourceowner",
			ctx:  CTX,
			dep: func(ctx context.Context, t *testing.T) *oidc_pb.CreateCallbackRequest {
				projectID, clientID := createOIDCApplication(ctx, t, true, false)
				orgResp := Instance.CreateOrganization(ctx, "oidc-permission-"+gofakeit.AppName(), gofakeit.Email())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), gofakeit.Email(), gofakeit.Phone())
				Instance.CreateProjectUserGrant(t, ctx, projectID, user.GetUserId())

				return createSessionAndAuthRequestForCallback(ctx, t, clientID, Instance.Users.Get(integration.UserTypeOrgOwner).ID, user.GetUserId())
			},
			want: &oidc_pb.CreateCallbackResponse{
				CallbackUrl: `oidcintegrationtest:\/\/callback\?code=(.*)&state=state`,
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
		},
		{
			name: "projectRoleCheck, no usergrant and different resourceowner",
			ctx:  CTX,
			dep: func(ctx context.Context, t *testing.T) *oidc_pb.CreateCallbackRequest {
				_, clientID := createOIDCApplication(ctx, t, true, false)
				orgResp := Instance.CreateOrganization(ctx, "oidc-permission-"+gofakeit.AppName(), gofakeit.Email())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), gofakeit.Email(), gofakeit.Phone())

				return createSessionAndAuthRequestForCallback(ctx, t, clientID, Instance.Users.Get(integration.UserTypeOrgOwner).ID, user.GetUserId())
			},
			wantErr: true,
		},
		{
			name: "projectRoleCheck, usergrant on project grant and different resourceowner",
			ctx:  CTX,
			dep: func(ctx context.Context, t *testing.T) *oidc_pb.CreateCallbackRequest {
				projectID, clientID := createOIDCApplication(ctx, t, true, false)

				orgResp := Instance.CreateOrganization(ctx, "oidc-permission-"+gofakeit.AppName(), gofakeit.Email())
				projectGrantResp := Instance.CreateProjectGrant(ctx, projectID, orgResp.GetOrganizationId())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), gofakeit.Email(), gofakeit.Phone())
				Instance.CreateProjectGrantUserGrant(ctx, orgResp.GetOrganizationId(), projectID, projectGrantResp.GetGrantId(), user.GetUserId())
				return createSessionAndAuthRequestForCallback(ctx, t, clientID, Instance.Users.Get(integration.UserTypeOrgOwner).ID, user.GetUserId())
			},
			want: &oidc_pb.CreateCallbackResponse{
				CallbackUrl: `oidcintegrationtest:\/\/callback\?code=(.*)&state=state`,
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
		},
		{
			name: "projectRoleCheck, no usergrant on project grant and different resourceowner",
			ctx:  CTX,
			dep: func(ctx context.Context, t *testing.T) *oidc_pb.CreateCallbackRequest {
				projectID, clientID := createOIDCApplication(ctx, t, true, false)

				orgResp := Instance.CreateOrganization(ctx, "oidc-permission-"+gofakeit.AppName(), gofakeit.Email())
				Instance.CreateProjectGrant(ctx, projectID, orgResp.GetOrganizationId())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), gofakeit.Email(), gofakeit.Phone())
				return createSessionAndAuthRequestForCallback(ctx, t, clientID, Instance.Users.Get(integration.UserTypeOrgOwner).ID, user.GetUserId())
			},
			wantErr: true,
		},
		{
			name: "hasProjectCheck, same resourceowner",
			ctx:  CTX,
			dep: func(ctx context.Context, t *testing.T) *oidc_pb.CreateCallbackRequest {
				user := Instance.CreateHumanUser(ctx)
				_, clientID := createOIDCApplication(ctx, t, false, true)

				return createSessionAndAuthRequestForCallback(ctx, t, clientID, Instance.Users.Get(integration.UserTypeOrgOwner).ID, user.GetUserId())
			},
			want: &oidc_pb.CreateCallbackResponse{
				CallbackUrl: `oidcintegrationtest:\/\/callback\?code=(.*)&state=state`,
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
		},
		{
			name: "hasProjectCheck, different resourceowner",
			ctx:  CTX,
			dep: func(ctx context.Context, t *testing.T) *oidc_pb.CreateCallbackRequest {
				_, clientID := createOIDCApplication(ctx, t, false, true)
				orgResp := Instance.CreateOrganization(ctx, "oidc-permission-"+gofakeit.AppName(), gofakeit.Email())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), gofakeit.Email(), gofakeit.Phone())

				return createSessionAndAuthRequestForCallback(ctx, t, clientID, Instance.Users.Get(integration.UserTypeOrgOwner).ID, user.GetUserId())
			},
			wantErr: true,
		},
		{
			name: "hasProjectCheck, different resourceowner with project grant",
			ctx:  CTX,
			dep: func(ctx context.Context, t *testing.T) *oidc_pb.CreateCallbackRequest {
				projectID, clientID := createOIDCApplication(ctx, t, false, true)

				orgResp := Instance.CreateOrganization(ctx, "oidc-permission-"+gofakeit.AppName(), gofakeit.Email())
				Instance.CreateProjectGrant(ctx, projectID, orgResp.GetOrganizationId())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), gofakeit.Email(), gofakeit.Phone())

				return createSessionAndAuthRequestForCallback(ctx, t, clientID, Instance.Users.Get(integration.UserTypeOrgOwner).ID, user.GetUserId())
			},
			want: &oidc_pb.CreateCallbackResponse{
				CallbackUrl: `oidcintegrationtest:\/\/callback\?code=(.*)&state=state`,
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := tt.dep(IAMCTX, t)

			got, err := Client.CreateCallback(tt.ctx, req)
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

func createSession(t *testing.T, ctx context.Context, userID string) *session.CreateSessionResponse {
	sessionResp, err := Instance.Client.SessionV2beta.CreateSession(ctx, &session.CreateSessionRequest{
		Checks: &session.Checks{
			User: &session.CheckUser{
				Search: &session.CheckUser_UserId{
					UserId: userID,
				},
			},
		},
	})
	require.NoError(t, err)
	return sessionResp
}

func createSessionAndAuthRequestForCallback(ctx context.Context, t *testing.T, clientID, loginClient, userID string) *oidc_pb.CreateCallbackRequest {
	_, authRequestID, err := Instance.CreateOIDCAuthRequest(ctx, clientID, loginClient, redirectURI)
	require.NoError(t, err)
	sessionResp := createSession(t, ctx, userID)
	return &oidc_pb.CreateCallbackRequest{
		AuthRequestId: authRequestID,
		CallbackKind: &oidc_pb.CreateCallbackRequest_Session{
			Session: &oidc_pb.Session{
				SessionId:    sessionResp.GetSessionId(),
				SessionToken: sessionResp.GetSessionToken(),
			},
		},
	}
}

func createOIDCApplication(ctx context.Context, t *testing.T, projectRoleCheck, hasProjectCheck bool) (string, string) {
	project, err := Instance.CreateProjectWithPermissionCheck(ctx, projectRoleCheck, hasProjectCheck)
	require.NoError(t, err)
	clientV2, err := Instance.CreateOIDCClientLoginVersion(ctx, redirectURI, logoutRedirectURI, project.GetId(), app.OIDCAppType_OIDC_APP_TYPE_NATIVE, app.OIDCAuthMethodType_OIDC_AUTH_METHOD_TYPE_NONE, false, loginV2)
	require.NoError(t, err)
	return project.GetId(), clientV2.GetClientId()
}
