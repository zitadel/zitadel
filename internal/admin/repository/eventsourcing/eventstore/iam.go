package eventstore

import (
	"context"
	"strings"

	caos_errs "github.com/caos/zitadel/internal/errors"

	"github.com/caos/logging"
	admin_view "github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	es_sdk "github.com/caos/zitadel/internal/eventstore/sdk"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es "github.com/caos/zitadel/internal/iam/repository/eventsourcing"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	org_es "github.com/caos/zitadel/internal/org/repository/eventsourcing"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	usr_model "github.com/caos/zitadel/internal/user/model"
	usr_es "github.com/caos/zitadel/internal/user/repository/eventsourcing"
	iam_business "github.com/caos/zitadel/internal/v2/business/iam"
)

type IAMRepository struct {
	SearchLimit uint64
	*iam_es.IAMEventstore
	OrgEvents      *org_es.OrgEventstore
	UserEvents     *usr_es.UserEventstore
	View           *admin_view.View
	SystemDefaults systemdefaults.SystemDefaults
	Roles          []string

	IAMV2 *iam_business.Repository
}

func (repo *IAMRepository) IAMMemberByID(ctx context.Context, iamID, userID string) (*iam_model.IAMMemberView, error) {
	member, err := repo.View.IAMMemberByIDs(iamID, userID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.IAMMemberToModel(member), nil
}

func (repo *IAMRepository) AddIAMMember(ctx context.Context, member *iam_model.IAMMember) (*iam_model.IAMMember, error) {
	member.AggregateID = repo.SystemDefaults.IamID
	if repo.IAMV2 != nil {
		return repo.IAMV2.AddMember(ctx, member)
	}
	return repo.IAMEventstore.AddIAMMember(ctx, member)
}

func (repo *IAMRepository) ChangeIAMMember(ctx context.Context, member *iam_model.IAMMember) (*iam_model.IAMMember, error) {
	member.AggregateID = repo.SystemDefaults.IamID
	if repo.IAMV2 != nil {
		return repo.IAMV2.ChangeMember(ctx, member)
	}
	return repo.IAMEventstore.ChangeIAMMember(ctx, member)
}

func (repo *IAMRepository) RemoveIAMMember(ctx context.Context, userID string) error {
	member := iam_model.NewIAMMember(repo.SystemDefaults.IamID, userID)
	if repo.IAMV2 != nil {
		return repo.IAMV2.RemoveMember(ctx, member)
	}
	return repo.IAMEventstore.RemoveIAMMember(ctx, member)
}

func (repo *IAMRepository) SearchIAMMembers(ctx context.Context, request *iam_model.IAMMemberSearchRequest) (*iam_model.IAMMemberSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, err := repo.View.GetLatestIAMMemberSequence()
	logging.Log("EVENT-Slkci").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Warn("could not read latest iam sequence")
	members, count, err := repo.View.SearchIAMMembers(request)
	if err != nil {
		return nil, err
	}
	result := &iam_model.IAMMemberSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      iam_es_model.IAMMembersToModel(members),
	}
	if err == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *IAMRepository) GetIAMMemberRoles() []string {
	roles := make([]string, 0)
	for _, roleMap := range repo.Roles {
		if strings.HasPrefix(roleMap, "IAM") {
			roles = append(roles, roleMap)
		}
	}
	return roles
}

