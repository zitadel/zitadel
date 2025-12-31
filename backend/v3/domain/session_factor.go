package domain

import (
	"time"

	"github.com/golang/mock/gomock"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
)

type SessionFactorType int

const (
	SessionFactorTypeUnknown SessionFactorType = iota
	SessionFactorTypeUser
	SessionFactorTypePassword
	SessionFactorTypePasskey // formerly known as WebAuthn
	SessionFactorTypeIdentityProviderIntent
	SessionFactorTypeTOTP
	SessionFactorTypeOTPSMS
	SessionFactorTypeOTPEmail
)

// TODO(IAM-Marco): Add gomock.Matcher interface and implement it on all factors
type SessionFactor interface {
	sessionFactorType() SessionFactorType
}

type SessionFactorUser struct {
	UserID         string
	LastVerifiedAt time.Time
}

func (s *SessionFactorUser) Matches(in any) bool {
	inputFactor, isFactorUser := in.(*SessionFactorUser)
	if !isFactorUser {
		return false
	}

	isUserMatching := inputFactor.UserID == s.UserID

	// Using a fixed delta, can be changed later to make it configurable if needed
	timeDelta := 800 * time.Millisecond
	// Due to LastVerifiedAt not being easily controllable through the tested code,
	// we check that it's in a certain range rather than doing an exact match on time.
	lowerThreshold := s.LastVerifiedAt.Add(-timeDelta)
	upperThreshold := s.LastVerifiedAt.Add(timeDelta)

	isLastVerifiedAtInRange := !inputFactor.LastVerifiedAt.Before(lowerThreshold) && !inputFactor.LastVerifiedAt.After(upperThreshold)

	return isUserMatching && isLastVerifiedAtInRange
}

func (s *SessionFactorUser) String() string {
	return "SessionFactorUser"
}

// sessionFactorType implements [SessionFactor].
func (s *SessionFactorUser) sessionFactorType() SessionFactorType {
	return SessionFactorTypeUser
}

type SessionFactorPassword struct {
	LastVerifiedAt time.Time
	LastFailedAt   time.Time
}

// Matches implements [gomock.Matcher].
func (s *SessionFactorPassword) Matches(in any) bool {
	inputFactor, isFactorPassword := in.(*SessionFactorPassword)
	if !isFactorPassword {
		return false
	}

	// Using a fixed delta, can be changed later to make it configurable if needed
	timeDelta := 800 * time.Millisecond
	// Due to LastVerifiedAt/LastFailedAt not being easily controllable through the tested code,
	// we check that their in a certain range rather than doing an exact match on time.
	var isLastVerifiedAtInRange, isLastFailedAtInRange bool
	if !inputFactor.LastVerifiedAt.IsZero() {
		lowerThreshold := s.LastVerifiedAt.Add(-timeDelta)
		upperThreshold := s.LastVerifiedAt.Add(timeDelta)
		isLastVerifiedAtInRange = !inputFactor.LastVerifiedAt.Before(lowerThreshold) && !inputFactor.LastVerifiedAt.After(upperThreshold)
	}

	if !inputFactor.LastFailedAt.IsZero() {
		lowerThreshold := s.LastFailedAt.Add(-timeDelta)
		upperThreshold := s.LastFailedAt.Add(timeDelta)
		isLastFailedAtInRange = !inputFactor.LastFailedAt.Before(lowerThreshold) && !inputFactor.LastFailedAt.After(upperThreshold)
	}

	return isLastVerifiedAtInRange || isLastFailedAtInRange
}

// String implements [gomock.Matcher].
func (s *SessionFactorPassword) String() string {
	return "SessionFactorPassword"
}

// sessionFactorType implements [SessionFactor].
func (s *SessionFactorPassword) sessionFactorType() SessionFactorType {
	return SessionFactorTypePassword
}

type SessionFactorIdentityProviderIntent struct {
	LastVerifiedAt time.Time
}

// sessionFactorType implements [SessionFactor].
func (s *SessionFactorIdentityProviderIntent) sessionFactorType() SessionFactorType {
	return SessionFactorTypeIdentityProviderIntent
}

type SessionFactorPasskey struct {
	LastVerifiedAt time.Time
	UserVerified   bool
}

