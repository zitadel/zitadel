package domain

import (
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
)

//go:generate enumer -type=SessionFactorType -trimprefix SessionFactorType -transform snake -linecomment -sql
type SessionFactorType uint8

const (
	SessionFactorTypeUnknown SessionFactorType = iota
	SessionFactorTypeUser
	SessionFactorTypePassword
	// SessionFactorTypePasskey was formerly known as WebAuthn
	SessionFactorTypePasskey
	SessionFactorTypeIdentityProviderIntent
	// SessionFactorTypeTOTP used to be called OTP in some places
	SessionFactorTypeTOTP
	SessionFactorTypeOTPSMS // otp_sms
	SessionFactorTypeOTPEmail
	SessionFactorTypeRecoveryCode
)

type SessionFactor interface {
	SessionFactorType() SessionFactorType
}

type SessionFactorUser struct {
	UserID         string    `json:"userId,omitempty" db:"user_id"`
	LastVerifiedAt time.Time `json:"lastVerifiedAt,omitzero" db:"last_verified_at"`
}

// SessionFactorType implements [SessionFactor].
func (s *SessionFactorUser) SessionFactorType() SessionFactorType {
	return SessionFactorTypeUser
}

type SessionFactorPassword struct {
	LastVerifiedAt time.Time `json:"lastVerifiedAt,omitzero" db:"last_verified_at"`
	LastFailedAt   time.Time `json:"lastFailedAt,omitzero" db:"last_failed_at"`
}

// SessionFactorType implements [SessionFactor].
func (s *SessionFactorPassword) SessionFactorType() SessionFactorType {
	return SessionFactorTypePassword
}

type SessionFactorIdentityProviderIntent struct {
	LastVerifiedAt time.Time `json:"lastVerifiedAt,omitzero" db:"last_verified_at"`
}

// SessionFactorType implements [SessionFactor].
func (s *SessionFactorIdentityProviderIntent) SessionFactorType() SessionFactorType {
	return SessionFactorTypeIdentityProviderIntent
}

type SessionFactorPasskey struct {
	LastVerifiedAt time.Time `json:"lastVerifiedAt,omitzero" db:"last_verified_at"`
	UserVerified   bool      `json:"userVerified,omitempty" db:"user_verified"`
}

// SessionFactorType implements [SessionFactor].
func (s *SessionFactorPasskey) SessionFactorType() SessionFactorType {
	return SessionFactorTypePasskey
}

type SessionFactorTOTP struct {
	LastVerifiedAt time.Time `json:"lastVerifiedAt,omitzero" db:"last_verified_at"`
	LastFailedAt   time.Time `json:"lastFailedAt,omitzero" db:"last_failed_at"`
}

// SessionFactorType implements [SessionFactor].
func (s *SessionFactorTOTP) SessionFactorType() SessionFactorType {
	return SessionFactorTypeTOTP
}

type SessionFactorOTPSMS struct {
	LastVerifiedAt time.Time `json:"lastVerifiedAt,omitzero" db:"last_verified_at"`
	LastFailedAt   time.Time `json:"lastFailedAt,omitzero" db:"last_failed_at"`
}

// SessionFactorType implements [SessionFactor].
func (s *SessionFactorOTPSMS) SessionFactorType() SessionFactorType {
	return SessionFactorTypeOTPSMS
}

type SessionFactorOTPEmail struct {
	LastVerifiedAt time.Time `json:"lastVerifiedAt,omitzero" db:"last_verified_at"`
	LastFailedAt   time.Time `json:"lastFailedAt,omitzero" db:"last_failed_at"`
}

// SessionFactorType implements [SessionFactor].
func (s *SessionFactorOTPEmail) SessionFactorType() SessionFactorType {
	return SessionFactorTypeOTPEmail
}

type SessionFactorRecoveryCode struct {
	LastVerifiedAt time.Time `json:"lastVerifiedAt,omitzero" db:"last_verified_at"`
	LastFailedAt   time.Time `json:"lastFailedAt,omitzero" db:"last_failed_at"`
}

// SessionFactorType implements [SessionFactor].
func (s *SessionFactorRecoveryCode) SessionFactorType() SessionFactorType {
	return SessionFactorTypeRecoveryCode
}

type SessionFactors []SessionFactor

func (s *SessionFactors) AppendTo(factor SessionFactor) {
	*s = append(*s, factor)
}

// GetUserFactor returns the user factor from the factors.
func (s SessionFactors) GetUserFactor() *SessionFactorUser {
	factor, _ := GetFactor[*SessionFactorUser](s)
	return factor
}

// GetPasswordFactor returns the password factor from the factors.
func (s SessionFactors) GetPasswordFactor() *SessionFactorPassword {
	factor, _ := GetFactor[*SessionFactorPassword](s)
	return factor
}

// GetIDPIntentFactor returns the IDP Intent factor from the factors.
func (s SessionFactors) GetIDPIntentFactor() *SessionFactorIdentityProviderIntent {
	factor, _ := GetFactor[*SessionFactorIdentityProviderIntent](s)
	return factor
}

// GetPasskeyFactor returns the passkey factor from the factors.
func (s SessionFactors) GetPasskeyFactor() *SessionFactorPasskey {
	factor, _ := GetFactor[*SessionFactorPasskey](s)
	return factor
}