func (repo *IAMRepository) IDPConfigByID(ctx context.Context, idpConfigID string) (*iam_model.IDPConfigView, error) {
	if repo.IAMV2 != nil {
		return repo.IAMV2.IDPConfigByID(ctx, repo.SystemDefaults.IamID, idpConfigID)
	}

	idp, err := repo.View.IDPConfigByID(idpConfigID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.IDPConfigViewToModel(idp), nil
}

func (repo *IAMRepository) AddOIDCIDPConfig(ctx context.Context, idp *iam_model.IDPConfig) (*iam_model.IDPConfig, error) {
	idp.AggregateID = repo.SystemDefaults.IamID
	if repo.IAMV2 != nil {
		return repo.IAMV2.AddIDPConfig(ctx, idp)
	}
	return repo.IAMEventstore.AddIDPConfig(ctx, idp)
}

func (repo *IAMRepository) ChangeIDPConfig(ctx context.Context, idp *iam_model.IDPConfig) (*iam_model.IDPConfig, error) {
	idp.AggregateID = repo.SystemDefaults.IamID
	if repo.IAMV2 != nil {
		return repo.IAMV2.ChangeIDPConfig(ctx, idp)
	}
	return repo.IAMEventstore.ChangeIDPConfig(ctx, idp)
}

func (repo *IAMRepository) DeactivateIDPConfig(ctx context.Context, idpConfigID string) (*iam_model.IDPConfig, error) {
	if repo.IAMV2 != nil {
		return repo.IAMV2.DeactivateIDPConfig(ctx, repo.SystemDefaults.IamID, idpConfigID)
	}
	return repo.IAMEventstore.DeactivateIDPConfig(ctx, repo.SystemDefaults.IamID, idpConfigID)
}

func (repo *IAMRepository) ReactivateIDPConfig(ctx context.Context, idpConfigID string) (*iam_model.IDPConfig, error) {
	if repo.IAMV2 != nil {
		return repo.IAMV2.ReactivateIDPConfig(ctx, repo.SystemDefaults.IamID, idpConfigID)
	}
	return repo.IAMEventstore.ReactivateIDPConfig(ctx, repo.SystemDefaults.IamID, idpConfigID)
}

func (repo *IAMRepository) RemoveIDPConfig(ctx context.Context, idpConfigID string) error {
	// if repo.IAMV2 != nil {
	// 	return repo.IAMV2.
	// }
	aggregates := make([]*es_models.Aggregate, 0)
	idp := iam_model.NewIDPConfig(repo.SystemDefaults.IamID, idpConfigID)
	_, agg, err := repo.IAMEventstore.PrepareRemoveIDPConfig(ctx, idp)
	if err != nil {
		return err
	}
	aggregates = append(aggregates, agg)

	providers, err := repo.View.IDPProvidersByIdpConfigID(idpConfigID)
	if err != nil {
		return err
	}
	for _, p := range providers {
		if p.AggregateID == repo.SystemDefaults.IamID {
			continue
		}
		provider := &iam_model.IDPProvider{ObjectRoot: es_models.ObjectRoot{AggregateID: p.AggregateID}, IdpConfigID: p.IDPConfigID}
		providerAgg := new(es_models.Aggregate)
		_, providerAgg, err = repo.OrgEvents.PrepareRemoveIDPProviderFromLoginPolicy(ctx, provider, true)
		if err != nil {
			return err
		}
		aggregates = append(aggregates, providerAgg)
	}
	externalIDPs, err := repo.View.ExternalIDPsByIDPConfigID(idpConfigID)
	if err != nil {
		return err
	}
	for _, externalIDP := range externalIDPs {
		idpRemove := &usr_model.ExternalIDP{ObjectRoot: es_models.ObjectRoot{AggregateID: externalIDP.UserID}, IDPConfigID: externalIDP.IDPConfigID, UserID: externalIDP.ExternalUserID}
		idpAgg := make([]*es_models.Aggregate, 0)
		_, idpAgg, err = repo.UserEvents.PrepareRemoveExternalIDP(ctx, idpRemove, true)
		if err != nil {
			return err
		}
		aggregates = append(aggregates, idpAgg...)
	}
	return es_sdk.PushAggregates(ctx, repo.Eventstore.PushAggregates, nil, aggregates...)
}

func (repo *IAMRepository) ChangeOidcIDPConfig(ctx context.Context, oidcConfig *iam_model.OIDCIDPConfig) (*iam_model.OIDCIDPConfig, error) {
	oidcConfig.AggregateID = repo.SystemDefaults.IamID
	if repo.IAMV2 != nil {
		return repo.IAMV2.ChangeIDPOIDCConfig(ctx, oidcConfig)
	}
	return repo.IAMEventstore.ChangeIDPOIDCConfig(ctx, oidcConfig)
}

func (repo *IAMRepository) SearchIDPConfigs(ctx context.Context, request *iam_model.IDPConfigSearchRequest) (*iam_model.IDPConfigSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, err := repo.View.GetLatestIDPConfigSequence()
	logging.Log("EVENT-Dk8si").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Warn("could not read latest idp config sequence")
	idps, count, err := repo.View.SearchIDPConfigs(request)
	if err != nil {
		return nil, err
	}
	result := &iam_model.IDPConfigSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      iam_es_model.IdpConfigViewsToModel(idps),
	}
	if err == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *IAMRepository) GetDefaultLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error) {
	policy, viewErr := repo.View.LabelPolicyByAggregateID(repo.SystemDefaults.IamID)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		policy = new(iam_es_model.LabelPolicyView)
	}
	events, esErr := repo.IAMEventstore.IAMEventsByID(ctx, repo.SystemDefaults.IamID, policy.Sequence)
	if caos_errs.IsNotFound(viewErr) && len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-4bM0s", "Errors.IAM.LabelPolicy.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-3M0xs").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.LabelPolicyViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.LabelPolicyViewToModel(policy), nil
		}
	}
	return iam_es_model.LabelPolicyViewToModel(policy), nil
}

