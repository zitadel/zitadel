package user

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/google/go-cmp/cmp/cmpopts"
	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/assert"
	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func Test_patchHumanUserToCommand(t *testing.T) {
	type args struct {
		userId   string
		userName *string
		human    *user.UpdateUserRequest_Human
	}
	tests := []struct {
		name    string
		args    args
		want    *command.ChangeHuman
		wantErr assert.ErrorAssertionFunc
	}{{
		name: "single property",
		args: args{
			userId: "userId",
			human: &user.UpdateUserRequest_Human{
				Profile: &user.UpdateUserRequest_Human_Profile{
					GivenName: gu.Ptr("givenName"),
				},
			},
		},
		want: &command.ChangeHuman{
			ID: "userId",
			Profile: &command.Profile{
				FirstName: gu.Ptr("givenName"),
			},
		},
		wantErr: assert.NoError,
	}, {
		name: "all properties",
		args: args{
			userId:   "userId",
			userName: gu.Ptr("userName"),
			human: &user.UpdateUserRequest_Human{
				Profile: &user.UpdateUserRequest_Human_Profile{
					GivenName:         gu.Ptr("givenName"),
					FamilyName:        gu.Ptr("familyName"),
					NickName:          gu.Ptr("nickName"),
					DisplayName:       gu.Ptr("displayName"),
					PreferredLanguage: gu.Ptr("en-US"),
					Gender:            gu.Ptr(user.Gender_GENDER_FEMALE),
				},
				Email: &user.SetHumanEmail{
					Email: "email@example.com",
					Verification: &user.SetHumanEmail_IsVerified{
						IsVerified: true,
					},
				},
				Phone: &user.SetHumanPhone{
					Phone: "+123456789",
					Verification: &user.SetHumanPhone_IsVerified{
						IsVerified: true,
					},
				},
				Password: &user.SetPassword{
					Verification: &user.SetPassword_CurrentPassword{
						CurrentPassword: "currentPassword",
					},
					PasswordType: &user.SetPassword_Password{
						Password: &user.Password{
							Password:       "newPassword",
							ChangeRequired: true,
						},
					},
				},
			},
		},
		want: &command.ChangeHuman{
			ID:       "userId",
			Username: gu.Ptr("userName"),
			Profile: &command.Profile{
				FirstName:         gu.Ptr("givenName"),
				LastName:          gu.Ptr("familyName"),
				NickName:          gu.Ptr("nickName"),
				DisplayName:       gu.Ptr("displayName"),
				PreferredLanguage: &language.AmericanEnglish,
				Gender:            gu.Ptr(domain.GenderFemale),
			},
			Email: &command.Email{
				Address:  "email@example.com",
				Verified: true,
			},
			Phone: &command.Phone{
				Number:   "+123456789",
				Verified: true,
			},
			Password: &command.Password{
				OldPassword:    "currentPassword",
				Password:       "newPassword",
				ChangeRequired: true,
			},
		},
		wantErr: assert.NoError,
	}, {
		name: "set email and request code",
		args: args{
			userId: "userId",
			human: &user.UpdateUserRequest_Human{
				Email: &user.SetHumanEmail{
					Email: "email@example.com",
					Verification: &user.SetHumanEmail_ReturnCode{
						ReturnCode: &user.ReturnEmailVerificationCode{},
					},
				},
			},
		},
		want: &command.ChangeHuman{
			ID: "userId",
			Email: &command.Email{
				Address:    "email@example.com",
				ReturnCode: true,
			},
		},
		wantErr: assert.NoError,
	}, {
		name: "set email and send code",
		args: args{
			userId: "userId",
			human: &user.UpdateUserRequest_Human{
				Email: &user.SetHumanEmail{
					Email: "email@example.com",
					Verification: &user.SetHumanEmail_SendCode{
						SendCode: &user.SendEmailVerificationCode{},
					},
				},
			},
		},
		want: &command.ChangeHuman{
			ID: "userId",
			Email: &command.Email{
				Address: "email@example.com",
			},
		},
		wantErr: assert.NoError,
	}, {
		name: "set email and send code with template",
		args: args{
			userId: "userId",
			human: &user.UpdateUserRequest_Human{
				Email: &user.SetHumanEmail{
					Email: "email@example.com",
					Verification: &user.SetHumanEmail_SendCode{
						SendCode: &user.SendEmailVerificationCode{
							UrlTemplate: gu.Ptr("Code: {{.Code}}"),
						},
					},
				},
			},
		},
		want: &command.ChangeHuman{
			ID: "userId",
			Email: &command.Email{
				Address:     "email@example.com",
				URLTemplate: "Code: {{.Code}}",
			},
		},
		wantErr: assert.NoError,
	}, {
		name: "set phone and request code",
		args: args{
			userId: "userId",
			human: &user.UpdateUserRequest_Human{
				Phone: &user.SetHumanPhone{
					Phone: "+123456789",
					Verification: &user.SetHumanPhone_ReturnCode{
						ReturnCode: &user.ReturnPhoneVerificationCode{},
					},
				},
			},
		},
		want: &command.ChangeHuman{
			ID: "userId",
			Phone: &command.Phone{
				Number:     "+123456789",
				ReturnCode: true,
			},
		},
		wantErr: assert.NoError,
	}, {
		name: "set phone and send code",
		args: args{
			userId: "userId",
			human: &user.UpdateUserRequest_Human{
				Phone: &user.SetHumanPhone{
					Phone: "+123456789",
					Verification: &user.SetHumanPhone_SendCode{
						SendCode: &user.SendPhoneVerificationCode{},
					},
				},
			},
		},
		want: &command.ChangeHuman{
			ID: "userId",
			Phone: &command.Phone{
				Number: "+123456789",
			},
		},
		wantErr: assert.NoError,
	}, {
		name: "remove phone, ok",
		args: args{
			userId: "userId",
			human: &user.UpdateUserRequest_Human{
				Phone: &user.SetHumanPhone{},
			},
		},
		want: &command.ChangeHuman{
			ID: "userId",
			Phone: &command.Phone{
				Remove: true,
			},
		},
		wantErr: assert.NoError,
	}, {
		name: "remove phone with verification, error",
		args: args{
			userId: "userId",
			human: &user.UpdateUserRequest_Human{
				Phone: &user.SetHumanPhone{
					Verification: &user.SetHumanPhone_ReturnCode{},
				},
			},
		},
		wantErr: assert.Error,
	}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := updateHumanUserToCommand(tt.args.userId, tt.args.userName, tt.args.human)
			if !tt.wantErr(t, err, fmt.Sprintf("patchHumanUserToCommand(%v, %v, %v)", tt.args.userId, tt.args.userName, tt.args.human)) {
				return
			}
			if diff := cmp.Diff(tt.want, got, cmpopts.EquateComparable(language.Tag{})); diff != "" {
				t.Errorf("patchHumanUserToCommand() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
