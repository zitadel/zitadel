package admin

import (
	"github.com/caos/logging"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
)

func createOidcIdpToModel(idp *admin.OidcIdpConfigCreate) *iam_model.IdpConfig {
	return &iam_model.IdpConfig{
		Name:    idp.Name,
		LogoSrc: idp.LogoSrc,
		Type:    iam_model.IDPConfigTypeOIDC,
		OIDCConfig: &iam_model.OidcIdpConfig{
			ClientID:           idp.ClientId,
			ClientSecretString: idp.ClientSecret,
			Issuer:             idp.Issuer,
			Scopes:             idp.Scopes,
		},
	}
}

func updateIdpToModel(idp *admin.IdpUpdate) *iam_model.IdpConfig {
	return &iam_model.IdpConfig{
		IDPConfigID: idp.Id,
		Name:        idp.Name,
		LogoSrc:     idp.LogoSrc,
	}
}

func updateOidcIdpToModel(idp *admin.OidcIdpConfigUpdate) *iam_model.OidcIdpConfig {
	return &iam_model.OidcIdpConfig{
		IDPConfigID:        idp.IdpId,
		ClientID:           idp.ClientId,
		ClientSecretString: idp.ClientSecret,
		Issuer:             idp.Issuer,
		Scopes:             idp.Scopes,
	}
}

func idpFromModel(idp *iam_model.IdpConfig) *admin.Idp {
	creationDate, err := ptypes.TimestampProto(idp.CreationDate)
	logging.Log("GRPC-8dju8").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(idp.ChangeDate)
	logging.Log("GRPC-Dsj8i").OnError(err).Debug("date parse failed")

	return &admin.Idp{
		Id:           idp.IDPConfigID,
		CreationDate: creationDate,
		ChangeDate:   changeDate,
		Sequence:     idp.Sequence,
		Name:         idp.Name,
		LogoSrc:      idp.LogoSrc,
		State:        idpConfigStateFromModel(idp.State),
		IdpConfig:    idpConfigFromModel(idp),
	}
}

func idpConfigFromModel(idp *iam_model.IdpConfig) *admin.Idp_OidcConfig {
	if idp.Type == iam_model.IDPConfigTypeOIDC {
		return &admin.Idp_OidcConfig{
			OidcConfig: oidcIdpConfigFromModel(idp.OIDCConfig),
		}
	}
	return nil
}

func oidcIdpConfigFromModel(idp *iam_model.OidcIdpConfig) *admin.OidcIdpConfig {
	return &admin.OidcIdpConfig{
		ClientId: idp.ClientID,
		Issuer:   idp.Issuer,
		Scopes:   idp.Scopes,
	}
}

func idpConfigStateFromModel(state iam_model.IdpConfigState) admin.IdpState {
	switch state {
	case iam_model.IdpConfigStateActive:
		return admin.IdpState_IDPCONFIGSTATE_ACTIVE
	case iam_model.IdpConfigStateInactive:
		return admin.IdpState_IDPCONFIGSTATE_INACTIVE
	default:
		return admin.IdpState_IDPCONFIGSTATE_UNSPECIFIED
	}
}
