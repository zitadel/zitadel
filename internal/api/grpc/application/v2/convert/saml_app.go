package convert

import (
	"github.com/muhlemmer/gu"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	"github.com/zitadel/zitadel/internal/query"
	"github.com/zitadel/zitadel/pkg/grpc/application/v2"
)

func CreateSAMLAppRequestToDomain(name, projectID string, req *application.CreateSAMLApplicationRequest) (*domain.SAMLApp, error) {
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

func UpdateSAMLAppConfigRequestToDomain(appID, projectID string, app *application.UpdateSAMLApplicationConfigurationRequest) (*domain.SAMLApp, error) {
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

func metasToDomain(metas application.MetaType) ([]byte, *string) {
	switch t := metas.(type) {
	case *application.UpdateSAMLApplicationConfigurationRequest_MetadataXml:
		return t.MetadataXml, nil
	case *application.UpdateSAMLApplicationConfigurationRequest_MetadataUrl:
		return nil, &t.MetadataUrl
	case nil:
		return nil, nil
	default:
		return nil, nil
	}
}

func appSAMLConfigToPb(samlApp *query.SAMLApp) application.IsApplicationConfiguration {
	if samlApp == nil {
		return &application.Application_SamlConfiguration{
			SamlConfiguration: &application.SAMLConfiguration{
				LoginVersion: &application.LoginVersion{},
			},
		}
	}

	return &application.Application_SamlConfiguration{
		SamlConfiguration: &application.SAMLConfiguration{
			MetadataXml:  samlApp.Metadata,
			MetadataUrl:  samlApp.MetadataURL,
			LoginVersion: loginVersionToPb(samlApp.LoginVersion, samlApp.LoginBaseURI),
		},
	}
}
