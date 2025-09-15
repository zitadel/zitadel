//go:build integration

package user_test

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc"
	"github.com/zitadel/zitadel/internal/integration"
	"github.com/zitadel/zitadel/internal/integration/sink"
	"github.com/zitadel/zitadel/pkg/grpc/idp"
	mgmt "github.com/zitadel/zitadel/pkg/grpc/management"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

var (
	CTX       context.Context
	IamCTX    context.Context
	LoginCTX  context.Context
	UserCTX   context.Context
	SystemCTX context.Context
	Instance  *integration.Instance
	Client    user.UserServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, cancel := context.WithTimeout(context.Background(), 15*time.Minute)
		defer cancel()

		Instance = integration.NewInstance(ctx)

		UserCTX = Instance.WithAuthorizationToken(ctx, integration.UserTypeNoPermission)
		IamCTX = Instance.WithAuthorizationToken(ctx, integration.UserTypeIAMOwner)
		LoginCTX = Instance.WithAuthorizationToken(ctx, integration.UserTypeLogin)
		SystemCTX = integration.WithSystemAuthorization(ctx)
		CTX = Instance.WithAuthorizationToken(ctx, integration.UserTypeOrgOwner)
		Client = Instance.Client.UserV2beta
		return m.Run()
	}())
}

