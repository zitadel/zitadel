package admin

import (
	"github.com/caos/logging"
	"golang.org/x/text/language"

	user_grpc "github.com/zitadel/zitadel/internal/api/grpc/user"
	"github.com/zitadel/zitadel/internal/command"
	admin_grpc "github.com/zitadel/zitadel/pkg/grpc/admin"
)

func setUpOrgHumanToCommand(human *admin_grpc.SetUpOrgRequest_Human) command.AddHuman {
	var lang language.Tag
	lang, err := language.Parse(human.Profile.PreferredLanguage)
	logging.OnError(err).Debug("unable to parse language")
	return command.AddHuman{
		Username:      human.UserName,
		FirstName:     human.Profile.FirstName,
		LastName:      human.Profile.LastName,
		NickName:      human.Profile.NickName,
		DisplayName:   human.Profile.DisplayName,
		PreferredLang: lang,
		Gender:        user_grpc.GenderToDomain(human.Profile.Gender),
		Email:         setUpOrgHumanEmailToDomain(human.Email),
		Phone:         setUpOrgHumanPhoneToDomain(human.Phone),
		Password:      human.Password,
	}
}

func setUpOrgHumanEmailToDomain(email *admin_grpc.SetUpOrgRequest_Human_Email) command.Email {
	return command.Email{
		Address:  email.Email,
		Verified: email.IsEmailVerified,
	}
}

func setUpOrgHumanPhoneToDomain(phone *admin_grpc.SetUpOrgRequest_Human_Phone) command.Phone {
	if phone == nil {
		return command.Phone{}
	}
	return command.Phone{
		Number:   phone.Phone,
		Verified: phone.IsPhoneVerified,
	}
}
