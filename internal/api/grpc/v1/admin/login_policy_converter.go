package admin

import (
	"github.com/caos/logging"
	"github.com/caos/zitadel/interna/domain"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/pkg/grpc/admin"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func loginPolicyToDomain(policy *admin.DefaultLoginPolicyRequest) *domain.LoginPolicy {
	return &domain.LoginPolicy{
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowExternalIDP:      policy.AllowExternalIdp,
		AllowRegister:         policy.AllowRegister,
		ForceMFA:              policy.ForceMfa,
		PasswordlessType:      passwordlessTypeToDomain(policy.PasswordlessType),
	}
}

func loginPolicyFromDomain(policy *domain.LoginPolicy) *admin.DefaultLoginPolicy {
	return &admin.DefaultLoginPolicy{
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowExternalIdp:      policy.AllowExternalIDP,
		AllowRegister:         policy.AllowRegister,
		ForceMfa:              policy.ForceMFA,
		PasswordlessType:      passwordlessTypeFromDomain(policy.PasswordlessType),
		ChangeDate:            timestamppb.New(policy.ChangeDate),
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
		PasswordlessType:      admin.PasswordlessType(policy.PasswordlessType),
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

func idpProviderToDomain(provider *admin.IdpProviderID) *domain.IDPProvider {
	return &domain.IDPProvider{
		IDPConfigID: provider.IdpConfigId,
		Type:        domain.IdentityProviderTypeSystem,
	}
}

func idpProviderToModel(provider *admin.IdpProviderID) *iam_model.IDPProvider {
	return &iam_model.IDPProvider{
		IDPConfigID: provider.IdpConfigId,
		Type:        iam_model.IDPProviderTypeSystem,
	}
}

func idpProviderFromDomain(provider *domain.IDPProvider) *admin.IdpProviderID {
	return &admin.IdpProviderID{
		IdpConfigId: provider.IDPConfigID,
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

func secondFactorsResultFromModel(result *iam_model.SecondFactorsSearchResponse) *admin.SecondFactorsResult {
	converted := make([]admin.SecondFactorType, len(result.Result))
	for i, mfaType := range result.Result {
		converted[i] = secondFactorTypeFromModel(mfaType)
	}
	return &admin.SecondFactorsResult{
		SecondFactors: converted,
	}
}

func secondFactorFromDomain(mfaType domain.SecondFactorType) *admin.SecondFactor {
	return &admin.SecondFactor{
		SecondFactor: secondFactorTypeFromDomain(mfaType),
	}
}

func secondFactorTypeFromDomain(mfaType domain.SecondFactorType) admin.SecondFactorType {
	switch mfaType {
	case domain.SecondFactorTypeOTP:
		return admin.SecondFactorType_SECONDFACTORTYPE_OTP
	case domain.SecondFactorTypeU2F:
		return admin.SecondFactorType_SECONDFACTORTYPE_U2F
	default:
		return admin.SecondFactorType_SECONDFACTORTYPE_UNSPECIFIED
	}
}

func secondFactorTypeFromModel(mfaType iam_model.SecondFactorType) admin.SecondFactorType {
	switch mfaType {
	case iam_model.SecondFactorTypeOTP:
		return admin.SecondFactorType_SECONDFACTORTYPE_OTP
	case iam_model.SecondFactorTypeU2F:
		return admin.SecondFactorType_SECONDFACTORTYPE_U2F
	default:
		return admin.SecondFactorType_SECONDFACTORTYPE_UNSPECIFIED
	}
}

func secondFactorTypeToDomain(mfaType *admin.SecondFactor) domain.SecondFactorType {
	switch mfaType.SecondFactor {
	case admin.SecondFactorType_SECONDFACTORTYPE_OTP:
		return domain.SecondFactorTypeOTP
	case admin.SecondFactorType_SECONDFACTORTYPE_U2F:
		return domain.SecondFactorTypeU2F
	default:
		return domain.SecondFactorTypeUnspecified
	}
}

func passwordlessTypeFromDomain(passwordlessType domain.PasswordlessType) admin.PasswordlessType {
	switch passwordlessType {
	case domain.PasswordlessTypeAllowed:
		return admin.PasswordlessType_PASSWORDLESSTYPE_ALLOWED
	default:
		return admin.PasswordlessType_PASSWORDLESSTYPE_NOT_ALLOWED
	}
}

func passwordlessTypeToDomain(passwordlessType admin.PasswordlessType) domain.PasswordlessType {
	switch passwordlessType {
	case admin.PasswordlessType_PASSWORDLESSTYPE_ALLOWED:
		return domain.PasswordlessTypeAllowed
	default:
		return domain.PasswordlessTypeNotAllowed
	}
}

func multiFactorResultFromModel(result *iam_model.MultiFactorsSearchResponse) *admin.MultiFactorsResult {
	converted := make([]admin.MultiFactorType, len(result.Result))
	for i, mfaType := range result.Result {
		converted[i] = multiFactorTypeFromModel(mfaType)
	}
	return &admin.MultiFactorsResult{
		MultiFactors: converted,
	}
}

func multiFactorFromDomain(mfaType domain.MultiFactorType) *admin.MultiFactor {
	return &admin.MultiFactor{
		MultiFactor: multiFactorTypeFromDomain(mfaType),
	}
}

func multiFactorTypeFromDomain(mfaType domain.MultiFactorType) admin.MultiFactorType {
	switch mfaType {
	case domain.MultiFactorTypeU2FWithPIN:
		return admin.MultiFactorType_MULTIFACTORTYPE_U2F_WITH_PIN
	default:
		return admin.MultiFactorType_MULTIFACTORTYPE_UNSPECIFIED
	}
}

func multiFactorTypeFromModel(mfaType iam_model.MultiFactorType) admin.MultiFactorType {
	switch mfaType {
	case iam_model.MultiFactorTypeU2FWithPIN:
		return admin.MultiFactorType_MULTIFACTORTYPE_U2F_WITH_PIN
	default:
		return admin.MultiFactorType_MULTIFACTORTYPE_UNSPECIFIED
	}
}

func multiFactorTypeToDomain(mfaType *admin.MultiFactor) domain.MultiFactorType {
	switch mfaType.MultiFactor {
	case admin.MultiFactorType_MULTIFACTORTYPE_U2F_WITH_PIN:
		return domain.MultiFactorTypeU2FWithPIN
	default:
		return domain.MultiFactorTypeUnspecified
	}
}
