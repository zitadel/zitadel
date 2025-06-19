package convert

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func CreateSAMLAppRequestToDomain(name, projectID string, req *app.CreateSAMLApplicationRequest) (*domain.SAMLApp, error) {
	loginVersion, loginBaseURI, err := loginVersionToDomain(req.GetLoginVersion())
	if err != nil {
		return nil, err
	}
	return &domain.SAMLApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: projectID,
		},
		AppName:      name,
		Metadata:     req.GetMetadataXml(),
		MetadataURL:  req.GetMetadataUrl(),
		LoginVersion: loginVersion,
		LoginBaseURI: loginBaseURI,
	}, nil
}

func PatchSAMLAppConfigRequestToDomain(appID, projectID string, app *app.UpdateSAMLApplicationConfigurationRequest) (*domain.SAMLApp, error) {
	loginVersion, loginBaseURI, err := loginVersionToDomain(app.GetLoginVersion())
	if err != nil {
		return nil, err
	}
	return &domain.SAMLApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: projectID,
		},
		AppID:        appID,
		Metadata:     app.GetMetadataXml(),
		MetadataURL:  app.GetMetadataUrl(),
		LoginVersion: loginVersion,
		LoginBaseURI: loginBaseURI,
	}, nil
}

func appSAMLConfigToPb(samlApp *query.SAMLApp) app.ApplicationConfig {
	if samlApp == nil {
		return &app.Application_SamlConfig{
			SamlConfig: &app.SAMLConfig{
				Metadata:     &app.SAMLConfig_MetadataXml{},
				LoginVersion: &app.LoginVersion{},
			},
		}
	}

	return &app.Application_SamlConfig{
		SamlConfig: &app.SAMLConfig{
			Metadata:     &app.SAMLConfig_MetadataXml{MetadataXml: samlApp.Metadata},
			LoginVersion: loginVersionToPb(samlApp.LoginVersion, samlApp.LoginBaseURI),
		},
	}
}
