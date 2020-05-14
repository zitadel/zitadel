package handler

import (
	"encoding/json"
	"time"

	caos_errs "github.com/caos/zitadel/internal/errors"
	es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"

	"github.com/caos/logging"

	"github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	"github.com/caos/zitadel/internal/user/repository/eventsourcing"
	user_events "github.com/caos/zitadel/internal/user/repository/eventsourcing"
)

type Token struct {
	handler
	userEvents *user_events.UserEventstore
}

const (
	tokenTable = "auth.token"
)

func (u *Token) MinimumCycleDuration() time.Duration { return u.cycleDuration }

func (u *Token) ViewModel() string {
	return tokenTable
}

func (u *Token) EventQuery() (*models.SearchQuery, error) {
	sequence, err := u.view.GetLatestTokenSequence()
	if err != nil {
		return nil, err
	}
	return eventsourcing.UserQuery(sequence), nil
}

func (u *Token) Process(event *models.Event) (err error) {
	switch event.Type {
	case es_model.SignedOut:
		id, err := agentIDFromSession(event)
		if err != nil {
			return err
		}
		err = u.view.DeleteSessionTokens(id, event.AggregateID, event.Sequence)
		if err != nil {
			return u.view.ProcessedTokenSequence(event.Sequence)
		}
	default:
		return u.view.ProcessedTokenSequence(event.Sequence)
	}
	return nil
}

func (u *Token) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-3jkl4", "id", event.AggregateID).WithError(err).Warn("something went wrong in token handler")
	return spooler.HandleError(event, err, u.view.GetLatestUserFailedEvent, u.view.ProcessedUserFailedEvent, u.view.ProcessedUserSequence, u.errorCountUntilSkip)
}

func agentIDFromSession(event *models.Event) (string, error) {
	session := make(map[string]interface{})
	if err := json.Unmarshal(event.Data, session); err != nil {
		logging.Log("EVEN-s3bq9").WithError(err).Error("could not unmarshal event data")
		return "", caos_errs.ThrowInternal(nil, "MODEL-sd325", "could not unmarshal data")
	}
	return session["agentID"].(string), nil
}
