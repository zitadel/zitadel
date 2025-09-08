//go:build integration

package saml_test

import (
	"context"
	"net/url"
	"regexp"
	"testing"
	"time"

	"github.com/crewjam/saml"
	"github.com/crewjam/saml/samlsp"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	mgmt "github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	saml_pb "github.com/zitadel/zitadel/pkg/grpc/saml/v2"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

func TestServer_GetSAMLRequest(t *testing.T) {
	idpMetadata, err := Instance.GetSAMLIDPMetadata()
	require.NoError(t, err)

	acsRedirect := idpMetadata.IDPSSODescriptors[0].SingleSignOnServices[0]
	acsPost := idpMetadata.IDPSSODescriptors[0].SingleSignOnServices[1]

	_, _, spMiddlewareRedirect := createSAMLApplication(CTX, t, idpMetadata, saml.HTTPRedirectBinding, false, false)
	_, _, spMiddlewarePost := createSAMLApplication(CTX, t, idpMetadata, saml.HTTPPostBinding, false, false)

	tests := []struct {
		name    string
		dep     func() (time.Time, string, error)
		wantErr bool
	}{
		{
			name: "Not found",
			dep: func() (time.Time, string, error) {
				return time.Time{}, "123", nil
			},
			wantErr: true,
		},
		{
			name: "success, redirect binding",
			dep: func() (time.Time, string, error) {
				return Instance.CreateSAMLAuthRequest(spMiddlewareRedirect, Instance.Users[integration.UserTypeLogin].ID, acsRedirect, integration.RelayState(), saml.HTTPRedirectBinding)
			},
		},
		{
			name: "success, post binding",
			dep: func() (time.Time, string, error) {
				return Instance.CreateSAMLAuthRequest(spMiddlewarePost, Instance.Users[integration.UserTypeLogin].ID, acsPost, integration.RelayState(), saml.HTTPPostBinding)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			creationTime, authRequestID, err := tt.dep()
			require.NoError(t, err)

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(LoginCTX, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.GetSAMLRequest(LoginCTX, &saml_pb.GetSAMLRequestRequest{
					SamlRequestId: authRequestID,
				})
				if tt.wantErr {
					assert.Error(ttt, err)
					return
				}
				assert.NoError(ttt, err)
				authRequest := got.GetSamlRequest()
				assert.NotNil(ttt, authRequest)
				assert.Equal(ttt, authRequestID, authRequest.GetId())
				assert.WithinRange(ttt, authRequest.GetCreationDate().AsTime(), creationTime.Add(-time.Second), creationTime.Add(time.Second))
			}, retryDuration, tick, "timeout waiting for expected saml request result")
		})
	}
}

