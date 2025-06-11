package convert

import (
	"net/url"

	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/query"
	app "github.com/zitadel/zitadel/pkg/grpc/app/v2beta"
)

func AppToPb(query_app *query.App) *app.Application {
	return &app.Application{
		Id:           query_app.ID,
		CreationDate: timestamppb.New(query_app.CreationDate),
		ChangeDate:   timestamppb.New(query_app.ChangeDate),
		State:        appStateToPb(query_app.State),
		Name:         query_app.Name,
		Config:       appConfigToPb(query_app),
	}
}

func appStateToPb(state domain.AppState) app.AppState {
	switch state {
	case domain.AppStateActive:
		return app.AppState_APP_STATE_ACTIVE
	case domain.AppStateInactive:
		return app.AppState_APP_STATE_INACTIVE
	default:
		return app.AppState_APP_STATE_UNSPECIFIED
	}
}

func appConfigToPb(app *query.App) app.ApplicationConfig {
	if app.OIDCConfig != nil {
		return appOIDCConfigToPb(app.OIDCConfig)
	}
	if app.SAMLConfig != nil {
		return appSAMLConfigToPb(app.SAMLConfig)
	}
	return appAPIConfigToPb(app.APIConfig)
}

func loginVersionToDomain(version *app.LoginVersion) (domain.LoginVersion, string, error) {
	switch v := version.GetVersion().(type) {
	case nil:
		return domain.LoginVersionUnspecified, "", nil
	case *app.LoginVersion_LoginV1:
		return domain.LoginVersion1, "", nil
	case *app.LoginVersion_LoginV2:
		_, err := url.Parse(v.LoginV2.GetBaseUri())
		return domain.LoginVersion2, v.LoginV2.GetBaseUri(), err
	default:
		return domain.LoginVersionUnspecified, "", nil
	}
}

func loginVersionToPb(version domain.LoginVersion, baseURI *string) *app.LoginVersion {
	switch version {
	case domain.LoginVersionUnspecified:
		return nil
	case domain.LoginVersion1:
		return &app.LoginVersion{Version: &app.LoginVersion_LoginV1{LoginV1: &app.LoginV1{}}}
	case domain.LoginVersion2:
		return &app.LoginVersion{Version: &app.LoginVersion_LoginV2{LoginV2: &app.LoginV2{BaseUri: baseURI}}}
	default:
		return nil
	}
}
