package command

import (
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/internal/v2/repository/user"
)

func writeModelToHuman(wm *HumanWriteModel) *domain.Human {
	return &domain.Human{
		ObjectRoot: writeModelToObjectRoot(wm.WriteModel),
		Username:   wm.UserName,
		State:      wm.UserState,
		Profile: &domain.Profile{
			FirstName:         wm.FirstName,
			LastName:          wm.LastName,
			NickName:          wm.NickName,
			DisplayName:       wm.DisplayName,
			PreferredLanguage: wm.PreferredLanguage,
			Gender:            wm.Gender,
		},
		Email: &domain.Email{
			EmailAddress:    wm.Email,
			IsEmailVerified: wm.IsEmailVerified,
		},
		Address: &domain.Address{
			Country:       wm.Country,
			Locality:      wm.Locality,
			PostalCode:    wm.PostalCode,
			Region:        wm.Region,
			StreetAddress: wm.StreetAddress,
		},
	}
}

func writeModelToProfile(wm *HumanProfileWriteModel) *domain.Profile {
	return &domain.Profile{
		ObjectRoot:        writeModelToObjectRoot(wm.WriteModel),
		FirstName:         wm.FirstName,
		LastName:          wm.LastName,
		NickName:          wm.NickName,
		DisplayName:       wm.DisplayName,
		PreferredLanguage: wm.PreferredLanguage,
		Gender:            wm.Gender,
	}
}

func writeModelToEmail(wm *HumanEmailWriteModel) *domain.Email {
	return &domain.Email{
		ObjectRoot:      writeModelToObjectRoot(wm.WriteModel),
		EmailAddress:    wm.Email,
		IsEmailVerified: wm.IsEmailVerified,
	}
}

func writeModelToPhone(wm *HumanPhoneWriteModel) *domain.Phone {
	return &domain.Phone{
		ObjectRoot:      writeModelToObjectRoot(wm.WriteModel),
		PhoneNumber:     wm.Phone,
		IsPhoneVerified: wm.IsPhoneVerified,
	}
}
func writeModelToAddress(wm *HumanAddressWriteModel) *domain.Address {
	return &domain.Address{
		ObjectRoot:    writeModelToObjectRoot(wm.WriteModel),
		Country:       wm.Country,
		Locality:      wm.Locality,
		PostalCode:    wm.PostalCode,
		Region:        wm.Region,
		StreetAddress: wm.StreetAddress,
	}
}

func writeModelToMachine(wm *MachineWriteModel) *domain.Machine {
	return &domain.Machine{
		ObjectRoot:  writeModelToObjectRoot(wm.WriteModel),
		Name:        wm.Name,
		Description: wm.Description,
	}
}

func keyWriteModelToMachineKey(wm *MachineKeyWriteModel) *domain.MachineKey {
	return &domain.MachineKey{
		ObjectRoot:     writeModelToObjectRoot(wm.WriteModel),
		KeyID:          wm.KeyID,
		Type:           wm.KeyType,
		ExpirationDate: wm.ExpirationDate,
	}
}

func readModelToU2FTokens(wm *HumanU2FTokensReadModel) []*domain.WebAuthNToken {
	tokens := make([]*domain.WebAuthNToken, len(wm.WebAuthNTokens))
	for i, token := range wm.WebAuthNTokens {
		tokens[i] = writeModelToWebAuthN(token)
	}
	return tokens
}

func readModelToPasswordlessTokens(wm *HumanPasswordlessTokensReadModel) []*domain.WebAuthNToken {
	tokens := make([]*domain.WebAuthNToken, len(wm.WebAuthNTokens))
	for i, token := range wm.WebAuthNTokens {
		tokens[i] = writeModelToWebAuthN(token)
	}
	return tokens
}

func writeModelToWebAuthN(wm *HumanWebAuthNWriteModel) *domain.WebAuthNToken {
	return &domain.WebAuthNToken{
		ObjectRoot:        writeModelToObjectRoot(wm.WriteModel),
		WebAuthNTokenID:   wm.WebauthNTokenID,
		Challenge:         wm.Challenge,
		KeyID:             wm.KeyID,
		PublicKey:         wm.PublicKey,
		AttestationType:   wm.AttestationType,
		AAGUID:            wm.AAGUID,
		SignCount:         wm.SignCount,
		WebAuthNTokenName: wm.WebAuthNTokenName,
		State:             wm.State,
	}
}

func authRequestDomainToAuthRequestInfo(authRequest *domain.AuthRequest) *user.AuthRequestInfo {
	info := &user.AuthRequestInfo{
		ID:                  authRequest.ID,
		UserAgentID:         authRequest.AgentID,
		SelectedIDPConfigID: authRequest.SelectedIDPConfigID,
	}
	if authRequest.BrowserInfo != nil {
		info.BrowserInfo = &user.BrowserInfo{
			UserAgent:      authRequest.BrowserInfo.UserAgent,
			AcceptLanguage: authRequest.BrowserInfo.AcceptLanguage,
			RemoteIP:       authRequest.BrowserInfo.RemoteIP,
		}
	}
	return info
}
