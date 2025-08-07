package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	user_repo "github.com/zitadel/zitadel/internal/repository/user"
	usr_model "github.com/zitadel/zitadel/internal/user/model"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	TokenKeyTokenID           = "id"
	TokenKeyUserID            = "user_id"
	TokenKeyRefreshTokenID    = "refresh_token_id"
	TokenKeyApplicationID     = "application_id"
	TokenKeyUserAgentID       = "user_agent_id"
	TokenKeyExpiration        = "expiration"
	TokenKeyResourceOwner     = "resource_owner"
	TokenKeyInstanceID        = "instance_id"
	TokenKeyCreationDate      = "creation_date"
	TokenKeyChangeDate        = "change_date"
	TokenKeySequence          = "sequence"
	TokenKeyActor             = "actor"
	TokenKeyID                = "id"
	TokenKeyAudience          = "audience"
	TokenKeyPreferredLanguage = "preferred_language"
	TokenKeyScopes            = "scopes"
	TokenKeyIsPat             = "is_pat"
)

type TokenView struct {
	ID                string                     `json:"tokenId" gorm:"column:id;primary_key"`
	CreationDate      time.Time                  `json:"-" gorm:"column:creation_date"`
	ChangeDate        time.Time                  `json:"-" gorm:"column:change_date"`
	ResourceOwner     string                     `json:"-" gorm:"column:resource_owner"`
	UserID            string                     `json:"-" gorm:"column:user_id"`
	ApplicationID     string                     `json:"applicationId" gorm:"column:application_id"`
	UserAgentID       string                     `json:"userAgentId" gorm:"column:user_agent_id"`
	Audience          database.TextArray[string] `json:"audience" gorm:"column:audience"`
	Scopes            database.TextArray[string] `json:"scopes" gorm:"column:scopes"`
	Expiration        time.Time                  `json:"expiration" gorm:"column:expiration"`
	Sequence          uint64                     `json:"-" gorm:"column:sequence"`
	PreferredLanguage string                     `json:"preferredLanguage" gorm:"column:preferred_language"`
	RefreshTokenID    string                     `json:"refreshTokenID,omitempty" gorm:"refresh_token_id"`
	IsPAT             bool                       `json:"-" gorm:"is_pat"`
	Deactivated       bool                       `json:"-" gorm:"-"`
	InstanceID        string                     `json:"instanceID" gorm:"column:instance_id;primary_key"`
	Actor             TokenActor                 `json:"actor" gorm:"column:actor"`
}

type TokenActor struct {
	*domain.TokenActor
}

func (a *TokenActor) Scan(value any) error {
	var data []byte
	switch v := value.(type) {
	case nil:
		a.TokenActor = nil
	case string:
		data = []byte(v)
	case []byte:
		data = v
	default:
		return zerrors.ThrowInternalf(nil, "MODEL-yo8Ae", "cannot scan type %T into %T", v, a)
	}
	if err := json.Unmarshal(data, &a.TokenActor); err != nil {
		return zerrors.ThrowInternal(nil, "MODEL-yo8Ae", "cannot unmarshal token actor")
	}
	return nil
}

func (a TokenActor) Value() (driver.Value, error) {
	if a.TokenActor == nil {
		return nil, nil
	}
	data, err := json.Marshal(a.TokenActor)
	if err != nil {
		return nil, zerrors.ThrowInternal(nil, "MODEL-oD2mi", "cannot marshal token actor")
	}
	return data, nil
}

func TokenViewToModel(token *TokenView) *usr_model.TokenView {
	return &usr_model.TokenView{
		ID:                token.ID,
		CreationDate:      token.CreationDate,
		ChangeDate:        token.ChangeDate,
		ResourceOwner:     token.ResourceOwner,
		UserID:            token.UserID,
		ApplicationID:     token.ApplicationID,
		UserAgentID:       token.UserAgentID,
		Audience:          token.Audience,
		Scopes:            token.Scopes,
		Expiration:        token.Expiration,
		Sequence:          token.Sequence,
		PreferredLanguage: token.PreferredLanguage,
		RefreshTokenID:    token.RefreshTokenID,
		IsPAT:             token.IsPAT,
		Actor:             token.Actor.TokenActor,
	}
}

