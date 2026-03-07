package domain

import (
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/eventstore"
)

var (
	pool                  database.Pool
	legacyEventstore      eventstore.LegacyEventstore
	sessionTokenVerifier  SessionTokenVerifier
	sessionTokenDecryptor SessionTokenDecryptor
)

func SetPool(p database.Pool) {
	pool = p
}

func SetLegacyEventstore(es eventstore.LegacyEventstore) {
	legacyEventstore = es
}

func SetSessionTokenVerifier(verifier SessionTokenVerifier) {
	sessionTokenVerifier = verifier
}

func SetSessionTokenDecryptor(decryptor SessionTokenDecryptor) {
	sessionTokenDecryptor = decryptor
}
