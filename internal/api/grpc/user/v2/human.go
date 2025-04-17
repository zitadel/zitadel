package user

import (
	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/pkg/grpc/user/v2"
	"golang.org/x/text/language"
	"io"
)

func patchHumanUserToCommand(userId string, userName *string, human *user.UpdateUserRequest_Human) (*command.ChangeHuman, error) {
	email, err := SetHumanEmailToEmail(human.Email, userId)
	if err != nil {
		return nil, err
	}
	return &command.ChangeHuman{
		ID:       userId,
		Username: userName,
		Profile:  SetHumanProfileToProfile(human.Profile),
		Email:    email,
		Phone:    SetHumanPhoneToPhone(human.Phone),
		Password: SetHumanPasswordToPassword(human.Password),
	}, nil
}

func UpdateHumanUserRequestToChangeHuman(req *user.UpdateHumanUserRequest) (*command.ChangeHuman, error) {
	email, err := SetHumanEmailToEmail(req.Email, req.GetUserId())
	if err != nil {
		return nil, err
	}
	changeHuman := &command.ChangeHuman{
		ID:       req.GetUserId(),
		Username: req.Username,
		Email:    email,
		Phone:    SetHumanPhoneToPhone(req.Phone),
		Password: SetHumanPasswordToPassword(req.Password),
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

func SetHumanEmailToEmail(email *user.SetHumanEmail, userID string) (*command.Email, error) {
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

func SetHumanPhoneToPhone(phone *user.SetHumanPhone) *command.Phone {
	if phone == nil {
		return nil
	}
	return &command.Phone{
		Number:     domain.PhoneNumber(phone.GetPhone()),
		Verified:   phone.GetIsVerified(),
		ReturnCode: phone.GetReturnCode() != nil,
	}
}

func SetHumanPasswordToPassword(password *user.SetPassword) *command.Password {
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
