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

func idpViewFromModel(idp *iam_model.IdpConfigView) *admin.IdpView {
	creationDate, err := ptypes.TimestampProto(idp.CreationDate)
	logging.Log("GRPC-8dju8").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(idp.ChangeDate)
	logging.Log("GRPC-Dsj8i").OnError(err).Debug("date parse failed")

	return &admin.IdpView{
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

func idpConfigViewFromModel(idp *iam_model.IdpConfigView) *admin.IdpView_OidcConfig {
	if idp.IsOidc {
		return &admin.IdpView_OidcConfig{
			OidcConfig: oidcIdpConfigViewFromModel(idp),
		}
	}
	return nil
}

func oidcIdpConfigViewFromModel(idp *iam_model.IdpConfigView) *admin.OidcIdpConfigView {
	return &admin.OidcIdpConfigView{
		ClientId: idp.OidcClientID,
		Issuer:   idp.OidcIssuer,
		Scopes:   idp.OidcScopes,
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

func idpConfigSearchRequestToModel(request *admin.IdpSearchRequest) *iam_model.IdpConfigSearchRequest {
	return &iam_model.IdpConfigSearchRequest{
		Limit:   request.Limit,
		Offset:  request.Offset,
		Queries: idpConfigSearchQueriesToModel(request.Queries),
	}
}

func idpConfigSearchQueriesToModel(queries []*admin.IdpSearchQuery) []*iam_model.IdpConfigSearchQuery {
	modelQueries := make([]*iam_model.IdpConfigSearchQuery, len(queries))
	for i, query := range queries {
		modelQueries[i] = idpConfigSearchQueryToModel(query)
	}

	return modelQueries
}

func idpConfigSearchQueryToModel(query *admin.IdpSearchQuery) *iam_model.IdpConfigSearchQuery {
	return &iam_model.IdpConfigSearchQuery{
		Key:    idpConfigSearchKeyToModel(query.Key),
		Method: searchMethodToModel(query.Method),
		Value:  query.Value,
	}
}

func idpConfigSearchKeyToModel(key admin.IdpSearchKey) iam_model.IdpConfigSearchKey {
	switch key {
	case admin.IdpSearchKey_IDPSEARCHKEY_IDP_CONFIG_ID:
		return iam_model.IdpConfigSearchKeyIdpConfigID
	case admin.IdpSearchKey_IDPSEARCHKEY_NAME:
		return iam_model.IdpConfigSearchKeyName
	default:
		return iam_model.IdpConfigSearchKeyUnspecified
	}
}

func idpConfigSearchResponseFromModel(resp *iam_model.IdpConfigSearchResponse) *admin.IdpSearchResponse {
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

func idpConfigsFromView(viewIdps []*iam_model.IdpConfigView) []*admin.IdpView {
	idps := make([]*admin.IdpView, len(viewIdps))
	for i, idp := range viewIdps {
		idps[i] = idpViewFromModel(idp)
	}
	return idps
}
