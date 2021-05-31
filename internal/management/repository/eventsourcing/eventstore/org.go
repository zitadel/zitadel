package eventstore

import (
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"

	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	iam_view "github.com/caos/zitadel/internal/iam/repository/view"
	org_view "github.com/caos/zitadel/internal/org/repository/view"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/errors"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	mgmt_view "github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	"github.com/caos/zitadel/internal/org/repository/view/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

const (
	orgOwnerRole = "ORG_OWNER"
)

type OrgRepository struct {
	SearchLimit    uint64
	Eventstore     v1.Eventstore
	View           *mgmt_view.View
	Roles          []string
	SystemDefaults systemdefaults.SystemDefaults
}

func (repo *OrgRepository) OrgByID(ctx context.Context, id string) (*org_model.OrgView, error) {
	org, err := repo.View.OrgByID(id)
	if err != nil {
		return nil, err
	}
	return model.OrgToModel(org), nil
}

func (repo *OrgRepository) OrgByDomainGlobal(ctx context.Context, domain string) (*org_model.OrgView, error) {
	verifiedDomain, err := repo.View.VerifiedOrgDomain(domain)
	if err != nil {
		return nil, err
	}
	return repo.OrgByID(ctx, verifiedDomain.OrgID)
}

func (repo *OrgRepository) GetMyOrgIamPolicy(ctx context.Context) (*iam_model.OrgIAMPolicyView, error) {
	policy, err := repo.View.OrgIAMPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if errors.IsNotFound(err) {
		policy, err = repo.View.OrgIAMPolicyByAggregateID(repo.SystemDefaults.IamID)
		if err != nil {
			return nil, err
		}
		policy.Default = true
	}
	if err != nil {
		return nil, err
	}
	return iam_es_model.OrgIAMViewToModel(policy), err
}

func (repo *OrgRepository) SearchMyOrgDomains(ctx context.Context, request *org_model.OrgDomainSearchRequest) (*org_model.OrgDomainSearchResponse, error) {
	err := request.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
	request.Queries = append(request.Queries, &org_model.OrgDomainSearchQuery{Key: org_model.OrgDomainSearchKeyOrgID, Method: domain.SearchMethodEquals, Value: authz.GetCtxData(ctx).OrgID})
	sequence, sequenceErr := repo.View.GetLatestOrgDomainSequence()
	logging.Log("EVENT-SLowp").OnError(sequenceErr).WithField("traceID", tracing.TraceIDFromCtx(ctx)).Warn("could not read latest org domain sequence")
	domains, count, err := repo.View.SearchOrgDomains(request)
	if err != nil {
		return nil, err
	}
	result := &org_model.OrgDomainSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: uint64(count),
		Result:      model.OrgDomainsToModel(domains),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *OrgRepository) OrgChanges(ctx context.Context, id string, lastSequence uint64, limit uint64, sortAscending bool, auditLogRetention time.Duration) (*org_model.OrgChanges, error) {
	changes, err := repo.getOrgChanges(ctx, id, lastSequence, limit, sortAscending, auditLogRetention)
	if err != nil {
		return nil, err
	}
	for _, change := range changes.Changes {
		change.ModifierName = change.ModifierId
		user, _ := repo.userByID(ctx, change.ModifierId)
		if user != nil {
			if user.HumanView != nil {
				change.ModifierName = user.DisplayName
			}
			if user.MachineView != nil {
				change.ModifierName = user.MachineView.Name
			}
		}
	}
	return changes, nil
}

func (repo *OrgRepository) OrgMemberByID(ctx context.Context, orgID, userID string) (*org_model.OrgMemberView, error) {
	member, err := repo.View.OrgMemberByIDs(orgID, userID)
	if err != nil {
		return nil, err
	}
	return model.OrgMemberToModel(member), nil
}

func (repo *OrgRepository) SearchMyOrgMembers(ctx context.Context, request *org_model.OrgMemberSearchRequest) (*org_model.OrgMemberSearchResponse, error) {
	err := request.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
	request.Queries = append(request.Queries, &org_model.OrgMemberSearchQuery{Key: org_model.OrgMemberSearchKeyOrgID, Method: domain.SearchMethodEquals, Value: authz.GetCtxData(ctx).OrgID})
	sequence, sequenceErr := repo.View.GetLatestOrgMemberSequence()
	logging.Log("EVENT-Smu3d").OnError(sequenceErr).Warn("could not read latest org member sequence")
	members, count, err := repo.View.SearchOrgMembers(request)
	if err != nil {
		return nil, err
	}
	result := &org_model.OrgMemberSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      model.OrgMembersToModel(members),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *OrgRepository) GetOrgMemberRoles() []string {
	roles := make([]string, 0)
	for _, roleMap := range repo.Roles {
		if strings.HasPrefix(roleMap, "ORG") {
			roles = append(roles, roleMap)
		}
	}
	return roles
}

func (repo *OrgRepository) IDPConfigByID(ctx context.Context, idpConfigID string) (*iam_model.IDPConfigView, error) {
	idp, err := repo.View.IDPConfigByID(idpConfigID)
	if err != nil {
		return nil, err
	}
	return iam_view_model.IDPConfigViewToModel(idp), nil
}

func (repo *OrgRepository) SearchIDPConfigs(ctx context.Context, request *iam_model.IDPConfigSearchRequest) (*iam_model.IDPConfigSearchResponse, error) {
	err := request.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
	request.AppendMyOrgQuery(authz.GetCtxData(ctx).OrgID, repo.SystemDefaults.IamID)

	sequence, sequenceErr := repo.View.GetLatestIDPConfigSequence()
	logging.Log("EVENT-Dk8si").OnError(sequenceErr).Warn("could not read latest idp config sequence")
	idps, count, err := repo.View.SearchIDPConfigs(request)
	if err != nil {
		return nil, err
	}
	result := &iam_model.IDPConfigSearchResponse{
		Offset:      request.Offset,
		Limit:       request.Limit,
		TotalResult: count,
		Result:      iam_view_model.IdpConfigViewsToModel(idps),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *OrgRepository) GetLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error) {
	policy, err := repo.View.LabelPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if errors.IsNotFound(err) {
		policy, err = repo.View.LabelPolicyByAggregateID(repo.SystemDefaults.IamID)
		if err != nil {
			return nil, err
		}
		policy.Default = true
	}
	if err != nil {
		return nil, err
	}
	return iam_es_model.LabelPolicyViewToModel(policy), err
}

func (repo *OrgRepository) GetDefaultLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error) {
	policy, viewErr := repo.View.LabelPolicyByAggregateID(repo.SystemDefaults.IamID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		policy = new(iam_es_model.LabelPolicyView)
	}
	events, esErr := repo.getIAMEvents(ctx, policy.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-3Nf8sd", "Errors.IAM.LabelPolicy.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-28uLp").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.LabelPolicyViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.LabelPolicyViewToModel(policy), nil
		}
	}
	policy.Default = true
	return iam_es_model.LabelPolicyViewToModel(policy), nil
}

