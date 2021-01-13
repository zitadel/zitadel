package management

import (
	"context"

	"github.com/caos/logging"
	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/management"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func loginPolicyRequestToDomain(ctx context.Context, policy *management.LoginPolicyRequest) *domain.LoginPolicy {
	return &domain.LoginPolicy{
		ObjectRoot: models.ObjectRoot{
			AggregateID: authz.GetCtxData(ctx).OrgID,
		},
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowExternalIdp:      policy.AllowExternalIdp,
		AllowRegister:         policy.AllowRegister,
		ForceMFA:              policy.ForceMfa,
		PasswordlessType:      passwordlessTypeToDomain(policy.PasswordlessType),
	}
}

func loginPolicyFromDomain(policy *domain.LoginPolicy) *management.LoginPolicy {
	return &management.LoginPolicy{
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowExternalIdp:      policy.AllowExternalIdp,
		AllowRegister:         policy.AllowRegister,
		CreationDate:          timestamppb.New(policy.CreationDate),
		ChangeDate:            timestamppb.New(policy.ChangeDate),
		ForceMfa:              policy.ForceMFA,
		PasswordlessType:      passwordlessTypeFromDomain(policy.PasswordlessType),
	}
}

func loginPolicyViewFromModel(policy *iam_model.LoginPolicyView) *management.LoginPolicyView {
	creationDate, err := ptypes.TimestampProto(policy.CreationDate)
	logging.Log("GRPC-5Tsm8").OnError(err).Debug("date parse failed")

	changeDate, err := ptypes.TimestampProto(policy.ChangeDate)
	logging.Log("GRPC-8dJgs").OnError(err).Debug("date parse failed")

	return &management.LoginPolicyView{
		Default:               policy.Default,
		AllowUsernamePassword: policy.AllowUsernamePassword,
		AllowExternalIdp:      policy.AllowExternalIDP,
		AllowRegister:         policy.AllowRegister,
		CreationDate:          creationDate,
		ChangeDate:            changeDate,
		ForceMfa:              policy.ForceMFA,
		PasswordlessType:      passwordlessTypeFromModel(policy.PasswordlessType),
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
		IDPConfigID: provider.IdpConfigId,
		Type:        iam_model.IDPProviderTypeSystem,
	}
}

func idpProviderAddToModel(provider *management.IdpProviderAdd) *iam_model.IDPProvider {
	return &iam_model.IDPProvider{
		IDPConfigID: provider.IdpConfigId,
		Type:        idpProviderTypeToModel(provider.IdpProviderType),
	}
}

func idpProviderIDFromModel(provider *iam_model.IDPProvider) *management.IdpProviderID {
	return &management.IdpProviderID{
		IdpConfigId: provider.IDPConfigID,
	}
}

func idpProviderFromModel(provider *iam_model.IDPProvider) *management.IdpProvider {
	return &management.IdpProvider{
		IdpConfigId:      provider.IDPConfigID,
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

func secondFactorResultFromModel(result *iam_model.SecondFactorsSearchResponse) *management.SecondFactorsResult {
	converted := make([]management.SecondFactorType, len(result.Result))
	for i, mfaType := range result.Result {
		converted[i] = secondFactorTypeFromModel(mfaType)
	}
	return &management.SecondFactorsResult{
		SecondFactors: converted,
	}
}

func secondFactorFromModel(mfaType iam_model.SecondFactorType) *management.SecondFactor {
	return &management.SecondFactor{
		SecondFactor: secondFactorTypeFromModel(mfaType),
	}
}

func secondFactorTypeFromModel(mfaType iam_model.SecondFactorType) management.SecondFactorType {
	switch mfaType {
	case iam_model.SecondFactorTypeOTP:
		return management.SecondFactorType_SECONDFACTORTYPE_OTP
	case iam_model.SecondFactorTypeU2F:
		return management.SecondFactorType_SECONDFACTORTYPE_U2F
	default:
		return management.SecondFactorType_SECONDFACTORTYPE_UNSPECIFIED
	}
}

func secondFactorTypeToModel(mfaType *management.SecondFactor) iam_model.SecondFactorType {
	switch mfaType.SecondFactor {
	case management.SecondFactorType_SECONDFACTORTYPE_OTP:
		return iam_model.SecondFactorTypeOTP
	case management.SecondFactorType_SECONDFACTORTYPE_U2F:
		return iam_model.SecondFactorTypeU2F
	default:
		return iam_model.SecondFactorTypeUnspecified
	}
}

func multiFactorResultFromModel(result *iam_model.MultiFactorsSearchResponse) *management.MultiFactorsResult {
	converted := make([]management.MultiFactorType, len(result.Result))
	for i, mfaType := range result.Result {
		converted[i] = multiFactorTypeFromModel(mfaType)
	}
	return &management.MultiFactorsResult{
		MultiFactors: converted,
	}
}

func multiFactorFromModel(mfaType iam_model.MultiFactorType) *management.MultiFactor {
	return &management.MultiFactor{
		MultiFactor: multiFactorTypeFromModel(mfaType),
	}
}

func multiFactorTypeFromModel(mfaType iam_model.MultiFactorType) management.MultiFactorType {
	switch mfaType {
	case iam_model.MultiFactorTypeU2FWithPIN:
		return management.MultiFactorType_MULTIFACTORTYPE_U2F_WITH_PIN
	default:
		return management.MultiFactorType_MULTIFACTORTYPE_UNSPECIFIED
	}
}

func multiFactorTypeToModel(mfaType *management.MultiFactor) iam_model.MultiFactorType {
	switch mfaType.MultiFactor {
	case management.MultiFactorType_MULTIFACTORTYPE_U2F_WITH_PIN:
		return iam_model.MultiFactorTypeU2FWithPIN
	default:
		return iam_model.MultiFactorTypeUnspecified
	}
}

func passwordlessTypeFromModel(passwordlessType iam_model.PasswordlessType) management.PasswordlessType {
	switch passwordlessType {
	case iam_model.PasswordlessTypeAllowed:
		return management.PasswordlessType_PASSWORDLESSTYPE_ALLOWED
	default:
		return management.PasswordlessType_PASSWORDLESSTYPE_NOT_ALLOWED
	}
}

func passwordlessTypeFromDomain(passwordlessType domain.PasswordlessType) management.PasswordlessType {
	switch passwordlessType {
	case domain.PasswordlessTypeAllowed:
		return management.PasswordlessType_PASSWORDLESSTYPE_ALLOWED
	default:
		return management.PasswordlessType_PASSWORDLESSTYPE_NOT_ALLOWED
	}
}

func passwordlessTypeToDomain(passwordlessType management.PasswordlessType) domain.PasswordlessType {
	switch passwordlessType {
	case management.PasswordlessType_PASSWORDLESSTYPE_ALLOWED:
		return domain.PasswordlessTypeAllowed
	default:
		return domain.PasswordlessTypeNotAllowed
	}
}
