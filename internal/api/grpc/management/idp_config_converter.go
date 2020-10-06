package management

import (
	"github.com/caos/logging"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes"
)

func createOidcIdpToModel(idp *management.OidcIdpConfigCreate) *iam_model.IDPConfig {
	return &iam_model.IDPConfig{
		Name:    idp.Name,
		LogoSrc: idp.LogoSrc,
		Type:    iam_model.IDPConfigTypeOIDC,
		OIDCConfig: &iam_model.OIDCIDPConfig{
			ClientID:              idp.ClientId,
			ClientSecretString:    idp.ClientSecret,
			Issuer:                idp.Issuer,
			Scopes:                idp.Scopes,
			IDPDisplayNameMapping: oidcMappingFieldToModel(idp.IdpDisplayNameMapping),
			UsernameMapping:       oidcMappingFieldToModel(idp.UsernameMapping),
		},
	}
}

func updateIdpToModel(idp *management.IdpUpdate) *iam_model.IDPConfig {
	return &iam_model.IDPConfig{
		IDPConfigID: idp.Id,
		Name:        idp.Name,
		LogoSrc:     idp.LogoSrc,
	}
}

func updateOidcIdpToModel(idp *management.OidcIdpConfigUpdate) *iam_model.OIDCIDPConfig {
	return &iam_model.OIDCIDPConfig{
		IDPConfigID:           idp.IdpId,
		ClientID:              idp.ClientId,
		ClientSecretString:    idp.ClientSecret,
		Issuer:                idp.Issuer,
		Scopes:                idp.Scopes,
		IDPDisplayNameMapping: oidcMappingFieldToModel(idp.IdpDisplayNameMapping),
		UsernameMapping:       oidcMappingFieldToModel(idp.UsernameMapping),
	}
}

