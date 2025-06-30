package convert

import (
	"github.com/muhlemmer/gu"

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
		MetadataURL:  gu.Ptr(req.GetMetadataUrl()),
		LoginVersion: loginVersion,
		LoginBaseURI: loginBaseURI,
	}, nil
}

func UpdateSAMLAppConfigRequestToDomain(appID, projectID string, app *app.UpdateSAMLApplicationConfigurationRequest) (*domain.SAMLApp, error) {
	loginVersion, loginBaseURI, err := loginVersionToDomain(app.GetLoginVersion())
	if err != nil {
		return nil, err
	}

	metasXML, metasURL := metasToDomain(app.GetMetadata())
	return &domain.SAMLApp{
		ObjectRoot: models.ObjectRoot{
			AggregateID: projectID,
		},
		AppID:        appID,
		Metadata:     metasXML,
		MetadataURL:  metasURL,
		LoginVersion: loginVersion,
		LoginBaseURI: loginBaseURI,
	}, nil
}

func metasToDomain(metas app.MetaType) ([]byte, *string) {
	switch t := metas.(type) {
	case *app.UpdateSAMLApplicationConfigurationRequest_MetadataXml:
		return t.MetadataXml, nil
	case *app.UpdateSAMLApplicationConfigurationRequest_MetadataUrl:
		return nil, &t.MetadataUrl
	case nil:
		return nil, nil
	default:
		return nil, nil
	}
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
