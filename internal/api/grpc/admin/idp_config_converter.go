package admin

import (
	"github.com/caos/logging"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
)

func createOIDCIDPToDomain(idp *admin.OidcIdpConfigCreate) *domain.IDPConfig {
	return &domain.IDPConfig{
		Name:        idp.Name,
		StylingType: idpConfigStylingTypeToDomain(idp.StylingType),
		Type:        domain.IDPConfigTypeOIDC,
		OIDCConfig: &domain.OIDCIDPConfig{
			ClientID:              idp.ClientId,
			ClientSecretString:    idp.ClientSecret,
			Issuer:                idp.Issuer,
			Scopes:                idp.Scopes,
			IDPDisplayNameMapping: oidcMappingFieldToDomain(idp.IdpDisplayNameMapping),
			UsernameMapping:       oidcMappingFieldToDomain(idp.UsernameMapping),
		},
	}
}

func updateIdpToDomain(idp *admin.IdpUpdate) *domain.IDPConfig {
	return &domain.IDPConfig{
		IDPConfigID: idp.Id,
		Name:        idp.Name,
		StylingType: idpConfigStylingTypeToDomain(idp.StylingType),
	}
}

func updateOIDCIDPToDomain(idp *admin.OidcIdpConfigUpdate) *domain.OIDCIDPConfig {
	return &domain.OIDCIDPConfig{
		IDPConfigID:           idp.IdpId,
		ClientID:              idp.ClientId,
		ClientSecretString:    idp.ClientSecret,
		Issuer:                idp.Issuer,
		Scopes:                idp.Scopes,
		IDPDisplayNameMapping: oidcMappingFieldToDomain(idp.IdpDisplayNameMapping),
		UsernameMapping:       oidcMappingFieldToDomain(idp.UsernameMapping),
	}
}

func idpFromDomain(idp *domain.IDPConfig) *admin.Idp {
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
		StylingType:  idpConfigStylingTypeFromDomain(idp.StylingType),
		State:        idpConfigStateFromDomain(idp.State),
		IdpConfig:    idpConfigFromDomain(idp),
	}
}

func idpViewFromDomain(idp *domain.IDPConfigView) *admin.IdpView {
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
		StylingType:   admin.IdpStylingType(idp.StylingType),
		State:         admin.IdpState(idp.State),
		IdpConfigView: idpConfigViewFromDomain(idp),
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
		StylingType:   admin.IdpStylingType(idp.StylingType),
		State:         admin.IdpState(idp.State),
		IdpConfigView: idpConfigViewFromModel(idp),
	}
}

func idpConfigFromDomain(idp *domain.IDPConfig) *admin.Idp_OidcConfig {
	if idp.Type == domain.IDPConfigTypeOIDC {
		return &admin.Idp_OidcConfig{
			OidcConfig: oidcIDPConfigFromDomain(idp.OIDCConfig),
		}
	}
	return nil
}

func oidcIDPConfigFromDomain(idp *domain.OIDCIDPConfig) *admin.OidcIdpConfig {
	return &admin.OidcIdpConfig{
		ClientId: idp.ClientID,
		Issuer:   idp.Issuer,
		Scopes:   idp.Scopes,
	}
}

func idpConfigViewFromDomain(idp *domain.IDPConfigView) *admin.IdpView_OidcConfig {
	if idp.IsOIDC {
		return &admin.IdpView_OidcConfig{
			OidcConfig: oidcIdpConfigViewFromDomain(idp),
		}
	}
	return nil
}

func idpConfigViewFromModel(idp *iam_model.IDPConfigView) *admin.IdpView_OidcConfig {
	if idp.IsOIDC {
		return &admin.IdpView_OidcConfig{
			OidcConfig: oidcIdpConfigViewFromModel(idp),
		}
	}
	return nil
}

func oidcIdpConfigViewFromDomain(idp *domain.IDPConfigView) *admin.OidcIdpConfigView {
	return &admin.OidcIdpConfigView{
		ClientId:              idp.OIDCClientID,
		Issuer:                idp.OIDCIssuer,
		Scopes:                idp.OIDCScopes,
		IdpDisplayNameMapping: oidcMappingFieldFromDomain(idp.OIDCIDPDisplayNameMapping),
		UsernameMapping:       oidcMappingFieldFromDomain(idp.OIDCUsernameMapping),
	}
}

func oidcIdpConfigViewFromModel(idp *iam_model.IDPConfigView) *admin.OidcIdpConfigView {
	return &admin.OidcIdpConfigView{
		ClientId:              idp.OIDCClientID,
		Issuer:                idp.OIDCIssuer,
		Scopes:                idp.OIDCScopes,
		IdpDisplayNameMapping: admin.OIDCMappingField(idp.OIDCIDPDisplayNameMapping),
		UsernameMapping:       admin.OIDCMappingField(idp.OIDCUsernameMapping),
	}
}

func idpConfigStateFromDomain(state domain.IDPConfigState) admin.IdpState {
	switch state {
	case domain.IDPConfigStateActive:
		return admin.IdpState_IDPCONFIGSTATE_ACTIVE
	case domain.IDPConfigStateInactive:
		return admin.IdpState_IDPCONFIGSTATE_INACTIVE
	default:
		return admin.IdpState_IDPCONFIGSTATE_UNSPECIFIED
	}
}

func oidcMappingFieldFromDomain(field domain.OIDCMappingField) admin.OIDCMappingField {
	switch field {
	case domain.OIDCMappingFieldPreferredLoginName:
		return admin.OIDCMappingField_OIDCMAPPINGFIELD_PREFERRED_USERNAME
	case domain.OIDCMappingFieldEmail:
		return admin.OIDCMappingField_OIDCMAPPINGFIELD_EMAIL
	default:
		return admin.OIDCMappingField_OIDCMAPPINGFIELD_UNSPECIFIED
	}
}

func oidcMappingFieldToDomain(field admin.OIDCMappingField) domain.OIDCMappingField {
	switch field {
	case admin.OIDCMappingField_OIDCMAPPINGFIELD_PREFERRED_USERNAME:
		return domain.OIDCMappingFieldPreferredLoginName
	case admin.OIDCMappingField_OIDCMAPPINGFIELD_EMAIL:
		return domain.OIDCMappingFieldEmail
	default:
		return domain.OIDCMappingFieldPreferredLoginName
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

func idpConfigStylingTypeFromDomain(stylingType domain.IDPConfigStylingType) admin.IdpStylingType {
	switch stylingType {
	case domain.IDPConfigStylingTypeGoogle:
		return admin.IdpStylingType_IDPSTYLINGTYPE_GOOGLE
	default:
		return admin.IdpStylingType_IDPSTYLINGTYPE_UNSPECIFIED
	}
}

func idpConfigStylingTypeToDomain(stylingType admin.IdpStylingType) domain.IDPConfigStylingType {
	switch stylingType {
	case admin.IdpStylingType_IDPSTYLINGTYPE_GOOGLE:
		return domain.IDPConfigStylingTypeGoogle
	default:
		return domain.IDPConfigStylingTypeUnspecified
	}
}
