package command

import (
	"encoding/base64"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/repository/user"
)

func writeModelToHuman(wm *HumanWriteModel) *domain.Human {
	human := &domain.Human{
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
	}
	if wm.Phone != "" {
		human.Phone = &domain.Phone{
			PhoneNumber: wm.Phone,
		}
	}
	if wm.Country != "" || wm.Locality != "" || wm.PostalCode != "" || wm.Region != "" || wm.StreetAddress != "" {
		human.Address = &domain.Address{
			Country:       wm.Country,
			Locality:      wm.Locality,
			PostalCode:    wm.PostalCode,
			Region:        wm.Region,
			StreetAddress: wm.StreetAddress,
		}
	}
	return human
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
		Username:    wm.UserName,
		Name:        wm.Name,
		Description: wm.Description,
		State:       wm.UserState,
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

func personalTokenWriteModelToToken(wm *PersonalAccessTokenWriteModel, algorithm crypto.EncryptionAlgorithm) (*domain.Token, string, error) {
	encrypted, err := algorithm.Encrypt([]byte(wm.TokenID + ":" + wm.AggregateID))
	if err != nil {
		return nil, "", err
	}
	return &domain.Token{
		ObjectRoot: writeModelToObjectRoot(wm.WriteModel),
		TokenID:    wm.TokenID,
		Expiration: wm.ExpirationDate,
	}, base64.RawURLEncoding.EncodeToString(encrypted), nil
}

func readModelToWebAuthNTokens(readModel HumanWebAuthNTokensReadModel) []*domain.WebAuthNToken {
	tokens := make([]*domain.WebAuthNToken, len(readModel.GetWebAuthNTokens()))
	for i, token := range readModel.GetWebAuthNTokens() {
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
		RPID:              wm.RPID,
	}
}

func authRequestDomainToAuthRequestInfo(authRequest *domain.AuthRequest) *user.AuthRequestInfo {
	if authRequest == nil {
		return nil
	}
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

func writeModelToPasswordlessInitCode(initCodeModel *HumanPasswordlessInitCodeWriteModel, code string) *domain.PasswordlessInitCode {
	return &domain.PasswordlessInitCode{
		ObjectRoot: writeModelToObjectRoot(initCodeModel.WriteModel),
		CodeID:     initCodeModel.CodeID,
		Code:       code,
		Expiration: initCodeModel.Expiration,
		State:      initCodeModel.State,
	}
}

func writeModelToUserMetadata(wm *UserMetadataWriteModel) *domain.Metadata {
	return &domain.Metadata{
		ObjectRoot: writeModelToObjectRoot(wm.WriteModel),
		Key:        wm.Key,
		Value:      wm.Value,
		State:      wm.State,
	}
}
