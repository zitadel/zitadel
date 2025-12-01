package repository

import (
	"encoding/json"
	"time"

	"github.com/zitadel/zitadel/backend/v3/domain"
	"github.com/zitadel/zitadel/backend/v3/storage/database"
)

type sessionFactor struct {
}

func (s sessionFactor) qualifiedTableName() string {
	return "zitadel.session_factors"
}

func (s sessionFactor) unqualifiedTableName() string {
	return "session_factors"
}

func (s sessionFactor) PrimaryKeyColumns() []database.Column {
	return []database.Column{
		s.instanceIDColumn(),
		s.sessionIDColumn(),
		s.typeColumn(),
	}
}

func (s sessionFactor) instanceIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "instance_id")
}

func (s sessionFactor) sessionIDColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "session_id")
}

func (s sessionFactor) typeColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "type")
}

func (s sessionFactor) lastChallengedAtColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "last_challenged_at")
}

func (s sessionFactor) lastFailedAtColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "last_failed_at")
}

func (s sessionFactor) lastVerifiedAtColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "last_verified_at")
}

func (s sessionFactor) challengedPayloadColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "challenged_payload")
}

func (s sessionFactor) verifiedPayloadColumn() database.Column {
	return database.NewColumn(s.unqualifiedTableName(), "verified_payload")
}

// -------------------------------------------------------------
// conditions
// -------------------------------------------------------------

func (s sessionFactor) FactorTypeCondition(factorType domain.SessionFactorType) database.Condition {
	return database.NewNumberCondition(s.typeColumn(), database.NumberOperationEqual, factorType)
}

func (s sessionFactor) LastVerifiedBeforeCondition(lastVerifiedAt time.Time) database.Condition {
	return database.NewNumberCondition(s.lastVerifiedAtColumn(), database.NumberOperationLessThan, lastVerifiedAt)
}

// -------------------------------------------------------------
// scanners
// -------------------------------------------------------------

type rawFactor struct {
	Type              string                `json:"type" db:"type"`
	LastChallengedAt  time.Time             `json:"lastChallengedAt" db:"last_challenged_at"`
	ChallengedPayload JSON[json.RawMessage] `json:"challengedPayload" db:"challenged_payload"`
	LastVerifiedAt    time.Time             `json:"lastVerifiedAt" db:"last_verified_at"`
	VerifiedPayload   JSON[json.RawMessage] `json:"verifiedPayload" db:"verified_payload"`
}

func (f *rawFactor) ToDomain() (domain.SessionFactor, domain.SessionChallenge, error) {
	factorType, err := domain.SessionFactorTypeString(f.Type)
	if err != nil {
		return nil, nil, err
	}
	switch factorType {
	case domain.SessionFactorTypeUser:
		return f.userFactorToDomain()
	case domain.SessionFactorTypePassword:
		return f.passwordFactorToDomain()
	case domain.SessionFactorTypePasskey:
		return f.passkeyFactorToDomain()
	case domain.SessionFactorTypeIdentityProviderIntent:
		return f.identityProviderIntentFactorToDomain()
	case domain.SessionFactorTypeTOTP:
		return f.totpFactorToDomain()
	case domain.SessionFactorTypeOTPSMS:
		return f.otpSMSFactorToDomain()
	case domain.SessionFactorTypeOTPEmail:
		return f.otpEmailFactorToDomain()
	default:
		return nil, nil, nil
	}
}

func (f *rawFactor) userFactorToDomain() (factor domain.SessionFactor, _ domain.SessionChallenge, err error) {
	if f.LastVerifiedAt.IsZero() {
		return nil, nil, nil
	}
	factor = new(domain.SessionFactorUser)
	if err := json.Unmarshal(f.VerifiedPayload.Value, factor); err != nil {
		return nil, nil, err
	}
	return factor, nil, nil
}

func (f *rawFactor) passwordFactorToDomain() (domain.SessionFactor, domain.SessionChallenge, error) {
	if f.LastVerifiedAt.IsZero() {
		return nil, nil, nil
	}
	return &domain.SessionFactorPassword{
		LastVerifiedAt: f.LastVerifiedAt,
	}, nil, nil
}

func (f *rawFactor) passkeyFactorToDomain() (factor domain.SessionFactor, challenge domain.SessionChallenge, err error) {
	if !f.LastChallengedAt.IsZero() && f.LastChallengedAt.After(f.LastVerifiedAt) {
		challenge = new(domain.SessionChallengePasskey)
		if err := json.Unmarshal(f.ChallengedPayload.Value, &challenge); err != nil {
			return nil, nil, err
		}
	}
	if !f.LastVerifiedAt.IsZero() {
		factor = new(domain.SessionFactorPasskey)
		if err := json.Unmarshal(f.VerifiedPayload.Value, factor); err != nil {
			return nil, nil, err
		}
	}
	return factor, challenge, nil
}

func (f *rawFactor) identityProviderIntentFactorToDomain() (domain.SessionFactor, domain.SessionChallenge, error) {
	if f.LastVerifiedAt.IsZero() {
		return nil, nil, nil
	}
	return &domain.SessionFactorIdentityProviderIntent{
		LastVerifiedAt: f.LastVerifiedAt,
	}, nil, nil
}

func (f *rawFactor) totpFactorToDomain() (domain.SessionFactor, domain.SessionChallenge, error) {
	if f.LastVerifiedAt.IsZero() {
		return nil, nil, nil
	}
	return &domain.SessionFactorTOTP{
		LastVerifiedAt: f.LastVerifiedAt,
	}, nil, nil
}

func (f *rawFactor) otpSMSFactorToDomain() (factor domain.SessionFactor, challenge domain.SessionChallenge, err error) {
	if !f.LastChallengedAt.IsZero() && f.LastChallengedAt.After(f.LastVerifiedAt) {
		challenge = new(domain.SessionChallengeOTPSMS)
		if err := json.Unmarshal(f.ChallengedPayload.Value, &challenge); err != nil {
			return nil, nil, err
		}
	}
	if !f.LastVerifiedAt.IsZero() {
		factor = &domain.SessionFactorOTPSMS{
			LastVerifiedAt: f.LastVerifiedAt,
		}
	}
	return factor, challenge, nil
}

func (f *rawFactor) otpEmailFactorToDomain() (factor domain.SessionFactor, challenge domain.SessionChallenge, err error) {
	if !f.LastChallengedAt.IsZero() && f.LastChallengedAt.After(f.LastVerifiedAt) {
		challenge = new(domain.SessionChallengeOTPEmail)
		if err := json.Unmarshal(f.ChallengedPayload.Value, &challenge); err != nil {
			return nil, nil, err
		}
	}
	if !f.LastVerifiedAt.IsZero() {
		factor = &domain.SessionFactorOTPEmail{
			LastVerifiedAt: f.LastVerifiedAt,
		}
	}
	return factor, challenge, nil
}