func idpFromModel(idp *iam_model.IDPConfig) *management.Idp {
	creationDate, err := ptypes.TimestampProto(idp.CreationDate)
	logging.Log("GRPC-8dju8").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(idp.ChangeDate)
	logging.Log("GRPC-Dsj8i").OnError(err).Debug("date parse failed")

	return &management.Idp{
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

func idpViewFromModel(idp *iam_model.IDPConfigView) *management.IdpView {
	creationDate, err := ptypes.TimestampProto(idp.CreationDate)
	logging.Log("GRPC-8dju8").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(idp.ChangeDate)
	logging.Log("GRPC-Dsj8i").OnError(err).Debug("date parse failed")

	return &management.IdpView{
		Id:            idp.IDPConfigID,
		CreationDate:  creationDate,
		ChangeDate:    changeDate,
		Sequence:      idp.Sequence,
		ProviderType:  idpProviderTypeFromModel(idp.IDPProviderType),
		Name:          idp.Name,
		LogoSrc:       idp.LogoSrc,
		State:         idpConfigStateFromModel(idp.State),
		IdpConfigView: idpConfigViewFromModel(idp),
	}
}

func idpConfigFromModel(idp *iam_model.IDPConfig) *management.Idp_OidcConfig {
	if idp.Type == iam_model.IDPConfigTypeOIDC {
		return &management.Idp_OidcConfig{
			OidcConfig: oidcIdpConfigFromModel(idp.OIDCConfig),
		}
	}
	return nil
}

func oidcIdpConfigFromModel(idp *iam_model.OIDCIDPConfig) *management.OidcIdpConfig {
	return &management.OidcIdpConfig{
		ClientId:              idp.ClientID,
		Issuer:                idp.Issuer,
		Scopes:                idp.Scopes,
		IdpDisplayNameMapping: oidcMappingFieldFromModel(idp.IDPDisplayNameMapping),
		UsernameMapping:       oidcMappingFieldFromModel(idp.UsernameMapping),
	}
}

func idpConfigViewFromModel(idp *iam_model.IDPConfigView) *management.IdpView_OidcConfig {
	if idp.IsOIDC {
		return &management.IdpView_OidcConfig{
			OidcConfig: oidcIdpConfigViewFromModel(idp),
		}
	}
	return nil
}

func oidcIdpConfigViewFromModel(idp *iam_model.IDPConfigView) *management.OidcIdpConfigView {
	return &management.OidcIdpConfigView{
		ClientId:              idp.OIDCClientID,
		Issuer:                idp.OIDCIssuer,
		Scopes:                idp.OIDCScopes,
		IdpDisplayNameMapping: oidcMappingFieldFromModel(idp.OIDCIDPDisplayNameMapping),
		UsernameMapping:       oidcMappingFieldFromModel(idp.OIDCUsernameMapping),
	}
}

func idpConfigStateFromModel(state iam_model.IDPConfigState) management.IdpState {
	switch state {
	case iam_model.IDPConfigStateActive:
		return management.IdpState_IDPCONFIGSTATE_ACTIVE
	case iam_model.IDPConfigStateInactive:
		return management.IdpState_IDPCONFIGSTATE_INACTIVE
	default:
		return management.IdpState_IDPCONFIGSTATE_UNSPECIFIED
	}
}

func idpConfigSearchRequestToModel(request *management.IdpSearchRequest) *iam_model.IDPConfigSearchRequest {
	return &iam_model.IDPConfigSearchRequest{
		Limit:   request.Limit,
		Offset:  request.Offset,
		Queries: idpConfigSearchQueriesToModel(request.Queries),
	}
}

func idpConfigSearchQueriesToModel(queries []*management.IdpSearchQuery) []*iam_model.IDPConfigSearchQuery {
	modelQueries := make([]*iam_model.IDPConfigSearchQuery, len(queries))
	for i, query := range queries {
		modelQueries[i] = idpConfigSearchQueryToModel(query)
	}

	return modelQueries
}

func idpConfigSearchQueryToModel(query *management.IdpSearchQuery) *iam_model.IDPConfigSearchQuery {
	return &iam_model.IDPConfigSearchQuery{
		Key:    idpConfigSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func idpConfigSearchKeyToModel(key management.IdpSearchKey) iam_model.IDPConfigSearchKey {
	switch key {
	case management.IdpSearchKey_IDPSEARCHKEY_IDP_CONFIG_ID:
		return iam_model.IDPConfigSearchKeyIdpConfigID
	case management.IdpSearchKey_IDPSEARCHKEY_NAME:
		return iam_model.IDPConfigSearchKeyName
	case management.IdpSearchKey_IDPSEARCHKEY_PROVIDER_TYPE:
		return iam_model.IDPConfigSearchKeyIdpProviderType
	default:
		return iam_model.IDPConfigSearchKeyUnspecified
	}
}

func idpConfigSearchResponseFromModel(resp *iam_model.IDPConfigSearchResponse) *management.IdpSearchResponse {
	timestamp, err := ptypes.TimestampProto(resp.Timestamp)
	logging.Log("GRPC-KSi8c").OnError(err).Debug("date parse failed")
	return &management.IdpSearchResponse{
		Limit:             resp.Limit,
		Offset:            resp.Offset,
		TotalResult:       resp.TotalResult,
		Result:            idpConfigsFromView(resp.Result),
		ProcessedSequence: resp.Sequence,
		ViewTimestamp:     timestamp,
	}
}

func idpConfigsFromView(viewIdps []*iam_model.IDPConfigView) []*management.IdpView {
	idps := make([]*management.IdpView, len(viewIdps))
	for i, idp := range viewIdps {
		idps[i] = idpViewFromModel(idp)
	}
	return idps
}

func oidcMappingFieldFromModel(field iam_model.OIDCMappingField) management.OIDCMappingField {
	switch field {
	case iam_model.OIDCMappingFieldPreferredLoginName:
		return management.OIDCMappingField_OIDCMAPPINGFIELD_PREFERRED_USERNAME
	case iam_model.OIDCMappingFieldEmail:
		return management.OIDCMappingField_OIDCMAPPINGFIELD_EMAIL
	default:
		return management.OIDCMappingField_OIDCMAPPINGFIELD_UNSPECIFIED
	}
}

func oidcMappingFieldToModel(field management.OIDCMappingField) iam_model.OIDCMappingField {
	switch field {
	case management.OIDCMappingField_OIDCMAPPINGFIELD_PREFERRED_USERNAME:
		return iam_model.OIDCMappingFieldPreferredLoginName
	case management.OIDCMappingField_OIDCMAPPINGFIELD_EMAIL:
		return iam_model.OIDCMappingFieldEmail
	default:
		return iam_model.OIDCMappingFieldUnspecified
	}
}
