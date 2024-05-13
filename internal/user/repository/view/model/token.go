package model

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	user_repo "github.com/zitadel/zitadel/internal/repository/user"
	usr_model "github.com/zitadel/zitadel/internal/user/model"
	"github.com/zitadel/zitadel/internal/view/repository"
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
	TokenKeyPreferredLanguage = "preferred_language"
)

type TokenView struct {
	ID         string `json:"tokenId" gorm:"column:id;primary_key"`
	InstanceID string `json:"instanceID" gorm:"column:instance_id;primary_key"`

	Deactivated bool `json:"-" gorm:"-"`

	CreationDate      repository.Field[time.Time]                  `json:"-" gorm:"column:creation_date"`
	ChangeDate        repository.Field[time.Time]                  `json:"-" gorm:"column:change_date"`
	ResourceOwner     repository.Field[string]                     `json:"-" gorm:"column:resource_owner"`
	UserID            repository.Field[string]                     `json:"-" gorm:"column:user_id"`
	ApplicationID     repository.Field[string]                     `json:"applicationId" gorm:"column:application_id"`
	UserAgentID       repository.Field[string]                     `json:"userAgentId" gorm:"column:user_agent_id"`
	Audience          repository.Field[database.TextArray[string]] `json:"audience" gorm:"column:audience"`
	Scopes            repository.Field[database.TextArray[string]] `json:"scopes" gorm:"column:scopes"`
	Expiration        repository.Field[time.Time]                  `json:"expiration" gorm:"column:expiration"`
	Sequence          repository.Field[uint64]                     `json:"-" gorm:"column:sequence"`
	PreferredLanguage repository.Field[string]                     `json:"preferredLanguage" gorm:"column:preferred_language"`
	RefreshTokenID    repository.Field[string]                     `json:"refreshTokenID,omitempty" gorm:"refresh_token_id"`
	IsPAT             repository.Field[bool]                       `json:"-" gorm:"is_pat"`
	Actor             repository.Field[TokenActor]                 `json:"actor" gorm:"column:actor"`
}

func (v *TokenView) PKColumns() []handler.Column {
	return []handler.Column{
		handler.NewCol("id", v.ID),
		handler.NewCol("instance_id", v.InstanceID),
	}
}

func (v *TokenView) PKConditions() []handler.Condition {
	return []handler.Condition{
		handler.NewCond("id", v.ID),
		handler.NewCond("instance_id", v.InstanceID),
	}
}

func (v *TokenView) Changes() []handler.Column {
	changes := make([]handler.Column, 0, 12)

	if v.CreationDate.DidChange() {
		changes = append(changes, handler.NewCol("creation_date", v.CreationDate.Value()))
	}
	if v.ChangeDate.DidChange() {
		changes = append(changes, handler.NewCol("change_date", v.ChangeDate.Value()))
	}
	if v.ResourceOwner.DidChange() {
		changes = append(changes, handler.NewCol("resource_owner", v.ResourceOwner.Value()))
	}
	if v.UserID.DidChange() {
		changes = append(changes, handler.NewCol("user_id", v.UserID.Value()))
	}
	if v.ApplicationID.DidChange() {
		changes = append(changes, handler.NewCol("application_id", v.ApplicationID.Value()))
	}
	if v.UserAgentID.DidChange() {
		changes = append(changes, handler.NewCol("user_agent_id", v.UserAgentID.Value()))
	}
	if v.Audience.DidChange() {
		changes = append(changes, handler.NewCol("audience", v.Audience.Value()))
	}
	if v.Scopes.DidChange() {
		changes = append(changes, handler.NewCol("scopes", v.Scopes.Value()))
	}
	if v.Expiration.DidChange() {
		changes = append(changes, handler.NewCol("expiration", v.Expiration.Value()))
	}
	if v.Sequence.DidChange() {
		changes = append(changes, handler.NewCol("sequence", v.Sequence.Value()))
	}
	if v.PreferredLanguage.DidChange() {
		changes = append(changes, handler.NewCol("preferred_language", v.PreferredLanguage.Value()))
	}
	if v.RefreshTokenID.DidChange() {
		changes = append(changes, handler.NewCol("refresh_token_id", v.RefreshTokenID.Value()))
	}
	if v.IsPAT.DidChange() {
		changes = append(changes, handler.NewCol("is_pat", v.IsPAT.Value()))
	}
	if v.Actor.DidChange() {
		changes = append(changes, handler.NewCol("actor", v.Actor.Value()))
	}

	return changes
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
		CreationDate:      token.CreationDate.Value(),
		ChangeDate:        token.ChangeDate.Value(),
		ResourceOwner:     token.ResourceOwner.Value(),
		UserID:            token.UserID.Value(),
		ApplicationID:     token.ApplicationID.Value(),
		UserAgentID:       token.UserAgentID.Value(),
		Audience:          token.Audience.Value(),
		Scopes:            token.Scopes.Value(),
		Expiration:        token.Expiration.Value(),
		Sequence:          token.Sequence.Value(),
		PreferredLanguage: token.PreferredLanguage.Value(),
		RefreshTokenID:    token.RefreshTokenID.Value(),
		IsPAT:             token.IsPAT.Value(),
		Actor:             token.Actor.Value().TokenActor,
	}
}

