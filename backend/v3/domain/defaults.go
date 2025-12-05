package domain

import (
	"log/slog"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/backend/v3/storage/eventstore"
	"github.com/zitadel/zitadel/backend/v3/telemetry/logging"
	"github.com/zitadel/zitadel/backend/v3/telemetry/tracing"
	"github.com/zitadel/zitadel/internal/config/systemdefaults"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/webauthn"
)

var (
	pool              database.Pool
	tracer            tracing.Tracer
	logger            logging.Logger = *logging.NewLogger(slog.Default())
	legacyEventstore  eventstore.LegacyEventstore
	sysConfig         systemdefaults.SystemDefaults
	passwordHasher    *crypto.Hasher
	idpEncryptionAlgo crypto.EncryptionAlgorithm
	mfaEncryptionAlgo crypto.EncryptionAlgorithm
	otpEncryptionAlgo crypto.EncryptionAlgorithm

	webauthnConfig *webauthn.Config
)

func SetPool(p database.Pool) {
	pool = p
}

func SetTracer(t tracing.Tracer) {
	tracer = t
}

func SetLogger(l logging.Logger) {
	logger = l
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

func SetWebAuthNConfig(cfg *webauthn.Config) {
	webauthnConfig = cfg
}

func SetMFAEncryptionAlgorithm(mfaEncryptionAlg crypto.EncryptionAlgorithm) {
	mfaEncryptionAlgo = mfaEncryptionAlg
}

func SetOTPEncryptionAlgorithm(otpEncryptionAlg crypto.EncryptionAlgorithm) {
	otpEncryptionAlgo = otpEncryptionAlg
}
