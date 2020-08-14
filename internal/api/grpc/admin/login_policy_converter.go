package admin

import (
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
)

func loginPolicyToModel(policy *admin.DefaultLoginPolicy) *iam_model.LoginPolicy {
	return &iam_model.LoginPolicy{
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowExternalIdp:      policy.AllowExternalIdp,
		AllowRegister:         policy.AllowRegister,
	}
}

func loginPolicyFromModel(policy *iam_model.LoginPolicy) *admin.DefaultLoginPolicy {
	return &admin.DefaultLoginPolicy{
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowExternalIdp:      policy.AllowExternalIdp,
		AllowRegister:         policy.AllowRegister,
	}
}

func loginPolicyViewFromModel(policy *iam_model.LoginPolicyView) *admin.DefaultLoginPolicyView {
	return &admin.DefaultLoginPolicyView{
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowExternalIdp:      policy.AllowExternalIdp,
		AllowRegister:         policy.AllowRegister,
	}
}

func idpProviderSearchRequestToModel(request *admin.IdpProviderSearchRequest) *iam_model.IdpProviderSearchRequest {
	return &iam_model.IdpProviderSearchRequest{
		Limit:  request.Limit,
		Offset: request.Offset,
	}
}

func idpProviderSearchResponseFromModel(response *iam_model.IdpProviderSearchResponse) *admin.IdpProviderSearchResponse {
	return &admin.IdpProviderSearchResponse{
		Limit:       response.Limit,
		Offset:      response.Offset,
		TotalResult: response.TotalResult,
		Result:      idpProviderViewsFromModel(response.Result),
	}
}

func idpProviderToModel(provider *admin.IdpProviderID) *iam_model.IdpProvider {
	return &iam_model.IdpProvider{
		IdpConfigID: provider.IdpConfigId,
		Type:        iam_model.IdpProviderTypeSystem,
	}
}

func idpProviderFromModel(provider *iam_model.IdpProvider) *admin.IdpProviderID {
	return &admin.IdpProviderID{
		IdpConfigId: provider.IdpConfigID,
	}
}

func idpProviderViewsFromModel(providers []*iam_model.IdpProviderView) []*admin.IdpProviderView {
	converted := make([]*admin.IdpProviderView, len(providers))
	for i, provider := range providers {
		converted[i] = idpProviderViewFromModel(provider)
	}

	return converted
}

func idpProviderViewFromModel(provider *iam_model.IdpProviderView) *admin.IdpProviderView {
	return &admin.IdpProviderView{
		IdpConfigId: provider.IdpConfigID,
		Name:        provider.Name,
		Type:        idpConfigTypeToModel(provider.IdpConfigType),
	}
}

func idpConfigTypeToModel(providerType iam_model.IdpConfigType) admin.IdpType {
	switch providerType {
	case iam_model.IDPConfigTypeOIDC:
		return admin.IdpType_IDPTYPE_OIDC
	case iam_model.IDPConfigTypeSAML:
		return admin.IdpType_IDPTYPE_SAML
	default:
		return admin.IdpType_IDPTYPE_UNSPECIFIED
	}
}