func (repo *OrgRepository) GetLoginPolicy(ctx context.Context) (*iam_model.LoginPolicyView, error) {
	policy, viewErr := repo.View.LoginPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		policy = new(iam_es_model.LoginPolicyView)
	}
	events, esErr := repo.getOrgEvents(ctx, repo.SystemDefaults.IamID, policy.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return repo.GetDefaultLoginPolicy(ctx)
	}
	if esErr != nil {
		logging.Log("EVENT-38iTr").WithError(esErr).Debug("error retrieving new events")
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

func (repo *OrgRepository) GetIDPProvidersByIDPConfigID(ctx context.Context, aggregateID, idpConfigID string) ([]*iam_model.IDPProviderView, error) {
	idpProviders, err := repo.View.IDPProvidersByIdpConfigID(aggregateID, idpConfigID)
	if err != nil {
		return nil, err
	}
	return iam_view_model.IDPProviderViewsToModel(idpProviders), err
}

func (repo *OrgRepository) GetDefaultLoginPolicy(ctx context.Context) (*iam_model.LoginPolicyView, error) {
	policy, viewErr := repo.View.LoginPolicyByAggregateID(domain.IAMID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		policy = new(iam_es_model.LoginPolicyView)
	}
	events, esErr := repo.getIAMEvents(ctx, policy.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-cmO9s", "Errors.IAM.LoginPolicy.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-28uLp").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.LoginPolicyViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.LoginPolicyViewToModel(policy), nil
		}
	}
	policy.Default = true
	return iam_es_model.LoginPolicyViewToModel(policy), nil
}

