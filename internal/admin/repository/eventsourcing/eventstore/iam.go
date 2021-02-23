package eventstore

import (
	"context"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_view "github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/user/repository/view/model"
	"strings"

	caos_errs "github.com/caos/zitadel/internal/errors"

	"github.com/caos/logging"
	admin_view "github.com/caos/zitadel/internal/admin/repository/eventsourcing/view"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	usr_model "github.com/caos/zitadel/internal/user/model"
)

type IAMRepository struct {
	Eventstore     v1.Eventstore
	SearchLimit    uint64
	View           *admin_view.View
	SystemDefaults systemdefaults.SystemDefaults
	Roles          []string
}

func (repo *IAMRepository) IAMMemberByID(ctx context.Context, iamID, userID string) (*iam_model.IAMMemberView, error) {
	member, err := repo.View.IAMMemberByIDs(iamID, userID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.IAMMemberToModel(member), nil
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

func (repo *IAMRepository) IDPProvidersByIDPConfigID(ctx context.Context, idpConfigID string) ([]*iam_model.IDPProviderView, error) {
	providers, err := repo.View.IDPProvidersByIdpConfigID(idpConfigID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.IDPProviderViewsToModel(providers), nil
}

func (repo *IAMRepository) ExternalIDPsByIDPConfigID(ctx context.Context, idpConfigID string) ([]*usr_model.ExternalIDPView, error) {
	externalIDPs, err := repo.View.ExternalIDPsByIDPConfigID(idpConfigID)
	if err != nil {
		return nil, err
	}
	return model.ExternalIDPViewsToModel(externalIDPs), nil
}

func (repo *IAMRepository) ExternalIDPsByIDPConfigIDFromDefaultPolicy(ctx context.Context, idpConfigID string) ([]*usr_model.ExternalIDPView, error) {
	policies, err := repo.View.AllDefaultLoginPolicies()
	if err != nil {
		return nil, err
	}
	resourceOwners := make([]string, len(policies))
	for i, policy := range policies {
		resourceOwners[i] = policy.AggregateID
	}

	externalIDPs, err := repo.View.ExternalIDPsByIDPConfigIDAndResourceOwners(idpConfigID, resourceOwners)
	if err != nil {
		return nil, err
	}
	return model.ExternalIDPViewsToModel(externalIDPs), nil
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

func (repo *IAMRepository) GetDefaultLoginPolicy(ctx context.Context) (*iam_model.LoginPolicyView, error) {
	policy, viewErr := repo.View.LoginPolicyByAggregateID(repo.SystemDefaults.IamID)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		policy = new(iam_es_model.LoginPolicyView)
	}
	events, esErr := repo.getIAMEvents(ctx, policy.Sequence)
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

func (repo *IAMRepository) GetDefaultPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error) {
	policy, viewErr := repo.View.PasswordComplexityPolicyByAggregateID(repo.SystemDefaults.IamID)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		policy = new(iam_es_model.PasswordComplexityPolicyView)
	}
	events, esErr := repo.getIAMEvents(ctx, policy.Sequence)
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

func (repo *IAMRepository) GetDefaultPasswordAgePolicy(ctx context.Context) (*iam_model.PasswordAgePolicyView, error) {
	policy, viewErr := repo.View.PasswordAgePolicyByAggregateID(repo.SystemDefaults.IamID)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		policy = new(iam_es_model.PasswordAgePolicyView)
	}
	events, esErr := repo.getIAMEvents(ctx, policy.Sequence)
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

func (repo *IAMRepository) GetDefaultPasswordLockoutPolicy(ctx context.Context) (*iam_model.PasswordLockoutPolicyView, error) {
	policy, viewErr := repo.View.PasswordLockoutPolicyByAggregateID(repo.SystemDefaults.IamID)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		policy = new(iam_es_model.PasswordLockoutPolicyView)
	}
	events, esErr := repo.getIAMEvents(ctx, policy.Sequence)
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

func (repo *IAMRepository) GetOrgIAMPolicy(ctx context.Context) (*iam_model.OrgIAMPolicyView, error) {
	policy, viewErr := repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
	if viewErr != nil && !caos_errs.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if caos_errs.IsNotFound(viewErr) {
		policy = new(iam_es_model.OrgIAMPolicyView)
	}
	events, esErr := repo.getIAMEvents(ctx, policy.Sequence)
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

func (repo *IAMRepository) GetDefaultLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error) {
	policy, err := repo.View.LabelPolicyByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.LabelPolicyViewToModel(policy), err
}

func (repo *IAMRepository) GetDefaultMailTemplate(ctx context.Context) (*iam_model.MailTemplateView, error) {
	template, err := repo.View.MailTemplateByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.MailTemplateViewToModel(template), err
}

func (repo *IAMRepository) SearchIAMMembersx(ctx context.Context, request *iam_model.IAMMemberSearchRequest) (*iam_model.IAMMemberSearchResponse, error) {
	request.EnsureLimit(repo.SearchLimit)
	sequence, err := repo.View.GetLatestIAMMemberSequence()
	logging.Log("EVENT-Slkci").OnError(err).Warn("could not read latest iam sequence")
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
		result.Timestamp = result.Timestamp
	}
	return result, nil
}

func (repo *IAMRepository) GetDefaultMailTexts(ctx context.Context) (*iam_model.MailTextsView, error) {
	text, err := repo.View.MailTexts(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.MailTextsViewToModel(text, true), err
}

func (repo *IAMRepository) GetDefaultMailText(ctx context.Context, textType string, language string) (*iam_model.MailTextView, error) {
	text, err := repo.View.MailTextByIDs(repo.SystemDefaults.IamID, textType, language)
	if err != nil {
		return nil, err
	}
	text.Default = true
	return iam_es_model.MailTextViewToModel(text), err
}

func (repo *IAMRepository) getIAMEvents(ctx context.Context, sequence uint64) ([]*models.Event, error) {
	query, err := iam_view.IAMByIDQuery(domain.IAMID, sequence)
	if err != nil {
		return nil, err
	}
	return repo.Eventstore.FilterEvents(ctx, query)
}
