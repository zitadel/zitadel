package command

import (
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
		ApplicationType:          writeModel.ApplicationType,
		AuthMethodType:           writeModel.AuthMethodType,
		PostLogoutRedirectUris:   writeModel.PostLogoutRedirectUris,
		OIDCVersion:              writeModel.OIDCVersion,
		DevMode:                  writeModel.DevMode,
		AccessTokenType:          writeModel.AccessTokenType,
		AccessTokenRoleAssertion: writeModel.AccessTokenRoleAssertion,
		IDTokenRoleAssertion:     writeModel.IDTokenRoleAssertion,
		IDTokenUserinfoAssertion: writeModel.IDTokenUserinfoAssertion,
		ClockSkew:                writeModel.ClockSkew,
		AdditionalOrigins:        writeModel.AdditionalOrigins,
		SkipNativeAppSuccessPage: writeModel.SkipNativeAppSuccessPage,
		BackChannelLogoutURI:     writeModel.BackChannelLogoutURI,
		LoginVersion:             writeModel.LoginVersion,
		LoginBaseURI:             writeModel.LoginBaseURI,
	}
}

func samlWriteModelToSAMLConfig(writeModel *SAMLApplicationWriteModel) *domain.SAMLApp {
	return &domain.SAMLApp{
		ObjectRoot:  writeModelToObjectRoot(writeModel.WriteModel),
		AppID:       writeModel.AppID,
		AppName:     writeModel.AppName,
		State:       writeModel.State,
		Metadata:    writeModel.Metadata,
		MetadataURL: writeModel.MetadataURL,
		EntityID:    writeModel.EntityID,
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

func roleWriteModelToRole(writeModel *ProjectRoleWriteModel) *domain.ProjectRole {
	return &domain.ProjectRole{
		ObjectRoot:  writeModelToObjectRoot(writeModel.WriteModel),
		Key:         writeModel.Key,
		DisplayName: writeModel.DisplayName,
		Group:       writeModel.Group,
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