// Matches implements [gomock.Matcher].
func (s *SessionFactorPasskey) Matches(in any) bool {
	inputFactor, isFactorPasskey := in.(*SessionFactorPasskey)
	if !isFactorPasskey {
		return false
	}

	// Using a fixed delta, can be changed later to make it configurable if needed
	timeDelta := 800 * time.Millisecond
	// Due to LastVerifiedAt not being easily controllable through the tested code,
	// we check that it's in a certain range rather than doing an exact match on time.
	var isLastVerifiedAtInRange bool
	if !inputFactor.LastVerifiedAt.IsZero() {
		lowerThreshold := s.LastVerifiedAt.Add(-timeDelta)
		upperThreshold := s.LastVerifiedAt.Add(timeDelta)
		isLastVerifiedAtInRange = !inputFactor.LastVerifiedAt.Before(lowerThreshold) && !inputFactor.LastVerifiedAt.After(upperThreshold)
	}

	return isLastVerifiedAtInRange
}

// String implements [gomock.Matcher].
func (s *SessionFactorPasskey) String() string {
	return "SessionFactorPasskey"
}

var _ gomock.Matcher = (*SessionFactorPasskey)(nil)

// sessionFactorType implements [SessionFactor].
func (s *SessionFactorPasskey) sessionFactorType() SessionFactorType {
	return SessionFactorTypePasskey
}

type SessionFactorTOTP struct {
	LastVerifiedAt time.Time
	LastFailedAt   time.Time
}

// sessionFactorType implements [SessionFactor].
func (s *SessionFactorTOTP) sessionFactorType() SessionFactorType {
	return SessionFactorTypeTOTP
}

type SessionFactorOTPSMS struct {
	LastVerifiedAt time.Time
	LastFailedAt   time.Time
}

// sessionFactorType implements [SessionFactor].
func (s *SessionFactorOTPSMS) sessionFactorType() SessionFactorType {
	return SessionFactorTypeOTPSMS
}

type SessionFactorOTPEmail struct {
	LastVerifiedAt time.Time
	LastFailedAt   time.Time
}

// sessionFactorType implements [SessionFactor].
func (s *SessionFactorOTPEmail) sessionFactorType() SessionFactorType {
	return SessionFactorTypeOTPEmail
}

type SessionFactors []SessionFactor

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
	sessionChallengeType() SessionFactorType
}

type SessionChallengePasskey struct {
	LastChallengedAt     time.Time
	Challenge            string
	AllowedCredentialIDs [][]byte
	UserVerification     domain.UserVerificationRequirement
	RPID                 string
}

// Matches implements [gomock.Matcher].
func (s *SessionChallengePasskey) Matches(in any) bool {
	challengePasskey, ok := in.(*SessionChallengePasskey)
	if !ok {
		return false
	}

	// Using a fixed delta, can be changed later to make it configurable if needed
	timeDelta := 800 * time.Millisecond
	// Due to LastChallengedAt not being easily controllable through the tested code,
	// we check that it's in a certain range rather than doing an exact match on time.
	var isLastVerifiedAtInRange bool
	if !challengePasskey.LastChallengedAt.IsZero() {
		lowerThreshold := s.LastChallengedAt.Add(-timeDelta)
		upperThreshold := s.LastChallengedAt.Add(timeDelta)
		isLastVerifiedAtInRange = !challengePasskey.LastChallengedAt.Before(lowerThreshold) && !challengePasskey.LastChallengedAt.After(upperThreshold)
	}
	return isLastVerifiedAtInRange
}

// String implements [gomock.Matcher].
func (s *SessionChallengePasskey) String() string {
	return "SessionChallengePasskey"
}

// sessionChallengeType implements [SessionChallenge].
func (s *SessionChallengePasskey) sessionChallengeType() SessionFactorType {
	return SessionFactorTypePasskey
}

type SessionChallengeOTPSMS struct {
	LastChallengedAt  time.Time
	Code              *crypto.CryptoValue
	Expiry            time.Duration
	CodeReturned      bool
	GeneratorID       string
	TriggeredAtOrigin string
}

