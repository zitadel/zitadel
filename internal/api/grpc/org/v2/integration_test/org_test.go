//go:build integration

package org_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/org/v2"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

var (
	CTX, OwnerCTX, UserCTX context.Context
	Instance               *integration.Instance
	Client                 org.OrganizationServiceClient
	User                   *user.AddHumanUserResponse
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)
		Client = Instance.Client.OrgV2

		CTX = Instance.WithAuthorization(ctx, integration.UserTypeIAMOwner)
		OwnerCTX = Instance.WithAuthorization(ctx, integration.UserTypeOrgOwner)
		UserCTX = Instance.WithAuthorization(ctx, integration.UserTypeNoPermission)
		User = Instance.CreateHumanUser(CTX)
		return m.Run()
	}())
}

func TestServer_AddOrganization(t *testing.T) {
	idpResp := Instance.AddGenericOAuthProvider(CTX, Instance.DefaultOrg.Id)
	userId := "userID"

	tests := []struct {
		name    string
		ctx     context.Context
		req     *org.AddOrganizationRequest
		want    *org.AddOrganizationResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Instance.WithAuthorization(context.Background(), integration.UserTypeOrgOwner),
			req: &org.AddOrganizationRequest{
				Name:   "name",
				Admins: nil,
			},
			wantErr: true,
		},
		{
			name: "empty name",
			ctx:  CTX,
			req: &org.AddOrganizationRequest{
				Name:   "",
				Admins: nil,
			},
			wantErr: true,
		},
		{
			name: "invalid admin type",
			ctx:  CTX,
			req: &org.AddOrganizationRequest{
				Name: integration.OrganizationName(),
				Admins: []*org.AddOrganizationRequest_Admin{
					{},
				},
			},
			wantErr: true,
		},
		{
			name: "no admin, custom org ID",
			ctx:  CTX,
			req: &org.AddOrganizationRequest{
				Name:  integration.OrganizationName(),
				OrgId: gu.Ptr("custom-org-ID"),
			},
			want: &org.AddOrganizationResponse{
				OrganizationId: "custom-org-ID",
				CreatedAdmins:  []*org.AddOrganizationResponse_CreatedAdmin{},
			},
		},
		{
			name: "admin with init with userID passed for Human admin",
			ctx:  CTX,
			req: &org.AddOrganizationRequest{
				Name: integration.OrganizationName(),
				Admins: []*org.AddOrganizationRequest_Admin{
					{
						UserType: &org.AddOrganizationRequest_Admin_Human{
							Human: &user.AddHumanUserRequest{
								UserId: &userId,
								Profile: &user.SetHumanProfile{
									GivenName:  "firstname",
									FamilyName: "lastname",
								},
								Email: &user.SetHumanEmail{
									Email: integration.Email(),
									Verification: &user.SetHumanEmail_ReturnCode{
										ReturnCode: &user.ReturnEmailVerificationCode{},
									},
								},
							},
						},
					},
				},
			},
			want: &org.AddOrganizationResponse{
				OrganizationId: integration.NotEmpty,
				CreatedAdmins: []*org.AddOrganizationResponse_CreatedAdmin{
					{
						UserId:    userId,
						EmailCode: gu.Ptr(integration.NotEmpty),
						PhoneCode: nil,
					},
				},
			},
		},
		{
			name: "existing user and new human with idp",
			ctx:  CTX,
			req: &org.AddOrganizationRequest{
				Name: integration.OrganizationName(),
				Admins: []*org.AddOrganizationRequest_Admin{
					{
						UserType: &org.AddOrganizationRequest_Admin_UserId{UserId: User.GetUserId()},
					},
					{
						UserType: &org.AddOrganizationRequest_Admin_Human{
							Human: &user.AddHumanUserRequest{
								Profile: &user.SetHumanProfile{
									GivenName:  "firstname",
									FamilyName: "lastname",
								},
								Email: &user.SetHumanEmail{
									Email: integration.Email(),
									Verification: &user.SetHumanEmail_IsVerified{
										IsVerified: true,
									},
								},
								IdpLinks: []*user.IDPLink{
									{
										IdpId:    idpResp.Id,
										UserId:   userId,
										UserName: "username",
									},
								},
							},
						},
					},
				},
			},
			want: &org.AddOrganizationResponse{
				CreatedAdmins: []*org.AddOrganizationResponse_CreatedAdmin{
					// a single admin is expected, because the first provided already exists
					{
						UserId: integration.NotEmpty,
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.AddOrganization(tt.ctx, tt.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			// check details
			assert.NotZero(t, got.GetDetails().GetSequence())
			gotCD := got.GetDetails().GetChangeDate().AsTime()
			now := time.Now()
			assert.WithinRange(t, gotCD, now.Add(-time.Minute), now.Add(time.Minute))
			assert.NotEmpty(t, got.GetDetails().GetResourceOwner())

			// organization id must be the same as the resourceOwner
			assert.Equal(t, got.GetDetails().GetResourceOwner(), got.GetOrganizationId())

			// check the admins
			require.Len(t, got.GetCreatedAdmins(), len(tt.want.GetCreatedAdmins()))
			for i, admin := range tt.want.GetCreatedAdmins() {
				gotAdmin := got.GetCreatedAdmins()[i]
				assertCreatedAdmin(t, admin, gotAdmin)
			}
		})
	}
}

func assertCreatedAdmin(t *testing.T, expected, got *org.AddOrganizationResponse_CreatedAdmin) {
	if expected.GetUserId() != "" {
		assert.NotEmpty(t, got.GetUserId())
	} else {
		assert.Empty(t, got.GetUserId())
	}
	if expected.GetEmailCode() != "" {
		assert.NotEmpty(t, got.GetEmailCode())
	} else {
		assert.Empty(t, got.GetEmailCode())
	}
	if expected.GetPhoneCode() != "" {
		assert.NotEmpty(t, got.GetPhoneCode())
	} else {
		assert.Empty(t, got.GetPhoneCode())
	}
}