func TestServer_AddHumanUser(t *testing.T) {
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
				CTX,
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
			name: "return email verification code",
			args: args{
				CTX,
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
				CTX,
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
				CTX,
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
				CTX,
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
				CTX,
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
				CTX,
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
				CTX,
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
				CTX,
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
				CTX,
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
				CTX,
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
				CTX,
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
				CTX,
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

func TestServer_AddHumanUser_Permission(t *testing.T) {
	newOrg := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
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
				CTX,
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

func TestServer_UpdateHumanUser(t *testing.T) {
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
				CTX,
				&user.UpdateHumanUserRequest{
					Username: gu.Ptr("changed"),
				},
			},
			wantErr: true,
		},
		{
			name: "change username, ok",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				userID := Instance.CreateHumanUser(CTX).GetUserId()
				request.UserId = userID
				return nil
			},
			args: args{
				CTX,
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
				userID := Instance.CreateHumanUser(CTX).GetUserId()
				request.UserId = userID
				return nil
			},
			args: args{
				CTX,
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
				userID := Instance.CreateHumanUser(CTX).GetUserId()
				request.UserId = userID
				return nil
			},
			args: args{
				CTX,
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
				userID := Instance.CreateHumanUser(CTX).GetUserId()
				request.UserId = userID
				return nil
			},
			args: args{
				CTX,
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
				userID := Instance.CreateHumanUser(CTX).GetUserId()
				request.UserId = userID
				return nil
			},
			args: args{
				CTX,
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
				userID := Instance.CreateHumanUser(CTX).GetUserId()
				request.UserId = userID
				return nil
			},
			args: args{
				CTX,
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
				userID := Instance.CreateHumanUser(CTX).GetUserId()
				request.UserId = userID
				resp, err := Client.PasswordReset(CTX, &user.PasswordResetRequest{
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
				CTX,
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
				userID := Instance.CreateHumanUser(CTX).GetUserId()
				request.UserId = userID
				resp, err := Client.PasswordReset(CTX, &user.PasswordResetRequest{
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
				CTX,
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
				userID := Instance.CreateHumanUser(CTX).GetUserId()
				request.UserId = userID
				resp, err := Client.PasswordReset(CTX, &user.PasswordResetRequest{
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
				CTX,
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
				userID := Instance.CreateHumanUser(CTX).GetUserId()
				request.UserId = userID

				resp, err := Client.PasswordReset(CTX, &user.PasswordResetRequest{
					UserId: userID,
					Medium: &user.PasswordResetRequest_ReturnCode{
						ReturnCode: &user.ReturnPasswordResetCode{},
					},
				})
				if err != nil {
					return err
				}
				pw := "Password1."
				_, err = Client.SetPassword(CTX, &user.SetPasswordRequest{
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
				CTX,
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

func TestServer_UpdateHumanUser_Permission(t *testing.T) {
	newOrg := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
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
				CTX,
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
				return
			}
			require.NoError(t, err)
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
				CTX,
				&user.LockUserRequest{
					UserId: "notexisting",
				},
				func(request *user.LockUserRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "lock, ok",
			args: args{
				CTX,
				&user.LockUserRequest{},
				func(request *user.LockUserRequest) error {
					resp := Instance.CreateHumanUser(CTX)
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
				CTX,
				&user.LockUserRequest{},
				func(request *user.LockUserRequest) error {
					resp := Instance.CreateMachineUser(CTX)
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
				CTX,
				&user.LockUserRequest{},
				func(request *user.LockUserRequest) error {
					resp := Instance.CreateHumanUser(CTX)
					request.UserId = resp.GetUserId()
					_, err := Client.LockUser(CTX, &user.LockUserRequest{
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
				CTX,
				&user.LockUserRequest{},
				func(request *user.LockUserRequest) error {
					resp := Instance.CreateMachineUser(CTX)
					request.UserId = resp.GetUserId()
					_, err := Client.LockUser(CTX, &user.LockUserRequest{
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
			} else {
				require.NoError(t, err)
			}
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
				CTX,
				&user.UnlockUserRequest{
					UserId: "notexisting",
				},
				func(request *user.UnlockUserRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "unlock, not locked",
			args: args{
				ctx: CTX,
				req: &user.UnlockUserRequest{},
				prepare: func(request *user.UnlockUserRequest) error {
					resp := Instance.CreateHumanUser(CTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "unlock machine, not locked",
			args: args{
				ctx: CTX,
				req: &user.UnlockUserRequest{},
				prepare: func(request *user.UnlockUserRequest) error {
					resp := Instance.CreateMachineUser(CTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "unlock, ok",
			args: args{
				ctx: CTX,
				req: &user.UnlockUserRequest{},
				prepare: func(request *user.UnlockUserRequest) error {
					resp := Instance.CreateHumanUser(CTX)
					request.UserId = resp.GetUserId()
					_, err := Client.LockUser(CTX, &user.LockUserRequest{
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
				ctx: CTX,
				req: &user.UnlockUserRequest{},
				prepare: func(request *user.UnlockUserRequest) error {
					resp := Instance.CreateMachineUser(CTX)
					request.UserId = resp.GetUserId()
					_, err := Client.LockUser(CTX, &user.LockUserRequest{
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
				CTX,
				&user.DeactivateUserRequest{
					UserId: "notexisting",
				},
				func(request *user.DeactivateUserRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "deactivate, ok",
			args: args{
				CTX,
				&user.DeactivateUserRequest{},
				func(request *user.DeactivateUserRequest) error {
					resp := Instance.CreateHumanUser(CTX)
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
				CTX,
				&user.DeactivateUserRequest{},
				func(request *user.DeactivateUserRequest) error {
					resp := Instance.CreateMachineUser(CTX)
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
				CTX,
				&user.DeactivateUserRequest{},
				func(request *user.DeactivateUserRequest) error {
					resp := Instance.CreateHumanUser(CTX)
					request.UserId = resp.GetUserId()
					_, err := Client.DeactivateUser(CTX, &user.DeactivateUserRequest{
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
				CTX,
				&user.DeactivateUserRequest{},
				func(request *user.DeactivateUserRequest) error {
					resp := Instance.CreateMachineUser(CTX)
					request.UserId = resp.GetUserId()
					_, err := Client.DeactivateUser(CTX, &user.DeactivateUserRequest{
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
				CTX,
				&user.ReactivateUserRequest{
					UserId: "notexisting",
				},
				func(request *user.ReactivateUserRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "reactivate, not deactivated",
			args: args{
				ctx: CTX,
				req: &user.ReactivateUserRequest{},
				prepare: func(request *user.ReactivateUserRequest) error {
					resp := Instance.CreateHumanUser(CTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "reactivate machine, not deactivated",
			args: args{
				ctx: CTX,
				req: &user.ReactivateUserRequest{},
				prepare: func(request *user.ReactivateUserRequest) error {
					resp := Instance.CreateMachineUser(CTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			wantErr: true,
		},
		{
			name: "reactivate, ok",
			args: args{
				ctx: CTX,
				req: &user.ReactivateUserRequest{},
				prepare: func(request *user.ReactivateUserRequest) error {
					resp := Instance.CreateHumanUser(CTX)
					request.UserId = resp.GetUserId()
					_, err := Client.DeactivateUser(CTX, &user.DeactivateUserRequest{
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
				ctx: CTX,
				req: &user.ReactivateUserRequest{},
				prepare: func(request *user.ReactivateUserRequest) error {
					resp := Instance.CreateMachineUser(CTX)
					request.UserId = resp.GetUserId()
					_, err := Client.DeactivateUser(CTX, &user.DeactivateUserRequest{
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
	projectResp := Instance.CreateProject(CTX, t, Instance.DefaultOrg.GetId(), integration.ProjectName(), false, false)
	type args struct {
		ctx     context.Context
		req     *user.DeleteUserRequest
		prepare func(request *user.DeleteUserRequest)
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
				CTX,
				&user.DeleteUserRequest{
					UserId: "notexisting",
				},
				nil,
			},
			wantErr: true,
		},
		{
			name: "remove human, ok",
			args: args{
				ctx: CTX,
				req: &user.DeleteUserRequest{},
				prepare: func(request *user.DeleteUserRequest) {
					resp := Instance.CreateHumanUser(CTX)
					request.UserId = resp.GetUserId()
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
				ctx: CTX,
				req: &user.DeleteUserRequest{},
				prepare: func(request *user.DeleteUserRequest) {
					resp := Instance.CreateMachineUser(CTX)
					request.UserId = resp.GetUserId()
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
				ctx: CTX,
				req: &user.DeleteUserRequest{},
				prepare: func(request *user.DeleteUserRequest) {
					resp := Instance.CreateHumanUser(CTX)
					request.UserId = resp.GetUserId()
					Instance.CreateProjectUserGrant(t, CTX, Instance.DefaultOrg.GetId(), projectResp.GetId(), request.UserId)
					Instance.CreateProjectMembership(t, CTX, projectResp.GetId(), request.UserId)
					Instance.CreateOrgMembership(t, CTX, Instance.DefaultOrg.Id, request.UserId)
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
			if tt.args.prepare != nil {
				tt.args.prepare(tt.args.req)
			}

			got, err := Client.DeleteUser(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_AddIDPLink(t *testing.T) {
	idpResp := Instance.AddGenericOAuthProvider(IamCTX, Instance.DefaultOrg.Id)
	type args struct {
		ctx context.Context
		req *user.AddIDPLinkRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.AddIDPLinkResponse
		wantErr bool
	}{
		{
			name: "user does not exist",
			args: args{
				CTX,
				&user.AddIDPLinkRequest{
					UserId: "userID",
					IdpLink: &user.IDPLink{
						IdpId:    idpResp.Id,
						UserId:   "userID",
						UserName: "username",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "idp does not exist",
			args: args{
				CTX,
				&user.AddIDPLinkRequest{
					UserId: Instance.Users.Get(integration.UserTypeOrgOwner).ID,
					IdpLink: &user.IDPLink{
						IdpId:    "idpID",
						UserId:   "userID",
						UserName: "username",
					},
				},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name: "add link",
			args: args{
				CTX,
				&user.AddIDPLinkRequest{
					UserId: Instance.Users.Get(integration.UserTypeOrgOwner).ID,
					IdpLink: &user.IDPLink{
						IdpId:    idpResp.Id,
						UserId:   "userID",
						UserName: "username",
					},
				},
			},
			want: &user.AddIDPLinkResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.DefaultOrg.Id,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.AddIDPLink(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_StartIdentityProviderIntent(t *testing.T) {
	idpResp := Instance.AddGenericOAuthProvider(IamCTX, Instance.DefaultOrg.Id)
	orgIdpID := Instance.AddOrgGenericOAuthProvider(CTX, Instance.DefaultOrg.Id)
	orgResp := Instance.CreateOrganization(IamCTX, integration.OrganizationName(), integration.Email())
	notDefaultOrgIdpID := Instance.AddOrgGenericOAuthProvider(IamCTX, orgResp.OrganizationId)
	samlIdpID := Instance.AddSAMLProvider(IamCTX)
	samlRedirectIdpID := Instance.AddSAMLRedirectProvider(IamCTX, "")
	samlPostIdpID := Instance.AddSAMLPostProvider(IamCTX)
	type args struct {
		ctx context.Context
		req *user.StartIdentityProviderIntentRequest
	}
	type want struct {
		details            *object.Details
		url                string
		parametersExisting []string
		parametersEqual    map[string]string
		postForm           bool
	}
	tests := []struct {
		name    string
		args    args
		want    want
		wantErr bool
	}{
		{
			name: "missing urls",
			args: args{
				CTX,
				&user.StartIdentityProviderIntentRequest{
					IdpId: idpResp.Id,
				},
			},
			wantErr: true,
		},
		{
			name: "next step oauth auth url",
			args: args{
				CTX,
				&user.StartIdentityProviderIntentRequest{
					IdpId: idpResp.Id,
					Content: &user.StartIdentityProviderIntentRequest_Urls{
						Urls: &user.RedirectURLs{
							SuccessUrl: "https://example.com/success",
							FailureUrl: "https://example.com/failure",
						},
					},
				},
			},
			want: want{
				details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
				url: "https://example.com/oauth/v2/authorize",
				parametersEqual: map[string]string{
					"client_id":     "clientID",
					"prompt":        "select_account",
					"redirect_uri":  "http://" + Instance.Domain + ":8080/idps/callback",
					"response_type": "code",
					"scope":         "openid profile email",
				},
				parametersExisting: []string{"state"},
			},
			wantErr: false,
		},
		{
			name: "next step oauth auth url, default org",
			args: args{
				CTX,
				&user.StartIdentityProviderIntentRequest{
					IdpId: orgIdpID.Id,
					Content: &user.StartIdentityProviderIntentRequest_Urls{
						Urls: &user.RedirectURLs{
							SuccessUrl: "https://example.com/success",
							FailureUrl: "https://example.com/failure",
						},
					},
				},
			},
			want: want{
				details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
				url: "https://example.com/oauth/v2/authorize",
				parametersEqual: map[string]string{
					"client_id":     "clientID",
					"prompt":        "select_account",
					"redirect_uri":  "http://" + Instance.Domain + ":8080/idps/callback",
					"response_type": "code",
					"scope":         "openid profile email",
				},
				parametersExisting: []string{"state"},
			},
			wantErr: false,
		},
		{
			name: "next step oauth auth url, default org",
			args: args{
				CTX,
				&user.StartIdentityProviderIntentRequest{
					IdpId: notDefaultOrgIdpID.Id,
					Content: &user.StartIdentityProviderIntentRequest_Urls{
						Urls: &user.RedirectURLs{
							SuccessUrl: "https://example.com/success",
							FailureUrl: "https://example.com/failure",
						},
					},
				},
			},
			want: want{
				details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
				url: "https://example.com/oauth/v2/authorize",
				parametersEqual: map[string]string{
					"client_id":     "clientID",
					"prompt":        "select_account",
					"redirect_uri":  "http://" + Instance.Domain + ":8080/idps/callback",
					"response_type": "code",
					"scope":         "openid profile email",
				},
				parametersExisting: []string{"state"},
			},
			wantErr: false,
		},
		{
			name: "next step oauth auth url org",
			args: args{
				CTX,
				&user.StartIdentityProviderIntentRequest{
					IdpId: orgIdpID.Id,
					Content: &user.StartIdentityProviderIntentRequest_Urls{
						Urls: &user.RedirectURLs{
							SuccessUrl: "https://example.com/success",
							FailureUrl: "https://example.com/failure",
						},
					},
				},
			},
			want: want{
				details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
				url: "https://example.com/oauth/v2/authorize",
				parametersEqual: map[string]string{
					"client_id":     "clientID",
					"prompt":        "select_account",
					"redirect_uri":  "http://" + Instance.Domain + ":8080/idps/callback",
					"response_type": "code",
					"scope":         "openid profile email",
				},
				parametersExisting: []string{"state"},
			},
			wantErr: false,
		},
		{
			name: "next step saml default",
			args: args{
				CTX,
				&user.StartIdentityProviderIntentRequest{
					IdpId: samlIdpID,
					Content: &user.StartIdentityProviderIntentRequest_Urls{
						Urls: &user.RedirectURLs{
							SuccessUrl: "https://example.com/success",
							FailureUrl: "https://example.com/failure",
						},
					},
				},
			},
			want: want{
				details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
				url:                "http://localhost:8000/sso",
				parametersExisting: []string{"RelayState", "SAMLRequest"},
			},
			wantErr: false,
		},
		{
			name: "next step saml auth url",
			args: args{
				CTX,
				&user.StartIdentityProviderIntentRequest{
					IdpId: samlRedirectIdpID,
					Content: &user.StartIdentityProviderIntentRequest_Urls{
						Urls: &user.RedirectURLs{
							SuccessUrl: "https://example.com/success",
							FailureUrl: "https://example.com/failure",
						},
					},
				},
			},
			want: want{
				details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
				url:                "http://localhost:8000/sso",
				parametersExisting: []string{"RelayState", "SAMLRequest"},
			},
			wantErr: false,
		},
		{
			name: "next step saml form",
			args: args{
				CTX,
				&user.StartIdentityProviderIntentRequest{
					IdpId: samlPostIdpID,
					Content: &user.StartIdentityProviderIntentRequest_Urls{
						Urls: &user.RedirectURLs{
							SuccessUrl: "https://example.com/success",
							FailureUrl: "https://example.com/failure",
						},
					},
				},
			},
			want: want{
				details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Instance.ID(),
				},
				url:                "http://localhost:8000/sso",
				parametersExisting: []string{"RelayState", "SAMLRequest"},
				postForm:           true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.StartIdentityProviderIntent(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
				return
			}
			require.NoError(t, err)

			if tt.want.url != "" && !tt.want.postForm {
				authUrl, err := url.Parse(got.GetAuthUrl())
				require.NoError(t, err)

				assert.Equal(t, tt.want.url, authUrl.Scheme+"://"+authUrl.Host+authUrl.Path)
				require.Len(t, authUrl.Query(), len(tt.want.parametersEqual)+len(tt.want.parametersExisting))

				for _, existing := range tt.want.parametersExisting {
					assert.True(t, authUrl.Query().Has(existing))
				}
				for key, equal := range tt.want.parametersEqual {
					assert.Equal(t, equal, authUrl.Query().Get(key))
				}
			}
			if tt.want.postForm {
				assert.Equal(t, tt.want.url, got.GetFormData().GetUrl())

				require.Len(t, got.GetFormData().GetFields(), len(tt.want.parametersEqual)+len(tt.want.parametersExisting))
				for _, existing := range tt.want.parametersExisting {
					assert.Contains(t, got.GetFormData().GetFields(), existing)
				}
				for key, equal := range tt.want.parametersEqual {
					assert.Equal(t, got.GetFormData().GetFields()[key], equal)
				}
			}
			integration.AssertDetails(t, &user.StartIdentityProviderIntentResponse{
				Details: tt.want.details,
			}, got)
		})
	}
}

func TestServer_RetrieveIdentityProviderIntent(t *testing.T) {
	oauthIdpID := Instance.AddGenericOAuthProvider(IamCTX, integration.IDPName()).GetId()
	oidcIdpID := Instance.AddGenericOIDCProvider(IamCTX, integration.IDPName()).GetId()
	samlIdpID := Instance.AddSAMLPostProvider(IamCTX)
	ldapIdpID := Instance.AddLDAPProvider(IamCTX)
	authURL, err := url.Parse(Instance.CreateIntent(CTX, oauthIdpID).GetAuthUrl())
	require.NoError(t, err)
	intentID := authURL.Query().Get("state")
	expiry := time.Now().Add(1 * time.Hour)
	expiryFormatted := expiry.Round(time.Millisecond).UTC().Format("2006-01-02T15:04:05.999Z07:00")

	intentUser := Instance.CreateHumanUser(IamCTX)
	_, err = Instance.CreateUserIDPlink(IamCTX, intentUser.GetUserId(), "idpUserID", oauthIdpID, "username")
	require.NoError(t, err)

	successfulID, token, changeDate, sequence, err := sink.SuccessfulOAuthIntent(Instance.ID(), oauthIdpID, "id", "", expiry)
	require.NoError(t, err)
	successfulWithUserID, withUsertoken, withUserchangeDate, withUsersequence, err := sink.SuccessfulOAuthIntent(Instance.ID(), oauthIdpID, "id", "user", expiry)
	require.NoError(t, err)
	successfulExpiredID, expiredToken, _, _, err := sink.SuccessfulOAuthIntent(Instance.ID(), oauthIdpID, "id", "user", time.Now().Add(time.Second))
	require.NoError(t, err)
	// make sure the intent is expired
	time.Sleep(2 * time.Second)
	successfulConsumedID, consumedToken, _, _, err := sink.SuccessfulOAuthIntent(Instance.ID(), oauthIdpID, "idpUserID", intentUser.GetUserId(), expiry)
	require.NoError(t, err)
	// make sure the intent is consumed
	Instance.CreateIntentSession(t, IamCTX, intentUser.GetUserId(), successfulConsumedID, consumedToken)
	oidcSuccessful, oidcToken, oidcChangeDate, oidcSequence, err := sink.SuccessfulOIDCIntent(Instance.ID(), oidcIdpID, "id", "", expiry)
	require.NoError(t, err)
	oidcSuccessfulWithUserID, oidcWithUserIDToken, oidcWithUserIDChangeDate, oidcWithUserIDSequence, err := sink.SuccessfulOIDCIntent(Instance.ID(), oidcIdpID, "id", "user", expiry)
	require.NoError(t, err)
	ldapSuccessfulID, ldapToken, ldapChangeDate, ldapSequence, err := sink.SuccessfulLDAPIntent(Instance.ID(), ldapIdpID, "id", "")
	require.NoError(t, err)
	ldapSuccessfulWithUserID, ldapWithUserToken, ldapWithUserChangeDate, ldapWithUserSequence, err := sink.SuccessfulLDAPIntent(Instance.ID(), ldapIdpID, "id", "user")
	require.NoError(t, err)
	samlSuccessfulID, samlToken, samlChangeDate, samlSequence, err := sink.SuccessfulSAMLIntent(Instance.ID(), samlIdpID, "id", "", expiry)
	require.NoError(t, err)
	samlSuccessfulWithUserID, samlWithUserToken, samlWithUserChangeDate, samlWithUserSequence, err := sink.SuccessfulSAMLIntent(Instance.ID(), samlIdpID, "id", "user", expiry)
	require.NoError(t, err)
	type args struct {
		ctx context.Context
		req *user.RetrieveIdentityProviderIntentRequest
	}
	tests := []struct {
		name    string
		args    args
		want    *user.RetrieveIdentityProviderIntentResponse
		wantErr bool
	}{
		{
			name: "failed intent",
			args: args{
				CTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    intentID,
					IdpIntentToken: "",
				},
			},
			wantErr: true,
		},
		{
			name: "wrong token",
			args: args{
				CTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    successfulID,
					IdpIntentToken: "wrong token",
				},
			},
			wantErr: true,
		},
		{
			name: "retrieve successful oauth intent",
			args: args{
				CTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    successfulID,
					IdpIntentToken: token,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(changeDate),
					ResourceOwner: Instance.ID(),
					Sequence:      sequence,
				},
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Oauth{
						Oauth: &user.IDPOAuthAccessInformation{
							AccessToken: "accessToken",
							IdToken:     gu.Ptr("idToken"),
						},
					},
					IdpId:    oauthIdpID,
					UserId:   "id",
					UserName: "",
					RawInformation: func() *structpb.Struct {
						s, err := structpb.NewStruct(map[string]interface{}{
							"RawInfo": map[string]interface{}{
								"id":                 "id",
								"preferred_username": "username",
							},
						})
						require.NoError(t, err)
						return s
					}(),
				},
			},
			wantErr: false,
		},
		{
			name: "retrieve successful intent with linked user",
			args: args{
				CTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    successfulWithUserID,
					IdpIntentToken: withUsertoken,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(withUserchangeDate),
					ResourceOwner: Instance.ID(),
					Sequence:      withUsersequence,
				},
				UserId: "user",
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Oauth{
						Oauth: &user.IDPOAuthAccessInformation{
							AccessToken: "accessToken",
							IdToken:     gu.Ptr("idToken"),
						},
					},
					IdpId:    oauthIdpID,
					UserId:   "id",
					UserName: "",
					RawInformation: func() *structpb.Struct {
						s, err := structpb.NewStruct(map[string]interface{}{
							"RawInfo": map[string]interface{}{
								"id":                 "id",
								"preferred_username": "username",
							},
						})
						require.NoError(t, err)
						return s
					}(),
				},
			},
			wantErr: false,
		},
		{
			name: "retrieve successful expired intent",
			args: args{
				CTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    successfulExpiredID,
					IdpIntentToken: expiredToken,
				},
			},
			wantErr: true,
		},
		{
			name: "retrieve successful consumed intent",
			args: args{
				CTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    successfulConsumedID,
					IdpIntentToken: consumedToken,
				},
			},
			wantErr: true,
		},
		{
			name: "retrieve successful oidc intent",
			args: args{
				CTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    oidcSuccessful,
					IdpIntentToken: oidcToken,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(oidcChangeDate),
					ResourceOwner: Instance.ID(),
					Sequence:      oidcSequence,
				},
				UserId: "",
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Oauth{
						Oauth: &user.IDPOAuthAccessInformation{
							AccessToken: "accessToken",
							IdToken:     gu.Ptr("idToken"),
						},
					},
					IdpId:    oidcIdpID,
					UserId:   "id",
					UserName: "username",
					RawInformation: func() *structpb.Struct {
						s, err := structpb.NewStruct(map[string]interface{}{
							"sub":                "id",
							"preferred_username": "username",
						})
						require.NoError(t, err)
						return s
					}(),
				},
			},
			wantErr: false,
		},
		{
			name: "retrieve successful oidc intent with linked user",
			args: args{
				CTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    oidcSuccessfulWithUserID,
					IdpIntentToken: oidcWithUserIDToken,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(oidcWithUserIDChangeDate),
					ResourceOwner: Instance.ID(),
					Sequence:      oidcWithUserIDSequence,
				},
				UserId: "user",
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Oauth{
						Oauth: &user.IDPOAuthAccessInformation{
							AccessToken: "accessToken",
							IdToken:     gu.Ptr("idToken"),
						},
					},
					IdpId:    oidcIdpID,
					UserId:   "id",
					UserName: "username",
					RawInformation: func() *structpb.Struct {
						s, err := structpb.NewStruct(map[string]interface{}{
							"sub":                "id",
							"preferred_username": "username",
						})
						require.NoError(t, err)
						return s
					}(),
				},
			},
			wantErr: false,
		},
		{
			name: "retrieve successful ldap intent",
			args: args{
				CTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    ldapSuccessfulID,
					IdpIntentToken: ldapToken,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(ldapChangeDate),
					ResourceOwner: Instance.ID(),
					Sequence:      ldapSequence,
				},
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Ldap{
						Ldap: &user.IDPLDAPAccessInformation{
							Attributes: func() *structpb.Struct {
								s, err := structpb.NewStruct(map[string]interface{}{
									"id":       []interface{}{"id"},
									"username": []interface{}{"username"},
									"language": []interface{}{"en"},
								})
								require.NoError(t, err)
								return s
							}(),
						},
					},
					IdpId:    ldapIdpID,
					UserId:   "id",
					UserName: "username",
					RawInformation: func() *structpb.Struct {
						s, err := structpb.NewStruct(map[string]interface{}{
							"id":                "id",
							"preferredUsername": "username",
							"preferredLanguage": "en",
						})
						require.NoError(t, err)
						return s
					}(),
				},
			},
			wantErr: false,
		},
		{
			name: "retrieve successful ldap intent with linked user",
			args: args{
				CTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    ldapSuccessfulWithUserID,
					IdpIntentToken: ldapWithUserToken,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(ldapWithUserChangeDate),
					ResourceOwner: Instance.ID(),
					Sequence:      ldapWithUserSequence,
				},
				UserId: "user",
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Ldap{
						Ldap: &user.IDPLDAPAccessInformation{
							Attributes: func() *structpb.Struct {
								s, err := structpb.NewStruct(map[string]interface{}{
									"id":       []interface{}{"id"},
									"username": []interface{}{"username"},
									"language": []interface{}{"en"},
								})
								require.NoError(t, err)
								return s
							}(),
						},
					},
					IdpId:    ldapIdpID,
					UserId:   "id",
					UserName: "username",
					RawInformation: func() *structpb.Struct {
						s, err := structpb.NewStruct(map[string]interface{}{
							"id":                "id",
							"preferredUsername": "username",
							"preferredLanguage": "en",
						})
						require.NoError(t, err)
						return s
					}(),
				},
			},
			wantErr: false,
		},
		{
			name: "retrieve successful saml intent",
			args: args{
				CTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    samlSuccessfulID,
					IdpIntentToken: samlToken,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(samlChangeDate),
					ResourceOwner: Instance.ID(),
					Sequence:      samlSequence,
				},
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Saml{
						Saml: &user.IDPSAMLAccessInformation{
							Assertion: []byte(fmt.Sprintf(`<Assertion xmlns="urn:oasis:names:tc:SAML:2.0:assertion" ID="id" IssueInstant="0001-01-01T00:00:00Z" Version=""><Issuer xmlns="urn:oasis:names:tc:SAML:2.0:assertion" NameQualifier="" SPNameQualifier="" Format="" SPProvidedID=""></Issuer><Conditions NotBefore="0001-01-01T00:00:00Z" NotOnOrAfter="%s"></Conditions></Assertion>`, expiryFormatted)),
						},
					},
					IdpId:    samlIdpID,
					UserId:   "id",
					UserName: "",
					RawInformation: func() *structpb.Struct {
						s, err := structpb.NewStruct(map[string]interface{}{
							"id": "id",
							"attributes": map[string]interface{}{
								"attribute1": []interface{}{"value1"},
							},
						})
						require.NoError(t, err)
						return s
					}(),
				},
			},
			wantErr: false,
		},
		{
			name: "retrieve successful saml intent with linked user",
			args: args{
				CTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    samlSuccessfulWithUserID,
					IdpIntentToken: samlWithUserToken,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(samlWithUserChangeDate),
					ResourceOwner: Instance.ID(),
					Sequence:      samlWithUserSequence,
				},
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Saml{
						Saml: &user.IDPSAMLAccessInformation{
							Assertion: []byte(fmt.Sprintf(`<Assertion xmlns="urn:oasis:names:tc:SAML:2.0:assertion" ID="id" IssueInstant="0001-01-01T00:00:00Z" Version=""><Issuer xmlns="urn:oasis:names:tc:SAML:2.0:assertion" NameQualifier="" SPNameQualifier="" Format="" SPProvidedID=""></Issuer><Conditions NotBefore="0001-01-01T00:00:00Z" NotOnOrAfter="%s"></Conditions></Assertion>`, expiryFormatted)),
						},
					},
					IdpId:    samlIdpID,
					UserId:   "id",
					UserName: "",
					RawInformation: func() *structpb.Struct {
						s, err := structpb.NewStruct(map[string]interface{}{
							"id": "id",
							"attributes": map[string]interface{}{
								"attribute1": []interface{}{"value1"},
							},
						})
						require.NoError(t, err)
						return s
					}(),
				},
				UserId: "user",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.RetrieveIdentityProviderIntent(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			grpc.AllFieldsEqual(t, tt.want.ProtoReflect(), got.ProtoReflect(), grpc.CustomMappers)
		})
	}
}

func TestServer_ListAuthenticationMethodTypes(t *testing.T) {
	userIDWithoutAuth := Instance.CreateHumanUser(CTX).GetUserId()

	userIDWithPasskey := Instance.CreateHumanUser(CTX).GetUserId()
	Instance.RegisterUserPasskey(CTX, userIDWithPasskey)

	userMultipleAuth := Instance.CreateHumanUser(CTX).GetUserId()
	Instance.RegisterUserPasskey(CTX, userMultipleAuth)
	provider, err := Instance.Client.Mgmt.AddGenericOIDCProvider(CTX, &mgmt.AddGenericOIDCProviderRequest{
		Name:         "ListAuthenticationMethodTypes",
		Issuer:       "https://example.com",
		ClientId:     "client_id",
		ClientSecret: "client_secret",
	})
	require.NoError(t, err)
	_, err = Instance.Client.Mgmt.AddCustomLoginPolicy(CTX, &mgmt.AddCustomLoginPolicyRequest{})
	require.Condition(t, func() bool {
		code := status.Convert(err).Code()
		return code == codes.AlreadyExists || code == codes.OK
	})
	_, err = Instance.Client.Mgmt.AddIDPToLoginPolicy(CTX, &mgmt.AddIDPToLoginPolicyRequest{
		IdpId:     provider.GetId(),
		OwnerType: idp.IDPOwnerType_IDP_OWNER_TYPE_ORG,
	})
	require.NoError(t, err)
	_, err = Client.AddIDPLink(CTX, &user.AddIDPLinkRequest{UserId: userMultipleAuth, IdpLink: &user.IDPLink{
		IdpId:    provider.GetId(),
		UserId:   "external-id",
		UserName: "displayName",
	}})
	require.NoError(t, err)
	// This should not remove the user IDP links
	_, err = Instance.Client.Mgmt.RemoveIDPFromLoginPolicy(CTX, &mgmt.RemoveIDPFromLoginPolicyRequest{
		IdpId: provider.GetId(),
	})
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
				CTX,
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
				CTX,
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
			name: "multiple auth",
			args: args{
				CTX,
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
