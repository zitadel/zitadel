//go:build integration

package org_test

import (
	"context"
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/zitadel/zitadel/internal/integration"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2alpha"
	org "github.com/zitadel/zitadel/pkg/grpc/organisation/v2beta"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2alpha"
)

var (
	CTX    context.Context
	Tester *integration.Tester
	Client org.OrganisationServiceClient
	User   *user.AddHumanUserResponse
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, errCtx, cancel := integration.Contexts(5 * time.Minute)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()
		Client = Tester.Client.Orgv2

		CTX, _ = Tester.WithAuthorization(ctx, integration.IAMOwner), errCtx
		User = Tester.CreateHumanUser(CTX)
		return m.Run()
	}())
}

func TestServer_AddOrganisation(t *testing.T) {
	tests := []struct {
		name    string
		ctx     context.Context
		req     *org.AddOrganisationRequest
		want    *org.AddOrganisationResponse
		wantErr bool
	}{
		{
			name: "missing permission",
			ctx:  Tester.WithAuthorization(context.Background(), integration.OrgOwner),
			req: &org.AddOrganisationRequest{
				Name:   "name",
				Admins: nil,
			},
			wantErr: true,
		},
		{
			name: "empty name",
			ctx:  CTX,
			req: &org.AddOrganisationRequest{
				Name:   "",
				Admins: nil,
			},
			wantErr: true,
		},
		{
			name: "invalid admin type",
			ctx:  CTX,
			req: &org.AddOrganisationRequest{
				Name: fmt.Sprintf("%d", time.Now().UnixNano()),
				Admins: []*org.AddOrganisationRequest_Admin{
					{},
				},
			},
			wantErr: true,
		},
		{
			name: "admin with init",
			ctx:  CTX,
			req: &org.AddOrganisationRequest{
				Name: fmt.Sprintf("%d", time.Now().UnixNano()),
				Admins: []*org.AddOrganisationRequest_Admin{
					{
						UserType: &org.AddOrganisationRequest_Admin_Human{
							Human: &user.AddHumanUserRequest{
								Profile: &user.SetHumanProfile{
									FirstName: "firstname",
									LastName:  "lastname",
								},
								Email: &user.SetHumanEmail{
									Email: fmt.Sprintf("%d@mouse.com", time.Now().UnixNano()),
									Verification: &user.SetHumanEmail_ReturnCode{
										ReturnCode: &user.ReturnEmailVerificationCode{},
									},
								},
							},
						},
					},
				},
			},
			want: &org.AddOrganisationResponse{
				Details: &object.Details{
					Sequence:      0,
					ChangeDate:    nil,
					ResourceOwner: "orgID",
				},
				OrganisationId: "",
				CreatedAdmins: []*org.AddOrganisationResponse_CreatedAdmin{
					{
						UserId:     "userID",
						EmailCode:  gu.Ptr("code"),
						PhoneCode:  nil,
						Pat:        nil,
						MachineKey: nil,
					},
				},
			},
		},
		{
			name: "existing user, new human and machine with pat and key",
			ctx:  CTX,
			req: &org.AddOrganisationRequest{
				Name: fmt.Sprintf("%d", time.Now().UnixNano()),
				Admins: []*org.AddOrganisationRequest_Admin{
					{
						UserType: &org.AddOrganisationRequest_Admin_UserId{UserId: User.GetUserId()},
					},
					{
						UserType: &org.AddOrganisationRequest_Admin_Human{
							Human: &user.AddHumanUserRequest{
								Profile: &user.SetHumanProfile{
									FirstName: "firstname",
									LastName:  "lastname",
								},
								Email: &user.SetHumanEmail{
									Email: fmt.Sprintf("%d@mouse.com", time.Now().UnixNano()),
									Verification: &user.SetHumanEmail_ReturnCode{
										ReturnCode: &user.ReturnEmailVerificationCode{},
									},
								},
							},
						},
					},
					{
						UserType: &org.AddOrganisationRequest_Admin_Machine{
							Machine: &org.AddMachineUserRequest{
								Username:   fmt.Sprintf("%d", time.Now().UnixNano()),
								Name:       "name",
								Pat:        true,
								MachineKey: true,
							},
						},
					},
				},
			},
			want: &org.AddOrganisationResponse{
				CreatedAdmins: []*org.AddOrganisationResponse_CreatedAdmin{
					{
						UserId:     integration.NotEmpty,
						Pat:        nil,
						MachineKey: nil,
					},
					{
						UserId:     integration.NotEmpty,
						Pat:        gu.Ptr(integration.NotEmpty),
						MachineKey: []byte(integration.NotEmpty),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.AddOrganisation(tt.ctx, tt.req)
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

			// organisation id must be the same as the resourceOwner
			assert.Equal(t, got.GetDetails().GetResourceOwner(), got.GetOrganisationId())

			// check the admins
			require.Len(t, got.GetCreatedAdmins(), len(tt.want.GetCreatedAdmins()))
			for i, admin := range tt.want.GetCreatedAdmins() {
				gotAdmin := got.GetCreatedAdmins()[i]
				assertCreatedAdmin(t, admin, gotAdmin)
			}
		})
	}
}

func assertCreatedAdmin(t *testing.T, expected, got *org.AddOrganisationResponse_CreatedAdmin) {
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
	if expected.GetPat() != "" {
		assert.NotEmpty(t, got.GetPat())
	} else {
		assert.Empty(t, got.GetPat())
	}
	if expected.GetMachineKey() != nil {
		assert.NotEmpty(t, got.GetMachineKey())
	} else {
		assert.Empty(t, got.GetMachineKey())
	}
}