// GetTOTPFactor returns the TOTP factor from the factors.
func (s SessionFactors) GetTOTPFactor() *SessionFactorTOTP {
	factor, _ := GetFactor[*SessionFactorTOTP](s)
	return factor
}

// GetOTPSMSFactor returns the OTP SMS factor from the factors.
func (s SessionFactors) GetOTPSMSFactor() *SessionFactorOTPSMS {
	factor, _ := GetFactor[*SessionFactorOTPSMS](s)
	return factor
}

// GetOTPEmailFactor returns the OTP Email factor from the factors.
func (s SessionFactors) GetOTPEmailFactor() *SessionFactorOTPEmail {
	factor, _ := GetFactor[*SessionFactorOTPEmail](s)
	return factor
}

// GetRecoveryCodeFactor returns the recovery code factor from the factors.
func (s SessionFactors) GetRecoveryCodeFactor() *SessionFactorRecoveryCode {
	factor, _ := GetFactor[*SessionFactorRecoveryCode](s)
	return factor
}

// GetFactor returns the factor of the given type from the factors.
func GetFactor[T SessionFactor](factors SessionFactors) (T, bool) {
	var nilT T
	for _, factor := range factors {
		if f, ok := factor.(T); ok {
			return f, true
		}
	}
	return nilT, false
}

type SessionChallenge interface {
	SessionChallengeType() SessionFactorType
}

type SessionChallengePasskey struct {
	LastChallengedAt     time.Time                          `json:"lastChallengedAt,omitzero" db:"last_challenged_at"`
	Challenge            string                             `json:"challenge,omitempty" db:"challenge"`
	AllowedCredentialIDs [][]byte                           `json:"allowedCredentialIDs,omitempty" db:"allowed_credential_ids"`
	UserVerification     domain.UserVerificationRequirement `json:"userVerification,omitempty" db:"user_verification"`
	RPID                 string                             `json:"rpId,omitempty" db:"rp_id"`
}

// SessionChallengeType implements [SessionChallenge].
func (s *SessionChallengePasskey) SessionChallengeType() SessionFactorType {
	return SessionFactorTypePasskey
}

type SessionChallengeOTPSMS struct {
	LastChallengedAt  time.Time           `json:"lastChallengedAt,omitzero" db:"last_challenged_at"`
	Code              *crypto.CryptoValue `json:"code,omitempty" db:"code"`
	Expiry            time.Duration       `json:"expiry,omitzero" db:"expiry"`
	CodeReturned      bool                `json:"codeReturned,omitempty" db:"code_returned"`
	GeneratorID       string              `json:"generatorId,omitempty" db:"generator_id"`
	TriggeredAtOrigin string              `json:"triggeredAtOrigin,omitempty" db:"triggered_at_origin"`
}

// SessionChallengeType implements [SessionChallenge].
func (s *SessionChallengeOTPSMS) SessionChallengeType() SessionFactorType {
	return SessionFactorTypeOTPSMS
}

type SessionChallengeOTPEmail struct {
	LastChallengedAt  time.Time           `json:"lastChallengedAt,omitzero" db:"last_challenged_at"`
	Code              *crypto.CryptoValue `json:"code,omitempty" db:"code"`
	Expiry            time.Duration       `json:"expiry,omitzero" db:"expiry"`
	CodeReturned      bool                `json:"codeReturned,omitempty" db:"code_returned"`
	URLTmpl           string              `json:"urlTmpl,omitempty" db:"url_template"`
	TriggeredAtOrigin string              `json:"triggeredAtOrigin,omitempty" db:"triggered_at_origin"`
}

// SessionChallengeType implements [SessionChallenge].
func (s *SessionChallengeOTPEmail) SessionChallengeType() SessionFactorType {
	return SessionFactorTypeOTPEmail
}

type SessionChallenges []SessionChallenge

func (s *SessionChallenges) AppendTo(challenge SessionChallenge) {
	*s = append(*s, challenge)
}

// GetPasskeyChallenge returns the passkey challenge from the challenges.
func (s SessionChallenges) GetPasskeyChallenge() *SessionChallengePasskey {
	challenge, _ := GetChallenge[*SessionChallengePasskey](s)
	return challenge
}

// GetOTPSMSChallenge returns the OTP SMS challenge from the challenges.
func (s SessionChallenges) GetOTPSMSChallenge() *SessionChallengeOTPSMS {
	challenge, _ := GetChallenge[*SessionChallengeOTPSMS](s)
	return challenge
}

// GetOTPEmailChallenge returns the OTP Email challenge from the challenges.
func (s SessionChallenges) GetOTPEmailChallenge() *SessionChallengeOTPEmail {
	challenge, _ := GetChallenge[*SessionChallengeOTPEmail](s)
	return challenge
}

// GetChallenge returns the challenge of the given type from the challenges.
func GetChallenge[T SessionChallenge](challenges SessionChallenges) (T, bool) {
	var nilT T
	for _, challenge := range challenges {
		if c, ok := challenge.(T); ok {
			return c, true
		}
	}
	return nilT, false
}

type SessionFactorConditions interface {
	// FactorTypeCondition returns an equal filter on the factor type field.
	FactorTypeCondition(factorType SessionFactorType) database.Condition
	// LastVerifiedBeforeCondition returns a filter on the factor last verified field before the given time.
	LastVerifiedBeforeCondition(lastVerifiedAt time.Time) database.Condition
}
