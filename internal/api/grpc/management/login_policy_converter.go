package management

import (
	"context"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/eventstore/models"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	"github.com/caos/zitadel/internal/v2/domain"
	"github.com/caos/zitadel/pkg/grpc/management"
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

func idpProviderIDToDomain(ctx context.Context, provider *management.IdpProviderID) *domain.IDPProvider {
	return &domain.IDPProvider{
		ObjectRoot: models.ObjectRoot{
			AggregateID: authz.GetCtxData(ctx).OrgID,
		},
		IDPConfigID: provider.IdpConfigId,
	}
}

func idpProviderAddToDomain(ctx context.Context, provider *management.IdpProviderAdd) *domain.IDPProvider {
	return &domain.IDPProvider{
		ObjectRoot: models.ObjectRoot{
			AggregateID: authz.GetCtxData(ctx).OrgID,
		},
		IDPConfigID: provider.IdpConfigId,
		Type:        idpProviderTypeToDomain(provider.IdpProviderType),
	}
}

func idpProviderIDFromModel(provider *iam_model.IDPProvider) *management.IdpProviderID {
	return &management.IdpProviderID{
		IdpConfigId: provider.IDPConfigID,
	}
}

func idpProviderFromDomain(provider *domain.IDPProvider) *management.IdpProvider {
	return &management.IdpProvider{
		IdpConfigId:      provider.IDPConfigID,
		IdpProvider_Type: idpProviderTypeFromDomain(provider.Type),
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

func idpProviderTypeToDomain(providerType management.IdpProviderType) domain.IdentityProviderType {
	switch providerType {
	case management.IdpProviderType_IDPPROVIDERTYPE_SYSTEM:
		return domain.IdentityProviderTypeSystem
	case management.IdpProviderType_IDPPROVIDERTYPE_ORG:
		return domain.IdentityProviderTypeOrg
	default:
		return domain.IdentityProviderTypeSystem
	}
}

func idpProviderTypeFromDomain(providerType domain.IdentityProviderType) management.IdpProviderType {
	switch providerType {
	case domain.IdentityProviderTypeSystem:
		return management.IdpProviderType_IDPPROVIDERTYPE_SYSTEM
	case domain.IdentityProviderTypeOrg:
		return management.IdpProviderType_IDPPROVIDERTYPE_ORG
	default:
		return management.IdpProviderType_IDPPROVIDERTYPE_UNSPECIFIED
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

func secondFactorFromDomain(mfaType domain.SecondFactorType) *management.SecondFactor {
	return &management.SecondFactor{
		SecondFactor: secondFactorTypeFromDomain(mfaType),
	}
}

func secondFactorFromModel(mfaType iam_model.SecondFactorType) *management.SecondFactor {
	return &management.SecondFactor{
		SecondFactor: secondFactorTypeFromModel(mfaType),
	}
}

func secondFactorTypeFromDomain(mfaType domain.SecondFactorType) management.SecondFactorType {
	switch mfaType {
	case domain.SecondFactorTypeOTP:
		return management.SecondFactorType_SECONDFACTORTYPE_OTP
	case domain.SecondFactorTypeU2F:
		return management.SecondFactorType_SECONDFACTORTYPE_U2F
	default:
		return management.SecondFactorType_SECONDFACTORTYPE_UNSPECIFIED
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

func secondFactorTypeToDomain(mfaType *management.SecondFactor) domain.SecondFactorType {
	switch mfaType.SecondFactor {
	case management.SecondFactorType_SECONDFACTORTYPE_OTP:
		return domain.SecondFactorTypeOTP
	case management.SecondFactorType_SECONDFACTORTYPE_U2F:
		return domain.SecondFactorTypeU2F
	default:
		return domain.SecondFactorTypeUnspecified
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

func multiFactorFromDomain(mfaType domain.MultiFactorType) *management.MultiFactor {
	return &management.MultiFactor{
		MultiFactor: multiFactorTypeFromDomain(mfaType),
	}
}

func multiFactorTypeFromDomain(mfaType domain.MultiFactorType) management.MultiFactorType {
	switch mfaType {
	case domain.MultiFactorTypeU2FWithPIN:
		return management.MultiFactorType_MULTIFACTORTYPE_U2F_WITH_PIN
	default:
		return management.MultiFactorType_MULTIFACTORTYPE_UNSPECIFIED
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

func multiFactorTypeToDomain(mfaType *management.MultiFactor) domain.MultiFactorType {
	switch mfaType.MultiFactor {
	case management.MultiFactorType_MULTIFACTORTYPE_U2F_WITH_PIN:
		return domain.MultiFactorTypeU2FWithPIN
	default:
		return domain.MultiFactorTypeUnspecified
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
