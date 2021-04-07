package user

import (
	"github.com/caos/zitadel/internal/api/grpc/object"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/user/model"
	usr_grant_model "github.com/caos/zitadel/internal/usergrant/model"
	user_pb "github.com/caos/zitadel/pkg/grpc/user"
)

func UsersToPb(users []*model.UserView) []*user_pb.User {
	u := make([]*user_pb.User, len(users))
	for i, user := range users {
		u[i] = UserToPb(user)
	}
	return u
}
func UserToPb(user *model.UserView) *user_pb.User {
	return &user_pb.User{
		Id:                 user.ID,
		State:              ModelUserStateToPb(user.State),
		UserName:           user.UserName,
		LoginNames:         user.LoginNames,
		PreferredLoginName: user.PreferredLoginName,
		Type:               UserTypeToPb(user),
		Details: object.ToViewDetailsPb(
			user.Sequence,
			user.CreationDate,
			user.ChangeDate,
			user.ResourceOwner,
		),
	}
}

func UserTypeToPb(user *model.UserView) user_pb.UserType {
	if user.HumanView != nil {
		return &user_pb.User_Human{
			Human: HumanToPb(user.HumanView),
		}
	}
	if user.MachineView != nil {
		return &user_pb.User_Machine{
			Machine: MachineToPb(user.MachineView),
		}
	}
	return nil
}

func HumanToPb(view *model.HumanView) *user_pb.Human {
	return &user_pb.Human{
		Profile: &user_pb.Profile{
			FirstName:         view.FirstName,
			LastName:          view.LastName,
			NickName:          view.NickName,
			DisplayName:       view.DisplayName,
			PreferredLanguage: view.PreferredLanguage,
			Gender:            GenderToPb(view.Gender),
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

func MachineToPb(view *model.MachineView) *user_pb.Machine {
	return &user_pb.Machine{
		Name:        view.Name,
		Description: view.Description,
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

func EmailToPb(email *model.Email) *user_pb.Email {
	return &user_pb.Email{
		Email:           email.EmailAddress,
		IsEmailVerified: email.IsEmailVerified,
	}
}

func PhoneToPb(phone *model.Phone) *user_pb.Phone {
	return &user_pb.Phone{
		Phone:           phone.PhoneNumber,
		IsPhoneVerified: phone.IsPhoneVerified,
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

func ModelUserGrantStateToPb(state usr_grant_model.UserGrantState) user_pb.UserGrantState {
	switch state {
	case usr_grant_model.UserGrantStateActive:
		return user_pb.UserGrantState_USER_GRANT_STATE_ACTIVE
	case usr_grant_model.UserGrantStateInactive:
		return user_pb.UserGrantState_USER_GRANT_STATE_INACTIVE
	default:
		return user_pb.UserGrantState_USER_GRANT_STATE_UNSPECIFIED
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

func AuthFactorsToPb(mfas []*model.MultiFactor) []*user_pb.AuthFactor {
	factors := make([]*user_pb.AuthFactor, len(mfas))
	for i, mfa := range mfas {
		factors[i] = AuthFactorToPb(mfa)
	}
	return factors
}

func AuthFactorToPb(mfa *model.MultiFactor) *user_pb.AuthFactor {
	factor := &user_pb.AuthFactor{
		State: MFAStateToPb(mfa.State),
	}
	switch mfa.Type {
	case model.MFATypeOTP:
		factor.Type = &user_pb.AuthFactor_Otp{
			Otp: &user_pb.AuthFactorOTP{},
		}
	case model.MFATypeU2F:
		factor.Type = &user_pb.AuthFactor_U2F{
			U2F: &user_pb.AuthFactorU2F{
				Id:   mfa.ID,
				Name: mfa.Attribute,
			},
		}
	}
	return factor
}

func MFAStateToPb(state model.MFAState) user_pb.AuthFactorState {
	switch state {
	case model.MFAStateNotReady:
		return user_pb.AuthFactorState_AUTH_FACTOR_STATE_NOT_READY
	case model.MFAStateReady:
		return user_pb.AuthFactorState_AUTH_FACTOR_STATE_READY
	default:
		return user_pb.AuthFactorState_AUTH_FACTOR_STATE_UNSPECIFIED
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
		PublicKey: token.PublicKey,
	}
}

func ExternalIDPViewsToExternalIDPs(externalIDPs []*model.ExternalIDPView) []*domain.ExternalIDP {
	idps := make([]*domain.ExternalIDP, len(externalIDPs))
	for i, idp := range externalIDPs {
		idps[i] = &domain.ExternalIDP{
			ObjectRoot: models.ObjectRoot{
				AggregateID:   idp.UserID,
				ResourceOwner: idp.ResourceOwner,
			},
			IDPConfigID:    idp.IDPConfigID,
			ExternalUserID: idp.ExternalUserID,
			DisplayName:    idp.UserDisplayName,
		}
	}
	return idps
}
