package user

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/user/model"
	user_pb "github.com/caos/zitadel/pkg/grpc/user"
)

func UserToPb(user *model.UserView) *user_pb.User {
	return &user_pb.User{
		Id:                 user.ID,
		State:              ModelUserStateToPb(user.State),
		UserName:           user.UserName,
		LoginNames:         user.LoginNames,
		PreferredLoginName: user.PreferredLoginName,
		Details: object.ToDetailsPb(
			user.Sequence,
			user.CreationDate,
			user.ChangeDate,
			user.ResourceOwner,
		),
	}
}

func ProfileToPb(profile *model.Profile) *user_pb.Profile {
	return &user_pb.Profile{
		FirstName:         profile.FirstName,
		LastName:          profile.LastName,
		NickName:          profile.NickName,
		DisplayName:       profile.DisplayName,
		PreferredLanguage: profile.PreferredLanguage.String(),
		Gender:            GenderToPb(profile.Gender),
	}
}

func ModelEmailToPb(email *model.Email) *user_pb.Email {
	return &user_pb.Email{
		Email:           email.EmailAddress,
		IsEmailVerified: email.IsEmailVerified,
	}
}

func ModelPhoneToPb(phone *model.Phone) *user_pb.Phone {
	return &user_pb.Phone{
		Phone:           phone.PhoneNumber,
		IsPhoneVerified: phone.IsPhoneVerified,
	}
}

func ModelAddressToPb(address *model.Address) *user_pb.Address {
	return &user_pb.Address{
		Country:       address.Country,
		Locality:      address.Locality,
		PostalCode:    address.PostalCode,
		Region:        address.Region,
		StreetAddress: address.StreetAddress,
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

func ModelUserStateToPb(state model.UserState) user_pb.UserState {
	switch state {
	case model.UserStateActive:
		return user_pb.UserState_USER_STATE_ACTIVE
	case model.UserStateInactive:
		return user_pb.UserState_USER_STATE_INACTIVE
	case model.UserStateDeleted:
		return user_pb.UserState_USER_STATE_DELETED
	case model.UserStateInitial:
		return user_pb.UserState_USER_STATE_INITIAL
	case model.UserStateLocked:
		return user_pb.UserState_USER_STATE_LOCKED
	case model.UserStateSuspend:
		return user_pb.UserState_USER_STATE_SUSPEND
	default:
		return user_pb.UserState_USER_STATE_UNSPECIFIED
	}
}

func GenderToPb(gender model.Gender) user_pb.Gender {
	switch gender {
	case model.GenderDiverse:
		return user_pb.Gender_GENDER_DIVERSE
	case model.GenderFemale:
		return user_pb.Gender_GENDER_FEMALE
	case model.GenderMale:
		return user_pb.Gender_GENDER_MALE
	default:
		return user_pb.Gender_GENDER_UNSPECIFIED
	}
}

func MultiFactorsToPb(mfas []*model.MultiFactor) []*user_pb.MultiFactor {
	factors := make([]*user_pb.MultiFactor, len(mfas))
	for i, mfa := range mfas {
		factors[i] = MultiFactorToPb(mfa)
	}
	return factors
}

func MultiFactorToPb(mfa *model.MultiFactor) *user_pb.MultiFactor {
	factor := &user_pb.MultiFactor{
		State: MFAStateToPb(mfa.State),
	}
	switch mfa.Type {
	case model.MFATypeOTP:
		factor.Type = &user_pb.MultiFactor_Otp{
			Otp: &user_pb.MultiFactorOTP{},
		}
	case model.MFATypeU2F:
		factor.Type = &user_pb.MultiFactor_U2F{
			U2F: &user_pb.MultiFactorU2F{
				Id:   mfa.ID,
				Name: mfa.Attribute,
			},
		}
	}
	return factor
}

func MFAStateToPb(state model.MFAState) user_pb.MultiFactorState {
	switch state {
	case model.MFAStateNotReady:
		return user_pb.MultiFactorState_MULTI_FACTOR_STATE_NOT_READY
	case model.MFAStateReady:
		return user_pb.MultiFactorState_MULTI_FACTOR_STATE_READY
	default:
		return user_pb.MultiFactorState_MULTI_FACTOR_STATE_UNSPECIFIED
	}
}

func WebAuthNTokensViewToPb(tokens []*model.WebAuthNView) []*user_pb.WebAuthNToken {
	t := make([]*user_pb.WebAuthNToken, len(tokens))
	for i, token := range tokens {
		t[i] = WebAuthNTokenViewToPb(token)
	}
	return t
}

func WebAuthNTokenViewToPb(token *model.WebAuthNView) *user_pb.WebAuthNToken {
	return &user_pb.WebAuthNToken{
		Id:    token.TokenID,
		State: MFAStateToPb(token.State),
		Name:  token.Name,
	}
}

func WebAuthNTokenToWebAuthNKeyPb(token *domain.WebAuthNToken) *user_pb.WebAuthNKey {
	return &user_pb.WebAuthNKey{
		Id:        string(token.KeyID), //TODO: ask if it's the correct id?
		PublicKey: token.PublicKey,
	}
}
