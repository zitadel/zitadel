package systemdefaults

import (
	"fmt"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
	"github.com/zitadel/zitadel/internal/zerrors"
)

type SystemDefaults struct {
	SecretGenerators     SecretGenerators
	PasswordHasher       crypto.HashConfig
	SecretHasher         crypto.HashConfig
	Multifactors         MultifactorConfig
	Tarpit               TarpitConfig
	DomainVerification   DomainVerification
	Notifications        Notifications
	KeyConfig            KeyConfig
	DefaultQueryLimit    uint64
	MaxQueryLimit        uint64
	MaxIdPIntentLifetime time.Duration
}

func (s SystemDefaults) Validate() error {
	// NB: currently we only validate the recovery codes config,
	// but we may want to add more validations in the future or refactor
	return s.Multifactors.RecoveryCodes.Validate()
}

// QueryLimit returns the limit to apply to a list query.
// If no limit is requested (0), the DefaultQueryLimit is returned.
// An error is returned if the requested limit exceeds the MaxQueryLimit.
// A DefaultQueryLimit or MaxQueryLimit of 0 disables the respective restriction.
func (s SystemDefaults) QueryLimit(requested uint64) (uint64, error) {
	if s.MaxQueryLimit > 0 && requested > s.MaxQueryLimit {
		return 0, zerrors.ThrowInvalidArgumentf(fmt.Errorf("given: %d, allowed: %d", requested, s.MaxQueryLimit), "QUERY-4M0fs", "Errors.Query.LimitExceeded")
	}
	if requested == 0 {
		return s.DefaultQueryLimit, nil
	}
	return requested, nil
}

type SecretGenerators struct {
	MachineKeySize     uint32
	ApplicationKeySize uint32
}

type MultifactorConfig struct {
	OTP           OTPConfig
	RecoveryCodes RecoveryCodesConfig
}

type OTPConfig struct {
	Issuer string
}

type RecoveryCodesConfig struct {
	MaxCount   int
	Format     string
	Length     int
	WithHyphen bool
}

func (r RecoveryCodesConfig) Validate() error {
	if r.MaxCount < 1 || r.MaxCount > 100 {
		return fmt.Errorf("RecoveryCodes.MaxCount must be between 1 and 100, got %d", r.MaxCount)
	}

	switch r.Format {
	case string(domain.RecoveryCodeFormatUUID):
		// pass
	case string(domain.RecoveryCodeFormatAlphanumeric):
		if r.Length < 8 || r.Length > 60 {
			return fmt.Errorf("RecoveryCodes.Length must be between 8 and 60 for alphanumeric format, got %d", r.Length)
		}
	default:
		return fmt.Errorf("RecoveryCodes.Format must be 'uuid' or 'alphanumeric', got '%s'", r.Format)
	}

	return nil
}

type DomainVerification struct {
	VerificationGenerator crypto.GeneratorConfig
}

type Notifications struct {
	FileSystemPath string
}

type KeyConfig struct {
	Size                int
	PrivateKeyLifetime  time.Duration
	PublicKeyLifetime   time.Duration
	CertificateSize     int
	CertificateLifetime time.Duration
}
