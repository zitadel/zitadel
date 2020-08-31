package handler

import (
	"context"
	"encoding/json"

	"github.com/caos/logging"
	caos_errs "github.com/caos/zitadel/internal/errors"
	"github.com/caos/zitadel/internal/eventstore/models"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/eventstore/spooler"
	proj_event "github.com/caos/zitadel/internal/project/repository/eventsourcing"
	project_es_model "github.com/caos/zitadel/internal/project/repository/eventsourcing/model"
	user_es_model "github.com/caos/zitadel/internal/user/repository/eventsourcing/model"
	view_model "github.com/caos/zitadel/internal/user/repository/view/model"
)

type Token struct {
	handler
	ProjectEvents *proj_event.ProjectEventstore
}

const (
	tokenTable = "auth.tokens"
)

func (u *Token) ViewModel() string {
	return tokenTable
}

func (u *Token) EventQuery() (*models.SearchQuery, error) {
	sequence, err := u.view.GetLatestTokenSequence()
	if err != nil {
		return nil, err
	}
	if err != nil {
		return nil, err
	}
	return es_models.NewSearchQuery().
		AggregateTypeFilter(user_es_model.UserAggregate, project_es_model.ProjectAggregate).
		LatestSequenceFilter(sequence.CurrentSequence), nil
}

func (u *Token) Reduce(event *models.Event) (err error) {
	switch event.Type {
	case user_es_model.UserProfileChanged,
		user_es_model.HumanProfileChanged:
		user := new(view_model.UserView)
		user.AppendEvent(event)
		tokens, err := u.view.TokensByUserID(event.AggregateID)
		if err != nil {
			return err
		}
		for _, token := range tokens {
			token.PreferredLanguage = user.PreferredLanguage
		}
		return u.view.PutTokens(tokens, event.Sequence)
	case user_es_model.SignedOut,
		user_es_model.HumanSignedOut:
		id, err := agentIDFromSession(event)
		if err != nil {
			return err
		}
		return u.view.DeleteSessionTokens(id, event.AggregateID, event.Sequence)
	case user_es_model.UserLocked,
		user_es_model.UserDeactivated,
		user_es_model.UserRemoved:
		return u.view.DeleteUserTokens(event.AggregateID, event.Sequence)
	case project_es_model.ApplicationDeactivated,
		project_es_model.ApplicationRemoved:
		application, err := applicationFromSession(event)
		if err != nil {
			return err
		}
		return u.view.DeleteApplicationTokens(event.Sequence, application.AppID)
	case project_es_model.ProjectDeactivated,
		project_es_model.ProjectRemoved:
		project, err := u.ProjectEvents.ProjectByID(context.Background(), event.AggregateID)
		if err != nil {
			return err
		}
		applicationsIDs := make([]string, 0, len(project.Applications))
		for _, app := range project.Applications {
			applicationsIDs = append(applicationsIDs, app.AppID)
		}
		return u.view.DeleteApplicationTokens(event.Sequence, applicationsIDs...)
	default:
		return u.view.ProcessedTokenSequence(event.Sequence)
	}
}

func (u *Token) OnError(event *models.Event, err error) error {
	logging.LogWithFields("SPOOL-3jkl4", "id", event.AggregateID).WithError(err).Warn("something went wrong in token handler")
	return spooler.HandleError(event, err, u.view.GetLatestTokenFailedEvent, u.view.ProcessedTokenFailedEvent, u.view.ProcessedTokenSequence, u.errorCountUntilSkip)
}

func agentIDFromSession(event *models.Event) (string, error) {
	session := make(map[string]interface{})
	if err := json.Unmarshal(event.Data, &session); err != nil {
		logging.Log("EVEN-s3bq9").WithError(err).Error("could not unmarshal event data")
		return "", caos_errs.ThrowInternal(nil, "MODEL-sd325", "could not unmarshal data")
	}
	return session["userAgentID"].(string), nil
}

func applicationFromSession(event *models.Event) (*project_es_model.Application, error) {
	application := new(project_es_model.Application)
	if err := json.Unmarshal(event.Data, &application); err != nil {
		logging.Log("EVEN-GRE2q").WithError(err).Error("could not unmarshal event data")
		return nil, caos_errs.ThrowInternal(nil, "MODEL-Hrw1q", "could not unmarshal data")
	}
	return application, nil
}