func (repo *IAMRepository) AddDefaultLabelPolicy(ctx context.Context, policy *iam_model.LabelPolicy) (*iam_model.LabelPolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.AddLabelPolicy(ctx, policy)
}

func (repo *IAMRepository) ChangeDefaultLabelPolicy(ctx context.Context, policy *iam_model.LabelPolicy) (*iam_model.LabelPolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.ChangeLabelPolicy(ctx, policy)
}

func (repo *IAMRepository) GetDefaultLoginPolicy(ctx context.Context) (*iam_model.LoginPolicyView, error) {
	policy, viewErr := repo.View.LoginPolicyByAggregateID(repo.SystemDefaults.IamID)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		policy = new(iam_es_model.LoginPolicyView)
	}
	events, esErr := repo.IAMEventstore.IAMEventsByID(ctx, repo.SystemDefaults.IamID, policy.Sequence)
	if caos_errs.IsNotFound(viewErr) && len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-cmO9s", "Errors.IAM.LoginPolicy.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-2Mi8s").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.LoginPolicyViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.LoginPolicyViewToModel(policy), nil
		}
	}
	return iam_es_model.LoginPolicyViewToModel(policy), nil
}

func (repo *IAMRepository) AddDefaultLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	if repo.IAMV2 != nil {
		return repo.IAMV2.AddLoginPolicy(ctx, policy)
	}
	return repo.IAMEventstore.AddLoginPolicy(ctx, policy)
}

func (repo *IAMRepository) ChangeDefaultLoginPolicy(ctx context.Context, policy *iam_model.LoginPolicy) (*iam_model.LoginPolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	if repo.IAMV2 != nil {
		return repo.IAMV2.ChangeLoginPolicy(ctx, policy)
	}
	return repo.IAMEventstore.ChangeLoginPolicy(ctx, policy)
}

func (repo *IAMRepository) SearchDefaultIDPProviders(ctx context.Context, request *iam_model.IDPProviderSearchRequest) (*iam_model.IDPProviderSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	request.AppendAggregateIDQuery(repo.SystemDefaults.IamID)
	sequence, err := repo.View.GetLatestIDPProviderSequence()
	logging.Log("EVENT-Tuiks").OnError(err).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Warn("could not read latest iam sequence")
	providers, count, err := repo.View.SearchIDPProviders(request)
	if err != nil {
		return nil, err
	}
	result := &iam_model.IDPProviderSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      iam_es_model.IDPProviderViewsToModel(providers),
	}
	if err == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *IAMRepository) AddIDPProviderToLoginPolicy(ctx context.Context, provider *iam_model.IDPProvider) (*iam_model.IDPProvider, error) {
	provider.AggregateID = repo.SystemDefaults.IamID
	if repo.IAMV2 != nil {
		return repo.IAMV2.AddIDPProviderToLoginPolicy(ctx, provider)
	}
	return repo.IAMEventstore.AddIDPProviderToLoginPolicy(ctx, provider)
}

