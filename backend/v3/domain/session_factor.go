package domain

import (
	"time"

	"github.com/zitadel/zitadel/backend/v3/storage/database"
	"github.com/zitadel/zitadel/internal/crypto"
	"github.com/zitadel/zitadel/internal/domain"
)

type SessionFactorType interface {
	isSessionFactorType()
}

type SessionFactorTypeUser struct {
	UserID         string
	LastVerifiedAt time.Time
}

func (s *SessionFactorTypeUser) isSessionFactorType() {}

type SessionFactorTypePassword struct {
	LastVerifiedAt time.Time
}

func (s SessionFactorTypePassword) isSessionFactorType() {}

type SessionFactorTypeIDPIntent struct {
	LastVerifiedAt time.Time
}

func (s SessionFactorTypeIDPIntent) isSessionFactorType() {}

type SessionFactorTypePasskey struct {
	LastVerifiedAt time.Time
	UserVerified   bool
}

func (s SessionFactorTypePasskey) isSessionFactorType() {}

type SessionFactorTypeTOTP struct {
	LastVerifiedAt time.Time
}

func (s SessionFactorTypeTOTP) isSessionFactorType() {}

type SessionFactorTypeOTPSMS struct {
	LastVerifiedAt time.Time
}

func (s SessionFactorTypeOTPSMS) isSessionFactorType() {}

type SessionFactorTypeOTPEmail struct {
	LastVerifiedAt time.Time
}

func (s SessionFactorTypeOTPEmail) isSessionFactorType() {}

type SessionFactors []SessionFactorType

func (s SessionFactors) GetUserFactor() *SessionFactorTypeUser {
	factor, _ := GetFactorType[*SessionFactorTypeUser](s)
	return factor
}

func (s SessionFactors) GetPasswordFactor() *SessionFactorTypePassword {
	factor, _ := GetFactorType[*SessionFactorTypePassword](s)
	return factor
}

func (s SessionFactors) GetIDPIntentFactor() *SessionFactorTypeIDPIntent {
	factor, _ := GetFactorType[*SessionFactorTypeIDPIntent](s)
	return factor
}

func (s SessionFactors) GetPasskeyFactor() *SessionFactorTypePasskey {
	factor, _ := GetFactorType[*SessionFactorTypePasskey](s)
	return factor
}

func (s SessionFactors) GetTOTPFactor() *SessionFactorTypeTOTP {
	factor, _ := GetFactorType[*SessionFactorTypeTOTP](s)
	return factor
}

func (s SessionFactors) GetOTPSMSFactor() *SessionFactorTypeOTPSMS {
	factor, _ := GetFactorType[*SessionFactorTypeOTPSMS](s)
	return factor
}

func (s SessionFactors) GetOTPEmailFactor() *SessionFactorTypeOTPEmail {
	factor, _ := GetFactorType[*SessionFactorTypeOTPEmail](s)
	return factor
}

func GetFactorType[T SessionFactorType](s SessionFactors) (T, bool) {
	var nilT T
	for _, factorType := range s {
		if ft, ok := factorType.(T); ok {
			return ft, true
		}
	}
	return nilT, false
}

type SessionChallengeType interface {
	isSessionChallengeType()
}

type SessionChallengeTypePasskey struct {
	LastChallengedAt     time.Time
	Challenge            string
	AllowedCredentialIDs [][]byte
	UserVerification     domain.UserVerificationRequirement
	RPID                 string
}

func (s *SessionChallengeTypePasskey) isSessionChallengeType() {}

type SessionChallengeTypeOTPSMS struct {
	LastChallengedAt  time.Time
	Code              *crypto.CryptoValue
	Expiry            time.Duration
	CodeReturned      bool
	GeneratorID       string
	TriggeredAtOrigin string
}

func (s *SessionChallengeTypeOTPSMS) isSessionChallengeType() {}

type SessionChallengeTypeOTPEmail struct {
	LastChallengedAt  time.Time
	Code              *crypto.CryptoValue
	Expiry            time.Duration
	CodeReturned      bool
	URLTmpl           string
	TriggeredAtOrigin string
}

func (s *SessionChallengeTypeOTPEmail) isSessionChallengeType() {}

/*
type SessionFactorType int

const (

	SessionFactorTypeUnknown SessionFactorType = iota
	SessionFactorTypeUser
	SessionFactorTypePassword
	SessionFactorTypePasskey
	SessionFactorTypeIDPIntent
	SessionFactorTypeTOTP
	SessionFactorTypeOTPSMS
	SessionFactorTypeOTPEmail

)
*/
//type SessionFactorRepository interface {
//	sessionFactorColumns
//	sessionFactorConditions
//}

type sessionFactorColumns interface {
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

type sessionFactorConditions interface {
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