// Matches implements [gomock.Matcher].
func (s *SessionChallengeOTPSMS) Matches(in any) bool {
	challengeOTPSMS, ok := in.(*SessionChallengeOTPSMS)
	if !ok {
		return false
	}

	// Using a fixed delta, can be changed later to make it configurable if needed
	timeDelta := 800 * time.Millisecond
	// Due to LastChallengedAt not being easily controllable through the tested code,
	// we check that it's in a certain range rather than doing an exact match on time.
	var isLastVerifiedAtInRange bool
	if !challengeOTPSMS.LastChallengedAt.IsZero() {
		lowerThreshold := s.LastChallengedAt.Add(-timeDelta)
		upperThreshold := s.LastChallengedAt.Add(timeDelta)
		isLastVerifiedAtInRange = !challengeOTPSMS.LastChallengedAt.Before(lowerThreshold) && !challengeOTPSMS.LastChallengedAt.After(upperThreshold)
	}
	return isLastVerifiedAtInRange
}

// String implements [gomock.Matcher].
func (s *SessionChallengeOTPSMS) String() string {
	return "SessionChallengeOTPSMS"
}

// sessionChallengeType implements [SessionChallenge].
func (s *SessionChallengeOTPSMS) sessionChallengeType() SessionFactorType {
	return SessionFactorTypeOTPSMS
}

type SessionChallengeOTPEmail struct {
	LastChallengedAt  time.Time
	Code              *crypto.CryptoValue
	Expiry            time.Duration
	CodeReturned      bool
	URLTmpl           string
	TriggeredAtOrigin string
}

// Matches implements [gomock.Matcher].
func (s *SessionChallengeOTPEmail) Matches(in any) bool {
	challengeOTPEmail, ok := in.(*SessionChallengeOTPEmail)
	if !ok {
		return false
	}

	// Using a fixed delta, can be changed later to make it configurable if needed
	timeDelta := 800 * time.Millisecond
	// Due to LastChallengedAt not being easily controllable through the tested code,
	// we check that it's in a certain range rather than doing an exact match on time.
	var isLastVerifiedAtInRange bool
	if !challengeOTPEmail.LastChallengedAt.IsZero() {
		lowerThreshold := s.LastChallengedAt.Add(-timeDelta)
		upperThreshold := s.LastChallengedAt.Add(timeDelta)
		isLastVerifiedAtInRange = !challengeOTPEmail.LastChallengedAt.Before(lowerThreshold) && !challengeOTPEmail.LastChallengedAt.After(upperThreshold)
	}
	return isLastVerifiedAtInRange
}

// String implements [gomock.Matcher].
func (s *SessionChallengeOTPEmail) String() string {
	return "SessionChallengeOTPEmail"
}

// sessionChallengeType implements [SessionChallenge].
func (s *SessionChallengeOTPEmail) sessionChallengeType() SessionFactorType {
	return SessionFactorTypeOTPEmail
}

type SessionChallenges []SessionChallenge

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

type SessionFactorColumns interface {
	// InstanceIDColumn returns the column for the instance id field.
	InstanceIDColumn() database.Column
	// SessionIDColumn returns the column for the session id field.
	SessionIDColumn() database.Column
	// FactorTypeColumn returns the column for the factor type field.
	FactorTypeColumn() database.Column
	// LastChallengedAtColumn returns the column for the last challenged at field.
	LastChallengedAtColumn() database.Column
	// LastFailedAtColumn returns the column for the last failed at field.
	LastFailedAtColumn() database.Column
	// LastVerifiedAtColumn returns the column for the last verified at field.
	LastVerifiedAtColumn() database.Column
}

type SessionFactorConditions interface {
	// PrimaryKeyCondition returns a filter on the primary key fields.
	PrimaryKeyCondition(instanceID, sessionID string, factorType SessionFactorType) database.Condition
	// InstanceIDCondition returns an equal filter on the instance id field.
	InstanceIDCondition(instanceID string) database.Condition
	// SessionIDCondition returns an equal filter on the session id field.
	SessionIDCondition(sessionID string) database.Condition
	// FactorTypeCondition returns an equal filter on the factor type field.
	FactorTypeCondition(factorType SessionFactorType) database.Condition
	// FactorLastVerifiedBeforeCondition returns a filter on the factor last verified field before the given time.
	FactorLastVerifiedBeforeCondition(lastVerifiedAt time.Time) database.Condition
}
