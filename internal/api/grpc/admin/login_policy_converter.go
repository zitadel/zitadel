package admin

import (
	"github.com/caos/logging"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
)

func loginPolicyToModel(policy *admin.DefaultLoginPolicyRequest) *iam_model.LoginPolicy {
	return &iam_model.LoginPolicy{
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowExternalIdp:      policy.AllowExternalIdp,
		AllowRegister:         policy.AllowRegister,
		ForceMFA:              policy.ForceMfa,
	}
}

func loginPolicyFromModel(policy *iam_model.LoginPolicy) *admin.DefaultLoginPolicy {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-3Fsm9").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-5Gsko").OnError(err).Debug("date parse failed")

	return &admin.DefaultLoginPolicy{
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowExternalIdp:      policy.AllowExternalIdp,
		AllowRegister:         policy.AllowRegister,
		ForceMfa:              policy.ForceMFA,
		CreationDate:          creationDate,
		ChangeDate:            changeDate,
	}
}

func loginPolicyViewFromModel(policy *iam_model.LoginPolicyView) *admin.DefaultLoginPolicyView {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-3Gk9s").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-6Jlos").OnError(err).Debug("date parse failed")

	return &admin.DefaultLoginPolicyView{
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowExternalIdp:      policy.AllowExternalIDP,
		AllowRegister:         policy.AllowRegister,
		ForceMfa:              policy.ForceMFA,
		CreationDate:          creationDate,
		ChangeDate:            changeDate,
	}
}

func idpProviderSearchRequestToModel(request *admin.IdpProviderSearchRequest) *iam_model.IDPProviderSearchRequest {
	return &iam_model.IDPProviderSearchRequest{
		Limit:  request.Limit,
		Offset: request.Offset,
	}
}

func idpProviderSearchResponseFromModel(response *iam_model.IDPProviderSearchResponse) *admin.IdpProviderSearchResponse {
	return &admin.IdpProviderSearchResponse{
		Limit:       response.Limit,
		Offset:      response.Offset,
		TotalResult: response.TotalResult,
		Result:      idpProviderViewsFromModel(response.Result),
	}
}

func idpProviderToModel(provider *admin.IdpProviderID) *iam_model.IDPProvider {
	return &iam_model.IDPProvider{
		IdpConfigID: provider.IdpConfigId,
		Type:        iam_model.IDPProviderTypeSystem,
	}
}

func idpProviderFromModel(provider *iam_model.IDPProvider) *admin.IdpProviderID {
	return &admin.IdpProviderID{
		IdpConfigId: provider.IdpConfigID,
	}
}

func idpProviderViewsFromModel(providers []*iam_model.IDPProviderView) []*admin.IdpProviderView {
	converted := make([]*admin.IdpProviderView, len(providers))
	for i, provider := range providers {
		converted[i] = idpProviderViewFromModel(provider)
	}

	return converted
}

func idpProviderViewFromModel(provider *iam_model.IDPProviderView) *admin.IdpProviderView {
	return &admin.IdpProviderView{
		IdpConfigId: provider.IDPConfigID,
		Name:        provider.Name,
		Type:        idpConfigTypeToModel(provider.IDPConfigType),
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

func softwareMFAResultFromModel(result *iam_model.SoftwareMFASearchResponse) *admin.SoftwareMFAResult {
	converted := make([]admin.SoftwareMFAType, len(result.Result))
	for i, mfaType := range result.Result {
		converted[i] = softwareMFATypeFromModel(mfaType)
	}
	return &admin.SoftwareMFAResult{
		Mfas: converted,
	}
}

func softwareMFAFromModel(mfaType iam_model.SoftwareMFAType) *admin.SoftwareMFA {
	return &admin.SoftwareMFA{
		Mfa: softwareMFATypeFromModel(mfaType),
	}
}

func softwareMFATypeFromModel(mfaType iam_model.SoftwareMFAType) admin.SoftwareMFAType {
	switch mfaType {
	case iam_model.SoftwareMFATypeOTP:
		return admin.SoftwareMFAType_SOFTWAREMFATYPE_OTP
	default:
		return admin.SoftwareMFAType_SOFTWAREMFATYPE_UNSPECIFIED
	}
}

func softwareMFATypeToModel(mfaType *admin.SoftwareMFA) iam_model.SoftwareMFAType {
	switch mfaType.Mfa {
	case admin.SoftwareMFAType_SOFTWAREMFATYPE_OTP:
		return iam_model.SoftwareMFATypeOTP
	default:
		return iam_model.SoftwareMFATypeUnspecified
	}
}

func hardwareMFAResultFromModel(result *iam_model.HardwareMFASearchResponse) *admin.HardwareMFAResult {
	converted := make([]admin.HardwareMFAType, len(result.Result))
	for i, mfaType := range result.Result {
		converted[i] = hardwareMFATypeFromModel(mfaType)
	}
	return &admin.HardwareMFAResult{
		Mfas: converted,
	}
}

func hardwareMFAFromModel(mfaType iam_model.HardwareMFAType) *admin.HardwareMFA {
	return &admin.HardwareMFA{
		Mfa: hardwareMFATypeFromModel(mfaType),
	}
}

func hardwareMFATypeFromModel(mfaType iam_model.HardwareMFAType) admin.HardwareMFAType {
	switch mfaType {
	case iam_model.HardwareMFATypeU2F:
		return admin.HardwareMFAType_HARDWAREMFATYPE_U2F
	default:
		return admin.HardwareMFAType_HARDWAREMFATYPE_UNSPECIFIED
	}
}

func hardwareMFATypeToModel(mfaType *admin.HardwareMFA) iam_model.HardwareMFAType {
	switch mfaType.Mfa {
	case admin.HardwareMFAType_HARDWAREMFATYPE_U2F:
		return iam_model.HardwareMFATypeU2F
	default:
		return iam_model.HardwareMFATypeUnspecified
	}
}
