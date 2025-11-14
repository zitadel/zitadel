package command

import (
	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/domain"
)

func projectWriteModelToProject(writeModel *ProjectWriteModel) *domain.Project {
	return &domain.Project{
		ObjectRoot:             writeModelToObjectRoot(writeModel.WriteModel),
		Name:                   writeModel.Name,
		ProjectRoleAssertion:   writeModel.ProjectRoleAssertion,
		ProjectRoleCheck:       writeModel.ProjectRoleCheck,
		HasProjectCheck:        writeModel.HasProjectCheck,
		PrivateLabelingSetting: writeModel.PrivateLabelingSetting,
	}
}

func projectGrantWriteModelToProjectGrant(writeModel *ProjectGrantWriteModel) *domain.ProjectGrant {
	return &domain.ProjectGrant{
		ObjectRoot:   writeModelToObjectRoot(writeModel.WriteModel),
		GrantID:      writeModel.GrantID,
		GrantedOrgID: writeModel.GrantedOrgID,
		RoleKeys:     writeModel.RoleKeys,
		State:        writeModel.State,
	}
}

func oidcWriteModelToOIDCConfig(writeModel *OIDCApplicationWriteModel) *domain.OIDCApp {
	return &domain.OIDCApp{
		ObjectRoot:               writeModelToObjectRoot(writeModel.WriteModel),
		AppID:                    writeModel.AppID,
		AppName:                  writeModel.AppName,
		State:                    writeModel.State,
		ClientID:                 writeModel.ClientID,
		RedirectUris:             writeModel.RedirectUris,
		ResponseTypes:            writeModel.ResponseTypes,
		GrantTypes:               writeModel.GrantTypes,
		ApplicationType:          gu.Ptr(writeModel.ApplicationType),
		AuthMethodType:           gu.Ptr(writeModel.AuthMethodType),
		PostLogoutRedirectUris:   writeModel.PostLogoutRedirectUris,
		OIDCVersion:              gu.Ptr(writeModel.OIDCVersion),
		DevMode:                  gu.Ptr(writeModel.DevMode),
		AccessTokenType:          gu.Ptr(writeModel.AccessTokenType),
		AccessTokenRoleAssertion: gu.Ptr(writeModel.AccessTokenRoleAssertion),
		IDTokenRoleAssertion:     gu.Ptr(writeModel.IDTokenRoleAssertion),
		IDTokenUserinfoAssertion: gu.Ptr(writeModel.IDTokenUserinfoAssertion),
		ClockSkew:                gu.Ptr(writeModel.ClockSkew),
		AdditionalOrigins:        writeModel.AdditionalOrigins,
		SkipNativeAppSuccessPage: gu.Ptr(writeModel.SkipNativeAppSuccessPage),
		BackChannelLogoutURI:     gu.Ptr(writeModel.BackChannelLogoutURI),
		LoginVersion:             gu.Ptr(writeModel.LoginVersion),
		LoginBaseURI:             gu.Ptr(writeModel.LoginBaseURI),
	}
}

func samlWriteModelToSAMLConfig(writeModel *SAMLApplicationWriteModel) *domain.SAMLApp {
	return &domain.SAMLApp{
		ObjectRoot:   writeModelToObjectRoot(writeModel.WriteModel),
		AppID:        writeModel.AppID,
		AppName:      writeModel.AppName,
		State:        writeModel.State,
		Metadata:     writeModel.Metadata,
		MetadataURL:  gu.Ptr(writeModel.MetadataURL),
		EntityID:     writeModel.EntityID,
		LoginVersion: gu.Ptr(writeModel.LoginVersion),
		LoginBaseURI: gu.Ptr(writeModel.LoginBaseURI),
	}
}

func apiWriteModelToAPIConfig(writeModel *APIApplicationWriteModel) *domain.APIApp {
	return &domain.APIApp{
		ObjectRoot:     writeModelToObjectRoot(writeModel.WriteModel),
		AppID:          writeModel.AppID,
		AppName:        writeModel.AppName,
		State:          writeModel.State,
		ClientID:       writeModel.ClientID,
		AuthMethodType: writeModel.AuthMethodType,
	}
}

func memberWriteModelToProjectGrantMember(writeModel *ProjectGrantMemberWriteModel) *domain.ProjectGrantMember {
	return &domain.ProjectGrantMember{
		ObjectRoot: writeModelToObjectRoot(writeModel.WriteModel),
		Roles:      writeModel.Roles,
		GrantID:    writeModel.GrantID,
		UserID:     writeModel.UserID,
	}
}

func applicationKeyWriteModelToKey(wm *ApplicationKeyWriteModel) *domain.ApplicationKey {
	return &domain.ApplicationKey{
		ObjectRoot:     writeModelToObjectRoot(wm.WriteModel),
		ApplicationID:  wm.AppID,
		ClientID:       wm.ClientID,
		KeyID:          wm.KeyID,
		Type:           wm.KeyType,
		ExpirationDate: wm.ExpirationDate,
	}
}
