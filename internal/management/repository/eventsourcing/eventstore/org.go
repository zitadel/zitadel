package eventstore

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/caos/logging"
	"github.com/golang/protobuf/ptypes"
	"golang.org/x/text/language"

	"github.com/caos/zitadel/internal/config/systemdefaults"
	"github.com/caos/zitadel/internal/domain"
	"github.com/caos/zitadel/internal/errors"
	v1 "github.com/caos/zitadel/internal/eventstore/v1"
	"github.com/caos/zitadel/internal/i18n"
	org_model "github.com/caos/zitadel/internal/org/model"
	org_es_model "github.com/caos/zitadel/internal/org/repository/eventsourcing/model"
	org_view "github.com/caos/zitadel/internal/org/repository/view"
	"github.com/caos/zitadel/internal/query"
)

type OrgRepository struct {
	Query                               *query.Queries
	SearchLimit                         uint64
	Eventstore                          v1.Eventstore
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
		user, _ := repo.Query.GetUserByID(ctx, change.ModifierId)
		if user != nil {
			change.ModifierLoginName = user.PreferredLoginName
			if user.Human != nil {
				change.ModifierName = user.Human.DisplayName
				change.ModifierAvatarURL = domain.AvatarURL(repo.PrefixAvatarURL, user.ResourceOwner, user.Human.AvatarKey)
			}
			if user.Machine != nil {
				change.ModifierName = user.Machine.Name
			}
		}
	}
	return changes, nil
}

func (repo *OrgRepository) GetOrgMemberRoles(isGlobal bool) []string {
	roles := make([]string, 0)
	for _, roleMap := range repo.Roles {
		if strings.HasPrefix(roleMap, "ORG") {
			roles = append(roles, roleMap)
		}
	}
	if isGlobal {
		roles = append(roles, domain.RoleSelfManagementGlobal)
	}
	return roles
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
