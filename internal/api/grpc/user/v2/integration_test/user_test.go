//go:build integration

package user_test

import (
	"context"
	"encoding/base64"
	"fmt"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/zitadel/logging"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/pkg/grpc/auth"
	"github.com/zitadel/zitadel/pkg/grpc/idp"
	mgmt "github.com/zitadel/zitadel/pkg/grpc/management"
	"github.com/zitadel/zitadel/pkg/grpc/object/v2"
	user_v1 "github.com/zitadel/zitadel/pkg/grpc/user"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func TestServer_Deprecated_AddHumanUser(t *testing.T) {
	idpResp := Instance.AddGenericOAuthProvider(IamCTX, Instance.DefaultOrg.Id)
	type args struct {
		ctx context.Context
		req *user.AddHumanUserRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.AddHumanUserResponse
		wantErr bool
	}{
		{
			name: "default verification",
			args: args{
				OrgCTX,
				&user.AddHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Instance.DefaultOrg.Id,
						},
					},
					Profile: &user.SetHumanProfile{
						GivenName:         "Donald",
						FamilyName:        "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{},
					Phone: &user.SetHumanPhone{},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			want: &user.AddHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "default verification (org domain ctx)",
			args: args{
				OrgCTX,
				&user.AddHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgDomain{
							OrgDomain: Instance.DefaultOrg.PrimaryDomain,
						},
					},
					Profile: &user.SetHumanProfile{
						GivenName:         "Donald",
						FamilyName:        "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{},
					Phone: &user.SetHumanPhone{},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			want: &user.AddHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "return email verification code",
			args: args{
				OrgCTX,
				&user.AddHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Instance.DefaultOrg.Id,
						},
					},
					Profile: &user.SetHumanProfile{
						GivenName:         "Donald",
						FamilyName:        "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{
						Verification: &user.SetHumanEmail_ReturnCode{
							ReturnCode: &user.ReturnEmailVerificationCode{},
						},
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			want: &user.AddHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
				EmailCode: gu.Ptr("something"),
			},
		},
		{
			name: "custom template",
			args: args{
				OrgCTX,
				&user.AddHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Instance.DefaultOrg.Id,
						},
					},
					Profile: &user.SetHumanProfile{
						GivenName:         "Donald",
						FamilyName:        "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{
						Verification: &user.SetHumanEmail_SendCode{
							SendCode: &user.SendEmailVerificationCode{
								UrlTemplate: gu.Ptr("https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}"),
							},
						},
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			want: &user.AddHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "return phone verification code",
			args: args{
				OrgCTX,
				&user.AddHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Instance.DefaultOrg.Id,
						},
					},
					Profile: &user.SetHumanProfile{
						GivenName:         "Donald",
						FamilyName:        "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{},
					Phone: &user.SetHumanPhone{
						Phone: "+41791234567",
						Verification: &user.SetHumanPhone_ReturnCode{
							ReturnCode: &user.ReturnPhoneVerificationCode{},
						},
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			want: &user.AddHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
				PhoneCode: gu.Ptr("something"),
			},
		},
		{
			name: "custom template error",
			args: args{
				OrgCTX,
				&user.AddHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Instance.DefaultOrg.Id,
						},
					},
					Profile: &user.SetHumanProfile{
						GivenName:         "Donald",
						FamilyName:        "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{
						Verification: &user.SetHumanEmail_SendCode{
							SendCode: &user.SendEmailVerificationCode{
								UrlTemplate: gu.Ptr("{{"),
							},
						},
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing REQUIRED profile",
			args: args{
				OrgCTX,
				&user.AddHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Instance.DefaultOrg.Id,
						},
					},
					Email: &user.SetHumanEmail{
						Verification: &user.SetHumanEmail_ReturnCode{
							ReturnCode: &user.ReturnEmailVerificationCode{},
						},
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing REQUIRED email",
			args: args{
				OrgCTX,
				&user.AddHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Instance.DefaultOrg.Id,
						},
					},
					Profile: &user.SetHumanProfile{
						GivenName:         "Donald",
						FamilyName:        "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "missing idp",
			args: args{
				OrgCTX,
				&user.AddHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Instance.DefaultOrg.Id,
						},
					},
					Profile: &user.SetHumanProfile{
						GivenName:         "Donald",
						FamilyName:        "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{
						Email: "livio@zitadel.com",
						Verification: &user.SetHumanEmail_IsVerified{
							IsVerified: true,
						},
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: false,
						},
					},
					IdpLinks: []*user.IDPLink{
						{
							IdpId:    "idpID",
							UserId:   "userID",
							UserName: "username",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "with idp",
			args: args{
				OrgCTX,
				&user.AddHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Instance.DefaultOrg.Id,
						},
					},
					Profile: &user.SetHumanProfile{
						GivenName:         "Donald",
						FamilyName:        "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{
						Email: "livio@zitadel.com",
						Verification: &user.SetHumanEmail_IsVerified{
							IsVerified: true,
						},
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: false,
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
			want: &user.AddHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "with totp",
			args: args{
				OrgCTX,
				&user.AddHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Instance.DefaultOrg.Id,
						},
					},
					Profile: &user.SetHumanProfile{
						GivenName:         "Donald",
						FamilyName:        "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{
						Email: "livio@zitadel.com",
						Verification: &user.SetHumanEmail_IsVerified{
							IsVerified: true,
						},
					},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: false,
						},
					},
					TotpSecret: gu.Ptr("secret"),
				},
			},
			want: &user.AddHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "password not complexity conform",
			args: args{
				OrgCTX,
				&user.AddHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Instance.DefaultOrg.Id,
						},
					},
					Profile: &user.SetHumanProfile{
						GivenName:         "Donald",
						FamilyName:        "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password: "insufficient",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "hashed password",
			args: args{
				OrgCTX,
				&user.AddHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Instance.DefaultOrg.Id,
						},
					},
					Profile: &user.SetHumanProfile{
						GivenName:         "Donald",
						FamilyName:        "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_HashedPassword{
						HashedPassword: &user.HashedPassword{
							Hash: "$2y$12$hXUrnqdq1RIIYZ2HPytIIe5lXdIvbhqrTvdPsSF7o.jFh817Z6lwm",
						},
					},
				},
			},
			want: &user.AddHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "unsupported hashed password",
			args: args{
				OrgCTX,
				&user.AddHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Instance.DefaultOrg.Id,
						},
					},
					Profile: &user.SetHumanProfile{
						GivenName:         "Donald",
						FamilyName:        "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_HashedPassword{
						HashedPassword: &user.HashedPassword{
							Hash: "$scrypt$ln=16,r=8,p=1$cmFuZG9tc2FsdGlzaGFyZA$Rh+NnJNo1I6nRwaNqbDm6kmADswD1+7FTKZ7Ln9D8nQ",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := fmt.Sprint(time.Now().UnixNano() + int64(i))
			tt.args.req.UserId = &userID
			// In order to prevent unique constraint errors, we set the email to a unique value
			if email := tt.args.req.GetEmail(); email != nil {
				email.Email = fmt.Sprintf("%s@me.now", userID)
			}

			if tt.want != nil {
				tt.want.UserId = userID
			}

			got, err := Client.AddHumanUser(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want.GetUserId(), got.GetUserId())
			if tt.want.GetEmailCode() != "" {
				assert.NotEmpty(t, got.GetEmailCode())
			} else {
				assert.Empty(t, got.GetEmailCode())
			}
			if tt.want.GetPhoneCode() != "" {
				assert.NotEmpty(t, got.GetPhoneCode())
			} else {
				assert.Empty(t, got.GetPhoneCode())
			}
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_Deprecated_AddHumanUser_Permission(t *testing.T) {
	newOrgOwnerEmail := integration.Email()
	newOrg := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), newOrgOwnerEmail)
	type args struct {
		ctx context.Context
		req *user.AddHumanUserRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.AddHumanUserResponse
		wantErr bool
	}{
		{
			name: "System, ok",
			args: args{
				SystemCTX,
				&user.AddHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: newOrg.GetOrganizationId(),
						},
					},
					Profile: &user.SetHumanProfile{
						GivenName:         "Donald",
						FamilyName:        "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{},
					Phone: &user.SetHumanPhone{},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			want: &user.AddHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: newOrg.GetOrganizationId(),
				},
			},
		},
		{
			name: "Instance, ok",
			args: args{
				IamCTX,
				&user.AddHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: newOrg.GetOrganizationId(),
						},
					},
					Profile: &user.SetHumanProfile{
						GivenName:         "Donald",
						FamilyName:        "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{},
					Phone: &user.SetHumanPhone{},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			want: &user.AddHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: newOrg.GetOrganizationId(),
				},
			},
		},
		{
			name: "Org, error",
			args: args{
				OrgCTX,
				&user.AddHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: newOrg.GetOrganizationId(),
						},
					},
					Profile: &user.SetHumanProfile{
						GivenName:         "Donald",
						FamilyName:        "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{},
					Phone: &user.SetHumanPhone{},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "User, error",
			args: args{
				UserCTX,
				&user.AddHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: newOrg.GetOrganizationId(),
						},
					},
					Profile: &user.SetHumanProfile{
						GivenName:         "Donald",
						FamilyName:        "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
					Email: &user.SetHumanEmail{},
					Phone: &user.SetHumanPhone{},
					Metadata: []*user.SetMetadataEntry{
						{
							Key:   "somekey",
							Value: []byte("somevalue"),
						},
					},
					PasswordType: &user.AddHumanUserRequest_Password{
						Password: &user.Password{
							Password:       "DifficultPW666!",
							ChangeRequired: true,
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := fmt.Sprint(time.Now().UnixNano() + int64(i))
			tt.args.req.UserId = &userID
			if email := tt.args.req.GetEmail(); email != nil {
				email.Email = fmt.Sprintf("%s@me.now", userID)
			}

			if tt.want != nil {
				tt.want.UserId = userID
			}

			got, err := Client.AddHumanUser(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			assert.Equal(t, tt.want.GetUserId(), got.GetUserId())
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_Deprecated_UpdateHumanUser(t *testing.T) {
	type args struct {
		ctx context.Context
		req *user.UpdateHumanUserRequest
	}
	tests := []struct {
		name    string
		prepare func(request *user.UpdateHumanUserRequest) error
		args    args
		want    *user.UpdateHumanUserResponse
		wantErr bool
	}{
		{
			name: "not exisiting",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				request.UserId = "notexisiting"
				return nil
			},
			args: args{
				OrgCTX,
				&user.UpdateHumanUserRequest{
					Username: gu.Ptr("changed"),
				},
			},
			wantErr: true,
		},
		{
			name: "change username, ok",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				userID := Instance.CreateHumanUser(OrgCTX).GetUserId()
				request.UserId = userID
				return nil
			},
			args: args{
				OrgCTX,
				&user.UpdateHumanUserRequest{
					Username: gu.Ptr(integration.Username()),
				},
			},
			want: &user.UpdateHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "change profile, ok",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				userID := Instance.CreateHumanUser(OrgCTX).GetUserId()
				request.UserId = userID
				return nil
			},
			args: args{
				OrgCTX,
				&user.UpdateHumanUserRequest{
					Profile: &user.SetHumanProfile{
						GivenName:         "Donald",
						FamilyName:        "Duck",
						NickName:          gu.Ptr("Dukkie"),
						DisplayName:       gu.Ptr("Donald Duck"),
						PreferredLanguage: gu.Ptr("en"),
						Gender:            user.Gender_GENDER_DIVERSE.Enum(),
					},
				},
			},
			want: &user.UpdateHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "change email, ok",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				userID := Instance.CreateHumanUser(OrgCTX).GetUserId()
				request.UserId = userID
				return nil
			},
			args: args{
				OrgCTX,
				&user.UpdateHumanUserRequest{
					Email: &user.SetHumanEmail{
						Email:        "changed@test.com",
						Verification: &user.SetHumanEmail_IsVerified{IsVerified: true},
					},
				},
			},
			want: &user.UpdateHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "change email, code, ok",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				userID := Instance.CreateHumanUser(OrgCTX).GetUserId()
				request.UserId = userID
				return nil
			},
			args: args{
				OrgCTX,
				&user.UpdateHumanUserRequest{
					Email: &user.SetHumanEmail{
						Email:        "changed@test.com",
						Verification: &user.SetHumanEmail_ReturnCode{},
					},
				},
			},
			want: &user.UpdateHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
				EmailCode: gu.Ptr("something"),
			},
		},
		{
			name: "change phone, ok",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				userID := Instance.CreateHumanUser(OrgCTX).GetUserId()
				request.UserId = userID
				return nil
			},
			args: args{
				OrgCTX,
				&user.UpdateHumanUserRequest{
					Phone: &user.SetHumanPhone{
						Phone:        "+41791234567",
						Verification: &user.SetHumanPhone_IsVerified{IsVerified: true},
					},
				},
			},
			want: &user.UpdateHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "change phone, code, ok",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				userID := Instance.CreateHumanUser(OrgCTX).GetUserId()
				request.UserId = userID
				return nil
			},
			args: args{
				OrgCTX,
				&user.UpdateHumanUserRequest{
					Phone: &user.SetHumanPhone{
						Phone:        "+41791234568",
						Verification: &user.SetHumanPhone_ReturnCode{},
					},
				},
			},
			want: &user.UpdateHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
				PhoneCode: gu.Ptr("something"),
			},
		},
		{
			name: "change password, code, ok",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				userID := Instance.CreateHumanUser(OrgCTX).GetUserId()
				request.UserId = userID
				resp, err := Client.PasswordReset(OrgCTX, &user.PasswordResetRequest{
					UserId: userID,
					Medium: &user.PasswordResetRequest_ReturnCode{
						ReturnCode: &user.ReturnPasswordResetCode{},
					},
				})
				if err != nil {
					return err
				}
				request.Password.Verification = &user.SetPassword_VerificationCode{
					VerificationCode: resp.GetVerificationCode(),
				}
				return nil
			},
			args: args{
				OrgCTX,
				&user.UpdateHumanUserRequest{
					Password: &user.SetPassword{
						PasswordType: &user.SetPassword_Password{
							Password: &user.Password{
								Password:       "Password1!",
								ChangeRequired: true,
							},
						},
					},
				},
			},
			want: &user.UpdateHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "change hashed password, code, ok",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				userID := Instance.CreateHumanUser(OrgCTX).GetUserId()
				request.UserId = userID
				resp, err := Client.PasswordReset(OrgCTX, &user.PasswordResetRequest{
					UserId: userID,
					Medium: &user.PasswordResetRequest_ReturnCode{
						ReturnCode: &user.ReturnPasswordResetCode{},
					},
				})
				if err != nil {
					return err
				}
				request.Password.Verification = &user.SetPassword_VerificationCode{
					VerificationCode: resp.GetVerificationCode(),
				}
				return nil
			},
			args: args{
				OrgCTX,
				&user.UpdateHumanUserRequest{
					Password: &user.SetPassword{
						PasswordType: &user.SetPassword_HashedPassword{
							HashedPassword: &user.HashedPassword{
								Hash: "$2y$12$hXUrnqdq1RIIYZ2HPytIIe5lXdIvbhqrTvdPsSF7o.jFh817Z6lwm",
							},
						},
					},
				},
			},
			want: &user.UpdateHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "change hashed password, code, not supported",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				userID := Instance.CreateHumanUser(OrgCTX).GetUserId()
				request.UserId = userID
				resp, err := Client.PasswordReset(OrgCTX, &user.PasswordResetRequest{
					UserId: userID,
					Medium: &user.PasswordResetRequest_ReturnCode{
						ReturnCode: &user.ReturnPasswordResetCode{},
					},
				})
				if err != nil {
					return err
				}
				request.Password = &user.SetPassword{
					Verification: &user.SetPassword_VerificationCode{
						VerificationCode: resp.GetVerificationCode(),
					},
				}
				return nil
			},
			args: args{
				OrgCTX,
				&user.UpdateHumanUserRequest{
					Password: &user.SetPassword{
						PasswordType: &user.SetPassword_HashedPassword{
							HashedPassword: &user.HashedPassword{
								Hash: "$scrypt$ln=16,r=8,p=1$cmFuZG9tc2FsdGlzaGFyZA$Rh+NnJNo1I6nRwaNqbDm6kmADswD1+7FTKZ7Ln9D8nQ",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "change password, old password, ok",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				userID := Instance.CreateHumanUser(OrgCTX).GetUserId()
				request.UserId = userID

				resp, err := Client.PasswordReset(OrgCTX, &user.PasswordResetRequest{
					UserId: userID,
					Medium: &user.PasswordResetRequest_ReturnCode{
						ReturnCode: &user.ReturnPasswordResetCode{},
					},
				})
				if err != nil {
					return err
				}
				pw := "Password1."
				_, err = Client.SetPassword(OrgCTX, &user.SetPasswordRequest{
					UserId: userID,
					NewPassword: &user.Password{
						Password:       pw,
						ChangeRequired: true,
					},
					Verification: &user.SetPasswordRequest_VerificationCode{
						VerificationCode: resp.GetVerificationCode(),
					},
				})
				if err != nil {
					return err
				}
				request.Password.Verification = &user.SetPassword_CurrentPassword{
					CurrentPassword: pw,
				}
				return nil
			},
			args: args{
				OrgCTX,
				&user.UpdateHumanUserRequest{
					Password: &user.SetPassword{
						PasswordType: &user.SetPassword_Password{
							Password: &user.Password{
								Password:       "Password1!",
								ChangeRequired: true,
							},
						},
					},
				},
			},
			want: &user.UpdateHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.prepare(tt.args.req)
			require.NoError(t, err)

			got, err := Client.UpdateHumanUser(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if tt.want.GetEmailCode() != "" {
				assert.NotEmpty(t, got.GetEmailCode())
			} else {
				assert.Empty(t, got.GetEmailCode())
			}
			if tt.want.GetPhoneCode() != "" {
				assert.NotEmpty(t, got.GetPhoneCode())
			} else {
				assert.Empty(t, got.GetPhoneCode())
			}
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_Deprecated_UpdateHumanUser_Permission(t *testing.T) {
	newOrgOwnerEmail := integration.Email()
	newOrg := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), newOrgOwnerEmail)
	newUserID := newOrg.CreatedAdmins[0].GetUserId()
	type args struct {
		ctx context.Context
		req *user.UpdateHumanUserRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.UpdateHumanUserResponse
		wantErr bool
	}{
		{
			name: "system, ok",
			args: args{
				SystemCTX,
				&user.UpdateHumanUserRequest{
					UserId:   newUserID,
					Username: gu.Ptr(integration.Username()),
				},
			},
			want: &user.UpdateHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: newOrg.GetOrganizationId(),
				},
			},
		},
		{
			name: "instance, ok",
			args: args{
				IamCTX,
				&user.UpdateHumanUserRequest{
					UserId:   newUserID,
					Username: gu.Ptr(integration.Username()),
				},
			},
			want: &user.UpdateHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: newOrg.GetOrganizationId(),
				},
			},
		},
		{
			name: "org, error",
			args: args{
				OrgCTX,
				&user.UpdateHumanUserRequest{
					UserId:   newUserID,
					Username: gu.Ptr(integration.Username()),
				},
			},
			wantErr: true,
		},
		{
			name: "user, error",
			args: args{
				UserCTX,
				&user.UpdateHumanUserRequest{
					UserId:   newUserID,
					Username: gu.Ptr(integration.Username()),
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.UpdateHumanUser(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_LockUser(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     *user.LockUserRequest
		prepare func(request *user.LockUserRequest) error
	}
	tests := []struct {
		name    string
		args    args
		want    *user.LockUserResponse
		wantErr bool
	}{
		{
			name: "lock, not existing",
			args: args{
				OrgCTX,
				&user.LockUserRequest{
					UserId: "notexisting",
				},
				func(request *user.LockUserRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "no permission, error",
			args: args{
				UserCTX,
				&user.LockUserRequest{
					UserId: Instance.Users.Get(integration.UserTypeNoPermission).ID,
				},
				func(request *user.LockUserRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "lock, ok",
			args: args{
				OrgCTX,
				&user.LockUserRequest{},
				func(request *user.LockUserRequest) error {
					resp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			want: &user.LockUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "lock machine, ok",
			args: args{
				OrgCTX,
				&user.LockUserRequest{},
				func(request *user.LockUserRequest) error {
					resp := Instance.CreateMachineUser(OrgCTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			want: &user.LockUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "lock, already locked",
			args: args{
				OrgCTX,
				&user.LockUserRequest{},
				func(request *user.LockUserRequest) error {
					resp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = resp.GetUserId()
					_, err := Client.LockUser(OrgCTX, &user.LockUserRequest{
						UserId: resp.GetUserId(),
					})
					return err
				},
			},
			wantErr: true,
		},
		{
			name: "lock machine, already locked",
			args: args{
				OrgCTX,
				&user.LockUserRequest{},
				func(request *user.LockUserRequest) error {
					resp := Instance.CreateMachineUser(OrgCTX)
					request.UserId = resp.GetUserId()
					_, err := Client.LockUser(OrgCTX, &user.LockUserRequest{
						UserId: resp.GetUserId(),
					})
					return err
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.prepare(tt.args.req)
			require.NoError(t, err)

			got, err := Client.LockUser(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_UnLockUser(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     *user.UnlockUserRequest
		prepare func(request *user.UnlockUserRequest) error
	}
	tests := []struct {
		name    string
		args    args
		want    *user.UnlockUserResponse
		wantErr bool
	}{
		{
			name: "unlock, not existing",
			args: args{
				OrgCTX,
				&user.UnlockUserRequest{
					UserId: "notexisting",
				},
				func(request *user.UnlockUserRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "no permission, error",
			args: args{
				UserCTX,
				&user.UnlockUserRequest{
					UserId: Instance.Users.Get(integration.UserTypeNoPermission).ID,
				},
				func(request *user.UnlockUserRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "unlock, not locked",
			args: args{
				ctx: OrgCTX,
				req: &user.UnlockUserRequest{},
				prepare: func(request *user.UnlockUserRequest) error {
					resp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "unlock machine, not locked",
			args: args{
				ctx: OrgCTX,
				req: &user.UnlockUserRequest{},
				prepare: func(request *user.UnlockUserRequest) error {
					resp := Instance.CreateMachineUser(OrgCTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "unlock, ok",
			args: args{
				ctx: OrgCTX,
				req: &user.UnlockUserRequest{},
				prepare: func(request *user.UnlockUserRequest) error {
					resp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = resp.GetUserId()
					_, err := Client.LockUser(OrgCTX, &user.LockUserRequest{
						UserId: resp.GetUserId(),
					})
					return err
				},
			},
			want: &user.UnlockUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "unlock machine, ok",
			args: args{
				ctx: OrgCTX,
				req: &user.UnlockUserRequest{},
				prepare: func(request *user.UnlockUserRequest) error {
					resp := Instance.CreateMachineUser(OrgCTX)
					request.UserId = resp.GetUserId()
					_, err := Client.LockUser(OrgCTX, &user.LockUserRequest{
						UserId: resp.GetUserId(),
					})
					return err
				},
			},
			want: &user.UnlockUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.prepare(tt.args.req)
			require.NoError(t, err)

			got, err := Client.UnlockUser(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_DeactivateUser(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     *user.DeactivateUserRequest
		prepare func(request *user.DeactivateUserRequest) error
	}
	tests := []struct {
		name    string
		args    args
		want    *user.DeactivateUserResponse
		wantErr bool
	}{
		{
			name: "deactivate, not existing",
			args: args{
				OrgCTX,
				&user.DeactivateUserRequest{
					UserId: "notexisting",
				},
				func(request *user.DeactivateUserRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "no permission, error",
			args: args{
				UserCTX,
				&user.DeactivateUserRequest{
					UserId: Instance.Users.Get(integration.UserTypeNoPermission).ID,
				},
				func(request *user.DeactivateUserRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "deactivate, ok",
			args: args{
				OrgCTX,
				&user.DeactivateUserRequest{},
				func(request *user.DeactivateUserRequest) error {
					resp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			want: &user.DeactivateUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "deactivate machine, ok",
			args: args{
				OrgCTX,
				&user.DeactivateUserRequest{},
				func(request *user.DeactivateUserRequest) error {
					resp := Instance.CreateMachineUser(OrgCTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			want: &user.DeactivateUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "deactivate, already deactivated",
			args: args{
				OrgCTX,
				&user.DeactivateUserRequest{},
				func(request *user.DeactivateUserRequest) error {
					resp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = resp.GetUserId()
					_, err := Client.DeactivateUser(OrgCTX, &user.DeactivateUserRequest{
						UserId: resp.GetUserId(),
					})
					return err
				},
			},
			wantErr: true,
		},
		{
			name: "deactivate machine, already deactivated",
			args: args{
				OrgCTX,
				&user.DeactivateUserRequest{},
				func(request *user.DeactivateUserRequest) error {
					resp := Instance.CreateMachineUser(OrgCTX)
					request.UserId = resp.GetUserId()
					_, err := Client.DeactivateUser(OrgCTX, &user.DeactivateUserRequest{
						UserId: resp.GetUserId(),
					})
					return err
				},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.prepare(tt.args.req)
			require.NoError(t, err)

			got, err := Client.DeactivateUser(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_ReactivateUser(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     *user.ReactivateUserRequest
		prepare func(request *user.ReactivateUserRequest) error
	}
	tests := []struct {
		name    string
		args    args
		want    *user.ReactivateUserResponse
		wantErr bool
	}{
		{
			name: "reactivate, not existing",
			args: args{
				OrgCTX,
				&user.ReactivateUserRequest{
					UserId: "notexisting",
				},
				func(request *user.ReactivateUserRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "no permission, error",
			args: args{
				UserCTX,
				&user.ReactivateUserRequest{
					UserId: Instance.Users.Get(integration.UserTypeNoPermission).ID,
				},
				func(request *user.ReactivateUserRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "reactivate, not deactivated",
			args: args{
				ctx: OrgCTX,
				req: &user.ReactivateUserRequest{},
				prepare: func(request *user.ReactivateUserRequest) error {
					resp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "reactivate machine, not deactivated",
			args: args{
				ctx: OrgCTX,
				req: &user.ReactivateUserRequest{},
				prepare: func(request *user.ReactivateUserRequest) error {
					resp := Instance.CreateMachineUser(OrgCTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "reactivate, ok",
			args: args{
				ctx: OrgCTX,
				req: &user.ReactivateUserRequest{},
				prepare: func(request *user.ReactivateUserRequest) error {
					resp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = resp.GetUserId()
					_, err := Client.DeactivateUser(OrgCTX, &user.DeactivateUserRequest{
						UserId: resp.GetUserId(),
					})
					return err
				},
			},
			want: &user.ReactivateUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "reactivate machine, ok",
			args: args{
				ctx: OrgCTX,
				req: &user.ReactivateUserRequest{},
				prepare: func(request *user.ReactivateUserRequest) error {
					resp := Instance.CreateMachineUser(OrgCTX)
					request.UserId = resp.GetUserId()
					_, err := Client.DeactivateUser(OrgCTX, &user.DeactivateUserRequest{
						UserId: resp.GetUserId(),
					})
					return err
				},
			},
			want: &user.ReactivateUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.prepare(tt.args.req)
			require.NoError(t, err)

			got, err := Client.ReactivateUser(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_DeleteUser(t *testing.T) {
	projectResp := Instance.CreateProject(OrgCTX, t, Instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)

	type args struct {
		req     *user.DeleteUserRequest
		prepare func(*testing.T, *user.DeleteUserRequest) context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    *user.DeleteUserResponse
		wantErr bool
	}{
		{
			name: "remove, not existing",
			args: args{
				&user.DeleteUserRequest{
					UserId: "notexisting",
				},
				func(*testing.T, *user.DeleteUserRequest) context.Context { return OrgCTX },
			},
			wantErr: true,
		},
		{
			name: "remove human, ok",
			args: args{
				req: &user.DeleteUserRequest{},
				prepare: func(_ *testing.T, request *user.DeleteUserRequest) context.Context {
					resp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = resp.GetUserId()
					return OrgCTX
				},
			},
			want: &user.DeleteUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "remove machine, ok",
			args: args{
				req: &user.DeleteUserRequest{},
				prepare: func(_ *testing.T, request *user.DeleteUserRequest) context.Context {
					resp := Instance.CreateMachineUser(OrgCTX)
					request.UserId = resp.GetUserId()
					return OrgCTX
				},
			},
			want: &user.DeleteUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "remove dependencies, ok",
			args: args{
				req: &user.DeleteUserRequest{},
				prepare: func(_ *testing.T, request *user.DeleteUserRequest) context.Context {
					resp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = resp.GetUserId()
					Instance.CreateProjectUserGrant(t, OrgCTX, Instance.DefaultOrg.GetId(), projectResp.GetId(), request.UserId)
					Instance.CreateProjectMembership(t, OrgCTX, projectResp.GetId(), request.UserId)
					Instance.CreateOrgMembership(t, OrgCTX, Instance.DefaultOrg.Id, request.UserId)
					return OrgCTX
				},
			},
			want: &user.DeleteUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "remove self, ok",
			args: args{
				req: &user.DeleteUserRequest{},
				prepare: func(t *testing.T, request *user.DeleteUserRequest) context.Context {
					removeUser, err := Client.CreateUser(OrgCTX, &user.CreateUserRequest{
						OrganizationId: Instance.DefaultOrg.Id,
						UserType: &user.CreateUserRequest_Human_{
							Human: &user.CreateUserRequest_Human{
								Profile: &user.SetHumanProfile{
									GivenName:  "givenName",
									FamilyName: "familyName",
								},
								Email: &user.SetHumanEmail{
									Email:        integration.Email(),
									Verification: &user.SetHumanEmail_IsVerified{IsVerified: true},
								},
							},
						},
					})
					require.NoError(t, err)

					// allow self management incl. deletion
					Instance.CreateOrgMembership(t, OrgCTX, Instance.DefaultOrg.Id, removeUser.Id, "ORG_USER_SELF_MANAGER")

					request.UserId = removeUser.Id
					Instance.RegisterUserPasskey(OrgCTX, removeUser.Id)
					token := createVerifiedWebAuthNSession(LoginCTX, t, removeUser.Id)
					return integration.WithAuthorizationToken(UserCTX, token)
				},
			},
			want: &user.DeleteUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctx := tt.args.prepare(t, tt.args.req)
			got, err := Client.DeleteUser(ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}
func createVerifiedWebAuthNSession(ctx context.Context, t *testing.T, userID string) string {
	// check if user is already processed
	retryDuration, tick := integration.WaitForAndTickWithMaxDuration(ctx, time.Minute)
	require.EventuallyWithT(t, func(collect *assert.CollectT) {
		_, err := Client.GetUserByID(ctx, &user.GetUserByIDRequest{UserId: userID})
		require.NoError(collect, err)
	}, retryDuration, tick)

	_, token, _, _ := Instance.CreateVerifiedWebAuthNSession(t, ctx, userID)
	return token
}

func ctxFromNewUserWithRegisteredPasswordlessLegacy(t *testing.T) (context.Context, string, *auth.AddMyPasswordlessResponse) {
	userID := Instance.CreateHumanUser(OrgCTX).GetUserId()
	Instance.RegisterUserPasskey(OrgCTX, userID)
	sessionToken := createVerifiedWebAuthNSession(LoginCTX, t, userID)
	ctx := integration.WithAuthorizationToken(OrgCTX, sessionToken)

	pkr, err := Instance.Client.Auth.AddMyPasswordless(ctx, &auth.AddMyPasswordlessRequest{})
	require.NoError(t, err)
	require.NotEmpty(t, pkr.GetKey())
	return ctx, userID, pkr
}

func ctxFromNewUserWithVerifiedPasswordlessLegacy(t *testing.T) (context.Context, string) {
	ctx, userID, pkr := ctxFromNewUserWithRegisteredPasswordlessLegacy(t)

	attestationResponse, err := Instance.WebAuthN.CreateAttestationResponseData(pkr.GetKey().GetPublicKey())
	require.NoError(t, err)

	_, err = Instance.Client.Auth.VerifyMyPasswordless(ctx, &auth.VerifyMyPasswordlessRequest{
		Verification: &user_v1.WebAuthNVerification{
			TokenName:           "Mickey",
			PublicKeyCredential: attestationResponse,
		},
	})
	require.NoError(t, err)
	return ctx, userID
}

func TestServer_ListAuthenticationMethodTypes(t *testing.T) {
	userIDWithoutAuth := Instance.CreateHumanUser(OrgCTX).GetUserId()

	userIDWithPasskey := Instance.CreateHumanUser(OrgCTX).GetUserId()
	Instance.RegisterUserPasskey(OrgCTX, userIDWithPasskey)

	userMultipleAuth := Instance.CreateHumanUser(OrgCTX).GetUserId()
	Instance.RegisterUserPasskey(OrgCTX, userMultipleAuth)
	provider, err := Instance.Client.Mgmt.AddGenericOIDCProvider(OrgCTX, &mgmt.AddGenericOIDCProviderRequest{
		Name:         "ListAuthenticationMethodTypes",
		Issuer:       "https://example.com",
		ClientId:     "client_id",
		ClientSecret: "client_secret",
	})
	require.NoError(t, err)
	_, err = Instance.Client.Mgmt.AddCustomLoginPolicy(OrgCTX, &mgmt.AddCustomLoginPolicyRequest{})
	require.Condition(t, func() bool {
		code := status.Convert(err).Code()
		return code == codes.AlreadyExists || code == codes.OK
	})
	_, err = Instance.Client.Mgmt.AddIDPToLoginPolicy(OrgCTX, &mgmt.AddIDPToLoginPolicyRequest{
		IdpId:     provider.GetId(),
		OwnerType: idp.IDPOwnerType_IDP_OWNER_TYPE_ORG,
	})
	require.NoError(t, err)
	_, err = Instance.Client.UserV2.AddIDPLink(OrgCTX, &user.AddIDPLinkRequest{UserId: userMultipleAuth, IdpLink: &user.IDPLink{
		IdpId:    provider.GetId(),
		UserId:   "external-id",
		UserName: "displayName",
	}})
	require.NoError(t, err)
	// This should not remove the user IDP links
	_, err = Instance.Client.Mgmt.RemoveIDPFromLoginPolicy(OrgCTX, &mgmt.RemoveIDPFromLoginPolicyRequest{
		IdpId: provider.GetId(),
	})
	require.NoError(t, err)

	_, userLegacyID := ctxFromNewUserWithVerifiedPasswordlessLegacy(t)
	require.NoError(t, err)

	type args struct {
		ctx context.Context
		req *user.ListAuthenticationMethodTypesRequest
	}
	tests := []struct {
		name string
		args args
		want *user.ListAuthenticationMethodTypesResponse
	}{
		{
			name: "no auth",
			args: args{
				OrgCTX,
				&user.ListAuthenticationMethodTypesRequest{
					UserId: userIDWithoutAuth,
				},
			},
			want: &user.ListAuthenticationMethodTypesResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
				},
			},
		},
		{
			name: "with auth (passkey)",
			args: args{
				OrgCTX,
				&user.ListAuthenticationMethodTypesRequest{
					UserId: userIDWithPasskey,
				},
			},
			want: &user.ListAuthenticationMethodTypesResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
				},
				AuthMethodTypes: []user.AuthenticationMethodType{
					user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_PASSKEY,
				},
			},
		},
		{
			name: "with auth (passkey) with domain",
			args: args{
				OrgCTX,
				&user.ListAuthenticationMethodTypesRequest{
					UserId: userIDWithPasskey,
					DomainQuery: &user.DomainQuery{
						Domain: Instance.Domain,
					},
				},
			},
			want: &user.ListAuthenticationMethodTypesResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
				},
				AuthMethodTypes: []user.AuthenticationMethodType{
					user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_PASSKEY,
				},
			},
		},
		{
			name: "with auth (passkey) with wrong domain",
			args: args{
				OrgCTX,
				&user.ListAuthenticationMethodTypesRequest{
					UserId: userIDWithPasskey,
					DomainQuery: &user.DomainQuery{
						Domain: "notexistent",
					},
				},
			},
			want: &user.ListAuthenticationMethodTypesResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
				},
			},
		},
		{
			name: "with auth (passkey) with legacy",
			args: args{
				OrgCTX,
				&user.ListAuthenticationMethodTypesRequest{
					UserId: userLegacyID,
					DomainQuery: &user.DomainQuery{
						Domain: "notexistent",
					},
				},
			},
			want: &user.ListAuthenticationMethodTypesResponse{
				Details: &object.ListDetails{
					TotalResult: 0,
				},
			},
		},
		{
			name: "with auth (passkey) with legacy included",
			args: args{
				OrgCTX,
				&user.ListAuthenticationMethodTypesRequest{
					UserId: userLegacyID,
					DomainQuery: &user.DomainQuery{
						Domain:               "notexistent",
						IncludeWithoutDomain: true,
					},
				},
			},
			want: &user.ListAuthenticationMethodTypesResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
				},
				AuthMethodTypes: []user.AuthenticationMethodType{
					user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_PASSKEY,
				},
			},
		},
		{
			name: "multiple auth",
			args: args{
				OrgCTX,
				&user.ListAuthenticationMethodTypesRequest{
					UserId: userMultipleAuth,
				},
			},
			want: &user.ListAuthenticationMethodTypesResponse{
				Details: &object.ListDetails{
					TotalResult: 2,
				},
				AuthMethodTypes: []user.AuthenticationMethodType{
					user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_PASSKEY,
					user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_IDP,
				},
			},
		},
		{
			name: "multiple auth with domain",
			args: args{
				OrgCTX,
				&user.ListAuthenticationMethodTypesRequest{
					UserId: userMultipleAuth,
					DomainQuery: &user.DomainQuery{
						Domain: Instance.Domain,
					},
				},
			},
			want: &user.ListAuthenticationMethodTypesResponse{
				Details: &object.ListDetails{
					TotalResult: 2,
				},
				AuthMethodTypes: []user.AuthenticationMethodType{
					user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_PASSKEY,
					user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_IDP,
				},
			},
		},
		{
			name: "multiple auth with wrong domain",
			args: args{
				OrgCTX,
				&user.ListAuthenticationMethodTypesRequest{
					UserId: userMultipleAuth,
					DomainQuery: &user.DomainQuery{
						Domain: "notexistent",
					},
				},
			},
			want: &user.ListAuthenticationMethodTypesResponse{
				Details: &object.ListDetails{
					TotalResult: 1,
				},
				AuthMethodTypes: []user.AuthenticationMethodType{
					user.AuthenticationMethodType_AUTHENTICATION_METHOD_TYPE_IDP,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(tt.args.ctx, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.ListAuthenticationMethodTypes(tt.args.ctx, tt.args.req)
				require.NoError(ttt, err)
				if !assert.Equal(ttt, tt.want.GetDetails().GetTotalResult(), got.GetDetails().GetTotalResult()) {
					return
				}
				assert.Equal(ttt, tt.want.GetAuthMethodTypes(), got.GetAuthMethodTypes())
				integration.AssertListDetails(ttt, tt.want, got)
			}, retryDuration, tick, "timeout waiting for expected auth methods result")
		})
	}
}

func TestServer_ListAuthenticationFactors(t *testing.T) {
	tests := []struct {
		name    string
		args    *user.ListAuthenticationFactorsRequest
		want    *user.ListAuthenticationFactorsResponse
		dep     func(args *user.ListAuthenticationFactorsRequest, want *user.ListAuthenticationFactorsResponse)
		wantErr bool
		ctx     context.Context
	}{
		{
			name: "no auth",
			args: &user.ListAuthenticationFactorsRequest{},
			want: &user.ListAuthenticationFactorsResponse{
				Result: nil,
			},
			dep: func(args *user.ListAuthenticationFactorsRequest, want *user.ListAuthenticationFactorsResponse) {
				userIDWithoutAuth := Instance.CreateHumanUser(OrgCTX).GetUserId()
				args.UserId = userIDWithoutAuth
			},
			ctx: OrgCTX,
		},
		{
			name: "with u2f",
			args: &user.ListAuthenticationFactorsRequest{},
			want: &user.ListAuthenticationFactorsResponse{
				Result: []*user.AuthFactor{
					{
						State: user.AuthFactorState_AUTH_FACTOR_STATE_READY,
					},
				},
			},
			dep: func(args *user.ListAuthenticationFactorsRequest, want *user.ListAuthenticationFactorsResponse) {
				userWithU2F := Instance.CreateHumanUser(OrgCTX).GetUserId()
				U2FId := Instance.RegisterUserU2F(OrgCTX, userWithU2F)

				args.UserId = userWithU2F
				want.Result[0].Type = &user.AuthFactor_U2F{
					U2F: &user.AuthFactorU2F{
						Id:   U2FId,
						Name: "nice name",
					},
				}
			},
			ctx: OrgCTX,
		},
		{
			name: "with totp, u2f",
			args: &user.ListAuthenticationFactorsRequest{},
			want: &user.ListAuthenticationFactorsResponse{
				Result: []*user.AuthFactor{
					{
						State: user.AuthFactorState_AUTH_FACTOR_STATE_READY,
						Type: &user.AuthFactor_Otp{
							Otp: &user.AuthFactorOTP{},
						},
					},
					{
						State: user.AuthFactorState_AUTH_FACTOR_STATE_READY,
					},
				},
			},
			dep: func(args *user.ListAuthenticationFactorsRequest, want *user.ListAuthenticationFactorsResponse) {
				userWithTOTP := Instance.CreateHumanUserWithTOTP(OrgCTX, "secret").GetUserId()
				U2FIdWithTOTP := Instance.RegisterUserU2F(OrgCTX, userWithTOTP)

				args.UserId = userWithTOTP
				want.Result[1].Type = &user.AuthFactor_U2F{
					U2F: &user.AuthFactorU2F{
						Id:   U2FIdWithTOTP,
						Name: "nice name",
					},
				}
			},
			ctx: OrgCTX,
		},
		{
			name: "with totp, u2f filtered",
			args: &user.ListAuthenticationFactorsRequest{
				AuthFactors: []user.AuthFactors{user.AuthFactors_U2F},
			},
			want: &user.ListAuthenticationFactorsResponse{
				Result: []*user.AuthFactor{
					{
						State: user.AuthFactorState_AUTH_FACTOR_STATE_READY,
					},
				},
			},
			dep: func(args *user.ListAuthenticationFactorsRequest, want *user.ListAuthenticationFactorsResponse) {
				userWithTOTP := Instance.CreateHumanUserWithTOTP(OrgCTX, "secret").GetUserId()
				U2FIdWithTOTP := Instance.RegisterUserU2F(OrgCTX, userWithTOTP)

				args.UserId = userWithTOTP
				want.Result[0].Type = &user.AuthFactor_U2F{
					U2F: &user.AuthFactorU2F{
						Id:   U2FIdWithTOTP,
						Name: "nice name",
					},
				}
			},
			ctx: OrgCTX,
		},
		{
			name: "with sms",
			args: &user.ListAuthenticationFactorsRequest{},
			want: &user.ListAuthenticationFactorsResponse{
				Result: []*user.AuthFactor{
					{
						State: user.AuthFactorState_AUTH_FACTOR_STATE_READY,
						Type: &user.AuthFactor_OtpSms{
							OtpSms: &user.AuthFactorOTPSMS{},
						},
					},
				},
			},
			dep: func(args *user.ListAuthenticationFactorsRequest, want *user.ListAuthenticationFactorsResponse) {
				userWithSMS := Instance.CreateHumanUserVerified(OrgCTX, Instance.DefaultOrg.GetId(), integration.Email(), integration.Phone()).GetUserId()
				Instance.RegisterUserOTPSMS(OrgCTX, userWithSMS)

				args.UserId = userWithSMS
			},
			ctx: OrgCTX,
		},
		{
			name: "with email",
			args: &user.ListAuthenticationFactorsRequest{},
			want: &user.ListAuthenticationFactorsResponse{
				Result: []*user.AuthFactor{
					{
						State: user.AuthFactorState_AUTH_FACTOR_STATE_READY,
						Type: &user.AuthFactor_OtpEmail{
							OtpEmail: &user.AuthFactorOTPEmail{},
						},
					},
				},
			},
			dep: func(args *user.ListAuthenticationFactorsRequest, want *user.ListAuthenticationFactorsResponse) {
				userWithEmail := Instance.CreateHumanUserVerified(OrgCTX, Instance.DefaultOrg.GetId(), integration.Email(), integration.Phone()).GetUserId()
				Instance.RegisterUserOTPEmail(OrgCTX, userWithEmail)

				args.UserId = userWithEmail
			},
			ctx: OrgCTX,
		},
		{
			name: "with not ready u2f",
			args: &user.ListAuthenticationFactorsRequest{},
			want: &user.ListAuthenticationFactorsResponse{
				Result: []*user.AuthFactor{},
			},
			dep: func(args *user.ListAuthenticationFactorsRequest, want *user.ListAuthenticationFactorsResponse) {
				userWithNotReadyU2F := Instance.CreateHumanUser(OrgCTX).GetUserId()
				_, err := Instance.Client.UserV2.RegisterU2F(OrgCTX, &user.RegisterU2FRequest{
					UserId: userWithNotReadyU2F,
					Domain: Instance.Domain,
				})
				logging.OnError(err).Panic("Could not register u2f")

				args.UserId = userWithNotReadyU2F
			},
			ctx: OrgCTX,
		},
		{
			name: "with not ready u2f state filtered",
			args: &user.ListAuthenticationFactorsRequest{
				States: []user.AuthFactorState{user.AuthFactorState_AUTH_FACTOR_STATE_NOT_READY},
			},
			want: &user.ListAuthenticationFactorsResponse{
				Result: []*user.AuthFactor{
					{
						State: user.AuthFactorState_AUTH_FACTOR_STATE_NOT_READY,
					},
				},
			},
			dep: func(args *user.ListAuthenticationFactorsRequest, want *user.ListAuthenticationFactorsResponse) {
				userWithNotReadyU2F := Instance.CreateHumanUser(OrgCTX).GetUserId()
				U2FNotReady, err := Instance.Client.UserV2.RegisterU2F(OrgCTX, &user.RegisterU2FRequest{
					UserId: userWithNotReadyU2F,
					Domain: Instance.Domain,
				})
				logging.OnError(err).Panic("Could not register u2f")

				args.UserId = userWithNotReadyU2F
				want.Result[0].Type = &user.AuthFactor_U2F{
					U2F: &user.AuthFactorU2F{
						Id:   U2FNotReady.GetU2FId(),
						Name: "",
					},
				}
			},
			ctx: OrgCTX,
		},
		{
			name: "with no userId",
			args: &user.ListAuthenticationFactorsRequest{
				UserId: "",
			},
			ctx:     OrgCTX,
			wantErr: true,
		},
		{
			name: "with no permission",
			args: &user.ListAuthenticationFactorsRequest{},
			dep: func(args *user.ListAuthenticationFactorsRequest, want *user.ListAuthenticationFactorsResponse) {
				userWithTOTP := Instance.CreateHumanUserWithTOTP(OrgCTX, "totp").GetUserId()

				args.UserId = userWithTOTP
			},
			ctx:     UserCTX,
			wantErr: true,
		},
		{
			name: "with unknown user",
			args: &user.ListAuthenticationFactorsRequest{
				UserId: "unknown",
			},
			want: &user.ListAuthenticationFactorsResponse{},
			ctx:  OrgCTX,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.dep != nil {
				tt.dep(tt.args, tt.want)
			}

			retryDuration, tick := integration.WaitForAndTickWithMaxDuration(OrgCTX, time.Minute)
			require.EventuallyWithT(t, func(ttt *assert.CollectT) {
				got, err := Client.ListAuthenticationFactors(tt.ctx, tt.args)
				if tt.wantErr {
					require.Error(ttt, err)
					return
				}
				require.NoError(ttt, err)

				assert.ElementsMatch(ttt, tt.want.GetResult(), got.GetResult())
			}, retryDuration, tick, "timeout waiting for expected auth methods result")
		})
	}
}

func TestServer_CreateInviteCode(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     *user.CreateInviteCodeRequest
		prepare func(request *user.CreateInviteCodeRequest) error
	}
	tests := []struct {
		name    string
		args    args
		want    *user.CreateInviteCodeResponse
		wantErr bool
	}{
		{
			name: "create, not existing",
			args: args{
				OrgCTX,
				&user.CreateInviteCodeRequest{
					UserId: "notexisting",
				},
				func(request *user.CreateInviteCodeRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "create, ok",
			args: args{
				ctx: OrgCTX,
				req: &user.CreateInviteCodeRequest{},
				prepare: func(request *user.CreateInviteCodeRequest) error {
					resp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			want: &user.CreateInviteCodeResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "create, invalid template",
			args: args{
				ctx: OrgCTX,
				req: &user.CreateInviteCodeRequest{
					Verification: &user.CreateInviteCodeRequest_SendCode{
						SendCode: &user.SendInviteCode{
							UrlTemplate: gu.Ptr("{{"),
						},
					},
				},
				prepare: func(request *user.CreateInviteCodeRequest) error {
					resp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "create, valid template",
			args: args{
				ctx: OrgCTX,
				req: &user.CreateInviteCodeRequest{
					Verification: &user.CreateInviteCodeRequest_SendCode{
						SendCode: &user.SendInviteCode{
							UrlTemplate:     gu.Ptr("https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}"),
							ApplicationName: gu.Ptr("TestApp"),
						},
					},
				},
				prepare: func(request *user.CreateInviteCodeRequest) error {
					resp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			want: &user.CreateInviteCodeResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "recreate",
			args: args{
				ctx: OrgCTX,
				req: &user.CreateInviteCodeRequest{},
				prepare: func(request *user.CreateInviteCodeRequest) error {
					resp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = resp.GetUserId()
					_, err := Instance.Client.UserV2.CreateInviteCode(OrgCTX, &user.CreateInviteCodeRequest{
						UserId: resp.GetUserId(),
						Verification: &user.CreateInviteCodeRequest_SendCode{
							SendCode: &user.SendInviteCode{
								UrlTemplate:     gu.Ptr("https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}"),
								ApplicationName: gu.Ptr("TestApp"),
							},
						},
					})
					return err
				},
			},
			want: &user.CreateInviteCodeResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
		{
			name: "create, return code, ok",
			args: args{
				ctx: OrgCTX,
				req: &user.CreateInviteCodeRequest{
					Verification: &user.CreateInviteCodeRequest_ReturnCode{
						ReturnCode: &user.ReturnInviteCode{},
					},
				},
				prepare: func(request *user.CreateInviteCodeRequest) error {
					resp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			want: &user.CreateInviteCodeResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
				InviteCode: gu.Ptr("something"),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.prepare(tt.args.req)
			require.NoError(t, err)

			got, err := Client.CreateInviteCode(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
			if tt.want.GetInviteCode() != "" {
				assert.NotEmpty(t, got.GetInviteCode())
			} else {
				assert.Empty(t, got.GetInviteCode())
			}
		})
	}
}

func TestServer_ResendInviteCode(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     *user.ResendInviteCodeRequest
		prepare func(request *user.ResendInviteCodeRequest) error
	}
	tests := []struct {
		name    string
		args    args
		want    *user.ResendInviteCodeResponse
		wantErr bool
	}{
		{
			name: "user not existing",
			args: args{
				OrgCTX,
				&user.ResendInviteCodeRequest{
					UserId: "notexisting",
				},
				func(request *user.ResendInviteCodeRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "code not existing",
			args: args{
				ctx: OrgCTX,
				req: &user.ResendInviteCodeRequest{},
				prepare: func(request *user.ResendInviteCodeRequest) error {
					resp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "code not sent before",
			args: args{
				ctx: OrgCTX,
				req: &user.ResendInviteCodeRequest{},
				prepare: func(request *user.ResendInviteCodeRequest) error {
					userResp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = userResp.GetUserId()
					Instance.CreateInviteCode(OrgCTX, userResp.GetUserId())
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "resend, ok",
			args: args{
				ctx: OrgCTX,
				req: &user.ResendInviteCodeRequest{},
				prepare: func(request *user.ResendInviteCodeRequest) error {
					resp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = resp.GetUserId()
					_, err := Instance.Client.UserV2.CreateInviteCode(OrgCTX, &user.CreateInviteCodeRequest{
						UserId: resp.GetUserId(),
					})
					return err
				},
			},
			want: &user.ResendInviteCodeResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.prepare(tt.args.req)
			require.NoError(t, err)

			got, err := Client.ResendInviteCode(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_VerifyInviteCode(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     *user.VerifyInviteCodeRequest
		prepare func(request *user.VerifyInviteCodeRequest) error
	}
	tests := []struct {
		name    string
		args    args
		want    *user.VerifyInviteCodeResponse
		wantErr bool
	}{
		{
			name: "user not existing",
			args: args{
				OrgCTX,
				&user.VerifyInviteCodeRequest{
					UserId: "notexisting",
				},
				func(request *user.VerifyInviteCodeRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "code not existing",
			args: args{
				ctx: OrgCTX,
				req: &user.VerifyInviteCodeRequest{},
				prepare: func(request *user.VerifyInviteCodeRequest) error {
					resp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "invalid code",
			args: args{
				ctx: OrgCTX,
				req: &user.VerifyInviteCodeRequest{
					VerificationCode: "invalid",
				},
				prepare: func(request *user.VerifyInviteCodeRequest) error {
					userResp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = userResp.GetUserId()
					Instance.CreateInviteCode(OrgCTX, userResp.GetUserId())
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "verify, ok",
			args: args{
				ctx: OrgCTX,
				req: &user.VerifyInviteCodeRequest{},
				prepare: func(request *user.VerifyInviteCodeRequest) error {
					userResp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = userResp.GetUserId()
					codeResp := Instance.CreateInviteCode(OrgCTX, userResp.GetUserId())
					request.VerificationCode = codeResp.GetInviteCode()
					return nil
				},
			},
			want: &user.VerifyInviteCodeResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.prepare(tt.args.req)
			require.NoError(t, err)

			got, err := Client.VerifyInviteCode(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_HumanMFAInitSkipped(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     *user.HumanMFAInitSkippedRequest
		prepare func(request *user.HumanMFAInitSkippedRequest) error
	}
	tests := []struct {
		name       string
		args       args
		want       *user.HumanMFAInitSkippedResponse
		checkState func(t *testing.T, userID string, resp *user.HumanMFAInitSkippedResponse)
		wantErr    bool
	}{
		{
			name: "user not existing",
			args: args{
				OrgCTX,
				&user.HumanMFAInitSkippedRequest{
					UserId: "notexisting",
				},
				func(request *user.HumanMFAInitSkippedRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "ok",
			args: args{
				OrgCTX,
				&user.HumanMFAInitSkippedRequest{},
				func(request *user.HumanMFAInitSkippedRequest) error {
					resp := Instance.CreateHumanUser(OrgCTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			want: &user.HumanMFAInitSkippedResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
			checkState: func(t *testing.T, userID string, resp *user.HumanMFAInitSkippedResponse) {
				state, err := Client.GetUserByID(OrgCTX, &user.GetUserByIDRequest{
					UserId: userID,
				})
				require.NoError(t, err)
				integration.EqualProto(t,
					state.GetUser().GetHuman().GetMfaInitSkipped(),
					resp.GetDetails().GetChangeDate(),
				)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.prepare(tt.args.req)
			require.NoError(t, err)
			got, err := Client.HumanMFAInitSkipped(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
			if tt.checkState != nil {
				tt.checkState(t, tt.args.req.GetUserId(), got)
			}
		})
	}
}

func TestServer_CreateUser(t *testing.T) {
	type args struct {
		ctx context.Context
		req *user.CreateUserRequest
	}
	type testCase struct {
		args    args
		want    *user.CreateUserResponse
		wantErr bool
	}
	tests := []struct {
		name     string
		testCase func(runId string) testCase
	}{
		{
			name: "default verification",
			testCase: func(runId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				email := username + "@example.com"
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							OrganizationId: Instance.DefaultOrg.Id,
							Username:       &username,
							UserType: &user.CreateUserRequest_Human_{
								Human: &user.CreateUserRequest_Human{
									Profile: &user.SetHumanProfile{
										GivenName:  "Donald",
										FamilyName: "Duck",
									},
									Email: &user.SetHumanEmail{
										Email: email,
									},
								},
							},
						},
					},
					want: &user.CreateUserResponse{
						Id: "is generated",
					},
					wantErr: false,
				}
			},
		},
		{
			name: "return email verification code",
			testCase: func(runId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				email := username + "@example.com"
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							OrganizationId: Instance.DefaultOrg.Id,
							Username:       &username,
							UserType: &user.CreateUserRequest_Human_{
								Human: &user.CreateUserRequest_Human{
									Profile: &user.SetHumanProfile{
										GivenName:  "Donald",
										FamilyName: "Duck",
									},
									Email: &user.SetHumanEmail{
										Email: email,
										Verification: &user.SetHumanEmail_ReturnCode{
											ReturnCode: &user.ReturnEmailVerificationCode{},
										},
									},
								},
							},
						},
					},
					want: &user.CreateUserResponse{
						Id:        "is generated",
						EmailCode: gu.Ptr("something"),
					},
				}
			},
		},
		{
			name: "custom template",
			testCase: func(runId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				email := username + "@example.com"
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							OrganizationId: Instance.DefaultOrg.Id,
							Username:       &username,
							UserType: &user.CreateUserRequest_Human_{
								Human: &user.CreateUserRequest_Human{
									Profile: &user.SetHumanProfile{
										GivenName:  "Donald",
										FamilyName: "Duck",
									},
									Email: &user.SetHumanEmail{
										Email: email,
										Verification: &user.SetHumanEmail_SendCode{
											SendCode: &user.SendEmailVerificationCode{
												UrlTemplate: gu.Ptr("https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}"),
											},
										},
									},
								},
							},
						},
					},
					want: &user.CreateUserResponse{
						Id: "is generated",
					},
				}
			},
		},
		{
			name: "return phone verification code",
			testCase: func(runId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				email := username + "@example.com"
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							OrganizationId: Instance.DefaultOrg.Id,
							Username:       &username,
							UserType: &user.CreateUserRequest_Human_{
								Human: &user.CreateUserRequest_Human{
									Profile: &user.SetHumanProfile{
										GivenName:  "Donald",
										FamilyName: "Duck",
									},
									Email: &user.SetHumanEmail{
										Email: email,
									},
									Phone: &user.SetHumanPhone{
										Phone: "+41791234567",
										Verification: &user.SetHumanPhone_ReturnCode{
											ReturnCode: &user.ReturnPhoneVerificationCode{},
										},
									},
								},
							},
						},
					},
					want: &user.CreateUserResponse{
						Id:        "is generated",
						PhoneCode: gu.Ptr("something"),
					},
				}
			},
		},
		{
			name: "custom template error",
			testCase: func(runId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				email := username + "@example.com"
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							OrganizationId: Instance.DefaultOrg.Id,
							Username:       &username,
							UserType: &user.CreateUserRequest_Human_{
								Human: &user.CreateUserRequest_Human{
									Profile: &user.SetHumanProfile{
										GivenName:  "Donald",
										FamilyName: "Duck",
									},
									Email: &user.SetHumanEmail{
										Email: email,
										Verification: &user.SetHumanEmail_SendCode{
											SendCode: &user.SendEmailVerificationCode{
												UrlTemplate: gu.Ptr("{{"),
											},
										},
									},
								},
							},
						},
					},
					wantErr: true,
				}
			},
		},
		{
			name: "missing REQUIRED profile",
			testCase: func(runId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				email := username + "@example.com"
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							OrganizationId: Instance.DefaultOrg.Id,
							Username:       &username,
							UserType: &user.CreateUserRequest_Human_{
								Human: &user.CreateUserRequest_Human{
									Email: &user.SetHumanEmail{
										Email: email,
										Verification: &user.SetHumanEmail_ReturnCode{
											ReturnCode: &user.ReturnEmailVerificationCode{},
										},
									},
								},
							},
						},
					},
					wantErr: true,
				}
			},
		},
		{
			name: "missing REQUIRED email",
			testCase: func(runId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							OrganizationId: Instance.DefaultOrg.Id,
							Username:       &username,
							UserType: &user.CreateUserRequest_Human_{
								Human: &user.CreateUserRequest_Human{
									Profile: &user.SetHumanProfile{
										GivenName:  "Donald",
										FamilyName: "Duck",
									},
								},
							},
						},
					},
					wantErr: true,
				}
			},
		},
		{
			name: "missing empty email",
			testCase: func(runId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							OrganizationId: Instance.DefaultOrg.Id,
							Username:       &username,
							UserType: &user.CreateUserRequest_Human_{
								Human: &user.CreateUserRequest_Human{
									Profile: &user.SetHumanProfile{
										GivenName:  "Donald",
										FamilyName: "Duck",
									},
									Email: &user.SetHumanEmail{},
								},
							},
						},
					},
					wantErr: true,
				}
			},
		},
		{
			name: "missing idp",
			testCase: func(runId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				email := username + "@example.com"
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							OrganizationId: Instance.DefaultOrg.Id,
							Username:       &username,
							UserType: &user.CreateUserRequest_Human_{
								Human: &user.CreateUserRequest_Human{
									Profile: &user.SetHumanProfile{
										GivenName:  "Donald",
										FamilyName: "Duck",
									},
									Email: &user.SetHumanEmail{
										Email: email,
										Verification: &user.SetHumanEmail_IsVerified{
											IsVerified: true,
										},
									},
									IdpLinks: []*user.IDPLink{
										{
											IdpId:    "idpID",
											UserId:   "userID",
											UserName: "username",
										},
									},
								},
							},
						},
					},
					wantErr: true,
				}
			},
		},
		{
			name: "with metadata",
			testCase: func(runId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				email := username + "@example.com"
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							OrganizationId: Instance.DefaultOrg.Id,
							Username:       &username,
							UserType: &user.CreateUserRequest_Human_{
								Human: &user.CreateUserRequest_Human{
									Profile: &user.SetHumanProfile{
										GivenName:  "Donald",
										FamilyName: "Duck",
									},
									Email: &user.SetHumanEmail{
										Email: email,
										Verification: &user.SetHumanEmail_IsVerified{
											IsVerified: true,
										},
									},
									Metadata: []*user.Metadata{
										{Key: "key1", Value: []byte(base64.StdEncoding.EncodeToString([]byte("value1")))},
										{Key: "key2", Value: []byte(base64.StdEncoding.EncodeToString([]byte("value2")))},
										{Key: "key3", Value: []byte(base64.StdEncoding.EncodeToString([]byte("value3")))},
									},
								},
							},
						},
					},
					want: &user.CreateUserResponse{
						Id: "is generated",
					},
				}
			},
		},
		{
			name: "with idp",
			testCase: func(runId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				email := username + "@example.com"
				idpResp := Instance.AddGenericOAuthProvider(IamCTX, Instance.DefaultOrg.Id)
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							OrganizationId: Instance.DefaultOrg.Id,
							Username:       &username,
							UserType: &user.CreateUserRequest_Human_{
								Human: &user.CreateUserRequest_Human{
									Profile: &user.SetHumanProfile{
										GivenName:  "Donald",
										FamilyName: "Duck",
									},
									Email: &user.SetHumanEmail{
										Email: email,
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
					want: &user.CreateUserResponse{
						Id: "is generated",
					},
				}
			},
		},
		{
			name: "with totp",
			testCase: func(runId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				email := username + "@example.com"
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							OrganizationId: Instance.DefaultOrg.Id,
							Username:       &username,
							UserType: &user.CreateUserRequest_Human_{
								Human: &user.CreateUserRequest_Human{
									Profile: &user.SetHumanProfile{
										GivenName:  "Donald",
										FamilyName: "Duck",
									},
									Email: &user.SetHumanEmail{
										Email: email,
										Verification: &user.SetHumanEmail_IsVerified{
											IsVerified: true,
										},
									},
									TotpSecret: gu.Ptr("secret"),
								},
							},
						},
					},
					want: &user.CreateUserResponse{
						Id: "is generated",
					},
				}
			},
		},
		{
			name: "password not complexity conform",
			testCase: func(runId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				email := username + "@example.com"
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							OrganizationId: Instance.DefaultOrg.Id,
							Username:       &username,
							UserType: &user.CreateUserRequest_Human_{
								Human: &user.CreateUserRequest_Human{
									Profile: &user.SetHumanProfile{
										GivenName:         "Donald",
										FamilyName:        "Duck",
										NickName:          gu.Ptr("Dukkie"),
										DisplayName:       gu.Ptr("Donald Duck"),
										PreferredLanguage: gu.Ptr("en"),
										Gender:            user.Gender_GENDER_DIVERSE.Enum(),
									},
									Email: &user.SetHumanEmail{
										Email: email,
									},
									PasswordType: &user.CreateUserRequest_Human_Password{
										Password: &user.Password{
											Password: "insufficient",
										},
									},
								},
							},
						},
					},
					wantErr: true,
				}
			},
		},
		{
			name: "hashed password",
			testCase: func(runId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				email := username + "@example.com"
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							OrganizationId: Instance.DefaultOrg.Id,
							Username:       &username,
							UserType: &user.CreateUserRequest_Human_{
								Human: &user.CreateUserRequest_Human{
									Profile: &user.SetHumanProfile{
										GivenName:  "Donald",
										FamilyName: "Duck",
									},
									Email: &user.SetHumanEmail{
										Email: email,
									},
									PasswordType: &user.CreateUserRequest_Human_HashedPassword{
										HashedPassword: &user.HashedPassword{
											Hash: "$2y$12$hXUrnqdq1RIIYZ2HPytIIe5lXdIvbhqrTvdPsSF7o.jFh817Z6lwm",
										},
									},
								},
							},
						},
					},
					want: &user.CreateUserResponse{
						Id: "is generated",
					},
				}
			},
		},
		{
			name: "unsupported hashed password",
			testCase: func(runId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				email := username + "@example.com"
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							OrganizationId: Instance.DefaultOrg.Id,
							Username:       &username,
							UserType: &user.CreateUserRequest_Human_{
								Human: &user.CreateUserRequest_Human{
									Profile: &user.SetHumanProfile{
										GivenName:  "Donald",
										FamilyName: "Duck",
									},
									Email: &user.SetHumanEmail{
										Email: email,
									},
									PasswordType: &user.CreateUserRequest_Human_HashedPassword{
										HashedPassword: &user.HashedPassword{
											Hash: "$scrypt$ln=16,r=8,p=1$cmFuZG9tc2FsdGlzaGFyZA$Rh+NnJNo1I6nRwaNqbDm6kmADswD1+7FTKZ7Ln9D8nQ",
										},
									},
								},
							},
						},
					},
					wantErr: true,
				}
			},
		},
		{
			name: "human default username",
			testCase: func(runId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				email := username + "@example.com"
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							OrganizationId: Instance.DefaultOrg.Id,
							UserType: &user.CreateUserRequest_Human_{
								Human: &user.CreateUserRequest_Human{
									Profile: &user.SetHumanProfile{
										GivenName:  "Donald",
										FamilyName: "Duck",
									},
									Email: &user.SetHumanEmail{
										Email: email,
									},
								},
							},
						},
					},
					want: &user.CreateUserResponse{
						Id: "is generated",
					},
				}
			},
		},
		{
			name: "machine user",
			testCase: func(runId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							OrganizationId: Instance.DefaultOrg.Id,
							Username:       &username,
							UserType: &user.CreateUserRequest_Machine_{
								Machine: &user.CreateUserRequest_Machine{
									Name: "donald",
								},
							},
						},
					},
					want: &user.CreateUserResponse{
						Id: "is generated",
					},
				}
			},
		},
		{
			name: "machine default username to generated id",
			testCase: func(runId string) testCase {
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							OrganizationId: Instance.DefaultOrg.Id,
							UserType: &user.CreateUserRequest_Machine_{
								Machine: &user.CreateUserRequest_Machine{
									Name: "donald",
								},
							},
						},
					},
					want: &user.CreateUserResponse{
						Id: "is generated",
					},
				}
			},
		},
		{
			name: "machine default username to given id",
			testCase: func(runId string) testCase {
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							UserId:         &runId,
							OrganizationId: Instance.DefaultOrg.Id,
							UserType: &user.CreateUserRequest_Machine_{
								Machine: &user.CreateUserRequest_Machine{
									Name: "donald",
								},
							},
						},
					},
					want: &user.CreateUserResponse{
						Id: runId,
					},
				}
			},
		},
		{
			name: "org does not exist human, error",
			testCase: func(runId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				email := username + "@example.com"
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							OrganizationId: "does not exist",
							Username:       &username,
							UserType: &user.CreateUserRequest_Human_{
								Human: &user.CreateUserRequest_Human{
									Profile: &user.SetHumanProfile{
										GivenName:  "Donald",
										FamilyName: "Duck",
									},
									Email: &user.SetHumanEmail{
										Email: email,
									},
								},
							},
						},
					},
					wantErr: true,
				}
			},
		},
		{
			name: "org does not exist machine, error",
			testCase: func(runId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				return testCase{
					args: args{
						OrgCTX,
						&user.CreateUserRequest{
							OrganizationId: "does not exist",
							Username:       &username,
							UserType: &user.CreateUserRequest_Machine_{
								Machine: &user.CreateUserRequest_Machine{
									Name: integration.Username(),
								},
							},
						},
					},
					wantErr: true,
				}
			},
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			runId := fmt.Sprint(now.UnixNano() + int64(i))
			test := tt.testCase(runId)
			got, err := Client.CreateUser(test.args.ctx, test.args.req)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			creationDate := got.CreationDate.AsTime()
			assert.Greater(t, creationDate, now, "creation date is before the test started")
			assert.Less(t, creationDate, time.Now(), "creation date is in the future")
			if test.want.GetEmailCode() != "" {
				assert.NotEmpty(t, got.GetEmailCode(), "email code is empty")
			} else {
				assert.Empty(t, got.GetEmailCode(), "email code is not empty")
			}
			if test.want.GetPhoneCode() != "" {
				assert.NotEmpty(t, got.GetPhoneCode(), "phone code is empty")
			} else {
				assert.Empty(t, got.GetPhoneCode(), "phone code is not empty")
			}
			if test.want.GetId() == "is generated" {
				assert.Len(t, got.GetId(), 18, "ID is not 18 characters")
			} else {
				assert.Equal(t, test.want.GetId(), got.GetId(), "ID is not the same")
			}
		})
	}
}

func TestServer_CreateUser_And_Compare(t *testing.T) {
	type args struct {
		ctx context.Context
		req *user.CreateUserRequest
	}
	type testCase struct {
		name   string
		args   args
		assert func(t *testing.T, createResponse *user.CreateUserResponse, getResponse *user.GetUserByIDResponse)
	}
	tests := []struct {
		name     string
		testCase func(runId string) testCase
	}{{
		name: "human given username",
		testCase: func(runId string) testCase {
			username := fmt.Sprintf("donald.duck+%s", runId)
			email := username + "@example.com"
			return testCase{
				args: args{
					ctx: OrgCTX,
					req: &user.CreateUserRequest{
						OrganizationId: Instance.DefaultOrg.Id,
						Username:       &username,
						UserType: &user.CreateUserRequest_Human_{
							Human: &user.CreateUserRequest_Human{
								Profile: &user.SetHumanProfile{
									GivenName:  "Donald",
									FamilyName: "Duck",
								},
								Email: &user.SetHumanEmail{
									Email: email,
								},
							},
						},
					},
				},
				assert: func(t *testing.T, _ *user.CreateUserResponse, getResponse *user.GetUserByIDResponse) {
					assert.Equal(t, username, getResponse.GetUser().GetUsername())
				},
			}
		},
	}, {
		name: "human username default to email",
		testCase: func(runId string) testCase {
			username := fmt.Sprintf("donald.duck+%s", runId)
			email := username + "@example.com"
			return testCase{
				args: args{
					ctx: OrgCTX,
					req: &user.CreateUserRequest{
						OrganizationId: Instance.DefaultOrg.Id,
						UserType: &user.CreateUserRequest_Human_{
							Human: &user.CreateUserRequest_Human{
								Profile: &user.SetHumanProfile{
									GivenName:  "Donald",
									FamilyName: "Duck",
								},
								Email: &user.SetHumanEmail{
									Email: email,
								},
							},
						},
					},
				},
				assert: func(t *testing.T, _ *user.CreateUserResponse, getResponse *user.GetUserByIDResponse) {
					assert.Equal(t, email, getResponse.GetUser().GetUsername())
				},
			}
		},
	}, {
		name: "machine username given",
		testCase: func(runId string) testCase {
			username := fmt.Sprintf("donald.duck+%s", runId)
			return testCase{
				args: args{
					ctx: OrgCTX,
					req: &user.CreateUserRequest{
						OrganizationId: Instance.DefaultOrg.Id,
						Username:       &username,
						UserType: &user.CreateUserRequest_Machine_{
							Machine: &user.CreateUserRequest_Machine{
								Name: "donald",
							},
						},
					},
				},
				assert: func(t *testing.T, _ *user.CreateUserResponse, getResponse *user.GetUserByIDResponse) {
					assert.Equal(t, username, getResponse.GetUser().GetUsername())
				},
			}
		},
	}, {
		name: "machine username default to generated id",
		testCase: func(runId string) testCase {
			return testCase{
				args: args{
					ctx: OrgCTX,
					req: &user.CreateUserRequest{
						OrganizationId: Instance.DefaultOrg.Id,
						UserType: &user.CreateUserRequest_Machine_{
							Machine: &user.CreateUserRequest_Machine{
								Name: "donald",
							},
						},
					},
				},
				assert: func(t *testing.T, createResponse *user.CreateUserResponse, getResponse *user.GetUserByIDResponse) {
					assert.Equal(t, createResponse.GetId(), getResponse.GetUser().GetUsername())
				},
			}
		},
	}, {
		name: "machine username default to given id",
		testCase: func(runId string) testCase {
			return testCase{
				args: args{
					ctx: OrgCTX,
					req: &user.CreateUserRequest{
						OrganizationId: Instance.DefaultOrg.Id,
						UserId:         &runId,
						UserType: &user.CreateUserRequest_Machine_{
							Machine: &user.CreateUserRequest_Machine{
								Name: "donald",
							},
						},
					},
				},
				assert: func(t *testing.T, createResponse *user.CreateUserResponse, getResponse *user.GetUserByIDResponse) {
					assert.Equal(t, runId, getResponse.GetUser().GetUsername())
				},
			}
		},
	}}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			runId := fmt.Sprint(now.UnixNano() + int64(i))
			test := tt.testCase(runId)
			createResponse, err := Client.CreateUser(test.args.ctx, test.args.req)
			require.NoError(t, err)
			Instance.TriggerUserByID(test.args.ctx, createResponse.GetId())
			getResponse, err := Client.GetUserByID(test.args.ctx, &user.GetUserByIDRequest{
				UserId: createResponse.GetId(),
			})
			require.NoError(t, err)
			test.assert(t, createResponse, getResponse)
		})
	}
}

func TestServer_CreateUser_Permission(t *testing.T) {
	newOrgOwnerEmail := integration.Email()
	newOrg := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), newOrgOwnerEmail)
	type args struct {
		ctx context.Context
		req *user.CreateUserRequest
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "human system, ok",
			args: args{
				SystemCTX,
				&user.CreateUserRequest{
					OrganizationId: newOrg.GetOrganizationId(),
					UserType: &user.CreateUserRequest_Human_{
						Human: &user.CreateUserRequest_Human{
							Profile: &user.SetHumanProfile{
								GivenName:  "Donald",
								FamilyName: "Duck",
							},
							Email: &user.SetHumanEmail{
								Email: "this is overwritten with a unique address",
							},
						},
					},
				},
			},
		},
		{
			name: "human instance, ok",
			args: args{
				IamCTX,
				&user.CreateUserRequest{
					OrganizationId: newOrg.GetOrganizationId(),
					UserType: &user.CreateUserRequest_Human_{
						Human: &user.CreateUserRequest_Human{
							Profile: &user.SetHumanProfile{
								GivenName:  "Donald",
								FamilyName: "Duck",
							},
							Email: &user.SetHumanEmail{
								Email: "this is overwritten with a unique address",
							},
						},
					},
				},
			},
		},
		{
			name: "human org, error",
			args: args{
				OrgCTX,
				&user.CreateUserRequest{
					OrganizationId: newOrg.GetOrganizationId(),
					UserType: &user.CreateUserRequest_Human_{
						Human: &user.CreateUserRequest_Human{
							Profile: &user.SetHumanProfile{
								GivenName:  "Donald",
								FamilyName: "Duck",
							},
							Email: &user.SetHumanEmail{
								Email: "this is overwritten with a unique address",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "human user, error",
			args: args{
				UserCTX,
				&user.CreateUserRequest{
					OrganizationId: newOrg.GetOrganizationId(),
					UserType: &user.CreateUserRequest_Human_{
						Human: &user.CreateUserRequest_Human{
							Profile: &user.SetHumanProfile{
								GivenName:  "Donald",
								FamilyName: "Duck",
							},
							Email: &user.SetHumanEmail{
								Email: "this is overwritten with a unique address",
							},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "machine system, ok",
			args: args{
				SystemCTX,
				&user.CreateUserRequest{
					OrganizationId: newOrg.GetOrganizationId(),
					UserType: &user.CreateUserRequest_Machine_{
						Machine: &user.CreateUserRequest_Machine{
							Name: "donald",
						},
					},
				},
			},
		},
		{
			name: "machine instance, ok",
			args: args{
				IamCTX,
				&user.CreateUserRequest{
					OrganizationId: newOrg.GetOrganizationId(),
					UserType: &user.CreateUserRequest_Machine_{
						Machine: &user.CreateUserRequest_Machine{
							Name: "donald",
						},
					},
				},
			},
		},
		{
			name: "machine org, error",
			args: args{
				OrgCTX,
				&user.CreateUserRequest{
					OrganizationId: newOrg.GetOrganizationId(),
					UserType: &user.CreateUserRequest_Machine_{
						Machine: &user.CreateUserRequest_Machine{
							Name: "donald",
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "machine user, error",
			args: args{
				UserCTX,
				&user.CreateUserRequest{
					OrganizationId: newOrg.GetOrganizationId(),
					UserType: &user.CreateUserRequest_Machine_{
						Machine: &user.CreateUserRequest_Machine{
							Name: "donald",
						},
					},
				},
			},
			wantErr: true,
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			userID := fmt.Sprint(time.Now().UnixNano() + int64(i))
			tt.args.req.UserId = &userID
			if email := tt.args.req.GetHuman().GetEmail(); email != nil {
				email.Email = fmt.Sprintf("%s@example.com", userID)
			}
			_, err := Client.CreateUser(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}

func TestServer_UpdateUserTypeHuman(t *testing.T) {
	type args struct {
		ctx context.Context
		req *user.UpdateUserRequest
	}
	type testCase struct {
		args    args
		want    *user.UpdateUserResponse
		wantErr bool
	}
	tests := []struct {
		name     string
		testCase func(runId, userId string) testCase
	}{
		{
			name: "default verification",
			testCase: func(runId, userId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				email := username + "@example.com"
				return testCase{
					args: args{
						OrgCTX,
						&user.UpdateUserRequest{
							UserId: userId,
							UserType: &user.UpdateUserRequest_Human_{
								Human: &user.UpdateUserRequest_Human{
									Email: &user.SetHumanEmail{
										Email: email,
									},
								},
							},
						},
					},
					want:    &user.UpdateUserResponse{},
					wantErr: false,
				}
			},
		},
		{
			name: "return email verification code",
			testCase: func(runId, userId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				email := username + "@example.com"
				return testCase{
					args: args{
						OrgCTX,
						&user.UpdateUserRequest{
							UserId: userId,
							UserType: &user.UpdateUserRequest_Human_{
								Human: &user.UpdateUserRequest_Human{
									Email: &user.SetHumanEmail{
										Email: email,
										Verification: &user.SetHumanEmail_ReturnCode{
											ReturnCode: &user.ReturnEmailVerificationCode{},
										},
									},
								},
							},
						},
					},
					want: &user.UpdateUserResponse{
						EmailCode: gu.Ptr("something"),
					},
				}
			},
		},
		{
			name: "custom template",
			testCase: func(runId, userId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				email := username + "@example.com"
				return testCase{
					args: args{
						OrgCTX,
						&user.UpdateUserRequest{
							UserId: userId,
							UserType: &user.UpdateUserRequest_Human_{
								Human: &user.UpdateUserRequest_Human{
									Email: &user.SetHumanEmail{
										Email: email,
										Verification: &user.SetHumanEmail_SendCode{
											SendCode: &user.SendEmailVerificationCode{
												UrlTemplate: gu.Ptr("https://example.com/email/verify?userID={{.UserID}}&code={{.Code}}&orgID={{.OrgID}}"),
											},
										},
									},
								},
							},
						},
					},
					want: &user.UpdateUserResponse{},
				}
			},
		},
		{
			name: "return phone verification code",
			testCase: func(runId, userId string) testCase {
				return testCase{
					args: args{
						OrgCTX,
						&user.UpdateUserRequest{
							UserId: userId,
							UserType: &user.UpdateUserRequest_Human_{
								Human: &user.UpdateUserRequest_Human{
									Phone: &user.SetHumanPhone{
										Phone: "+41791234568",
										Verification: &user.SetHumanPhone_ReturnCode{
											ReturnCode: &user.ReturnPhoneVerificationCode{},
										},
									},
								},
							},
						},
					},
					want: &user.UpdateUserResponse{
						PhoneCode: gu.Ptr("something"),
					},
				}
			},
		},
		{
			name: "custom template error",
			testCase: func(runId, userId string) testCase {
				username := fmt.Sprintf("donald.duck+%s", runId)
				email := username + "@example.com"
				return testCase{
					args: args{
						OrgCTX,
						&user.UpdateUserRequest{
							UserId: userId,
							UserType: &user.UpdateUserRequest_Human_{
								Human: &user.UpdateUserRequest_Human{
									Email: &user.SetHumanEmail{
										Email: email,
										Verification: &user.SetHumanEmail_SendCode{
											SendCode: &user.SendEmailVerificationCode{
												UrlTemplate: gu.Ptr("{{"),
											},
										},
									},
								},
							},
						},
					},
					wantErr: true,
				}
			},
		},
		{
			name: "missing empty email",
			testCase: func(runId, userId string) testCase {
				return testCase{
					args: args{
						OrgCTX,
						&user.UpdateUserRequest{
							UserId: userId,
							UserType: &user.UpdateUserRequest_Human_{
								Human: &user.UpdateUserRequest_Human{
									Email: &user.SetHumanEmail{},
								},
							},
						},
					},
					wantErr: true,
				}
			},
		},
		{
			name: "password not complexity conform",
			testCase: func(runId, userId string) testCase {
				return testCase{
					args: args{
						OrgCTX,
						&user.UpdateUserRequest{
							UserId: userId,
							UserType: &user.UpdateUserRequest_Human_{
								Human: &user.UpdateUserRequest_Human{
									Password: &user.SetPassword{
										PasswordType: &user.SetPassword_Password{
											Password: &user.Password{
												Password: "insufficient",
											},
										},
									},
								},
							},
						},
					},
					wantErr: true,
				}
			},
		},
		{
			name: "hashed password",
			testCase: func(runId, userId string) testCase {
				return testCase{
					args: args{
						OrgCTX,
						&user.UpdateUserRequest{
							UserId: userId,
							UserType: &user.UpdateUserRequest_Human_{
								Human: &user.UpdateUserRequest_Human{
									Password: &user.SetPassword{
										PasswordType: &user.SetPassword_HashedPassword{
											HashedPassword: &user.HashedPassword{
												Hash: "$2y$12$hXUrnqdq1RIIYZ2HPytIIe5lXdIvbhqrTvdPsSF7o.jFh817Z6lwm",
											},
										},
									},
								},
							},
						},
					},
					want: &user.UpdateUserResponse{},
				}
			},
		},
		{
			name: "unsupported hashed password",
			testCase: func(runId, userId string) testCase {
				return testCase{
					args: args{
						OrgCTX,
						&user.UpdateUserRequest{
							UserId: userId,
							UserType: &user.UpdateUserRequest_Human_{
								Human: &user.UpdateUserRequest_Human{
									Password: &user.SetPassword{
										PasswordType: &user.SetPassword_HashedPassword{
											HashedPassword: &user.HashedPassword{
												Hash: "$scrypt$ln=16,r=8,p=1$cmFuZG9tc2FsdGlzaGFyZA$Rh+NnJNo1I6nRwaNqbDm6kmADswD1+7FTKZ7Ln9D8nQ",
											},
										},
									},
								},
							},
						},
					},
					wantErr: true,
				}
			},
		},
		{
			name: "update human user with machine fields, error",
			testCase: func(runId, userId string) testCase {
				return testCase{
					args: args{
						OrgCTX,
						&user.UpdateUserRequest{
							UserId: userId,
							UserType: &user.UpdateUserRequest_Machine_{
								Machine: &user.UpdateUserRequest_Machine{
									Name: &runId,
								},
							},
						},
					},
					wantErr: true,
				}
			},
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			runId := fmt.Sprint(now.UnixNano() + int64(i))
			userId := Instance.CreateUserTypeHuman(OrgCTX, integration.Email()).GetId()
			test := tt.testCase(runId, userId)
			got, err := Client.UpdateUser(test.args.ctx, test.args.req)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			changeDate := got.ChangeDate.AsTime()
			assert.Greater(t, changeDate, now, "change date is before the test started")
			assert.Less(t, changeDate, time.Now(), "change date is in the future")
			if test.want.GetEmailCode() != "" {
				assert.NotEmpty(t, got.GetEmailCode(), "email code is empty")
			} else {
				assert.Empty(t, got.GetEmailCode(), "email code is not empty")
			}
			if test.want.GetPhoneCode() != "" {
				assert.NotEmpty(t, got.GetPhoneCode(), "phone code is empty")
			} else {
				assert.Empty(t, got.GetPhoneCode(), "phone code is not empty")
			}
		})
	}
}

func TestServer_UpdateUserTypeMachine(t *testing.T) {
	type args struct {
		ctx context.Context
		req *user.UpdateUserRequest
	}
	type testCase struct {
		args    args
		wantErr bool
	}
	tests := []struct {
		name     string
		testCase func(runId, userId string) testCase
	}{
		{
			name: "update machine, ok",
			testCase: func(runId, userId string) testCase {
				return testCase{
					args: args{
						OrgCTX,
						&user.UpdateUserRequest{
							UserId: userId,
							UserType: &user.UpdateUserRequest_Machine_{
								Machine: &user.UpdateUserRequest_Machine{
									Name: gu.Ptr("donald"),
								},
							},
						},
					},
				}
			},
		},
		{
			name: "update machine user with human fields, error",
			testCase: func(runId, userId string) testCase {
				return testCase{
					args: args{
						OrgCTX,
						&user.UpdateUserRequest{
							UserId: userId,
							UserType: &user.UpdateUserRequest_Human_{
								Human: &user.UpdateUserRequest_Human{
									Profile: &user.UpdateUserRequest_Human_Profile{
										GivenName: gu.Ptr("Donald"),
									},
								},
							},
						},
					},
					wantErr: true,
				}
			},
		},
	}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			runId := fmt.Sprint(now.UnixNano() + int64(i))
			userId := Instance.CreateUserTypeMachine(OrgCTX, Instance.DefaultOrg.Id).GetId()
			test := tt.testCase(runId, userId)
			got, err := Client.UpdateUser(test.args.ctx, test.args.req)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			changeDate := got.ChangeDate.AsTime()
			assert.Greater(t, changeDate, now, "change date is before the test started")
			assert.Less(t, changeDate, time.Now(), "change date is in the future")
		})
	}
}

func TestServer_UpdateUser_And_Compare(t *testing.T) {
	type args struct {
		ctx    context.Context
		create *user.CreateUserRequest
		update *user.UpdateUserRequest
	}
	type testCase struct {
		args   args
		assert func(t *testing.T, getResponse *user.GetUserByIDResponse)
	}
	tests := []struct {
		name     string
		testCase func(runId string) testCase
	}{{
		name: "human remove phone",
		testCase: func(runId string) testCase {
			username := fmt.Sprintf("donald.duck+%s", runId)
			email := username + "@example.com"
			return testCase{
				args: args{
					ctx: OrgCTX,
					create: &user.CreateUserRequest{
						OrganizationId: Instance.DefaultOrg.Id,
						UserId:         &runId,
						UserType: &user.CreateUserRequest_Human_{
							Human: &user.CreateUserRequest_Human{
								Profile: &user.SetHumanProfile{
									GivenName:  "Donald",
									FamilyName: "Duck",
								},
								Email: &user.SetHumanEmail{
									Email: email,
								},
								Phone: &user.SetHumanPhone{
									Phone: "+1234567890",
								},
							},
						},
					},
					update: &user.UpdateUserRequest{
						UserId: runId,
						UserType: &user.UpdateUserRequest_Human_{
							Human: &user.UpdateUserRequest_Human{
								Phone: &user.SetHumanPhone{},
							},
						},
					},
				},
				assert: func(t *testing.T, getResponse *user.GetUserByIDResponse) {
					assert.Empty(t, getResponse.GetUser().GetHuman().GetPhone().GetPhone(), "phone is not empty")
				},
			}
		},
	}, {
		name: "human username",
		testCase: func(runId string) testCase {
			username := fmt.Sprintf("donald.duck+%s", runId)
			email := username + "@example.com"
			return testCase{
				args: args{
					ctx: OrgCTX,
					create: &user.CreateUserRequest{
						OrganizationId: Instance.DefaultOrg.Id,
						UserId:         &runId,
						UserType: &user.CreateUserRequest_Human_{
							Human: &user.CreateUserRequest_Human{
								Profile: &user.SetHumanProfile{
									GivenName:  "Donald",
									FamilyName: "Duck",
								},
								Email: &user.SetHumanEmail{
									Email: email,
								},
							},
						},
					},
					update: &user.UpdateUserRequest{
						UserId:   runId,
						Username: &username,
						UserType: &user.UpdateUserRequest_Human_{
							Human: &user.UpdateUserRequest_Human{},
						},
					},
				},
				assert: func(t *testing.T, getResponse *user.GetUserByIDResponse) {
					assert.Equal(t, username, getResponse.GetUser().GetUsername())
				},
			}
		},
	}, {
		name: "machine username",
		testCase: func(runId string) testCase {
			username := fmt.Sprintf("donald.duck+%s", runId)
			return testCase{
				args: args{
					ctx: OrgCTX,
					create: &user.CreateUserRequest{
						OrganizationId: Instance.DefaultOrg.Id,
						UserId:         &runId,
						UserType: &user.CreateUserRequest_Machine_{
							Machine: &user.CreateUserRequest_Machine{
								Name: "Donald",
							},
						},
					},
					update: &user.UpdateUserRequest{
						UserId:   runId,
						Username: &username,
						UserType: &user.UpdateUserRequest_Machine_{
							Machine: &user.UpdateUserRequest_Machine{},
						},
					},
				},
				assert: func(t *testing.T, getResponse *user.GetUserByIDResponse) {
					assert.Equal(t, username, getResponse.GetUser().GetUsername())
				},
			}
		},
	}}
	for i, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			now := time.Now()
			runId := fmt.Sprint(now.UnixNano() + int64(i))
			test := tt.testCase(runId)
			createResponse, err := Client.CreateUser(test.args.ctx, test.args.create)
			require.NoError(t, err)
			_, err = Client.UpdateUser(test.args.ctx, test.args.update)
			require.NoError(t, err)
			Instance.TriggerUserByID(test.args.ctx, createResponse.GetId())
			getResponse, err := Client.GetUserByID(test.args.ctx, &user.GetUserByIDRequest{
				UserId: createResponse.GetId(),
			})
			require.NoError(t, err)
			test.assert(t, getResponse)
		})
	}
}

func TestServer_UpdateUser_Permission(t *testing.T) {
	newOrgOwnerEmail := integration.Email()
	newOrg := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), newOrgOwnerEmail)
	newHumanUserID := newOrg.CreatedAdmins[0].GetUserId()
	machineUserResp, err := Instance.Client.UserV2.CreateUser(IamCTX, &user.CreateUserRequest{
		OrganizationId: newOrg.OrganizationId,
		UserType: &user.CreateUserRequest_Machine_{
			Machine: &user.CreateUserRequest_Machine{
				Name: "Donald",
			},
		},
	})
	require.NoError(t, err)
	newMachineUserID := machineUserResp.GetId()
	Instance.TriggerUserByID(IamCTX, newMachineUserID)
	type args struct {
		ctx context.Context
		req *user.UpdateUserRequest
	}
	type testCase struct {
		args    args
		wantErr bool
	}
	tests := []struct {
		name     string
		testCase func() testCase
	}{
		{
			name: "human, system, ok",
			testCase: func() testCase {
				return testCase{
					args: args{
						SystemCTX,
						&user.UpdateUserRequest{
							UserId: newHumanUserID,
							UserType: &user.UpdateUserRequest_Human_{
								Human: &user.UpdateUserRequest_Human{
									Profile: &user.UpdateUserRequest_Human_Profile{
										GivenName: gu.Ptr("Donald"),
									},
								},
							},
						},
					},
				}
			},
		},
		{
			name: "human instance, ok",
			testCase: func() testCase {
				return testCase{
					args: args{
						IamCTX,
						&user.UpdateUserRequest{
							UserId: newHumanUserID,
							UserType: &user.UpdateUserRequest_Human_{
								Human: &user.UpdateUserRequest_Human{
									Profile: &user.UpdateUserRequest_Human_Profile{
										GivenName: gu.Ptr("Donald"),
									},
								},
							},
						},
					},
				}
			},
		},
		{
			name: "human org, error",
			testCase: func() testCase {
				return testCase{
					args: args{
						OrgCTX,
						&user.UpdateUserRequest{
							UserId: newHumanUserID,
							UserType: &user.UpdateUserRequest_Human_{
								Human: &user.UpdateUserRequest_Human{
									Profile: &user.UpdateUserRequest_Human_Profile{
										GivenName: gu.Ptr("Donald"),
									},
								},
							},
						},
					},
					wantErr: true,
				}
			},
		},
		{
			name: "human user, error",
			testCase: func() testCase {
				return testCase{
					args: args{
						UserCTX,
						&user.UpdateUserRequest{
							UserId: newHumanUserID,
							UserType: &user.UpdateUserRequest_Human_{
								Human: &user.UpdateUserRequest_Human{
									Profile: &user.UpdateUserRequest_Human_Profile{
										GivenName: gu.Ptr("Donald"),
									},
								},
							},
						},
					},
					wantErr: true,
				}
			},
		},
		{
			name: "machine system, ok",
			testCase: func() testCase {
				return testCase{
					args: args{
						SystemCTX,
						&user.UpdateUserRequest{
							UserId: newMachineUserID,
							UserType: &user.UpdateUserRequest_Machine_{
								Machine: &user.UpdateUserRequest_Machine{
									Name: gu.Ptr("Donald"),
								},
							},
						},
					},
				}
			},
		},
		{
			name: "machine instance, ok",
			testCase: func() testCase {
				return testCase{
					args: args{
						IamCTX,
						&user.UpdateUserRequest{
							UserId: newMachineUserID,
							UserType: &user.UpdateUserRequest_Machine_{
								Machine: &user.UpdateUserRequest_Machine{
									Name: gu.Ptr("Donald"),
								},
							},
						},
					},
				}
			},
		},
		{
			name: "machine org, error",
			testCase: func() testCase {
				return testCase{
					args: args{
						OrgCTX,
						&user.UpdateUserRequest{
							UserId: newMachineUserID,
							UserType: &user.UpdateUserRequest_Machine_{
								Machine: &user.UpdateUserRequest_Machine{
									Name: gu.Ptr("Donald"),
								},
							},
						},
					},
					wantErr: true,
				}
			},
		},
		{
			name: "machine user, error",
			testCase: func() testCase {
				return testCase{
					args: args{
						UserCTX,
						&user.UpdateUserRequest{
							UserId: newMachineUserID,
							UserType: &user.UpdateUserRequest_Machine_{
								Machine: &user.UpdateUserRequest_Machine{
									Name: gu.Ptr("Donald"),
								},
							},
						},
					},
					wantErr: true,
				}
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			test := tt.testCase()
			_, err := Client.UpdateUser(test.args.ctx, test.args.req)
			if test.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
		})
	}
}
