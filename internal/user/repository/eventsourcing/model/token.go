package model

import (
	"encoding/json"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	caos_errs "github.com/zitadel/zitadel/internal/errors"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
	user_repo "github.com/zitadel/zitadel/internal/repository/user"
)

type Token struct {
	es_models.ObjectRoot

	TokenID           string                     `json:"tokenId" gorm:"column:token_id"`
	ApplicationID     string                     `json:"applicationId" gorm:"column:application_id"`
	UserAgentID       string                     `json:"userAgentId" gorm:"column:user_agent_id"`
	Audience          database.TextArray[string] `json:"audience" gorm:"column:audience"`
	Scopes            database.TextArray[string] `json:"scopes" gorm:"column:scopes"`
	Expiration        time.Time                  `json:"expiration" gorm:"column:expiration"`
	PreferredLanguage string                     `json:"preferredLanguage" gorm:"column:preferred_language"`
}

func (t *Token) AppendEvents(events ...*es_models.Event) error {
	for _, event := range events {
		if err := t.AppendEvent(event); err != nil {
			return err
		}
	}

	return nil
}

func (t *Token) AppendEvent(event *es_models.Event) error {
	if event.Typ == user_repo.UserTokenAddedType {
		err := t.setData(event)
		if err != nil {
			return err
		}
		t.CreationDate = event.CreationDate
	}
	return nil
}

func (t *Token) setData(event *es_models.Event) error {
	t.ObjectRoot.AppendEvent(event)
	if err := json.Unmarshal(event.Data, t); err != nil {
		logging.Log("EVEN-4Fm9s").WithError(err).Error("could not unmarshal event data")
		return caos_errs.ThrowInternal(err, "MODEL-5Gms9", "could not unmarshal event")
	}
	return nil
}
