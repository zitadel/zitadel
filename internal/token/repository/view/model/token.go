package model

import (
	"time"

	"github.com/lib/pq"

	"github.com/caos/zitadel/internal/token/model"
)

const (
	TokenKeyTokenID       = "id"
	TokenKeyUserID        = "user_id"
	TokenKeyApplicationID = "application_id"
	TokenKeyUserAgentID   = "user_agent_id"
	TokenKeyExpiration    = "expiration"
	TokenKeyResourceOwner = "resource_owner"
)

type Token struct {
	ID            string         `json:"-" gorm:"column:id;primary_key"`
	CreationDate  time.Time      `json:"-" gorm:"column:creation_date"`
	ChangeDate    time.Time      `json:"-" gorm:"column:change_date"`
	ResourceOwner string         `json:"-" gorm:"column:resource_owner"`
	UserID        string         `json:"-" gorm:"column:user_id"`
	ApplicationID string         `json:"-" gorm:"column:application_id"`
	UserAgentID   string         `json:"-" gorm:"column:user_agent_id"`
	Scopes        pq.StringArray `json:"-" gorm:"column:scopes"`
	Expiration    time.Time      `json:"-" gorm:"column:expiration"`
	Sequence      uint64         `json:"-" gorm:"column:sequence"`
}

func TokenFromModel(token *model.Token) *Token {
	return &Token{
		ID:            token.ID,
		CreationDate:  token.CreationDate,
		ChangeDate:    token.ChangeDate,
		ResourceOwner: token.ResourceOwner,
		UserID:        token.UserID,
		ApplicationID: token.ApplicationID,
		UserAgentID:   token.UserAgentID,
		Scopes:        token.Scopes,
		Expiration:    token.Expiration,
		Sequence:      token.Sequence,
	}
}

func TokenToModel(token *Token) *model.Token {
	return &model.Token{
		ID:            token.ID,
		CreationDate:  token.CreationDate,
		ChangeDate:    token.ChangeDate,
		ResourceOwner: token.ResourceOwner,
		UserID:        token.UserID,
		ApplicationID: token.ApplicationID,
		UserAgentID:   token.UserAgentID,
		Scopes:        token.Scopes,
		Expiration:    token.Expiration,
		Sequence:      token.Sequence,
	}
}
