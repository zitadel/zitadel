package model

import (
	"time"

	"github.com/zitadel/logging"

	"github.com/zitadel/zitadel/internal/database"
	"github.com/zitadel/zitadel/internal/eventstore"
	"github.com/zitadel/zitadel/internal/eventstore/handler/v2"
	user_repo "github.com/zitadel/zitadel/internal/repository/user"
	usr_model "github.com/zitadel/zitadel/internal/user/model"
	"github.com/zitadel/zitadel/internal/view/repository"
	"github.com/zitadel/zitadel/internal/zerrors"
)

const (
	RefreshTokenKeyTokenID       = "id"
	RefreshTokenKeyUserID        = "user_id"
	RefreshTokenKeyApplicationID = "application_id"
	RefreshTokenKeyUserAgentID   = "user_agent_id"
	RefreshTokenKeyExpiration    = "expiration"
	RefreshTokenKeyResourceOwner = "resource_owner"
	RefreshTokenKeyInstanceID    = "instance_id"
)

type RefreshTokenView struct {
	ID         string `json:"tokenId" gorm:"column:id;primary_key"`
	InstanceID string `json:"instanceID" gorm:"column:instance_id;primary_key"`

	CreationDate          repository.Field[time.Time]                  `json:"-" gorm:"column:creation_date"`
	ChangeDate            repository.Field[time.Time]                  `json:"-" gorm:"column:change_date"`
	ResourceOwner         repository.Field[string]                     `json:"-" gorm:"column:resource_owner"`
	Token                 repository.Field[string]                     `json:"-" gorm:"column:token"`
	UserID                repository.Field[string]                     `json:"-" gorm:"column:user_id"`
	ClientID              repository.Field[string]                     `json:"clientID" gorm:"column:client_id"`
	UserAgentID           repository.Field[string]                     `json:"userAgentId" gorm:"column:user_agent_id"`
	Audience              repository.Field[database.TextArray[string]] `json:"audience" gorm:"column:audience"`
	Scopes                repository.Field[database.TextArray[string]] `json:"scopes" gorm:"column:scopes"`
	AuthMethodsReferences repository.Field[database.TextArray[string]] `json:"authMethodsReference" gorm:"column:amr"`
	AuthTime              repository.Field[time.Time]                  `json:"authTime" gorm:"column:auth_time"`
	IdleExpiration        repository.Field[time.Time]                  `json:"-" gorm:"column:idle_expiration"`
	Expiration            repository.Field[time.Time]                  `json:"-" gorm:"column:expiration"`
	Sequence              repository.Field[uint64]                     `json:"-" gorm:"column:sequence"`
	Actor                 repository.Field[TokenActor]                 `json:"actor" gorm:"column:actor"`
}

func (v *RefreshTokenView) PKColumns() []handler.Column {
	return []handler.Column{
		handler.NewCol(RefreshTokenKeyTokenID, v.ID),
		handler.NewCol(RefreshTokenKeyInstanceID, v.InstanceID),
	}
}

func (v *RefreshTokenView) PKConditions() []handler.Condition {
	return []handler.Condition{
		handler.NewCond(RefreshTokenKeyTokenID, v.ID),
		handler.NewCond(RefreshTokenKeyInstanceID, v.InstanceID),
	}
}

