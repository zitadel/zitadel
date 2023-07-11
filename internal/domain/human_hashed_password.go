package domain

import (
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type HashedPassword struct {
	es_models.ObjectRoot

	SecretString  string
	EncodedSecret string
}

func NewHashedPassword(password, algorithm string) *HashedPassword {
	return &HashedPassword{
		SecretString:  password,
		EncodedSecret: password,
	}
}
