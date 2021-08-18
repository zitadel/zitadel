package eventstore

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/caos/logging"
	"github.com/ghodss/yaml"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/api/authz"
	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/i18n"
	iam_model "github.com/caos/zitadel/internal/iam/model"
	iam_view "github.com/caos/zitadel/internal/iam/repository/view"
	iam_es_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	mgmt_view "github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	org_view "github.com/caos/zitadel/internal/org/repository/view"
	"github.com/caos/zitadel/internal/org/repository/view/model"
	"github.com/caos/zitadel/internal/telemetry/tracing"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type OrgRepository struct {
	SearchLimit                         uint64
	Eventstore                          v1.Eventstore
	View                                *mgmt_view.View
	Roles                               []string
	SystemDefaults                      systemdefaults.SystemDefaults
	PrefixAvatarURL                     string
	LoginDir                            http.FileSystem
	NotificationDir                     http.FileSystem
	LoginTranslationFileContents        map[string][]byte
	NotificationTranslationFileContents map[string][]byte
	mutex                               sync.Mutex
	supportedLangs                      []language.Tag
}

func (repo *OrgRepository) Languages(ctx context.Context) ([]language.Tag, error) {
	if len(repo.supportedLangs) == 0 {
		langs, err := i18n.SupportedLanguages(repo.LoginDir)
		if err != nil {
			logging.Log("ADMIN-tiMWs").WithError(err).Debug("unable to parse language")
			return nil, err
		}
		repo.supportedLangs = langs
	}
	return repo.supportedLangs, nil
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
		change.ModifierLoginName = change.ModifierId
		user, _ := repo.userByID(ctx, change.ModifierId)
		if user != nil {
			change.ModifierLoginName = user.PreferredLoginName
			if user.HumanView != nil {
				change.ModifierName = user.HumanView.DisplayName
				change.ModifierAvatarURL = user.HumanView.AvatarURL
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
	return model.OrgMemberToModel(member, repo.PrefixAvatarURL), nil
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
		Result:      model.OrgMembersToModel(members, repo.PrefixAvatarURL),
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
	policy, err := repo.View.LabelPolicyByAggregateIDAndState(authz.GetCtxData(ctx).OrgID, int32(domain.LabelPolicyStateActive))
	if errors.IsNotFound(err) {
		policy, err = repo.View.LabelPolicyByAggregateIDAndState(repo.SystemDefaults.IamID, int32(domain.LabelPolicyStateActive))
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

func (repo *OrgRepository) GetPreviewLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error) {
	policy, err := repo.View.LabelPolicyByAggregateIDAndState(authz.GetCtxData(ctx).OrgID, int32(domain.LabelPolicyStatePreview))
	if errors.IsNotFound(err) {
		policy, err = repo.View.LabelPolicyByAggregateIDAndState(repo.SystemDefaults.IamID, int32(domain.LabelPolicyStatePreview))
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
	return repo.getDefaultLabelPolicy(ctx, domain.LabelPolicyStateActive)
}

func (repo *OrgRepository) GetPreviewDefaultLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error) {
	return repo.getDefaultLabelPolicy(ctx, domain.LabelPolicyStatePreview)
}

func (repo *OrgRepository) getDefaultLabelPolicy(ctx context.Context, state domain.LabelPolicyState) (*iam_model.LabelPolicyView, error) {
	policy, viewErr := repo.View.LabelPolicyByAggregateIDAndState(repo.SystemDefaults.IamID, int32(state))
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

func (repo *OrgRepository) GetLockoutPolicy(ctx context.Context) (*iam_model.LockoutPolicyView, error) {
	policy, viewErr := repo.View.LockoutPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		policy = new(iam_es_model.LockoutPolicyView)
	}
	events, esErr := repo.getOrgEvents(ctx, repo.SystemDefaults.IamID, policy.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return repo.GetDefaultLockoutPolicy(ctx)
	}
	if esErr != nil {
		logging.Log("EVENT-mS9od").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.LockoutViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.LockoutViewToModel(policy), nil
		}
	}
	return iam_es_model.LockoutViewToModel(policy), nil
}

func (repo *OrgRepository) GetDefaultLockoutPolicy(ctx context.Context) (*iam_model.LockoutPolicyView, error) {
	policy, viewErr := repo.View.LockoutPolicyByAggregateID(repo.SystemDefaults.IamID)
	if viewErr != nil && !errors.IsNotFound(viewErr) {
		return nil, viewErr
	}
	if errors.IsNotFound(viewErr) {
		policy = new(iam_es_model.LockoutPolicyView)
	}
	events, esErr := repo.getIAMEvents(ctx, policy.Sequence)
	if errors.IsNotFound(viewErr) && len(events) == 0 {
		return nil, errors.ThrowNotFound(nil, "EVENT-cmO9s", "Errors.IAM.LockoutPolicy.NotFound")
	}
	if esErr != nil {
		logging.Log("EVENT-2Ms9f").WithError(esErr).Debug("error retrieving new events")
		return iam_es_model.LockoutViewToModel(policy), nil
	}
	policyCopy := *policy
	for _, event := range events {
		if err := policyCopy.AppendEvent(event); err != nil {
			return iam_es_model.LockoutViewToModel(policy), nil
		}
	}
	policy.Default = true
	return iam_es_model.LockoutViewToModel(policy), nil
}

func (repo *OrgRepository) GetPrivacyPolicy(ctx context.Context) (*iam_model.PrivacyPolicyView, error) {
	policy, err := repo.View.PrivacyPolicyByAggregateID(authz.GetCtxData(ctx).OrgID)
	if errors.IsNotFound(err) {
		return repo.GetDefaultPrivacyPolicy(ctx)
	}
	return iam_es_model.PrivacyViewToModel(policy), nil
}

func (repo *OrgRepository) GetDefaultPrivacyPolicy(ctx context.Context) (*iam_model.PrivacyPolicyView, error) {
	policy, err := repo.View.PrivacyPolicyByAggregateID(repo.SystemDefaults.IamID)
	if err != nil {
		return nil, err
	}
	policy.Default = true
	return iam_es_model.PrivacyViewToModel(policy), nil
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

func (repo *OrgRepository) GetDefaultMessageText(ctx context.Context, textType, lang string) (*domain.CustomMessageText, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	var err error
	contents, ok := repo.NotificationTranslationFileContents[lang]
	if !ok {
		contents, err = repo.readTranslationFile(repo.NotificationDir, fmt.Sprintf("/i18n/%s.yaml", lang))
		if errors.IsNotFound(err) {
			contents, err = repo.readTranslationFile(repo.NotificationDir, fmt.Sprintf("/i18n/%s.yaml", repo.SystemDefaults.DefaultLanguage.String()))
		}
		if err != nil {
			return nil, err
		}
		repo.NotificationTranslationFileContents[lang] = contents
	}
	notificationTextMap := make(map[string]interface{})
	if err := yaml.Unmarshal(contents, &notificationTextMap); err != nil {
		return nil, errors.ThrowInternal(err, "TEXT-093sd", "Errors.TranslationFile.ReadError")
	}
	texts, err := repo.View.CustomTextsByAggregateIDAndTemplateAndLand(repo.SystemDefaults.IamID, textType, lang)
	if err != nil {
		return nil, err
	}
	for _, text := range texts {
		messageTextMap, ok := notificationTextMap[textType].(map[string]interface{})
		if !ok {
			continue
		}
		messageTextMap[text.Key] = text.Text
	}
	jsonbody, err := json.Marshal(notificationTextMap)
	if err != nil {
		return nil, errors.ThrowInternal(err, "TEXT-02m8f", "Errors.TranslationFile.MergeError")
	}
	notificationText := new(domain.MessageTexts)
	if err := json.Unmarshal(jsonbody, &notificationText); err != nil {
		return nil, errors.ThrowInternal(err, "TEXT-20ops", "Errors.TranslationFile.MergeError")
	}
	result := notificationText.GetMessageTextByType(textType)
	result.Default = true
	return result, nil
}

func (repo *OrgRepository) GetMessageText(ctx context.Context, orgID, textType, lang string) (*domain.CustomMessageText, error) {
	texts, err := repo.View.CustomTextsByAggregateIDAndTemplateAndLand(orgID, textType, lang)
	if err != nil {
		return nil, err
	}
	if len(texts) == 0 {
		return repo.GetDefaultMessageText(ctx, textType, lang)
	}
	return iam_es_model.CustomTextViewsToMessageDomain(repo.SystemDefaults.IamID, lang, texts), err
}

func (repo *OrgRepository) GetDefaultLoginTexts(ctx context.Context, lang string) (*domain.CustomLoginText, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	contents, ok := repo.LoginTranslationFileContents[lang]
	var err error
	if !ok {
		contents, err = repo.readTranslationFile(repo.LoginDir, fmt.Sprintf("/i18n/%s.yaml", lang))
		if errors.IsNotFound(err) {
			contents, err = repo.readTranslationFile(repo.LoginDir, fmt.Sprintf("/i18n/%s.yaml", repo.SystemDefaults.DefaultLanguage.String()))
		}
		if err != nil {
			return nil, err
		}
		repo.LoginTranslationFileContents[lang] = contents
	}
	loginTextMap := make(map[string]interface{})
	if err := yaml.Unmarshal(contents, &loginTextMap); err != nil {
		return nil, errors.ThrowInternal(err, "TEXT-l0fse", "Errors.TranslationFile.ReadError")
	}
	texts, err := repo.View.CustomTextsByAggregateIDAndTemplateAndLand(repo.SystemDefaults.IamID, domain.LoginCustomText, lang)
	if err != nil {
		return nil, err
	}
	for _, text := range texts {
		keys := strings.Split(text.Key, ".")
		screenTextMap, ok := loginTextMap[keys[0]].(map[string]interface{})
		if !ok {
			continue
		}
		screenTextMap[keys[1]] = text.Text
	}
	jsonbody, err := json.Marshal(loginTextMap)
	if err != nil {
		return nil, errors.ThrowInternal(err, "TEXT-2n8fs", "Errors.TranslationFile.MergeError")
	}
	loginText := new(domain.CustomLoginText)
	if err := json.Unmarshal(jsonbody, &loginText); err != nil {
		return nil, errors.ThrowInternal(err, "TEXT-2n8fs", "Errors.TranslationFile.MergeError")
	}
	return loginText, nil
}

func (repo *OrgRepository) GetLoginTexts(ctx context.Context, orgID, lang string) (*domain.CustomLoginText, error) {
	texts, err := repo.View.CustomTextsByAggregateIDAndTemplateAndLand(orgID, domain.LoginCustomText, lang)
	if err != nil {
		return nil, err
	}
	return iam_es_model.CustomTextViewsToLoginDomain(repo.SystemDefaults.IamID, lang, texts), err
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
		return usr_es_model.UserToModel(user, repo.PrefixAvatarURL), nil
	}
	userCopy := *user
	for _, event := range events {
		if err := userCopy.AppendEvent(event); err != nil {
			return usr_es_model.UserToModel(user, repo.PrefixAvatarURL), nil
		}
	}
	if userCopy.State == int32(usr_es_model.UserStateDeleted) {
		return nil, errors.ThrowNotFound(nil, "EVENT-3n8Fs", "Errors.User.NotFound")
	}
	return usr_es_model.UserToModel(&userCopy, repo.PrefixAvatarURL), nil
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

func (repo *OrgRepository) readTranslationFile(dir http.FileSystem, filename string) ([]byte, error) {
	r, err := dir.Open(filename)
	if os.IsNotExist(err) {
		return nil, errors.ThrowNotFound(err, "TEXT-93nfl", "Errors.TranslationFile.NotFound")
	}
	if err != nil {
		return nil, errors.ThrowInternal(err, "TEXT-3n8fs", "Errors.TranslationFile.ReadError")
	}
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, errors.ThrowInternal(err, "TEXT-322fs", "Errors.TranslationFile.ReadError")
	}
	return contents, nil
}