func TestServer_CreateResponse(t *testing.T) {
	idpMetadata, err := Instance.GetSAMLIDPMetadata()
	require.NoError(t, err)
	acsRedirect := idpMetadata.IDPSSODescriptors[0].SingleSignOnServices[0]
	acsPost := idpMetadata.IDPSSODescriptors[0].SingleSignOnServices[1]

	_, rootURLPost, spMiddlewarePost := createSAMLApplication(CTX, t, idpMetadata, saml.HTTPPostBinding, false, false)
	_, rootURLRedirect, spMiddlewareRedirect := createSAMLApplication(CTX, t, idpMetadata, saml.HTTPRedirectBinding, false, false)
	sessionResp := createSession(LoginCTX, t, Instance.Users[integration.UserTypeLogin].ID)

	tests := []struct {
		name      string
		ctx       context.Context
		req       *saml_pb.CreateResponseRequest
		AuthError string
		want      *saml_pb.CreateResponseResponse
		wantURL   *url.URL
		wantErr   bool
	}{
		{
			name: "Not found",
			ctx:  LoginCTX,
			req: &saml_pb.CreateResponseRequest{
				SamlRequestId: "123",
				ResponseKind: &saml_pb.CreateResponseRequest_Session{
					Session: &saml_pb.Session{
						SessionId:    sessionResp.GetSessionId(),
						SessionToken: sessionResp.GetSessionToken(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "session not found",
			ctx:  LoginCTX,
			req: &saml_pb.CreateResponseRequest{
				SamlRequestId: func() string {
					_, authRequestID, err := Instance.CreateSAMLAuthRequest(spMiddlewareRedirect, Instance.Users[integration.UserTypeLogin].ID, acsRedirect, integration.RelayState(), saml.HTTPRedirectBinding)
					require.NoError(t, err)
					return authRequestID
				}(),
				ResponseKind: &saml_pb.CreateResponseRequest_Session{
					Session: &saml_pb.Session{
						SessionId:    "foo",
						SessionToken: "bar",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "session token invalid",
			ctx:  LoginCTX,
			req: &saml_pb.CreateResponseRequest{
				SamlRequestId: func() string {
					_, authRequestID, err := Instance.CreateSAMLAuthRequest(spMiddlewareRedirect, Instance.Users[integration.UserTypeLogin].ID, acsRedirect, integration.RelayState(), saml.HTTPRedirectBinding)
					require.NoError(t, err)
					return authRequestID
				}(),
				ResponseKind: &saml_pb.CreateResponseRequest_Session{
					Session: &saml_pb.Session{
						SessionId:    sessionResp.GetSessionId(),
						SessionToken: "bar",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fail callback, post",
			ctx:  LoginCTX,
			req: &saml_pb.CreateResponseRequest{
				SamlRequestId: func() string {
					_, authRequestID, err := Instance.CreateSAMLAuthRequest(spMiddlewarePost, Instance.Users[integration.UserTypeLogin].ID, acsPost, integration.RelayState(), saml.HTTPPostBinding)
					require.NoError(t, err)
					return authRequestID
				}(),
				ResponseKind: &saml_pb.CreateResponseRequest_Error{
					Error: &saml_pb.AuthorizationError{
						Error:            saml_pb.ErrorReason_ERROR_REASON_REQUEST_DENIED,
						ErrorDescription: gu.Ptr("nope"),
					},
				},
			},
			want: &saml_pb.CreateResponseResponse{
				Url: regexp.QuoteMeta(`https://` + rootURLPost + `/saml/acs`),
				Binding: &saml_pb.CreateResponseResponse_Post{Post: &saml_pb.PostResponse{
					RelayState:   "notempty",
					SamlResponse: "notempty",
				}},
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
			wantErr: false,
		},
		{
			name: "fail callback, post, already failed",
			ctx:  LoginCTX,
			req: &saml_pb.CreateResponseRequest{
				SamlRequestId: func() string {
					_, authRequestID, err := Instance.CreateSAMLAuthRequest(spMiddlewarePost, Instance.Users[integration.UserTypeLogin].ID, acsPost, integration.RelayState(), saml.HTTPPostBinding)
					require.NoError(t, err)
					Instance.FailSAMLAuthRequest(LoginCTX, authRequestID, saml_pb.ErrorReason_ERROR_REASON_AUTH_N_FAILED)
					return authRequestID
				}(),
				ResponseKind: &saml_pb.CreateResponseRequest_Error{
					Error: &saml_pb.AuthorizationError{
						Error:            saml_pb.ErrorReason_ERROR_REASON_REQUEST_DENIED,
						ErrorDescription: gu.Ptr("nope"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fail callback, redirect",
			ctx:  LoginCTX,
			req: &saml_pb.CreateResponseRequest{
				SamlRequestId: func() string {
					_, authRequestID, err := Instance.CreateSAMLAuthRequest(spMiddlewareRedirect, Instance.Users[integration.UserTypeLogin].ID, acsPost, integration.RelayState(), saml.HTTPPostBinding)
					require.NoError(t, err)
					return authRequestID
				}(),
				ResponseKind: &saml_pb.CreateResponseRequest_Error{
					Error: &saml_pb.AuthorizationError{
						Error:            saml_pb.ErrorReason_ERROR_REASON_REQUEST_DENIED,
						ErrorDescription: gu.Ptr("nope"),
					},
				},
			},
			want: &saml_pb.CreateResponseResponse{
				Url:     `https:\/\/` + rootURLRedirect + `\/saml\/acs\?RelayState=(.*)&SAMLResponse=(.*)`,
				Binding: &saml_pb.CreateResponseResponse_Redirect{Redirect: &saml_pb.RedirectResponse{}},
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
			wantErr: false,
		},
		{
			name: "fail callback, no permission, error",
			ctx:  CTX,
			req: &saml_pb.CreateResponseRequest{
				SamlRequestId: func() string {
					_, authRequestID, err := Instance.CreateSAMLAuthRequest(spMiddlewareRedirect, Instance.Users[integration.UserTypeLogin].ID, acsPost, integration.RelayState(), saml.HTTPPostBinding)
					require.NoError(t, err)
					return authRequestID
				}(),
				ResponseKind: &saml_pb.CreateResponseRequest_Error{
					Error: &saml_pb.AuthorizationError{
						Error:            saml_pb.ErrorReason_ERROR_REASON_REQUEST_DENIED,
						ErrorDescription: gu.Ptr("nope"),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "callback, redirect",
			ctx:  LoginCTX,
			req: &saml_pb.CreateResponseRequest{
				SamlRequestId: func() string {
					_, authRequestID, err := Instance.CreateSAMLAuthRequest(spMiddlewareRedirect, Instance.Users[integration.UserTypeLogin].ID, acsRedirect, integration.RelayState(), saml.HTTPRedirectBinding)
					require.NoError(t, err)
					return authRequestID
				}(),
				ResponseKind: &saml_pb.CreateResponseRequest_Session{
					Session: &saml_pb.Session{
						SessionId:    sessionResp.GetSessionId(),
						SessionToken: sessionResp.GetSessionToken(),
					},
				},
			},
			want: &saml_pb.CreateResponseResponse{
				Url:     `https:\/\/` + rootURLRedirect + `\/saml\/acs\?RelayState=(.*)&SAMLResponse=(.*)&SigAlg=(.*)&Signature=(.*)`,
				Binding: &saml_pb.CreateResponseResponse_Redirect{Redirect: &saml_pb.RedirectResponse{}},
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
			wantErr: false,
		},
		{
			name: "callback, post",
			ctx:  LoginCTX,
			req: &saml_pb.CreateResponseRequest{
				SamlRequestId: func() string {
					_, authRequestID, err := Instance.CreateSAMLAuthRequest(spMiddlewarePost, Instance.Users[integration.UserTypeLogin].ID, acsPost, integration.RelayState(), saml.HTTPPostBinding)
					require.NoError(t, err)
					return authRequestID
				}(),
				ResponseKind: &saml_pb.CreateResponseRequest_Session{
					Session: &saml_pb.Session{
						SessionId:    sessionResp.GetSessionId(),
						SessionToken: sessionResp.GetSessionToken(),
					},
				},
			},
			want: &saml_pb.CreateResponseResponse{
				Url: regexp.QuoteMeta(`https://` + rootURLPost + `/saml/acs`),
				Binding: &saml_pb.CreateResponseResponse_Post{Post: &saml_pb.PostResponse{
					RelayState:   "notempty",
					SamlResponse: "notempty",
				}},
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
			wantErr: false,
		},
		{
			name: "callback, post",
			ctx:  LoginCTX,
			req: &saml_pb.CreateResponseRequest{
				SamlRequestId: func() string {
					_, authRequestID, err := Instance.CreateSAMLAuthRequest(spMiddlewarePost, Instance.Users[integration.UserTypeLogin].ID, acsPost, integration.RelayState(), saml.HTTPPostBinding)
					require.NoError(t, err)
					Instance.SuccessfulSAMLAuthRequest(LoginCTX, Instance.Users[integration.UserTypeLogin].ID, authRequestID)
					return authRequestID
				}(),
				ResponseKind: &saml_pb.CreateResponseRequest_Session{
					Session: &saml_pb.Session{
						SessionId:    sessionResp.GetSessionId(),
						SessionToken: sessionResp.GetSessionToken(),
					},
				},
			},
			wantErr: true,
		},
		{
			name: "callback, no permission, error",
			ctx:  CTX,
			req: &saml_pb.CreateResponseRequest{
				SamlRequestId: func() string {
					_, authRequestID, err := Instance.CreateSAMLAuthRequest(spMiddlewarePost, Instance.Users[integration.UserTypeLogin].ID, acsPost, integration.RelayState(), saml.HTTPPostBinding)
					require.NoError(t, err)
					return authRequestID
				}(),
				ResponseKind: &saml_pb.CreateResponseRequest_Session{
					Session: &saml_pb.Session{
						SessionId:    sessionResp.GetSessionId(),
						SessionToken: sessionResp.GetSessionToken(),
					},
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.CreateResponse(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
			if tt.want != nil {
				assert.Regexp(t, regexp.MustCompile(tt.want.Url), got.GetUrl())
				if tt.want.GetPost() != nil {
					assert.NotEmpty(t, got.GetPost().GetRelayState())
					assert.NotEmpty(t, got.GetPost().GetSamlResponse())
				}
				if tt.want.GetRedirect() != nil {
					assert.NotNil(t, got.GetRedirect())
				}
			}
		})
	}
}

func TestServer_CreateResponse_Permission(t *testing.T) {
	idpMetadata, err := Instance.GetSAMLIDPMetadata()
	require.NoError(t, err)
	acsRedirect := idpMetadata.IDPSSODescriptors[0].SingleSignOnServices[0]

	tests := []struct {
		name    string
		dep     func(ctx context.Context, t *testing.T) *saml_pb.CreateResponseRequest
		want    *saml_pb.CreateResponseResponse
		wantURL *url.URL
		wantErr bool
	}{
		{
			name: "usergrant to project and different resourceowner with different project grant",
			dep: func(ctx context.Context, t *testing.T) *saml_pb.CreateResponseRequest {
				projectID, _, sp := createSAMLApplication(ctx, t, idpMetadata, saml.HTTPRedirectBinding, true, true)
				projectID2, _, _ := createSAMLApplication(ctx, t, idpMetadata, saml.HTTPRedirectBinding, true, true)

				orgResp := Instance.CreateOrganization(ctx, integration.OrganizationName(), integration.Email())
				Instance.CreateProjectGrant(ctx, t, projectID2, orgResp.GetOrganizationId())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), integration.Email(), integration.Phone())
				createProjectUserGrant(ctx, t, Instance.DefaultOrg.GetId(), projectID, user.GetUserId())

				return createSessionAndSmlRequestForCallback(ctx, t, sp, Instance.Users[integration.UserTypeLogin].ID, acsRedirect, user.GetUserId(), saml.HTTPRedirectBinding)
			},
			wantErr: true,
		},
		{
			name: "usergrant to project and different resourceowner with project grant",
			dep: func(ctx context.Context, t *testing.T) *saml_pb.CreateResponseRequest {
				projectID, _, sp := createSAMLApplication(ctx, t, idpMetadata, saml.HTTPRedirectBinding, true, true)

				orgResp := Instance.CreateOrganization(ctx, integration.OrganizationName(), integration.Email())
				Instance.CreateProjectGrant(ctx, t, projectID, orgResp.GetOrganizationId())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), integration.Email(), integration.Phone())
				createProjectUserGrant(ctx, t, Instance.DefaultOrg.GetId(), projectID, user.GetUserId())

				return createSessionAndSmlRequestForCallback(ctx, t, sp, Instance.Users[integration.UserTypeLogin].ID, acsRedirect, user.GetUserId(), saml.HTTPRedirectBinding)
			},
			want: &saml_pb.CreateResponseResponse{
				Url:     `https:\/\/(.*)\/saml\/acs\?RelayState=(.*)&SAMLResponse=(.*)&SigAlg=(.*)&Signature=(.*)`,
				Binding: &saml_pb.CreateResponseResponse_Redirect{Redirect: &saml_pb.RedirectResponse{}},
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
			wantErr: false,
		},
		{
			name: "usergrant to project grant and different resourceowner with project grant",
			dep: func(ctx context.Context, t *testing.T) *saml_pb.CreateResponseRequest {
				projectID, _, sp := createSAMLApplication(ctx, t, idpMetadata, saml.HTTPRedirectBinding, true, true)

				orgResp := Instance.CreateOrganization(ctx, integration.OrganizationName(), integration.Email())
				Instance.CreateProjectGrant(ctx, t, projectID, orgResp.GetOrganizationId())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), integration.Email(), integration.Phone())
				createProjectGrantUserGrant(ctx, t, orgResp.GetOrganizationId(), projectID, orgResp.GetOrganizationId(), user.GetUserId())

				return createSessionAndSmlRequestForCallback(ctx, t, sp, Instance.Users[integration.UserTypeLogin].ID, acsRedirect, user.GetUserId(), saml.HTTPRedirectBinding)
			},
			want: &saml_pb.CreateResponseResponse{
				Url:     `https:\/\/(.*)\/saml\/acs\?RelayState=(.*)&SAMLResponse=(.*)&SigAlg=(.*)&Signature=(.*)`,
				Binding: &saml_pb.CreateResponseResponse_Redirect{Redirect: &saml_pb.RedirectResponse{}},
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
			wantErr: false,
		},
		{
			name: "no usergrant and different resourceowner",
			dep: func(ctx context.Context, t *testing.T) *saml_pb.CreateResponseRequest {
				_, _, sp := createSAMLApplication(ctx, t, idpMetadata, saml.HTTPRedirectBinding, true, true)

				orgResp := Instance.CreateOrganization(ctx, integration.OrganizationName(), integration.Email())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), integration.Email(), integration.Phone())
				return createSessionAndSmlRequestForCallback(ctx, t, sp, Instance.Users[integration.UserTypeLogin].ID, acsRedirect, user.GetUserId(), saml.HTTPRedirectBinding)
			},
			wantErr: true,
		},
		{
			name: "no usergrant and same resourceowner",
			dep: func(ctx context.Context, t *testing.T) *saml_pb.CreateResponseRequest {
				_, _, sp := createSAMLApplication(ctx, t, idpMetadata, saml.HTTPRedirectBinding, true, true)
				user := Instance.CreateHumanUser(ctx)

				return createSessionAndSmlRequestForCallback(ctx, t, sp, Instance.Users[integration.UserTypeLogin].ID, acsRedirect, user.GetUserId(), saml.HTTPRedirectBinding)
			},
			wantErr: true,
		},
		{
			name: "usergrant and different resourceowner",
			dep: func(ctx context.Context, t *testing.T) *saml_pb.CreateResponseRequest {
				projectID, _, sp := createSAMLApplication(ctx, t, idpMetadata, saml.HTTPRedirectBinding, true, true)

				orgResp := Instance.CreateOrganization(ctx, integration.OrganizationName(), integration.Email())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), integration.Email(), integration.Phone())
				createProjectUserGrant(ctx, t, Instance.DefaultOrg.GetId(), projectID, user.GetUserId())

				return createSessionAndSmlRequestForCallback(ctx, t, sp, Instance.Users[integration.UserTypeLogin].ID, acsRedirect, user.GetUserId(), saml.HTTPRedirectBinding)
			},
			wantErr: true,
		},
		{
			name: "usergrant and same resourceowner",
			dep: func(ctx context.Context, t *testing.T) *saml_pb.CreateResponseRequest {
				projectID, _, sp := createSAMLApplication(ctx, t, idpMetadata, saml.HTTPRedirectBinding, true, true)

				user := Instance.CreateHumanUser(ctx)
				createProjectUserGrant(ctx, t, Instance.DefaultOrg.GetId(), projectID, user.GetUserId())

				return createSessionAndSmlRequestForCallback(ctx, t, sp, Instance.Users[integration.UserTypeLogin].ID, acsRedirect, user.GetUserId(), saml.HTTPRedirectBinding)
			},
			want: &saml_pb.CreateResponseResponse{
				Url:     `https:\/\/(.*)\/saml\/acs\?RelayState=(.*)&SAMLResponse=(.*)&SigAlg=(.*)&Signature=(.*)`,
				Binding: &saml_pb.CreateResponseResponse_Redirect{Redirect: &saml_pb.RedirectResponse{}},
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
		},
		{
			name: "projectRoleCheck, usergrant and same resourceowner",
			dep: func(ctx context.Context, t *testing.T) *saml_pb.CreateResponseRequest {
				projectID, _, sp := createSAMLApplication(ctx, t, idpMetadata, saml.HTTPRedirectBinding, true, false)

				user := Instance.CreateHumanUser(ctx)
				createProjectUserGrant(ctx, t, Instance.DefaultOrg.GetId(), projectID, user.GetUserId())

				return createSessionAndSmlRequestForCallback(ctx, t, sp, Instance.Users[integration.UserTypeLogin].ID, acsRedirect, user.GetUserId(), saml.HTTPRedirectBinding)
			},
			want: &saml_pb.CreateResponseResponse{
				Url:     `https:\/\/(.*)\/saml\/acs\?RelayState=(.*)&SAMLResponse=(.*)&SigAlg=(.*)&Signature=(.*)`,
				Binding: &saml_pb.CreateResponseResponse_Redirect{Redirect: &saml_pb.RedirectResponse{}},
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
		},
		{
			name: "projectRoleCheck, no usergrant and same resourceowner",
			dep: func(ctx context.Context, t *testing.T) *saml_pb.CreateResponseRequest {
				_, _, sp := createSAMLApplication(ctx, t, idpMetadata, saml.HTTPRedirectBinding, true, false)
				user := Instance.CreateHumanUser(ctx)

				return createSessionAndSmlRequestForCallback(ctx, t, sp, Instance.Users[integration.UserTypeLogin].ID, acsRedirect, user.GetUserId(), saml.HTTPRedirectBinding)
			},
			wantErr: true,
		},
		{
			name: "projectRoleCheck, usergrant and different resourceowner",
			dep: func(ctx context.Context, t *testing.T) *saml_pb.CreateResponseRequest {
				projectID, _, sp := createSAMLApplication(ctx, t, idpMetadata, saml.HTTPRedirectBinding, true, false)
				orgResp := Instance.CreateOrganization(ctx, integration.OrganizationName(), integration.Email())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), integration.Email(), integration.Phone())
				createProjectUserGrant(ctx, t, Instance.DefaultOrg.GetId(), projectID, user.GetUserId())

				return createSessionAndSmlRequestForCallback(ctx, t, sp, Instance.Users[integration.UserTypeLogin].ID, acsRedirect, user.GetUserId(), saml.HTTPRedirectBinding)
			},
			want: &saml_pb.CreateResponseResponse{
				Url:     `https:\/\/(.*)\/saml\/acs\?RelayState=(.*)&SAMLResponse=(.*)&SigAlg=(.*)&Signature=(.*)`,
				Binding: &saml_pb.CreateResponseResponse_Redirect{Redirect: &saml_pb.RedirectResponse{}},
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
		},
		{
			name: "projectRoleCheck, no usergrant and different resourceowner",
			dep: func(ctx context.Context, t *testing.T) *saml_pb.CreateResponseRequest {
				_, _, sp := createSAMLApplication(ctx, t, idpMetadata, saml.HTTPRedirectBinding, true, false)
				orgResp := Instance.CreateOrganization(ctx, integration.OrganizationName(), integration.Email())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), integration.Email(), integration.Phone())

				return createSessionAndSmlRequestForCallback(ctx, t, sp, Instance.Users[integration.UserTypeLogin].ID, acsRedirect, user.GetUserId(), saml.HTTPRedirectBinding)
			},
			wantErr: true,
		},
		{
			name: "projectRoleCheck, usergrant on project grant and different resourceowner",
			dep: func(ctx context.Context, t *testing.T) *saml_pb.CreateResponseRequest {
				projectID, _, sp := createSAMLApplication(ctx, t, idpMetadata, saml.HTTPRedirectBinding, true, false)

				orgResp := Instance.CreateOrganization(ctx, integration.OrganizationName(), integration.Email())
				Instance.CreateProjectGrant(ctx, t, projectID, orgResp.GetOrganizationId())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), integration.Email(), integration.Phone())
				createProjectGrantUserGrant(ctx, t, orgResp.GetOrganizationId(), projectID, orgResp.GetOrganizationId(), user.GetUserId())

				return createSessionAndSmlRequestForCallback(ctx, t, sp, Instance.Users[integration.UserTypeLogin].ID, acsRedirect, user.GetUserId(), saml.HTTPRedirectBinding)
			},
			want: &saml_pb.CreateResponseResponse{
				Url:     `https:\/\/(.*)\/saml\/acs\?RelayState=(.*)&SAMLResponse=(.*)&SigAlg=(.*)&Signature=(.*)`,
				Binding: &saml_pb.CreateResponseResponse_Redirect{Redirect: &saml_pb.RedirectResponse{}},
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
		},
		{
			name: "projectRoleCheck, no usergrant on project grant and different resourceowner",
			dep: func(ctx context.Context, t *testing.T) *saml_pb.CreateResponseRequest {
				projectID, _, sp := createSAMLApplication(ctx, t, idpMetadata, saml.HTTPRedirectBinding, true, false)

				orgResp := Instance.CreateOrganization(ctx, integration.OrganizationName(), integration.Email())
				Instance.CreateProjectGrant(ctx, t, projectID, orgResp.GetOrganizationId())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), integration.Email(), integration.Phone())

				return createSessionAndSmlRequestForCallback(ctx, t, sp, Instance.Users[integration.UserTypeLogin].ID, acsRedirect, user.GetUserId(), saml.HTTPRedirectBinding)
			},
			wantErr: true,
		},
		{
			name: "hasProjectCheck, same resourceowner",
			dep: func(ctx context.Context, t *testing.T) *saml_pb.CreateResponseRequest {
				_, _, sp := createSAMLApplication(ctx, t, idpMetadata, saml.HTTPRedirectBinding, false, true)
				user := Instance.CreateHumanUser(ctx)

				return createSessionAndSmlRequestForCallback(ctx, t, sp, Instance.Users[integration.UserTypeLogin].ID, acsRedirect, user.GetUserId(), saml.HTTPRedirectBinding)
			},
			want: &saml_pb.CreateResponseResponse{
				Url:     `https:\/\/(.*)\/saml\/acs\?RelayState=(.*)&SAMLResponse=(.*)&SigAlg=(.*)&Signature=(.*)`,
				Binding: &saml_pb.CreateResponseResponse_Redirect{Redirect: &saml_pb.RedirectResponse{}},
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
			},
		},
		{
			name: "hasProjectCheck, different resourceowner",
			dep: func(ctx context.Context, t *testing.T) *saml_pb.CreateResponseRequest {
				_, _, sp := createSAMLApplication(ctx, t, idpMetadata, saml.HTTPRedirectBinding, false, true)
				orgResp := Instance.CreateOrganization(ctx, integration.OrganizationName(), integration.Email())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), integration.Email(), integration.Phone())

				return createSessionAndSmlRequestForCallback(ctx, t, sp, Instance.Users[integration.UserTypeLogin].ID, acsRedirect, user.GetUserId(), saml.HTTPRedirectBinding)
			},
			wantErr: true,
		},
		{
			name: "hasProjectCheck, different resourceowner with project grant",
			dep: func(ctx context.Context, t *testing.T) *saml_pb.CreateResponseRequest {
				projectID, _, sp := createSAMLApplication(ctx, t, idpMetadata, saml.HTTPRedirectBinding, false, true)
				orgResp := Instance.CreateOrganization(ctx, integration.OrganizationName(), integration.Email())
				Instance.CreateProjectGrant(ctx, t, projectID, orgResp.GetOrganizationId())
				user := Instance.CreateHumanUserVerified(ctx, orgResp.GetOrganizationId(), integration.Email(), integration.Phone())

				return createSessionAndSmlRequestForCallback(ctx, t, sp, Instance.Users[integration.UserTypeLogin].ID, acsRedirect, user.GetUserId(), saml.HTTPRedirectBinding)
			},
			want: &saml_pb.CreateResponseResponse{
				Url:     `https:\/\/(.*)\/saml\/acs\?RelayState=(.*)&SAMLResponse=(.*)&SigAlg=(.*)&Signature=(.*)`,
				Binding: &saml_pb.CreateResponseResponse_Redirect{Redirect: &saml_pb.RedirectResponse{}},
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

			got, err := Client.CreateResponse(LoginCTX, req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
			if tt.want != nil {
				assert.Regexp(t, regexp.MustCompile(tt.want.Url), got.GetUrl())
				if tt.want.GetPost() != nil {
					assert.NotEmpty(t, got.GetPost().GetRelayState())
					assert.NotEmpty(t, got.GetPost().GetSamlResponse())
				}
				if tt.want.GetRedirect() != nil {
					assert.NotNil(t, got.GetRedirect())
				}
			}
		})
	}
}

func createSession(ctx context.Context, t *testing.T, userID string) *session.CreateSessionResponse {
	sessionResp, err := Instance.Client.SessionV2.CreateSession(ctx, &session.CreateSessionRequest{
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

func createSessionAndSmlRequestForCallback(ctx context.Context, t *testing.T, sp *samlsp.Middleware, loginClient string, acsRedirect saml.Endpoint, userID, binding string) *saml_pb.CreateResponseRequest {
	_, authRequestID, err := Instance.CreateSAMLAuthRequest(sp, loginClient, acsRedirect, integration.RelayState(), binding)
	require.NoError(t, err)
	sessionResp := createSession(ctx, t, userID)
	return &saml_pb.CreateResponseRequest{
		SamlRequestId: authRequestID,
		ResponseKind: &saml_pb.CreateResponseRequest_Session{
			Session: &saml_pb.Session{
				SessionId:    sessionResp.GetSessionId(),
				SessionToken: sessionResp.GetSessionToken(),
			},
		},
	}
}

func createSAMLSP(t *testing.T, idpMetadata *saml.EntityDescriptor, binding string) (string, *samlsp.Middleware) {
	rootURL := "example." + integration.DomainName()
	spMiddleware, err := integration.CreateSAMLSP("https://"+rootURL, idpMetadata, binding)
	require.NoError(t, err)
	return rootURL, spMiddleware
}

func createSAMLApplication(ctx context.Context, t *testing.T, idpMetadata *saml.EntityDescriptor, binding string, projectRoleCheck, hasProjectCheck bool) (string, string, *samlsp.Middleware) {
	project := Instance.CreateProject(ctx, t, "", integration.ProjectName(), projectRoleCheck, hasProjectCheck)
	rootURL, sp := createSAMLSP(t, idpMetadata, binding)
	_, err := Instance.CreateSAMLClient(ctx, project.GetId(), sp)
	require.NoError(t, err)
	return project.GetId(), rootURL, sp
}

func createProjectUserGrant(ctx context.Context, t *testing.T, orgID, projectID, userID string) {
	resp := Instance.CreateProjectUserGrant(t, ctx, orgID, projectID, userID)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(collect *assert.CollectT) {
		_, err := Instance.Client.Mgmt.GetUserGrantByID(integration.SetOrgID(ctx, orgID), &mgmt.GetUserGrantByIDRequest{
			UserId:  userID,
			GrantId: resp.GetUserGrantId(),
		})
		assert.NoError(collect, err)
	}, retryDuration, tick)
}

func createProjectGrantUserGrant(ctx context.Context, t *testing.T, orgID, projectID, projectGrantID, userID string) {
	resp := Instance.CreateProjectGrantUserGrant(ctx, orgID, projectID, projectGrantID, userID)

	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(collect *assert.CollectT) {
		_, err := Instance.Client.Mgmt.GetUserGrantByID(integration.SetOrgID(ctx, orgID), &mgmt.GetUserGrantByIDRequest{
			UserId:  userID,
			GrantId: resp.GetUserGrantId(),
		})
		assert.NoError(collect, err)
	}, retryDuration, tick)
}
