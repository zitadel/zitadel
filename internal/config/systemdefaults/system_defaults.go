package systemdefaults

import (
	"fmt"
	"strings"
	"time"

	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
)

type SystemDefaults struct {
	SecretGenerators     SecretGenerators
	PasswordHasher       crypto.HashConfig
	SecretHasher         crypto.HashConfig
	Multifactors         MultifactorConfig
	Tarpit               TarpitConfig
	Risk                 RiskConfig
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
	if err := s.Multifactors.RecoveryCodes.Validate(); err != nil {
		return err
	}
	return s.Risk.Validate()
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

type RiskConfig struct {
	Enabled               bool
	FailOpen              bool
	FailureBurstThreshold int
	HistoryWindow         time.Duration
	ContextChangeWindow   time.Duration
	MaxSignalsPerUser     int
	MaxSignalsPerSession  int
	LLM                   RiskLLMConfig
}

type RiskLLMConfig struct {
	Mode               string
	Endpoint           string
	Model              string
	Timeout            time.Duration
	MaxEvents          int
	NumPredict         int
	Temperature        *float64
	TopK               int
	TopP               float64
	KeepAlive          string
	HighRiskConfidence float64
	LogPrompts         bool
	CircuitBreaker     *RiskCBConfig
}

type RiskCBConfig struct {
	Interval               time.Duration
	MaxConsecutiveFailures uint32
	MaxFailureRatio        float64
	Timeout                time.Duration
	MaxRetryRequests       uint32
	FailOpen               bool
}

func (r RiskConfig) Validate() error {
	if !r.Enabled {
		return nil
	}
	if r.FailureBurstThreshold <= 0 {
		return fmt.Errorf("Risk.FailureBurstThreshold must be greater than 0, got %d", r.FailureBurstThreshold)
	}
	if r.HistoryWindow <= 0 {
		return fmt.Errorf("Risk.HistoryWindow must be greater than 0, got %s", r.HistoryWindow)
	}
	if r.ContextChangeWindow <= 0 {
		return fmt.Errorf("Risk.ContextChangeWindow must be greater than 0, got %s", r.ContextChangeWindow)
	}
	if r.MaxSignalsPerUser <= 0 {
		return fmt.Errorf("Risk.MaxSignalsPerUser must be greater than 0, got %d", r.MaxSignalsPerUser)
	}
	if r.MaxSignalsPerSession <= 0 {
		return fmt.Errorf("Risk.MaxSignalsPerSession must be greater than 0, got %d", r.MaxSignalsPerSession)
	}
	return r.LLM.Validate()
}

func (r RiskLLMConfig) Validate() error {
	switch strings.ToLower(r.Mode) {
	case "", "disabled":
		return nil
	case "observe", "enforce":
	default:
		return fmt.Errorf("Risk.LLM.Mode must be one of 'disabled', 'observe', or 'enforce', got %q", r.Mode)
	}
	if r.Endpoint == "" {
		return fmt.Errorf("Risk.LLM.Endpoint must not be empty when Risk.LLM.Mode is %q", r.Mode)
	}
	if r.Model == "" {
		return fmt.Errorf("Risk.LLM.Model must not be empty when Risk.LLM.Mode is %q", r.Mode)
	}
	if r.Timeout <= 0 {
		return fmt.Errorf("Risk.LLM.Timeout must be greater than 0, got %s", r.Timeout)
	}
	if r.MaxEvents <= 0 {
		return fmt.Errorf("Risk.LLM.MaxEvents must be greater than 0, got %d", r.MaxEvents)
	}
	if r.HighRiskConfidence <= 0 || r.HighRiskConfidence > 1 {
		return fmt.Errorf("Risk.LLM.HighRiskConfidence must be in (0,1], got %f", r.HighRiskConfidence)
	}
	return nil
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
