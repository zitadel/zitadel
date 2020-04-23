package model

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
	"time"
)

type Phone struct {
	es_models.ObjectRoot

	PhoneNumber     string `json:"phone,omitempty"`
	IsPhoneVerified bool   `json:"-"`
}

type PhoneCode struct {
	es_models.ObjectRoot

	Code   *crypto.CryptoValue `json:"code,omitempty"`
	Expiry time.Duration       `json:"expiry,omitempty"`
}

func (p *Phone) Changes(changed *Phone) map[string]interface{} {
	changes := make(map[string]interface{}, 1)
	if changed.PhoneNumber != "" && p.PhoneNumber != changed.PhoneNumber {
		changes["phone"] = changed.PhoneNumber
	}
	return changes
}

func PhoneFromModel(phone *model.Phone) *Phone {
	return &Phone{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  phone.ObjectRoot.AggregateID,
			Sequence:     phone.Sequence,
			ChangeDate:   phone.ChangeDate,
			CreationDate: phone.CreationDate,
		},
		PhoneNumber:     phone.PhoneNumber,
		IsPhoneVerified: phone.IsPhoneVerified,
	}
}

func PhoneToModel(phone *Phone) *model.Phone {
	return &model.Phone{
		ObjectRoot: es_models.ObjectRoot{
			AggregateID:  phone.ObjectRoot.AggregateID,
			Sequence:     phone.Sequence,
			ChangeDate:   phone.ChangeDate,
			CreationDate: phone.CreationDate,
		},
		PhoneNumber:     phone.PhoneNumber,
		IsPhoneVerified: phone.IsPhoneVerified,
	}
}