func (repo *IAMRepository) RemoveIDPProviderFromLoginPolicy(ctx context.Context, provider *iam_model.IDPProvider) error {
	aggregates := make([]*es_models.Aggregate, 0)
	provider.AggregateID = repo.SystemDefaults.IamID
	_, removeAgg, err := repo.IAMEventstore.PrepareRemoveIDPProviderFromLoginPolicy(ctx, provider)
	if err != nil {
		return err
	}
	aggregates = append(aggregates, removeAgg)

	externalIDPs, err := repo.View.ExternalIDPsByIDPConfigID(provider.IdpConfigID)
	if err != nil {
		return err
	}
	for _, externalIDP := range externalIDPs {
		idpRemove := &usr_model.ExternalIDP{ObjectRoot: es_models.ObjectRoot{AggregateID: externalIDP.UserID}, IDPConfigID: externalIDP.IDPConfigID, UserID: externalIDP.ExternalUserID}
		idpAgg := make([]*es_models.Aggregate, 0)
		_, idpAgg, err = repo.UserEvents.PrepareRemoveExternalIDP(ctx, idpRemove, true)
		if err != nil {
			return err
		}
		aggregates = append(aggregates, idpAgg...)
	}
	return es_sdk.PushAggregates(ctx, repo.Eventstore.PushAggregates, nil, aggregates...)
}

func (repo *IAMRepository) SearchDefaultSecondFactors(ctx context.Context) (*iam_model.SecondFactorsSearchResponse, error) {
	policy, err := repo.GetDefaultLoginPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &iam_model.SecondFactorsSearchResponse{
		TotalResult: uint64(len(policy.SecondFactors)),
		Result:      policy.SecondFactors,
	}, nil
}

func (repo *IAMRepository) AddSecondFactorToLoginPolicy(ctx context.Context, mfa iam_model.SecondFactorType) (iam_model.SecondFactorType, error) {
	return repo.IAMEventstore.AddSecondFactorToLoginPolicy(ctx, repo.SystemDefaults.IamID, mfa)
}

func (repo *IAMRepository) RemoveSecondFactorFromLoginPolicy(ctx context.Context, mfa iam_model.SecondFactorType) error {
	return repo.IAMEventstore.RemoveSecondFactorFromLoginPolicy(ctx, repo.SystemDefaults.IamID, mfa)
}

func (repo *IAMRepository) SearchDefaultMultiFactors(ctx context.Context) (*iam_model.MultiFactorsSearchResponse, error) {
	policy, err := repo.GetDefaultLoginPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &iam_model.MultiFactorsSearchResponse{
		TotalResult: uint64(len(policy.MultiFactors)),
		Result:      policy.MultiFactors,
	}, nil
}

func (repo *IAMRepository) AddMultiFactorToLoginPolicy(ctx context.Context, mfa iam_model.MultiFactorType) (iam_model.MultiFactorType, error) {
	return repo.IAMEventstore.AddMultiFactorToLoginPolicy(ctx, repo.SystemDefaults.IamID, mfa)
}

func (repo *IAMRepository) RemoveMultiFactorFromLoginPolicy(ctx context.Context, mfa iam_model.MultiFactorType) error {
	return repo.IAMEventstore.RemoveMultiFactorFromLoginPolicy(ctx, repo.SystemDefaults.IamID, mfa)
}

func (repo *IAMRepository) GetDefaultPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error) {
	policy, viewErr := repo.View.PasswordComplexityPolicyByAggregateID(repo.SystemDefaults.IamID)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		policy = new(iam_es_model.PasswordComplexityPolicyView)
	}
	events, esErr := repo.IAMEventstore.IAMEventsByID(ctx, repo.SystemDefaults.IamID, policy.Sequence)
	if caos_errs.IsNotFound(viewErr) && len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-1Mc0s", "Errors.IAM.PasswordComplexityPolicy.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-3M0xs").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.PasswordComplexityViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.PasswordComplexityViewToModel(policy), nil
		}
	}
	return iam_es_model.PasswordComplexityViewToModel(policy), nil
}

func (repo *IAMRepository) AddDefaultPasswordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) (*iam_model.PasswordComplexityPolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.AddPasswordComplexityPolicy(ctx, policy)
}

