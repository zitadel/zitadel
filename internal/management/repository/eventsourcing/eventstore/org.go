package eventstore

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/caos/logging"
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
	iam_view_model "github.com/caos/zitadel/internal/iam/repository/view/model"
	mgmt_view "github.com/caos/zitadel/internal/management/repository/eventsourcing/view"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	org_view "github.com/caos/zitadel/internal/org/repository/view"
	"github.com/caos/zitadel/internal/org/repository/view/model"
	"github.com/caos/zitadel/internal/query"
	usr_model "github.com/caos/zitadel/internal/user/model"
	"github.com/caos/zitadel/internal/user/repository/view"
	usr_es_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type OrgRepository struct {
	Query                               *query.Queries
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

func (repo *OrgRepository) GetIDPProvidersByIDPConfigID(ctx context.Context, aggregateID, idpConfigID string) ([]*iam_model.IDPProviderView, error) {
	idpProviders, err := repo.View.IDPProvidersByIdpConfigID(aggregateID, idpConfigID)
	if err != nil {
		return nil, err
	}
	return iam_view_model.IDPProviderViewsToModel(idpProviders), err
}

func (repo *OrgRepository) SearchIDPProviders(ctx context.Context, request *iam_model.IDPProviderSearchRequest) (*iam_model.IDPProviderSearchResponse, error) {
	policy, err := repo.Query.LoginPolicyByID(ctx, authz.GetCtxData(ctx).OrgID)
	if err != nil {
		return nil, err
	}
	if policy.IsDefault {
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
		Result:      iam_view_model.IDPProviderViewsToModel(providers),
	}
	if sequenceErr == nil {
		result.Sequence = sequence.CurrentSequence
		result.Timestamp = sequence.LastSuccessfulSpoolerRun
	}
	return result, nil
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