func (repo *OrgRepository) SearchIDPProviders(ctx context.Context, request *iam_model.IDPProviderSearchRequest) (*iam_model.IDPProviderSearchResponse, error) {
	policy, err := repo.View.LoginPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	if policy.Default {
		request.AppendAggregateIDQuery(domain.IAMID)
	} else {
		request.AppendAggregateIDQuery(authz.GetCtxData(ctx).OrgID)
	}
	err = request.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
	sequence, sequenceErr := repo.View.GetLatestIDPProviderSequence()
	logging.Log("EVENT-Tuiks").OnError(sequenceErr).Warn("could not read latest iam sequence")
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
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
}

func (repo *OrgRepository) SearchSecondFactors(ctx context.Context) (*iam_model.SecondFactorsSearchResponse, error) {
	policy, err := repo.GetLoginPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &iam_model.SecondFactorsSearchResponse{
		TotalResult: uint64(len(policy.SecondFactors)),
		Result:      policy.SecondFactors,
	}, nil
}

func (repo *OrgRepository) SearchMultiFactors(ctx context.Context) (*iam_model.MultiFactorsSearchResponse, error) {
	policy, err := repo.GetLoginPolicy(ctx)
	if err != nil {
		return nil, err
	}
	return &iam_model.MultiFactorsSearchResponse{
		TotalResult: uint64(len(policy.MultiFactors)),
		Result:      policy.MultiFactors,
	}, nil
}

func (repo *OrgRepository) GetPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error) {
	policy, viewErr := repo.View.PasswordComplexityPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		policy = new(iam_es_model.PasswordComplexityPolicyView)
	}
	events, esErr := repo.getOrgEvents(ctx, repo.SystemDefaults.IamID, policy.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return repo.GetDefaultPasswordComplexityPolicy(ctx)
	}
	if esErr != nil {
		logging.Log("EVENT-1Bx8s").WithError(esErr).Debug("error retrieving new events")
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

func (repo *OrgRepository) GetDefaultPasswordComplexityPolicy(ctx context.Context) (*iam_model.PasswordComplexityPolicyView, error) {
	policy, viewErr := repo.View.PasswordComplexityPolicyByAggregateID(repo.SystemDefaults.IamID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		policy = new(iam_es_model.PasswordComplexityPolicyView)
	}
	events, esErr := repo.getIAMEvents(ctx, policy.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-cmO9s", "Errors.IAM.PasswordComplexityPolicy.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-pL9sw").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.PasswordComplexityViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.PasswordComplexityViewToModel(policy), nil
		}
	}
	policy.Default = true
	return iam_es_model.PasswordComplexityViewToModel(policy), nil
}

