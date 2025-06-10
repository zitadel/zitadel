package convert

import (
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore/v1/models"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func CreateSAMLAppRequestToDomain(name, projectID string, req *app.CreateSAMLApplicationRequest) (*domain.SAMLApp, error) {
	loginVersion, loginBaseURI, err := LoginVersionToDomain(req.GetLoginVersion())
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

func PatchSAMLAppConfigRequestToDomain(appID, projectID string, app *app.PatchSAMLApplicationConfigurationRequest) (*domain.SAMLApp, error) {
	loginVersion, loginBaseURI, err := LoginVersionToDomain(app.GetLoginVersion())
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