func (t *TokenView) AppendEventIfMyToken(event eventstore.Event) (err error) {
	// in case anything needs to be change here check if the Reduce function needs the change as well
	view := new(TokenView)
	switch event.Type() {
	case user_repo.UserTokenAddedType,
		user_repo.PersonalAccessTokenAddedType:
		view.setRootData(event)
		err = view.setData(event)
	case user_repo.UserTokenRemovedType:
		return t.appendTokenRemoved(event)
	case user_repo.HumanRefreshTokenRemovedType:
		return t.appendRefreshTokenRemoved(event)
	case user_repo.UserV1SignedOutType,
		user_repo.HumanSignedOutType:
		id, err := UserAgentIDFromEvent(event)
		if err != nil {
			return err
		}
		if t.UserAgentID == id {
			t.Deactivated = true
		}
		return nil
	case user_repo.UserRemovedType,
		user_repo.UserDeactivatedType,
		user_repo.UserLockedType:
		t.Deactivated = true
		return nil
	case user_repo.UserUnlockedType,
		user_repo.UserReactivatedType:
		if t.ID != "" && event.CreatedAt().Before(t.CreationDate) {
			t.Deactivated = false
		}
		return nil
	case user_repo.PersonalAccessTokenRemovedType:
		return t.appendPATRemoved(event)
	default:
		return nil
	}
	if err != nil {
		return err
	}
	if view.ID == t.ID {
		return t.AppendEvent(event)
	}
	return nil
}

func (t *TokenView) AppendEvent(event eventstore.Event) error {
	// in case anything needs to be change here check if the Reduce function needs the change as well
	t.ChangeDate = event.CreatedAt()
	t.Sequence = event.Sequence()
	switch event.Type() {
	case user_repo.UserTokenAddedType,
		user_repo.PersonalAccessTokenAddedType:
		t.setRootData(event)
		err := t.setData(event)
		if err != nil {
			return err
		}
		t.CreationDate = event.CreatedAt()
		t.IsPAT = event.Type() == user_repo.PersonalAccessTokenAddedType
	}
	return nil
}

func (t *TokenView) setRootData(event eventstore.Event) {
	t.UserID = event.Aggregate().ID
	t.ResourceOwner = event.Aggregate().ResourceOwner
	t.InstanceID = event.Aggregate().InstanceID
}

func (t *TokenView) setData(event eventstore.Event) error {
	if err := event.Unmarshal(t); err != nil {
		logging.WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(err, "MODEL-5Gms9", "could not unmarshal event")
	}
	return nil
}

func (t *TokenView) appendTokenRemoved(event eventstore.Event) error {
	tokenID, err := tokenIDFromEvent(event)
	if err != nil {
		return err
	}
	if tokenID == t.ID {
		t.Deactivated = true
	}
	return nil
}

func (t *TokenView) appendRefreshTokenRemoved(event eventstore.Event) error {
	tokenID, err := tokenIDFromEvent(event)
	if err != nil {
		return err
	}
	if tokenID == t.RefreshTokenID {
		t.Deactivated = true
	}
	return nil
}

func (t *TokenView) appendPATRemoved(event eventstore.Event) error {
	tokenID, err := tokenIDFromEvent(event)
	if err != nil {
		return err
	}
	if tokenID == t.ID && t.IsPAT {
		t.Deactivated = true
	}
	return nil
}

func (t *TokenView) GetRelevantEventTypes() []eventstore.EventType {
	return []eventstore.EventType{
		user_repo.UserTokenAddedType,
		user_repo.PersonalAccessTokenAddedType,
		user_repo.UserTokenRemovedType,
		user_repo.HumanRefreshTokenRemovedType,
		user_repo.UserV1SignedOutType,
		user_repo.HumanSignedOutType,
		user_repo.UserRemovedType,
		user_repo.UserDeactivatedType,
		user_repo.UserLockedType,
		user_repo.UserLockedType,
		user_repo.UserReactivatedType,
		user_repo.PersonalAccessTokenRemovedType,
	}
}

type tokenIDPayload struct {
	ID string `json:"tokenId"`
}

func tokenIDFromEvent(event eventstore.Event) (string, error) {
	m := new(tokenIDPayload)
	if err := event.Unmarshal(&m); err != nil {
		logging.WithError(err).Error("could not unmarshal event data")
		return "", zerrors.ThrowInternal(nil, "MODEL-SDAfw", "could not unmarshal data")
	}
	return m.ID, nil
}
