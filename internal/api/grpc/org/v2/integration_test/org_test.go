//go:build integration

package org_test

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/brianvoe/gofakeit/v6"
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
				Name: gofakeit.AppName(),
				Admins: []*org.AddOrganizationRequest_Admin{
					{},
				},
			},
			wantErr: true,
		},
		{
			name: "admin with init with userID passed for Human admin",
			ctx:  CTX,
			req: &org.AddOrganizationRequest{
				Name: gofakeit.AppName(),
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
									Email: gofakeit.Email(),
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
				OrganizationAdmins: []*org.OrganizationAdmin{
					{
						OrganizationAdmin: &org.OrganizationAdmin_CreatedAdmin{
							CreatedAdmin: &org.CreatedAdmin{
								UserId:    userId,
								EmailCode: gu.Ptr(integration.NotEmpty),
								PhoneCode: nil,
							},
						},
					},
				},
			},
		},
		{
			name: "existing user and new human with idp",
			ctx:  CTX,
			req: &org.AddOrganizationRequest{
				Name: gofakeit.AppName(),
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
									Email: gofakeit.Email(),
									Verification: &user.SetHumanEmail_IsVerified{
										IsVerified: true,
									},
								},
								IdpLinks: []*user.IDPLink{
									{
										IdpId:    idpResp.Id,
										UserId:   "userID",
										UserName: "username",
									},
								},
							},
						},
					},
				},
			},
			want: &org.AddOrganizationResponse{
				OrganizationAdmins: []*org.OrganizationAdmin{
					{
						OrganizationAdmin: &org.OrganizationAdmin_AssignedAdmin{
							AssignedAdmin: &org.AssignedAdmin{
								UserId: User.GetUserId(),
							},
						},
					},
					{
						OrganizationAdmin: &org.OrganizationAdmin_CreatedAdmin{
							CreatedAdmin: &org.CreatedAdmin{
								UserId: integration.NotEmpty,
							},
						},
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
			require.Len(t, got.GetOrganizationAdmins(), len(tt.want.GetOrganizationAdmins()))
			for i, admin := range tt.want.GetOrganizationAdmins() {
				gotAdmin := got.GetOrganizationAdmins()[i].OrganizationAdmin
				switch admin := admin.OrganizationAdmin.(type) {
				case *org.OrganizationAdmin_CreatedAdmin:
					assertCreatedAdmin(t, admin.CreatedAdmin, gotAdmin.(*org.OrganizationAdmin_CreatedAdmin).CreatedAdmin)
				case *org.OrganizationAdmin_AssignedAdmin:
					assert.Equal(t, admin.AssignedAdmin.GetUserId(), gotAdmin.(*org.OrganizationAdmin_AssignedAdmin).AssignedAdmin.GetUserId())
				}
			}
		})
	}
}

func assertCreatedAdmin(t *testing.T, expected, got *org.CreatedAdmin) {
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
