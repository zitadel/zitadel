package management

import (
	"github.com/caos/logging"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes"
)

func createOidcIdpToModel(idp *management.OidcIdpConfigCreate) *iam_model.IdpConfig {
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

func updateIdpToModel(idp *management.IdpUpdate) *iam_model.IdpConfig {
	return &iam_model.IdpConfig{
		IDPConfigID: idp.Id,
		Name:        idp.Name,
		LogoSrc:     idp.LogoSrc,
	}
}

func updateOidcIdpToModel(idp *management.OidcIdpConfigUpdate) *iam_model.OidcIdpConfig {
	return &iam_model.OidcIdpConfig{
		IDPConfigID:        idp.IdpId,
		ClientID:           idp.ClientId,
		ClientSecretString: idp.ClientSecret,
		Issuer:             idp.Issuer,
		Scopes:             idp.Scopes,
	}
}

func idpFromModel(idp *iam_model.IdpConfig) *management.Idp {
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

func idpViewFromModel(idp *iam_model.IdpConfigView) *management.IdpView {
	creationDate, err := ptypes.TimestampProto(idp.CreationDate)
	logging.Log("GRPC-8dju8").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(idp.ChangeDate)
	logging.Log("GRPC-Dsj8i").OnError(err).Debug("date parse failed")

	return &management.IdpView{
		Id:            idp.IdpConfigID,
		CreationDate:  creationDate,
		ChangeDate:    changeDate,
		Sequence:      idp.Sequence,
		Name:          idp.Name,
		LogoSrc:       idp.LogoSrc,
		State:         idpConfigStateFromModel(idp.State),
		IdpConfigView: idpConfigViewFromModel(idp),
	}
}

func idpConfigFromModel(idp *iam_model.IdpConfig) *management.Idp_OidcConfig {
	if idp.Type == iam_model.IDPConfigTypeOIDC {
		return &management.Idp_OidcConfig{
			OidcConfig: oidcIdpConfigFromModel(idp.OIDCConfig),
		}
	}
	return nil
}

func oidcIdpConfigFromModel(idp *iam_model.OidcIdpConfig) *management.OidcIdpConfig {
	return &management.OidcIdpConfig{
		ClientId: idp.ClientID,
		Issuer:   idp.Issuer,
		Scopes:   idp.Scopes,
	}
}

func idpConfigViewFromModel(idp *iam_model.IdpConfigView) *management.IdpView_OidcConfig {
	if idp.IsOidc {
		return &management.IdpView_OidcConfig{
			OidcConfig: oidcIdpConfigViewFromModel(idp),
		}
	}
	return nil
}

func oidcIdpConfigViewFromModel(idp *iam_model.IdpConfigView) *management.OidcIdpConfigView {
	return &management.OidcIdpConfigView{
		ClientId: idp.OidcClientID,
		Issuer:   idp.OidcIssuer,
		Scopes:   idp.OidcScopes,
	}
}

func idpConfigStateFromModel(state iam_model.IdpConfigState) management.IdpState {
	switch state {
	case iam_model.IdpConfigStateActive:
		return management.IdpState_IDPCONFIGSTATE_ACTIVE
	case iam_model.IdpConfigStateInactive:
		return management.IdpState_IDPCONFIGSTATE_INACTIVE
	default:
		return management.IdpState_IDPCONFIGSTATE_UNSPECIFIED
	}
}

func idpConfigSearchRequestToModel(request *management.IdpSearchRequest) *iam_model.IdpConfigSearchRequest {
	return &iam_model.IdpConfigSearchRequest{
		Limit:   request.Limit,
		Offset:  request.Offset,
		Queries: idpConfigSearchQueriesToModel(request.Queries),
	}
}

func idpConfigSearchQueriesToModel(queries []*management.IdpSearchQuery) []*iam_model.IdpConfigSearchQuery {
	modelQueries := make([]*iam_model.IdpConfigSearchQuery, len(queries))
	for i, query := range queries {
		modelQueries[i] = idpConfigSearchQueryToModel(query)
	}

	return modelQueries
}

func idpConfigSearchQueryToModel(query *management.IdpSearchQuery) *iam_model.IdpConfigSearchQuery {
	return &iam_model.IdpConfigSearchQuery{
		Key:    idpConfigSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func idpConfigSearchKeyToModel(key management.IdpSearchKey) iam_model.IdpConfigSearchKey {
	switch key {
	case management.IdpSearchKey_IDPSEARCHKEY_IDP_CONFIG_ID:
		return iam_model.IdpConfigSearchKeyIdpConfigID
	case management.IdpSearchKey_IDPSEARCHKEY_NAME:
		return iam_model.IdpConfigSearchKeyName
	default:
		return iam_model.IdpConfigSearchKeyUnspecified
	}
}

func idpConfigSearchResponseFromModel(resp *iam_model.IdpConfigSearchResponse) *management.IdpSearchResponse {
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

func idpConfigsFromView(viewIdps []*iam_model.IdpConfigView) []*management.IdpView {
	idps := make([]*management.IdpView, len(viewIdps))
	for i, idp := range viewIdps {
		idps[i] = idpViewFromModel(idp)
	}
	return idps
}