func (repo *IAMRepository) ChangeDefaultPasswordComplexityPolicy(ctx context.Context, policy *iam_model.PasswordComplexityPolicy) (*iam_model.PasswordComplexityPolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.ChangePasswordComplexityPolicy(ctx, policy)
}

func (repo *IAMRepository) GetDefaultPasswordAgePolicy(ctx context.Context) (*iam_model.PasswordAgePolicyView, error) {
	policy, viewErr := repo.View.PasswordAgePolicyByAggregateID(repo.SystemDefaults.IamID)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		policy = new(iam_es_model.PasswordAgePolicyView)
	}
	events, esErr := repo.IAMEventstore.IAMEventsByID(ctx, repo.SystemDefaults.IamID, policy.Sequence)
	if caos_errs.IsNotFound(viewErr) && len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-vMyS3", "Errors.IAM.PasswordAgePolicy.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-3M0xs").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.PasswordAgeViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.PasswordAgeViewToModel(policy), nil
		}
	}
	return iam_es_model.PasswordAgeViewToModel(policy), nil
}

func (repo *IAMRepository) AddDefaultPasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*iam_model.PasswordAgePolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.AddPasswordAgePolicy(ctx, policy)
}

func (repo *IAMRepository) ChangeDefaultPasswordAgePolicy(ctx context.Context, policy *iam_model.PasswordAgePolicy) (*iam_model.PasswordAgePolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.ChangePasswordAgePolicy(ctx, policy)
}

func (repo *IAMRepository) GetDefaultPasswordLockoutPolicy(ctx context.Context) (*iam_model.PasswordLockoutPolicyView, error) {
	policy, viewErr := repo.View.PasswordLockoutPolicyByAggregateID(repo.SystemDefaults.IamID)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		policy = new(iam_es_model.PasswordLockoutPolicyView)
	}
	events, esErr := repo.IAMEventstore.IAMEventsByID(ctx, repo.SystemDefaults.IamID, policy.Sequence)
	if caos_errs.IsNotFound(viewErr) && len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-2M9oP", "Errors.IAM.PasswordLockoutPolicy.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-3M0xs").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.PasswordLockoutViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.PasswordLockoutViewToModel(policy), nil
		}
	}
	return iam_es_model.PasswordLockoutViewToModel(policy), nil
}

func (repo *IAMRepository) AddDefaultPasswordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) (*iam_model.PasswordLockoutPolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.AddPasswordLockoutPolicy(ctx, policy)
}

func (repo *IAMRepository) ChangeDefaultPasswordLockoutPolicy(ctx context.Context, policy *iam_model.PasswordLockoutPolicy) (*iam_model.PasswordLockoutPolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.ChangePasswordLockoutPolicy(ctx, policy)
}

func (repo *IAMRepository) GetOrgIAMPolicy(ctx context.Context) (*iam_model.OrgIAMPolicyView, error) {
	policy, viewErr := repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		policy = new(iam_es_model.OrgIAMPolicyView)
	}
	events, esErr := repo.IAMEventstore.IAMEventsByID(ctx, repo.SystemDefaults.IamID, policy.Sequence)
	if caos_errs.IsNotFound(viewErr) && len(events) == 0 {
		return nil, caos_errs.ThrowNotFound(nil, "EVENT-MkoL0", "Errors.IAM.OrgIAMPolicy.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-3M0xs").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.OrgIAMViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.OrgIAMViewToModel(policy), nil
		}
	}
	return iam_es_model.OrgIAMViewToModel(policy), nil
}

func (repo *IAMRepository) AddDefaultOrgIAMPolicy(ctx context.Context, policy *iam_model.OrgIAMPolicy) (*iam_model.OrgIAMPolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.AddOrgIAMPolicy(ctx, policy)
}

func (repo *IAMRepository) ChangeDefaultOrgIAMPolicy(ctx context.Context, policy *iam_model.OrgIAMPolicy) (*iam_model.OrgIAMPolicy, error) {
	policy.AggregateID = repo.SystemDefaults.IamID
	return repo.IAMEventstore.ChangeOrgIAMPolicy(ctx, policy)
}
