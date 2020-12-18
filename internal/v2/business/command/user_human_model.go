package command

import (
	"github.com/caos/zitadel/internal/eventstore/v2"
	"github.com/caos/zitadel/internal/v2/business/domain"
	"golang.org/x/text/language"
)

type HumanWriteModel struct {
	eventstore.WriteModel

	UserName string

	FirstName         string
	LastName          string
	NickName          string
	DisplayName       string
	PreferredLanguage language.Tag
	Gender            domain.Gender

	Email           string
	IsEmailVerified bool

	Phone           string
	IsPhoneVerified bool

	Country       string
	Locality      string
	PostalCode    string
	Region        string
	StreetAddress string

	UserState domain.UserState
}

func NewHumanWriteModel(userID string) *HumanWriteModel {
	return &HumanWriteModel{
		WriteModel: eventstore.WriteModel{
			AggregateID: userID,
		},
	}
}
