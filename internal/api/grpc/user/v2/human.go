package user

import (
	"io"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
)

func patchHumanUserToCommand(userId string, userName *string, human *user.UpdateUserRequest_Human) (*command.ChangeHuman, error) {
	phone := human.GetPhone()
	if phone != nil && phone.Phone == "" && phone.GetVerification() != nil {
		// TODO: Translate
		return nil, zerrors.ThrowInvalidArgument(nil, "USERv2-4f3d6", "Errors.User.Phone.RemoveWithVerification")
	}
	email, err := setHumanEmailToEmail(human.Email, userId)
	if err != nil {
		return nil, err
	}
	return &command.ChangeHuman{
		ID:       userId,
		Username: userName,
		Profile:  SetHumanProfileToProfile(human.Profile),
		Email:    email,
		Phone:    setHumanPhoneToPhone(human.Phone, true),
		Password: setHumanPasswordToPassword(human.Password),
	}, nil
}

func updateHumanUserRequestToChangeHuman(req *user.UpdateHumanUserRequest) (*command.ChangeHuman, error) {
	email, err := setHumanEmailToEmail(req.Email, req.GetUserId())
	if err != nil {
		return nil, err
	}
	changeHuman := &command.ChangeHuman{
		ID:       req.GetUserId(),
		Username: req.Username,
		Email:    email,
		Phone:    setHumanPhoneToPhone(req.Phone, false),
		Password: setHumanPasswordToPassword(req.Password),
	}
	if profile := req.GetProfile(); profile != nil {
		var firstName *string
		if profile.GivenName != "" {
			firstName = &profile.GivenName
		}
		var lastName *string
		if profile.FamilyName != "" {
			lastName = &profile.FamilyName
		}
		changeHuman.Profile = SetHumanProfileToProfile(&user.UpdateUserRequest_Human_Profile{
			GivenName:         firstName,
			FamilyName:        lastName,
			NickName:          profile.NickName,
			DisplayName:       profile.DisplayName,
			PreferredLanguage: profile.PreferredLanguage,
			Gender:            profile.Gender,
		})
	}
	return changeHuman, nil
}

func SetHumanProfileToProfile(profile *user.UpdateUserRequest_Human_Profile) *command.Profile {
	if profile == nil {
		return nil
	}
	return &command.Profile{
		FirstName:         profile.GivenName,
		LastName:          profile.FamilyName,
		NickName:          profile.NickName,
		DisplayName:       profile.DisplayName,
		PreferredLanguage: ifNotNilPtr(profile.PreferredLanguage, language.Make),
		Gender:            ifNotNilPtr(profile.Gender, genderToDomain),
	}
}

func setHumanEmailToEmail(email *user.SetHumanEmail, userID string) (*command.Email, error) {
	if email == nil {
		return nil, nil
	}
	var urlTemplate string
	if email.GetSendCode() != nil && email.GetSendCode().UrlTemplate != nil {
		urlTemplate = *email.GetSendCode().UrlTemplate
		if err := domain.RenderConfirmURLTemplate(io.Discard, urlTemplate, userID, "code", "orgID"); err != nil {
			return nil, err
		}
	}
	return &command.Email{
		Address:     domain.EmailAddress(email.Email),
		Verified:    email.GetIsVerified(),
		ReturnCode:  email.GetReturnCode() != nil,
		URLTemplate: urlTemplate,
	}, nil
}

func setHumanPhoneToPhone(phone *user.SetHumanPhone, withRemove bool) *command.Phone {
	if phone == nil {
		return nil
	}
	number := phone.GetPhone()
	return &command.Phone{
		Number:     domain.PhoneNumber(number),
		Verified:   phone.GetIsVerified(),
		ReturnCode: phone.GetReturnCode() != nil,
		Remove:     withRemove && number == "",
	}
}

func setHumanPasswordToPassword(password *user.SetPassword) *command.Password {
	if password == nil {
		return nil
	}
	return &command.Password{
		PasswordCode:        password.GetVerificationCode(),
		OldPassword:         password.GetCurrentPassword(),
		Password:            password.GetPassword().GetPassword(),
		EncodedPasswordHash: password.GetHashedPassword().GetHash(),
		ChangeRequired:      password.GetPassword().GetChangeRequired() || password.GetHashedPassword().GetChangeRequired(),
	}
}
