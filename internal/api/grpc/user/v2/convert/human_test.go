package convert

import (
	"testing"
	"time"

	"github.com/muhlemmer/gu"
	"github.com/stretchr/testify/require"
	"golang.org/x/text/language"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func Test_genderToPb(t *testing.T) {
	t.Parallel()

	tt := []struct {
		name   string
		input  domain.Gender
		expect user.Gender
	}{
		{
			name:   "Diverse",
			input:  domain.GenderDiverse,
			expect: user.Gender_GENDER_DIVERSE,
		},
		{
			name:   "Female",
			input:  domain.GenderFemale,
			expect: user.Gender_GENDER_FEMALE,
		},
		{
			name:   "Male",
			input:  domain.GenderMale,
			expect: user.Gender_GENDER_MALE,
		},
		{
			name:   "Unspecified",
			input:  domain.GenderUnspecified,
			expect: user.Gender_GENDER_UNSPECIFIED,
		},
		{
			name:   "Unknown value",
			input:  domain.Gender(999),
			expect: user.Gender_GENDER_UNSPECIFIED,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got := genderToPb(tc.input)
			require.Equal(t, tc.expect, got)
		})
	}
}

func Test_humanToPb(t *testing.T) {
	t.Parallel()

	now := time.Now().UTC()
	assetPrefix := "prefix"
	owner := "owner"
	avatarKey := "avatar-key"

	tt := []struct {
		name     string
		input    *query.Human
		expected *user.HumanUser
	}{
		{
			name: "all fields set",
			input: &query.Human{
				FirstName:              "John",
				LastName:               "Doe",
				NickName:               "JD",
				DisplayName:            "Johnny",
				PreferredLanguage:      language.English,
				Gender:                 domain.GenderMale,
				AvatarKey:              avatarKey,
				Email:                  "john.doe@example.com",
				IsEmailVerified:        true,
				Phone:                  "+123456789",
				IsPhoneVerified:        true,
				PasswordChangeRequired: true,
				PasswordChanged:        now,
				MFAInitSkipped:         now,
			},
			expected: &user.HumanUser{
				Profile: &user.HumanProfile{
					GivenName:         "John",
					FamilyName:        "Doe",
					NickName:          gu.Ptr("JD"),
					DisplayName:       gu.Ptr("Johnny"),
					PreferredLanguage: gu.Ptr("en"),
					Gender:            gu.Ptr(user.Gender_GENDER_MALE),
					AvatarUrl:         domain.AvatarURL(assetPrefix, owner, avatarKey),
				},
				Email: &user.HumanEmail{
					Email:      "john.doe@example.com",
					IsVerified: true,
				},
				Phone: &user.HumanPhone{
					Phone:      "+123456789",
					IsVerified: true,
				},
				PasswordChangeRequired: true,
				PasswordChanged:        timestamppb.New(now),
				MfaInitSkipped:         timestamppb.New(now),
			},
		},
		{
			name: "zero times, not verified, unspecified gender",
			input: &query.Human{
				FirstName:              "Jane",
				LastName:               "Smith",
				NickName:               "",
				DisplayName:            "",
				PreferredLanguage:      language.German,
				Gender:                 domain.GenderUnspecified,
				AvatarKey:              "",
				Email:                  "jane.smith@example.com",
				IsEmailVerified:        false,
				Phone:                  "",
				IsPhoneVerified:        false,
				PasswordChangeRequired: false,
				PasswordChanged:        time.Time{},
				MFAInitSkipped:         time.Time{},
			},
			expected: &user.HumanUser{
				Profile: &user.HumanProfile{
					GivenName:         "Jane",
					FamilyName:        "Smith",
					NickName:          gu.Ptr(""),
					DisplayName:       gu.Ptr(""),
					PreferredLanguage: gu.Ptr("de"),
					Gender:            gu.Ptr(user.Gender_GENDER_UNSPECIFIED),
					AvatarUrl:         domain.AvatarURL(assetPrefix, owner, ""),
				},
				Email: &user.HumanEmail{
					Email:      "jane.smith@example.com",
					IsVerified: false,
				},
				Phone: &user.HumanPhone{
					Phone:      "",
					IsVerified: false,
				},
				PasswordChangeRequired: false,
				PasswordChanged:        nil,
				MfaInitSkipped:         nil,
			},
		},
		{
			name:     "nil input",
			input:    nil,
			expected: nil,
		},
	}

	for _, tc := range tt {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			if tc.input == nil {
				require.Nil(t, humanToPb(nil, assetPrefix, owner))
				return
			}
			got := humanToPb(tc.input, assetPrefix, owner)
			require.Equal(t, tc.expected.Profile, got.Profile)
			require.Equal(t, tc.expected.Email, got.Email)
			require.Equal(t, tc.expected.Phone, got.Phone)
			require.Equal(t, tc.expected.PasswordChangeRequired, got.PasswordChangeRequired)
			require.Equal(t, tc.expected.PasswordChanged, got.PasswordChanged)
			require.Equal(t, tc.expected.MfaInitSkipped, got.MfaInitSkipped)
		})
	}
}
