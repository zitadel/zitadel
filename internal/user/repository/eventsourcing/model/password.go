package model

import (
	"github.com/caos/zitadel/internal/crypto"
	es_models "github.com/caos/zitadel/internal/eventstore/models"
	"github.com/caos/zitadel/internal/user/model"
)

type Password struct {
	es_models.ObjectRoot

	Secret         *crypto.CryptoValue `json:"secret,omitempty"`
	ChangeRequired bool                `json:"changeRequired,omitempty"`
}

func PasswordFromModel(password *model.Password) *Password {
	return &Password{
		ObjectRoot: es_models.ObjectRoot{
			ID:           password.ObjectRoot.ID,
			Sequence:     password.Sequence,
			ChangeDate:   password.ChangeDate,
			CreationDate: password.CreationDate,
		},
		Secret:         password.SecretCrypto,
		ChangeRequired: password.ChangeRequired,
	}
}

func PasswordToModel(password *Password) *model.Password {
	return &model.Password{
		ObjectRoot: es_models.ObjectRoot{
			ID:           password.ObjectRoot.ID,
			Sequence:     password.Sequence,
			ChangeDate:   password.ChangeDate,
			CreationDate: password.CreationDate,
		},
		SecretCrypto:   password.Secret,
		ChangeRequired: password.ChangeRequired,
	}
}