func (repo *OrgRepository) GetPasswordAgePolicy(ctx context.Context) (*iam_model.PasswordAgePolicyView, error) {
	policy, viewErr := repo.View.PasswordAgePolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		policy = new(iam_es_model.PasswordAgePolicyView)
	}
	events, esErr := repo.getOrgEvents(ctx, repo.SystemDefaults.IamID, policy.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return repo.GetDefaultPasswordAgePolicy(ctx)
	}
	if esErr != nil {
		logging.Log("EVENT-5Mx7s").WithError(esErr).Debug("error retrieving new events")
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

func (repo *OrgRepository) GetDefaultPasswordAgePolicy(ctx context.Context) (*iam_model.PasswordAgePolicyView, error) {
	policy, viewErr := repo.View.PasswordAgePolicyByAggregateID(repo.SystemDefaults.IamID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		policy = new(iam_es_model.PasswordAgePolicyView)
	}
	events, esErr := repo.getIAMEvents(ctx, policy.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-cmO9s", "Errors.IAM.PasswordAgePolicy.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-3I90s").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.PasswordAgeViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.PasswordAgeViewToModel(policy), nil
		}
	}
	policy.Default = true
	return iam_es_model.PasswordAgeViewToModel(policy), nil
}

func (repo *OrgRepository) GetPasswordLockoutPolicy(ctx context.Context) (*iam_model.PasswordLockoutPolicyView, error) {
	policy, viewErr := repo.View.PasswordLockoutPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		policy = new(iam_es_model.PasswordLockoutPolicyView)
	}
	events, esErr := repo.getOrgEvents(ctx, repo.SystemDefaults.IamID, policy.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return repo.GetDefaultPasswordLockoutPolicy(ctx)
	}
	if esErr != nil {
		logging.Log("EVENT-mS9od").WithError(esErr).Debug("error retrieving new events")
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

func (repo *OrgRepository) GetDefaultPasswordLockoutPolicy(ctx context.Context) (*iam_model.PasswordLockoutPolicyView, error) {
	policy, viewErr := repo.View.PasswordLockoutPolicyByAggregateID(repo.SystemDefaults.IamID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		policy = new(iam_es_model.PasswordLockoutPolicyView)
	}
	events, esErr := repo.getIAMEvents(ctx, policy.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-cmO9s", "Errors.IAM.PasswordLockoutPolicy.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-2Ms9f").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.PasswordLockoutViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.PasswordLockoutViewToModel(policy), nil
		}
	}
	policy.Default = true
	return iam_es_model.PasswordLockoutViewToModel(policy), nil
}

func (repo *OrgRepository) GetDefaultMailTemplate(ctx context.Context) (*iam_model.MailTemplateView, error) {
	template, err := repo.View.MailTemplateByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	template.Default = true
	return iam_es_model.MailTemplateViewToModel(template), err
}

func (repo *OrgRepository) GetMailTemplate(ctx context.Context) (*iam_model.MailTemplateView, error) {
	template, err := repo.View.MailTemplateByAggregateID(authz.GetCtxData(ctx).OrgID)
	if errors.IsNotFound(err) {
		template, err = repo.View.MailTemplateByAggregateID(repo.SystemDefaults.IamID)
		if err != nil {
			return nil, err
		}
		template.Default = true
	}
	if err != nil {
		return nil, err
	}
	return iam_es_model.MailTemplateViewToModel(template), err
}

func (repo *OrgRepository) GetDefaultMessageTexts(ctx context.Context) (*iam_model.MessageTextsView, error) {
	texts, err := repo.View.MessageTextsByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.MessageTextsViewToModel(texts, true), err
}

func (repo *OrgRepository) GetMessageTexts(ctx context.Context) (*iam_model.MessageTextsView, error) {
	defaultIn := false
	texts, err := repo.View.MessageTextsByAggregateID(authz.GetCtxData(ctx).OrgID)
	if errors.IsNotFound(err) || len(texts) == 0 {
		texts, err = repo.View.MessageTextsByAggregateID(repo.SystemDefaults.IamID)
		if err != nil {
			return nil, err
		}
		defaultIn = true
	}
	if err != nil {
		return nil, err
	}
	return iam_es_model.MessageTextsViewToModel(texts, defaultIn), err
}

func (repo *OrgRepository) GetDefaultMessageText(ctx context.Context, textType, lang string) (*iam_model.MessageTextView, error) {
	text, err := repo.View.MessageTextByIDs(repo.SystemDefaults.IamID, textType, lang)
	if err != nil {
		return nil, err
	}
	text.Default = true
	return iam_es_model.MessageTextViewToModel(text), err
}

