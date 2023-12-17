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
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/api/grpc"
	"github.com/zitadel/zitadel/internal/integration"
	mgmt "github.com/zitadel/zitadel/pkg/grpc/management"
	object "github.com/zitadel/zitadel/pkg/grpc/object/v2beta"
	user "github.com/zitadel/zitadel/pkg/grpc/user/v2beta"
)

var (
	CTX    context.Context
	ErrCTX context.Context
	Tester *integration.Tester
	Client user.UserServiceClient
)

func TestMain(m *testing.M) {
	os.Exit(func() int {
		ctx, errCtx, cancel := integration.Contexts(time.Hour)
		defer cancel()

		Tester = integration.NewTester(ctx)
		defer Tester.Done()

		CTX, ErrCTX = Tester.WithAuthorization(ctx, integration.OrgOwner), errCtx
		Client = Tester.Client.UserV2
		return m.Run()
	}())
}

func TestServer_AddHumanUser(t *testing.T) {
	idpID := Tester.AddGenericOAuthProvider(t)
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
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
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
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "return email verification code",
			args: args{
				CTX,
				&user.AddHumanUserRequest{
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
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
					ResourceOwner: Tester.Organisation.ID,
				},
				EmailCode: gu.Ptr("something"),
			},
		},
		{
			name: "custom template",
			args: args{
				CTX,
				&user.AddHumanUserRequest{
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
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
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "return phone verification code",
			args: args{
				CTX,
				&user.AddHumanUserRequest{
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
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
					ResourceOwner: Tester.Organisation.ID,
				},
				PhoneCode: gu.Ptr("something"),
			},
		},
		{
			name: "custom template error",
			args: args{
				CTX,
				&user.AddHumanUserRequest{
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
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
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
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
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
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
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
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
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
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
							IdpId:    idpID,
							UserId:   "userID",
							UserName: "username",
						},
					},
				},
			},
			want: &user.AddHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "hashed password",
			args: args{
				CTX,
				&user.AddHumanUserRequest{
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
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
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "unsupported hashed password",
			args: args{
				CTX,
				&user.AddHumanUserRequest{
					Organisation: &object.Organisation{
						Org: &object.Organisation_OrgId{
							OrgId: Tester.Organisation.ID,
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
			} else {
				require.NoError(t, err)
			}

			assert.Equal(t, tt.want.GetUserId(), got.GetUserId())
			if tt.want.GetEmailCode() != "" {
				assert.NotEmpty(t, got.GetEmailCode())
			}
			if tt.want.GetPhoneCode() != "" {
				assert.NotEmpty(t, got.GetPhoneCode())
			}
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
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
					Username: gu.Ptr("changed"),
				},
			},
			wantErr: true,
		},
		{
			name: "change username, ok",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				userID := Tester.CreateHumanUser(CTX).GetUserId()
				request.UserId = userID
				return nil
			},
			args: args{
				CTX,
				&user.UpdateHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
					Username: gu.Ptr(fmt.Sprint(time.Now().UnixNano() + 1)),
				},
			},
			want: &user.UpdateHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "change profile, ok",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				userID := Tester.CreateHumanUser(CTX).GetUserId()
				request.UserId = userID
				return nil
			},
			args: args{
				CTX,
				&user.UpdateHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
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
				},
			},
			want: &user.UpdateHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "change email, ok",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				userID := Tester.CreateHumanUser(CTX).GetUserId()
				request.UserId = userID
				return nil
			},
			args: args{
				CTX,
				&user.UpdateHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
					Email: &user.SetHumanEmail{
						Email:        "changed@test.com",
						Verification: &user.SetHumanEmail_IsVerified{IsVerified: true},
					},
				},
			},
			want: &user.UpdateHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "change email, code, ok",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				userID := Tester.CreateHumanUser(CTX).GetUserId()
				request.UserId = userID
				return nil
			},
			args: args{
				CTX,
				&user.UpdateHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
					Email: &user.SetHumanEmail{
						Email:        "changed@test.com",
						Verification: &user.SetHumanEmail_ReturnCode{},
					},
				},
			},
			want: &user.UpdateHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
				EmailCode: gu.Ptr("something"),
			},
		},
		{
			name: "change phone, ok",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				userID := Tester.CreateHumanUser(CTX).GetUserId()
				request.UserId = userID
				return nil
			},
			args: args{
				CTX,
				&user.UpdateHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
					Phone: &user.SetHumanPhone{
						Phone:        "+41791234567",
						Verification: &user.SetHumanPhone_IsVerified{IsVerified: true},
					},
				},
			},
			want: &user.UpdateHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "change phone, code, ok",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				userID := Tester.CreateHumanUser(CTX).GetUserId()
				request.UserId = userID
				return nil
			},
			args: args{
				CTX,
				&user.UpdateHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
					Phone: &user.SetHumanPhone{
						Phone:        "+41791234568",
						Verification: &user.SetHumanPhone_ReturnCode{},
					},
				},
			},
			want: &user.UpdateHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
				PhoneCode: gu.Ptr("something"),
			},
		},
		{
			name: "change password, code, ok",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				userID := Tester.CreateHumanUser(CTX).GetUserId()
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
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
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
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "change hashed password, code, ok",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				userID := Tester.CreateHumanUser(CTX).GetUserId()
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
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
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
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "change hashed password, code, not supported",
			prepare: func(request *user.UpdateHumanUserRequest) error {
				userID := Tester.CreateHumanUser(CTX).GetUserId()
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
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
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
				userID := Tester.CreateHumanUser(CTX).GetUserId()
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
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
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
					ResourceOwner: Tester.Organisation.ID,
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
			} else {
				require.NoError(t, err)
			}
			if tt.want.GetEmailCode() != "" {
				assert.NotEmpty(t, got.GetEmailCode())
			}
			if tt.want.GetPhoneCode() != "" {
				assert.NotEmpty(t, got.GetPhoneCode())
			}
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_LockHumanUser(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     *user.LockHumanUserRequest
		prepare func(request *user.LockHumanUserRequest) error
	}
	tests := []struct {
		name    string
		args    args
		want    *user.LockHumanUserResponse
		wantErr bool
	}{
		{
			name: "lock, not existing",
			args: args{
				CTX,
				&user.LockHumanUserRequest{
					UserId: "notexisting",
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
				},
				func(request *user.LockHumanUserRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "lock, ok",
			args: args{
				CTX,
				&user.LockHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
				},
				func(request *user.LockHumanUserRequest) error {
					resp := Tester.CreateHumanUser(CTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			want: &user.LockHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "lock, already locked",
			args: args{
				CTX,
				&user.LockHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
				},
				func(request *user.LockHumanUserRequest) error {
					resp := Tester.CreateHumanUser(CTX)
					request.UserId = resp.GetUserId()
					_, err := Client.LockHumanUser(CTX, &user.LockHumanUserRequest{
						UserId: resp.GetUserId(),
						Organization: &object.Organization{
							Org: &object.Organization_OrgId{
								OrgId: Tester.Organisation.ID,
							},
						},
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

			got, err := Client.LockHumanUser(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_UnlockHumanUser(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     *user.UnlockHumanUserRequest
		prepare func(request *user.UnlockHumanUserRequest) error
	}
	tests := []struct {
		name    string
		args    args
		want    *user.UnlockHumanUserResponse
		wantErr bool
	}{
		{
			name: "unlock, not existing",
			args: args{
				CTX,
				&user.UnlockHumanUserRequest{
					UserId: "notexisting",
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
				},
				func(request *user.UnlockHumanUserRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "unlock, not locked",
			args: args{
				ctx: CTX,
				req: &user.UnlockHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
				},
				prepare: func(request *user.UnlockHumanUserRequest) error {
					resp := Tester.CreateHumanUser(CTX)
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
				req: &user.UnlockHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
				},
				prepare: func(request *user.UnlockHumanUserRequest) error {
					resp := Tester.CreateHumanUser(CTX)
					request.UserId = resp.GetUserId()
					_, err := Client.LockHumanUser(CTX, &user.LockHumanUserRequest{
						UserId: resp.GetUserId(),
						Organization: &object.Organization{
							Org: &object.Organization_OrgId{
								OrgId: Tester.Organisation.ID,
							},
						},
					})
					return err
				},
			},
			want: &user.UnlockHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.prepare(tt.args.req)
			require.NoError(t, err)

			got, err := Client.UnlockHumanUser(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_DeactivateHumanUser(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     *user.DeactivateHumanUserRequest
		prepare func(request *user.DeactivateHumanUserRequest) error
	}
	tests := []struct {
		name    string
		args    args
		want    *user.DeactivateHumanUserResponse
		wantErr bool
	}{
		{
			name: "deactivate, not existing",
			args: args{
				CTX,
				&user.DeactivateHumanUserRequest{
					UserId: "notexisting",
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
				},
				func(request *user.DeactivateHumanUserRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "deactivate, ok",
			args: args{
				CTX,
				&user.DeactivateHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
				},
				func(request *user.DeactivateHumanUserRequest) error {
					resp := Tester.CreateHumanUser(CTX)
					request.UserId = resp.GetUserId()
					return nil
				},
			},
			want: &user.DeactivateHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "deactivate, already deactivated",
			args: args{
				CTX,
				&user.DeactivateHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
				},
				func(request *user.DeactivateHumanUserRequest) error {
					resp := Tester.CreateHumanUser(CTX)
					request.UserId = resp.GetUserId()
					_, err := Client.DeactivateHumanUser(CTX, &user.DeactivateHumanUserRequest{
						UserId: resp.GetUserId(),
						Organization: &object.Organization{
							Org: &object.Organization_OrgId{
								OrgId: Tester.Organisation.ID,
							},
						},
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

			got, err := Client.DeactivateHumanUser(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_ReactivateHumanUser(t *testing.T) {
	type args struct {
		ctx     context.Context
		req     *user.ReactivateHumanUserRequest
		prepare func(request *user.ReactivateHumanUserRequest) error
	}
	tests := []struct {
		name    string
		args    args
		want    *user.ReactivateHumanUserResponse
		wantErr bool
	}{
		{
			name: "reactivate, not existing",
			args: args{
				CTX,
				&user.ReactivateHumanUserRequest{
					UserId: "notexisting",
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
				},
				func(request *user.ReactivateHumanUserRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "reactivate, not deactivated",
			args: args{
				ctx: CTX,
				req: &user.ReactivateHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
				},
				prepare: func(request *user.ReactivateHumanUserRequest) error {
					resp := Tester.CreateHumanUser(CTX)
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
				req: &user.ReactivateHumanUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
				},
				prepare: func(request *user.ReactivateHumanUserRequest) error {
					resp := Tester.CreateHumanUser(CTX)
					request.UserId = resp.GetUserId()
					_, err := Client.DeactivateHumanUser(CTX, &user.DeactivateHumanUserRequest{
						UserId: resp.GetUserId(),
						Organization: &object.Organization{
							Org: &object.Organization_OrgId{
								OrgId: Tester.Organisation.ID,
							},
						},
					})
					return err
				},
			},
			want: &user.ReactivateHumanUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.prepare(tt.args.req)
			require.NoError(t, err)

			got, err := Client.ReactivateHumanUser(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_DeleteUser(t *testing.T) {
	projectResp, err := Tester.CreateProject(CTX)
	require.NoError(t, err)
	type args struct {
		ctx     context.Context
		req     *user.DeleteUserRequest
		prepare func(request *user.DeleteUserRequest) error
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
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
				},
				func(request *user.DeleteUserRequest) error { return nil },
			},
			wantErr: true,
		},
		{
			name: "remove, ok",
			args: args{
				ctx: CTX,
				req: &user.DeleteUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
				},
				prepare: func(request *user.DeleteUserRequest) error {
					resp := Tester.CreateHumanUser(CTX)
					request.UserId = resp.GetUserId()
					return err
				},
			},
			want: &user.DeleteUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
		{
			name: "remove dependencies, ok",
			args: args{
				ctx: CTX,
				req: &user.DeleteUserRequest{
					Organization: &object.Organization{
						Org: &object.Organization_OrgId{
							OrgId: Tester.Organisation.ID,
						},
					},
				},
				prepare: func(request *user.DeleteUserRequest) error {
					resp := Tester.CreateHumanUser(CTX)
					request.UserId = resp.GetUserId()
					Tester.CreateProjectUserGrant(t, CTX, projectResp.GetId(), request.UserId)
					Tester.CreateProjectMembership(t, CTX, projectResp.GetId(), request.UserId)
					Tester.CreateOrgMembership(t, CTX, request.UserId)
					return err
				},
			},
			want: &user.DeleteUserResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.args.prepare(tt.args.req)
			require.NoError(t, err)

			got, err := Client.DeleteUser(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_AddIDPLink(t *testing.T) {
	idpID := Tester.AddGenericOAuthProvider(t)
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
						IdpId:    idpID,
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
					UserId: Tester.Users[integration.FirstInstanceUsersKey][integration.OrgOwner].ID,
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
					UserId: Tester.Users[integration.FirstInstanceUsersKey][integration.OrgOwner].ID,
					IdpLink: &user.IDPLink{
						IdpId:    idpID,
						UserId:   "userID",
						UserName: "username",
					},
				},
			},
			want: &user.AddIDPLinkResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.Now(),
					ResourceOwner: Tester.Organisation.ID,
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
			} else {
				require.NoError(t, err)
			}

			integration.AssertDetails(t, tt.want, got)
		})
	}
}

func TestServer_StartIdentityProviderIntent(t *testing.T) {
	idpID := Tester.AddGenericOAuthProvider(t)
	samlIdpID := Tester.AddSAMLProvider(t)
	samlRedirectIdpID := Tester.AddSAMLRedirectProvider(t)
	samlPostIdpID := Tester.AddSAMLPostProvider(t)
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
					IdpId: idpID,
				},
			},
			wantErr: true,
		},
		{
			name: "next step oauth auth url",
			args: args{
				CTX,
				&user.StartIdentityProviderIntentRequest{
					IdpId: idpID,
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
					ResourceOwner: Tester.Organisation.ID,
				},
				url: "https://example.com/oauth/v2/authorize",
				parametersEqual: map[string]string{
					"client_id":     "clientID",
					"prompt":        "select_account",
					"redirect_uri":  "http://" + Tester.Config.ExternalDomain + ":8080/idps/callback",
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
					ResourceOwner: Tester.Organisation.ID,
				},
				url:                "http://" + Tester.Config.ExternalDomain + ":8000/sso",
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
					ResourceOwner: Tester.Organisation.ID,
				},
				url:                "http://" + Tester.Config.ExternalDomain + ":8000/sso",
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
					ResourceOwner: Tester.Organisation.ID,
				},
				postForm: true,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Client.StartIdentityProviderIntent(tt.args.ctx, tt.args.req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			if tt.want.url != "" {
				authUrl, err := url.Parse(got.GetAuthUrl())
				assert.NoError(t, err)

				assert.Len(t, authUrl.Query(), len(tt.want.parametersEqual)+len(tt.want.parametersExisting))

				for _, existing := range tt.want.parametersExisting {
					assert.True(t, authUrl.Query().Has(existing))
				}
				for key, equal := range tt.want.parametersEqual {
					assert.Equal(t, equal, authUrl.Query().Get(key))
				}
			}
			if tt.want.postForm {
				assert.NotEmpty(t, got.GetPostForm())
			}
			integration.AssertDetails(t, &user.StartIdentityProviderIntentResponse{
				Details: tt.want.details,
			}, got)
		})
	}
}

func TestServer_RetrieveIdentityProviderIntent(t *testing.T) {
	idpID := Tester.AddGenericOAuthProvider(t)
	intentID := Tester.CreateIntent(t, idpID)
	successfulID, token, changeDate, sequence := Tester.CreateSuccessfulOAuthIntent(t, idpID, "", "id")
	successfulWithUserID, WithUsertoken, WithUserchangeDate, WithUsersequence := Tester.CreateSuccessfulOAuthIntent(t, idpID, "user", "id")
	ldapSuccessfulID, ldapToken, ldapChangeDate, ldapSequence := Tester.CreateSuccessfulLDAPIntent(t, idpID, "", "id")
	ldapSuccessfulWithUserID, ldapWithUserToken, ldapWithUserChangeDate, ldapWithUserSequence := Tester.CreateSuccessfulLDAPIntent(t, idpID, "user", "id")
	samlSuccessfulID, samlToken, samlChangeDate, samlSequence := Tester.CreateSuccessfulSAMLIntent(t, idpID, "", "id")
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
			name: "retrieve successful intent",
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
					ResourceOwner: Tester.Organisation.ID,
					Sequence:      sequence,
				},
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Oauth{
						Oauth: &user.IDPOAuthAccessInformation{
							AccessToken: "accessToken",
							IdToken:     gu.Ptr("idToken"),
						},
					},
					IdpId:    idpID,
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
			name: "retrieve successful intent with linked user",
			args: args{
				CTX,
				&user.RetrieveIdentityProviderIntentRequest{
					IdpIntentId:    successfulWithUserID,
					IdpIntentToken: WithUsertoken,
				},
			},
			want: &user.RetrieveIdentityProviderIntentResponse{
				Details: &object.Details{
					ChangeDate:    timestamppb.New(WithUserchangeDate),
					ResourceOwner: Tester.Organisation.ID,
					Sequence:      WithUsersequence,
				},
				UserId: "user",
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Oauth{
						Oauth: &user.IDPOAuthAccessInformation{
							AccessToken: "accessToken",
							IdToken:     gu.Ptr("idToken"),
						},
					},
					IdpId:    idpID,
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
					ResourceOwner: Tester.Organisation.ID,
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
					IdpId:    idpID,
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
					ResourceOwner: Tester.Organisation.ID,
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
					IdpId:    idpID,
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
					ResourceOwner: Tester.Organisation.ID,
					Sequence:      samlSequence,
				},
				IdpInformation: &user.IDPInformation{
					Access: &user.IDPInformation_Saml{
						Saml: &user.IDPSAMLAccessInformation{
							Assertion: []byte("<Assertion xmlns=\"urn:oasis:names:tc:SAML:2.0:assertion\" ID=\"id\" IssueInstant=\"0001-01-01T00:00:00Z\" Version=\"\"><Issuer xmlns=\"urn:oasis:names:tc:SAML:2.0:assertion\" NameQualifier=\"\" SPNameQualifier=\"\" Format=\"\" SPProvidedID=\"\"></Issuer></Assertion>"),
						},
					},
					IdpId:    idpID,
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
	userIDWithoutAuth := Tester.CreateHumanUser(CTX).GetUserId()

	userIDWithPasskey := Tester.CreateHumanUser(CTX).GetUserId()
	Tester.RegisterUserPasskey(CTX, userIDWithPasskey)

	userMultipleAuth := Tester.CreateHumanUser(CTX).GetUserId()
	Tester.RegisterUserPasskey(CTX, userMultipleAuth)
	provider, err := Tester.Client.Mgmt.AddGenericOIDCProvider(CTX, &mgmt.AddGenericOIDCProviderRequest{
		Name:         "ListAuthenticationMethodTypes",
		Issuer:       "https://example.com",
		ClientId:     "client_id",
		ClientSecret: "client_secret",
	})
	require.NoError(t, err)
	idpLink, err := Tester.Client.UserV2.AddIDPLink(CTX, &user.AddIDPLinkRequest{UserId: userMultipleAuth, IdpLink: &user.IDPLink{
		IdpId:    provider.GetId(),
		UserId:   "external-id",
		UserName: "displayName",
	}})
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
			var got *user.ListAuthenticationMethodTypesResponse
			var err error

			for {
				got, err = Client.ListAuthenticationMethodTypes(tt.args.ctx, tt.args.req)
				if err == nil && !got.GetDetails().GetTimestamp().AsTime().Before(idpLink.GetDetails().GetChangeDate().AsTime()) {
					break
				}
				select {
				case <-CTX.Done():
					t.Fatal(CTX.Err(), err)
				case <-time.After(time.Second):
					t.Log("retrying ListAuthenticationMethodTypes")
					continue
				}
			}
			require.NoError(t, err)
			assert.Equal(t, tt.want.GetDetails().GetTotalResult(), got.GetDetails().GetTotalResult())
			require.Equal(t, tt.want.GetAuthMethodTypes(), got.GetAuthMethodTypes())
		})
	}
}
