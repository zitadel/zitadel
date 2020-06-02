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
)

type Token struct {
	handler
}

const (
	tokenTable = "authz.tokens"
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
			return err
		}
		return u.view.ProcessedTokenSequence(event.Sequence)
	default:
		return u.view.ProcessedTokenSequence(event.Sequence)
	}
	return nil
}

func (u *Token) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-udhAD", "id", event.AggregateID).WithError(err).Warn("something went wrong in token handler")
	return spooler.HandleError(event, err, u.view.GetLatestTokenFailedEvent, u.view.ProcessedTokenFailedEvent, u.view.ProcessedTokenSequence, u.errorCountUntilSkip)
}

func agentIDFromSession(event *models.Event) (string, error) {
	session := make(map[string]interface{})
	if err := json.Unmarshal(event.Data, session); err != nil {
		logging.Log("EVEN-SsdHa").WithError(err).Error("could not unmarshal event data")
		return "", caos_errs.ThrowInternal(nil, "MODEL-KDL5e", "could not unmarshal data")
	}
	return session["agentID"].(string), nil
}