func (t *TokenView) AppendEventIfMyToken(event eventstore.Event) (err error) {
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
		id, err := agentIDFromSession(event)
		if err != nil {
			return err
		}
		if t.UserAgentID.Value() == id {
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
		if t.ID != "" && event.CreatedAt().Before(t.CreationDate.Value()) {
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
	t.ChangeDate.Set(event.CreatedAt())
	t.Sequence.Set(event.Sequence())
	switch event.Type() {
	case user_repo.UserTokenAddedType,
		user_repo.PersonalAccessTokenAddedType:
		t.setRootData(event)
		err := t.setData(event)
		if err != nil {
			return err
		}
		t.CreationDate.Set(event.CreatedAt())
		t.IsPAT.Set(event.Type() == user_repo.PersonalAccessTokenAddedType)
	}
	return nil
}

func (t *TokenView) setRootData(event eventstore.Event) {
	t.UserID.Set(event.Aggregate().ID)
	t.ResourceOwner.Set(event.Aggregate().ResourceOwner)
	t.InstanceID = event.Aggregate().InstanceID
}

func (t *TokenView) setData(event eventstore.Event) error {
	if err := event.Unmarshal(t); err != nil {
		logging.Log("EVEN-3Gm9s").WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(err, "MODEL-5Gms9", "could not unmarshal event")
	}
	return nil
}

func agentIDFromSession(event eventstore.Event) (string, error) {
	session := make(map[string]interface{})
	if err := event.Unmarshal(&session); err != nil {
		logging.Log("EVEN-Ghgt3").WithError(err).Error("could not unmarshal event data")
		return "", zerrors.ThrowInternal(nil, "MODEL-GBf32", "could not unmarshal data")
	}
	return session["userAgentID"].(string), nil
}

func (t *TokenView) appendTokenRemoved(event eventstore.Event) error {
	token, err := eventToMap(event)
	if err != nil {
		return err
	}
	if token["tokenId"] == t.ID {
		t.Deactivated = true
	}
	return nil
}

func (t *TokenView) appendRefreshTokenRemoved(event eventstore.Event) error {
	refreshToken, err := eventToMap(event)
	if err != nil {
		return err
	}
	if refreshToken["tokenId"] == t.RefreshTokenID {
		t.Deactivated = true
	}
	return nil
}

func (t *TokenView) appendPATRemoved(event eventstore.Event) error {
	pat, err := eventToMap(event)
	if err != nil {
		return err
	}
	if pat["tokenId"] == t.ID && t.IsPAT.Value() {
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

func eventToMap(event eventstore.Event) (map[string]interface{}, error) {
	m := make(map[string]interface{})
	if err := event.Unmarshal(&m); err != nil {
		logging.Log("EVEN-Dbffe").WithError(err).Error("could not unmarshal event data")
		return nil, zerrors.ThrowInternal(nil, "MODEL-SDAfw", "could not unmarshal data")
	}
	return m, nil
}
