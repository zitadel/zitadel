package repository

import (
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

//TODO: below conditions need to be fixed

func (s sessionFactor) FactorTypeCondition(factorType domain.SessionFactorType) database.Condition {
	return database.NewNumberCondition(s.typeColumn(), database.NumberOperationEqual, factorType)
}

func (s sessionFactor) LastVerifiedBeforeCondition(lastVerifiedAt time.Time) database.Condition {
	return database.NewNumberCondition(s.lastVerifiedAtColumn(), database.NumberOperationLessThan, lastVerifiedAt)
}
