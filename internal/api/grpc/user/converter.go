package user

import (
	"github.com/zitadel/zitadel/internal/api/grpc/object"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	user_pb "github.com/zitadel/zitadel/pkg/grpc/user"
)

func UsersToPb(users []*query.User, assetPrefix string) []*user_pb.User {
	u := make([]*user_pb.User, len(users))
	for i, user := range users {
		u[i] = UserToPb(user, assetPrefix)
	}
	return u
}
func UserToPb(user *query.User, assetPrefix string) *user_pb.User {
	return &user_pb.User{
		Id:                 user.ID,
		State:              UserStateToPb(user.State),
		UserName:           user.Username,
		LoginNames:         user.LoginNames,
		PreferredLoginName: user.PreferredLoginName,
		Type:               UserTypeToPb(user, assetPrefix),
		Details: object.ToViewDetailsPb(
			user.Sequence,
			user.CreationDate,
			user.ChangeDate,
			user.ResourceOwner,
		),
	}
}

func UserTypeToPb(user *query.User, assetPrefix string) user_pb.UserType {
	if user.Human != nil {
		return &user_pb.User_Human{
			Human: HumanToPb(user.Human, assetPrefix, user.ResourceOwner),
		}
	}
	if user.Machine != nil {
		return &user_pb.User_Machine{
			Machine: MachineToPb(user.Machine),
		}
	}
	return nil
}

func HumanToPb(view *query.Human, assetPrefix, owner string) *user_pb.Human {
	return &user_pb.Human{
		Profile: &user_pb.Profile{
			FirstName:         view.FirstName,
			LastName:          view.LastName,
			NickName:          view.NickName,
			DisplayName:       view.DisplayName,
			PreferredLanguage: view.PreferredLanguage.String(),
			Gender:            GenderToPb(view.Gender),
			AvatarUrl:         domain.AvatarURL(assetPrefix, owner, view.AvatarKey),
		},
		Email: &user_pb.Email{
			Email:           view.Email,
			IsEmailVerified: view.IsEmailVerified,
		},
		Phone: &user_pb.Phone{
			Phone:           view.Phone,
			IsPhoneVerified: view.IsPhoneVerified,
		},
	}
}

func MachineToPb(view *query.Machine) *user_pb.Machine {
	return &user_pb.Machine{
		Name:            view.Name,
		Description:     view.Description,
		HasSecret:       view.HasSecret,
		AccessTokenType: AccessTokenTypeToPb(view.AccessTokenType),
	}
}

func ProfileToPb(profile *query.Profile, assetPrefix string) *user_pb.Profile {
	return &user_pb.Profile{
		FirstName:         profile.FirstName,
		LastName:          profile.LastName,
		NickName:          profile.NickName,
		DisplayName:       profile.DisplayName,
		PreferredLanguage: profile.PreferredLanguage.String(),
		Gender:            GenderToPb(profile.Gender),
		AvatarUrl:         domain.AvatarURL(assetPrefix, profile.ResourceOwner, profile.AvatarKey),
	}
}

func EmailToPb(email *query.Email) *user_pb.Email {
	return &user_pb.Email{
		Email:           email.Email,
		IsEmailVerified: email.IsVerified,
	}
}

func PhoneToPb(phone *query.Phone) *user_pb.Phone {
	return &user_pb.Phone{
		Phone:           phone.Phone,
		IsPhoneVerified: phone.IsVerified,
	}
}

func ModelEmailToPb(email *query.Email) *user_pb.Email {
	return &user_pb.Email{
		Email:           email.Email,
		IsEmailVerified: email.IsVerified,
	}
}

func ModelPhoneToPb(phone *query.Phone) *user_pb.Phone {
	return &user_pb.Phone{
		Phone:           phone.Phone,
		IsPhoneVerified: phone.IsVerified,
	}
}

func GenderToDomain(gender user_pb.Gender) domain.Gender {
	switch gender {
	case user_pb.Gender_GENDER_DIVERSE:
		return domain.GenderDiverse
	case user_pb.Gender_GENDER_MALE:
		return domain.GenderMale
	case user_pb.Gender_GENDER_FEMALE:
		return domain.GenderFemale
	default:
		return -1
	}
}

func AccessTokenTypeToDomain(accessTokenType user_pb.AccessTokenType) domain.OIDCTokenType {
	switch accessTokenType {
	case user_pb.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER:
		return domain.OIDCTokenTypeBearer
	case user_pb.AccessTokenType_ACCESS_TOKEN_TYPE_JWT:
		return domain.OIDCTokenTypeJWT
	default:
		return -1
	}
}

