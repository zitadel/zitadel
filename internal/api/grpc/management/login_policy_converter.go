package management

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/management"
)

func loginPolicyAddToModel(policy *management.LoginPolicyAdd) *iam_model.LoginPolicy {
	return &iam_model.LoginPolicy{
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowExternalIdp:      policy.AllowExternalIdp,
		AllowRegister:         policy.AllowRegister,
	}
}
func loginPolicyToModel(policy *management.LoginPolicy) *iam_model.LoginPolicy {
	return &iam_model.LoginPolicy{
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowExternalIdp:      policy.AllowExternalIdp,
		AllowRegister:         policy.AllowRegister,
	}
}

func loginPolicyFromModel(policy *iam_model.LoginPolicy) *management.LoginPolicy {
	return &management.LoginPolicy{
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowExternalIdp:      policy.AllowExternalIdp,
		AllowRegister:         policy.AllowRegister,
	}
}

func loginPolicyViewFromModel(policy *iam_model.LoginPolicyView) *management.LoginPolicyView {
	return &management.LoginPolicyView{
		Default:               policy.Default,
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowExternalIdp:      policy.AllowExternalIDP,
		AllowRegister:         policy.AllowRegister,
	}
}

func idpProviderSearchRequestToModel(request *management.IdpProviderSearchRequest) *iam_model.IDPProviderSearchRequest {
	return &iam_model.IDPProviderSearchRequest{
		Limit:  request.Limit,
		Offset: request.Offset,
	}
}

func idpProviderSearchResponseFromModel(response *iam_model.IDPProviderSearchResponse) *management.IdpProviderSearchResponse {
	return &management.IdpProviderSearchResponse{
		Limit:       response.Limit,
		Offset:      response.Offset,
		TotalResult: response.TotalResult,
		Result:      idpProviderViewsFromModel(response.Result),
	}
}

func idpProviderToModel(provider *management.IdpProviderID) *iam_model.IDPProvider {
	return &iam_model.IDPProvider{
		IdpConfigID: provider.IdpConfigId,
		Type:        iam_model.IDPProviderTypeSystem,
	}
}

func idpProviderAddToModel(provider *management.IdpProviderAdd) *iam_model.IDPProvider {
	return &iam_model.IDPProvider{
		IdpConfigID: provider.IdpConfigId,
		Type:        idpProviderTypeToModel(provider.IdpProviderType),
	}
}

func idpProviderIDFromModel(provider *iam_model.IDPProvider) *management.IdpProviderID {
	return &management.IdpProviderID{
		IdpConfigId: provider.IdpConfigID,
	}
}

func idpProviderFromModel(provider *iam_model.IDPProvider) *management.IdpProvider {
	return &management.IdpProvider{
		IdpConfigId:      provider.IdpConfigID,
		IdpProvider_Type: idpProviderTypeFromModel(provider.Type),
	}
}

func idpProviderViewsFromModel(providers []*iam_model.IDPProviderView) []*management.IdpProviderView {
	converted := make([]*management.IdpProviderView, len(providers))
	for i, provider := range providers {
		converted[i] = idpProviderViewFromModel(provider)
	}

	return converted
}

func idpProviderViewFromModel(provider *iam_model.IDPProviderView) *management.IdpProviderView {
	return &management.IdpProviderView{
		IdpConfigId: provider.IDPConfigID,
		Name:        provider.Name,
		Type:        idpConfigTypeToModel(provider.IDPConfigType),
	}
}

func idpConfigTypeToModel(providerType iam_model.IdpConfigType) management.IdpType {
	switch providerType {
	case iam_model.IDPConfigTypeOIDC:
		return management.IdpType_IDPTYPE_OIDC
	case iam_model.IDPConfigTypeSAML:
		return management.IdpType_IDPTYPE_SAML
	default:
		return management.IdpType_IDPTYPE_UNSPECIFIED
	}
}

func idpProviderTypeToModel(providerType management.IdpProviderType) iam_model.IDPProviderType {
	switch providerType {
	case management.IdpProviderType_IDPPROVIDERTYPE_SYSTEM:
		return iam_model.IDPProviderTypeSystem
	case management.IdpProviderType_IDPPROVIDERTYPE_ORG:
		return iam_model.IDPProviderTypeOrg
	default:
		return iam_model.IDPProviderTypeSystem
	}
}

func idpProviderTypeFromModel(providerType iam_model.IDPProviderType) management.IdpProviderType {
	switch providerType {
	case iam_model.IDPProviderTypeSystem:
		return management.IdpProviderType_IDPPROVIDERTYPE_SYSTEM
	case iam_model.IDPProviderTypeOrg:
		return management.IdpProviderType_IDPPROVIDERTYPE_ORG
	default:
		return management.IdpProviderType_IDPPROVIDERTYPE_UNSPECIFIED
	}
}
