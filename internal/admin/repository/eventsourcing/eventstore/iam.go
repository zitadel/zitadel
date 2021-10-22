package eventstore

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"

	"github.com/ghodss/yaml"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/domain"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/eventstore/v1/models"
	"github.com/caos/zitadel/internal/i18n"
	iam_view "github.com/caos/zitadel/internal/iam/repository/view"
	"github.com/caos/zitadel/internal/query"
	"github.com/caos/zitadel/internal/user/repository/view/model"

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
	Query                               *query.Queries
	Eventstore                          v1.Eventstore
	SearchLimit                         uint64
	View                                *admin_view.View
	SystemDefaults                      systemdefaults.SystemDefaults
	Roles                               []string
	PrefixAvatarURL                     string
	LoginDir                            http.FileSystem
	NotificationDir                     http.FileSystem
	LoginTranslationFileContents        map[string][]byte
	NotificationTranslationFileContents map[string][]byte
	mutex                               sync.Mutex
	supportedLangs                      []language.Tag
}

func (repo *IAMRepository) Languages(ctx context.Context) ([]language.Tag, error) {
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

func (repo *IAMRepository) IAMMemberByID(ctx context.Context, iamID, userID string) (*iam_model.IAMMemberView, error) {
	member, err := repo.View.IAMMemberByIDs(iamID, userID)
	if err != nil {
		return nil, err
	}
	return iam_es_model.IAMMemberToModel(member, repo.PrefixAvatarURL), nil
}

func (repo *IAMRepository) SearchIAMMembers(ctx context.Context, request *iam_model.IAMMemberSearchRequest) (*iam_model.IAMMemberSearchResponse, error) {
	err := request.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
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
		Result:      iam_es_model.IAMMembersToModel(members, repo.PrefixAvatarURL),
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

func (repo *IAMRepository) SearchIDPConfigs(ctx context.Context, request *iam_model.IDPConfigSearchRequest) (*iam_model.IDPConfigSearchResponse, error) {
	err := request.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
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

func (repo *IAMRepository) SearchDefaultIDPProviders(ctx context.Context, request *iam_model.IDPProviderSearchRequest) (*iam_model.IDPProviderSearchResponse, error) {
	err := request.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
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

func (repo *IAMRepository) GetDefaultLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error) {
	policy, err := repo.View.LabelPolicyByAggregateIDAndState(repo.SystemDefaults.IamID, int32(domain.LabelPolicyStateActive))
	if err != nil {
		return nil, err
	}
	return iam_es_model.LabelPolicyViewToModel(policy), err
}

func (repo *IAMRepository) GetDefaultPreviewLabelPolicy(ctx context.Context) (*iam_model.LabelPolicyView, error) {
	policy, err := repo.View.LabelPolicyByAggregateIDAndState(repo.SystemDefaults.IamID, int32(domain.LabelPolicyStatePreview))
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
	err := request.EnsureLimit(repo.SearchLimit)
	if err != nil {
		return nil, err
	}
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
		Result:      iam_es_model.IAMMembersToModel(members, repo.PrefixAvatarURL),
	}
	if err == nil {
		result.Sequence = sequence.CurrentSequence
	}
	return result, nil
}

func (repo *IAMRepository) GetDefaultMessageText(ctx context.Context, textType, lang string) (*domain.CustomMessageText, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	var err error
	contents, ok := repo.NotificationTranslationFileContents[lang]
	if !ok {
		contents, err = repo.readTranslationFile(repo.NotificationDir, fmt.Sprintf("/i18n/%s.yaml", lang))
		if caos_errs.IsNotFound(err) {
			contents, err = repo.readTranslationFile(repo.NotificationDir, fmt.Sprintf("/i18n/%s.yaml", repo.SystemDefaults.DefaultLanguage.String()))
		}
		if err != nil {
			return nil, err
		}
		repo.NotificationTranslationFileContents[lang] = contents
	}
	messageTexts := new(domain.MessageTexts)
	if err := yaml.Unmarshal(contents, messageTexts); err != nil {
		return nil, caos_errs.ThrowInternal(err, "TEXT-3N9fs", "Errors.TranslationFile.ReadError")
	}
	return messageTexts.GetMessageTextByType(textType), nil
}

func (repo *IAMRepository) GetCustomMessageText(ctx context.Context, textType, lang string) (*domain.CustomMessageText, error) {
	texts, err := repo.View.CustomTextsByAggregateIDAndTemplateAndLand(repo.SystemDefaults.IamID, textType, lang)
	if err != nil {
		return nil, err
	}
	result := iam_es_model.CustomTextViewsToMessageDomain(repo.SystemDefaults.IamID, lang, texts)
	result.Default = true
	return result, err
}

func (repo *IAMRepository) GetDefaultLoginTexts(ctx context.Context, lang string) (*domain.CustomLoginText, error) {
	repo.mutex.Lock()
	defer repo.mutex.Unlock()
	contents, ok := repo.LoginTranslationFileContents[lang]
	var err error
	if !ok {
		contents, err = repo.readTranslationFile(repo.LoginDir, fmt.Sprintf("/i18n/%s.yaml", lang))
		if caos_errs.IsNotFound(err) {
			contents, err = repo.readTranslationFile(repo.LoginDir, fmt.Sprintf("/i18n/%s.yaml", repo.SystemDefaults.DefaultLanguage.String()))
		}
		if err != nil {
			return nil, err
		}
		repo.LoginTranslationFileContents[lang] = contents
	}
	loginText := new(domain.CustomLoginText)
	if err := yaml.Unmarshal(contents, loginText); err != nil {
		return nil, caos_errs.ThrowInternal(err, "TEXT-GHR3Q", "Errors.TranslationFile.ReadError")
	}
	return loginText, nil
}

func (repo *IAMRepository) GetCustomLoginTexts(ctx context.Context, lang string) (*domain.CustomLoginText, error) {
	texts, err := repo.View.CustomTextsByAggregateIDAndTemplateAndLand(repo.SystemDefaults.IamID, domain.LoginCustomText, lang)
	if err != nil {
		return nil, err
	}
	return iam_es_model.CustomTextViewsToLoginDomain(repo.SystemDefaults.IamID, lang, texts), err
}

func (repo *IAMRepository) getIAMEvents(ctx context.Context, sequence uint64) ([]*models.Event, error) {
	query, err := iam_view.IAMByIDQuery(domain.IAMID, sequence)
	if err != nil {
		return nil, err
	}
	return repo.Eventstore.FilterEvents(ctx, query)
}

func (repo *IAMRepository) readTranslationFile(dir http.FileSystem, filename string) ([]byte, error) {
	r, err := dir.Open(filename)
	if os.IsNotExist(err) {
		return nil, caos_errs.ThrowNotFound(err, "TEXT-3n9fs", "Errors.TranslationFile.NotFound")
	}
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "TEXT-93njw", "Errors.TranslationFile.ReadError")
	}
	contents, err := ioutil.ReadAll(r)
	if err != nil {
		return nil, caos_errs.ThrowInternal(err, "TEXT-l0fse", "Errors.TranslationFile.ReadError")
	}
	return contents, nil
}
