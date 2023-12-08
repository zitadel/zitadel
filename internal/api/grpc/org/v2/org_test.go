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
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	org "github.com/zitadel/zitadel/pkg/grpc/org/v2beta"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

func Test_addOrganizationRequestToCommand(t *testing.T) {
	type args struct {
		request *org.AddOrganizationRequest
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
				request: &org.AddOrganizationRequest{
					Name: "name",
					Admins: []*org.AddOrganizationRequest_Admin{
						{},
					},
				},
			},
			wantErr: zerrors.ThrowUnimplementedf(nil, "ORGv2-SD2r1", "userType oneOf %T in method AddOrganization not implemented", nil),
		},
		{
			name: "user ID",
			args: args{
				request: &org.AddOrganizationRequest{
					Name: "name",
					Admins: []*org.AddOrganizationRequest_Admin{
						{
							UserType: &org.AddOrganizationRequest_Admin_UserId{
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
				request: &org.AddOrganizationRequest{
					Name: "name",
					Admins: []*org.AddOrganizationRequest_Admin{
						{
							UserType: &org.AddOrganizationRequest_Admin_Human{
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
			got, err := addOrganizationRequestToCommand(tt.args.request)
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
		want    *org.AddOrganizationResponse
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
			want: &org.AddOrganizationResponse{
				Details: &object.Details{
					Sequence:      1,
					ChangeDate:    timestamppb.New(now),
					ResourceOwner: "orgID",
				},
				OrganizationId: "orgID",
				CreatedAdmins: []*org.AddOrganizationResponse_CreatedAdmin{
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
