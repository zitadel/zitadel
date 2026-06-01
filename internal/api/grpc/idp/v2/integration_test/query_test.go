//go:build integration

package idp_test

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/idp/v2"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
)

type idpAttr struct {
	ID      string
	Name    string
	Details *object.Details
}

func TestServer_GetIDPByID(t *testing.T) {
	type args struct {
		ctx context.Context
		req *idp.GetIDPByIDRequest
		dep func(ctx context.Context, request *idp.GetIDPByIDRequest) *idpAttr
	}
	tests := []struct {
		name    string
		args    args
		want    *idp.GetIDPByIDResponse
		wantErr bool
	}{
		{
			name: "idp by ID, no id provided",
			args: args{
				IamCTX,
				&idp.GetIDPByIDRequest{
					Id: "",
				},
				func(ctx context.Context, request *idp.GetIDPByIDRequest) *idpAttr {
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "idp by ID, not found",
			args: args{
				IamCTX,
				&idp.GetIDPByIDRequest{
					Id: "unknown",
				},
				func(ctx context.Context, request *idp.GetIDPByIDRequest) *idpAttr {
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "idp by ID, instance, ok",
			args: args{
				IamCTX,
				&idp.GetIDPByIDRequest{},
				func(ctx context.Context, request *idp.GetIDPByIDRequest) *idpAttr {
					name := integration.IDPName()
					resp := Instance.AddGenericOAuthProvider(ctx, name)
					request.Id = resp.Id
					return &idpAttr{
						resp.GetId(),
						name,
						&object.Details{
							Sequence:      resp.Details.Sequence,
							CreationDate:  resp.Details.CreationDate,
							ChangeDate:    resp.Details.ChangeDate,
							ResourceOwner: resp.Details.ResourceOwner,
						}}
				},
			},
			want: &idp.GetIDPByIDResponse{
				Idp: &idp.IDP{
					Details: &object.Details{
						ChangeDate: timestamppb.Now(),
					},
					State: idp.IDPState_IDP_STATE_ACTIVE,
					Type:  idp.IDPType_IDP_TYPE_OAUTH,
					Config: &idp.IDPConfig{
						Config: &idp.IDPConfig_Oauth{
							Oauth: &idp.OAuthConfig{
								ClientId:              "clientID",
								AuthorizationEndpoint: "https://example.com/oauth/v2/authorize",
								TokenEndpoint:         "https://example.com/oauth/v2/token",
								UserEndpoint:          "https://api.example.com/user",
								Scopes:                []string{"openid", "profile", "email"},
								IdAttribute:           "id",
							},
						},
						Options: &idp.Options{
							IsLinkingAllowed:  true,
							IsCreationAllowed: true,
							IsAutoCreation:    true,
							IsAutoUpdate:      true,
							AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
						},
					},
				},
			},
		},
		{
			name: "idp by ID, instance, no permission",
			args: args{
				UserCTX,
				&idp.GetIDPByIDRequest{},
				func(ctx context.Context, request *idp.GetIDPByIDRequest) *idpAttr {
					name := integration.IDPName()
					resp := Instance.AddGenericOAuthProvider(IamCTX, name)
					request.Id = resp.Id
					return &idpAttr{
						resp.GetId(),
						name,
						&object.Details{
							Sequence:      resp.Details.Sequence,
							CreationDate:  resp.Details.CreationDate,
							ChangeDate:    resp.Details.ChangeDate,
							ResourceOwner: resp.Details.ResourceOwner,
						}}
				},
			},
			wantErr: true,
		},
		{
			name: "idp by ID, org, ok",
			args: args{
				CTX,
				&idp.GetIDPByIDRequest{},
				func(ctx context.Context, request *idp.GetIDPByIDRequest) *idpAttr {
					name := integration.IDPName()
					resp := Instance.AddOrgGenericOAuthProvider(ctx, name)
					request.Id = resp.Id
					return &idpAttr{
						resp.GetId(),
						name,
						&object.Details{
							Sequence:      resp.Details.Sequence,
							CreationDate:  resp.Details.CreationDate,
							ChangeDate:    resp.Details.ChangeDate,
							ResourceOwner: resp.Details.ResourceOwner,
						}}
				},
			},
			want: &idp.GetIDPByIDResponse{
				Idp: &idp.IDP{
					Details: &object.Details{
						ChangeDate: timestamppb.Now(),
					},
					State: idp.IDPState_IDP_STATE_ACTIVE,
					Type:  idp.IDPType_IDP_TYPE_OAUTH,
					Config: &idp.IDPConfig{
						Config: &idp.IDPConfig_Oauth{
							Oauth: &idp.OAuthConfig{
								ClientId:              "clientID",
								AuthorizationEndpoint: "https://example.com/oauth/v2/authorize",
								TokenEndpoint:         "https://example.com/oauth/v2/token",
								UserEndpoint:          "https://api.example.com/user",
								Scopes:                []string{"openid", "profile", "email"},
								IdAttribute:           "id",
							},
						},
						Options: &idp.Options{
							IsLinkingAllowed:  true,
							IsCreationAllowed: true,
							IsAutoCreation:    true,
							IsAutoUpdate:      true,
							AutoLinking:       idp.AutoLinkingOption_AUTO_LINKING_OPTION_USERNAME,
						},
					},
				},
			},
		},
		{
			name: "idp by ID, org, no permission",
			args: args{
				UserCTX,
				&idp.GetIDPByIDRequest{},
				func(ctx context.Context, request *idp.GetIDPByIDRequest) *idpAttr {
					name := integration.IDPName()
					resp := Instance.AddOrgGenericOAuthProvider(CTX, name)
					request.Id = resp.Id
					return &idpAttr{
						resp.GetId(),
						name,
						&object.Details{
							Sequence:      resp.Details.Sequence,
							CreationDate:  resp.Details.CreationDate,
							ChangeDate:    resp.Details.ChangeDate,
							ResourceOwner: resp.Details.ResourceOwner,
						}}
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			idpAttr := tt.args.dep(tt.args.ctx, tt.args.req)
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(CTX, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.GetIDPByID(tt.args.ctx, tt.args.req)
				if tt.wantErr {
					require.Error(ttt, err)
					return
				}
				require.NoError(ttt, err)

				// set provided info from creation
				tt.want.Idp.Details = idpAttr.Details
				tt.want.Idp.Name = idpAttr.Name
				tt.want.Idp.Id = idpAttr.ID

				// first check for details, mgmt and admin api don't fill the details correctly
				integration.AssertDetails(t, tt.want.Idp, got.Idp)
				// then set details
				tt.want.Idp.Details = got.Idp.Details
				// to check the rest of the content
				assert.Equal(ttt, tt.want.Idp, got.Idp)
			}, retryDuration, tick)
		})
	}
}
