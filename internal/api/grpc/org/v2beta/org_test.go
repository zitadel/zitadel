package org

import (
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
	v2beta_org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func Test_createOrganizationRequestToCommand(t *testing.T) {
	type args struct {
		request *v2beta_org.CreateOrganizationRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *command.OrgSetup
		wantErr error
	}{
		{
			name: "nil user",
			args: args{
				request: &v2beta_org.CreateOrganizationRequest{
					Name: "name",
					Admins: []*v2beta_org.CreateOrganizationRequest_Admin{
						{},
					},
				},
			},
			wantErr: zerrors.ThrowUnimplementedf(nil, "ORGv2-SL2r8", "userType oneOf %T in method AddOrganization not implemented", nil),
		},
		{
			name: "user ID",
			args: args{
				request: &v2beta_org.CreateOrganizationRequest{
					Name: "name",
					Admins: []*v2beta_org.CreateOrganizationRequest_Admin{
						{
							UserType: &v2beta_org.CreateOrganizationRequest_Admin_UserId{
								UserId: "userID",
							},
							Roles: nil,
						},
					},
				},
			},
			want: &command.OrgSetup{
				Name:         "name",
				CustomDomain: "",
				Admins: []*command.OrgSetupAdmin{
					{
						ID: "userID",
					},
				},
			},
		},
		{
			name: "human user",
			args: args{
				request: &v2beta_org.CreateOrganizationRequest{
					Name: "name",
					Admins: []*v2beta_org.CreateOrganizationRequest_Admin{
						{
							UserType: &v2beta_org.CreateOrganizationRequest_Admin_Human{
								Human: &user.AddHumanUserRequest{
									Profile: &user.SetHumanProfile{
										GivenName:  "firstname",
										FamilyName: "lastname",
									},
									Email: &user.SetHumanEmail{
										Email: "email@test.com",
									},
								},
							},
							Roles: nil,
						},
					},
				},
			},
			want: &command.OrgSetup{
				Name:         "name",
				CustomDomain: "",
				Admins: []*command.OrgSetupAdmin{
					{
						Human: &command.AddHuman{
							Username:  "email@test.com",
							FirstName: "firstname",
							LastName:  "lastname",
							Email: command.Email{
								Address: "email@test.com",
							},
							Metadata: make([]*command.AddMetadataEntry, 0),
							Links:    make([]*command.AddLink, 0),
						},
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createOrganizationRequestToCommand(tt.args.request)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}

func Test_createdOrganizationToPb(t *testing.T) {
	now := time.Now()
	type args struct {
		createdOrg *command.CreatedOrg
	}
	tests := []struct {
		name    string
		args    args
		want    *v2beta_org.CreateOrganizationResponse
		wantErr error
	}{
		{
			name: "human user with phone and email code",
			args: args{
				createdOrg: &command.CreatedOrg{
					ObjectDetails: &domain.ObjectDetails{
						Sequence:      1,
						EventDate:     now,
						ResourceOwner: "orgID",
					},
					CreatedAdmins: []*command.CreatedOrgAdmin{
						{
							ID:        "id",
							EmailCode: gu.Ptr("emailCode"),
							PhoneCode: gu.Ptr("phoneCode"),
						},
					},
				},
			},
			want: &v2beta_org.CreateOrganizationResponse{
				CreationDate: timestamppb.New(now),
				Id:           "orgID",
				CreatedAdmins: []*v2beta_org.CreatedAdmin{
					{
						UserId:    "id",
						EmailCode: gu.Ptr("emailCode"),
						PhoneCode: gu.Ptr("phoneCode"),
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := createdOrganizationToPb(tt.args.createdOrg)
			require.ErrorIs(t, err, tt.wantErr)
			assert.Equal(t, tt.want, got)
		})
	}
}
