package domain

import (
	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/eventstore"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/crypto"
)

var (
	pool                        database.Pool
	legacyEventstore            eventstore.LegacyEventstore
	sysConfig                   systemdefaults.SystemDefaults
	passwordHasher              *crypto.Hasher
	idpEncryptionAlgo           crypto.EncryptionAlgorithm
	sessionTokenDecryptor       SessionTokenDecryptor
	mfaEncryptionAlgo           crypto.EncryptionAlgorithm
	otpSMSSecretGeneratorConfig *crypto.GeneratorConfig
)

func SetPool(p database.Pool) {
	pool = p
}

func SetLegacyEventstore(es eventstore.LegacyEventstore) {
	legacyEventstore = es
}

func SetSystemConfig(cfg systemdefaults.SystemDefaults) {
	sysConfig = cfg
}

func SetPasswordHasher(hasher *crypto.Hasher) {
	passwordHasher = hasher
}

func SetIDPEncryptionAlgorithm(idpEncryptionAlg crypto.EncryptionAlgorithm) {
	idpEncryptionAlgo = idpEncryptionAlg
}

func SetSessionTokenDecryptor(decryptor SessionTokenDecryptor) {
	sessionTokenDecryptor = decryptor
}

func SetOTPSMSSecretGeneratorConfig(cfg *crypto.GeneratorConfig) {
	otpSMSSecretGeneratorConfig = cfg
}

func SetMFAEncryptionAlgorithm(mfaEncryptionAlg crypto.EncryptionAlgorithm) {
	mfaEncryptionAlgo = mfaEncryptionAlg
}