func (repo *OrgRepository) GetMessageText(ctx context.Context, orgID, textType, lang string) (*iam_model.MessageTextView, error) {
	text, err := repo.View.MessageTextByIDs(orgID, textType, lang)
	if errors.IsNotFound(err) {
		result, err := repo.GetDefaultMessageText(ctx, textType, lang)
		if err != nil {
			return nil, err
		}
		return result, nil
	}
	if err != nil {
		return nil, err
	}
	return iam_es_model.MessageTextViewToModel(text), err
}

func (repo *OrgRepository) getOrgChanges(ctx context.Context, orgID string, lastSequence uint64, limit uint64, sortAscending bool, auditLogRetention time.Duration) (*org_model.OrgChanges, error) {
	query := org_view.ChangesQuery(orgID, lastSequence, limit, sortAscending, auditLogRetention)

	events, err := repo.Eventstore.FilterEvents(context.Background(), query)
	if err != nil {
		logging.Log("EVENT-ZRffs").WithError(err).Warn("eventstore unavailable")
		return nil, errors.ThrowInternal(err, "EVENT-328b1", "Errors.Org.NotFound")
	}
	if len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-FpQqK", "Errors.Changes.NotFound")
	}

	changes := make([]*org_model.OrgChange, len(events))

	for i, event := range events {
		creationDate, err := ptypes.TimestampProto(event.CreationDate)
		logging.Log("EVENT-qxIR7").OnError(err).Debug("unable to parse timestamp")
		change := &org_model.OrgChange{
			ChangeDate: creationDate,
			EventType:  event.Type.String(),
			ModifierId: event.EditorUser,
			Sequence:   event.Sequence,
		}

		if event.Data != nil {
			org := new(org_es_model.Org)
			err := json.Unmarshal(event.Data, org)
			logging.Log("EVENT-XCLEm").OnError(err).Debug("unable to unmarshal data")
			change.Data = org
		}

		changes[i] = change
		if lastSequence < event.Sequence {
			lastSequence = event.Sequence
		}
	}

	return &org_model.OrgChanges{
		Changes:      changes,
		LastSequence: lastSequence,
	}, nil
}

func (repo *OrgRepository) userByID(ctx context.Context, id string) (*usr_model.UserView, error) {
	user, viewErr := repo.View.UserByID(id)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		user = new(usr_es_model.UserView)
	}
	events, esErr := repo.getUserEvents(ctx, id, user.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-3nF8s", "Errors.User.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-PSoc3").WithError(esErr).Debug("error retrieving new events")
		return usr_es_model.UserToModel(user), nil
	}
	userCopy := *user
	for _, event := range events {
		if err := userCopy.AppendEvent(event); err != nil {
			return usr_es_model.UserToModel(user), nil
		}
	}
	if userCopy.State == int32(usr_es_model.UserStateDeleted) {
		return nil, errors.ThrowNotFound(nil, "EVENT-3n8Fs", "Errors.User.NotFound")
	}
	return usr_es_model.UserToModel(&userCopy), nil
}

func (r *OrgRepository) getUserEvents(ctx context.Context, userID string, sequence uint64) ([]*models.Event, error) {
	query, err := view.UserByIDQuery(userID, sequence)
	if err != nil {
		return nil, err
	}

	return r.Eventstore.FilterEvents(ctx, query)
}

func (es *OrgRepository) getOrgEvents(ctx context.Context, id string, sequence uint64) ([]*models.Event, error) {
	query, err := org_view.OrgByIDQuery(id, sequence)
	if err != nil {
		return nil, err
	}
	return es.Eventstore.FilterEvents(ctx, query)
}

func (repo *OrgRepository) getIAMEvents(ctx context.Context, sequence uint64) ([]*models.Event, error) {
	query, err := iam_view.IAMByIDQuery(domain.IAMID, sequence)
	if err != nil {
		return nil, err
	}
	return repo.Eventstore.FilterEvents(ctx, query)
}
