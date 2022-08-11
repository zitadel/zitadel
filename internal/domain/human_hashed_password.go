package domain

import (
	"github.com/zitadel/zitadel/internal/crypto"
	es_models "github.com/zitadel/zitadel/internal/eventstore/v1/models"
)

type HashedPassword struct {
	es_models.ObjectRoot

	SecretString string
	SecretCrypto *crypto.CryptoValue
}

func NewHashedPassword(password, algorithm string) *HashedPassword {
	return &HashedPassword{
		SecretString: password,
		SecretCrypto: &crypto.CryptoValue{
			CryptoType: crypto.TypeHash,
			Algorithm:  algorithm,
			Crypted:    []byte(password),
		},
	}
}