func UserStateToPb(state domain.UserState) user_pb.UserState {
	switch state {
	case domain.UserStateActive:
		return user_pb.UserState_USER_STATE_ACTIVE
	case domain.UserStateInactive:
		return user_pb.UserState_USER_STATE_INACTIVE
	case domain.UserStateDeleted:
		return user_pb.UserState_USER_STATE_DELETED
	case domain.UserStateInitial:
		return user_pb.UserState_USER_STATE_INITIAL
	case domain.UserStateLocked:
		return user_pb.UserState_USER_STATE_LOCKED
	case domain.UserStateSuspend:
		return user_pb.UserState_USER_STATE_SUSPEND
	default:
		return user_pb.UserState_USER_STATE_UNSPECIFIED
	}
}

func GenderToPb(gender domain.Gender) user_pb.Gender {
	switch gender {
	case domain.GenderDiverse:
		return user_pb.Gender_GENDER_DIVERSE
	case domain.GenderFemale:
		return user_pb.Gender_GENDER_FEMALE
	case domain.GenderMale:
		return user_pb.Gender_GENDER_MALE
	default:
		return user_pb.Gender_GENDER_UNSPECIFIED
	}
}

func AccessTokenTypeToPb(accessTokenType domain.OIDCTokenType) user_pb.AccessTokenType {
	switch accessTokenType {
	case domain.OIDCTokenTypeBearer:
		return user_pb.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER
	case domain.OIDCTokenTypeJWT:
		return user_pb.AccessTokenType_ACCESS_TOKEN_TYPE_JWT
	default:
		return user_pb.AccessTokenType_ACCESS_TOKEN_TYPE_BEARER
	}
}

func AuthMethodsToPb(mfas *query.AuthMethods) []*user_pb.AuthFactor {
	factors := make([]*user_pb.AuthFactor, len(mfas.AuthMethods))
	for i, mfa := range mfas.AuthMethods {
		factors[i] = AuthMethodToPb(mfa)
	}
	return factors
}

func AuthMethodToPb(mfa *query.AuthMethod) *user_pb.AuthFactor {
	factor := &user_pb.AuthFactor{
		State: MFAStateToPb(mfa.State),
	}
	switch mfa.Type {
	case domain.UserAuthMethodTypeOTP:
		factor.Type = &user_pb.AuthFactor_Otp{
			Otp: &user_pb.AuthFactorOTP{},
		}
	case domain.UserAuthMethodTypeU2F:
		factor.Type = &user_pb.AuthFactor_U2F{
			U2F: &user_pb.AuthFactorU2F{
				Id:   mfa.TokenID,
				Name: mfa.Name,
			},
		}
	}
	return factor
}

func MFAStateToPb(state domain.MFAState) user_pb.AuthFactorState {
	switch state {
	case domain.MFAStateNotReady:
		return user_pb.AuthFactorState_AUTH_FACTOR_STATE_NOT_READY
	case domain.MFAStateReady:
		return user_pb.AuthFactorState_AUTH_FACTOR_STATE_READY
	default:
		return user_pb.AuthFactorState_AUTH_FACTOR_STATE_UNSPECIFIED
	}
}

func UserAuthMethodsToWebAuthNTokenPb(methods *query.AuthMethods) []*user_pb.WebAuthNToken {
	t := make([]*user_pb.WebAuthNToken, len(methods.AuthMethods))
	for i, token := range methods.AuthMethods {
		t[i] = UserAuthMethodToWebAuthNTokenPb(token)
	}
	return t
}

func UserAuthMethodToWebAuthNTokenPb(token *query.AuthMethod) *user_pb.WebAuthNToken {
	return &user_pb.WebAuthNToken{
		Id:    token.TokenID,
		State: MFAStateToPb(token.State),
		Name:  token.Name,
	}
}

func ExternalIDPViewsToExternalIDPs(externalIDPs []*query.IDPUserLink) []*domain.UserIDPLink {
	idps := make([]*domain.UserIDPLink, len(externalIDPs))
	for i, idp := range externalIDPs {
		idps[i] = &domain.UserIDPLink{
			ObjectRoot: models.ObjectRoot{
				AggregateID:   idp.UserID,
				ResourceOwner: idp.ResourceOwner,
			},
			IDPConfigID:    idp.IDPID,
			ExternalUserID: idp.ProvidedUserID,
			DisplayName:    idp.ProvidedUsername,
		}
	}
	return idps
}
