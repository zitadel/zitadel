package convert

import (
	"github.com/muhlemmer/gu"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func humanToPb(userQ *query.Human, assetPrefix, owner string) *user.HumanUser {
	if userQ == nil {
		return nil
	}

	var passwordChanged, mfaInitSkipped *timestamppb.Timestamp
	if !userQ.PasswordChanged.IsZero() {
		passwordChanged = timestamppb.New(userQ.PasswordChanged)
	}
	if !userQ.MFAInitSkipped.IsZero() {
		mfaInitSkipped = timestamppb.New(userQ.MFAInitSkipped)
	}
	return &user.HumanUser{
		Profile: &user.HumanProfile{
			GivenName:         userQ.FirstName,
			FamilyName:        userQ.LastName,
			NickName:          gu.Ptr(userQ.NickName),
			DisplayName:       gu.Ptr(userQ.DisplayName),
			PreferredLanguage: gu.Ptr(userQ.PreferredLanguage.String()),
			Gender:            gu.Ptr(genderToPb(userQ.Gender)),
			AvatarUrl:         domain.AvatarURL(assetPrefix, owner, userQ.AvatarKey),
		},
		Email: &user.HumanEmail{
			Email:      string(userQ.Email),
			IsVerified: userQ.IsEmailVerified,
		},
		Phone: &user.HumanPhone{
			Phone:      string(userQ.Phone),
			IsVerified: userQ.IsPhoneVerified,
		},
		PasswordChangeRequired: userQ.PasswordChangeRequired,
		PasswordChanged:        passwordChanged,
		MfaInitSkipped:         mfaInitSkipped,
	}
}

func genderToPb(gender domain.Gender) user.Gender {
	switch gender {
	case domain.GenderDiverse:
		return user.Gender_GENDER_DIVERSE
	case domain.GenderFemale:
		return user.Gender_GENDER_FEMALE
	case domain.GenderMale:
		return user.Gender_GENDER_MALE
	case domain.GenderUnspecified:
		return user.Gender_GENDER_UNSPECIFIED
	default:
		return user.Gender_GENDER_UNSPECIFIED
	}
}
