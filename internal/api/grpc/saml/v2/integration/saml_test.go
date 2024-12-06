//go:build integration

package saml_test

import (
	"context"
	"net/url"
	"os"
	"regexp"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	oidc_pb "github.com/zitadel/zitadel/pkg/grpc/oidc/v2"
	saml_pb "github.com/zitadel/zitadel/pkg/grpc/saml/v2"
	"github.com/zitadel/zitadel/pkg/grpc/session/v2"
)

var (
	CTX      context.Context
	Instance *integration.Instance
	Client   saml_pb.SAMLServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)
		Client = Instance.Client.SAMLv2

		CTX = Instance.WithAuthorization(ctx, integration.UserTypeOrgOwner)
		return m.Run()
	}())
}

func TestServer_GetAuthRequest(t *testing.T) {
	rootURL := "https://sp.example.com"
	idpMetadata, err := Instance.GetSAMLIDPMetadata()
	require.NoError(t, err)
	spMiddleware, err := integration.CreateSAMLSP(rootURL, idpMetadata)
	require.NoError(t, err)

	acsRedirect := idpMetadata.IDPSSODescriptors[0].SingleSignOnServices[0]
	acsPost := idpMetadata.IDPSSODescriptors[0].SingleSignOnServices[1]

	project, err := Instance.CreateProject(CTX)
	require.NoError(t, err)
	_, err = Instance.CreateSAMLClient(CTX, project.GetId(), spMiddleware)
	require.NoError(t, err)

	now := time.Now()

	tests := []struct {
		name    string
		dep     func() (string, error)
		want    *oidc_pb.GetAuthRequestResponse
		wantErr bool
	}{
		{
			name: "Not found",
			dep: func() (string, error) {
				return "123", nil
			},
			wantErr: true,
		},
		{
			name: "success, redirect binding",
			dep: func() (string, error) {
				return Instance.CreateSAMLAuthRequest(spMiddleware, Instance.Users[integration.UserTypeOrgOwner].ID, acsRedirect, gofakeit.BitcoinAddress())
			},
		},
		{
			name: "success, post binding",
			dep: func() (string, error) {
				return Instance.CreateSAMLAuthRequest(spMiddleware, Instance.Users[integration.UserTypeOrgOwner].ID, acsPost, gofakeit.BitcoinAddress())
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			authRequestID, err := tt.dep()
			require.NoError(t, err)

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.GetAuthRequest(CTX, &saml_pb.GetSAMLRequestRequest{
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
				assert.WithinRange(ttt, authRequest.GetCreationDate().AsTime(), now.Add(-time.Second), now.Add(time.Second))
			}, retryDuration, tick, "timeout waiting for expected saml request result")
		})
	}
}

func TestServer_CreateCallback(t *testing.T) {
	rootURL := "https://sp.example.com"
	idpMetadata, err := Instance.GetSAMLIDPMetadata()
	require.NoError(t, err)
	spMiddleware, err := integration.CreateSAMLSP(rootURL, idpMetadata)
	require.NoError(t, err)

	acsRedirect := idpMetadata.IDPSSODescriptors[0].SingleSignOnServices[0]
	//acsPost := idpMetadata.IDPSSODescriptors[0].SingleSignOnServices[1]

	project, err := Instance.CreateProject(CTX)
	require.NoError(t, err)
	_, err = Instance.CreateSAMLClient(CTX, project.GetId(), spMiddleware)
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
		req       *saml_pb.CreateCallbackRequest
		AuthError string
		want      *saml_pb.CreateCallbackResponse
		wantURL   *url.URL
		wantErr   bool
	}{
		{
			name: "Not found",
			req: &saml_pb.CreateCallbackRequest{
				SamlRequestId: "123",
				CallbackKind: &saml_pb.CreateCallbackRequest_Session{
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
			req: &saml_pb.CreateCallbackRequest{
				SamlRequestId: func() string {
					authRequestID, err := Instance.CreateSAMLAuthRequest(spMiddleware, Instance.Users[integration.UserTypeOrgOwner].ID, acsRedirect, gofakeit.BitcoinAddress())
					require.NoError(t, err)
					return authRequestID
				}(),
				CallbackKind: &saml_pb.CreateCallbackRequest_Session{
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
			req: &saml_pb.CreateCallbackRequest{
				SamlRequestId: func() string {
					authRequestID, err := Instance.CreateSAMLAuthRequest(spMiddleware, Instance.Users[integration.UserTypeOrgOwner].ID, acsRedirect, gofakeit.BitcoinAddress())
					require.NoError(t, err)
					return authRequestID
				}(),
				CallbackKind: &saml_pb.CreateCallbackRequest_Session{
					Session: &saml_pb.Session{
						SessionId:    sessionResp.GetSessionId(),
						SessionToken: "bar",
					},
				},
			},
			wantErr: true,
		},
		{
			name: "fail callback",
			req: &saml_pb.CreateCallbackRequest{
				SamlRequestId: func() string {
					authRequestID, err := Instance.CreateSAMLAuthRequest(spMiddleware, Instance.Users[integration.UserTypeOrgOwner].ID, acsRedirect, gofakeit.BitcoinAddress())
					require.NoError(t, err)
					return authRequestID
				}(),
				CallbackKind: &saml_pb.CreateCallbackRequest_Error{
					Error: &saml_pb.AuthorizationError{
						Error:            saml_pb.ErrorReason_ERROR_REASON_REQUEST_DENIED,
						ErrorDescription: gu.Ptr("nope"),
						ErrorUri:         gu.Ptr("https://example.com/docs"),
					},
				},
			},
			want: &saml_pb.CreateCallbackResponse{
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
			req: &saml_pb.CreateCallbackRequest{
				SamlRequestId: func() string {
					authRequestID, err := Instance.CreateSAMLAuthRequest(spMiddleware, Instance.Users[integration.UserTypeOrgOwner].ID, acsRedirect, gofakeit.BitcoinAddress())
					require.NoError(t, err)
					return authRequestID
				}(),
				CallbackKind: &saml_pb.CreateCallbackRequest_Session{
					Session: &saml_pb.Session{
						SessionId:    sessionResp.GetSessionId(),
						SessionToken: sessionResp.GetSessionToken(),
					},
				},
			},
			want: &saml_pb.CreateCallbackResponse{
				CallbackUrl: `oidcintegrationtest:\/\/callback\?code=(.*)&state=state`,
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