func (v *RefreshTokenView) Changes() []handler.Column {
	changes := make([]handler.Column, 0, 14)

	if v.CreationDate.DidChange() {
		changes = append(changes, handler.NewCol("creation_date", v.CreationDate.Value()))
	}
	if v.ChangeDate.DidChange() {
		changes = append(changes, handler.NewCol("change_date", v.ChangeDate.Value()))
	}
	if v.ResourceOwner.DidChange() {
		changes = append(changes, handler.NewCol("resource_owner", v.ResourceOwner.Value()))
	}
	if v.Token.DidChange() {
		changes = append(changes, handler.NewCol("token", v.Token.Value()))
	}
	if v.UserID.DidChange() {
		changes = append(changes, handler.NewCol("user_id", v.UserID.Value()))
	}
	if v.ClientID.DidChange() {
		changes = append(changes, handler.NewCol("client_id", v.ClientID.Value()))
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
	if v.AuthMethodsReferences.DidChange() {
		changes = append(changes, handler.NewCol("amr", v.AuthMethodsReferences.Value()))
	}
	if v.AuthTime.DidChange() {
		changes = append(changes, handler.NewCol("auth_time", v.AuthTime.Value()))
	}
	if v.IdleExpiration.DidChange() {
		changes = append(changes, handler.NewCol("idle_expiration", v.IdleExpiration.Value()))
	}
	if v.Expiration.DidChange() {
		changes = append(changes, handler.NewCol("expiration", v.Expiration.Value()))
	}
	if v.Sequence.DidChange() {
		changes = append(changes, handler.NewCol("sequence", v.Sequence.Value()))
	}
	if v.Actor.DidChange() {
		changes = append(changes, handler.NewCol("actor", v.Actor.Value()))
	}

	return changes
}

func RefreshTokenViewsToModel(tokens []*RefreshTokenView) []*usr_model.RefreshTokenView {
	result := make([]*usr_model.RefreshTokenView, len(tokens))
	for i, g := range tokens {
		result[i] = RefreshTokenViewToModel(g)
	}
	return result
}

func RefreshTokenViewToModel(token *RefreshTokenView) *usr_model.RefreshTokenView {
	return &usr_model.RefreshTokenView{
		ID:                    token.ID,
		CreationDate:          token.CreationDate.Value(),
		ChangeDate:            token.ChangeDate.Value(),
		ResourceOwner:         token.ResourceOwner.Value(),
		Token:                 token.Token.Value(),
		UserID:                token.UserID.Value(),
		ClientID:              token.ClientID.Value(),
		UserAgentID:           token.UserAgentID.Value(),
		Audience:              token.Audience.Value(),
		Scopes:                token.Scopes.Value(),
		AuthMethodsReferences: token.AuthMethodsReferences.Value(),
		AuthTime:              token.AuthTime.Value(),
		IdleExpiration:        token.IdleExpiration.Value(),
		Expiration:            token.Expiration.Value(),
		Sequence:              token.Sequence.Value(),
		Actor:                 token.Actor.Value().TokenActor,
	}
}

func (t *RefreshTokenView) AppendEventIfMyRefreshToken(event eventstore.Event) (err error) {
	view := new(RefreshTokenView)
	switch event.Type() {
	case user_repo.HumanRefreshTokenAddedType:
		view.setRootData(event)
		err = view.appendAddedEvent(event)
		if err != nil {
			return err
		}
	case user_repo.HumanRefreshTokenRenewedType:
		view.setRootData(event)
		err = view.appendRenewedEvent(event)
		if err != nil {
			return err
		}
	case user_repo.HumanRefreshTokenRemovedType,
		user_repo.UserRemovedType,
		user_repo.UserDeactivatedType,
		user_repo.UserLockedType:
		view.appendRemovedEvent(event)
	default:
		return nil
	}
	if view.ID == t.ID {
		return t.AppendEvent(event)
	}
	return nil
}

func (t *RefreshTokenView) AppendEvent(event eventstore.Event) error {
	t.ChangeDate.Set(event.CreatedAt())
	t.Sequence.Set(event.Sequence())
	switch event.Type() {
	case user_repo.HumanRefreshTokenAddedType:
		t.setRootData(event)
		return t.appendAddedEvent(event)
	case user_repo.HumanRefreshTokenRenewedType:
		t.setRootData(event)
		return t.appendRenewedEvent(event)
	}
	return nil
}

func (t *RefreshTokenView) setRootData(event eventstore.Event) {
	t.UserID.Set(event.Aggregate().ID)
	t.ResourceOwner.Set(event.Aggregate().ResourceOwner)
	t.InstanceID = event.Aggregate().InstanceID
}

func (t *RefreshTokenView) appendAddedEvent(event eventstore.Event) error {
	e := new(user_repo.HumanRefreshTokenAddedEvent)
	if err := event.Unmarshal(e); err != nil {
		logging.Log("EVEN-Dbb31").WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(err, "MODEL-Bbr42", "could not unmarshal event")
	}
	t.ID = e.TokenID
	t.CreationDate.Set(event.CreatedAt())
	t.AuthMethodsReferences.Set(e.AuthMethodsReferences)
	t.AuthTime.Set(e.AuthTime)
	t.Audience.Set(e.Audience)
	t.ClientID.Set(e.ClientID)
	t.Expiration.Set(event.CreatedAt().Add(e.Expiration))
	t.IdleExpiration.Set(event.CreatedAt().Add(e.IdleExpiration))
	t.Scopes.Set(e.Scopes)
	t.Token.Set(e.TokenID)
	t.UserAgentID.Set(e.UserAgentID)
	t.Actor.Set(TokenActor{e.Actor})
	return nil
}

func (t *RefreshTokenView) appendRenewedEvent(event eventstore.Event) error {
	e := new(user_repo.HumanRefreshTokenRenewedEvent)
	if err := event.Unmarshal(e); err != nil {
		logging.Log("EVEN-Vbbn2").WithError(err).Error("could not unmarshal event data")
		return zerrors.ThrowInternal(err, "MODEL-Bbrn4", "could not unmarshal event")
	}
	t.ID = e.TokenID
	t.IdleExpiration.Set(event.CreatedAt().Add(e.IdleExpiration))
	t.Token.Set(e.RefreshToken)
	return nil
}

func (t *RefreshTokenView) appendRemovedEvent(event eventstore.Event) {
	t.Expiration.Set(event.CreatedAt())
}

func (t *RefreshTokenView) GetRelevantEventTypes() []eventstore.EventType {
	return []eventstore.EventType{
		user_repo.HumanRefreshTokenAddedType,
		user_repo.HumanRefreshTokenRenewedType,
		user_repo.HumanRefreshTokenRemovedType,
		user_repo.UserRemovedType,
		user_repo.UserDeactivatedType,
		user_repo.UserLockedType,
	}
}
