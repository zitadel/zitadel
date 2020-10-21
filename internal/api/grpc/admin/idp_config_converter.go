package admin

import (
	"github.com/caos/logging"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
)

func createOidcIdpToModel(idp *admin.OidcIdpConfigCreate) *iam_model.IDPConfig {
	return &iam_model.IDPConfig{
		Name:        idp.Name,
		StylingType: idpConfigStylingTypeToModel(idp.StylingType),
		Type:        iam_model.IDPConfigTypeOIDC,
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

func updateIdpToModel(idp *admin.IdpUpdate) *iam_model.IDPConfig {
	return &iam_model.IDPConfig{
		IDPConfigID: idp.Id,
		Name:        idp.Name,
		StylingType: idpConfigStylingTypeToModel(idp.StylingType),
	}
}

func updateOidcIdpToModel(idp *admin.OidcIdpConfigUpdate) *iam_model.OIDCIDPConfig {
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

func idpFromModel(idp *iam_model.IDPConfig) *admin.Idp {
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
		StylingType:  idpConfigStylingTypeFromModel(idp.StylingType),
		State:        idpConfigStateFromModel(idp.State),
		IdpConfig:    idpConfigFromModel(idp),
	}
}

func idpViewFromModel(idp *iam_model.IDPConfigView) *admin.IdpView {
	creationDate, err := ptypes.TimestampProto(idp.CreationDate)
	logging.Log("GRPC-8dju8").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(idp.ChangeDate)
	logging.Log("GRPC-Dsj8i").OnError(err).Debug("date parse failed")

	return &admin.IdpView{
		Id:            idp.IDPConfigID,
		CreationDate:  creationDate,
		ChangeDate:    changeDate,
		Sequence:      idp.Sequence,
		Name:          idp.Name,
		StylingType:   idpConfigStylingTypeFromModel(idp.StylingType),
		State:         idpConfigStateFromModel(idp.State),
		IdpConfigView: idpConfigViewFromModel(idp),
	}
}

func idpConfigFromModel(idp *iam_model.IDPConfig) *admin.Idp_OidcConfig {
	if idp.Type == iam_model.IDPConfigTypeOIDC {
		return &admin.Idp_OidcConfig{
			OidcConfig: oidcIdpConfigFromModel(idp.OIDCConfig),
		}
	}
	return nil
}

func oidcIdpConfigFromModel(idp *iam_model.OIDCIDPConfig) *admin.OidcIdpConfig {
	return &admin.OidcIdpConfig{
		ClientId: idp.ClientID,
		Issuer:   idp.Issuer,
		Scopes:   idp.Scopes,
	}
}

func idpConfigViewFromModel(idp *iam_model.IDPConfigView) *admin.IdpView_OidcConfig {
	if idp.IsOIDC {
		return &admin.IdpView_OidcConfig{
			OidcConfig: oidcIdpConfigViewFromModel(idp),
		}
	}
	return nil
}

func oidcIdpConfigViewFromModel(idp *iam_model.IDPConfigView) *admin.OidcIdpConfigView {
	return &admin.OidcIdpConfigView{
		ClientId:              idp.OIDCClientID,
		Issuer:                idp.OIDCIssuer,
		Scopes:                idp.OIDCScopes,
		IdpDisplayNameMapping: oidcMappingFieldFromModel(idp.OIDCIDPDisplayNameMapping),
		UsernameMapping:       oidcMappingFieldFromModel(idp.OIDCUsernameMapping),
	}
}

func idpConfigStateFromModel(state iam_model.IDPConfigState) admin.IdpState {
	switch state {
	case iam_model.IDPConfigStateActive:
		return admin.IdpState_IDPCONFIGSTATE_ACTIVE
	case iam_model.IDPConfigStateInactive:
		return admin.IdpState_IDPCONFIGSTATE_INACTIVE
	default:
		return admin.IdpState_IDPCONFIGSTATE_UNSPECIFIED
	}
}

func oidcMappingFieldFromModel(field iam_model.OIDCMappingField) admin.OIDCMappingField {
	switch field {
	case iam_model.OIDCMappingFieldPreferredLoginName:
		return admin.OIDCMappingField_OIDCMAPPINGFIELD_PREFERRED_USERNAME
	case iam_model.OIDCMappingFieldEmail:
		return admin.OIDCMappingField_OIDCMAPPINGFIELD_EMAIL
	default:
		return admin.OIDCMappingField_OIDCMAPPINGFIELD_UNSPECIFIED
	}
}

func oidcMappingFieldToModel(field admin.OIDCMappingField) iam_model.OIDCMappingField {
	switch field {
	case admin.OIDCMappingField_OIDCMAPPINGFIELD_PREFERRED_USERNAME:
		return iam_model.OIDCMappingFieldPreferredLoginName
	case admin.OIDCMappingField_OIDCMAPPINGFIELD_EMAIL:
		return iam_model.OIDCMappingFieldEmail
	default:
		return iam_model.OIDCMappingFieldUnspecified
	}
}

func idpConfigSearchRequestToModel(request *admin.IdpSearchRequest) *iam_model.IDPConfigSearchRequest {
	return &iam_model.IDPConfigSearchRequest{
		Limit:   request.Limit,
		Offset:  request.Offset,
		Queries: idpConfigSearchQueriesToModel(request.Queries),
	}
}

func idpConfigSearchQueriesToModel(queries []*admin.IdpSearchQuery) []*iam_model.IDPConfigSearchQuery {
	modelQueries := make([]*iam_model.IDPConfigSearchQuery, len(queries))
	for i, query := range queries {
		modelQueries[i] = idpConfigSearchQueryToModel(query)
	}

	return modelQueries
}

func idpConfigSearchQueryToModel(query *admin.IdpSearchQuery) *iam_model.IDPConfigSearchQuery {
	return &iam_model.IDPConfigSearchQuery{
		Key:    idpConfigSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func idpConfigSearchKeyToModel(key admin.IdpSearchKey) iam_model.IDPConfigSearchKey {
	switch key {
	case admin.IdpSearchKey_IDPSEARCHKEY_IDP_CONFIG_ID:
		return iam_model.IDPConfigSearchKeyIdpConfigID
	case admin.IdpSearchKey_IDPSEARCHKEY_NAME:
		return iam_model.IDPConfigSearchKeyName
	default:
		return iam_model.IDPConfigSearchKeyUnspecified
	}
}

func idpConfigSearchResponseFromModel(resp *iam_model.IDPConfigSearchResponse) *admin.IdpSearchResponse {
	timestamp, err := ptypes.TimestampProto(resp.Timestamp)
	logging.Log("GRPC-KSi8c").OnError(err).Debug("date parse failed")
	return &admin.IdpSearchResponse{
		Limit:             resp.Limit,
		Offset:            resp.Offset,
		TotalResult:       resp.TotalResult,
		Result:            idpConfigsFromView(resp.Result),
		ProcessedSequence: resp.Sequence,
		ViewTimestamp:     timestamp,
	}
}

func idpConfigsFromView(viewIdps []*iam_model.IDPConfigView) []*admin.IdpView {
	idps := make([]*admin.IdpView, len(viewIdps))
	for i, idp := range viewIdps {
		idps[i] = idpViewFromModel(idp)
	}
	return idps
}

func idpConfigStylingTypeFromModel(stylingType iam_model.IDPStylingType) admin.IdpStylingType {
	switch stylingType {
	case iam_model.IDPStylingTypeGoogle:
		return admin.IdpStylingType_IDPSTYLINGTYPE_GOOGLE
	default:
		return admin.IdpStylingType_IDPSTYLINGTYPE_UNSPECIFIED
	}
}

func idpConfigStylingTypeToModel(stylingType admin.IdpStylingType) iam_model.IDPStylingType {
	switch stylingType {
	case admin.IdpStylingType_IDPSTYLINGTYPE_GOOGLE:
		return iam_model.IDPStylingTypeGoogle
	default:
		return iam_model.IDPStylingTypeUnspecified
	}
}
