package model

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	"time"
)

type Email struct {
	es_models.ObjectRoot

	EmailAddress    string `json:"email,omitempty"`
	IsEmailVerified bool   `json:"-"`

	isEmailUnique bool `json:"-"`
}

type EmailCode struct {
	es_models.ObjectRoot

	Code   *crypto.CryptoValue `json:"code,omitempty"`
	Expiry time.Duration       `json:"expiry,omitempty"`
}

func (e *Email) Changes(changed *Email) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.EmailAddress != "" && e.EmailAddress != changed.EmailAddress {
		changes["email"] = changed.EmailAddress
	}
	return changes
}

func EmailFromModel(email *model.Email) *Email {
	return &Email{
		ObjectRoot: es_models.ObjectRoot{
			ID:           email.ObjectRoot.ID,
			Sequence:     email.Sequence,
			ChangeDate:   email.ChangeDate,
			CreationDate: email.CreationDate,
		},
		EmailAddress:    email.EmailAddress,
		IsEmailVerified: email.IsEmailVerified,
	}
}

func EmailToModel(email *Email) *model.Email {
	return &model.Email{
		ObjectRoot: es_models.ObjectRoot{
			ID:           email.ObjectRoot.ID,
			Sequence:     email.Sequence,
			ChangeDate:   email.ChangeDate,
			CreationDate: email.CreationDate,
		},
		EmailAddress:    email.EmailAddress,
		IsEmailVerified: email.IsEmailVerified,
	}
}
