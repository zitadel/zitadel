package resources

import (
	"context"

	"golang.org/x/text/language"

	"github.com/zitadel/zitadel/internal/command"
	"github.com/zitadel/zitadel/internal/domain"
)

func (h *UsersHandler) mapToAddHuman(ctx context.Context, scimUser *ScimUser) (*command.AddHuman, error) {
	// zitadel has its own state mechanism
	// ignore scimUser.Active
	human := &command.AddHuman{
		Username:    scimUser.UserName,
		NickName:    scimUser.NickName,
		DisplayName: scimUser.DisplayName,
		Email:       h.mapPrimaryEmail(scimUser),
		Phone:       h.mapPrimaryPhone(scimUser),
	}

	md, err := h.mapMetadataToCommands(ctx, scimUser)
	if err != nil {
		return nil, err
	}
	human.Metadata = md

	if scimUser.Password != nil {
		human.Password = scimUser.Password.String()
		scimUser.Password = nil
	}

	if scimUser.Name != nil {
		human.FirstName = scimUser.Name.GivenName
		human.LastName = scimUser.Name.FamilyName

		// the direct mapping displayName => displayName has priority
		// over the formatted name assignment
		if human.DisplayName == "" {
			human.DisplayName = scimUser.Name.Formatted
		}
	}

	if err := domain.LanguageIsDefined(scimUser.PreferredLanguage); err != nil {
		human.PreferredLanguage = language.English
		scimUser.PreferredLanguage = language.English
	}

	return human, nil
}

func (h *UsersHandler) mapPrimaryEmail(scimUser *ScimUser) command.Email {
	for _, email := range scimUser.Emails {
		if !email.Primary {
			continue
		}

		return command.Email{
			Address:  domain.EmailAddress(email.Value),
			Verified: h.config.EmailVerified,
		}
	}

	return command.Email{}
}

func (h *UsersHandler) mapPrimaryPhone(scimUser *ScimUser) command.Phone {
	for _, phone := range scimUser.PhoneNumbers {
		if !phone.Primary {
			continue
		}

		return command.Phone{
			Number:   domain.PhoneNumber(phone.Value),
			Verified: h.config.PhoneVerified,
		}
	}

	return command.Phone{}
}
